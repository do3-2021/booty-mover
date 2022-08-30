package groupcreator

import (
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


//This command will store The User, give him a role that grants him acess to a fresh new channel
//
func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//user := i.Member

	grpNme := strings.Replace(i.ApplicationCommandData().Options[0].StringValue(), " ", "-", -1)

	s.GuildChannelCreate(
		i.GuildID, 
		grpNme,
		discordgo.ChannelTypeGuildText,
	)

	verbosity.Debug(i.ApplicationCommandData().Options[0].Value)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Group ",
		},
	})
}

var Summary = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
