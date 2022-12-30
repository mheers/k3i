package cluster

import (
	"fmt"

	"github.com/mheers/k3i/ignite"
)

func Shell(clusterName, nodeName string) error {
	cluster, err := Get(clusterName)
	if err != nil {
		return err
	}
	if len(cluster.Nodes) == 0 {
		return fmt.Errorf("no nodes found for cluster %s", clusterName)
	}

	for _, node := range cluster.Nodes {
		if node.Name == nodeName {
			return ignite.Shell(&node.VM.VM)
		}
	}

	return fmt.Errorf("could not find node '%s' in cluster '%s'", nodeName, clusterName)
}
