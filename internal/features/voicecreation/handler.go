package voicecreation

import (
	"strings"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
)

// when there is no empty voice channel in a category , create a new voice channel in this category
// when there is more than one empty voice channel in a category, delete all of them except one
func Handle(s *discordgo.Session, evt *discordgo.VoiceStateUpdate) {

	channelID := evt.ChannelID

	// if joining a channel
	if channelID != "" {
		go func() {
			err := updateVoiceOfCategory(s, channelID)
			if err != nil {
				verbosity.Error(err)
			}
		}()
	}

	// update the last channel the user was in
	if evt.BeforeUpdate != nil {
		go func() {

			err := updateVoiceOfCategory(s, evt.BeforeUpdate.ChannelID)
			if err != nil {
				verbosity.Error(err)
			}
		}()
	}

}

func updateVoiceOfCategory(session *discordgo.Session, channelID string) (err error) {

	channel, err := session.Channel(channelID)
	if err != nil {
		verbosity.Error(err)
		return
	}

	// we only want a voice channel with the emoji
	if channel.Type != discordgo.ChannelTypeGuildVoice || !strings.HasPrefix(channel.Name, "🔉") {
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
		if guildChannel.ParentID == category && guildChannel.Type == discordgo.ChannelTypeGuildVoice && strings.HasPrefix(guildChannel.Name, "🔉") {
			if !usedVoiceChannels[guildChannel.ID] {
				emptyChannels++
			}
			if emptyChannels > 1 {
				// delete this channel
				st, err := session.ChannelDelete(guildChannel.ID)
				verbosity.Debug("deleted channel ", st.Name, " ", st.ID)
				if err != nil {
					verbosity.Error(err)
				}
			}
		}
	}

	if emptyChannels == 0 {
		// create a new channel
		st, err := session.GuildChannelCreate(channel.GuildID, channel.Name, discordgo.ChannelTypeGuildVoice)
		if err != nil {
			verbosity.Error(err)
		}
		verbosity.Debug("created channel ", st.Name, " ", st.ID)
		_, err = session.ChannelEditComplex(st.ID, &discordgo.ChannelEdit{
			ParentID: category,
			Position: channel.Position + 1,
		})
		if err != nil {
			verbosity.Error(err)
		}
	}

	return
}
