package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectMethods(t *testing.T) {
	project := &Project{
		Name:        "test-project",
		ServiceName: "test-service",
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
}
