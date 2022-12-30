package cluster

import "github.com/mheers/k3i/ignite"

func Delete(clusterName string) error {
	cluster, err := Get(clusterName)
	if err != nil {
		return err
	}

	for _, node := range cluster.Nodes {
		err := ignite.Delete(string(node.UID))
		if err != nil {
			return err
		}
	}

	return nil
}
