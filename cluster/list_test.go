package cluster

import (
	"testing"

	"github.com/mheers/k3i/ignite"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	err := ignite.InitClient()
	require.NoError(t, err)

	clusters, err := List()
	require.NoError(t, err)
	require.NotNil(t, clusters)
}
