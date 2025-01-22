package cmd

import (
	"github.com/marcinhlybin/docker-env/app"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var RestartCommand = cli.Command{
	Name:    "restart",
	Aliases: []string{"r", "reboot"},
	Usage:   "Restart docker containers",
	Description: `Restart docker containers.
If project name is not specified, master branch is used.`,
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

	ctx, err := app.NewAppContext(c)
	if err != nil {
		return err
	}

	p, err := ctx.ActiveProject()
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	logger.SetPrefix(p.Name)

	reg := ctx.Registry
	return reg.RestartProject(p)
}
