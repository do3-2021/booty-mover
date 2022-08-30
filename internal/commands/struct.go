package commands

import (
	"fmt"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
)

type CommandsHandler struct {
	summaries map[string]*common.Summary
}

func New(summaries []common.Summary) *CommandsHandler {
	summariesMap := make(map[string]*common.Summary)

	for _, summary := range summaries {
		summariesMap[summary.Command.Name] = &summary
	}

	return &CommandsHandler{
		summaries: summariesMap,
	}
}

func (h *CommandsHandler) Register(session *discordgo.Session) (registeredCommands []*discordgo.ApplicationCommand) {
	registeredCommands = make([]*discordgo.ApplicationCommand, len(h.summaries))

	for _, summary := range h.summaries {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", summary.Command)
		if err != nil {
			verbosity.Fatal(err)

		}
		verbosity.Debug("Registered command:", cmd.Name)

		registeredCommands = append(registeredCommands, cmd)

	}

	return
}

func (h *CommandsHandler) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	verbosity.Debug(interaction.Type)
	verbosity.Debug(interaction.ID)

	if interaction.Type == discordgo.InteractionMessageComponent {

		verbosity.Debug(fmt.Sprintf("message : %+v", interaction.MessageComponentData()))
		return
	}

	name := interaction.ApplicationCommandData().Name

	summary, ok := h.summaries[name]
	if !ok {
		verbosity.Debug("Unknown command:", name)
		return
	}
	summary.Execute(session, interaction)

}
