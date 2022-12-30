package ignite

import (
	"github.com/mheers/k3i/config"
	"github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
)

func Exec(vm *ignite.VM, args []string) (string, string, error) {
	config := config.GetConfig()

	stdout, stderr, err := runSSH(vm, config.SSHPrivateKeyFile, args, false, constants.SSH_DEFAULT_TIMEOUT_SECONDS)

	defer func() {
		logrus.Debug(stdout)
		logrus.Debug(stderr)
	}()

	if err != nil && err.Error() != "EOF" {
		return "", "", err
	}

	return stdout, stderr, nil
}
