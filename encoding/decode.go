package encoding

import (
	"encoding/base64"
	"fmt"

	"github.com/go-chat-bot/bot"
)

func decode(command *bot.Cmd) (string, error) {

	if len(command.Args) != 2 {
		return invalidAmountOfParams, nil
	}

	var str string
	var err error
	switch command.Args[0] {
	case "base64":
		str, err = decodeBase64(command.Args[1])
	default:
		return invalidParam, nil
	}

	if err != nil {
		return fmt.Sprintf("Error: %s", err), nil
	}

	return str, nil
}

func decodeBase64(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%s", data), nil
}

func init() {
	bot.RegisterCommand(
		"decode",
		"Decodes the given string",
		"base64 VGhlIEdvIFByb2dyYW1taW5nIExhbmd1YWdl",
		decode)
}
