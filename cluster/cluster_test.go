package cluster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeStatusFromString(t *testing.T) {
	value := `Version: v1.25.4+k0s.0
Process ID: 1577
Role: controller
Workloads: false
SingleNode: false
`
	s := NodeStatus{}
	require.Equal(t, 0, s.ProcessID)

	err := s.fromString(value)
	require.NoError(t, err)
	require.Equal(t, "v1.25.4+k0s.0", s.Version)
	require.Equal(t, 1577, s.ProcessID)
	require.Equal(t, "controller", s.Role)
	require.Equal(t, false, s.Workloads)
	require.Equal(t, false, s.SingleNode)
}
