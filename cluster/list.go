package cluster

import (
	"github.com/mheers/k3i/ignite"
	"github.com/mheers/k3i/models"
)

func List() ([]*Cluster, error) {
	vms, err := ignite.GetAllClusterVMs()
	if err != nil {
		return nil, err
	}

	clusters, err := getClustersFromVMs(vms)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func ClusterMap() (map[string]*Cluster, error) {
	vms, err := ignite.GetAllClusterVMs()
	if err != nil {
		return nil, err
	}

	clusters, err := getClusterMapFromVMs(vms)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func getClusterMapFromVMs(vms []*models.VM) (map[string]*Cluster, error) {
	clusters := map[string]*Cluster{}
	for _, vm := range vms {

		clusterName := vm.Labels["clusterName"]
		if clusterName == "" {
			continue
		}

		node, err := convertVMToNode(vm)
		if err != nil {
			return nil, err
		}

		cluster, ok := clusters[clusterName]
		if !ok {
			cluster := &Cluster{
				Name:  clusterName,
				Nodes: []*Node{node},
			}
			clusters[clusterName] = cluster
		} else {
			cluster.Nodes = append(cluster.Nodes, node)
		}
	}

	return clusters, nil
}

func getClustersFromVMs(vms []*models.VM) ([]*Cluster, error) {
	clusters, err := getClusterMapFromVMs(vms)
	if err != nil {
		return nil, err
	}

	clusterList := []*Cluster{}
	for _, cluster := range clusters {
		clusterList = append(clusterList, cluster)
	}
	return clusterList, nil
}
