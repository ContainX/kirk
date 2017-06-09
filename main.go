package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"fmt"
	"strings"
	"unicode"
)

func main() {
	config := map[string]string{
		"SLACK_BOT_ACCESS_TOKEN": "xoxb-196008087415-z3m59xI3h3DMMZtLovYsZY40",
	}

	//jiraUrlBase := "https://jira.auction.com/browse/"
	jiraProjectKey := "XOEY"



	api := slack.New(config["SLACK_BOT_ACCESS_TOKEN"])
	logger := log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received")

		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Infos:", event.Info)
			fmt.Println("Connection counter:", event.ConnectionCount)
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("This is. your captain.  I. am. here.", "C5QF2Q4NL"))
		case *slack.MessageEvent:
			text := event.Text
			if strings.Contains(text, jiraProjectKey) {
				// TODO: Transformation not working
				indexOfProjectKey := strings.Index(text, jiraProjectKey)
				indexOfCharAfterIssueKey := strings.IndexFunc(
					text[indexOfProjectKey:len(text)-1],
					func(c rune) bool {
						return !unicode.IsDigit(c)
					},
				)
				jiraIssueId := text[indexOfProjectKey:indexOfCharAfterIssueKey]
				fmt.Println(text, indexOfProjectKey, indexOfCharAfterIssueKey)
				fmt.Printf("JIRA ISSUE FOUND: %v\n", jiraIssueId)
			}
		}
	}

	fmt.Println("Kirk. is started.")
}