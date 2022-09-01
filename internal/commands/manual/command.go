// manually configure a group
package manual

import (
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var command = &discordgo.ApplicationCommand{
	Name:        "manual-create-group",
	Description: "Create a new group from an existing role",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Description: "The role of the group",
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
			Required:    false,
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

}

var Descriptor = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
