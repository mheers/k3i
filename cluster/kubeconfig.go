package cluster

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
)

func Kubeconfig(clusterName string, write bool) (string, error) {
	cluster, err := Get(clusterName)
	if err != nil {
		return "", err
	}

	kc, err := cluster.kubeconfig()
	if err != nil {
		return "", err
	}

	kubeconfig, err := cluster.adjustKubeconfig(kc)
	if err != nil {
		return "", err
	}

	if write {
		err = writeKubeconfig(kubeconfig)
		if err != nil {
			return "", err
		}
	}

	return kubeconfig, nil
}

func (cluster *Cluster) kubeconfig() (string, error) {
	cns := cluster.ControllerNodes()
	initialControlNode := cns[0]

	// start k0s as getting kubeconfig requires k0s to be running
	_, _, err := initialControlNode.Exec([]string{"k0s start"})
	if err != nil {
		return "", err
	}

	out, _, err := initialControlNode.Exec([]string{"k0s kubeconfig admin"})
	if err != nil {
		return "", err
	}

	return out, nil
}

func (cluster *Cluster) adjustKubeconfig(kubeconfig string) (string, error) {
	clusterName := cluster.Name

	logrus.Debug("-------------")
	logrus.Debug(kubeconfig)
	logrus.Debug("-------------")

	kubeconfigParsed, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		return "", err
	}

	logrus.Debug(kubeconfigParsed.Clusters)

	currentClusterNames := make([]string, 0)
	for cn := range kubeconfigParsed.Clusters {
		currentClusterNames = append(currentClusterNames, cn)
	}

	defaultClusterName := currentClusterNames[0]
	defaultContextName := kubeconfigParsed.CurrentContext

	// set cluster name
	kubeconfigParsed.Contexts[defaultContextName].Cluster = clusterName
	kubeconfigParsed.Clusters[clusterName] = kubeconfigParsed.Clusters[defaultClusterName]
	delete(kubeconfigParsed.Clusters, defaultClusterName)
	result, err := clientcmd.Write(*kubeconfigParsed)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func writeKubeconfig(kubeconfig string) error {
	kubeconfigParsed, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		return err
	}
	clusterNames := make([]string, 0)
	for clusterName := range kubeconfigParsed.Clusters {
		clusterNames = append(clusterNames, clusterName)
	}
	clusterName := clusterNames[0]

	// write kubeconfig to file
	outputFile := fmt.Sprintf("kubeconfig.%s.yaml", clusterName)
	err = os.WriteFile(outputFile, []byte(kubeconfig), 0644)
	if err != nil {
		return err
	}

	logrus.Infof("kubeconfig written to file: %s", outputFile)

	return nil
}
