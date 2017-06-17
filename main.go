package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"fmt"
	"strings"
	"regexp"
)

type TeamConfigStruct struct {
	subscribedProjects []string
	jiraBaseUrl string
}
var configByTeamId = make(map[string]TeamConfigStruct)

func handleCommand (commandText string, teamConfig *TeamConfigStruct) string {
	userCommand := strings.Split(commandText, " ")
	
	jiraBaseUrlConfigKey := "jira-base-url"
	projectsConfigKey := "jira-projects"
	if len(userCommand) >= 2 {
		switch userCommand[0] {
		case "config":
			switch userCommand[1] {
			case "get":
				if len(userCommand) == 2 {
					return jiraBaseUrlConfigKey + ": " + teamConfig.jiraBaseUrl + "\n" + projectsConfigKey + ": " + strings.Join(teamConfig.subscribedProjects, ", ")
				} else if len(userCommand) >= 3 {
					switch userCommand[2] {
					case jiraBaseUrlConfigKey:
						return "The JIRA. Base url is: " +  teamConfig.jiraBaseUrl
					case projectsConfigKey:
						return "The follow. projects. are tracked: " + strings.Join(teamConfig.subscribedProjects, ", ")
					}
				}
			case "set":
				if len(userCommand) >= 4 {
					switch userCommand[2] {
					case jiraBaseUrlConfigKey:
						// Slack adds < and > when URLs are detected
						teamConfig.jiraBaseUrl = strings.Trim(userCommand[3], "<>")
						return "I've changed. Your JIRA Base. URL. to " + teamConfig.jiraBaseUrl
					}
				}
			case "add":
				if len(userCommand) >= 4 {
					switch userCommand[2] {
					case projectsConfigKey:
						for _, project := range teamConfig.subscribedProjects {
							if project == userCommand[3] {
								return project + " is. already tracked."
							}
						}
						teamConfig.subscribedProjects = append(teamConfig.subscribedProjects, userCommand[3])
						return userCommand[3] + " project. added. to tracked projects"
					}
				}
			case "remove":
				if len(userCommand) >= 4 {
					switch userCommand[2] {
					case projectsConfigKey:
						for index, project := range teamConfig.subscribedProjects {
							if project == userCommand[3] {
								teamConfig.subscribedProjects = append(teamConfig.subscribedProjects[:index], teamConfig.subscribedProjects[index+1:]...)
								return userCommand[3] + " project. no longer. tracked"
								break
							}
						}
					}
				}
			}
		}
	}

	return "I don't know. how to handle the. command. " + commandText
}

func main() {
	botConfig := map[string]string{
		"SLACK_BOT_ACCESS_TOKEN": "xoxb-196008087415-z3m59xI3h3DMMZtLovYsZY40",
	}

	api := slack.New(botConfig["SLACK_BOT_ACCESS_TOKEN"])
	logger := log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	teamInfo, teamInfoErr := api.GetTeamInfo()
	if teamInfoErr != nil {
		panic("Could not get team info")
	} else if _, ok := configByTeamId[teamInfo.ID]; !ok {
		configByTeamId[teamInfo.ID] = TeamConfigStruct{
			subscribedProjects: []string{"XOEY", "VOLT"},
			jiraBaseUrl: "https://jira.auction.com",
		}
	}

	teamConfig := configByTeamId[teamInfo.ID]

	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received")

		switch event := msg.Data.(type) {
		//case *slack.ConnectedEvent:
		//	fmt.Println("Kirk is connected to Slack")
		//	fmt.Println("Connection counter:", event.ConnectionCount)
		case *slack.MessageEvent:
			if strings.HasPrefix(event.Channel, "D") == true {
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
						responseText := handleCommand(event.Text, &teamConfig)
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
				issueIdRegex := "(" + strings.Join(teamConfig.subscribedProjects, "|") + `)-\d+`
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
						newText += teamConfig.jiraBaseUrl + "/browse/" + issueId + "\n"
					}

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

	fmt.Println("Kirk. is started.")
}