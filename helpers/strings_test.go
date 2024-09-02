package helpers_test

import (
	"testing"

	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/stretchr/testify/assert"
)

func TestToTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "multiple characters",
			input:    "abc def",
			expected: "Abc def",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := helpers.ToTitle(test.input)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestTrimToLastSlash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no slash", "path", "path"},
		{"single slash", "path/", ""},
		{"multiple slashes", "path/to/file", "file"},
		{"trailing slash", "path/to/file/", ""},
		{"empty string", "", ""},
		{"only slash", "/", ""},
		{"double slashes", "//", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := helpers.TrimToLastSlash(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
