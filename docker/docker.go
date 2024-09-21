package docker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type DockerCmd struct {
	Config *config.Config
	Cmd    string
	Args   []string
}

func NewDockerCmd(cfg *config.Config) *DockerCmd {
	logger.ShowCommands(cfg.ShowCommands)

	return &DockerCmd{
		Config: cfg,
		Cmd:    "docker",
		Args:   []string{},
	}
}

// Start building a docker command
func (dc *DockerCmd) DockerCommand() *DockerCmd {
	dc.Args = []string{}
	return dc
}

// Start building a docker compose command
func (dc *DockerCmd) DockerComposeCommand() *DockerCmd {
	dc.Args = []string{
		"compose",
		"--file", dc.Config.ComposeFile,
		"--profile", dc.Config.ComposeDefaultProfile,
		"--progress", "tty"}

	// Add the override file if it exists
	if _, err := os.Stat(dc.Config.ComposeFileOverride); err == nil {
		dc.Args = append(dc.Args, "--file", dc.Config.ComposeFileOverride)
	}

	// Add the env files
	for _, envFile := range dc.Config.EnvFiles {
		dc.Args = append(dc.Args, "--env-file", envFile)
	}

	return dc
}

func (dc *DockerCmd) String() string {
	return fmt.Sprintf("%s %s", dc.Cmd, strings.Join(dc.Args, " "))
}

func (dc *DockerCmd) Slice() []string {
	return append([]string{dc.Cmd}, dc.Args...)
}

func (dc *DockerCmd) WithArgs(args ...string) *DockerCmd {
	dc.Args = append(dc.Args, args...)
	return dc
}

func (dc *DockerCmd) WithProjectName(p *project.Project) *DockerCmd {
	projectName := dc.Config.ComposeProjectName + "-" + p.Name
	dc.Args = append(dc.Args, "--project-name", projectName)
	return dc
}

func (dc *DockerCmd) WithSidecarProfile() *DockerCmd {
	dc.Args = append(dc.Args, "--profile", dc.Config.ComposeSidecarProfile)
	return dc
}

func (dc *DockerCmd) WithProjectFilter() *DockerCmd {
	prefix := "name=" + dc.Config.ComposeProjectName + "-"
	dc.Args = append(dc.Args, "--filter", prefix)
	return dc
}

func (dc *DockerCmd) Execute() error {
	cmd := exec.Command(dc.Cmd, dc.Args...)
	logger.Execute(cmd.String())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// Capture stderr output
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Run(); err != nil {
		// Log the captured stderr output
		logger.Error(stderrBuf.String())
		return fmt.Errorf("error running docker command: %w", err)
	}

	return nil
}

func (dc *DockerCmd) ExecuteWithOutput() ([]string, error) {
	cmd := exec.Command(dc.Cmd, dc.Args...)
	logger.Execute(cmd.String())

	output, err := cmd.CombinedOutput()

	// Convert to string and remove trailing newline
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		return nil, fmt.Errorf("error running docker command: %v", outputStr)
	}

	// Handle empty output
	if outputStr == "" {
		return []string{}, nil
	}

	// Split the string into a slice of strings
	return strings.Split(outputStr, "\n"), nil
}
