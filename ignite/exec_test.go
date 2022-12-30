package ignite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	err := InitClient()
	require.NoError(t, err)

	clusterVMs, err := GetAllClusterVMs()
	require.NoError(t, err)
	require.NotEmpty(t, clusterVMs)

	for _, vm := range clusterVMs {
		stdout, stderr, err := Exec(&vm.VM, []string{"echo", "-n", "hello"})
		require.NoError(t, err)
		require.Equal(t, "hello", stdout)
		require.Equal(t, "", stderr)
	}
}

func TestExecPipe(t *testing.T) {
	err := InitClient()
	require.NoError(t, err)

	clusterVMs, err := GetAllClusterVMs()
	require.NoError(t, err)
	require.NotEmpty(t, clusterVMs)

	for _, vm := range clusterVMs {
		stdout, stderr, err := Exec(&vm.VM, []string{"ls / | grep etc"})
		require.NoError(t, err)
		require.Equal(t, "etc\n", stdout)
		require.Equal(t, "", stderr)
	}
}
