package project

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{"EmptyName", "", true},
		{"ValidName", "Valid_Project-123", false},
		{"InvalidNameSpecialChars", "Invalid@Name!", true},
		{"InvalidNameSpaces", "Invalid Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		wantErr     bool
	}{
		{"EmptyName", "", false},
		{"NonEmptyName", "ServiceName", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServiceName(tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateServiceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
