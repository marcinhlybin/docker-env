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
	p, err := NewProject(c)
	if err != nil {
		return err
	}

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(p.Name)

	cmd := c.Args().First()

	return reg.Terminal(p, cmd)
}
