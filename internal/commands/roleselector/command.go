package roleselector

import (
	"fmt"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

var perms = int64(discordgo.PermissionManageRoles)

var command = &discordgo.ApplicationCommand{
	Name:                     "add_role",
	Description:              "Add a role to be selected",
	DefaultMemberPermissions: &perms,
	Type:                     discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "role_name",
			Required:    true,
			Description: "The name of the role to be selected",
		},
	},
}

func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	verbosity.Debug(i.Type)

	verbosity.Debug(fmt.Sprintf("%+v", i.ApplicationCommandData().Options[0].Value))

	verbosity.Error(s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "role",

			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Add",
							CustomID: "add",
						},
					},
				},
			},
		},
	}))
}

var Summary = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
