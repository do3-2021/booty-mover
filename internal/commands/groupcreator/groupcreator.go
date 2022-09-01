package groupcreator

import (
	"fmt"
	"strings"
	"sync"

	"github.com/NilsPonsard/verbosity"

	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var command = &discordgo.ApplicationCommand{
	Name:        "create-group",
	Description: "Create a new group, optionnaly private",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Required:    true,
			Description: "Name of the group to be created",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "description",
			Required:    true,
			Description: "Description of the group",
			MaxLength:   200,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "visibility",
			Required:    false,
			Description: "Wether the group is private or not. A private group is joignable only when you are invited",
		},
	},
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// This command will store The User, give him a role that grants him acess to a fresh new channel
func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	member := i.Member

	// Handle the "join group" button press
	if i.Type == discordgo.InteractionMessageComponent {
		role := strings.Replace(i.MessageComponentData().CustomID, "create-group-", "", 1)

		s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "You have Sucessfully joined the group! 🎉",
			},
		})
		return
	}

	// get args

	description := ""
	groupName := ""
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "description":
			description = option.StringValue()
		case "name":
			groupName = strings.Replace(option.StringValue(), " ", "-", -1)
		}
	}

	guild, error := s.GuildChannels(i.GuildID)
	if error != nil {
		SendErrorMessage(s, i, error.Error())
		return
	}

	verbosity.Debug("User", member.User.ID, ", tries to create group '", i.ApplicationCommandData().Options[0].Value, "' in guild ", guild)

	var groups []string

	for i := 0; i < len(guild); i += 1 {
		if guild[i].Type == discordgo.ChannelTypeGuildCategory {
			verbosity.Debug(guild[i].Name)
			groups = append(groups, guild[i].Name)
		}
	}

	if contains(groups, groupName) {
		SendErrorMessage(s, i, fmt.Sprintf("group %v already exists", groupName))
		verbosity.Debug(fmt.Sprintf("group %v already exists", groupName))

		return
	}

	role, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name:        "grp-" + groupName,
		Color:       &ROLE_COLOR,
		Hoist:       &ROLE_HOIST,
		Permissions: &ROLE_PERMISSIONS,
		Mentionable: &ROLE_MENTIONNABLE,
	})

	if err != nil {
		SendErrorMessage(s, i, error.Error())
		return
	}

	var wg sync.WaitGroup

	err = createCategory(s, i, &wg, role.ID, groupName, description)

	if err != nil {
		SendErrorMessage(s, i, error.Error())
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role.ID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		error = ReferenceRoleInChannel(
			s,
			i,
			groupName,
			description,
			role.ID,
		)
	}()

	wg.Wait()

	// Sucess message
	SendSuccessMessage(s, i, groupName)

	verbosity.Debug("User", member.User.ID, ", created group '", i.ApplicationCommandData().Options[0].Value, "' in guild ", guild)
}

var Summary = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
