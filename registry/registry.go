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

func NewDockerProjectRegistry(cfg *config.Config) *DockerProjectRegistry {
	dc := docker.NewDockerCmd(cfg)

	return &DockerProjectRegistry{
		dockerCmd: dc,
		config:    cfg,
	}
}

func (reg *DockerProjectRegistry) Config() *config.Config {
	return reg.config
}

func (reg *DockerProjectRegistry) ProjectExists(p *project.Project) (bool, error) {
	includeStopped := true
	projects, err := reg.fetchProjects(includeStopped)
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

func (reg *DockerProjectRegistry) StartProject(p *project.Project, recreate, update bool) error {
	if err := reg.stopOtherActiveProjects(p); err != nil {
		return err
	}

	if reg.config.AwsLogin {
		logger.Info("Logging into AWS registry")
		if err := reg.dockerCmd.LoginAws(); err != nil {
			return err
		}
	}

	logger.Info("Starting", p.String())
	dc := reg.dockerCmd.CreateAndStartProjectCommand(p, recreate, update)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) stopOtherActiveProjects(p *project.Project) error {
	logger.Info("Stopping other active projects")
	includeStopped := false
	activeProjects, err := reg.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	for _, ap := range activeProjects {
		if ap.Name == p.Name {
			continue
		}
		logger.Debug("Stopping %s", ap.String())
		dc := reg.dockerCmd.StopProjectCommand(ap)
		if err := dc.Execute(); err != nil {
			logger.Warning(fmt.Sprintf("Could not stop %s", ap.String()), err)
		}
	}

	return nil
}

func (reg *DockerProjectRegistry) StopProject(p *project.Project) error {
	logger.Info("Stopping", p.String())

	exists, err := reg.ProjectExists(p)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := reg.dockerCmd.StopProjectCommand(p)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) RestartProject(p *project.Project) error {
	logger.Info("Restarting", p.String())

	exists, err := reg.ProjectExists(p)
	if err != nil {
		return err
	}

	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := reg.dockerCmd.RestartProjectCommand(p)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) RemoveProject(p *project.Project) error {
	logger.Info("Removing", p.String())

	exists, err := reg.ProjectExists(p)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warning("%s does not exist", helpers.ToTitle(p.String()))
		return nil
	}

	dc := reg.dockerCmd.RemoveProjectCommand(p)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) BuildProject(p *project.Project, noCache bool) error {
	logger.Info("Building", p.String())
	dc := reg.dockerCmd.BuildProjectCommand(p, noCache)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) Cleanup() error {
	logger.Info("Cleaning up")
	includeStopped := true
	projects, err := reg.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	isErr := false
	for _, p := range projects {
		dc := reg.dockerCmd.RemoveProjectCommand(p)
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

func (reg *DockerProjectRegistry) Terminal(p *project.Project, cmd string) error {
	logger.Info("Running terminal for", p.String())

	if !p.IsServiceDefined() {
		p.SetServiceName(reg.config.TerminalService)
	}

	if cmd == "" {
		cmd = reg.config.TerminalCommand
	}

	dc := reg.dockerCmd.TerminalCommand(p, cmd)

	return dc.Execute()
}
