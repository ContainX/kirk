package commands

import (
	"strings"
	"github.com/jeremyroberts0/kirk/config"
)

const unknownResponse string = "I. don't know. how to handle. that command.  Type `help`. for more. info."

func HandleCommand (commandText string, teamConfig *config.TeamConfigStruct) string {
	userCommand := strings.Split(commandText, " ")

	switch userCommand[0] {
	case "config":
		return configCommand(userCommand, teamConfig)
	case "help":
		return helpCommand()
	}

	return unknownResponse
}
