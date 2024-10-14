package cmd

import (
	"github.com/urfave/cli/v2"
)

var InfoCommand = cli.Command{
	Name:        "info",
	Aliases:     []string{"config", "show"},
	Usage:       "Show configuration",
	Description: `Show docker env configratuion.`,
	Action:      infoAction,
}

func infoAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	cfg := app.Config

	return cfg.ShowConfig()
}
