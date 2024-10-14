package docker

import (
	"github.com/marcinhlybin/docker-env/project"
)

func (dc *DockerCmd) CreateAndStartProjectCommand(p *project.Project, recreate, update bool) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("up", "--detach", "--wait")

	if update {
		dc.WithArgs("--pull", "always", "--force-recreate")
	} else if recreate {
		dc.WithArgs("--force-recreate")
	}

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) StopProjectCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithSidecarProfile()
	dc.WithProjectName(p)
	dc.WithArgs("stop")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) RestartProjectCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithSidecarProfile()
	dc.WithProjectName(p)
	dc.WithArgs("restart")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) RemoveProjectCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithSidecarProfile()
	dc.WithProjectName(p)
	dc.WithArgs("down", "--volumes")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) FetchProjectsCommand(includeStopped bool) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithArgs("ls", "--format", "json")
	dc.WithProjectFilter()

	if includeStopped {
		dc.WithArgs("-a")
	}

	return dc
}

func (dc *DockerCmd) FetchProjectContainersCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("ps", "-a", "--no-trunc", "--format", "json")
	return dc
}

// FetchAllContainersCommand uses docker ps command (not docker compose) to fetch containers
// Filtering is done by project prefix which must be a part of container name to make it work
// Docker compose configuration must define container name as:
// container_name: $COMPOSE_PROJECT_NAME-service_name_here
//
// It is not possible to list all containers using docker compose command
// because --project-name argument is mandatory
func (dc *DockerCmd) FetchAllContainersCommand(includeStopped bool) *DockerCmd {
	dc.DockerCommand()
	dc.WithArgs("ps", "--no-trunc", "--format", "json")
	dc.WithProjectFilter()

	if includeStopped {
		dc.WithArgs("-a")
	}

	return dc
}

func (dc *DockerCmd) BuildProjectCommand(p *project.Project, noCache bool) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithSidecarProfile()
	dc.WithArgs("build")

	if noCache {
		dc.WithArgs("--no-cache")
	}

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) TerminalCommand(p *project.Project, cmd []string) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)

	dc.WithArgs("exec", "-it", p.ServiceName)
	dc.WithArgs(cmd...)

	return dc
}

func (dc *DockerCmd) FetchImagesCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("images", "--format", "json")
	return dc
}

func (dc *DockerCmd) RemoveImageCommand(id string) *DockerCmd {
	dc.DockerCommand()
	dc.WithArgs("rmi", id)
	return dc
}

type LogsOptions struct {
	FollowOutput   bool
	ShowTimestamps bool
}

func (dc *DockerCmd) LogsCommand(p *project.Project, opts LogsOptions) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("logs")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	if opts.FollowOutput {
		dc.WithArgs("--follow")
	}

	if opts.ShowTimestamps {
		dc.WithArgs("--timestamps")
	}

	return dc
}
