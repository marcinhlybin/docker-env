package cmd

import (
	"github.com/marcinhlybin/docker-env/addons"
	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/marcinhlybin/docker-env/registry"
	"github.com/urfave/cli/v2"
)

type App struct {
	Config   *config.Config
	Project  *project.Project
	Registry *registry.DockerProjectRegistry
}

func NewApp(c *cli.Context) (*App, error) {
	cfg, err := initializeConfig(c)
	if err != nil {
		return nil, err
	}

	p, err := initializeProject(c)
	if err != nil {
		return nil, err
	}

	reg, err := initializeRegistry(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:   cfg,
		Project:  p,
		Registry: reg,
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

func initializeProject(c *cli.Context) (*project.Project, error) {
	// Read project arguments
	projectName := c.String("project")
	serviceName := c.String("service")

	return project.NewProject(projectName, serviceName)
}

func initializeRegistry(cfg *config.Config) (*registry.DockerProjectRegistry, error) {
	return registry.NewDockerProjectRegistry(cfg), nil
}

func (app *App) RunPreStartHooks() error {
	logger.Info("Running pre-start hooks")
	p, cfg := app.Project, app.Config
	for _, path := range cfg.PreStartHooks {
		hook := addons.NewPreStartHook(path, p.Name, p.ServiceName)
		if err := hook.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) RunPostStartHooks() error {
	logger.Info("Running post-start hooks")
	p, cfg := app.Project, app.Config
	for _, path := range cfg.PostStartHooks {
		hook := addons.NewPostStartHook(path, p.Name, p.ServiceName)
		if err := hook.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) RunPostStopHooks() error {
	logger.Info("Running post-stop hooks")
	p, cfg := app.Project, app.Config
	for _, path := range cfg.PostStopHooks {
		hook := addons.NewPostStopHook(path, p.Name, p.ServiceName)
		if err := hook.Run(); err != nil {
			return err
		}
	}
	return nil
}
