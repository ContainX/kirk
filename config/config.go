package config

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/jeremyroberts0/kirk/db"
	"fmt"
)

type TeamConfig struct{
	Id                 	bson.ObjectId `bson:"_id"`
	Team_id             string
	Jira_base_url       string
	Subscribed_projects []string
}

func GetConfigCollection () *mgo.Collection {
	return db.Instance.C("config")
}
func GetTeamConfig(teamId string) TeamConfig {
	teamConfig := TeamConfig{}
	collection := GetConfigCollection()
	query := bson.M{"team_id": teamId}
	err := collection.Find(query).One(&teamConfig)

	if err != nil {
		fmt.Println("Config not found", err)
		// Setup config if it is missing
		insertErr := collection.Insert(&TeamConfig{
			Team_id:             teamId,
			Jira_base_url:       "BASE_URL_NOT_SET",
			Subscribed_projects: make([]string, 0),
		})
		if insertErr != nil {
			fmt.Println("Error inserting new config for team", teamId, insertErr)
		}

		collection.Find(query).One(&teamConfig)
		fmt.Println("Created new config for team", teamId,  teamConfig)
	}

	return teamConfig
}
