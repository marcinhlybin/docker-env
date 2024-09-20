package config_test

import (
	"os"
	"testing"

	"github.com/marcinhlybin/docker-env/config"
	"github.com/stretchr/testify/assert"
)

func createTempConfigFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func TestNewConfig_DefaultConfig(t *testing.T) {
	configContent := `
project: test_project
secrets:
  - secret1
  - secret2
env_files:
  - .env
compose_file: docker-compose.yml
compose_default_profile: app
aws_login: true
aws_region: us-west-2
aws_repository: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo
`
	configPath := createTempConfigFile(t, configContent)
	defer os.Remove(configPath)

	cfg, err := config.NewConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "test_project", cfg.Project)
	assert.Equal(t, []string{"secret1", "secret2"}, cfg.Secrets)
	assert.Equal(t, []string{".env"}, cfg.EnvFiles)
	assert.Equal(t, "docker-compose.yml", cfg.ComposeFile)
	assert.Equal(t, "app", cfg.ComposeDefaultProfile)
	assert.True(t, cfg.AwsLogin)
	assert.Equal(t, "us-west-2", cfg.AwsRegion)
	assert.Equal(t, "123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo", cfg.AwsRepository)
}

func TestNewConfig_OverrideConfig(t *testing.T) {
	configContent := `
project: test_project
secrets:
  - secret1
  - secret2
env_files:
  - .env
compose_file: docker-compose.yml
compose_default_profile: app
aws_login: true
aws_region: us-west-2
aws_repository: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo
`
	overrideContent := `
project: override_project
aws_region: us-east-1
`
	configPath := createTempConfigFile(t, configContent)
	defer os.Remove(configPath)

	overridePath := createTempConfigFile(t, overrideContent)
	defer os.Remove(overridePath)

	// Temporarily override the default override config path
	config.DefaultOverrideConfigPath = overridePath

	cfg, err := config.NewConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "override_project", cfg.Project)
	assert.Equal(t, "us-east-1", cfg.AwsRegion)
}

func TestNewConfig_MissingConfigFile(t *testing.T) {
	_, err := config.NewConfig("nonexistent.yaml")
	assert.Error(t, err)
}

func TestNewConfig_InvalidConfigFile(t *testing.T) {
	configContent := `
invalid_yaml_content
`
	configPath := createTempConfigFile(t, configContent)
	defer os.Remove(configPath)

	_, err := config.NewConfig(configPath)
	assert.Error(t, err)
}
