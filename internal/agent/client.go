package agent

import (
	"time"

	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/cli/internal/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientInterface interface {
	InitInstance(
		initInstanceRequest *proto.InitInstanceRequest,
		streamHandler func(stream InitInstanceStream) error,
	) error

	InstallRuntimes(
		installRuntimesRequest *proto.InstallRuntimesRequest,
		streamHandler func(stream InstallRuntimesStream) error,
	) error

	ReconcileServedPortsState(
		serveRequest *proto.ReconcileServedPortsStateRequest,
		streamHandler func(stream ReconcileServedPortsStateStream) error,
	) error

	CheckDomainReachability(
		checkDomainReachabilityRequest *proto.CheckDomainReachabilityRequest,
		streamHandler func(stream CheckDomainReachabilityStream) error,
	) error
}

type ClientConfig struct {
	ServerRootUser           string
	ServerSSHPrivateKeyBytes []byte
	ServerAddr               string
	LocalAddr                string
	RemoteAddrProtocol       string
	RemoteAddr               string
}

type Client struct {
	config ClientConfig
}

func NewClient(config ClientConfig) Client {
	return Client{
		config: config,
	}
}

func (c Client) Execute(fnToExec func(agentGRPCClient proto.AgentClient) error) error {
	pollTimeoutChan := time.After(48 * time.Second)
	pollSleepDuration := 4 * time.Second

	var portForwarderReadyResp ssh.PortForwarderReadyResp
	var portForwarderRespChan chan error

	for {
		select {
		case <-pollTimeoutChan:
			return portForwarderReadyResp.Error
		default:
			portForwarderRespChan = make(chan error, 1)
			portForwarderReadyChan := make(chan ssh.PortForwarderReadyResp)
			portForwarder := ssh.NewPortForwarder()

			// Open an SSH tunnel to the GRPC server
			// from "localAddr" to "remoteAddr" inside server ("serverAddr")
			go func() {
				portForwarderRespChan <- portForwarder.Forward(
					portForwarderReadyChan,
					c.config.ServerSSHPrivateKeyBytes,
					c.config.ServerRootUser,
					c.config.ServerAddr,
					c.config.LocalAddr,
					c.config.RemoteAddrProtocol,
					c.config.RemoteAddr,
				)
			}()

			portForwarderReadyResp = <-portForwarderReadyChan
		}

		if portForwarderReadyResp.Error == nil {
			break
		}

		time.Sleep(pollSleepDuration)
	}

	grpcRespChan := make(chan error, 1)

	go func() {
		grpcConn, err := grpc.Dial(
			portForwarderReadyResp.LocalAddr,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)

		if err != nil {
			grpcRespChan <- err
			return
		}

		defer grpcConn.Close()

		agentGRPCClient := proto.NewAgentClient(grpcConn)

		grpcRespChan <- fnToExec(agentGRPCClient)
	}()

	select {
	case portForwarderErr := <-portForwarderRespChan:
		return portForwarderErr
	case grpcErr := <-grpcRespChan:
		return grpcErr
	}
}
