package twitter

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/go-chat-bot/bot"
	"os"
	"reflect"
	"testing"
)

func TestTwitter(t *testing.T) {
	t.Skip()
	//got := fmt.Printf("%s", var)
	cmd := &bot.PassiveCmd{}
	_, err := expandTweet(cmd)
	got := "foo" //<-result.Message
	want := "bar"
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestFindTweetIDs(t *testing.T) {
	t.Parallel()
	var cases = []struct {
		text   string
		result []int64
	}{
		{"this message has no links", nil},
		{"http://twitter.com/jbouie/status/1247273759632961537", []int64{1247273759632961537}},
		{"https://twitter.com/jbouie/status/1247273759632961537", []int64{1247273759632961537}},
		{"https://mobile.twitter.com/jbouie/status/1247273759632961537", []int64{1247273759632961537}},
		{"wow check out this tweet https://mobile.twitter.com/jbouie/status/1247273759632961537", []int64{1247273759632961537}},
		{"wow check out this tweethttps://mobile.twitter.com/jbouie/status/1247273759632961537", []int64{1247273759632961537}},
		{"wow check out this tweet https://mobile.twitter.com/jbouie/status/1247273759632961537super cool right?", []int64{1247273759632961537}},
		{"https://twitter.com/dmackdrwns/status/1217830568848764930/photo/1", []int64{1217830568848764930}},
		// FIXME these should all pass
		//{"https://mobile.twitter.com/jbouie/status/1247273759632961537 https://twitter.com/NicolleDWallace/status/1260032806832336900", []int64{1247273759632961537, 1260032806832336900}},
		//{"https://mobile.twitter.com/jbouie/status/1247273759632961537 https://mobile.twitter.com/NicolleDWallace/status/1260032806832336900", []int64{1247273759632961537, 1260032806832336900}},
		//{"https://mobile.twitter.com/NicolleDWallace/status/1260032806832336900 https://mobile.twitter.com/jbouie/status/1247273759632961537", []int64{1260032806832336900, 1247273759632961537}},
	}
	for _, c := range cases {
		t.Run("toot", func(t *testing.T) {
			got, err := findTweetIDs(c.text)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(c.result, got) {
				t.Errorf("got %+v; want %+v", got, c.result)
			}
		})

	}
}

func TestNewAuthenticatedTwitterClient(t *testing.T) {
	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	var cases = []struct {
		key, secret   string
		expectedError error
	}{
		{key: "", secret: "", expectedError: errors.New("missing API credentials")},
		// FIXME can't seem to capture the exact error, but it is a variation of 403 Forbidden
		//{key: "asdf", secret: "jklmnop", expectedError: errors.New("403 Forbidden")},
		{key: twitterConsumerKey, secret: twitterConsumerSecret, expectedError: nil},
	}
	for i, c := range cases {
		t.Run(string(i), func(t *testing.T) {
			_, err := newAuthenticatedTwitterClient(c.key, c.secret)
			if err != nil && err.Error() != c.expectedError.Error() {
				t.Errorf("got %s; want %s", err, c.expectedError)
			}
		})
	}
}

func TestFetchTweets(t *testing.T) {
	t.Skip()
	t.Parallel() // this test function can run in parallel with other tests
	// given a slice of tweet IDs, I should get back a slice of Tweets
	nilTweet := twitter.Tweet{}
	var cases = []struct {
		tweetIDs   []int64
		tweetsText []string
	}{
		{[]int64{123456789}, []string{}},
		// FIXME this should not pass
		{[]int64{1247273759632961537}, []string{"hi a status"}},
		{[]int64{1247273759632961537}, []string{"This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common peop… https://t.co/sH7mwGkaC4"}},
		{[]int64{1258736595953475584}, []string{"This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common peop… https://t.co/sH7mwGkaC4"}},
	}
	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		t.Run(string(c.tweetIDs[0]), func(t *testing.T) {
			t.Parallel() // test cases may run in parallel
			// TODO it should be possible to test this without making actual twitter calls
			for i := range c.tweetIDs {
				got, err := fetchTweets(client, c.tweetIDs)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(nilTweet, got[i]) && c.tweetsText != nil {
					t.Errorf("got %+v; want %+v", got[i], c.tweetsText[i])
				}
				if !reflect.DeepEqual(nilTweet, got[i]) && got[i].Text != c.tweetsText[i] {
					t.Errorf("got %+v; want %+v", got[i], c.tweetsText)
				}
			}
		})
	}
}

func TestFetchTweet(t *testing.T) {
	var cases = []struct {
		tweetID       int64
		tweetText     string
		expectedError error
	}{
		{123456789, "", errors.New("404 Not Found")},
		{1247273759632961537, `This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common people” instead reveal their cloistered, condescending view of ordinary people. https://t.co/KV8xnG2w48`, nil},
		{1259875673994338305, `This is the first U.S. presidential election in which "Vote Him Out Before He Kills You and Your Family" is a wholly reasonable slogan for the challenger`, nil},
	}
	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		t.Run(string(c.tweetID), func(t *testing.T) {
			// TODO it should be possible to test this without making actual twitter calls
			got, err := fetchTweet(client, c.tweetID)
			if err != nil && err.Error() != c.expectedError.Error() {
				t.Error(err)
			}
			if got.FullText != c.tweetText {
				t.Errorf("got %s; want %s", got.FullText, c.tweetText)
			}
		})
	}
}

func TestFormatTweets(t *testing.T) {
	t.Parallel() // this test function can run in parallel with other tests
	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	if err != nil {
		t.Fatal(err)
	}
	// TODO how to read in some data from testdata
	jbouieStatus, jbouieErr := fetchTweet(client, 1247273759632961537)
	if jbouieErr != nil {
		t.Fatal(err)
	}
	yamicheStatus, yamicheErr := fetchTweet(client, 1258736595953475584)
	if yamicheErr != nil {
		t.Fatal(err)
	}
	// given a slice of tweets, I should get back a slice of nicely formatted message strings

	var cases = []struct {
		tweets   []twitter.Tweet
		messages []string
	}{
		{[]twitter.Tweet{*jbouieStatus}, []string{"Tweet from @jbouie: This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common people” instead reveal their cloistered, condescending view of ordinary people. https://t.co/KV8xnG2w48"}},
		{[]twitter.Tweet{*yamicheStatus}, []string{"Tweet from @Yamiche: A quote that will be much talked about.   President Trump said on Fox News moments ago: “I learned a lot from Richard Nixon, don’t fire people. I learned a lot by watching Richard Nixon.”  He added, “I did nothing wrong and there are no tapes in my case.”"}},
	}
	for _, c := range cases {
		t.Run(c.tweets[0].User.ScreenName, func(t *testing.T) {
			got := formatTweets(c.tweets)
			want := c.messages
			if want[0] != got[0] {
				t.Errorf("got %s, want %s", got, want)
			}
		})

	}
}

func TestExpandTweet(t *testing.T) {
	t.Parallel() // this test function can run in parallel with other tests

	var (
		twitterConsumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
		twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	)
	client, err := newAuthenticatedTwitterClient(twitterConsumerKey, twitterConsumerSecret)
	if err != nil {
		t.Fatal(err)
	}
	// TODO how to read in some data from testdata
	jbouieStatus, jbouieErr := fetchTweet(client, 1247273759632961537)
	if jbouieErr != nil {
		t.Fatal(err)
	}
	yamicheStatus, yamicheErr := fetchTweet(client, 1258736595953475584)
	if yamicheErr != nil {
		t.Fatal(err)
	}
	// given a slice of tweets, I should get back a slice of nicely formatted message strings

	var cases = []struct {
		tweets   []twitter.Tweet
		messages []string
	}{
		{[]twitter.Tweet{*jbouieStatus}, []string{"Tweet from @jbouie: This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common people” instead reveal their cloistered, condescending view of ordinary people. https://t.co/KV8xnG2w48"}},
		{[]twitter.Tweet{*yamicheStatus}, []string{"Tweet from @Yamiche: A quote that will be much talked about.   President Trump said on Fox News moments ago: “I learned a lot from Richard Nixon, don’t fire people. I learned a lot by watching Richard Nixon.”  He added, “I did nothing wrong and there are no tapes in my case.”"}},
	}
	for _, c := range cases {
		t.Run(c.tweets[0].User.ScreenName, func(t *testing.T) {
			got := formatTweets(c.tweets)
			want := c.messages
			if want[0] != got[0] {
				t.Errorf("got %s, want %s", got, want)
			}
		})

	}
}
