# Cachet plugin

This plugin can provide notifications for services that are failed in
[Cachet](https://cachethq.io/).

## Configuration

There are two environment variables that this plugin needs:

* CACHET_API - URL of your Cachet top-level API endpoint.
  Example API URL `https://status.company.com/api`
* CACHET_ALERT_CONFIG - Path to file with notification configuration

Alert configuration is a JSON file which can be edited using bot commands. You
can also edit it manually (remember to restart the bot)

JSON file is a list of objects which have following keys:

* `channel` - name of channels the configuration is for
* `repeatGap` - number of minutes between repeated outage notifications
* `services` - cachet component names which will be notified (or `all` for any outage)

Example:

```json
[
  {
    "channel": "#outages",
    "services": [
      "all"
    ],
    "repeatGap": 120
  },
  {
    "channel": "#service1",
    "services": [
      "service"
    ],
    "repeatGap": 15
  },
  {
    "channel": "#team",
    "services": [
      "service2",
      "service3"
    ],
    "repeatGap": 60
  }
]
```

Above configuration would sent alerts to `#outage` for any outage every 2
hours. It would send alerts to `#service1` every 15 minutes for `service`
outages. And it would also send alerts to `#team` every 60 minutes if either
`service2` or `service3` are in outage.

## Commands

Bot recognizes following commands:

* `services` - list all services known to cachet
* `subscriptions` - list all active subscriptions for this channel
* `subscribe <service>` - subscribe to receive outage notification for `<service>`
* `unsubscribe <service>` - unsubscribe from outage notification for `<service>`
* `repeatgap <minutes>` - set how often alerts will be repeated (in minutes)

Configuration is automatically saved on each change through bot commands
