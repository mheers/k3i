package cmd

import (
	"github.com/mheers/k3i/assets"
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	// ControllerCountFlag is the flag for the number of controller nodes
	ControllerCountFlag int

	// WorkerCountFlag is the flag for the number of worker nodes
	WorkerCountFlag int

	createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "creates a cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) == 0 {
				return cmd.Help()
			}

			if !assets.OK() {
				err := assets.Download()
				if err != nil {
					return err
				}
			}

			name := args[0]

			c, err := cluster.Create(cluster.ClusterOptions{
				Name:            name,
				ControllerCount: ControllerCountFlag,
				WorkerCount:     WorkerCountFlag,
			})
			if err != nil {
				return err
			}

			return renderCluster(c)
		},
	}
)

func init() {
	createCmd.PersistentFlags().IntVarP(&ControllerCountFlag, "controller-count", "c", 1, "number of controller nodes")
	createCmd.PersistentFlags().IntVarP(&WorkerCountFlag, "worker-count", "w", 1, "number of worker nodes")
}
