package registry

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type DockerProjectRegistry struct {
	Config    *config.Config
	dockerCmd *docker.DockerCmd
}

func NewDockerProjectRegistry(cfg *config.Config) *DockerProjectRegistry {
	dc := docker.NewDockerCmd(cfg)

	return &DockerProjectRegistry{
		Config:    cfg,
		dockerCmd: dc,
	}
}

func (reg *DockerProjectRegistry) UpdateProjectStatus(p *project.Project) error {
	projects, err := reg.fetchProjects(true)
	if err != nil {
		return err
	}

	for _, proj := range projects {
		if proj.Name == p.Name {
			p.SetStatus(proj.Status)
			return nil
		}
	}

	return nil
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
	// Login to AWS registry
	if reg.Config.AwsLogin {
		logger.Info("Logging into AWS registry")
		if err := reg.dockerCmd.LoginAws(); err != nil {
			return err
		}
	}

	// Start project
	logger.Info("Starting %s", p.String())
	dc := reg.dockerCmd.CreateAndStartProjectCommand(p, recreate, update)

	return dc.Execute()
}

func (reg *DockerProjectRegistry) StopOtherActiveProjects(p *project.Project) error {
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
		if !ap.IsRunning() {
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
	logger.Info("Stopping %s", p.String())

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
	if err := reg.StopOtherActiveProjects(p); err != nil {
		return err
	}

	logger.Info("Restarting %s", p.String())

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
	logger.Info("Removing %s", p.String())

	exists, err := reg.ProjectExists(p)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warning("%s not found", helpers.ToTitle(p.String()))
		return nil
	}

	dc := reg.dockerCmd.RemoveProjectCommand(p)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) BuildProject(p *project.Project, noCache bool) error {
	logger.Info("Building %s", p.String())
	dc := reg.dockerCmd.BuildProjectCommand(p, noCache)
	return dc.Execute()
}

func (reg *DockerProjectRegistry) Terminal(p *project.Project, cmd []string) error {
	logger.Info("Running terminal for %s", p.String())

	// Set default service
	if !p.IsServiceDefined() {
		p.SetServiceName(reg.Config.TerminalDefaultService)
	}

	// Set default command
	if len(cmd) == 0 {
		cmd = strings.Split(reg.Config.TerminalDefaultCommand, " ")
	}

	dc := reg.dockerCmd.TerminalCommand(p, cmd)

	return dc.Execute()
}

func (reg *DockerProjectRegistry) Code(p *project.Project, dir string) error {
	logger.Info("Opening code editor for %s", p.String())

	// Set default service
	if !p.IsServiceDefined() {
		logger.Debug("Setting default service name")
		p.SetServiceName(reg.Config.VscodeDefaultService)
	}

	// Set default directory
	if dir == "" {
		dir = reg.Config.VscodeDefaultDir
	}

	container, err := reg.ServiceContainer(p)
	if err != nil {
		return err
	}
	if container == nil {
		logger.Warning("%s not found", helpers.ToTitle(p.String()))
		return nil
	}

	variant := reg.Config.VscodeBinary

	return reg.dockerCmd.OpenCode(container, dir, variant)
}

func (reg *DockerProjectRegistry) Logs(p *project.Project, opts docker.LogsOptions) error {
	logger.Info("Showing logs for %s", p.String())
	dc := reg.dockerCmd.LogsCommand(p, opts)
	return dc.Execute()
}
