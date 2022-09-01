package groupcreator

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Creates the category and the channels for the group
func createCategory(s *discordgo.Session, i *discordgo.InteractionCreate, wg *sync.WaitGroup, roleID string, groupName string, description string) (err error) {

	category, err := s.GuildChannelCreateComplex(
		i.GuildID,
		discordgo.GuildChannelCreateData{
			Name: groupName,
			Type: discordgo.ChannelTypeGuildCategory,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    roleID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: ROLE_PERMISSIONS,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: ROLE_PERMISSIONS,
				},
			},
		},
	)

	if err != nil {
		return
	}

	// Creating text channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GuildChannelCreateComplex(
			i.GuildID,
			discordgo.GuildChannelCreateData{
				Name:     "txt-" + groupName,
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: category.ID,
			},
		)
	}()

	// Creating voice channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GuildChannelCreateComplex(
			i.GuildID,
			discordgo.GuildChannelCreateData{
				Name:     "🔉voc-" + groupName,
				Type:     discordgo.ChannelTypeGuildVoice,
				ParentID: category.ID,
			},
		)
	}()

	return
}
