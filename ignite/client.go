package ignite

import (
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/providers"
	cniprovider "github.com/weaveworks/ignite/pkg/providers/cni"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func InitClient() error {
	providers.RuntimeName = "containerd"
	providers.NetworkPluginName = "cni"

	if err := config.SetAndPopulateProviders(providers.RuntimeName, providers.NetworkPluginName); err != nil {
		return err
	}

	err := providers.Populate(ignite.Preload)
	if err != nil {
		return err
	}

	err = cniprovider.SetCNINetworkPlugin()
	if err != nil {
		return err
	}

	return nil
}
