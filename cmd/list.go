package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list", "l", "ll"},
	Usage:       "List projects, 'll' to show containers.",
	Description: `List docker projects and containers.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "containers",
			Aliases: []string{"c"},
			Usage:   "show containers",
		},
		&cli.BoolFlag{
			Name:    "running",
			Aliases: []string{"r"},
			Usage:   "show only running projects",
		},
	},
	Action: listAction,
}

func listAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	reg, cfg := app.Registry, app.Config
	logger.SetPrefix(cfg.ComposeProjectName)
	logger.ShowExecutedCommands(false)

	containers := c.Bool("containers") || isAliasUsed("ll")
	includeStopped := !c.Bool("running")

	if containers {
		return reg.ListContainers(includeStopped)
	}

	return reg.ListProjects(includeStopped)
}
