package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	"github.com/do3-2021/booty-mover/internal/commands/ping"
	"github.com/do3-2021/booty-mover/internal/commands/roleselector"
)

var descriptors = []common.CommandDescriptor{
	ping.Summary,
	roleselector.Summary,
}

func Configure(session *discordgo.Session) (commandsHandler *CommandsHandler) {

	commandsHandler = New(descriptors)
	commandsHandler.Register(session)
	session.AddHandler(commandsHandler.Handle)

	return
}
