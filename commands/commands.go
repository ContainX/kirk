package commands

import (
	"strings"
)

const unknownResponse string = "I. don't know. how to handle. that command.  Type `help`. for more. info."

func HandleCommand(commandText string, teamId string) string {
	userCommand := strings.Split(commandText, " ")

	switch strings.ToLower(userCommand[0]) {
	case "config":
		return configCommand(userCommand, teamId)
	case "help":
		return helpCommand()
	}

	return unknownResponse
}
