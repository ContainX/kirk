package slackTeam

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ContainX/kirk/commands"
	"github.com/ContainX/kirk/config"
	"github.com/ContainX/kirk/tracer"
	"github.com/nlopes/slack"
)

func handleMessage(event *slack.MessageEvent, rtm slack.RTM, teamInfo slack.TeamInfo, botUserId string) {
	teamId := teamInfo.ID

	// Instrument time it takes to handle the message
	t := tracer.Timer("handleMessage.latency", []string{"team:" + teamId})
	defer t.End()

	if event.User != botUserId {
		// Don't respond to messages from self
		if strings.HasPrefix(event.Channel, "D") == true {
			// D prefix in channel ID is indicative of direct message
			imDetails, err := rtm.GetIMChannels()
			if err != nil {
				logger.Fatal("Could not get IM channels")
			} else {
				messageIsIm := false
				for _, imChannel := range imDetails {
					if imChannel.ID == event.Channel {
						messageIsIm = true
						break
					}
				}

				if messageIsIm == true {
					// Direct message to bot
					responseText := commands.HandleCommand(event.Text, teamId)
					rtm.SendMessage(
						rtm.NewOutgoingMessage(
							responseText,
							event.Channel,
						),
					)
				}
			}
		} else {
			// Message is not special, run through default formatter
			err, teamConfig := config.GetTeamConfig(teamInfo.ID)
			if err != nil {
				fmt.Println("Error getting team config when processing message")
			} else {
				tracer.Get().Incr("projects.watched", []string{"team:", teamId}, 1)
				issueIdRegex := "(?i)(" + strings.Join(teamConfig.Subscribed_projects, "|") + `)-\d+`
				issueIdRe := regexp.MustCompile(issueIdRegex)
				issueIds := issueIdRe.FindAllString(event.Text, -1)
				if issueIds != nil {
					idLinkPrefix := "/browse/"
					issueIdAlreadyLinksRe := regexp.MustCompile(idLinkPrefix + issueIdRegex)
					issueIdsAlreadyLinks := issueIdAlreadyLinksRe.FindAllString(event.Text, -1)
					newText := ""
					for _, issueId := range issueIds {
						// Detect if captured issueId is already wrapped in a link
						// Must do it this way because golang doesn't support negative lookahead in regex
						issueIdIsLink := false
						for _, linkFragment := range issueIdsAlreadyLinks {
							if idLinkPrefix+issueId == linkFragment {
								issueIdIsLink = true
								continue
							}
						}
						if issueIdIsLink == true {
							continue
						}
						newText += teamConfig.Jira_base_url + "/browse/" + issueId + "\n"
					}

					fmt.Println("Jira issues recongized")
					if newText != "" {
						rtm.SendMessage(
							rtm.NewOutgoingMessage(
								"Here are. some. links to JIRA:\n"+newText,
								event.Channel,
							),
						)
					}
				}
			}
		}
	}
}
