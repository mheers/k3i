package cmd

import (
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get [name]",
		Short: "gets a cluster",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			if len(args) == 0 {
				return cmd.Help()
			}
			name := args[0]
			cluster, err := cluster.Get(name)
			if err != nil {
				return err
			}

			return renderCluster(cluster)
		},
	}
)

func renderCluster(c *cluster.Cluster) error {
	if OutputFormatFlag == "table" {
		renderListTable([]*cluster.Cluster{c})
	}
	if OutputFormatFlag == "json" {
		err := helpers.PrintJSON(c)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "yaml" {
		err := helpers.PrintYAML(c)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "csv" {
		err := helpers.PrintCSV(c)
		if err != nil {
			return err
		}
	}
	return nil
}
