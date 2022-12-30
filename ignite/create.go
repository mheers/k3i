package ignite

import (
	"github.com/mheers/k3i/config"
	"github.com/mheers/k3i/models"
	"github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	// DefaultImage is the default image to use for the VM
	DefaultImage = "mheers/ignite-ubuntu-base:v1.0.5" // starting k3s fails after a short while
	// DefaultImage = "mheers/ignite-ubuntu-base-containerd:v1.0.5"
	// DefaultImage = "mheers/ignite-ubuntu-base-docker:v1.0.5"
	// DefaultImage = "mheers/ignite-alpine-base:v1.0.5" // does not have working ssh

	// SandboxImage is the image to use for the sandbox
	SandboxImage = "weaveworks/ignite:v0.10.0-amd64"
)

func Create(name string, labels []string) (*models.VM, error) {
	Name := name

	CPUs := uint64(2)
	KernelCmdLine := ""
	MemoryGB := 2
	DiskSizeGB := 10
	// KernelOCI := ""
	// SandboxOCI := SandboxImage
	// Storage := ""
	Labels := []string{}
	Labels = append(Labels, labels...)

	logrus.Infof("creating node %s with labels %v", name, labels)

	baseVM := providers.Client.VMs().New()
	baseVM.ObjectMeta.Name = Name

	baseVM.Status.Runtime.Name = providers.RuntimeName
	baseVM.Status.Network.Plugin = providers.NetworkPluginName

	imageOciRef, err := meta.NewOCIImageRef(DefaultImage)
	if err != nil {
		return nil, err
	}
	baseVM.Spec.Image.OCI = imageOciRef

	sandboxImageOciRef, err := meta.NewOCIImageRef(SandboxImage)
	if err != nil {
		return nil, err
	}
	baseVM.Spec.Sandbox.OCI = sandboxImageOciRef

	// Generate a VM name and UID if not set yet.
	if err := metadata.SetNameAndUID(baseVM, providers.Client); err != nil {
		return nil, err
	}

	// Validate the VM object.
	if err := validation.ValidateVM(baseVM).ToAggregate(); err != nil {
		return nil, err
	}

	// Get the image, or import it if it doesn't exist.
	image, err := operations.FindOrImportImage(providers.Client, baseVM.Spec.Image.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Image on the VM object.
	baseVM.SetImage(image)

	// Get the kernel, or import it if it doesn't exist.
	kernel, err := operations.FindOrImportKernel(providers.Client, baseVM.Spec.Kernel.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Kernel on the VM object.
	baseVM.SetKernel(kernel)

	// Get the sandbox, or import it if it doesn't exist.
	sandbox, err := operations.FindOrImportImage(providers.Client, baseVM.Spec.Sandbox.OCI)
	if err != nil {
		return nil, err
	}

	setSandBoxImage := func(vm *ignite.VM, sandbox *ignite.Image) error {
		// Populate relevant data from the Sandbox on the VM object.
		vm.Spec.Sandbox.OCI = sandbox.Spec.OCI
		return nil
	}

	// Set the sandbox image
	if err := setSandBoxImage(baseVM, sandbox); err != nil {
		return nil, err
	}

	baseVM.Spec.CPUs = CPUs
	baseVM.Spec.Kernel.CmdLine = KernelCmdLine
	baseVM.Spec.Memory = meta.NewSizeFromBytes(uint64(MemoryGB) * 1024 * 1024 * 1024)
	baseVM.Spec.DiskSize = meta.NewSizeFromBytes(uint64(DiskSizeGB) * 1024 * 1024 * 1024)
	// baseVM.Spec.Kernel.OCI = KernelOCI
	// baseVM.Spec.Sandbox.OCI = SandboxOCI
	// baseVM.Spec.Storage = Storage

	// TODO: Volumes
	// TODO: CopyFiles
	// TODO: PortMappings

	config := config.GetConfig()
	// publicKey, err := config.GetSSHPublicKey()
	baseVM.Spec.SSH = &ignite.SSH{
		Generate:  false,
		PublicKey: config.SSHPublicKeyFile,
	}

	// Generate a random UID and Name
	if err = metadata.SetNameAndUID(baseVM, providers.Client); err != nil {
		return nil, err
	}
	// Set VM labels.
	if err = metadata.SetLabels(baseVM, Labels); err != nil {
		return nil, err
	}
	defer util.DeferErr(&err, func() error { return metadata.Cleanup(baseVM, false) })

	if err = providers.Client.VMs().Set(baseVM); err != nil {
		return nil, err
	}

	// Allocate and populate the overlay file
	if err = dmlegacy.AllocateAndPopulateOverlay(baseVM); err != nil {
		return nil, err
	}

	if err = metadata.Success(baseVM); err != nil {
		return nil, err
	}

	debug := true
	if err = operations.StartVM(baseVM, debug); err != nil {
		return nil, err
	}

	if err = waitForSSH(baseVM, constants.SSH_DEFAULT_TIMEOUT_SECONDS, constants.IGNITE_SPAWN_TIMEOUT); err != nil {
		return nil, err
	}

	// get the VM object again to get the updated status
	baseVMCreated, err := providers.Client.VMs().Get(baseVM.GetUID())
	if err != nil {
		return nil, err
	}

	vm := &models.VM{
		VM: *baseVMCreated,
	}
	return vm, nil
}
