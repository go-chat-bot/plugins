package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chat-bot/bot"
)

const (
	pattern           = ".*?([A-Z]+)-([0-9]+)\\b"
	userEnv           = "JIRA_USER"
	passEnv           = "JIRA_PASS"
	tokenEnv          = "JIRA_TOKEN"
	baseURLEnv        = "JIRA_BASE_URL"
	channelConfigEnv  = "JIRA_CONFIG_FILE"
	notifyIntervalEnv = "JIRA_NOTIFY_INTERVAL"
	defaultTemplate   = "{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): " +
		"{{.Fields.Summary}} - {{.Self}}"
	defaultTemplateNew = "New {{.Fields.Type.Name}}: {{.Key}} " +
		"({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): " +
		"{{.Fields.Summary}} - {{.Self}}"
	defaultTemplateResolved = "Resolved {{.Fields.Type.Name}}: {{.Key}} " +
		"({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): " +
		"{{.Fields.Summary}} - {{.Self}}"
	verboseEnv = "JIRA_VERBOSE"
	threadEnv  = "JIRA_THREAD"
)

var (
	url              string
	projects         map[string]gojira.Project // project.Key -> project map
	channelConfigs   map[string]channelConfig  // channel -> channelConfig map
	notifyNewConfig  map[string][]string       // project.Key -> slice of channel names
	notifyResConfig  map[string][]string       // project.Key -> slice of channel names
	componentsConfig map[string][]string       // project.key -> slice of component names
	client           *gojira.Client
	re               = regexp.MustCompile(pattern)
	projectJQL       = "project in (%s) "
	componentJQL     = "AND component in (%s) "
	newJQL           = "AND resolution = Unresolved " +
		"AND created > '-%dm' " +
		"ORDER BY key ASC"
	resolvedJQL = "AND resolved > '-%dm' " +
		"ORDER BY key ASC"
	notifyInterval int
	verbose        bool
	thread         bool
)

type channelConfig struct {
	Channel          string   `json:"channel"`
	Thread           string   `json:"thread,omitempty"`
	Template         string   `json:"template,omitempty"`         // template format for issues being posted
	TemplateNew      string   `json:"templateNew,omitempty"`      // template format for newly created issues
	TemplateResolved string   `json:"templateResolved,omitempty"` // template format for resolved issues
	NotifyNew        []string `json:"notifyNew,omitempty"`        // list of JIRA projects to watch for new issues
	NotifyResolved   []string `json:"notifyResolved,omitempty"`   // list of JIRA projects to watch for resolved issues
	Components       []string `json:"components,omitempty"`       // list of JIRA project components to watch for
}

func getProjects() (map[string]gojira.Project, error) {
	req, err := client.NewRequest("GET", "rest/api/2/project", nil)
	if err != nil {
		return projects, fmt.Errorf("Error creating request object: %v", err)
	}

	projectObjects := new([]gojira.Project)
	projects = make(map[string]gojira.Project)
	_, err = client.Do(req, projectObjects)
	if err != nil {
		return projects, fmt.Errorf("Failed getting JIRA projects: %v", err)
	}
	for _, project := range *projectObjects {
		projects[project.Key] = project
	}
	return projects, nil
}

func getIssuesFromString(text string) [][2]string {
	matches := re.FindAllStringSubmatch(text, -1)
	var data [][2]string
	for _, match := range matches {
		// match[1] == project key
		// match[2] == issue number
		data = append(data, [2]string{match[1], match[2]})
	}
	return data
}

func provideDefaultValues(issue *gojira.Issue) {
	if issue.Fields.Assignee == nil {
		issue.Fields.Assignee = &gojira.User{Key: "no assignee"}
	}
	// we use Self as the web URL in template
	issue.Self = url + issue.Key
}

func formatIssue(issue *gojira.Issue, channel string, templ string) string {
	defaultRet := url + issue.Key
	provideDefaultValues(issue)

	tmpl, err := template.New("default").Parse(templ)
	if err != nil {
		log.Printf("Failed formatting for %s: %v\n", issue.Key, err)
		return defaultRet
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, issue)
	if err != nil {
		log.Printf("Failed formatting for %s: %s\n", issue.Key, err.Error())
		return defaultRet
	}
	return buf.String()
}

func jira(cmd *bot.PassiveCmd) (bot.CmdResultV3, error) {
	result := bot.CmdResultV3{
		Message: make(chan string),
		Done:    make(chan bool, 1)}
	result.Channel = cmd.Channel
	issues := getIssuesFromString(cmd.Raw)
	if issues != nil {
		go func() {
			for _, issue := range issues {
				project, num := issue[0], issue[1]
				key := project + "-" + num
				_, found := projects[project]
				if found {
					issue, _, err := client.Issue.Get(key, nil)
					if err != nil {
						log.Printf("Failed getting issue %s info: %v\n",
							key, err)
						continue
					}
					if verbose {
						log.Printf("Replying to %s about issue %s\n", cmd.Channel,
							key)
					}
					template := defaultTemplate
					config, found := channelConfigs[cmd.Channel]
					if found {
						template = config.Template
					}
					result.Message <- formatIssue(issue, cmd.Channel, template)
				}
			}
			result.Done <- true
		}()
	} else {
		result.Done <- true
	}

	return result, nil
}

func containsComponent(fromJira []*gojira.Component, fromConf []string) bool {
	for _, i := range fromConf {
		for _, j := range fromJira {
			if i == j.Name {
				return true
			}
		}
	}
	return false
}

func periodicJIRANotifyNew() (ret []bot.CmdResult, err error) {
	newProjectKeys := make([]string, 0, len(notifyNewConfig))
	for k := range notifyNewConfig {
		newProjectKeys = append(newProjectKeys, k)
	}
	componentsKeys := make([]string, 0, len(componentsConfig))
	for k := range componentsConfig {
		componentsKeys = append(componentsKeys, k)
	}

	query := fmt.Sprintf(projectJQL, strings.Join(newProjectKeys, ","))
	query = query + fmt.Sprintf(newJQL, notifyInterval)
	if verbose {
		log.Printf("New issues query: %s", query)
	}
	newIssues, _, err := client.Issue.Search(query, nil)
	if err != nil {
		log.Printf("Error querying JIRA for new issues: %v\n", err)
		return nil, err
	}
	for _, issue := range newIssues {
		channels := notifyNewConfig[issue.Fields.Project.Key]
		for _, notifyChan := range channels {
			if len(channelConfigs[notifyChan].Components) == 0 ||
				(len(channelConfigs[notifyChan].Components) > 0 &&
					containsComponent(issue.Fields.Components, channelConfigs[notifyChan].Components)) {
				// displays only if Components are not defined OR Components exist in Jira output
				threadName := channelConfigs[notifyChan].Thread
				if thread && (len(threadName) > 0) {
					notifyChan += ":" + notifyChan + "/" + threadName
				}
				if verbose {
					log.Printf("Notifying %s about new %s %s", notifyChan,
						issue.Fields.Type.Name,
						issue.Key)
				}
				template := defaultTemplateNew
				config, found := channelConfigs[notifyChan]
				if found {
					template = config.TemplateNew
				}
				ret = append(ret, bot.CmdResult{
					Message: formatIssue(&issue, notifyChan, template),
					Channel: notifyChan,
				})
			}
		}
	}

	return ret, nil
}

func periodicJIRANotifyResolved() (ret []bot.CmdResult, err error) {
	resolvedProjectKeys := make([]string, 0, len(notifyResConfig))
	for k := range notifyResConfig {
		resolvedProjectKeys = append(resolvedProjectKeys, k)
	}
	componentsKeys := make([]string, 0, len(componentsConfig))
	for k := range componentsConfig {
		componentsKeys = append(componentsKeys, k)
	}

	query := fmt.Sprintf(projectJQL, strings.Join(resolvedProjectKeys, ","))
	query = query + fmt.Sprintf(resolvedJQL, notifyInterval)
	if verbose {
		log.Printf("Resolved issues query: %s", query)
	}
	resolvedIssues, _, err := client.Issue.Search(query, nil)
	if err != nil {
		log.Printf("Error querying JIRA for resolved issues: %v\n", err)
		return nil, err
	}
	for _, issue := range resolvedIssues {
		channels := notifyResConfig[issue.Fields.Project.Key]
		if verbose {
			log.Printf("Resolved issues result: %s", spew.Sdump(issue.Fields.Components))
		}
		for _, notifyChan := range channels {
			if len(channelConfigs[notifyChan].Components) == 0 ||
				(len(channelConfigs[notifyChan].Components) > 0 &&
					containsComponent(issue.Fields.Components, channelConfigs[notifyChan].Components)) {
				// displays only if Components are not defined OR Components exist in Jira output
				threadName := channelConfigs[notifyChan].Thread
				if thread && (len(threadName) > 0) {
					notifyChan += ":" + notifyChan + "/" + threadName
				}
				if verbose {
					log.Printf("Notifying %s about resolved %s %s", notifyChan,
						issue.Fields.Type.Name,
						issue.Key)
				}
				template := defaultTemplateResolved
				config, found := channelConfigs[notifyChan]
				if found {
					template = config.TemplateResolved
				}
				ret = append(ret, bot.CmdResult{
					Message: formatIssue(&issue, notifyChan, template),
					Channel: notifyChan,
				})
			}
		}
	}

	return ret, nil
}

func initJIRAClient(baseURL, jiraUser, jiraPass, jiraToken string) error {
	var err error

	if len(jiraToken) > 0 {
		tpPATA := gojira.PATAuthTransport{
			Token: jiraToken,
		}
		client, err = gojira.NewClient(tpPATA.Client(), baseURL)
	} else {
		tpBA := gojira.BasicAuthTransport{
			Username: jiraUser,
			Password: jiraPass,
		}
		client, err = gojira.NewClient(tpBA.Client(), baseURL)
	}
	if err != nil {
		log.Printf("Error initializing JIRA client: %v\n", err)
		return err
	}
	return nil
}

func loadChannelConfigs(filename string) error {
	channelConfigs = make(map[string]channelConfig)
	notifyNewConfig = make(map[string][]string)
	notifyResConfig = make(map[string][]string)
	componentsConfig = make(map[string][]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed opening configuration file %s: %v\n", filename, err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configs := make([]channelConfig, 0)
	err = decoder.Decode(&configs)
	if err != nil {
		log.Printf("Error loading configuration: %v\n", err)
		return err
	}
	for _, chanConf := range configs {
		if chanConf.Channel == "" {
			log.Println("Configuration without channel found. Skipping")
			continue
		}
		if chanConf.Template == "" {
			chanConf.Template = defaultTemplate
		}
		if chanConf.TemplateNew == "" {
			chanConf.TemplateNew = defaultTemplateNew
		}
		if chanConf.TemplateResolved == "" {
			chanConf.TemplateResolved = defaultTemplateResolved
		}
		channelConfigs[chanConf.Channel] = chanConf
		for _, project := range chanConf.NotifyNew {
			notifyNewConfig[project] = append(notifyNewConfig[project],
				chanConf.Channel)
		}
		for _, project := range chanConf.NotifyResolved {
			notifyResConfig[project] = append(notifyResConfig[project],
				chanConf.Channel)
		}
		for _, project := range chanConf.Components {
			componentsConfig[project] = append(componentsConfig[project],
				chanConf.Channel)
		}
	}
	return nil
}

func init() {
	_, verbose = os.LookupEnv(verboseEnv)
	_, thread = os.LookupEnv(threadEnv)

	jiraUser := os.Getenv(userEnv)
	jiraPass := os.Getenv(passEnv)
	jiraToken := os.Getenv(tokenEnv)
	baseURL := os.Getenv(baseURLEnv)
	confFile := os.Getenv(channelConfigEnv)
	url = baseURL + "/browse/"

	err := initJIRAClient(baseURL, jiraUser, jiraPass, jiraToken)
	if err != nil {
		log.Printf("Error querying JIRA for projects: %v\n", err)
		return
	}

	if confFile != "" {
		err = loadChannelConfigs(confFile)
		if err != nil {
			log.Printf("Error loading channel configuration (non-fatal): %v\n", err)
		}
	}

	_, err = getProjects()
	if err != nil {
		log.Printf("Error querying JIRA for projects: %v\n", err)
		return
	}

	interval := os.Getenv(notifyIntervalEnv)
	if interval == "" {
		interval = "1"
	}
	notifyInterval, err = strconv.Atoi(interval)
	if err != nil {
		log.Printf("Error parsing interval from %s. Using default",
			interval)
		notifyInterval = 1
	}

	bot.RegisterPassiveCommandV2(
		"jira",
		jira)

	if len(notifyNewConfig) > 0 {
		bot.RegisterPeriodicCommandV2(
			"periodicJIRANotifyNew",
			bot.PeriodicConfig{
				CronSpec:  fmt.Sprintf("*/%d * * * *", notifyInterval),
				CmdFuncV2: periodicJIRANotifyNew,
			})
	}
	log.Printf("New issue notifications set up for %d JIRA projects", len(notifyNewConfig))
	if len(notifyResConfig) > 0 {
		bot.RegisterPeriodicCommandV2(
			"periodicJIRANotifyResolved",
			bot.PeriodicConfig{
				CronSpec:  fmt.Sprintf("*/%d * * * *", notifyInterval),
				CmdFuncV2: periodicJIRANotifyResolved,
			})
	}
	log.Printf("Resolved issue notifications set up for %d JIRA projects", len(notifyResConfig))
	log.Printf("JIRA plugin initialization successful")
}
