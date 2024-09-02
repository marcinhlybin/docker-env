package registry

import (
	"fmt"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type DockerProjectRegistry struct {
	config    *config.Config
	dockerCmd *docker.DockerCmd
}

func NewDockerProjectRegistry(cfg *config.Config) project.ProjectRegistry {
	dc := docker.NewDockerCmd(cfg)

	return &DockerProjectRegistry{
		dockerCmd: dc,
		config:    cfg,
	}
}

func (registry *DockerProjectRegistry) Config() *config.Config {
	return registry.config
}

func (registry *DockerProjectRegistry) ProjectExists(p *project.Project) (bool, error) {
	includeStopped := true
	projects, err := registry.fetchProjects(includeStopped)
	if err != nil {
		return false, err
	}

	for _, proj := range projects {
		if proj.Name == p.Name {
			return true, nil
		}
	}

	return false, nil
}

func (registry *DockerProjectRegistry) StartProject(p *project.Project, recreate, update bool) error {
	if err := registry.stopOtherActiveProjects(p); err != nil {
		return err
	}

	if registry.config.AwsLogin {
		logger.Info("Logging into AWS registry")
		if err := registry.dockerCmd.LoginAws(); err != nil {
			return err
		}
	}

	logger.Info("Starting", p.String())
	dc := registry.dockerCmd.CreateAndStartProjectCommand(p, recreate, update)
	return dc.Execute()
}

func (registry *DockerProjectRegistry) stopOtherActiveProjects(p *project.Project) error {
	logger.Info("Stopping other active projects")
	includeStopped := false
	activeProjects, err := registry.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	for _, ap := range activeProjects {
		if ap.Name == p.Name {
			continue
		}
		logger.Debug("Stopping %s", ap.String())
		dc := registry.dockerCmd.StopProjectCommand(ap)
		if err := dc.Execute(); err != nil {
			logger.Warning(fmt.Sprintf("Could not stop %s", ap.String()), err)
		}
	}

	return nil
}

func (registry *DockerProjectRegistry) StopProject(p *project.Project) error {
	logger.Info("Stopping", p.String())

	exists, err := registry.ProjectExists(p)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := registry.dockerCmd.StopProjectCommand(p)
	return dc.Execute()
}

func (registry *DockerProjectRegistry) RestartProject(p *project.Project) error {
	logger.Info("Restarting", p.String())

	exists, err := registry.ProjectExists(p)
	if err != nil {
		return err
	}

	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := registry.dockerCmd.RestartProjectCommand(p)
	return dc.Execute()
}

func (registry *DockerProjectRegistry) RemoveProject(p *project.Project) error {
	logger.Info("Removing", p.String())

	exists, err := registry.ProjectExists(p)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := registry.dockerCmd.RemoveProjectCommand(p)
	return dc.Execute()
}

func (registry *DockerProjectRegistry) BuildProject(p *project.Project, noCache bool) error {
	logger.Info("Building", p.String())
	dc := registry.dockerCmd.BuildProjectCommand(p, noCache)
	return dc.Execute()
}

func (registry *DockerProjectRegistry) Cleanup() error {
	logger.Info("Cleaning up")
	includeStopped := true
	projects, err := registry.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	isErr := false
	for _, p := range projects {
		dc := registry.dockerCmd.RemoveProjectCommand(p)
		err := dc.Execute()
		if err != nil {
			isErr = true
			logger.Warning("Could not remove %s: %v", p.String(), err)
		}
	}

	if isErr {
		return fmt.Errorf("one or more projects could not be removed")
	}

	return nil
}

func (registry *DockerProjectRegistry) RunTerminal(p *project.Project) error {
	logger.Info("Running terminal for", p.String())
	dc := registry.dockerCmd.RunTerminalCommand(p, registry.config.TerminalCommand)
	return dc.Execute()
}
