package cmd

import (
	"github.com/marcinhlybin/docker-env/app"
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

	ctx, err := app.NewAppContext(c)
	if err != nil {
		return err
	}

	cfg := ctx.Config

	return cfg.ShowConfig()
}
