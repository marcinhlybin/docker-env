package project

import (
	"fmt"

	"github.com/marcinhlybin/docker-env/git"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
)

type Project struct {
	Name        string
	ServiceName string
}

func NewProject(projectName, serviceName string) (*Project, error) {
	logger.Debug("Creating project project_name='%s' service_name='%s'", projectName, serviceName)

	p := &Project{}

	var err error
	if projectName != "" {
		p.SetProjectName(projectName)
	} else {
		p.SetProjectNameFromGitBranch()
	}
	if err != nil {
		return nil, err
	}

	if serviceName != "" {
		err = p.SetServiceName(serviceName)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Project) SetProjectName(name string) error {
	if err := validateProjectName(name); err != nil {
		return err
	}
	p.Name = name
	return nil
}

func (p *Project) SetProjectNameFromGitBranch() error {
	name, err := git.CurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get project name from git branch: %w", err)
	}
	// Simplify project name
	name = helpers.TrimToLastSlash(name)

	return p.SetProjectName(name)
}

func (p *Project) SetServiceName(name string) error {
	if err := validateServiceName(name); err != nil {
		return err
	}
	p.ServiceName = name
	return nil
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
