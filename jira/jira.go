package jira

import (
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/go-chat-bot/bot"
	"log"
	"os"
	"regexp"
)

const (
	pattern    = ".*?([A-Z]+)-([0-9]+)\\b"
	env        = "JIRA_ISSUES_URL"
	userEnv    = "JIRA_USER"
	passEnv    = "JIRA_PASS"
	baseURLEnv = "JIRA_BASE_URL"
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
					result.Message <- url + key + "-" + num
				}
			}
			result.Done <- true
		}()
	} else {
		result.Done <- true
	}

	return result, nil
}

func init() {
	var err error
	url = os.Getenv(env)
	jiraUser = os.Getenv(userEnv)
	jiraPass = os.Getenv(passEnv)
	baseURL = os.Getenv(baseURLEnv)

	tp := gojira.BasicAuthTransport{
		Username: jiraUser,
		Password: jiraPass,
	}

	client, err = gojira.NewClient(tp.Client(), baseURL)
	if err != nil {
		log.Printf("Error initializing JIRA client: %v\n", err)
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
