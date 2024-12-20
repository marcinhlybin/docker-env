package registry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type DockerComposeProject struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	ConfigFiles string `json:"configFiles"`
}

func (reg *DockerProjectRegistry) ListProjects(includeStopped bool) error {
	projects, err := reg.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	for _, p := range projects {
		if includeStopped && p.IsRunning() {
			fmt.Println(p.Name, "*")
		} else {
			fmt.Println(p.Name)
		}
	}

	return nil
}

func (reg *DockerProjectRegistry) fetchProjects(includeStopped bool) ([]*project.Project, error) {
	logger.Debug("Fetching project names")
	dc := reg.dockerCmd.FetchProjectsCommand(includeStopped)
	jsonOutput, err := dc.ExecuteWithOutput()
	if err != nil {
		return nil, err
	}
	jsonString := strings.Join(jsonOutput, "")

	return reg.createProjectsFromJson(jsonString)
}

func (reg *DockerProjectRegistry) createProjectsFromJson(jsonString string) ([]*project.Project, error) {
	var projects []*project.Project
	var dockerComposeProjects []*DockerComposeProject

	// Unmarshal json
	err := json.Unmarshal([]byte(jsonString), &dockerComposeProjects)
	if err != nil {
		return nil, err
	}

	for _, dcProject := range dockerComposeProjects {
		name := reg.trimComposeProjectNamePrefix(dcProject.Name)
		p, err := project.NewProject(name, "")
		if err != nil {
			logger.Warning("Skipping %s due to error: %v", name, err)
			continue
		}
		p.SetStatus(dcProject.Status)
		projects = append(projects, p)
	}

	return projects, nil
}

func (reg *DockerProjectRegistry) trimComposeProjectNamePrefix(name string) string {
	return strings.TrimPrefix(name, reg.Config.ComposeProjectName+"-")
}

func (reg *DockerProjectRegistry) removeProjects(projects []*project.Project) error {
	isErr := false
	for _, p := range projects {
		dc := reg.dockerCmd.RemoveProjectCommand(p)
		err := dc.Execute()
		if err != nil {
			isErr = true
			logger.Warning("Could not remove %s", p.String())
			continue
		}
	}

	if isErr {
		return fmt.Errorf("one or more projects could not be removed")
	}

	return nil
}
