package cmd

import (
	"github.com/urfave/cli/v2"
)

var InfoCommand = cli.Command{
	Name:        "info",
	Usage:       "Show configuration",
	Description: `Show docker env configratuion.`,
	Action:      infoAction,
}

func infoAction(c *cli.Context) error {
	cfg, err := initializeConfig(c)
	if err != nil {
		return err
	}

	return cfg.ShowConfig()
}
