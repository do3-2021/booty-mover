package move

import (
	"errors"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var command = &discordgo.ApplicationCommand{
	Name:        "move-voice-all",
	Description: "Move all users from one voice channel to another",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "source_channel",
			Type:         discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
			Description:  "Source voice channel",
			Required:     true,
		},
		{
			Name:         "destination_channel",
			Type:         discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
			Description:  "Destination voice channel",
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

func IsSameChannels(channelOne, channelTwo string) error {
	if channelOne == channelTwo {
		return errors.New("Same source and destination channels")
	}

	return nil
}

func GetChannelMember(s *discordgo.Session, i *discordgo.InteractionCreate, channelId string) ([]string, error) {
	guild, err := s.State.Guild(i.GuildID)

	if err != nil {
		return []string{}, errors.New("Cannot get guild state")
	}

	// get members
	voiceStates := guild.VoiceStates

	var members []string
	for _, voiceState := range voiceStates {
		if voiceState.ChannelID == channelId {
			members = append(members, voiceState.UserID)
		}
	}

	verbosity.Debug("Source voice channel members: ", members)

	return members, nil
}

func MoveUsers(s *discordgo.Session, i *discordgo.InteractionCreate, members []string, destinationChannelId string) error {
	for _, member := range members {
		err := s.GuildMemberMove(i.GuildID, member, &destinationChannelId)

		if err != nil {
			return errors.New("Cannot move member " + member)
		}

		memberName, _ := s.State.Member(i.GuildID, member)
		verbosity.Debug("User ", memberName.User.Username, " moved")
	}

	return nil
}

func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// get source and destination channels

	options := i.ApplicationCommandData().Options
	sourceChannelId := options[0].StringValue()
	destinationChannelId := options[1].StringValue()

	sameChannelError := IsSameChannels(sourceChannelId, destinationChannelId)
	if sameChannelError != nil {
		SendError(s, i, sameChannelError.Error())
		return
	}

	verbosity.Debug("Source channel ID: " + sourceChannelId)
	verbosity.Debug("Destination channel ID: " + destinationChannelId)

	// get members of source channel

	sourceChannelMembers, err := GetChannelMember(s, i, sourceChannelId)

	if err != nil {
		SendError(s, i, err.Error())
		return
	}

	// move members

	err = MoveUsers(s, i, sourceChannelMembers, destinationChannelId)

	if err != nil {
		SendError(s, i, err.Error())
		return
	}

	// response

	sourceChannelName, sourceChannelNameError := s.Channel(sourceChannelId)
	destinationChannelName, destinationChannelNameError := s.Channel(destinationChannelId)

	if sourceChannelNameError != nil || destinationChannelNameError != nil {
		SendError(s, i, "Cannot get source or destination channel name")
		return
	}

	message := "All users of the " + sourceChannelName.Name + " channel have been moved to " + destinationChannelName.Name
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

var Descriptor = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
