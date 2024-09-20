package addons

import (
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

	logger.Info("Running", s.Name, "scripts")
	return s.execute()
}

func (s *Script) execute() error {
	cmd := exec.Command("/bin/sh", s.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s script failed: %w", s.Name, err)
	}

	return nil
}
