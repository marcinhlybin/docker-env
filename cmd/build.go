package cmd

import (
	"github.com/urfave/cli/v2"
)

var BuildCommand = cli.Command{
	Name:      "build",
	Usage:     "Build docker images",
	ArgsUsage: "[PROJECT_NAME]",
	Description: `Build docker images.
If environment name is not specified current branch name is used.`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "start a single service",
		},
		&cli.BoolFlag{
			Name:  "no-cache",
			Usage: "do not use cache when building the image",
		},
	},
	Action: buildAction,
}

func buildAction(c *cli.Context) error {
	p, err := NewProject(c)
	if err != nil {
		return err
	}

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	noCache := c.Bool("no-cache")

	return reg.BuildProject(p, noCache)
}
