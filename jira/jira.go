package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/go-chat-bot/bot"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	pattern           = ".*?([A-Z]+)-([0-9]+)\\b"
	userEnv           = "JIRA_USER"
	passEnv           = "JIRA_PASS"
	baseURLEnv        = "JIRA_BASE_URL"
	channelConfigEnv  = "JIRA_CONFIG_FILE"
	notifyIntervalEnv = "JIRA_NOTIFY_INTERVAL"
	defaultTemplate   = "{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): " +
		"{{.Fields.Summary}} - {{.Self}}"
)

var (
	url             string
	baseURL         string
	jiraUser        string
	jiraPass        string
	projects        map[string]gojira.Project
	channelConfigs  map[string]channelConfig
	notifyNewConfig map[string][]string
	notifyResConfig map[string][]string
	client          *gojira.Client
	re              = regexp.MustCompile(pattern)
)

type channelConfig struct {
	Channel        string   `json:"channel"`
	Template       string   `json:"template,omitempty"`       // template format for issues being posted
	NotifyNew      []string `json:"notifyNew,omitempty"`      // list of JIRA projects to watch for new issues
	NotifyResolved []string `json:"notifyResolved,omitempty"` // list of JIRA projects to watch for resolved issues
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

func formatIssue(issue *gojira.Issue, channel string) string {
	defaultRet := url + issue.Key
	provideDefaultValues(issue)
	templ := defaultTemplate
	config, found := channelConfigs[channel]
	if found {
		templ = config.Template
	}

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
					log.Printf("Replying to %s about issue %s\n", cmd.Channel,
						key)
					result.Message <- formatIssue(issue, cmd.Channel)
				}
			}
			result.Done <- true
		}()
	} else {
		result.Done <- true
	}

	return result, nil
}

func initJIRAClient() error {
	var err error

	tp := gojira.BasicAuthTransport{
		Username: jiraUser,
		Password: jiraPass,
	}

	client, err = gojira.NewClient(tp.Client(), baseURL)
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
		channelConfigs[chanConf.Channel] = chanConf
		for _, project := range chanConf.NotifyNew {
			notifyNewConfig[project] = append(notifyNewConfig[project],
				chanConf.Channel)
		}
		for _, project := range chanConf.NotifyResolved {
			notifyResConfig[project] = append(notifyResConfig[project],
				chanConf.Channel)
		}
	}
	return nil
}

func init() {
	jiraUser = os.Getenv(userEnv)
	jiraPass = os.Getenv(passEnv)
	baseURL = os.Getenv(baseURLEnv)
	confFile := os.Getenv(channelConfigEnv)
	url = baseURL + "/browse/"

	err := initJIRAClient()
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

	bot.RegisterPassiveCommandV2(
		"jira",
		jira)
}
