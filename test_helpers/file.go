package test_helpers

import (
	"os"
	"testing"
)

func CreateTempFile(t *testing.T, content string) string {
	t.Helper() // mark as helper function

	tmpfile, err := os.CreateTemp("", "temp-*.txt")
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
