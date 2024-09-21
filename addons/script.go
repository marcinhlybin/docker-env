package addons

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/marcinhlybin/docker-env/logger"
)

type Script struct {
	Name string
	Path string
}

func NewScript(name, path string) *Script {
	return &Script{
		Name: name,
		Path: path,
	}
}

// Helper function to run a script
func RunScript(name, path string) error {
	s := NewScript(name, path)
	return s.Run()
}

func (s *Script) Run() error {
	if s.Path == "" {
		return nil
	}

	// Check if script exists
	if _, err := os.Stat(s.Path); err != nil {
		return fmt.Errorf("cannot open %s script '%s': %w", s.Name, s.Path, err)
	}

	logger.Info("Running %s scripts", s.Name)
	return s.execute()
}

func (s *Script) execute() error {
	cmd := exec.Command("/bin/sh", s.Path)

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
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)

	// Read and print each line with the prefix
	go func() {
		for scanner.Scan() {
			logger.Info(scanner.Text())
		}
	}()

	go func() {
		for scannerErr.Scan() {
			logger.Error(scannerErr.Text())
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	return nil
}
