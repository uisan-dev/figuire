package bot

import (
	"errors"
	"figuire/figuire/amiami"
	"figuire/figuire/colors"
	"figuire/figuire/store"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type pending struct {
	change  *store.Change
	watches []store.Watch
}

func (b *Bot) CheckAll(s *discordgo.Session) {
	products, err := b.Store.WatchedProducts()
	if err != nil {
		log.Printf("loading watched products: %v", err)
	}

	if len(products) == 0 {
		return
	}

	var queue []pending
	var failures int
	for _, p := range products {
		item, err := b.Client.FetchItem(p.Code)
		if err != nil {
			if errors.Is(err, amiami.ErrItemNotFound) {
				log.Printf("check: %s no longer on AmiAmi", p.Code)
			} else {
				failures++
				log.Printf("check %s: %v", p.Code, err)
			}
			continue
		}

		c, err := b.Store.Record(item)
		if err != nil {
			log.Printf("record %s: %v", p.Code, err)
			continue
		}

		if !c.Any() {
			continue
		}

		watches, err := b.Store.WatchesFor(p.ID)
		if err != nil {
			log.Printf("watches for %s: %v", p.Code, err)
			continue
		}

		if len(watches) == 0 {
			continue
		}

		queue = append(queue, pending{change: c, watches: watches})
	}

	log.Printf("Checked %d products, %d changed, %d failed", len(products), len(queue), failures)

	if failures > len(products)/2 {
		log.Printf("Aborting notifications: %d/%d fetches failed", failures, len(products))
		return
	}

	for _, p := range queue {
		embed := embedForChange(p.change)
		if embed == nil {
			continue
		}

		for _, w := range p.watches {
			b.post(s, w, embed)
		}
	}
}

func (b *Bot) post(s *discordgo.Session, w store.Watch, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(w.ChannelID, embed)
	if err == nil {
		return
	}

	var restErr *discordgo.RESTError
	if errors.As(err, &restErr) && restErr.Response != nil {
		switch restErr.Response.StatusCode {
		case 403, 404:
			log.Printf("Channel %s unreachable (%d), removing watching", w.ChannelID, restErr.Response.StatusCode)
			if _, dbErr := b.Store.RemoveWatch(w.ChannelID, w.ProductID); dbErr != nil {
				log.Printf("removing dead watch: %v", dbErr)
			}
			return
		}
	}

	log.Printf("posting to %s: %v", w.ChannelID, err)
}

func embedForChange(c *store.Change) *discordgo.MessageEmbed {
	if c.IsNew {
		return nil
	}

	p := c.Product
	e := &discordgo.MessageEmbed{
		URL:    p.URL,
		Footer: &discordgo.MessageEmbedFooter{Text: p.Code},
	}

	if p.ImageURL != "" {
		e.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: p.ImageURL}
	}

	switch {
	case c.BecameAvailable():
		e.Title = "In stock: " + p.Name
		e.Description = yen(p.PriceJPY)
		e.Color = colors.ColorGreen
	case c.PriceChanged():
		if c.OldPriceJPY > c.NewPriceJPY {
			e.Title = "Price drop: " + p.Name
			e.Description = fmt.Sprintf("%s -> **%s** (**%.2f%%** change)", yen(c.OldPriceJPY), yen(c.NewPriceJPY), float32(c.NewPriceJPY-c.OldPriceJPY)/float32(c.OldPriceJPY)*100)
			e.Color = colors.ColorGreen
		} else {
			e.Title = "Price increase: " + p.Name
			e.Description = fmt.Sprintf("%s -> **%s** (**%.2f%%** change)", yen(c.OldPriceJPY), yen(c.NewPriceJPY), float32(c.NewPriceJPY-c.OldPriceJPY)/float32(c.OldPriceJPY)*100)
			e.Color = colors.ColorOrange
		}
	case c.AvailabilityChanged():
		e.Title = "Availability changed: " + p.Name
		e.Description = fmt.Sprintf("%s -> **%s**", c.OldAvailability.AsDisplayString(), c.NewAvailability.AsDisplayString())
		e.Color = colors.ColorBlue
	default:
		return nil
	}

	e.Title = truncate(e.Title, 256)
	return e
}
