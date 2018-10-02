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

By default the plugin will output issues in following format:
```
<key> (<assignee>, <status>): <summary> - <url>
```
Optionally you can set up per-channel configuration by setting environment
variables of JIRA_CHAN_TEMPL_<channel> with value of new template. To see
which values are available for use in templates see
[go-jira](https://github.com/andygrunwald/go-jira/blob/master/issue.go).

The format used is go template notation on the issue object. If you want to just
post URL to the issue itself when the channel is #jirabot you can configure it
by setting environment as such:
```
export JIRA_CHAN_TEMPL_jirabot="{{.Self}}"
```

Default template looks as like this:
```
{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): {{.Fields.Summary}} - {{.Self}}
```

### TODO
* Notifications when new issues are created
