package main

import (
	"github.com/nlopes/slack"
	"github.com/jeremyroberts0/kirk/commands"
	"github.com/jeremyroberts0/kirk/config"
	"github.com/jeremyroberts0/kirk/db"
	"log"
	"os"
	"fmt"
	"strings"
	"regexp"
)


func main() {
	_, mongoSession := db.Init()
	defer mongoSession.Close()

	// Create botInstances

	botConfig := map[string]string{
		"SLACK_BOT_ACCESS_TOKEN": os.Getenv("SLACK_BOT_ACCESS_TOKEN"),
	}

	fmt.Println(botConfig)

	api := slack.New(botConfig["SLACK_BOT_ACCESS_TOKEN"])
	logger := log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	teamInfo, teamInfoErr := api.GetTeamInfo()
	teamId := teamInfo.ID
	authTest, botUserErr := api.AuthTest()
	if teamInfoErr != nil || botUserErr != nil {
		fmt.Println(teamInfoErr, botUserErr)
		panic("Initialization error")
	} else {
		fmt.Println("Your captain is here")
	}

	teamConfig := config.GetTeamConfig(teamInfo.ID)

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		//case *slack.ConnectedEvent:
		//	fmt.Println("Kirk is connected to Slack")
		//	fmt.Println("Connection counter:", event.ConnectionCount)
		case *slack.MessageEvent:
			if event.User != authTest.UserID {
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
					issueIdRegex := "(" + strings.Join(teamConfig.Subscribed_projects, "|") + `)-\d+`
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

	fmt.Println("Kirk. is started.")
}