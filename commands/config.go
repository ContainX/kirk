package commands

import (
	"strings"
	"github.com/jeremyroberts0/kirk/config"
	"fmt"
)
var jiraBaseUrlConfigKey string = "jira-base-url"
var projectsConfigKey string = "jira-projects"
const noConfigMessage string = "I. could not find. your. team's config"


// Struct values must be upper case for them to make it back and forth from mongo

func get (userCommand []string, teamId string) string {
	// TODO: Mongo query to update in place, without having to make extra trip to get the config
	err, teamConfig := config.GetTeamConfig(teamId)
	if err != nil {
		return noConfigMessage
	}

	if len(userCommand) == 2 {
		return jiraBaseUrlConfigKey + ": " + teamConfig.Jira_base_url + "\n" + projectsConfigKey + ": " + strings.Join(teamConfig.Subscribed_projects, ", ")
	} else if len(userCommand) >= 3 {
		switch userCommand[2] {
		case jiraBaseUrlConfigKey:
			return "The JIRA. Base url is: " +  teamConfig.Jira_base_url
		case projectsConfigKey:
			return "The follow. projects. are tracked: " + strings.Join(teamConfig.Subscribed_projects, ", ")
		}
		return "Could not. find. config value. " + userCommand[2]
	}

	return unknownResponse
}

func set (userCommand []string, teamId string) string {
	// TODO: Mongo query to update in place, without having to make extra trip to get the config
	err, teamConfig := config.GetTeamConfig(teamId)
	if err != nil {
		return noConfigMessage
	}
	configCollection := config.GetConfigCollection()

	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case jiraBaseUrlConfigKey:
			// Slack adds < and > when URLs are detected
			teamConfig.Jira_base_url = strings.Trim(userCommand[3], "<>")
			updateErr := configCollection.UpdateId(teamConfig.Id, teamConfig)
			if updateErr != nil {
				fmt.Println("Update error", updateErr)
				return "I. Encountered a problem. Updating your configuration"
			}
			return "I've changed. Your JIRA Base. URL. to " + teamConfig.Jira_base_url
		}

		return userCommand[2] + " is not. a. config value."
	} else if len(userCommand) >= 5 {
		return "Please pass. one value. after the config. key."
	}

	return unknownResponse
}

func add (userCommand []string, teamId string) string {
	err, teamConfig := config.GetTeamConfig(teamId)
	if err != nil {
		return noConfigMessage
	}
	configCollection := config.GetConfigCollection()

	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case projectsConfigKey:
			for _, project := range teamConfig.Subscribed_projects {
				if project == userCommand[3] {
					return project + " is. already tracked."
				}
			}
			teamConfig.Subscribed_projects = append(teamConfig.Subscribed_projects, userCommand[3])
			configCollection.UpdateId(teamConfig.Id, teamConfig)
			return userCommand[3] + " project. added. to tracked projects"
		}
	}

	return unknownResponse
}

func remove (userCommand []string, teamId string) string {
	err, teamConfig := config.GetTeamConfig(teamId)
	if err != nil {
		return noConfigMessage
	}
	configCollection := config.GetConfigCollection()

	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case projectsConfigKey:
			for index, project := range teamConfig.Subscribed_projects {
				if project == userCommand[3] {
					teamConfig.Subscribed_projects = append(teamConfig.Subscribed_projects[:index], teamConfig.Subscribed_projects[index+1:]...)
					configCollection.UpdateId(teamConfig.Id, teamConfig)
					return userCommand[3] + " project. no longer. tracked"
					break
				}
			}
		}
	}
	return unknownResponse
}

func configCommand (userCommand []string, teamId string) string {
	switch userCommand[1] {
	case "get":
		return get(userCommand, teamId)
	case "set":
		return set(userCommand, teamId)
	case "add":
		return add(userCommand, teamId)
	case "remove":
		return remove(userCommand, teamId)
	}

	return unknownResponse
}