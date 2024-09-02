package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var RestartCommand = cli.Command{
	Name:      "restart",
	Usage:     "Restart docker containers",
	ArgsUsage: "[PROJECT_NAME]",
	Description: `Restart docker containers.
If environment name is not specified current branch name is used.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "restart a single service",
		},
	},
	Action: restartAction,
}

func restartAction(c *cli.Context) error {
	p, err := NewProject(c)
	if err != nil {
		return err
	}

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(p.Name)

	return reg.RestartProject(p)
}
