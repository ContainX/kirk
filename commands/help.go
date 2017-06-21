package commands

func helpCommand() string {
	return `I. am kirk.
I will make interacting with JIRA easier in Slack.
If I see a JIRA ticket ID in a channel I'm a member of, I'll post links to the ticket in the channel.

Send me a direct message to access the following commands:
` + "`help`" + ` - This message
` + "`config get`" + ` Get my configuration
` + "`config set " + jiraBaseUrlConfigKey + ` <jira base url>` + "`" + ` - Set your JIRA base URL (e.g. http://jira.mydomain.com)
` + "`config add " + projectsConfigKey + ` <jira project ID>` + "`" + ` - Add a project for me to watch for (e.g. VOLT)
` + "`config remove " + projectsConfigKey + ` <jira project ID>` + "`" + ` - Remove a project I'm watching for (e.g. VOLT)
`
}