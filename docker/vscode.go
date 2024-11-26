package docker

import (
	"encoding/hex"
	"fmt"
	"os/exec"

	"github.com/marcinhlybin/docker-env/logger"
)

func (dc *DockerCmd) OpenCode(c *Container, dir string, binary string) error {
	// Encode the container name to hexadecimal
	encodedName := hex.EncodeToString([]byte(c.Name))

	// Format the command with the encoded container name and directory
	command := fmt.Sprintf("%s --folder-uri=vscode-remote://attached-container+%s/%s", binary, encodedName, dir)

	// Use shell to exeucute the command
	cmd := exec.Command("/bin/sh", "-c", command)

	logger.Execute(cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %v", err)
	}

	return nil
}
