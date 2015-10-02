package crypto

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	invalidParams         = "Invalid parameters"
)

func crypto(command *bot.Cmd) (string, error) {

	if len(command.Args) < 2 {
		return invalidAmountOfParams, nil
	}

	var str string
	var err error
	switch command.Args[0] {
	case "md5":
		s := strings.Join(command.Args[1:], " ")
		str, err = encryptMD5(s)
	default:
		return invalidParams, nil
	}

	if err != nil {
		return fmt.Sprintf("Error: %s", err), nil
	}

	return str, nil
}

func encryptMD5(str string) (string, error) {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func init() {
	bot.RegisterCommand(
		"crypto",
		"Encrypts the input data from its hash value",
		"md5 enter here text to encrypt",
		crypto)
}
