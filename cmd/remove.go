package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var RemoveCommand = cli.Command{
	Name:      "remove",
	Usage:     "Remove docker containers",
	ArgsUsage: "[PROJECT_NAME]",
	Description: `Remove docker containers.
If environment name is not specified current branch name is used.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "remove a single service",
		},
	},
	Action: removeAction,
}

func removeAction(c *cli.Context) error {
	p, err := initializeProject(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(p.Name)

	return p.Remove()
}
