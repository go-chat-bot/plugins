# Setting up Twitter API credentials

- Request Twitter development access [here](https://developer.twitter.com) (Note: the approval process takes about a week)
- Once your dev access is approved, create an App [here](https://developer.twitter.com/en/apps)
- Visit the App's "Keys and Tokens" tab
- Save the *Consumer API Key* and *Consumer API key secret* in a safe place (it is not necessary to generate *Access token* and *Access token secret* for this plugin)
- Export the required environment variables into the shell in which your go-chat-bot process will run:

```
export TWITTER_CONSUMER_KEY="yourconsumerkeyhere" \
       TWITTER_CONSUMER_SECRET="yourconsumersecrethere"
```
