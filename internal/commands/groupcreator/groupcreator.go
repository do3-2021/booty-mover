package groupcreator

import (
	"fmt"
	"strings"

	"github.com/NilsPonsard/verbosity"

	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var command = &discordgo.ApplicationCommand{
	Name:        "create-group",
	Description: "Create a new group, add the -p option to make it private",
	Options:[]*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "name",
			Required: true,
			Description: "Name of the group to be created",
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "visibility",
			Required: false,
			Description: "Wether the group is private or not. A private group is joignable only when you are invited",
		},
	} ,
}

func contains(s []string, e string) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}

//This command will store The User, give him a role that grants him acess to a fresh new channel
//
func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.Member

	verbosity.Debug("User", member.User.ID, ",tries to create group '", i.ApplicationCommandData().Options[0].Value, "'")
	groupName := strings.Replace(i.ApplicationCommandData().Options[0].StringValue(), " ", "-", -1)

	guild, error := s.GuildChannels(i.GuildID)
	var groups []string
	
	verbosity.Debug("Guild: ", guild)

	for i := 0; i < len(guild); i += 1 {
		if guild[i].Type == discordgo.ChannelTypeGuildCategory {
			verbosity.Debug(guild[i].Name)
			groups = append(groups, guild[i].Name)
		}
	}

	verbosity.Debug(groups)

	if (contains(groups, groupName)) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Error: group %v already exists", groupName),
			},
		})

		return;
	}

	category, error := s.GuildChannelCreateComplex(
		i.GuildID, 
		discordgo.GuildChannelCreateData{
			Name: groupName,
			Type: discordgo.ChannelTypeGuildCategory,
		},
	)

	if (error != nil) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%v", error),
			},
		})

		return;
	}

	s.GuildChannelCreateComplex(
		i.GuildID, 
		discordgo.GuildChannelCreateData{
			Name: "txt-" + groupName,
			Type: discordgo.ChannelTypeGuildText,
			ParentID: category.ID,
		},
	)

	s.GuildChannelCreateComplex(
		i.GuildID, 
		discordgo.GuildChannelCreateData{
			Name: "voc-" + groupName,
			Type: discordgo.ChannelTypeGuildVoice,
			ParentID: category.ID,
		},
	)


	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Created goup '%v'! 🎉🎉", groupName),
		},
	})
}

var Summary = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
