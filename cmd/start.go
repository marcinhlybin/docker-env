package cmd

import (
	"fmt"

	"github.com/marcinhlybin/docker-env/app"
	"github.com/marcinhlybin/docker-env/hooks"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"s", "up"},
	Usage:   "Start docker containers",
	Description: `Start docker containers.
If project name is not specified, master branch is used.
If project does not exist it will be created.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "set a project name",
		},
		&cli.BoolFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "use current git branch as project name",
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

	// Mutually exclusive flags -p and -b
	if c.IsSet("project") && c.IsSet("branch") {
		return fmt.Errorf("flags -p and -b are mutually exclusive")
	}

	// Get the app context
	ctx, err := app.NewAppContext(c)
	if err != nil {
		return err
	}

	// Create a new project
	p, err := ctx.CreateProject()
	if err != nil {
		return err
	}

	// Set logger
	logger.SetPrefix(p.Name)

	// Stop other active projects
	if err := ctx.Registry.StopOtherActiveProjects(p); err != nil {
		return err
	}

	// Pre-start hooks
	ctx.Registry.UpdateProjectStatus(p)
	withHooks := !c.Bool("no-hooks")
	if withHooks && !p.IsRunning() {
		if err := hooks.RunPreStartHooks(p, ctx); err != nil {
			return err
		}
	}

	// Start the project
	recreate := c.Bool("recreate")
	update := c.Bool("update")
	if err := ctx.Registry.StartProject(p, recreate, update); err != nil {
		return err
	}

	// Post-start hooks
	if withHooks && !p.IsRunning() {
		if err := hooks.RunPostStartHooks(p, ctx); err != nil {
			return err
		}
	}

	return nil
}
