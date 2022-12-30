package cmd

import (
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	shellCmd = &cobra.Command{
		Use:   "shell [clusterName] [nodeName]",
		Short: "execs into a clusters node",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) != 2 {
				return cmd.Help()
			}
			clusterName := args[0]
			nodeName := args[1]
			err := cluster.Shell(clusterName, nodeName)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
