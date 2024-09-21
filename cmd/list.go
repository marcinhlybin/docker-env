package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list", "l", "ll"},
	Usage:       "List projects. Use 'll' to show containers.",
	Description: `List docker projects and containers.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "show containers",
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

	verbose := c.Bool("verbose") || isAliasUsed("ll")

	logger.SetPrefix(ctx.Config.ComposeProjectName)

	return ctx.Registry.ListProjects(verbose)
}
