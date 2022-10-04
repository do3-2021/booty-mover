package groupcreator

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Creates the category and the channels for the group
func createCategory(s *discordgo.Session, i *discordgo.InteractionCreate, wg *sync.WaitGroup, roleID string, groupName string, description string, isPrivate bool) (err error) {

	categoryName := groupName

	if isPrivate {
		categoryName = "🔒 " + categoryName
	}

	category, err := s.GuildChannelCreateComplex(
		i.GuildID,
		discordgo.GuildChannelCreateData{
			Name: categoryName,
			Topic: groupName,
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
				Topic: description,
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
				Topic: description,
				Type:     discordgo.ChannelTypeGuildVoice,
				ParentID: category.ID,
			},
		)
	}()

	return
}
