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
)

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

// actually this should take an interface :o
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
	formattedTweets := formatTweets(tweets)
	if formattedTweets != nil {
		// FIXME only 1 tweet per message lol sry
		// join with newlines maybe?
		message = formattedTweets[0]
	}
	return message, err
}

func init() {
	// TODO initialize Twitter client here
	// we should only need to create a Twitter.Client once on startup
	bot.RegisterPassiveCommand(
		"twitter",
		expandTweet)
}
