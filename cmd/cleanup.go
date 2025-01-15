package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var ResetCommand = cli.Command{
	Name:        "reset",
	Aliases:     []string{"cleanup"},
	Usage:       "Removes all projects",
	Description: `Removes projects and images for this repository. Only images associated with existing projects will be removed.`,
	Action:      resetAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "hard",
			Aliases: []string{"i", "include-images", "with-images"},
			Usage:   "also remove images",
		},
	},
}

func resetAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	reg, cfg := app.Registry, app.Config
	logger.SetPrefix(cfg.ComposeProjectName)

	includeImages := c.Bool("with-images")

	return reg.Cleanup(includeImages)
}
