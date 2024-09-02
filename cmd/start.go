package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StartCommand = cli.Command{
	Name:      "start",
	Usage:     "Start docker containers",
	ArgsUsage: "[PROJECT_NAME]",
	Description: `Start docker containers.
If project name is not specified, current branch name is used.
If project does not exist it will be created.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "start a single service",
		},
		&cli.BoolFlag{
			Name:    "recreate",
			Aliases: []string{"r"},
			Usage:   "recreate the containers",
		},
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update the images and recreate the containers",
		},
	},
	Action: startAction,
}

func startAction(c *cli.Context) error {
	p, err := NewProject(c)
	if err != nil {
		return err
	}

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	recreate := c.Bool("recreate")
	update := c.Bool("update")

	logger.SetPrefix(p.Name)

	return reg.StartProject(p, recreate, update)
}
