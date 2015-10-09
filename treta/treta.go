package treta

import (
	"math/rand"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	msgInvalidAmountOfParams = "Invalid amount of parameters"
	msgInvalidParam          = "Invalid parameter"
)

var (
	quotes = map[string][]string{
		"DELPHI": {
			"Delphi. Now there's a name I haven't heard in a long time.",
			"Access Violation at address 00405772 in module 'Project1.exe'. Read of address 00000388.",
			"Develop iOS applications with RAD Studio",
			"It’s not difficult to read and listen about the wonders of Embarcadero DataSnap technology around the world.",
		},
		"JAVA": {
			"You're using Java? Well there's your problem.",
			"I had a problem so I thought to use Java. Now I have a ProblemFactory.",
		},
		"JAVASCRIPT": {
			"Javascript is not funny",
			"JavaScript why you no works?",
			"Brace yourself. A new Javascript framework is coming.",
		},
		"PYTHON": {
			"We'll can do cool things... even with Python",
		},
		"RUBY": {
			"Ruby is slower than Internet Explorer",
			"Can Rails Scale? NOOOOO!",
			"Why is Ruby so slow?",
			"I hate managing inventory and the game drops more weapon than the rails can handle the requests",
			"Ruby on Rails? Pleaaase. Do you even code, bro?",
			"The classic Hello, world! program is really easy with Ruby. You just need to know the name of the gem you want to install.",
			"Python > Ruby",
			"even PHP > Ruby",
		},
		"VIM": {
			"Emacs > VIM",
			"Sublime Text > VIM",
			"even Notepad > VIM",
			"VIM... Why can't I quit you?!",
			"Vim Is Too Mainstream. I'm Switching To Emacs",
		},
		"WINDOWS": {
			"If We Add A Start Menu To Windows 8 We Can Call It Windows 10",
			"Keyboard not responding. Press any key to continue.",
			"A system call that should never fail has failed.",
			"Bluescreen has performed an illegal operation. Bluescreen must be closed.",
			"An error occurred whilst trying to load the previous error.",
			"Help and Support Error: Windows cannot open Help and Support because a system service is not running. To fix this problems, start the service named Help and Support",
		},
	}
)

func treta(command *bot.Cmd) (string, error) {
	var key string
	switch len(command.Args) {
	case 0:
		key = randKey()
	case 1:
		key = strings.ToUpper(command.Args[0])
	default:
		return msgInvalidAmountOfParams, nil
	}

	q, found := quotes[key]
	if !found {
		return msgInvalidParam, nil
	}
	return q[rand.Intn(len(q))], nil
}

func randKey() string {
	keys := make([]string, 0, len(quotes))
	for k := range quotes {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

func init() {
	bot.RegisterCommand(
		"treta",
		"sowing discord",
		"",
		treta)
}
