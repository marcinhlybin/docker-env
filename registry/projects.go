package registry

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

func (registry *DockerProjectRegistry) ListProjects(verbose bool) error {
	if verbose {
		return registry.ListContainers()
	}

	includeStopped := true
	projects, err := registry.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	for _, p := range projects {
		fmt.Println(p.Name)
	}

	return nil
}

func (registry *DockerProjectRegistry) fetchProjects(includeStopped bool) ([]*project.Project, error) {
	logger.Debug("Fetching project names")
	dc := registry.dockerCmd.FetchProjectsCommand(includeStopped)
	names, err := dc.ExecuteWithOutput()
	if err != nil {
		return nil, err
	}

	logger.Debug("Found %d projects", len(names))

	return registry.createProjectsFromNames(names)
}

func (registry *DockerProjectRegistry) createProjectsFromNames(names []string) ([]*project.Project, error) {
	var projects []*project.Project

	for _, name := range names {
		name = registry.trimComposeProjectNamePrefix(name)

		p, err := project.NewProject(name, "")
		if err != nil {
			logger.Warning("Skipping %s due to error: %v", name, err)
			continue
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (registry *DockerProjectRegistry) trimComposeProjectNamePrefix(name string) string {
	return strings.TrimPrefix(name, registry.config.Project+"-")
}
