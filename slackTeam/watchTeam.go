package slackTeam

import (
	"fmt"
	"github.com/jeremyroberts0/kirk/commands"
	"github.com/jeremyroberts0/kirk/config"
	"github.com/nlopes/slack"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func Watch(token string) {
	var startedWg sync.WaitGroup
	startedWg.Add(1)
	go watchTeam(token, &startedWg)
	startedWg.Wait()
}

func watchTeam(token string, startedWg *sync.WaitGroup) {
	fmt.Println("Preparing to watch", token)
	logger := log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	api := slack.New(token)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// TODO: Could we save teamId in DB on install?  Then we don't have to call getTeamInfo on startup
	teamInfo, teamInfoErr := api.GetTeamInfo()
	teamId := teamInfo.ID
	authTest, botUserErr := api.AuthTest()
	if teamInfoErr != nil || botUserErr != nil {
		fmt.Println(teamInfoErr, botUserErr)
		panic("Initialization error")
	} else {
		fmt.Println("Listening to team", teamId)
	}

	// Don't start a listener if we can't find the team's config
	err, _ := config.GetTeamConfig(teamInfo.ID)
	if err != nil {
		fmt.Println(err)
		startedWg.Done()
	} else {
		startedWg.Done()
		for msg := range rtm.IncomingEvents {
			fmt.Printf("Event Received %+v\n", msg)
			switch event := msg.Data.(type) {
			//TODO: Listen for token revoked event and invalidate in DB
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
						err, teamConfig := config.GetTeamConfig(teamInfo.ID)
						if err != nil {
							fmt.Println("Error getting team config when processing message")
						} else {
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
		}
	}
}
