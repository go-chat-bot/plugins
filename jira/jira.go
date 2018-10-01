package jira

import (
	"bytes"
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/go-chat-bot/bot"
	"log"
	"os"
	"regexp"
	"text/template"
)

const (
	pattern         = ".*?([A-Z]+)-([0-9]+)\\b"
	userEnv         = "JIRA_USER"
	passEnv         = "JIRA_PASS"
	baseURLEnv      = "JIRA_BASE_URL"
	defaultTemplate = "{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): " +
		"{{.Fields.Summary}} - {{.Self}}"
)

var (
	url      string
	baseURL  string
	jiraUser string
	jiraPass string
	projects map[string]gojira.Project
	client   *gojira.Client
	re       = regexp.MustCompile(pattern)
)

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

func getIssues(text string) [][2]string {
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

func formatIssue(issueKey string, channel string) string {
	defaultRet := url + issueKey
	issue, _, err := client.Issue.Get(issueKey, nil)
	if err != nil {
		log.Printf("Failed getting issue %s info: %v\n", issueKey, err)
		return defaultRet
	}

	tmpl, err := template.New("default").Parse(defaultTemplate)
	if err != nil {
		log.Printf("Failed formatting for %s: %v\n", issueKey, err)
		return defaultRet
	}

	buf := &bytes.Buffer{}
	provideDefaultValues(issue)

	err = tmpl.Execute(buf, issue)
	if err != nil {
		log.Printf("Failed formatting for %s: %s\n", issueKey, err.Error())
		return defaultRet
	}
	return buf.String()
}

func jira(cmd *bot.PassiveCmd) (bot.CmdResultV3, error) {
	result := bot.CmdResultV3{
		Message: make(chan string),
		Done:    make(chan bool, 1)}
	result.Channel = cmd.Channel
	issues := getIssues(cmd.Raw)
	if issues != nil {
		go func() {
			for _, issue := range issues {
				key, num := issue[0], issue[1]
				_, found := projects[key]
				if found {
					result.Message <- formatIssue(key+"-"+num, cmd.Channel)
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

func init() {
	jiraUser = os.Getenv(userEnv)
	jiraPass = os.Getenv(passEnv)
	baseURL = os.Getenv(baseURLEnv)
	url = baseURL + "/browse/"

	err := initJIRAClient()
	if err != nil {
		log.Printf("Error querying JIRA for projects: %v\n", err)
		return
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
