package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var CodeCommand = cli.Command{
	Name:      "code",
	Aliases:   []string{"open"},
	ArgsUsage: "[DIR]",
	Usage:     "Open code editor",
	Description: `Open code editor for the project and attach to the container.
Directory is optional. By default it will open the / directory.`,
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
	},
	Action: codeAction,
}

func codeAction(c *cli.Context) error {
	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(ctx.Project.Name)

	dir := c.Args().First()

	return ctx.Registry.Code(ctx.Project, dir)
}
