package agent

import (
	"context"

	"github.com/eleven-sh/agent/proto"
)

type ReconcileServedPortsStateStream interface {
	Recv() (*proto.ReconcileServedPortsStateReply, error)
}

func (c Client) ReconcileServedPortsState(
	serveRequest *proto.ReconcileServedPortsStateRequest,
	streamHandler func(stream ReconcileServedPortsStateStream) error,
) error {

	return c.Execute(func(agentGRPCClient proto.AgentClient) error {
		serveStream, err := agentGRPCClient.ReconcileServedPortsState(
			context.TODO(),
			serveRequest,
		)

		if err != nil {
			return err
		}

		return streamHandler(serveStream)
	})
}
