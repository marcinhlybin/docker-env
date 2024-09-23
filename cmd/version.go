package cmd

import (
	"github.com/marcinhlybin/docker-env/version"
	"github.com/urfave/cli/v2"
)

var VersionCommand = cli.Command{
	Name:        "version",
	Aliases:     []string{"v"},
	Usage:       "Show version",
	Description: `Shows docker-env version, build date and commit hash.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Aliases: []string{"s"},
			Usage:   "only show the version string",
		},
		&cli.BoolFlag{
			Name:    "build-date",
			Aliases: []string{"b"},
			Usage:   "only show the build date",
		},
		&cli.BoolFlag{
			Name:    "commit-hash",
			Aliases: []string{"c"},
			Usage:   "only show the commit hash",
		},
	},
	Action: versionAction,
}

type VersionFlags map[string]func()

func versionAction(c *cli.Context) error {
	flags := VersionFlags{
		"short":       version.PrintVersionString,
		"build-date":  version.PrintBuildDateString,
		"commit-hash": version.PrintCommitHashString,
	}

	for flag, function := range flags {
		if c.Bool(flag) {
			function()
			return nil
		}
	}

	version.PrintFullVersion()
	return nil
}
