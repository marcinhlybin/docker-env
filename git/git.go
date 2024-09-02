package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func IsGitRepo() bool {
	_, err := git.PlainOpen(".")
	return err == nil
}

func GetCurrentBranch() (string, error) {
	// Open the current repository in the current directory
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %s", err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("no commits in the repository: %s", err)
	}

	// Extract and print the branch name
	branchName := ref.Name().Short()

	return branchName, nil
}
