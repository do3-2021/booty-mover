package voicecreation

import (
	"strings"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
)

func isCompatibleChannelName(name string) bool {
	return strings.HasPrefix(name, "🔉") || strings.HasPrefix(name, "🕐")
}

func isTemporaryChannelName(name string) bool {
	return strings.HasPrefix(name, "🕐")
}

// when there is no empty voice channel in a category , create a new voice channel in this category
// when there is more than one empty voice channel in a category, delete all of them except one
func Handle(s *discordgo.Session, evt *discordgo.VoiceStateUpdate) {

	channelID := evt.ChannelID

	// if joining a channel
	if channelID != "" {
		err := updateVoiceOfCategory(s, channelID)
		if err != nil {
			verbosity.Error(err)
		}
	}

	// update the last channel the user was in
	if evt.BeforeUpdate != nil {
		err := updateVoiceOfCategory(s, evt.BeforeUpdate.ChannelID)
		if err != nil {
			verbosity.Error(err)
		}
	}

}

func updateVoiceOfCategory(session *discordgo.Session, channelID string) (err error) {

	channel, err := session.Channel(channelID)
	if err != nil {
		verbosity.Error(err)
		return
	}

	// we only want a voice channel with the emoji
	if channel.Type != discordgo.ChannelTypeGuildVoice || !isCompatibleChannelName(channel.Name) {
		return
	}

	category := channel.ParentID

	// not in a category
	if category == "" {
		return
	}

	guild, err := session.State.Guild(channel.GuildID)

	if err != nil {
		verbosity.Error(err)
		return
	}

	// find all used voice channels
	usedVoiceChannels := make(map[string]bool)
	for _, state := range guild.VoiceStates {
		usedVoiceChannels[state.ChannelID] = true
	}

	// find all voice channels in the category

	emptyChannels := 0

	for _, guildChannel := range guild.Channels {
		if guildChannel.ParentID == category && guildChannel.Type == discordgo.ChannelTypeGuildVoice && isCompatibleChannelName(guildChannel.Name) {
			if !usedVoiceChannels[guildChannel.ID] {
				emptyChannels++
			}
			if emptyChannels > 1 && isTemporaryChannelName(guildChannel.Name) {
				// delete this channel
				st, err := session.ChannelDelete(guildChannel.ID)

				if err != nil {
					verbosity.Error("error deleting channel : ", err)
				} else {
					verbosity.Debug("deleted channel ", st.Name, " ", st.ID)
				}
			}
		}
	}

	if emptyChannels == 0 {

		baseName := strings.TrimPrefix(strings.TrimPrefix(channel.Name, "🔉"), "🕐")
		splittedName := strings.Split(baseName, "#")

		// create a new channel
		st, err := session.GuildChannelCreate(channel.GuildID, baseName, discordgo.ChannelTypeGuildVoice)
		if err != nil {
			verbosity.Error(err)
		}
		suffix := st.ID[len(st.ID)-4:]

		if len(splittedName) > 1 {
			baseName = strings.Join(splittedName[:len(splittedName)-1], "#")
		}

		name := "🕐" + baseName + "#" + suffix

		_, err = session.ChannelEditComplex(st.ID, &discordgo.ChannelEdit{
			ParentID: category,
			Position: channel.Position + 1,
			Name:     name,
		})
		if err != nil {
			verbosity.Error(err)
		}
		verbosity.Debug("created channel ", name, " ", st.ID)
	}

	return
}
