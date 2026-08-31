package commands

import (
	"figurebot/figuire"

	"github.com/bwmarrin/discordgo"
)

var handlers map[string]figuire.Handler = map[string]figuire.Handler{
	"ping":  HandlePing,
	"watch": HandleWatch,
}

func GetHandlers() map[string]figuire.Handler {
	return handlers
}

func GetCommands() []*discordgo.ApplicationCommand {
	commands := make([]*discordgo.ApplicationCommand, 0)
	commands = append(commands, StatusCommands...)
	commands = append(commands, WatchCommands...)
	return commands
}
