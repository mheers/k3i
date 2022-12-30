package cluster

import (
	"fmt"
	"strings"

	"github.com/mheers/k3i/ignite"
	"github.com/sirupsen/logrus"
)

// Copy copies a file from the host to each node in the cluster
func Copy(src, dest string) error {
	clusterName, err := getClusterName(src, dest)
	if err != nil {
		return err
	}

	isNodeName, node, err := isNodeName(clusterName)
	if err != nil {
		return err
	}

	if isNodeName {
		err = copy(node, src, dest)
		if err != nil {
			return err
		}
		return nil
	}

	cluster, err := Get(clusterName)
	if err != nil {
		return err
	}
	if len(cluster.Nodes) == 0 {
		return fmt.Errorf("no nodes found for cluster %s", clusterName)
	}

	for _, node := range cluster.Nodes {
		err = copy(node, src, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

func copy(node *Node, src, dest string) error {
	logrus.Debugf("copying %s to %s on node %s", src, dest, node.Name)
	vmSrc, vmDest, err := getVMSrcDest(node, src, dest)
	if err != nil {
		return err
	}
	err = ignite.Copy(&node.VM.VM, vmSrc, vmDest)
	if err != nil {
		return err
	}
	return nil
}

func getVMSrcDest(node *Node, src, dest string) (string, string, error) {
	clusterName, err := getClusterName(src, dest)
	if err != nil {
		return "", "", err
	}

	nodePath := fmt.Sprintf("%s://", node.GetUID())
	srcReplaced := strings.Replace(src, clusterName+"://", nodePath, 1)
	destReplaced := strings.Replace(dest, clusterName+"://", nodePath, 1)

	return srcReplaced, destReplaced, nil
}

func getClusterName(src, dest string) (string, error) {
	srcCluster, srcIsCluster := detectClusterName(src)
	destCluster, destIsCluster := detectClusterName(dest)

	if srcIsCluster && destIsCluster {
		return "", fmt.Errorf("copying between clusters is not supported")
	}

	if srcIsCluster {
		return srcCluster, nil
	}

	if destIsCluster {
		return destCluster, nil
	}

	return "", fmt.Errorf("no cluster name found in src or dest")
}

func detectClusterName(src string) (string, bool) {
	if strings.Index(src, "://") > 0 {
		parts := strings.Split(src, "://")
		if len(parts) == 2 {
			return parts[0], true
		}
	}
	return "", false
}

func isNodeName(name string) (bool, *Node, error) {
	vms, err := ignite.GetAllClusterVMs()
	if err != nil {
		return false, nil, err
	}

	for _, vm := range vms {
		if vm.Name == name {
			node, err := convertVMToNode(vm)
			if err != nil {
				return false, nil, err
			}
			return true, node, nil
		}
	}
	return false, nil, nil
}
