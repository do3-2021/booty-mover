package groupcreator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/NilsPonsard/verbosity"

	"github.com/bwmarrin/discordgo"
	"github.com/do3-2021/booty-mover/internal/commands/common"
	configurechannel "github.com/do3-2021/booty-mover/internal/commands/configure-channel"
	"github.com/do3-2021/booty-mover/internal/database"
	"github.com/do3-2021/booty-mover/internal/guild"
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

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func sendErrorMessage(s *discordgo.Session, i *discordgo.InteractionCreate, error string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("❌ Error: %v", error),
		},
	})
}

func referenceRoleInChannel(s *discordgo.Session, i *discordgo.InteractionCreate, group string, description string, roleID string) (error error) {
	db, error := database.GetDB()

	if error != nil {
		verbosity.Error(error.Error())
		return errors.New("could not contact database to get the role listing channel's ID.\nCancelling Group's creation")
	}

	channel, error := guild.GetGroupChannel(db, i.GuildID)
	if error != nil {
		verbosity.Error(error.Error())
		return errors.New("could not find the role listing channel's.\nDid you run the " + configurechannel.Descriptor.Command.Name + " command?")
	}

	s.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
		Content: fmt.Sprintf("%v: %v", group, description),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Join now!",
						CustomID: "create-group-" + roleID,
					},
				},
			},
		},
	})

	return
}

var ROLE_PERMISSIONS int64 = 0x0000000000000040 | // ADD_REACTIONS
	0x0000000000000200 | // STREAM
	0x0000000000000400 | // VIEW_CHANNEL
	0x0000000000000800 | // SEND_MESSAGES
	0x0000000000004000 | // EMBED_LINKS
	0x0000000000008000 | // ATTACH_FILES
	0x0000000000010000 | // READ_MESSAGE_HISTORY
	0x0000000000040000 | // USE_EXTERNAL_EMOJI
	0x0000000000100000 | // CONNECT
	0x0000000000200000 | // SPEAK
	0x0000000002000000 | // USE_VAD
	0x0000000004000000 | // CHANGE_NICKNAME
	0x0000000800000000 | // CREATE_PUBLIC_THREADS
	0x0000001000000000 | // CREATE_PRIVATE_THREADS
	0x0000004000000000 | // SEND_MESSAGES_IN_THREADS
	0x0000002000000000 // USE_EXTERNAL_STICKERS
var ROLE_COLOR int = 3447003
var ROLE_HOIST bool = false
var ROLE_MENTIONNABLE bool = true

// This command will store The User, give him a role that grants him acess to a fresh new channel
func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	member := i.Member
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
			groupName = strings.Replace(i.ApplicationCommandData().Options[0].StringValue(), " ", "-", -1)
		}
	}

	guild, error := s.GuildChannels(i.GuildID)
	if error != nil {
		sendErrorMessage(s, i, error.Error())
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
		sendErrorMessage(s, i, fmt.Sprintf("group %v already exists", groupName))
		verbosity.Debug(fmt.Sprintf("group %v already exists", groupName))

		return
	}

	role, _ := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name:        "grp-" + groupName,
		Color:       &ROLE_COLOR,
		Hoist:       &ROLE_HOIST,
		Permissions: &ROLE_PERMISSIONS,
		Mentionable: &ROLE_MENTIONNABLE,
	})

	error = referenceRoleInChannel(
		s,
		i,
		groupName,
		description,
		role.ID,
	)

	if error != nil {
		sendErrorMessage(s, i, error.Error())
		return
	}

	category, error := s.GuildChannelCreateComplex(
		i.GuildID,
		discordgo.GuildChannelCreateData{
			Name: groupName,
			Type: discordgo.ChannelTypeGuildCategory,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    role.ID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: ROLE_PERMISSIONS,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: ROLE_PERMISSIONS,
				},
			},
		},
	)

	if error != nil {
		sendErrorMessage(s, i, error.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role.ID)
	}()

	// Creating text channel
	go func() {
		defer wg.Done()
		s.GuildChannelCreateComplex(
			i.GuildID,
			discordgo.GuildChannelCreateData{
				Name:     "txt-" + groupName,
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: category.ID,
			},
		)
	}()

	// Creating voice channel
	go func() {
		defer wg.Done()
		s.GuildChannelCreateComplex(
			i.GuildID,
			discordgo.GuildChannelCreateData{
				Name:     "🔉voc-" + groupName,
				Type:     discordgo.ChannelTypeGuildVoice,
				ParentID: category.ID,
			},
		)
	}()

	// Sucess message
	go func() {
		defer wg.Done()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("✅ Created goup '%v'! 🎉", groupName),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}()

	wg.Wait()

	verbosity.Debug("User", member.User.ID, ", created group '", i.ApplicationCommandData().Options[0].Value, "' in guild ", guild)
}

var Summary = common.CommandDescriptor{
	Command: command,
	Execute: execute,
}
