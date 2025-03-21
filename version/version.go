package version

import (
	"fmt"
)

const Version = "2.1.2"

// Set at build time
var (
	BuildDate  = ""
	CommitHash = ""
)

func PrintFullVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Commit hash: %s\n", CommitHash)
}

func PrintVersionString() {
	fmt.Println(Version)
}

func PrintBuildDateString() {
	fmt.Println(BuildDate)
}

func PrintCommitHashString() {
	fmt.Println(CommitHash)
}
