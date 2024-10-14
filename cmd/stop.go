package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StopCommand = cli.Command{
	Name:    "stop",
	Aliases: []string{"ss", "down"},
	Usage:   "Stop docker containers",
	Description: `Stop docker containers.
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
			Usage:   "stop a single service",
		},
		&cli.BoolFlag{
			Name:    "no-hooks",
			Aliases: []string{"without-hooks"},
			Usage:   "do not run pre/post start hooks",
		},
	},
	Action: stopAction,
}

func stopAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	p, reg := app.Project, app.Registry
	logger.SetPrefix(p.Name)

	// Stop the project
	if err := reg.StopProject(p); err != nil {
		return err
	}

	// Post-stop hooks
	withHooks := !c.Bool("no-hooks")
	if withHooks {
		if err := app.RunPostStopHook(); err != nil {
			return err
		}
	}

	return nil
}
