package test_helpers

import (
	"testing"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/docker"
)

func SetupDockerCompose(t *testing.T) *docker.DockerCmd {
	t.Helper() // mark as helper function

	cfg := &config.Config{
		ComposeFile:    "docker-compose.yml",
		ComposeProfile: "default",
		EnvFiles:       []string{"env1", "env2"},
		Project:        "project",
	}

	return docker.NewDockerCmd(cfg)
}

func CheckCommand(t *testing.T, cmd, expectedCmd string) {
	t.Helper() // mark as helper function
	if cmd != expectedCmd {
		t.Errorf("Expected command mismatch. Expected %s, got %s", expectedCmd, cmd)
	}
}

func CheckCommandArgs(t *testing.T, args, expectedArgs []string) {
	t.Helper() // mark as helper function

	if !EqualSlices(args, expectedArgs) {
		t.Errorf("Expected arguments mismatch. Expected %v, got %v", len(expectedArgs), len(args))
	}
	// Check if all expected arguments are present
	for _, expected := range expectedArgs {
		if !Contains(args, expected) {
			t.Errorf("Expected argument '%s' not found in args: %s", expected, args)
		}
	}
}
