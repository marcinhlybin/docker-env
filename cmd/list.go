package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "ls",
	Usage:       "list projects",
	ArgsUsage:   "[PROJECT_NAME]",
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
	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	verbose := c.Bool("verbose")

	logger.ShowCommands(false)

	return reg.ListProjects(verbose)
}
