package ignite

import (
	"github.com/mheers/k3i/models"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
)

func GetAllVMs() ([]*models.VM, error) {
	allVMs, err := allVMs()
	if err != nil {
		return nil, err
	}

	return convertIgnitesToVMs(allVMs), nil
}

func allVMs() ([]*ignite.VM, error) {
	allVMs, err := providers.Client.VMs().FindAll(filter.NewVMFilterAll("", true))
	if err != nil {
		return nil, err
	}
	return allVMs, nil
}

func convertIgnitesToVMs(allVMs []*ignite.VM) []*models.VM {
	vms := []*models.VM{}
	for _, vm := range allVMs {
		vms = append(vms, &models.VM{
			VM: *vm,
		})
	}
	return vms
}

func GetAllClusterVMs() ([]*models.VM, error) {
	allVMs, err := allVMs()
	if err != nil {
		return nil, err
	}

	requiredLabel := "k3icluster"
	filteredVMs := []*ignite.VM{}
	for _, vm := range allVMs {
		if vm.Labels[requiredLabel] == "true" {
			filteredVMs = append(filteredVMs, vm)
		}
	}

	return convertIgnitesToVMs(filteredVMs), nil
}
