package entities

import (
	"io"

	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/cli/internal/agent"
	"github.com/eleven-sh/eleven/entities"
)

func ReconcileServedPortsState(
	env *entities.Env,
	agentClientBuilder agent.ClientBuilder,
) error {

	servedPorts := BuildProtoEnvServedPortsFromEnv(env)

	agentClient := agentClientBuilder.Build(
		agent.NewDefaultClientConfig(
			[]byte(env.SSHKeyPairPEMContent),
			env.InstancePublicIPAddress,
		),
	)

	return agentClient.ReconcileServedPortsState(
		&proto.ReconcileServedPortsStateRequest{
			ServedPorts: servedPorts,
		},
		func(stream agent.ReconcileServedPortsStateStream) error {

			for {
				_, err := stream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func BuildProtoEnvServedPortsFromEnv(
	env *entities.Env,
) map[string]*proto.EnvServedPortBindings {

	portsToServe := map[string]*proto.EnvServedPortBindings{}

	for port, bindings := range env.ServedPorts {
		allBindings := []*proto.EnvServedPortBinding{}

		for _, binding := range bindings {
			allBindings = append(allBindings, &proto.EnvServedPortBinding{
				Value:           binding.Value,
				Type:            string(binding.Type),
				RedirectToHttps: binding.RedirectToHTTPS,
			})
		}

		portsToServe[string(port)] = &proto.EnvServedPortBindings{
			Bindings: allBindings,
		}
	}

	return portsToServe
}
