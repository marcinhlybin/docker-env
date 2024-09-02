package cmd

import (
	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/marcinhlybin/docker-env/registry"
	"github.com/urfave/cli/v2"
)

func initializeProject(c *cli.Context) (*project.Project, error) {
	// Read project arguments
	projectName := c.Args().First()
	serviceName := c.String("service")

	registry, err := initializeProjectRegistry(c)
	if err != nil {
		return nil, err
	}

	return project.NewProject(projectName, serviceName, registry)
}

func initializeProjectRegistry(c *cli.Context) (project.ProjectRegistry, error) {
	cfg, err := initializeConfig(c)
	if err != nil {
		return nil, err
	}

	return registry.NewDockerProjectRegistry(cfg), nil

	// return project.NewProjectRegistry(registry.NewDockerProjectRegistry, cfg), nil
}

func initializeConfig(c *cli.Context) (*config.Config, error) {
	return config.NewConfig(c.String("config"))
}
