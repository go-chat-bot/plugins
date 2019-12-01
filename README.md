[![Circle CI](https://circleci.com/gh/go-chat-bot/plugins.svg?style=svg)](https://circleci.com/gh/go-chat-bot/plugins)

### Active

* **gif**: Posts a random gif url from [giphy.com][giphy.com]. Try it with: **!gif cat**
* **catgif**: Posts a random cat gif url from [thecatapi.com][thecatapi.com]
* **godoc**: Searches packages in godoc.org. Try it with: **!godoc net/http**
* **puppet**: Allows you to send messages through the bot: Try it with: **!puppet say #go-bot Hello!**
* **guid**: Generates a new guid
* **crypto**: Encrypts the input data using sha1 or md5
* **encode**: Encodes a string, currently only to base64
* **decode**: Decodes a string. currently ony from base64
* **treta**: Use it to sow discord on a channel

Tip: Use `!help <command>` to obtaing more info about these commands.

### Passive (triggers)

Passive commands receive all the text sent to the bot or the channels that the bot is in and can process it and reply.

These commands differ from the active commands as they are executed for every text that the bot receives. Ex: The Chuck Norris command, replies with a Chuck Norris fact every time the words "chuck" or "norris" are mentioned on a channel.

* **url**: Detects url and gets it's title (very naive implementation, works sometimes)
* **catfacts**: Tells a random cat fact based on some cat keywords
* **jira**: Detects jira issue numbers and posts information about it. Necessary
  to configure. See README.md in jira subdirectory for details
* **chucknorris**: Shows a random chuck norris quote every time the word "chuck" is mentioned
* **bitly**: Shortens URLs appearing in output of other plugins before they are sent to channels

### Periodic (triggers)

Periodic commands are run based on a [cron
specification](https://godoc.org/github.com/robfig/cron) passed to the
config. These commands are runned periodically, outputting a message to
the configured channel(s).

Look into the good morning [example
command](https://github.com/go-chat-bot/plugins/blob/master/example/goodmorning_command.go) for guidance on how to write and configure periodic commands.

* **cachet**: Notifies of service outages based on Cachet data

### Wish to write a new plugin?

Start with the example commands in the [example directory](https://github.com/go-chat-bot/plugins/tree/master/example).

It's dead simple, you just need to write a go function and register it on the bot.

Here's a Hello World plugin example:

```Go
package example

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

func hello(command *bot.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("Hello %s", command.User.RealName)
	return
}

func init() {
	bot.RegisterCommand(
		"hello",
		"Sends a 'Hello' message to you on the channel.",
		"",
		hello)
}
```
