package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Defaults for the config
var (
	DefaultConfigPath             = ".docker-env/config.yaml"
	DefaultOverrideConfigPath     = ".docker-env/config.override.yaml"
	DefaultComposeFile            = "docker-compose.yml"
	DefaultComposeProfile         = "app"
	DefaultTerminalDefaultService = "app"
	DefaultTerminalDefaultCommand = "/bin/bash"
	DefaultVscodeDefaultService   = "app"
	DefaultVscodeDefaultDir       = "/"
)

type Config struct {
	Project                string   `yaml:"project"`
	Secrets                []string `yaml:"secrets"`
	EnvFiles               []string `yaml:"env_files"`
	ComposeFile            string   `yaml:"compose_file"`
	ComposeProfile         string   `yaml:"compose_profile"`
	TerminalDefaultService string   `yaml:"terminal_default_service"`
	TerminalDefaultCommand string   `yaml:"terminal_default_command"`
	VscodeDefaultService   string   `yaml:"vscode_default_service"`
	VscodeDefaultDir       string   `yaml:"vscode_default_dir"`
	AwsLogin               bool     `yaml:"aws_login"`
	AwsRegion              string   `yaml:"aws_region"`
	AwsRepository          string   `yaml:"aws_repository"`
}

// Read and parse the config file with fields validation
func NewConfig(path string) (*Config, error) {
	var cfg Config

	if path == "" {
		path = DefaultConfigPath
	}

	// Read main config file
	if err := loadConfigFile(path, &cfg); err != nil {
		return nil, err
	}

	// Read override config file if it exists
	overridePath := DefaultOverrideConfigPath
	if _, err := os.Stat(overridePath); err == nil {
		if err := loadConfigFile(overridePath, &cfg); err != nil {
			return nil, err
		}
	}

	setDefaults(&cfg)

	// Validate config
	if err := ValidateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func loadConfigFile(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	return nil
}

func setDefaults(cfg *Config) {
	if cfg.ComposeFile == "" {
		cfg.ComposeFile = DefaultComposeFile
	}
	if cfg.ComposeProfile == "" {
		cfg.ComposeProfile = DefaultComposeProfile
	}
	if cfg.TerminalDefaultService == "" {
		cfg.TerminalDefaultService = DefaultTerminalDefaultService
	}
	if cfg.TerminalDefaultCommand == "" {
		cfg.TerminalDefaultCommand = DefaultTerminalDefaultCommand
	}
	if cfg.VscodeDefaultService == "" {
		cfg.VscodeDefaultService = DefaultVscodeDefaultService
	}
	if cfg.VscodeDefaultDir == "" {
		cfg.VscodeDefaultDir = DefaultVscodeDefaultDir
	}
}

func (c *Config) ShowConfig() error {
	fmt.Println("Project name:", c.Project)
	fmt.Println("Mandatory secrets:", strings.Join(c.Secrets, ", "))
	fmt.Println("Env files:", strings.Join(c.EnvFiles, ", "))
	fmt.Println("Compose file:", c.ComposeFile)
	fmt.Println("Compose profile:", c.ComposeProfile)
	fmt.Println("Terminal default service:", c.TerminalDefaultService)
	fmt.Println("Terminal default command:", c.TerminalDefaultCommand)
	fmt.Println("VSCode default service:", c.VscodeDefaultService)
	fmt.Println("VSCode default directory:", c.VscodeDefaultDir)
	fmt.Println("AWS login:", c.AwsLogin)
	fmt.Println("AWS region:", c.AwsRegion)
	fmt.Println("AWS repository:", c.AwsRepository)
	return nil
}
