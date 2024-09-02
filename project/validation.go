package project

import (
	"fmt"
	"regexp"
)

var (
	projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func validateProjectName(projectName string) error {
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if !projectNameRegex.MatchString(projectName) {
		return fmt.Errorf("project name can only contain letters and numbers and _")
	}

	return nil
}

func validateServiceName(serviceName string) error {
	if serviceName == "" {
		return nil
	}
	return nil
}
