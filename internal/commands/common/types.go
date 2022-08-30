package common

import "github.com/bwmarrin/discordgo"

type Summary struct {
	Command *discordgo.ApplicationCommand
	Execute func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
