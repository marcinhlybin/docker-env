package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

var alphaNumericAndUnderscoreRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (cfg *Config) validateConfig() error {
	if err := checkComposeProjectName(cfg.ComposeProjectName); err != nil {
		return err
	}

	if err := validateEnvFiles(cfg.EnvFiles); err != nil {
		return err
	}

	if len(cfg.RequiredVars) > 0 {
		if err := validateRequiredVars(cfg.RequiredVars, cfg.EnvFiles); err != nil {
			return err
		}
	}

	return nil
}

func checkComposeProjectName(v string) error {
	if v == "" {
		return fmt.Errorf("'project' option is required")
	}

	if !isAlphaNumericAndUnderscore(v) {
		return fmt.Errorf("'project' value can only contain letters and numbers and _")
	}

	return nil
}

func isAlphaNumericAndUnderscore(s string) bool {
	return alphaNumericAndUnderscoreRegex.MatchString(s)
}

func validateEnvFiles(envFiles []string) error {
	for _, f := range envFiles {
		path, err := filepath.Abs(f)
		if err != nil {
			return fmt.Errorf("error getting absolute path for env file %s: %s", f, err)
		}
		if _, err := godotenv.Read(path); err != nil {
			return fmt.Errorf("error reading env file %s: %s", path, err)
		}
	}
	return nil
}

func validateRequiredVars(vars []string, envFiles []string) error {
	var missingVars []string
	env, err := godotenv.Read(envFiles...)
	if err != nil {
		return err
	}

	for _, v := range vars {
		if _, ok := env[v]; !ok {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing vars in env files: %s", strings.Join(missingVars, ", "))
	}

	return nil
}
