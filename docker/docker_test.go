package docker_test

import (
	"testing"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/project"
	"github.com/stretchr/testify/assert"
)

func TestNewDockerCmd(t *testing.T) {
	cfg := &config.Config{
		ComposeFile:           "docker-compose.yml",
		ComposeDefaultProfile: "dev",
		EnvFiles:              []string{"env1", "env2"},
		ComposeProjectName:    "myproject",
	}
	dc := docker.NewDockerCmd(cfg)

	assert.NotNil(t, dc)
	assert.Equal(t, "docker", dc.Cmd)
	assert.Empty(t, dc.Args)
	assert.Equal(t, cfg, dc.Config)
}

func TestDockerCommand(t *testing.T) {
	cfg := &config.Config{}
	dc := docker.NewDockerCmd(cfg).DockerCommand()

	assert.Empty(t, dc.Args)
}

func TestDockerComposeCommand(t *testing.T) {
	cfg := &config.Config{
		ComposeFile:           "docker-compose.yml",
		ComposeDefaultProfile: "dev",
		EnvFiles:              []string{"env1", "env2"},
	}
	dc := docker.NewDockerCmd(cfg).DockerComposeCommand()

	expectedArgs := []string{
		"compose",
		"--file", "docker-compose.yml",
		"--profile", "dev",
		"--progress", "tty",
		"--env-file", "env1",
		"--env-file", "env2",
	}

	assert.Equal(t, expectedArgs, dc.Args)
}

func TestString(t *testing.T) {
	cfg := &config.Config{}
	dc := docker.NewDockerCmd(cfg).DockerCommand().WithArgs("ps")

	expectedString := "docker ps"
	assert.Equal(t, expectedString, dc.String())
}

func TestSlice(t *testing.T) {
	cfg := &config.Config{}
	dc := docker.NewDockerCmd(cfg).DockerCommand().WithArgs("ps")

	expectedSlice := []string{"docker", "ps"}
	assert.Equal(t, expectedSlice, dc.Slice())
}

func TestWithArgs(t *testing.T) {
	cfg := &config.Config{}
	dc := docker.NewDockerCmd(cfg).DockerCommand().WithArgs("ps", "-a")

	expectedArgs := []string{"ps", "-a"}
	assert.Equal(t, expectedArgs, dc.Args)
}

func TestWithProjectName(t *testing.T) {
	cfg := &config.Config{
		ComposeProjectName: "myproject",
	}
	p := &project.Project{
		Name: "service",
	}
	dc := docker.NewDockerCmd(cfg).DockerCommand().WithProjectName(p)

	expectedArgs := []string{"--project-name", "myproject-service"}
	assert.Equal(t, expectedArgs, dc.Args)
}

func TestWithProjectFilter(t *testing.T) {
	cfg := &config.Config{
		ComposeProjectName: "myproject",
	}
	dc := docker.NewDockerCmd(cfg).DockerCommand().WithProjectFilter()

	expectedArgs := []string{"--filter", "name=myproject-"}
	assert.Equal(t, expectedArgs, dc.Args)
}
