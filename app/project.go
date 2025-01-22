package app

import (
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

func (ctx *AppContext) CreateProject() (*project.Project, error) {
	p, err := project.NewProject(ctx.ProjectName, ctx.ServiceName)
	if err != nil {
		return nil, err
	}

	// Set project name from git branch
	if ctx.IsBranch {
		if err := p.SetProjectNameFromGitBranch(); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (ctx *AppContext) ActiveProject() (*project.Project, error) {
	if ctx.IsProjectNameSet {
		return ctx.CreateProject()
	}

	// Fetch active project
	reg := ctx.Registry
	p, err := reg.ActiveProject()
	if err != nil {
		return nil, err
	}

	if p == nil {
		logger.Warning("No active project found")
		return nil, nil
	}

	p.SetServiceName(ctx.ServiceName)
	return p, nil
}
