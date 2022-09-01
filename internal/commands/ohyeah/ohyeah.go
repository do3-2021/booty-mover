package ohyeah

import (
	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	"github.com/do3-2021/booty-mover/internal/commands/ohyeah/dgvoice"
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

	// play audio in a channel
	voiceConnection, err := s.ChannelVoiceJoin(i.GuildID, channelId, false, true)
	voiceConnection.Speaking(true)

	defer voiceConnection.Speaking(false)
	defer voiceConnection.Disconnect()

	if err != nil {
		SendError(s, i, "Booty-mover cannot join channel")
		return
	}

	dgvoice.PlayAudioFile(voiceConnection, "internal/commands/ohyeah/sound.mp3", make(chan bool))

	// sender := voiceConnection.OpusSend

	// data, err := os.ReadFile("internal/commands/ohyeah/sound.mp3")

	// verbosity.Info(data)

	// if err != nil {
	// 	SendError(s, i, "Cannot read sound file")
	// 	return
	// }

	// sender <- data

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
