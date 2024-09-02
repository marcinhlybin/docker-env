package project

import (
	"testing"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegistry is a mock implementation of the ProjectRegistry interface
type MockRegistry struct {
	mock.Mock
}

func (m *MockRegistry) StartProject(p *Project, recreate, update bool) error {
	args := m.Called(p, recreate, update)
	return args.Error(0)
}

func (m *MockRegistry) StopProject(p *Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRegistry) RestartProject(p *Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRegistry) RemoveProject(p *Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRegistry) BuildProject(p *Project, noCache bool) error {
	args := m.Called(p, noCache)
	return args.Error(0)
}

func (m *MockRegistry) ProjectExists(p *Project) (bool, error) {
	args := m.Called(p)
	return args.Bool(0), args.Error(1)
}

func (m *MockRegistry) ListProjects(verbose bool) error {
	args := m.Called(verbose)
	return args.Error(0)
}

func (m *MockRegistry) Cleanup() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistry) Config() *config.Config {
	args := m.Called()
	return args.Get(0).(*config.Config)
}

func (m *MockRegistry) RunTerminal(p *Project) error {
	args := m.Called(p)
	return args.Error(0)
}

func TestProjectMethods(t *testing.T) {
	mockRegistry := new(MockRegistry)
	project := &Project{
		Name:        "test-project",
		ServiceName: "test-service",
		registry:    mockRegistry,
	}

	t.Run("IsServiceDefined", func(t *testing.T) {
		assert.True(t, project.IsServiceDefined())
		project.ServiceName = ""
		assert.False(t, project.IsServiceDefined())
	})

	t.Run("String", func(t *testing.T) {
		project.ServiceName = "test-service"
		assert.Equal(t, "service test-project/test-service", project.String())
		project.ServiceName = ""
		assert.Equal(t, "project test-project", project.String())
	})

	t.Run("Start", func(t *testing.T) {
		mockRegistry.On("StartProject", project, true, false).Return(nil)
		err := project.Start(true, false)
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "StartProject", project, true, false)
	})

	t.Run("Stop", func(t *testing.T) {
		mockRegistry.On("StopProject", project).Return(nil)
		err := project.Stop()
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "StopProject", project)
	})

	t.Run("Restart", func(t *testing.T) {
		mockRegistry.On("RestartProject", project).Return(nil)
		err := project.Restart()
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "RestartProject", project)
	})

	t.Run("Remove", func(t *testing.T) {
		mockRegistry.On("RemoveProject", project).Return(nil)
		err := project.Remove()
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "RemoveProject", project)
	})

	t.Run("Build", func(t *testing.T) {
		mockRegistry.On("BuildProject", project, true).Return(nil)
		err := project.Build(true)
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "BuildProject", project, true)
	})

	t.Run("Exists", func(t *testing.T) {
		mockRegistry.On("ProjectExists", project).Return(true, nil)
		exists, err := project.Exists()
		assert.Nil(t, err)
		assert.True(t, exists)
		mockRegistry.AssertCalled(t, "ProjectExists", project)
	})

	t.Run("ListProjects", func(t *testing.T) {
		mockRegistry.On("ListProjects", true).Return(nil)
		err := mockRegistry.ListProjects(true)
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "ListProjects", true)
	})

	t.Run("Cleanup", func(t *testing.T) {
		mockRegistry.On("Cleanup").Return(nil)
		err := mockRegistry.Cleanup()
		assert.Nil(t, err)
		mockRegistry.AssertCalled(t, "Cleanup")
	})

	t.Run("Config", func(t *testing.T) {
		mockConfig := &config.Config{}
		mockRegistry.On("Config").Return(mockConfig)
		cfg := mockRegistry.Config()
		assert.Equal(t, mockConfig, cfg)
		mockRegistry.AssertCalled(t, "Config")
	})
}
