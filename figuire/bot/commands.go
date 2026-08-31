package bot

import (
	"figuire/figuire"

	"github.com/bwmarrin/discordgo"
)

var handlers map[string]figuire.Handler = map[string]figuire.Handler{
	"ping":  FiguireBot.HandlePing,
	"watch": FiguireBot.HandleWatch,
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
