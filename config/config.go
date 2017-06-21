package config

type TeamConfigStruct struct {
	SubscribedProjects []string
	JiraBaseUrl string
}
var ConfigByTeamId = make(map[string]TeamConfigStruct)
