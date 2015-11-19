package example

import (
	"testing"

	"github.com/go-chat-bot/bot"
)

func TestHelloworld(t *testing.T) {
	want := "Hello nick"
	bot := &bot.Cmd{
		Command: "helloworld",
		User: &bot.User{
			Nick: "nick",
		},
	}
	got, error := hello(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
