package cmd

import (
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	copyCmd = &cobra.Command{
		Use:     "copy [src] [dest]",
		Short:   "copies files from and to a clusters node",
		Aliases: []string{"cp"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) < 2 {
				return cmd.Help()
			}
			src := args[0]
			dest := args[1]
			err := cluster.Copy(src, dest)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
