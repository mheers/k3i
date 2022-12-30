package cmd

import (
	"github.com/mheers/k3i/assets"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	assetsCmd = &cobra.Command{
		Use:   "assets [name]",
		Short: "assets downloads the assets for a cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			err := assets.Download()
			if err != nil {
				return err
			}

			return nil
		},
	}
)
