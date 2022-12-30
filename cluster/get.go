package cluster

import (
	"errors"
)

func Get(clusterName string) (*Cluster, error) {
	clusters, err := ClusterMap()
	if err != nil {
		return nil, err
	}

	cluster, ok := clusters[clusterName]
	if !ok {
		return nil, errors.New("cluster not found")
	}
	return cluster, nil
}
