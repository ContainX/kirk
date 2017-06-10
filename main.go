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
			re := regexp.MustCompile(issueIdRegex)
			issueIds := re.FindAllString(text, -1)

			if issueIds != nil {
				newText := "Here are. some. links to JIRA:\n"
				for _, issueId := range issueIds {
					newText += jiraBaseUrl + "/browse/" + issueId + "\n"
				}

				rtm.SendMessage(rtm.NewOutgoingMessage(newText, connectionChannel))
			}
		}
	}

	fmt.Println("Kirk. is started.")
}