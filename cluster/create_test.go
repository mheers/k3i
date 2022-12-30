package cluster

import (
	"testing"

	"github.com/mheers/k3i/ignite"
	"github.com/stretchr/testify/require"
)

func TestInstallK0s(t *testing.T) {
	err := ignite.InitClient()
	require.NoError(t, err)

	clusterName := "some-name"
	possibleNodes, err := ignite.GetAllClusterVMs()
	require.NoError(t, err)
	require.NotEmpty(t, possibleNodes)

	node1 := possibleNodes[0]

	cluster := &Cluster{
		Name: clusterName,
		Nodes: []*Node{
			{
				VM: *node1,
			},
		},
	}

	err = cluster.installK0s()
	require.NoError(t, err)
}
