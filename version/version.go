package version

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const Version = "0.1.0"

// Set at build time
var (
	BuildDate  = ""
	CommitHash = ""
)

func VersionPrinter(c *cli.Context) {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Commit hash: %s\n", CommitHash)
}
