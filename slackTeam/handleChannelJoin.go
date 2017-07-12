package slackTeam

import "github.com/nlopes/slack"

const channelJoinMessage = `I've joined. This channel. I will provide. convenient links to. jira. whenever I see. an issue. id
To learn more. about me. Send me a. direct. message.`

func handleChannelJoin(rtm slack.RTM, event slack.ChannelJoinedEvent) {
	rtm.SendMessage(
		rtm.NewOutgoingMessage(
			channelJoinMessage,
			event.Channel.ID,
		),
	)
}
