package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var TerminalCommand = cli.Command{
	Name:        "terminal",
	Aliases:     []string{"term", "shell", "ssh"},
	ArgsUsage:   "[COMMAND]",
	Usage:       "Run terminal",
	Description: `Run terminal in the project.`,
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
	Action: terminalAction,
}

func terminalAction(c *cli.Context) error {
	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(ctx.Project.Name)

	cmd := c.Args().Slice()

	return ctx.Registry.Terminal(ctx.Project, cmd)
}
