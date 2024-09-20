package cmd

import (
	"os"
	"strings"

	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list", "l", "ll"},
	Usage:       "List projects. Use 'll' to show containers.",
	Description: `List docker projects and containers.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "show containers",
		},
	},
	Action: listAction,
}

func listAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	verbose := c.Bool("verbose") || isAliasUsed("ll")

	logger.SetPrefix(reg.Config().ComposeProjectName)
	// logger.ShowCommands(false)

	return reg.ListProjects(verbose)
}

func isAliasUsed(alias string) bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, alias) {
			return true
		}
	}
	return false
}
