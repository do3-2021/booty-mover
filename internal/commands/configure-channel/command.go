package configurechannel

import (
	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	"github.com/do3-2021/booty-mover/internal/database"
	"github.com/do3-2021/booty-mover/internal/guild"
)

var command = &discordgo.ApplicationCommand{
	Name:        "set-group-channel",
	Description: "Set this channel as the group channel",
}

func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	db, err := database.GetDB()

	message := ""

	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		verbosity.Error(err)

		message = "Error getting the db :/ "
		return
	}

	err = guild.SetGroupChannel(db, i.ChannelID, i.GuildID)

	if err != nil {
		verbosity.Error(err)

		message = "Error updating the db :/"

		return
	}

	message = "Channel confiured !"

}

var Descriptor = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
