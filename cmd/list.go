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
			Name:  "running",
			Usage: "show only running projects",
		},
	},
	Action: listAction,
}

func listAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(ctx.Config.ComposeProjectName)

	showContainers := c.Bool("containers") || isAliasUsed("ll")
	includeStopped := !c.Bool("running")

	return ctx.Registry.ListProjects(includeStopped, showContainers)
}
