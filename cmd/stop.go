package cmd

import (
	"github.com/marcinhlybin/docker-env/app"
	"github.com/marcinhlybin/docker-env/hooks"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StopCommand = cli.Command{
	Name:    "stop",
	Aliases: []string{"ss", "down"},
	Usage:   "Stop docker containers",
	Description: `Stop docker containers.
If project name is not specified master branch is used.`,
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

	// Stop the project
	if err := ctx.Registry.StopProject(p); err != nil {
		return err
	}

	// Post-stop hooks
	withHooks := !c.Bool("no-hooks")
	if withHooks {
		if err := hooks.RunPostStopHooks(p, ctx); err != nil {
			return err
		}
	}

	return nil
}
