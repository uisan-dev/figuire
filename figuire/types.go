package figuire

import "github.com/bwmarrin/discordgo"

type Handler func(*discordgo.Session, *discordgo.InteractionCreate)
