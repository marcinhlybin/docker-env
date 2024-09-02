package helpers

import "testing"

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
			actual := ToTitle(test.input)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
