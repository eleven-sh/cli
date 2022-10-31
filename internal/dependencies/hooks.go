// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/eleven-sh/cli/internal/hooks"
	"github.com/google/wire"
)

func ProvidePreRemoveHook() hooks.PreRemove {
	panic(
		wire.Build(
			sshConfigManagerSet,

			sshKnownHostsManagerSet,

			sshKeysManagerSet,

			userConfigManagerSet,

			githubManagerSet,

			hooks.NewPreRemove,
		),
	)
}

func ProvideDomainReachabilityChecker() hooks.DomainReachabilityChecker {
	panic(
		wire.Build(
			agentSet,

			hooks.NewDomainReachabilityChecker,
		),
	)
}
