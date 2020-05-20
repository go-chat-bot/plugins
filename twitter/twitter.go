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

// newAuthenticatedTwitterClient uses a provided consumer key and secret to authenticate
// against Twitter's Oauth2 endpoint, then validates the authentication by checking the
// current RateLimit against the provided account credentials.
// It returns a twitter.Client.
func newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret string) (*twitter.Client, error) {
	if twitterConsumerKey == "" || twitterConsumerSecret == "" {
		return nil, errors.New("missing API credentials")
	}
	// oauth2
	config := &clientcredentials.Config{
		ClientID:     twitterConsumerKey,
		ClientSecret: twitterConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	params := twitter.RateLimitParams{Resources: []string{"statuses"}}

	// query rate limit, to verify authentication and...make sure we haven't been rate limited :)
	// NOTE: calls to RateLimits apply against the Remaining calls for that endpoint
	rl, resp, err := client.RateLimits.Status(&params)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode/200 != 1 {
		err = errors.New(resp.Status)
	}

	remaining := rl.Resources.Statuses["/statuses/show/:id"].Remaining

	if remaining == 0 {
		err = errors.New("rate limit exceeded")
	}

	return client, err
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
	// populate FullText field
	params := twitter.StatusShowParams{TweetMode: "extended"}
	// FIXME what if we get an error here...
	tweet, resp, err := client.Statuses.Show(tweetID, &params)
	if resp.StatusCode/200 != 1 {
		// ...and also here?
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
		username := tweet.User.ScreenName
		text := newlines.ReplaceAllString(tweet.FullText, " ")
		newMessage := fmt.Sprintf(formatString, username, text)
		messages = append(messages, newMessage)
	}
	return messages
}

// expandTweet receives a bot.PassiveCmd and performs the full parse-and-fetch
// pipeline. It sets up a client, finds Tweet IDs in the message text, fetches
// the tweets, and formats them. If multiple Tweet IDs were found in the message,
// all formatted Tweets will be joined into a single message.
// It returns a single string suitable for sending as a chat message.
func expandTweet(cmd *bot.PassiveCmd) (string, error) {
	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	var message string

	messageText := cmd.MessageData.Text
	tweetIDs, err := findTweetIDs(messageText)
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
		expandTweet)
}
