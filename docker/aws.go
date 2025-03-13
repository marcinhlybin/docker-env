package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

type STSCredentials struct {
	Credentials struct {
		AccessKeyId     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		SessionToken    string `json:"SessionToken"`
		Expiration      string `json:"Expiration"`
	} `json:"Credentials"`
}

func (dc *DockerCmd) LoginAws() error {
	if err := checkAwsCliInstalled(); err != nil {
		return err
	}

	// Prepare AWS environment (optional MFA)
	awsEnv, err := dc.awsEnv()
	if err != nil {
		return err
	}

	// Retrieve the ECR login password using the environment
	password, err := getEcrLoginPassword(awsEnv, dc.Config.AwsRegion)
	if err != nil {
		return err
	}

	// Use the password to Docker login
	if err := dockerLogin(password, dc.Config.AwsRepository); err != nil {
		return err
	}

	return nil
}

func checkAwsCliInstalled() error {
	if _, err := exec.LookPath("aws"); err != nil {
		return fmt.Errorf("AWS CLI binary not found. Install it, e.g. 'brew install awscli' or 'sudo apt-get install awscli'")
	}
	return nil
}

func (dc *DockerCmd) awsEnv() ([]string, error) {
	// Attempt to load .env file if it exists
	_ = godotenv.Load(".env")

	env := os.Environ()

	// If no MFA serial is set, just return the current environment unmodified
	if !dc.Config.AwsMfa {
		return env, nil
	}

	// Get MFA serial and token
	mfaSerial, err := getMfaSerial()
	if err != nil {
		return nil, err
	}

	mfaToken, err := promptMfaToken()
	if err != nil {
		return nil, err
	}

	// Fetch STS credentials
	creds, err := stsCredentials(mfaSerial, mfaToken, dc.Config.AwsMfaDurationSeconds)
	if err != nil {
		return nil, err
	}

	// Append the short-lived STS credentials to our environment
	env = append(env,
		"AWS_ACCESS_KEY_ID="+creds.Credentials.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY="+creds.Credentials.SecretAccessKey,
		"AWS_SESSION_TOKEN="+creds.Credentials.SessionToken,
	)

	return env, nil
}

func getMfaSerial() (string, error) {
	if val, ok := os.LookupEnv("AWS_MFA_SERIAL"); ok {
		return val, nil
	}
	return "", fmt.Errorf("multi-factor authentication (MFA) enabled in config file but environmental variable AWS_MFA_SERIAL is not set")
}

func promptMfaToken() (string, error) {
	fmt.Print("Enter your MFA token code: ")
	var mfaCode string
	_, err := fmt.Scanln(&mfaCode)
	if err != nil {
		return "", fmt.Errorf("error reading MFA code: %v", err)
	}
	return mfaCode, nil
}

func stsCredentials(mfaSerial, mfaToken string, durationSeconds int) (STSCredentials, error) {
	stsCmd := exec.Command("aws", "sts", "get-session-token",
		"--serial-number", mfaSerial,
		"--token-code", mfaToken,
		"--duration-seconds", fmt.Sprintf("%d", durationSeconds))
	stsOutput, err := stsCmd.CombinedOutput()
	if err != nil {
		return STSCredentials{}, fmt.Errorf("error getting AWS STS session token: %v", strings.TrimSpace(string(stsOutput)))
	}

	var creds STSCredentials
	if err := json.Unmarshal(stsOutput, &creds); err != nil {
		return STSCredentials{}, fmt.Errorf("unable to parse STS credentials: %v", err)
	}

	return creds, nil
}

func getEcrLoginPassword(env []string, region string) (string, error) {
	cmd := exec.Command("aws", "ecr", "get-login-password", "--region", region)
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting AWS registry password: %v\nOutput: %s",
			err, strings.TrimSpace(string(output)))
	}

	return string(output), nil
}

func dockerLogin(password, repository string) error {
	loginCmd := exec.Command("docker", "login",
		"--username", "AWS",
		"--password-stdin", repository,
	)
	loginCmd.Stdin = bytes.NewReader([]byte(password))

	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("error running docker login: %v", err)
	}
	return nil
}
