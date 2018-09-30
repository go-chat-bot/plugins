package jira

import (
	"github.com/go-chat-bot/bot"
	"os"
	"regexp"
)

const (
	pattern = ".*?([A-Z]+)-([0-9]+)\\b"
	env     = "JIRA_ISSUES_URL"
)

var (
	url string
	re  = regexp.MustCompile(pattern)
)

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
				result.Message <- url + key + "-" + num
			}
			result.Done <- true
		}()
	} else {
		result.Done <- true
	}

	return result, nil
}

func init() {
	url = os.Getenv(env)
	bot.RegisterPassiveCommandV2(
		"jira",
		jira)
}
