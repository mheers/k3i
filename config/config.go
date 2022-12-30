package config

import "os"

type Config struct {
	SSHPrivateKeyFile string
	SSHPublicKeyFile  string
}

func GetConfig() *Config {
	return &Config{
		SSHPrivateKeyFile: "/home/marcel/.ssh/id_rsa",
		SSHPublicKeyFile:  "/home/marcel/.ssh/id_rsa.pub",
	}
}

func (config Config) GetSSHPublicKey() ([]byte, error) {
	return os.ReadFile(config.SSHPublicKeyFile)
}
