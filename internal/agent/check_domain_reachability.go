package agent

import (
	"context"

	"github.com/eleven-sh/agent/proto"
)

type CheckDomainReachabilityStream interface {
	Recv() (*proto.CheckDomainReachabilityReply, error)
}

func (c Client) CheckDomainReachability(
	checkDomainReachabilityRequest *proto.CheckDomainReachabilityRequest,
	streamHandler func(stream CheckDomainReachabilityStream) error,
) error {

	return c.Execute(func(agentGRPCClient proto.AgentClient) error {
		checkDomainReachabilityStream, err := agentGRPCClient.CheckDomainReachability(
			context.TODO(),
			checkDomainReachabilityRequest,
		)

		if err != nil {
			return err
		}

		return streamHandler(checkDomainReachabilityStream)
	})
}
