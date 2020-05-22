package twitter

import (
	"errors"
	"github.com/go-chat-bot/bot"
	"regexp"
	"testing"
)

func TestTwitter(t *testing.T) {
	// given a message string, I should get back a response message string
	// containing one or more parsed Tweets
	jbouieOutput := `Tweet from @jbouie: This falls into one of my favorite genres of tweets, bona fide elites whose pretenses to understanding “common people” instead reveal their cloistered, condescending view of ordinary people. https://t.co/KV8xnG2w48`
	//yamicheOutput := `Tweet from @Yamiche: A quote that will be much talked about.   President Trump said on Fox News moments ago: “I learned a lot from Richard Nixon, don’t fire people. I learned a lot by watching Richard Nixon.”  He added, “I did nothing wrong and there are no tapes in my case.”`
	sethAbramsonOutput := `Tweet from @SethAbramson: This is the first U.S. presidential election in which "Vote Him Out Before He Kills You and Your Family" is a wholly reasonable slogan for the challenger`
	dmackdrwnsOutput := `Tweet from @dmackdrwns: It was pretty fun to try to manifest creatures plucked right from the minds of manic children.  #georgiamuseumofart https://t.co/C983t6QjmT`

	var cases = []struct {
		input, output string
		expectedError error
	}{
		{
			input:         "this message has no links",
			output:        "",
			expectedError: nil,
		}, {
			input:         "http://twitter.com/jbouie/status/1247273759632961537",
			output:        jbouieOutput,
			expectedError: nil,
		}, {
			input:         "https://mobile.twitter.com/jbouie/status/1247273759632961537",
			output:        jbouieOutput,
			expectedError: nil,
		}, {
			input:         "wow check out this tweet https://mobile.twitter.com/jbouie/status/1247273759632961537",
			output:        jbouieOutput,
			expectedError: nil,
		}, {
			input:         "wow check out this tweethttps://mobile.twitter.com/jbouie/status/1247273759632961537",
			output:        jbouieOutput,
			expectedError: nil,
		}, {
			input:         "wow check out this tweet https://mobile.twitter.com/jbouie/status/1247273759632961537super cool right?",
			output:        jbouieOutput,
			expectedError: nil,
		}, {
			input:         "https://twitter.com/dmackdrwns/status/1217830568848764930/photo/1",
			output:        dmackdrwnsOutput,
			expectedError: nil,
		}, {
			input:         "http://twitter.com/jbouie/status/123456789",
			output:        "",
			expectedError: errors.New("twitter: 144 No status found with that ID."),
		}, {
			input:         "https://twitter.com/SethAbramson/status/1259875673994338305 lol bye",
			output:        sethAbramsonOutput,
			expectedError: nil,
		},

		// FIXME these should all pass
		//{"https://mobile.twitter.com/jbouie/status/1247273759632961537 https://twitter.com/NicolleDWallace/status/1260032806832336900", []int64{1247273759632961537, 1260032806832336900}},
		//{"https://mobile.twitter.com/jbouie/status/1247273759632961537 https://mobile.twitter.com/NicolleDWallace/status/1260032806832336900", []int64{1247273759632961537, 1260032806832336900}},
		//{"https://mobile.twitter.com/NicolleDWallace/status/1260032806832336900 https://mobile.twitter.com/jbouie/status/1247273759632961537", []int64{1260032806832336900, 1247273759632961537}},
	}
	for i, c := range cases {
		testingUser := bot.User{
			ID:       "test",
			Nick:     "test",
			RealName: "test",
			IsBot:    true,
		}
		testingMessage := bot.Message{
			Text:     c.input,
			IsAction: false,
		}
		testingCmd := bot.PassiveCmd{
			Raw:         c.input,
			Channel:     "test",
			User:        &testingUser,
			MessageData: &testingMessage,
		}
		t.Run(string(i), func(t *testing.T) {
			// these CANNOT run concurrently
			// FIXME panic here when no credentials
			got, err := expandTweets(&testingCmd)
			want := c.output
			if err != nil && err.Error() != c.expectedError.Error() {
				t.Error(err)
			}
			if got != want {
				t.Errorf("got %+v; want %+v", got, want)
			}
		})
	}
}

func TestNewAuthenticatedTwitterClient(t *testing.T) {
	// TODO test case for these envvars not being set
	key, secret, err := getCredentialsFromEnvironment()
	if err != nil {
		t.Error(err)
	}
	var cases = []struct {
		key, secret   string
		expectedError error
	}{
		{key: "", secret: "", expectedError: errors.New("missing API credentials")},
		{key: "asdf", secret: "jklmnop", expectedError: errors.New(`Get https://api.twitter.com/1.1/application/rate_limit_status.json?resources=statuses: oauth2: cannot fetch token: 403 Forbidden Response: {"errors":[{"code":99,"message":"Unable to verify your credentials","label":"authenticity_token_error"}]}`)},
		{key: key, secret: secret, expectedError: errors.New("")},
	}
	newlines := regexp.MustCompile(`\r?\n`)
	for i, c := range cases {
		t.Run(string(i), func(t *testing.T) {
			// these CANNOT run concurrently
			_, err := newAuthenticatedTwitterClient(c.key, c.secret)
			if err != nil {
				// eat newlines because they mess with our tests
				got := newlines.ReplaceAllString(err.Error(), " ")
				want := c.expectedError.Error()
				if got != want {
					t.Errorf("got %s; want %s", got, want)
				}
			}
		})
	}
}
