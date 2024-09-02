package config

import (
	"fmt"
	"regexp"
)

var alphaNumericAndUnderscoreRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateConfig(cfg *Config) error {
	if err := checkProjectValue(cfg.Project); err != nil {
		return err
	}
	return nil
}

func checkProjectValue(v string) error {
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
