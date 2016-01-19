package gag

import (
	"net/http"

	"github.com/go-chat-bot/bot"
)

const (
	rangomPage = "http://9gag.com/random"
)

func gag(command *bot.Cmd) (string, error) {
	res, err := http.Get(rangomPage)
	if err != nil {
		return "", err
	}
	return res.Request.URL.String(), nil
}

func init() {
	bot.RegisterCommand(
		"9gag",
		"Returns a random 9gag page.",
		"",
		gag)
}
