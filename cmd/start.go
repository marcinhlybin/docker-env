package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"s", "up"},
	Usage:   "Start docker containers",
	Description: `Start docker containers.
If project name is not specified, current branch name is used.
If project does not exist it will be created.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "set a project name",
		},
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "start a single service",
		},
		&cli.BoolFlag{
			Name:    "recreate",
			Aliases: []string{"r"},
			Usage:   "recreate the containers",
		},
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update the images and recreate the containers",
		},
	},
	Action: startAction,
}

func startAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}

	recreate := c.Bool("recreate")
	update := c.Bool("update")

	logger.SetPrefix(ctx.Project.Name)

	// Run pre-start script
	if err := ctx.PreStartHook(); err != nil {
		return err
	}

	if err := ctx.Registry.StartProject(ctx.Project, recreate, update); err != nil {
		return err
	}

	// Run post-start script
	return ctx.PostStartHook()
}
