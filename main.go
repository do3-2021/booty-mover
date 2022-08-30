package main

import (
	"os"
	"path"

	"github.com/NilsPonsard/verbosity"
	cli "github.com/jawher/mow.cli"
	"github.com/do3-2021/booty-mover/internal/commands"
	"github.com/do3-2021/booty-mover/pkg/files"
)

// Version will be set by the script build.sh
var version string

func main() {

	app := cli.App("booty-mover", "starter project")
	app.Version("v version", version)

	defaultPath := files.ParsePath("~/.booty-mover/")

	// arguments

	var (
		verbose     = app.BoolOpt("d debug", false, "Debug mode, more verbose operations")
		appPath     = app.StringOpt("c config", path.Join(defaultPath, "config.json"), "Path to the config file")
		disableLogs = app.BoolOpt("disable-logs", false, "Disable the saving of logs")
	)

	// Executed befor the commands

	app.Before = func() {

		parsedConfigPath := *appPath
		files.EnsureFolder(parsedConfigPath)

		// create the folder for the logs

		files.EnsureFolder(path.Join(defaultPath, "test"))

		// Configure the logs

		verbosity.SetupLog(*verbose, path.Join(defaultPath, "logs.txt"), version)

		verbosity.SetLogging(!*disableLogs)

	}

	// set subcommands

	commands.SetupCommands(app)

	// parse the args

	app.Run(os.Args)
}
