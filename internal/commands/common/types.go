package common

import "github.com/bwmarrin/discordgo"

type CommandDescriptor struct {
	Command *discordgo.ApplicationCommand
	Execute func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
