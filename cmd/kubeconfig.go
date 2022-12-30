package cmd

import (
	"fmt"

	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	// writeKubeconfig is the flag for writing the kubeconfig to a file
	writeKubeconfig bool

	kubeconfigCmd = &cobra.Command{
		Use:   "kubeconfig [name]",
		Short: "gets the kubeconfig of a cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) == 0 {
				return cmd.Help()
			}
			name := args[0]
			kubeconfig, err := cluster.Kubeconfig(name, writeKubeconfig)
			if err != nil {
				return err
			}

			fmt.Println(kubeconfig)
			return nil
		},
	}
)

func init() {
	kubeconfigCmd.PersistentFlags().BoolVarP(&writeKubeconfig, "write", "w", false, "write the kubeconfig to a file")
}
