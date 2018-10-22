package bitly

import (
	"github.com/go-chat-bot/bot"
	"github.com/mvdan/xurls"
	"github.com/zpnk/go-bitly"
	"log"
	"os"
	"strings"
)

const (
	bitlyTokenEnv = "BITLY_TOKEN"
)

var (
	bitlyClient *bitly.Client
	urlRegex    = xurls.Strict()
)

func bitlyFilter(cmd *bot.FilterCmd) (string, error) {
	urls := urlRegex.FindAllString(cmd.Message, -1)
	if urls == nil {
		// no urls to shorten
		return cmd.Message, nil
	}

	for _, url := range urls {
		shortURL, err := bitlyClient.Links.Shorten(url)
		if err != nil {
			log.Printf("Failed to shorten URL (%s): %s", url, err.Error())
			continue
		}
		cmd.Message = strings.Replace(cmd.Message,
			url, shortURL.URL, -1)
	}

	return cmd.Message, nil
}

func init() {
	token := os.Getenv(bitlyTokenEnv)
	bitlyClient = bitly.New(token)

	bot.RegisterFilterCommand(
		"bitly",
		bitlyFilter)
}
