package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var StatusCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping the bot",
	},
}

func HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
	if err != nil {
		log.Printf("responding to ping: %v", err)
	}
}
