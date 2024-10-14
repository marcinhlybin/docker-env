package cmd

import (
	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var LogsCommand = cli.Command{
	Name:        "logs",
	Aliases:     []string{"log"},
	Usage:       "Show container logs",
	Description: `Show container logs. If a service name is not specified, all logs are shown.`,
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
			Name:    "follow",
			Aliases: []string{"f"},
			Usage:   "follow log output",
		},
		&cli.BoolFlag{
			Name:    "timestamps",
			Aliases: []string{"t"},
			Usage:   "show timestamps",
		},
	},
	Action: logsAction,
}

func logsAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	app, err := NewApp(c)
	if err != nil {
		return err
	}

	p, reg := app.Project, app.Registry
	logger.SetPrefix(p.Name)

	opts := docker.LogsOptions{
		FollowOutput:   c.Bool("follow"),
		ShowTimestamps: c.Bool("timestamps"),
	}

	// Show logs
	if err := reg.Logs(p, opts); err != nil {
		return err
	}

	return nil
}
