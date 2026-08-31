package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
	if err != nil {
		log.Printf("responding: %v", err)
	}
}

func followup(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
	})
	if err != nil {
		log.Printf("following up: %v", err)
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func yen(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}

	return "¥" + b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "..."
}
