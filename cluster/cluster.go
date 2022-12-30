package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mheers/k3i/assets"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	Name                   string
	Status                 string
	Created                metav1.Time
	Nodes                  []*Node
	DesiredControllerCount int
	DesiredWorkerCount     int
}

func (c *Cluster) AddNode(node *Node) {
	c.Nodes = append(c.Nodes, node)
}

func (c *Cluster) RemoveNode(node *Node) {
	for i, n := range c.Nodes {
		if n.Name == node.Name {
			c.Nodes = append(c.Nodes[:i], c.Nodes[i+1:]...)
		}
	}
}

func (c *Cluster) GetNode(name string) *Node {
	for _, n := range c.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func (cluster *Cluster) installK0s() error {
	logrus.Infof("installing k0s on cluster '%s'", cluster.Name)

	for _, node := range cluster.Nodes {
		// copy k0s from assets
		paths := map[string]string{
			assets.K0sBinaryPath(): "/usr/local/bin/k0s",
		}
		err := node.CopyToMultiplePaths(paths)
		if err != nil {
			return err
		}
	}

	// join control plane nodes
	cns := cluster.ControllerNodes()
	for i, cn := range cns {
		if i == 0 {
			_, _, err := cn.Exec([]string{"k0s install controller"})
			if err != nil {
				return err
			}
		} else {
			err := cluster.writeTokenIntoNode(cn)
			if err != nil {
				return err
			}
			_, _, err = cn.Exec([]string{"k0s install controller --token-file /token"})
			if err != nil {
				return err
			}
			_, _, err = cn.Exec([]string{"k0s start"})
			if err != nil {
				return err
			}
		}
	}

	// join worker nodes
	wns := cluster.WorkerNodes()
	for _, wn := range wns {

		// copy k0s images from assets
		logrus.Infof("copying k0s images to node '%s'", wn.Name)
		paths := map[string]string{
			assets.ImagesPath(): "/var/lib/k0s/images/",
		}
		err := wn.CopyToMultiplePaths(paths)
		if err != nil {
			return err
		}

		// create a token for the worker node
		logrus.Infof("creating token for node '%s'", wn.Name)
		err = cluster.writeTokenIntoNode(wn)
		if err != nil {
			return err
		}

		// install k0s on the worker node
		logrus.Infof("installing k0s on node '%s'", wn.Name)
		_, _, err = wn.Exec([]string{"k0s install worker --token-file /token"})
		if err != nil {
			return err
		}

		// start k0s on the worker node
		logrus.Infof("starting k0s on node '%s'", wn.Name)
		_, _, err = wn.Exec([]string{"k0s start"})
		if err != nil {
			return err
		}
	}

	return nil
}

func (cluster *Cluster) ControllerNodes() []*Node {
	controllerNodes := []*Node{}
	for _, node := range cluster.Nodes {
		if node.Controller {
			controllerNodes = append(controllerNodes, node)
		}
	}
	return controllerNodes
}

func (cluster *Cluster) WorkerNodes() []*Node {
	workerNodes := []*Node{}
	for _, node := range cluster.Nodes {
		if node.Worker {
			workerNodes = append(workerNodes, node)
		}
	}
	return workerNodes
}

type Role int

const (
	Controller Role = iota
	Worker
)

func (cluster *Cluster) writeTokenIntoNode(node *Node) error {
	token, err := cluster.tokenForNode(node)
	if err != nil {
		return err
	}

	_, _, err = node.Exec([]string{fmt.Sprintf("echo '%s' > /token", token)})
	if err != nil {
		return err
	}

	return nil
}

func (cluster *Cluster) tokenForNode(node *Node) (string, error) {
	role := "worker"
	if node.Controller {
		role = "controller"
	}
	token, err := cluster.token(role)
	if err != nil {
		return "", err
	}

	return token, nil
}

type NodeStatus struct {
	Version    string
	ProcessID  int
	Role       string
	Workloads  bool
	SingleNode bool
}

func (cluster *Cluster) statusFromNode(node *Node) (NodeStatus, error) {
	stdout, stderr, err := node.Exec([]string{"k0s status"})
	if err != nil {
		return NodeStatus{}, err
	}
	if stderr != "" {
		return NodeStatus{}, fmt.Errorf("stderr: %s", stderr)
	}
	s := NodeStatus{}
	err = s.fromString(stdout)
	if err != nil {
		return NodeStatus{}, err
	}
	return s, nil
}

func (s *NodeStatus) fromString(status string) error {
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		value := strings.Split(line, ": ")[1]
		if strings.HasPrefix(line, "Version:") {
			s.Version = value
		}
		if strings.HasPrefix(line, "Process ID:") {
			pid, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			s.ProcessID = pid
		}
		if strings.HasPrefix(line, "Role:") {
			s.Role = value
		}
		if strings.HasPrefix(line, "Workloads:") {
			workloads, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			s.Workloads = workloads
		}
		if strings.HasPrefix(line, "SingleNode:") {
			singleNode, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			s.SingleNode = singleNode
		}
	}
	return nil
}

func (cluster *Cluster) token(role string) (string, error) {
	cns := cluster.ControllerNodes()
	initialControlNode := cns[0]

	status, err := cluster.statusFromNode(initialControlNode)
	if err != nil || status.ProcessID == 0 {
		logrus.Infof("k0s not running on %s, starting it", initialControlNode.Name)
		// start k0s as token creation requires k0s to be running
		_, _, err := initialControlNode.Exec([]string{"k0s start && sleep 5"})
		if err != nil {
			return "", err
		}
		status, err := cluster.statusFromNode(initialControlNode)
		if err != nil || status.ProcessID == 0 {
			return "", fmt.Errorf("k0s is still not running on %s", initialControlNode.Name)
		}
	}

	out, _, err := initialControlNode.Exec([]string{fmt.Sprintf("k0s token create --role=%s --expiry 10m", role)})
	if err != nil {
		return "", err
	}
	return out, nil
}
