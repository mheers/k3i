package cluster

import (
	"fmt"
	"strconv"

	"github.com/mheers/k3i/ignite"
	"github.com/mheers/k3i/models"
	"github.com/sirupsen/logrus"
)

type Node struct {
	models.VM
	Controller bool
	Worker     bool
	Number     int
}

func convertVMToNode(vm *models.VM) (*Node, error) {
	labels := vm.Labels
	controller := labels["controller"] == "true"
	worker := labels["worker"] == "true"
	number, err := strconv.Atoi(labels["number"])
	if err != nil {
		return nil, fmt.Errorf("could not read node number from node '%s': %s", vm.Name, err)
	}

	return &Node{
		VM:         *vm,
		Controller: controller,
		Worker:     worker,
		Number:     number,
	}, nil
}

func (node Node) Exec(args []string) (string, string, error) {
	logrus.Infof("executing command on node %s: %v", node.Name, args)
	return ignite.Exec(&node.VM.VM, args)
}

func (node Node) CopyToMultiplePaths(paths map[string]string) error {
	for local, remote := range paths {
		err := node.CopyTo(local, remote)
		if err != nil {
			return err
		}
	}
	return nil
}

func (node Node) CopyTo(local, remote string) error {
	remote = node.pathToNodePath(remote)
	return ignite.Copy(&node.VM.VM, local, remote)
}

func (node Node) CopyFrom(remote, local string) error {
	remote = node.pathToNodePath(remote)
	return ignite.Copy(&node.VM.VM, remote, local)
}

func (node Node) pathToNodePath(path string) string {
	return fmt.Sprintf("%s://%s", node.VM.VM.GetUID(), path)
}
