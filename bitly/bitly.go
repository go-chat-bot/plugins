package bitly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chat-bot/bot"
	"io/ioutil"
	"log"
	"mvdan.cc/xurls/v2"
	"net/http"
	"os"
	"strings"
)

type shortenRequest struct {
	LongURL string `json:"long_url"`
}

type shortenReply struct {
	Link string `json:"link"`
}

const (
	bitlyTokenEnv = "BITLY_TOKEN"
	shortenURLAPI = "https://api-ssl.bitly.com/v4/shorten"
)

var (
	urlRegex = xurls.Strict()
)

func shorten(longurl string) (string, error) {
	sr := shortenRequest{longurl}
	body, err := json.Marshal(sr)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", shortenURLAPI, bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s",
		os.Getenv(bitlyTokenEnv)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("bitly API request returned non-20x code: %d", resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body. ", err)
	}
	shortReply := shortenReply{}
	err = json.Unmarshal(body, &shortReply)
	if err != nil {
		return "", err
	}
	return shortReply.Link, nil
}

func bitlyFilter(cmd *bot.FilterCmd) (string, error) {
	urls := urlRegex.FindAllString(cmd.Message, -1)
	if urls == nil {
		// no urls to shorten
		return cmd.Message, nil
	}

	for _, url := range urls {
		shortURL, err := shorten(url)
		if err != nil {
			log.Printf("Failed to shorten URL (%s): %s", url, err.Error())
			continue
		}
		log.Printf("Succesfully shortened URL (%s) to %s", url, shortURL)
		cmd.Message = strings.Replace(cmd.Message,
			url, shortURL, -1)
	}

	return cmd.Message, nil
}

func init() {
	bot.RegisterFilterCommand(
		"bitly",
		bitlyFilter)
}
