package main

import (
	"os"
	"os/signal"

	"github.com/NilsPonsard/verbosity"
	"github.com/do3-2021/booty-mover/internal/bot"
	"github.com/do3-2021/booty-mover/internal/commands"
	"github.com/do3-2021/booty-mover/internal/database"
)

// Version will be set by the script build.sh
var version string

func main() {
	verbosity.SetupLog(true, "", version)

	session, err := bot.Connect()

	if err != nil {
		verbosity.Fatal(err)
	}

	defer session.Close()

	db, err := database.ConnectPostgres()

	if err != nil {
		verbosity.Fatal("Can’t connect to db : ", err)
	}
	defer db.Close()

	commands.Configure(session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	verbosity.Info("Press Ctrl+C to exit")
	<-stop
}
