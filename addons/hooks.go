package addons

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/marcinhlybin/docker-env/logger"
)

type Hook struct {
	Name         string
	Path         string
	Args         []string
	StdoutPrefix string
}

func NewHook(name string, path string, args ...string) *Hook {
	return &Hook{
		Name:         name,
		Path:         path,
		Args:         args,
		StdoutPrefix: "  ",
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

func (h *Hook) Run() error {
	if h.Path == "" {
		return nil
	}

	// Check if hook exists
	if _, err := os.Stat(h.Path); err != nil {
		return fmt.Errorf("cannot open %s hook '%s': %w", h.Name, h.Path, err)
	}

	logger.Debug("Executing %s hook %s", h.Name, h.Path)
	return h.executeCommand()
}

func (h *Hook) executeCommand() error {
	cmd := exec.Command(h.Path, h.Args...)

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

	var wg sync.WaitGroup

	// Print stdout and stderr
	wg.Add(1)
	go h.printStdout(stdout, &wg)
	wg.Add(1)
	go h.printStderr(stderr, &wg)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s hook execution failed: %w", h.Name, err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

func (h *Hook) printStdout(pipe io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		logger.Info("%s%s", h.StdoutPrefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error reading stdout: %v", err)
	}
}

func (h *Hook) printStderr(pipe io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		logger.Error(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error reading stderr: %v", err)
	}
}
