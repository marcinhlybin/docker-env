package cmd

import (
	"github.com/marcinhlybin/docker-env/app"
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

	ctx, err := app.NewAppContext(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(ctx.Config.ComposeProjectName)
	logger.ShowExecutedCommands(false)

	containers := c.Bool("containers") || IsAliasUsed("ll")
	includeStopped := !c.Bool("running")

	if containers {
		return ctx.Registry.ListContainers(includeStopped)
	}

	return ctx.Registry.ListProjects(includeStopped)
}
