package main

import (
	"os"

	"github.com/marcinhlybin/docker-env/cmd"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "docker-env",
		Usage:   "Docker environments manager",
		Version: version.Version,
		Description: `All commands must run in the git repository directory of the project.
If environment name is not specified current branch name is used.`,
		Commands: []*cli.Command{
			&cmd.StartCommand,
			&cmd.StopCommand,
			&cmd.RestartCommand,
			&cmd.RemoveCommand,
			&cmd.ListCommand,
			&cmd.CleanupCommand,
			&cmd.BuildCommand,
			&cmd.InfoCommand,
			&cmd.TerminalCommand,
			&cmd.CodeCommand,
			&cmd.VersionCommand,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file path",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug mode",
			},
		},
		Before: func(c *cli.Context) error {
			// Show help if no arguments
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			// Enable debug mode
			if c.Bool("debug") {
				logger.SetDebug(true)
			}

			return nil
		},
	}

	// Version
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "show version string, alias for 'version --short'",
	}
	cli.VersionPrinter = func(c *cli.Context) {
		version.PrintVersionString()
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error("%s", helpers.ToTitle(err.Error()))
		os.Exit(1)
	}
}
