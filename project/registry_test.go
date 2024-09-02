package project

import (
	"testing"

	"github.com/marcinhlybin/docker-env/config"
)

type mockProjectRegistry struct{}

func (m *mockProjectRegistry) StartProject(p *Project, recreate, update bool) error { return nil }
func (m *mockProjectRegistry) StopProject(p *Project) error                         { return nil }
func (m *mockProjectRegistry) RestartProject(p *Project) error                      { return nil }
func (m *mockProjectRegistry) RemoveProject(p *Project) error                       { return nil }
func (m *mockProjectRegistry) BuildProject(p *Project, noCache bool) error          { return nil }
func (m *mockProjectRegistry) ProjectExists(p *Project) (bool, error)               { return false, nil }
func (m *mockProjectRegistry) ListProjects(verbose bool) error                      { return nil }
func (m *mockProjectRegistry) RunTerminal(p *Project) error                         { return nil }
func (m *mockProjectRegistry) Cleanup() error                                       { return nil }
func (m *mockProjectRegistry) Config() *config.Config                               { return nil }

func TestNewProjectRegistry(t *testing.T) {
	cfg := &config.Config{}

	t.Run("ValidCreateFunc", func(t *testing.T) {
		createFunc := func(cfg *config.Config) ProjectRegistry {
			return &mockProjectRegistry{}
		}
		registry := NewProjectRegistry(createFunc, cfg)
		if registry == nil {
			t.Errorf("Expected non-nil ProjectRegistry")
		}
	})

	t.Run("NilCreateFunc", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when createFunc is nil")
			}
		}()
		NewProjectRegistry(nil, cfg)
	})
}
