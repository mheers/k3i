package ignite

import (
	"github.com/mheers/k3i/config"
	"github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
)

func Shell(vm *ignite.VM) error {
	config := config.GetConfig()

	stdout, stderr, err := runSSH(vm, config.SSHPrivateKeyFile, []string{}, true, constants.SSH_DEFAULT_TIMEOUT_SECONDS)

	defer func() {
		logrus.Println(stdout)
		logrus.Error(stderr)
	}()

	if err != nil && err.Error() != "EOF" {
		return err
	}

	return nil
}
