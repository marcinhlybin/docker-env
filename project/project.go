package project

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/git"
	"github.com/marcinhlybin/docker-env/logger"
)

type Project struct {
	Name        string
	ServiceName string
	registry    ProjectRegistry
}

func NewProject(projectName, serviceName string, registry ProjectRegistry) (*Project, error) {
	logger.Debug("Creating project project_name='%s' service_name='%s'", projectName, serviceName)

	projectName, err := resolveProjectName(projectName)
	if err != nil {
		return nil, err
	}

	if err = validateProjectName(projectName); err != nil {
		return nil, err

	}

	if err = validateServiceName(serviceName); err != nil {
		return nil, err
	}

	return &Project{
		Name:        projectName,
		ServiceName: serviceName,
		registry:    registry,
	}, nil
}

func resolveProjectName(projectName string) (string, error) {
	if projectName == "" {
		return getProjectNameFromGitBranch()
	}
	return projectName, nil
}

func getProjectNameFromGitBranch() (string, error) {
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get project name from git branch: %w", err)
	}
	return trimAfterLastSlash(branch), nil
}

func trimAfterLastSlash(s string) string {
	lastSlashIndex := strings.LastIndex(s, "/")
	if lastSlashIndex == -1 {
		return s // No slash found, return the original string
	}
	return s[lastSlashIndex+1:]
}

func (p *Project) IsServiceDefined() bool {
	return p.ServiceName != ""
}

func (p *Project) String() string {
	if p.IsServiceDefined() {
		return fmt.Sprintf("service %s/%s", p.Name, p.ServiceName)
	}
	return fmt.Sprintf("project %s", p.Name)
}

func (p *Project) Start(recreate, update bool) error {
	return p.registry.StartProject(p, recreate, update)
}

func (p *Project) Stop() error {
	return p.registry.StopProject(p)
}

func (p *Project) Restart() error {
	return p.registry.RestartProject(p)
}

func (p *Project) Remove() error {
	return p.registry.RemoveProject(p)
}

func (p *Project) Build(noCache bool) error {
	return p.registry.BuildProject(p, noCache)
}

func (p *Project) Exists() (bool, error) {
	return p.registry.ProjectExists(p)
}

func (p *Project) Terminal() error {
	return p.registry.RunTerminal(p)
}
