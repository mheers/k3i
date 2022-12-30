package cmd

import (
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:     "delete [name]",
		Short:   "deletes a cluster",
		Aliases: []string{"del", "rm", "remove"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) == 0 {
				return cmd.Help()
			}
			name := args[0]
			err := cluster.Delete(name)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
