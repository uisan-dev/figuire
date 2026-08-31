package bot

import (
	"errors"
	"figuire/figuire/amiami"
	"figuire/figuire/store"
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

func (b *Bot) HandleWatch(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.GuildID == "" {
		respond(s, i, "This command only works in a server channel")
		return
	}

	code := i.ApplicationCommandData().Options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("deferring %s: %v", code, err)
		return
	}

	item, err := FiguireBot.Client.FetchItem(code)
	if err != nil {
		msg := "Something went wrong reaching AmiAmi. Try again in a moment"
		if errors.Is(err, amiami.ErrItemNotFound) {
			msg = "No item found with code `" + code + "`"
		}
		log.Printf("watch %s: %v", code, err)
		followup(s, i, msg)
		return
	}

	c, err := b.Store.Record(item)
	if err != nil {
		log.Printf("watch record %s: %v", code, err)
		followup(s, i, "Couldnt add **"+item.Name+"** to watchlist. Try again in a moment")
		return
	}

	err = b.Store.AddWatch(&store.Watch{
		GuildID:   i.GuildID,
		ChannelID: i.ChannelID,
		ProductID: c.Product.ID,
		AddedBy:   i.Member.User.ID,
	})

	if err != nil {
		if errors.Is(err, store.ErrAlreadyWatching) {
			followup(s, i, "Already watching **"+item.Name+"** in this channel")
			return
		} else {
			followup(s, i, "Couldnt add **"+item.Name+"** to watchlist. Try again in a moment")
			return
		}
	}

	isPreowned := "Yes"

	if !item.PreOwned {
		isPreowned = "No"
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{{
			Title:       item.Name,
			URL:         item.URL,
			Color:       0x5865F2,
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: item.ImageURL},
			Description: "Watching: **" + item.Name + "**",
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Price", Value: yen(item.PriceJPY), Inline: true},
				{Name: "Status", Value: string(item.Availability.AsDisplayString()), Inline: true},
				{Name: "Preowned", Value: isPreowned, Inline: true},
				{Name: "Release", Value: orDash(item.ReleaseText), Inline: true},
			},
			Footer: &discordgo.MessageEmbedFooter{Text: item.Code},
		}},
	})

	if err != nil {
		log.Printf("watch followup %s: %v", code, err)
	}
}
