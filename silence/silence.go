package silence

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	"log"
	"strconv"
	"time"
)

var (
	silentMap map[string]time.Time
)

func silenceFilter(cmd *bot.FilterCmd) (string, error) {
	until, found := silentMap[cmd.Target]
	if !found || time.Now().After(until) {
		return cmd.Message, nil
	}
	log.Printf("Silencing message in %s\n", cmd.Target)
	return "", nil
}

func silence(cmd *bot.Cmd) (string, error) {
	if len(cmd.Args) != 1 {
		return "Argument must be exactly 1 number (of minutes to be silent)", nil
	}

	min, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Argument must be exactly 1 number (of minutes to be silent)!", nil
	}
	until := time.Now().UTC().Add(time.Duration(min) * time.Minute)

	go func() {
		// delay setting the timeout so plugin reply can go out first
		time.Sleep(1 * time.Second)
		silentMap[cmd.Channel] = until
	}()
	// disable map for plugin reply
	silentMap[cmd.Channel] = time.Now()
	return fmt.Sprintf("OK, I will be silent until %s\n",
		until.Format(time.RFC1123)), nil
}

func init() {
	silentMap = make(map[string]time.Time)
	bot.RegisterFilterCommand(
		"silence",
		silenceFilter)

	bot.RegisterCommand(
		"silence",
		"Makes the bot completely silent for X minutes (0 removes silence)",
		"5",
		silence)
}
