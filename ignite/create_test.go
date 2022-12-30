package ignite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	err := InitClient()
	require.NoError(t, err)

	vm, err := Create("some-name", nil)
	require.NoError(t, err)
	require.NotNil(t, vm)
}
