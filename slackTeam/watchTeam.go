package slackTeam

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ContainX/kirk/config"
	"github.com/ContainX/kirk/tracer"
	"github.com/nlopes/slack"
)

var logger = log.New(os.Stdout, "slack-bot:", log.Lshortfile|log.LstdFlags)

func Watch(token string) {
	var startedWg sync.WaitGroup
	startedWg.Add(1)
	go watchTeam(token, &startedWg)
	startedWg.Wait()
}

func watchTeam(token string, startedWg *sync.WaitGroup) {
	fmt.Println("Preparing to watch", token)

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

	botUserId := authTest.UserID

	// Don't start a listener if we can't find the team's config
	err, _ := config.GetTeamConfig(teamInfo.ID)
	if err != nil {
		fmt.Println(err)
		startedWg.Done()
	} else {
		startedWg.Done()
		for msg := range rtm.IncomingEvents {
			//Print all events, for debug purposes
			//fmt.Printf("Event Received %+v\n", msg)
			tracer.Get().Incr("teams.active", []string{"team:" + teamId}, 1)
			switch event := msg.Data.(type) {
			//TODO: Listen for token revoked event and invalidate in DB
			//case *slack.ConnectedEvent:
			//	fmt.Println("Kirk is connected to Slack")
			//	fmt.Println("Connection counter:", event.ConnectionCount)
			case *slack.ChannelJoinedEvent:
				handleChannelJoin(*rtm, *event)
			case *slack.MessageEvent:
				handleMessage(event, *rtm, *teamInfo, botUserId)
			}
		}
	}
}
