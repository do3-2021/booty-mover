package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	configurechannel "github.com/do3-2021/booty-mover/internal/commands/configure-channel"
	"github.com/do3-2021/booty-mover/internal/commands/groupcreator"
	"github.com/do3-2021/booty-mover/internal/commands/ping"
	"github.com/do3-2021/booty-mover/internal/commands/roleselector"
	"github.com/do3-2021/booty-mover/internal/features/voicecreation"
)

// Put your commands descriptors here
var descriptors = []common.CommandDescriptor{
	ping.Summary,
	roleselector.Summary,
	configurechannel.Descriptor,
	groupcreator.Summary,
}

// Configure the bot for commands
func Configure(session *discordgo.Session) (commandsHandler *CommandsHandler) {

	commandsHandler = New(descriptors)
	commandsHandler.Register(session)
	session.AddHandler(commandsHandler.Handle)
	session.AddHandler(voicecreation.Handle)

	return
}
