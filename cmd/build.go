package cmd

import (
	"fmt"

	"github.com/marcinhlybin/docker-env/app"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/urfave/cli/v2"
)

var BuildCommand = cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build docker images",
	Description: `Build docker images.
If project name is not specified, master branch is used.`,
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
			Name:  "no-cache",
			Usage: "do not use cache when building the image",
		},
	},
	Action: buildAction,
}

func buildAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	// Mutually exclusive flags -p and -b
	if c.IsSet("project") && c.IsSet("branch") {
		return fmt.Errorf("flags -p and -b are mutually exclusive")
	}

	ctx, err := app.NewAppContext(c)
	if err != nil {
		return err
	}

	var p *project.Project

	// Use active project name if no project or branch is specified
	if !c.IsSet("project") && !c.IsSet("branch") {
		p, err = ctx.ActiveProject()
		if err != nil {
			return err
		}
	}

	// No active project found
	if p == nil {
		p, err = ctx.CreateProject()
		if err != nil {
			return err
		}
	}

	noCache := c.Bool("no-cache")
	return ctx.Registry.BuildProject(p, noCache)
}
