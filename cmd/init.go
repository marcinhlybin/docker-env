package cmd

import (
	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/marcinhlybin/docker-env/registry"
	"github.com/urfave/cli/v2"
)

func NewProject(c *cli.Context) (*project.Project, error) {
	// Read project arguments
	projectName := c.Args().First()
	serviceName := c.String("service")

	return project.NewProject(projectName, serviceName)
}

func NewRegistry(c *cli.Context) (*registry.DockerProjectRegistry, error) {
	cfg, err := NewConfig(c)
	if err != nil {
		return nil, err
	}

	return registry.NewDockerProjectRegistry(cfg), nil
}

func NewConfig(c *cli.Context) (*config.Config, error) {
	return config.NewConfig(c.String("config"))
}
