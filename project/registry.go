package project

import "github.com/marcinhlybin/docker-env/config"

type ProjectRegistry interface {
	StartProject(p *Project, recreate, update bool) error
	StopProject(p *Project) error
	RestartProject(p *Project) error
	RemoveProject(p *Project) error
	BuildProject(p *Project, noCache bool) error
	ProjectExists(p *Project) (bool, error)
	ListProjects(verbose bool) error
	RunTerminal(p *Project) error
	Cleanup() error
	Config() *config.Config
}

func NewProjectRegistry(createFunc func(*config.Config) ProjectRegistry, cfg *config.Config) ProjectRegistry {
	if createFunc == nil {
		panic("createFunc is nil")
	}
	return createFunc(cfg)
}
