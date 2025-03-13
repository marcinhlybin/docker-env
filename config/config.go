package config

import (
	"fmt"
	"os"

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
	ComposeProgress        string   `yaml:"compose_progress"`
	ComposeDefaultProfile  string   `yaml:"compose_default_profile"`
	ComposeSidecarProfile  string   `yaml:"compose_sidecar_profile"`
	GitDefaultBranch       string   `yaml:"git_default_branch"`
	EnvFiles               []string `yaml:"env_files"`
	TerminalDefaultService string   `yaml:"terminal_default_service"`
	TerminalDefaultCommand string   `yaml:"terminal_default_command"`
	VscodeDefaultService   string   `yaml:"vscode_default_service"`
	VscodeDefaultDir       string   `yaml:"vscode_default_dir"`
	VscodeBinary           string   `yaml:"vscode_binary"`
	AwsLogin               bool     `yaml:"aws_login"`
	AwsRegion              string   `yaml:"aws_region"`
	AwsRepository          string   `yaml:"aws_repository"`
	AwsMfa                 bool     `yaml:"aws_mfa"`
	AwsMfaDurationSeconds  int      `yaml:"aws_mfa_duration_seconds"`
	PreStartHooks          []string `yaml:"pre_start_hooks"`
	PostStartHooks         []string `yaml:"post_start_hooks"`
	PostStopHooks          []string `yaml:"post_stop_hooks"`
	RequiredVars           []string `yaml:"required_vars"`
	ShowExecutedCommands   bool     `yaml:"show_executed_commands"`
}

func NewConfig() *Config {
	return &Config{
		ComposeFile:            "docker-compose.yml",
		ComposeFileOverride:    "docker-compose.override.yml",
		ComposeProgress:        "tty",
		ComposeDefaultProfile:  "app",
		ComposeSidecarProfile:  "sidecar",
		GitDefaultBranch:       "master",
		EnvFiles:               []string{},
		TerminalDefaultService: "app",
		TerminalDefaultCommand: "/bin/bash",
		VscodeDefaultService:   "app",
		VscodeDefaultDir:       "/",
		VscodeBinary:           "code",
		PreStartHooks:          []string{},
		PostStartHooks:         []string{},
		PostStopHooks:          []string{},
		RequiredVars:           []string{},
		ShowExecutedCommands:   true,
		AwsLogin:               false,
		AwsMfa:                 false,
		AwsMfaDurationSeconds:  3600,
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
	fmt.Println("Compose file override:", c.ComposeFileOverride)
	fmt.Println("Compose progress:", c.ComposeProgress)
	fmt.Println("Compose default profile:", c.ComposeDefaultProfile)
	fmt.Println("Compose sidecar profile:", c.ComposeSidecarProfile)
	fmt.Println()
	fmt.Println("Git default branch:", c.GitDefaultBranch)
	fmt.Println()
	printList("Env files:", c.EnvFiles)
	fmt.Println()
	printList("Required vars:", c.RequiredVars)
	fmt.Println()
	printList("Pre-start hooks:", c.PreStartHooks)
	fmt.Println()
	printList("Post-start hooks:", c.PostStartHooks)
	fmt.Println()
	printList("Post-stop hooks:", c.PostStopHooks)
	fmt.Println()
	fmt.Println("AWS login:", c.AwsLogin)
	fmt.Println("AWS region:", c.AwsRegion)
	fmt.Println("AWS repository:", c.AwsRepository)
	fmt.Println("AWS MFA:", c.AwsMfa)
	fmt.Println("AWS MFA duration seconds:", c.AwsMfaDurationSeconds)
	fmt.Println()
	fmt.Println("Terminal default service:", c.TerminalDefaultService)
	fmt.Println("Terminal default command:", c.TerminalDefaultCommand)
	fmt.Println("VSCode default service:", c.VscodeDefaultService)
	fmt.Println("VSCode default directory:", c.VscodeDefaultDir)
	fmt.Println("VSCode binary:", c.VscodeBinary)
	fmt.Println()
	fmt.Println("Show executed commands:", c.ShowExecutedCommands)
	fmt.Println()

	return nil
}

func printList(title string, items []string) {
	fmt.Println(title)
	for _, item := range items {
		fmt.Println("  -", item)
	}
}
