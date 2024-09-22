package cmd

import (
	"github.com/marcinhlybin/docker-env/addons"
	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/marcinhlybin/docker-env/registry"
	"github.com/urfave/cli/v2"
)

type AppContext struct {
	Config   *config.Config
	Project  *project.Project
	Registry *registry.DockerProjectRegistry
}

func NewAppContext(c *cli.Context) (*AppContext, error) {
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

	return &AppContext{
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

func (ctx *AppContext) PreStartHook() error {
	hook := addons.NewPreStartHook(ctx.Config.PreStartHook, ctx.Project.Name, ctx.Project.ServiceName)
	return hook.Run()
}

func (ctx *AppContext) PostStartHook() error {
	hook := addons.NewPostStartHook(ctx.Config.PostStartHook, ctx.Project.Name, ctx.Project.ServiceName)
	return hook.Run()
}

func (ctx *AppContext) PostStopHook() error {
	hook := addons.NewPostStopHook(ctx.Config.PostStopHook, ctx.Project.Name, ctx.Project.ServiceName)
	return hook.Run()
}
