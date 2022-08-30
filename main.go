package main

import (
	"os"
	"os/signal"

	"github.com/NilsPonsard/verbosity"
	"github.com/do3-2021/booty-mover/internal/bot"
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

	// commands.Configure(session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	verbosity.Info("Press Ctrl+C to exit")
	<-stop
}
