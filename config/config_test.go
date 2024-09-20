package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/marcinhlybin/docker-env/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig_DefaultConfig(t *testing.T) {
	configContent := `
compose_project_name: test_project
vars:
  - var1
  - var2
env_files:
  - %s
compose_file: docker-compose.yml
compose_default_profile: app
aws_login: true
aws_region: us-west-2
aws_repository: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo
`
	// Create temp env file
	envFile := test_helpers.CreateTempFile(t, "var1=value1\nvar2=value")
	defer os.Remove(envFile)

	// Inject the env file path into the config content
	configContent = fmt.Sprintf(configContent, envFile)

	configPath := test_helpers.CreateTempFile(t, configContent)
	defer os.Remove(configPath)

	cfg := config.NewConfig()
	err := cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, "test_project", cfg.ComposeProjectName)
	assert.Equal(t, []string{"var1", "var2"}, cfg.Vars)
	assert.Equal(t, 1, len(cfg.EnvFiles))
	assert.Equal(t, "docker-compose.yml", cfg.ComposeFile)
	assert.Equal(t, "app", cfg.ComposeDefaultProfile)
	assert.True(t, cfg.AwsLogin)
	assert.Equal(t, "us-west-2", cfg.AwsRegion)
	assert.Equal(t, "123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo", cfg.AwsRepository)
}

func TestNewConfig_OverrideConfig(t *testing.T) {
	configContent := `
compose_project_name: test_project
compose_file: docker-compose.yml
compose_default_profile: app
aws_login: true
aws_region: us-west-2
aws_repository: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo
`
	overrideContent := `
compose_project_name: override_project
aws_region: us-east-1
`
	configPath := test_helpers.CreateTempFile(t, configContent)
	defer os.Remove(configPath)

	overridePath := test_helpers.CreateTempFile(t, overrideContent)
	defer os.Remove(overridePath)

	// Temporarily override the default override config path
	config.OverrideConfigPath = overridePath

	cfg := config.NewConfig()
	err := cfg.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "override_project", cfg.ComposeProjectName)
	assert.Equal(t, "us-east-1", cfg.AwsRegion)
}

func TestNewConfig_MissingConfigFile(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
}

func TestNewConfig_InvalidConfigFile(t *testing.T) {
	configContent := `
invalid_yaml_content
`
	configPath := test_helpers.CreateTempFile(t, configContent)
	defer os.Remove(configPath)

	cfg := config.NewConfig()
	err := cfg.LoadConfig(configPath)
	assert.Error(t, err)
}
