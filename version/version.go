package version

import (
	"fmt"
)

const Version = "1.0"

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

func PrintBuildDate() {
	fmt.Println(BuildDate)
}

func PrintCommitHash() {
	fmt.Println(CommitHash)
}
