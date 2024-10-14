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
		&cli.BoolFlag{
			Name:    "no-hooks",
			Aliases: []string{"without-hooks"},
			Usage:   "do not run pre/post start hooks",
		},
	},
	Action: startAction,
}

func startAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	p, reg := app.Project, app.Registry
	reg.UpdateProjectStatus(p)

	recreate := c.Bool("recreate")
	update := c.Bool("update")

	logger.SetPrefix(p.Name)

	if err := reg.StopOtherActiveProjects(p); err != nil {
		return err
	}

	// Pre-start hooks
	withHooks := !c.Bool("no-hooks")
	if withHooks && !p.IsRunning() {
		if err := app.RunPreStartHook(); err != nil {
			return err
		}
	}

	// Start the project
	if err := reg.StartProject(p, recreate, update); err != nil {
		return err
	}

	// Post-start hooks
	if withHooks && !p.IsRunning() {
		if err := app.RunPostStartHook(); err != nil {
			return err
		}
	}

	return nil
}
