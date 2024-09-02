package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StopCommand = cli.Command{
	Name:      "stop",
	Usage:     "Stop docker containers",
	ArgsUsage: "[PROJECT_NAME]",
	Description: `Stop docker containers.
If environment name is not specified current branch name is used.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "stop a single service",
		},
	},
	Action: stopAction,
}

func stopAction(c *cli.Context) error {
	p, err := initializeProject(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(p.Name)

	return p.Stop()
}
