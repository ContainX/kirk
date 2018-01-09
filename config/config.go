package config

import (
	"errors"
	"fmt"

	"github.com/ContainX/kirk/db"
	"github.com/ContainX/kirk/tracer"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TeamConfig struct {
	Id                  bson.ObjectId `bson:"_id,omitempty"`
	Team_id             string        `bson:"team_id,omitempty"`
	Jira_base_url       string        `bson:"jira_base_url,omitempty"`
	Subscribed_projects []string      `bson:"subscribed_projects,omitempty"`
	Access_token        string        `bson:"access_token,omitempty"`
}

func GetConfigCollection() *mgo.Collection {
	return db.Instance.C("config")
}
func GetAllConfig() (error, []TeamConfig) {
	var results []TeamConfig
	collection := GetConfigCollection()
	err := collection.Find(nil).All(&results)
	return err, results
}
func GetTeamConfig(teamId string) (error, TeamConfig) {
	// TODO: Cache team config so we don't have to go to DB on every message
	// TODO: Caching will require updates run through this file as well
	teamConfig := TeamConfig{}
	collection := GetConfigCollection()
	query := bson.M{"team_id": teamId}
	err := collection.Find(query).One(&teamConfig)

	if err != nil {
		// Should never get here, team config should be created during oauth flow
		fmt.Println("Team Config not found", err)
		return errors.New("Error getting team config"), TeamConfig{}
	}

	return nil, teamConfig
}
func AddNewTeam(accessToken string) error {
	api := slack.New(accessToken)
	teamInfo, teamInfoErr := api.GetTeamInfo()

	if teamInfoErr != nil {
		fmt.Println("Error getting team info", teamInfoErr)
		return errors.New("Error getting team info")
	}

	collection := GetConfigCollection()
	teamConfig := TeamConfig{}

	_, err := collection.Upsert(
		bson.M{"team_id": teamInfo.ID},
		TeamConfig{
			Team_id:      teamInfo.ID,
			Access_token: accessToken,
		},
	)

	if err != nil {
		fmt.Println("Error creating team config", err)
		return errors.New("Error creating team config")
	}

	tracer.Get().Incr("teamAdded", []string{"team:" + teamInfo.ID}, 1)

	return nil
}
