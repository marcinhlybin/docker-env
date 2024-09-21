package cmd

import (
	"github.com/marcinhlybin/docker-env/addons"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var StopCommand = cli.Command{
	Name:    "stop",
	Aliases: []string{"ss", "down"},
	Usage:   "Stop docker containers",
	Description: `Stop docker containers.
If environment name is not specified current branch name is used.`,
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
	},
	Action: stopAction,
}

func stopAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	ctx, err := NewAppContext(c)
	if err != nil {
		return err
	}
	logger.SetPrefix(ctx.Project.Name)

	if err := ctx.Registry.StopProject(ctx.Project); err != nil {
		return err
	}

	// Run post-stop script
	return addons.RunScript("post-stop", ctx.Config.PostStopScript)

}
