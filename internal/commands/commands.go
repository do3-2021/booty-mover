package commands

import (
	cli "github.com/jawher/mow.cli"
	"github.com/do3-2021/booty-mover/internal/commands/ping"
)

// configure subcommands
func SetupCommands(app *cli.Cli) {
	app.Command("ping", "ping", ping.Ping)
}
