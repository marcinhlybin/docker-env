package docker

import "github.com/marcinhlybin/docker-env/project"

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
	dc.WithProjectName(p)
	dc.WithArgs("stop")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) RestartProjectCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("restart")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) RemoveProjectCommand(p *project.Project) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("down", "--volumes")

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) FetchProjectsCommand(includeStopped bool) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithArgs("ls", "-q")
	dc.WithProjectFilter()

	if includeStopped {
		dc.WithArgs("-a")
	}

	return dc
}

func (dc *DockerCmd) FetchContainersCommand() *DockerCmd {
	dc.DockerCommand()
	dc.WithArgs("ps", "-a", "--format", "json")
	dc.WithProjectFilter()

	return dc
}

func (dc *DockerCmd) BuildProjectCommand(p *project.Project, noCache bool) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)
	dc.WithArgs("build")

	if noCache {
		dc.WithArgs("--no-cache")
	}

	if p.IsServiceDefined() {
		dc.WithArgs(p.ServiceName)
	}

	return dc
}

func (dc *DockerCmd) RunTerminalCommand(p *project.Project, cmd string) *DockerCmd {
	dc.DockerComposeCommand()
	dc.WithProjectName(p)

	if p.IsServiceDefined() {
		dc.WithArgs("exec", "-it", p.ServiceName)
	} else {
		dc.WithArgs("run", "-it", "--rm")
	}

	dc.WithArgs(cmd)

	return dc
}
