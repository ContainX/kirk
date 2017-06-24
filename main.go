package main

import (
	"github.com/nlopes/slack"
	"github.com/jeremyroberts0/kirk/commands"
	"github.com/jeremyroberts0/kirk/config"
	"log"
	"os"
	"fmt"
	"strings"
	"regexp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type Person struct {
	Name string
	Phone string
}


func main() {
	mongoSession, mongoErr := mgo.Dial("mongodb://mongo:27017")
	db := mongoSession.DB("kirk")
	//defer mongoSession.Close()

	if mongoErr != nil {
		fmt.Println("Error connecting to mongo", mongoErr)
		panic(mongoErr)
	}

	// Mongo Testing
	testCollection := db.C("test")
	err := testCollection.Insert(&Person{"Jeremy", "408-555-4444"})
	if err != nil {
		panic(err)
	}

	result := Person{}
	query := testCollection.Find(bson.M{"name": "Jeremy"})
	err = query.One(&result)
	if err != nil {
		panic(err)
	}
	count, err := query.Count()
	if err != nil {
		panic(err)
	}

	fmt.Println(count, result)


	// End Mongo Testing



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
	authTest, botUserErr := api.AuthTest()
	if teamInfoErr != nil || botUserErr != nil {
		fmt.Println(teamInfoErr, botUserErr)
		panic("Initialization error")
	} else if _, ok := config.ConfigByTeamId[teamInfo.ID]; !ok {
		config.ConfigByTeamId[teamInfo.ID] = config.TeamConfigStruct{
			SubscribedProjects: []string{"XOEY", "VOLT"},
			JiraBaseUrl: "https://jira.auction.com",
		}
		fmt.Println("Your captain is here")
	}

	teamConfig := config.ConfigByTeamId[teamInfo.ID]

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
							responseText := commands.HandleCommand(event.Text, &teamConfig)
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
					issueIdRegex := "(" + strings.Join(teamConfig.SubscribedProjects, "|") + `)-\d+`
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
							newText += teamConfig.JiraBaseUrl + "/browse/" + issueId + "\n"
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