package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func (dc *DockerCmd) LoginAws() error {
	// Check if the AWS CLI binary is available
	if _, err := exec.LookPath("aws"); err != nil {
		return fmt.Errorf("AWS CLI binary not found. Install with 'brew install awscli' or 'sudo apt-get install awscli'")
	}

	// Get the password from AWS command
	awsCmd := exec.Command("aws", "ecr", "get-login-password", "--region", dc.Config.AwsRegion)
	awsOutput, err := awsCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error getting AWS registry password: %v", strings.TrimSpace(string(awsOutput)))
	}

	// Docker login
	loginCmd := exec.Command("docker", "login", "--username", "AWS", "--password-stdin", dc.Config.AwsRepository)
	loginCmd.Stdin = bytes.NewReader(awsOutput)

	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("error running docker login: %v", err)
	}

	return nil
}
