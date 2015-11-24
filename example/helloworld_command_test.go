package example

import (
	"testing"

	"github.com/go-chat-bot/bot"
)

func TestHelloworld(t *testing.T) {
	bot := &bot.Cmd{
		Command: "helloworld",
		User: &bot.User{
			Nick:     "nick",
			RealName: "Real Name",
		},
	}
	want := "Hello Real Name"
	got, error := hello(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
