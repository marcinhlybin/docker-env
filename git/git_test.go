package git

import (
	"os"
	"os/exec"
	"testing"
)

func setupTestRepo(t *testing.T) (string, func()) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("/tmp", "docker-env-testrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize a new Git repository
	runGitCommand(t, tempDir, "init")

	// Create a new file in the repository
	filePath := tempDir + "/README.md"
	content := []byte("# Initial Commit\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Add the file to the staging area
	runGitCommand(t, tempDir, "add", "README.md")

	// Commit the file
	runGitCommand(t, tempDir, "commit", "-m", "Initial commit")

	// Create a new branch named "test-branch"
	runGitCommand(t, tempDir, "checkout", "-b", "test-branch")

	// Return the temporary directory and a cleanup function
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run git command %v: %v", args, err)
	}
}

func TestIsGitRepo(t *testing.T) {
	tempDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if !IsGitRepo() {
		t.Errorf("Expected IsGitRepo to return true")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tempDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	branchName, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	expectedBranchName := "test-branch"
	if branchName != expectedBranchName {
		t.Errorf("Expected branch name %s, got %s", expectedBranchName, branchName)
	}
}
