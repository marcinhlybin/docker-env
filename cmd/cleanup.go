package cmd

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/urfave/cli/v2"
)

var CleanupCommand = cli.Command{
	Name:        "cleanup",
	Usage:       "Cleanup entire project",
	Description: `Remove everything related to the project.`,
	Action:      cleanupAction,
}

func cleanupAction(c *cli.Context) error {
	ExitWithErrorOnArgs(c)

	reg, err := NewRegistry(c)
	if err != nil {
		return err
	}

	logger.SetPrefix(reg.Config().Project)

	return reg.Cleanup()
}
