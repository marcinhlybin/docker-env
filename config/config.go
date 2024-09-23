package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Defaults for the config
var (
	ConfigPath         = ".docker-env/config.yml"
	ConfigPathOverride = ".docker-env/config.override.yml"
)

type Config struct {
	Path                   string
	ComposeProjectName     string   `yaml:"compose_project_name"`
	ComposeFile            string   `yaml:"compose_file"`
	ComposeFileOverride    string   `yaml:"compose_file_override"`
	ComposeDefaultProfile  string   `yaml:"compose_default_profile"`
	ComposeSidecarProfile  string   `yaml:"compose_sidecar_profile"`
	EnvFiles               []string `yaml:"env_files"`
	TerminalDefaultService string   `yaml:"terminal_default_service"`
	TerminalDefaultCommand string   `yaml:"terminal_default_command"`
	VscodeDefaultService   string   `yaml:"vscode_default_service"`
	VscodeDefaultDir       string   `yaml:"vscode_default_dir"`
	AwsLogin               bool     `yaml:"aws_login"`
	AwsRegion              string   `yaml:"aws_region"`
	AwsRepository          string   `yaml:"aws_repository"`
	PreStartHook           string   `yaml:"pre_start_hook"`
	PostStartHook          string   `yaml:"post_start_hook"`
	PostStopHook           string   `yaml:"post_stop_hook"`
	RequiredVars           []string `yaml:"required_vars"`
	ShowExecutedCommands   bool     `yaml:"show_executed_commands"`
}

func NewConfig() *Config {
	return &Config{
		ComposeFile:            "docker-compose.yml",
		ComposeFileOverride:    "docker-compose.override.yml",
		ComposeDefaultProfile:  "app",
		ComposeSidecarProfile:  "sidecar",
		EnvFiles:               []string{},
		TerminalDefaultService: "app",
		TerminalDefaultCommand: "/bin/bash",
		VscodeDefaultService:   "app",
		VscodeDefaultDir:       "/",
		PreStartHook:           "",
		PostStartHook:          "",
		PostStopHook:           "",
		RequiredVars:           []string{},
		ShowExecutedCommands:   true,
	}
}

func (cfg *Config) LoadConfig(path string) error {
	if path == "" {
		path = ConfigPath
	}

	// Read main config file
	if err := readConfigFile(path, cfg); err != nil {
		return err
	}

	// Read override config file if it exists
	if _, err := os.Stat(ConfigPathOverride); err == nil {
		if err := readConfigFile(ConfigPathOverride, cfg); err != nil {
			return err
		}
	}

	// Validate config
	if err := cfg.validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func readConfigFile(path string, cfg *Config) error {
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

func (c *Config) ShowConfig() error {
	fmt.Println("Compose project name:", c.ComposeProjectName)
	fmt.Println("Compose file:", c.ComposeFile)
	fmt.Println("Compose default profile:", c.ComposeDefaultProfile)
	fmt.Println("Compose sidecar profile:", c.ComposeSidecarProfile)
	fmt.Println()
	fmt.Println("Env files:", strings.Join(c.EnvFiles, ", "))
	fmt.Println("Required vars:", strings.Join(c.RequiredVars, ", "))
	fmt.Println()
	fmt.Println("Pre-start hook:", c.PreStartHook)
	fmt.Println("Post-start hook:", c.PostStartHook)
	fmt.Println("Post-stop hook:", c.PostStopHook)
	fmt.Println()
	fmt.Println("AWS login:", c.AwsLogin)
	fmt.Println("AWS region:", c.AwsRegion)
	fmt.Println("AWS repository:", c.AwsRepository)
	fmt.Println()
	fmt.Println("Terminal default service:", c.TerminalDefaultService)
	fmt.Println("Terminal default command:", c.TerminalDefaultCommand)
	fmt.Println("VSCode default service:", c.VscodeDefaultService)
	fmt.Println("VSCode default directory:", c.VscodeDefaultDir)
	fmt.Println()
	fmt.Println("Show executed commands:", c.ShowExecutedCommands)
	return nil
}
