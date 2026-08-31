package bot

import (
	"figuire/figuire"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) GetHandlers() map[string]figuire.Handler {
	return map[string]figuire.Handler{
		"ping":  b.HandlePing,
		"watch": b.HandleWatch,
	}
}

func (b *Bot) GetCommands() []*discordgo.ApplicationCommand {
	commands := make([]*discordgo.ApplicationCommand, 0)
	commands = append(commands, StatusCommands...)
	commands = append(commands, WatchCommands...)
	return commands
}
