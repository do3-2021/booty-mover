package move

import (
	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var command = &discordgo.ApplicationCommand{
	Name:        "ohyeah",
	Description: "shouts an oh yeah in a channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "channel",
			Type:         discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
			Description:  "Voice channel where an 'oh yeah' will be launched",
			Required:     true,
		},
	},
}

func SendError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// get channel

	options := i.ApplicationCommandData().Options
	channelId := options[0].Value.(string)

	verbosity.Debug("Channel ID: " + channelId)

	// response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "oooooooooooooooooooh Yeaaaaaaaaaaaaaaaah",
		},
	})
}

var Descriptor = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
