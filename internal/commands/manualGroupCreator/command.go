// manually configure a group
package manualgroupcreator

import (
	"fmt"
	"strings"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	"github.com/do3-2021/booty-mover/internal/commands/groupcreator"
)

var command = &discordgo.ApplicationCommand{
	Name:        "manual-create-group",
	Description: "Create a new group from an existing role",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Description: "The role of the group",
			Required:    true,
		},
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

// manually configure a group from an existing role
func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// get args

	roleID := ""
	description := ""
	groupName := ""
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "description":
			description = option.StringValue()
		case "name":
			groupName = strings.Replace(option.StringValue(), " ", "-", -1)
		case "role":
			roleID = option.RoleValue(nil, "").ID
		}
	}

	// create message

	err := groupcreator.ReferenceRoleInChannel(s, i, groupName, description, roleID)

	if err != nil {
		verbosity.Error(err)
		groupcreator.SendErrorMessage(s, i, fmt.Sprint(err))
		return
	}

	// send message
	groupcreator.SendSuccessMessage(s, i, groupName)

}

var Descriptor = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
