package cmd

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func ExitWithErrorOnArgs(c *cli.Context) {
	if c.Args().Present() {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}
}

func IsAliasUsed(alias string) bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, alias) {
			return true
		}
	}
	return false
}
