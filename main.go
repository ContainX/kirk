package main

import (
	"fmt"
	"github.com/ContainX/kirk/config"
	"github.com/ContainX/kirk/db"
	"github.com/ContainX/kirk/slackTeam"
)

func main() {

	// Connect to Mongo
	_, mongoSession := db.Init()
	defer mongoSession.Close()

	// Create botInstances
	fmt.Println("Recreating bot instances")
	err, configs := config.GetAllConfig()
	if err != nil {
		panic("Could not get team configs")
	}
	for _, config := range configs {
		if config.Access_token != "" {
			slackTeam.Watch(config.Access_token)
			fmt.Println("Watching team")
		} else {
			fmt.Println("Config missing access token", config)
		}
	}

	// Start HTTP Server
	router := getRouter()
	router.Run(":8080")
}
