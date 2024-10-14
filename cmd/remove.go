package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var RemoveCommand = cli.Command{
	Name:    "remove",
	Aliases: []string{"rm", "delete"},
	Usage:   "Remove docker containers",
	Description: `Remove docker containers.
If environment name is not specified current branch name is used.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "set a project name",
		},
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "remove a single service",
		},
	},
	Action: removeAction,
}

func removeAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	p, reg := app.Project, app.Registry
	logger.SetPrefix(p.Name)

	return reg.RemoveProject(p)
}
