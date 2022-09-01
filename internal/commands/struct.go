package commands

import (
	"strings"

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

	for _, descriptor := range h.descriptors {
		verbosity.Debug("Registering command:", descriptor.Command.Name)
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", descriptor.Command)
		if err != nil {
			verbosity.Fatal(err)

		}
		verbosity.Debug("Registered command:", cmd.Name)

		h.registeredCommands = append(h.registeredCommands, cmd)

	}

	return
}

// Handle a discord interaction event and dispatch them to the correspondin descriptor
// InteractionApplicationCommand are found by name
// InteractionMessageComponent are matched with the prefix of the custom id being the command name
func (h *CommandsHandler) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	switch interaction.Type {

	// slash commands
	case discordgo.InteractionApplicationCommand:
		name := interaction.ApplicationCommandData().Name

		descriptor, ok := h.descriptors[name]

		if !ok {
			verbosity.Error("No command found for name:", name)
			return
		}

		descriptor.Execute(session, interaction)

	// interactions on a message
	case discordgo.InteractionMessageComponent:

		customId := interaction.MessageComponentData().CustomID

		// search for the command corresponding to the custom id

		for key, descriptor := range h.descriptors {
			if strings.HasPrefix(customId, key) {
				descriptor.Execute(session, interaction)
				return
			}
		}

	default:
		verbosity.Debug("Unhandled interaction type:", interaction.Type)
	}

}
