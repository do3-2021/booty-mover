package commands

import (
	"fmt"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

type CommandsHandler struct {
	descriptors        map[string]*common.CommandDescriptor
	registeredCommands []*discordgo.ApplicationCommand
}

// Create a new command handler with given command descriptors
func New(summaries []common.CommandDescriptor) *CommandsHandler {
	summariesMap := make(map[string]*common.CommandDescriptor)

	for _, summary := range summaries {
		summariesMap[summary.Command.Name] = &summary
	}

	return &CommandsHandler{
		descriptors: summariesMap,
	}
}

// Register every commands to the session
func (h *CommandsHandler) Register(session *discordgo.Session) (registeredCommands []*discordgo.ApplicationCommand) {

	for _, summary := range h.descriptors {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", summary.Command)
		if err != nil {
			verbosity.Fatal(err)

		}
		verbosity.Debug("Registered command:", cmd.Name)

		h.registeredCommands = append(h.registeredCommands, cmd)

	}

	return
}

// Handle a discord interaction event and dispatch them to the correspondin descriptor
func (h *CommandsHandler) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:

	case discordgo.InteractionMessageComponent:
	default:
		verbosity.Debug("Unhandled interaction type:", interaction.Type)
	}

	if interaction.Type == discordgo.InteractionMessageComponent {
		verbosity.Debug(fmt.Sprintf("message : %+v", interaction.MessageComponentData()))
		return
	}

	name := interaction.ApplicationCommandData().Name

	summary, ok := h.descriptors[name]
	if !ok {
		verbosity.Debug("Unknown command:", name)
		return
	}
	summary.Execute(session, interaction)

}
