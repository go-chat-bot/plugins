package guid

import (
	uuid "github.com/beevik/guid"
	"github.com/go-chat-bot/bot"
)

func guid(command *bot.Cmd) (string, error) {
	return uuid.NewString(), nil
}

func init() {
	bot.RegisterCommand(
		"guid",
		"Generates UUID",
		"",
		guid)
}
