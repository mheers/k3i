package main

import (
	"github.com/mheers/k3i/cmd"
	"github.com/mheers/k3i/ignite"
	"github.com/sirupsen/logrus"
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

func main() {
	cmd.VERSION = VERSION
	cmd.BuildTime = BuildTime
	cmd.CommitHash = CommitHash
	cmd.GoVersion = GoVersion
	cmd.GitTag = GitTag
	cmd.GitBranch = GitBranch

	err := ignite.InitClient()
	if err != nil {
		logrus.Fatalf("InitClient failed: %+v", err)
	}

	// execeute the command
	err = cmd.Execute()
	if err != nil {
		logrus.Fatalf("Execute failed: %+v", err)
	}
}
