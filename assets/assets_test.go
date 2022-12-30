package assets

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAirgapImages(t *testing.T) {
	AssetDir = "./"
	images, err := getAirgapImages()
	require.NoError(t, err)
	require.Len(t, images, 10)
}

func TestDownloadAirgapImages(t *testing.T) {
	AssetDir = "./"
	err := mkAssetDir()
	require.NoError(t, err)

	err = downloadAirgapImages()
	require.NoError(t, err)
}

func TestPullAndExportImage(t *testing.T) {
	AssetDir = "./"
	err := mkAssetDir()
	require.NoError(t, err)

	err = pullAndExportImage("docker.io/library/alpine:3.13")
	require.NoError(t, err)
}

func TestOK(t *testing.T) {
	AssetDir = "./"
	ok := OK()
	require.True(t, ok)
}

func TestImageTarOK(t *testing.T) {
	AssetDir = "./"
	ok := imageTarOK("images/docker.io_library_alpine:3.13.tar")
	require.True(t, ok)
	nok := imageTarOK("images/docker.io_library_alpine:3.13.tarx")
	require.False(t, nok)
	nok2 := imageTarOK("test/broken_alpine.tar")
	require.False(t, nok2)
}
