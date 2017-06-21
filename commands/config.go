package commands

import (
	"strings"
	"github.com/jeremyroberts0/kirk/config"
)
var jiraBaseUrlConfigKey string = "jira-base-url"
var projectsConfigKey string = "jira-projects"

func get (userCommand []string, teamConfig *config.TeamConfigStruct) string {
	if len(userCommand) == 2 {
		return jiraBaseUrlConfigKey + ": " + teamConfig.JiraBaseUrl + "\n" + projectsConfigKey + ": " + strings.Join(teamConfig.SubscribedProjects, ", ")
	} else if len(userCommand) >= 3 {
		switch userCommand[2] {
		case jiraBaseUrlConfigKey:
			return "The JIRA. Base url is: " +  teamConfig.JiraBaseUrl
		case projectsConfigKey:
			return "The follow. projects. are tracked: " + strings.Join(teamConfig.SubscribedProjects, ", ")
		}
		return "Could not. find. config value. " + userCommand[2]
	}

	return unknownResponse
}

func set (userCommand []string, teamConfig *config.TeamConfigStruct) string {
	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case jiraBaseUrlConfigKey:
			// Slack adds < and > when URLs are detected
			teamConfig.JiraBaseUrl = strings.Trim(userCommand[3], "<>")
			return "I've changed. Your JIRA Base. URL. to " + teamConfig.JiraBaseUrl
		}

		return userCommand[2] + " is not. a. config value."
	} else if len(userCommand) >= 5 {
		return "Please pass. one value. after the config. key."
	}

	return unknownResponse
}

func add (userCommand []string, teamConfig *config.TeamConfigStruct) string {
	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case projectsConfigKey:
			for _, project := range teamConfig.SubscribedProjects {
				if project == userCommand[3] {
					return project + " is. already tracked."
				}
			}
			teamConfig.SubscribedProjects = append(teamConfig.SubscribedProjects, userCommand[3])
			return userCommand[3] + " project. added. to tracked projects"
		}
	}

	return unknownResponse
}

func remove (userCommand []string, teamConfig *config.TeamConfigStruct) string {
	if len(userCommand) >= 4 {
		switch userCommand[2] {
		case projectsConfigKey:
			for index, project := range teamConfig.SubscribedProjects {
				if project == userCommand[3] {
					teamConfig.SubscribedProjects = append(teamConfig.SubscribedProjects[:index], teamConfig.SubscribedProjects[index+1:]...)
					return userCommand[3] + " project. no longer. tracked"
					break
				}
			}
		}
	}
	return unknownResponse
}

func configCommand (userCommand []string, teamConfig *config.TeamConfigStruct) string {
	switch userCommand[1] {
	case "get":
		return get(userCommand, teamConfig)
	case "set":
		return set(userCommand, teamConfig)
	case "add":
		return add(userCommand, teamConfig)
	case "remove":
		return remove(userCommand, teamConfig)
	}

	return unknownResponse
}