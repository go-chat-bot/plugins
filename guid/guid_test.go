package guid

import (
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
)

const (
	guidSize = 36
)

func TestCPF(t *testing.T) {
	bot := &bot.Cmd{
		Command: "guid",
	}

	got, error := guid(bot)

	if len(got) != guidSize {
		t.Errorf("Expected GUID with '%v' characters got '%v'", guidSize, len(got))
	}

	if strings.Split(got, "")[14] != "4" {
		t.Errorf("Expected GUID version 4 got an invalid ('%v')", got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
