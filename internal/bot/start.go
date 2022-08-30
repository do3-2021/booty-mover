package bot

import (
	"errors"
	"fmt"
	"os"

	"github.com/NilsPonsard/verbosity"
	"github.com/bwmarrin/discordgo"
)

var (
	ErrNoToken = errors.New("no token in env")
)

func Connect() (session *discordgo.Session, err error) {

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		err = ErrNoToken
		return
	}

	session, err = discordgo.New("Bot " + token)
	if err != nil {
		return
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		verbosity.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	err = session.Open()

	session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "discordgo",
				Type: discordgo.ActivityTypeGame,
			},
		},
	})

	return
}
