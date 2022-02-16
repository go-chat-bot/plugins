### Overview

This is a plugin for [Atlassian JIRA](https://www.atlassian.com/software/jira)
issue tracking system.

### Features
* Simple authentication support for JIRA
* Outputs some of the issue details in the channels not just links
* Optional per-channel configuration of issue message format

### Setup
* Set up JIRA_BASE_URL env variable to your JIRA server URL. For example
  https://issues.jenkins-ci.org
* Set up JIRA_USER env variable to JIRA username for the bot account
* Set up JIRA_PASS env variable to JIRA password for the bot account
* Optional: set up JIRA_TOKEN env variable, if your instance requires
  using personal access tokens (user/pass no longer need to be defined).

In addition to the above channel-specific configuration variables can be defined
in a separate JSON configuration file loaded from path specified by environment
variable `JIRA_CONFIG_FILE`. Example file can be seen in
`example_config.json`. It is an array of channel configurations with each
configuration having:
 * `channel` for which the configuration is intended
 * `template` to override default issue template (see Issue Formatting)
 * `templateNew` to override default issue template for new issue notifications
 * `templateResolved` to override default issue template for resolved issue notifications
 * `notifyNew` is array of JIRA project keys to watch for new issues
 * `notifyResolved` is array JIRA project keys to watch for resolved issues

### Issue Formatting

By default the plugin will output issues in the following format:
```
<key> (<assignee>, <status>): <summary> - <url>
```
To see which values are available for use in templates see
[go-jira](https://github.com/andygrunwald/go-jira/blob/master/issue.go).

The format used is go template notation on the issue object. If you want to just
post URL to the issue itself you can configure it by setting the template to
`{{.Self}}` for given channel in the configuration file.

Default template looks like this:
```
{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): {{.Fields.Summary}} - {{.Self}}
```

`JIRA_NOTIFY_INTERVAL` environment variable can be used to control how often the
notification methods will be run. It defaults to be run every minute.

### Threaded notifications
**NOTE:** This feature has only been tested in Google Chat. The person who wrote this code
          does not use this bot in any other platform. Feel free to contribute to make it
          work in your prefered platform (in case it supports threads).

In Google Chat, each notification will create a new thread by default. In some cases, it might
be desirable to restrict to a single thread, for cleaningness. Due to a limitation in the API,
the thread must exist first. Once a thread is created, you must fetch the full URL. There are
many different methods for this, so use whatever is better for you, but a thread URL
should look similar to one of these examples:
* `https://chat.google.com/room/<roomID>/<threadID>`
* `https://mail.google.com/chat/u/0/#chat/space/<roomID>/<threadID>`

Once you have that information, your `JIRA_CONFIG_FILE` should look like
[example_config_thread.json](example_config_thread.json).

Also you, need to start the bot with `JIRA_THREAD=true` environment variable defined.

### Verbose log
If JIRA_VERBOSE variable is defined (any value) the bot generates a log
every time it queries JIRA.
