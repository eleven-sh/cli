package agent

import (
	"context"

	"github.com/eleven-sh/agent/proto"
)

type InstallRuntimesStream interface {
	Recv() (*proto.InstallRuntimesReply, error)
}

func (c Client) InstallRuntimes(
	installRuntimesRequest *proto.InstallRuntimesRequest,
	streamHandler func(stream InstallRuntimesStream) error,
) error {

	return c.Execute(func(agentGRPCClient proto.AgentClient) error {
		installRuntimesStream, err := agentGRPCClient.InstallRuntimes(
			context.TODO(),
			installRuntimesRequest,
		)

		if err != nil {
			return err
		}

		return streamHandler(installRuntimesStream)
	})
}
