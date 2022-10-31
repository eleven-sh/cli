package agent

import (
	"net"

	"github.com/eleven-sh/agent/config"
)

type ClientBuilder interface {
	Build(config ClientConfig) ClientInterface
}

type DefaultClientBuilder struct{}

func NewDefaultClientBuilder() DefaultClientBuilder {
	return DefaultClientBuilder{}
}

func (DefaultClientBuilder) Build(config ClientConfig) ClientInterface {
	return NewClient(config)
}

func NewDefaultClientConfig(
	sshPrivateKeyBytes []byte,
	instancePublicIPAddress string,
) ClientConfig {

	return ClientConfig{
		ServerRootUser:           config.ElevenUserName,
		ServerSSHPrivateKeyBytes: sshPrivateKeyBytes,
		ServerAddr: net.JoinHostPort(
			instancePublicIPAddress,
			config.SSHServerListenPort,
		),
		// Ends with ":" to let "net.listener"
		// choose a random available port for us
		LocalAddr:          "127.0.0.1:",
		RemoteAddrProtocol: config.GRPCServerAddrProtocol,
		RemoteAddr:         config.GRPCServerAddr,
	}
}
