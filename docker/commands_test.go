package docker_test

import (
	"testing"

	"github.com/marcinhlybin/docker-env/project"
	"github.com/marcinhlybin/docker-env/test_helpers"
)

func TestCreateAndStartProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		project      *project.Project
		recreate     bool
		update       bool
		expectedArgs []string
	}{
		{
			name:     "Basic command",
			project:  &project.Project{Name: "test"},
			recreate: false,
			update:   false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"up", "--detach", "--wait",
			},
		},
		{
			name:     "With recreate",
			project:  &project.Project{Name: "test"},
			recreate: true,
			update:   false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"up", "--detach", "--wait",
				"--force-recreate",
			},
		},
		{
			name:     "With update",
			project:  &project.Project{Name: "test"},
			recreate: false,
			update:   true,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"up", "--detach", "--wait",
				"--pull", "always",
				"--force-recreate",
			},
		},
		{
			name:     "With recreate and update",
			project:  &project.Project{Name: "test"},
			recreate: true,
			update:   true,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"up", "--detach", "--wait",
				"--pull", "always",
				"--force-recreate",
			},
		},
		{
			name: "With service",
			project: &project.Project{
				Name:        "test",
				ServiceName: "test-service",
			},
			recreate: false,
			update:   false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"up", "--detach", "--wait",
				"test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.CreateAndStartProjectCommand(tt.project, tt.recreate, tt.update)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestStopProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		project      *project.Project
		expectedArgs []string
	}{
		{
			name:    "Basic command",
			project: &project.Project{Name: "test"},
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"stop",
			},
		},
		{
			name: "With service",
			project: &project.Project{
				Name:        "test",
				ServiceName: "test-service",
			},
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"stop", "test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.StopProjectCommand(tt.project)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestRestartProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		project      *project.Project
		expectedArgs []string
	}{
		{
			name:    "Basic command",
			project: &project.Project{Name: "test"},
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"restart",
			},
		},
		{
			name: "With service",
			project: &project.Project{
				Name:        "test",
				ServiceName: "test-service",
			},
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"restart", "test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.RestartProjectCommand(tt.project)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestRemoveProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		project      *project.Project
		recreate     bool
		update       bool
		expectedArgs []string
	}{
		{
			name:     "Basic command",
			project:  &project.Project{Name: "test"},
			recreate: false,
			update:   false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"down", "--volumes",
			},
		},
		{
			name: "With service",
			project: &project.Project{
				Name:        "test",
				ServiceName: "test-service",
			},
			recreate: false,
			update:   false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"down", "--volumes",
				"test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.RemoveProjectCommand(tt.project)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestFetchProjectsCommand(t *testing.T) {
	tests := []struct {
		name           string
		includeStopped bool
		expectedArgs   []string
	}{
		{
			name:           "Basic command",
			includeStopped: false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"ls", "-q",
				"--filter", "name=project-",
			},
		},
		{
			name:           "Include stopped",
			includeStopped: true,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"ls", "-q",
				"--filter", "name=project-",
				"-a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.FetchProjectsCommand(tt.includeStopped)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestFetchContainersCommand(t *testing.T) {
	tests := []struct {
		name         string
		stopped      bool
		expectedArgs []string
	}{
		{
			name: "Basic command",
			expectedArgs: []string{
				"ps", "-a", "--format", "json",
				"--filter", "name=project-",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.FetchContainersCommand()
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}

func TestBuildCommand(t *testing.T) {
	tests := []struct {
		name         string
		project      *project.Project
		service      string
		noCache      bool
		expectedArgs []string
	}{
		{
			name:    "Basic command",
			project: &project.Project{Name: "test"},
			noCache: false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"build",
			},
		},
		{
			name:    "With no cache",
			project: &project.Project{Name: "test"},
			noCache: true,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"build", "--no-cache",
			},
		},
		{
			name: "With service",
			project: &project.Project{
				Name:        "test",
				ServiceName: "test-service",
			},
			noCache: false,
			expectedArgs: []string{
				"compose",
				"--file", "docker-compose.yml",
				"--profile", "default",
				"--progress", "tty",
				"--env-file", "env1",
				"--env-file", "env2",
				"--project-name", "project-test",
				"build",
				"test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := test_helpers.SetupDockerCompose(t)
			dc.BuildProjectCommand(tt.project, tt.noCache)
			test_helpers.CheckCommandArgs(t, dc.Args, tt.expectedArgs)
		})
	}
}
