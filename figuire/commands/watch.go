package commands

import (
	"errors"
	"figurebot/figuire"
	"figurebot/figuire/amiami"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var WatchCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "watch",
		Description: "Watch an AmiAmi product for price and stock changes",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "AmiAmi gcode, e.g FIGURE-207417",
				Required:    true,
			},
		},
	},
}

var WatchHandlers = map[string]figuire.Handler{
	"watch": HandleWatch,
}

func HandleWatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	code := i.ApplicationCommandData().Options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("deferring %s: %v", code, err)
	}

	item, err := figuire.Client.FetchItem(code)
	if err != nil {
		msg := "Something went wrong reaching AmiAmi. Try again in a moment"
		if errors.Is(err, amiami.ErrItemNotFound) {
			msg = "No item found with code `" + code + "`"
		}
		log.Printf("watch %s: %v", code, err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: msg})
		return
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{{
			Title:       item.Name,
			URL:         item.URL,
			Description: fmt.Sprintf("¥%d - %s", item.PriceJPY, item.Availability.AsDisplayString()),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: item.ImageURL},
		}},
	})
}
