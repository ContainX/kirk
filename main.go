package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"fmt"
	"strings"
	"regexp"
)

func main() {
	config := map[string]string{
		"SLACK_BOT_ACCESS_TOKEN": "xoxb-196008087415-z3m59xI3h3DMMZtLovYsZY40",
	}

	api := slack.New(config["SLACK_BOT_ACCESS_TOKEN"])
	logger := log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	const connectionChannel string = "C5QF2Q4NL"
	const jiraBaseUrl string = "https://jira.auction.com"
	projects := []string{"VOLT", "XOEY"}

	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received")

		switch event := msg.Data.(type) {
		//case *slack.ConnectedEvent:
		//	fmt.Println("Kirk is connected to Slack")
		//	fmt.Println("Connection counter:", event.ConnectionCount)
		case *slack.MessageEvent:
			text := event.Text

			issueIdRegex := "(" + strings.Join(projects, "|") + `)-\d+`
			issueIdRe := regexp.MustCompile(issueIdRegex)
			issueIds := issueIdRe.FindAllString(text, -1)

			// Message contains issueIds
			if issueIds != nil {
				idLinkPrefix := "/browse/"
				issueIdAlreadyLinksRe := regexp.MustCompile(idLinkPrefix + issueIdRegex)
				issueIdsAlreadyLinks := issueIdAlreadyLinksRe.FindAllString(text, -1)
				newText := ""
				for _, issueId := range issueIds {
					// Detect if captured issueId is already wrapped in a link
					// Must do it this way because golang doesn't support negative lookahead in regex
					issueIdIsLink := false
					for _, linkFragment := range issueIdsAlreadyLinks {
						if idLinkPrefix + issueId == linkFragment {
							issueIdIsLink = true
							continue
						}
					}
					if issueIdIsLink == true {
						continue
					}
					newText += jiraBaseUrl + "/browse/" + issueId + "\n"
				}

				if newText != "" {
					rtm.SendMessage(
						rtm.NewOutgoingMessage(
							"Here are. some. links to JIRA:\n" + newText,
							connectionChannel,
						),
					)
				}

			}
		}
	}

	fmt.Println("Kirk. is started.")
}