package cluster

import (
	"fmt"

	"github.com/bradfitz/iter"
	"github.com/mheers/k3i/ignite"
)

type ClusterOptions struct {
	Name            string
	ControllerCount int
	WorkerCount     int
}

func Create(options ClusterOptions) (*Cluster, error) {
	clusterName := options.Name

	cluster := &Cluster{
		Name:                   clusterName,
		Nodes:                  []*Node{},
		DesiredControllerCount: options.ControllerCount,
		DesiredWorkerCount:     options.WorkerCount,
	}

	// create control plane nodes
	for i := range iter.N(cluster.DesiredControllerCount) {
		controlPlaneNode, err := createNode(clusterName, i, true)
		if err != nil {
			return nil, err
		}
		cluster.Nodes = append(cluster.Nodes, controlPlaneNode)
	}

	// create worker nodes
	for i := range iter.N(cluster.DesiredWorkerCount) {
		workerNode, err := createNode(clusterName, i, false)
		if err != nil {
			return nil, err
		}
		cluster.Nodes = append(cluster.Nodes, workerNode)
	}

	err := cluster.installK0s()
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func createNode(clusterName string, nodeNumber int, controller bool) (*Node, error) {
	labels := []string{
		fmt.Sprintf("clusterName=%s", clusterName),
		fmt.Sprintf("k3icluster=%t", true),
		fmt.Sprintf("controller=%t", controller),
		fmt.Sprintf("worker=%t", !controller),
		fmt.Sprintf("number=%d", nodeNumber),
	}
	typeString := "c"
	if !controller {
		typeString = "w"
	}
	nodeName := fmt.Sprintf("%s-%s-%d", clusterName, typeString, nodeNumber)
	node, err := ignite.Create(nodeName, labels)
	if err != nil {
		return nil, err
	}
	return convertVMToNode(node)
}
