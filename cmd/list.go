package cmd

import (
	"math"
	"os"

	"github.com/bradfitz/iter"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mheers/k3i/cluster"
	"github.com/mheers/k3i/helpers"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "list all clusters",
		Aliases: []string{"ls", "ps"},
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the log level
			helpers.SetLogLevel(LogLevelFlag)

			clusters, err := cluster.List()
			if err != nil {
				return err
			}

			return renderClusters(clusters)
		},
	}
)

func renderClusters(clusters []*cluster.Cluster) error {
	if OutputFormatFlag == "table" {
		renderListTable(clusters)
	}
	if OutputFormatFlag == "json" {
		err := helpers.PrintJSON(clusters)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "yaml" {
		err := helpers.PrintYAML(clusters)
		if err != nil {
			return err
		}
	}
	if OutputFormatFlag == "csv" {
		err := helpers.PrintCSV(clusters)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderListTable(clusters []*cluster.Cluster) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Status", "ControllerNodes", "WorkerNodes", "Created"})
	for _, cluster := range clusters {
		controllerNodes := cluster.ControllerNodes()
		workerNodes := cluster.WorkerNodes()

		maxRows := math.Max(float64(len(controllerNodes)), float64(len(workerNodes)))

		for i := range iter.N(int(maxRows)) {
			clusterName := ""
			if i == 0 {
				clusterName = cluster.Name
			}

			status := ""
			if i == 0 {
				status = "running"
			}

			controllerNodeName := ""
			if i < len(controllerNodes) {
				controllerNodeName = controllerNodes[i].Name
			}

			workerNodeName := ""
			if i < len(workerNodes) {
				workerNodeName = workerNodes[i].Name
			}

			createdTime := ""
			if i == 0 {
				createdTime = cluster.Nodes[0].Created.Time.String()
			}

			t.AppendRow(
				table.Row{
					clusterName,
					status,
					controllerNodeName,
					workerNodeName,
					createdTime,
				},
			)
		}
		t.AppendSeparator()
	}
	t.Render()
}
