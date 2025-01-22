package app

import (
	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/registry"
	"github.com/urfave/cli/v2"
)

type AppContext struct {
	ProjectName      string
	ServiceName      string
	IsBranch         bool
	IsProjectNameSet bool
	Config           *config.Config
	Registry         *registry.DockerProjectRegistry
}

func NewAppContext(c *cli.Context) (*AppContext, error) {
	cfg, err := initializeConfig(c)
	if err != nil {
		return nil, err
	}

	reg, err := initializeRegistry(cfg)
	if err != nil {
		return nil, err
	}

	// Set default project name
	projectName := c.String("project")
	if projectName == "" {
		projectName = cfg.GitDefaultBranch
	}

	return &AppContext{
		ProjectName:      projectName,
		ServiceName:      c.String("service"),
		IsProjectNameSet: c.IsSet("project"),
		IsBranch:         c.Bool("branch"),
		Config:           cfg,
		Registry:         reg,
	}, nil
}

func initializeConfig(c *cli.Context) (*config.Config, error) {
	cfg := config.NewConfig()
	if err := cfg.LoadConfig(c.String("config")); err != nil {
		return nil, err
	}

	// Show executed commands
	logger.ShowExecutedCommands(cfg.ShowExecutedCommands)

	return cfg, nil
}

func initializeRegistry(cfg *config.Config) (*registry.DockerProjectRegistry, error) {
	return registry.NewDockerProjectRegistry(cfg), nil
}
