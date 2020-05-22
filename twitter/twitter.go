// Package twitter provides a plugin that scrapes messages for Twitter links,
// then expands them into chat messages.
package twitter

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/go-chat-bot/bot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// findTweetIDs checks a given message string for strings that look like Twitter links,
// then attempts to extract the Tweet ID from the link.
// It returns an array of Tweet IDs.
func findTweetIDs(message string) ([]int64, error) {
	re := regexp.MustCompile(`http(?:s)?://(?:mobile.)?twitter.com/(?:.*)/status/([0-9]*)`)
	// FIXME this is only returning the LAST match, should return ALL matches
	result := re.FindAllStringSubmatch(message, -1)
	var (
		tweetIDs []int64
		id       int64
		err      error
	)

	for i := range result {
		last := len(result[i]) - 1
		idStr := result[i][last]
		id, err = strconv.ParseInt(idStr, 10, 64)
		tweetIDs = append(tweetIDs, id)
	}
	return tweetIDs, err
}

// getCredentialsFromEnvironment attempts to extract the Twitter consumer key
// and consumer secret from the current process environment. If either the key
// or the secret is not found, it returns a pair of empty strings and a
// missingAPICredentialsError.
// If successful, it returns the consumer key and consumer secret.
func getCredentialsFromEnvironment() (string, string, error) {
	var err error
	key, keyOk := os.LookupEnv("TWITTER_CONSUMER_KEY")
	secret, secretOk := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if !keyOk || !secretOk {
		return "", "", errors.New("missing API credentials")
	}
	return key, secret, err
}

// newTwitterClientConfig takes a Twitter consumer key and consumer secret and
// attempts to create a clientcredentials.Config. If either the key or the secret
// is an empty string, no client is returned and a missingAPICredentialsError is returned.
// If successful, it returns a clientcredentials.Config.
func newTwitterClientConfig(twitterConsumerKey, twitterConsumerSecret string) (*clientcredentials.Config, error) {
	if twitterConsumerKey == "" || twitterConsumerSecret == "" {
		return nil, errors.New("missing API credentials")
	} else {
		config := &clientcredentials.Config{
			ClientID:     twitterConsumerKey,
			ClientSecret: twitterConsumerSecret,
			TokenURL:     "https://api.twitter.com/oauth2/token",
		}
		return config, nil
	}
}

// newAuthenticatedTwitterClient uses a provided consumer key and secret to authenticate
// against Twitter's Oauth2 endpoint, then validates the authentication by checking the
// current RateLimit against the provided account credentials.
// It returns a twitter.Client.
func newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret string) (*twitter.Client, error) {
	config, err := newTwitterClientConfig(twitterConsumerKey, twitterConsumerSecret)
	if err != nil {
		return nil, err
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	err = checkTwitterClientRateLimit(client)

	return client, err
}

// checkTwitterClientRateLimit uses the provided twitter.Client to check the remaining
// RateLimit.Status for that client.
// It returns an error if authentication failed or if the rate limit has been exceeded.
func checkTwitterClientRateLimit(client *twitter.Client) error {
	// NOTE: calls to RateLimits apply against the Remaining calls for that endpoint
	params := twitter.RateLimitParams{Resources: []string{"statuses"}}
	rl, resp, err := client.RateLimits.Status(&params)

	// FIXME if i don't return this err at this point and credentials are bad, a panic happens
	if err != nil {
		return err
	}

	remaining := rl.Resources.Statuses["/statuses/show/:id"].Remaining
	if resp.StatusCode/200 != 1 {
		return errors.New(resp.Status)
	}

	if remaining == 0 {
		return errors.New("rate limit exceeded")
	}
	return err
}

// fetchTweets takes an array of Tweet IDs and retrieves the corresponding
// Statuses.
// It returns an array of twitter.Tweets.
func fetchTweets(client *twitter.Client, tweetIDs []int64) ([]twitter.Tweet, error) {
	var tweets []twitter.Tweet
	var err error
	for _, tweetID := range tweetIDs {
		tweet, fetchErr := fetchTweet(client, tweetID)
		if fetchErr != nil {
			// TODO what about multiple rrors
			err = fetchErr
		}
		tweets = append(tweets, *tweet)
	}
	return tweets, err
}

// fetchTweet takes a single Tweet ID and fetches the corresponding Status.
// It returns a twitter.Tweet.
func fetchTweet(client *twitter.Client, tweetID int64) (*twitter.Tweet, error) {
	var err error
	// TODO get alt text
	// params: include_entities=true,include_ext_alt_text=true

	// populate FullText field
	params := twitter.StatusShowParams{TweetMode: "extended"}
	tweet, resp, err := client.Statuses.Show(tweetID, &params)

	if err != nil {
		return tweet, err
	}

	if resp.StatusCode/200 != 1 {
		err = errors.New(resp.Status)
	}

	return tweet, err
}

// formatTweets takes an array of twitter.Tweets and formats them in preparation for
// sending as a chat message.
// It returns an array of nicely formatted strings.
func formatTweets(tweets []twitter.Tweet) []string {
	formatString := "Tweet from @%s: %s"
	newlines := regexp.MustCompile(`\r?\n`)
	var messages []string
	for _, tweet := range tweets {
		// TODO get link title, eg: Tweet from @user: look at this cool thing https://thing.cool (Link title: A Cool Thing)
		// tweet.Entities.Urls contains []URLEntity
		// fetch title from urlEntity.URL
		// urls plugin already correctly handles t.co links
		username := tweet.User.ScreenName
		text := newlines.ReplaceAllString(tweet.FullText, " ")
		newMessage := fmt.Sprintf(formatString, username, text)
		messages = append(messages, newMessage)
	}
	return messages
}

// expandTweets receives a bot.PassiveCmd and performs the full parse-and-fetch
// pipeline. It sets up a client, finds Tweet IDs in the message text, fetches
// the tweets, and formats them. If multiple Tweet IDs were found in the message,
// all formatted Tweets will be joined into a single message.
// It returns a single string suitable for sending as a chat message.
func expandTweets(cmd *bot.PassiveCmd) (string, error) {
	var message string
	messageText := cmd.MessageData.Text
	// message text could be empty

	twitterConsumerKey, twitterConsumerSecret, err := getCredentialsFromEnvironment()
	// key or secret could be empty
	if err != nil {
		return message, err
	}

	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	// credentials could be non-empty but bad
	if err != nil {
		return message, err
	}

	tweetIDs, err := findTweetIDs(messageText)
	if err != nil {
		return message, err
	}

	tweets, err := fetchTweets(client, tweetIDs)
	if err != nil {
		return message, err
	}

	formattedTweets := formatTweets(tweets)
	if formattedTweets != nil {
		message = strings.Join(formattedTweets, "\n")
	}
	return message, err
}

// init initalizes a PassiveCommand for expanding Tweets.
func init() {
	// TODO initialize Twitter client here
	// we should only need to create a Twitter.Client once on startup
	bot.RegisterPassiveCommand(
		"twitter",
		expandTweets)

	// TODO !ratelimit command to expose current rate limits
}
