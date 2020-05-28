# Overview

The Twitter plugin scrapes message text for Twitter URLs, then attempts to fetch the linked Tweet and post it to the channel, like so:

```
08:29:19 <user> https://twitter.com/simonpierce/status/1265829199115218945
08:29:21 <chatbot> Tweet from @simonpierce: Gentoo penguins like to exercise their growing chicks by makingthem run around the colony, squawking hysterically, if they want to get fed.  I'm not saying it'd be fun to try this with your own kids if you're stuck at home... but I'm not not saying that. https://t.co/2Y0wewRKDw
```

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
