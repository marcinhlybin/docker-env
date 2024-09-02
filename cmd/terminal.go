package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var TerminalCommand = cli.Command{
	Name:        "terminal",
	Aliases:     []string{"term", "shell", "ssh"},
	Usage:       "Run terminal",
	ArgsUsage:   "[PROJECT_NAME]",
	Description: `Run terminal in the project.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "start a single service",
		},
		&cli.StringFlag{
			Name:    "cmd",
			Aliases: []string{"c"},
			Usage:   "run command in the terminal",
		},
	},
	Action: terminalAction,
}

func terminalAction(c *cli.Context) error {
	p, err := NewProject(c)
	if err != nil {
		return err
	}

	println(p.String())
	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(p.Name)

	cmd := c.String("cmd")

	return reg.Terminal(p, cmd)
}
