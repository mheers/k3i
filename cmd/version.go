package cmd

import (
	"fmt"

	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

// build flags
var (
	VERSION    string
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
	GitBranch  string
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints the version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			helpers.SetLogLevel(LogLevelFlag)
			helpers.PrintInfo()
			printVersion()
		},
	}
)

func printVersion() {
	fmt.Printf("Version: %s\n", VERSION)
	fmt.Printf("BuildTime: %s\n", BuildTime)
	fmt.Printf("CommitHash: %s\n", CommitHash)
	fmt.Printf("GoVersion: %s\n", GoVersion)
	fmt.Printf("GitTag: %s\n", GitTag)
	fmt.Printf("GitBranch: %s\n", GitBranch)
}
