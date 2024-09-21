package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var RestartCommand = cli.Command{
	Name:    "restart",
	Aliases: []string{"r", "reboot"},
	Usage:   "Restart docker containers",
	Description: `Restart docker containers.
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
			Usage:   "restart a single service",
		},
	},
	Action: restartAction,
}

func restartAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(ctx.Project.Name)

	return ctx.Registry.RestartProject(ctx.Project)
}
