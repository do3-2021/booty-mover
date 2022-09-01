package groupcreator

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendErrorMessage(s *discordgo.Session, i *discordgo.InteractionCreate, err string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("❌ Error: %v", err),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func SendSuccessMessage(s *discordgo.Session, i *discordgo.InteractionCreate, groupName string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Created goup '%v'! 🎉", groupName),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
