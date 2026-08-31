package main

import (
	"figurebot/figuire/commands"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmds := commands.GetCommands()
	handlers := commands.GetHandlers()

	token := os.Getenv("DISCORD_TOKEN")
	appID := os.Getenv("DISCORD_APP_ID")
	guildID := os.Getenv("DISCORD_GUILD_ID")

	if token == "" || appID == "" {
		log.Fatal("DISCORD_TOKEN and/or DISCORD_APP_ID is not set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("creating session: %v", err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name
		h, ok := handlers[name]
		if !ok {
			log.Printf("No handler for command %s", name)
			return
		}

		h(s, i)
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s#%s", r.User.Username, r.User.Discriminator)
	})

	if err := session.Open(); err != nil {
		log.Fatalf("opening connection: %v", err)
	}
	defer session.Close()

	registered := make([]*discordgo.ApplicationCommand, 0, len(commands.StatusCommands))
	for _, cmd := range cmds {
		rc, err := session.ApplicationCommandCreate(appID, guildID, cmd)
		if err != nil {
			log.Fatalf("registering %s: %v", cmd, err)
		}
		log.Printf("Registered %s", rc.Name)
		registered = append(registered, rc)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	for _, cmd := range registered {
		if err := session.ApplicationCommandDelete(appID, guildID, cmd.ID); err != nil {
			log.Printf("removing %s: %v", cmd.Name, err)
		}
	}

	log.Println("Shut down with no errors")
}
