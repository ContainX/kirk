package main

import (
	"github.com/jeremyroberts0/kirk/db"
	"os"
	"fmt"
)

func main() {

	oauthScopes := []string{
		"channels:read",
		"chat:write:user", // or bot?
		"team:read",
	}


	fmt.Println(oauthScopes)

	_, mongoSession := db.Init()
	defer mongoSession.Close()

	// Create botInstances
	botConfig := map[string]string{
		"SLACK_BOT_ACCESS_TOKEN": os.Getenv("SLACK_BOT_ACCESS_TOKEN"),
	}
	fmt.Println(botConfig)

	watchTeam(botConfig["SLACK_BOT_ACCESS_TOKEN"])

	fmt.Println("Kirk. is started.")
}