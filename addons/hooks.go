package addons

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/marcinhlybin/docker-env/logger"
)

type Hook struct {
	Name string
	Path string
	Args []string
}

func NewHook(name string, path string, args ...string) *Hook {
	return &Hook{
		Name: name,
		Path: path,
		Args: args,
	}
}

func NewPreStartHook(path string, args ...string) *Hook {
	return NewHook("pre-start", path, args...)
}

func NewPostStartHook(path string, args ...string) *Hook {
	return NewHook("post-start", path, args...)
}

func NewPostStopHook(path string, args ...string) *Hook {
	return NewHook("post-stop", path, args...)
}

func (s *Hook) Run() error {
	if s.Path == "" {
		return nil
	}

	// Check if hook exists
	if _, err := os.Stat(s.Path); err != nil {
		return fmt.Errorf("cannot open %s hook '%s': %w", s.Name, s.Path, err)
	}

	logger.Info("Running %s hook %s", s.Name, s.Path)
	return s.executeCommand()
}

func (s *Hook) executeCommand() error {
	cmd := exec.Command(s.Path, s.Args...)

	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)

	// Read and print each line with the prefix
	hookOutputPrefix := "  "
	go func() {
		for scanner.Scan() {
			logger.Info("%s%s", hookOutputPrefix, scanner.Text())
		}
	}()

	go func() {
		for scannerErr.Scan() {
			logger.Error(scannerErr.Text())
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s hook execution failed: %w", s.Name, err)
	}

	return nil
}
