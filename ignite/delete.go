package ignite

import (
	"errors"

	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/runtime"
)

func Delete(uid string) error {
	// get the vm
	vm, err := providers.Client.VMs().Get(runtime.UID(uid))
	if err != nil {
		return err
	}

	if vm == nil {
		return errors.New("vm not found")
	}

	// stop the vm
	err = operations.StopVM(vm, true, true)
	if err != nil {
		return err
	}

	// delete the vm
	return operations.DeleteVM(providers.Client, vm)
}
