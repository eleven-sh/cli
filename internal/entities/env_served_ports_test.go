package entities

import (
	"reflect"
	"testing"

	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/eleven/entities"
)

func TestBuildProtoEnvServedPortsFromEnv(t *testing.T) {
	testCases := []struct {
		test                     string
		servedPorts              entities.EnvServedPorts
		expectedProtoServedPorts map[string]*proto.EnvServedPortBindings
	}{
		{
			test: "with empty served ports",

			servedPorts: entities.EnvServedPorts{},

			expectedProtoServedPorts: map[string]*proto.EnvServedPortBindings{},
		},

		{
			test: "with empty served port bindings",

			servedPorts: entities.EnvServedPorts{
				"8000": {},
			},

			expectedProtoServedPorts: map[string]*proto.EnvServedPortBindings{
				"8000": {
					Bindings: []*proto.EnvServedPortBinding{},
				},
			},
		},

		{
			test: "with mixed served ports",

			servedPorts: entities.EnvServedPorts{
				"8000": {
					{
						Value:           "api.eleven.sh",
						Type:            entities.EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},

					{
						Value:           "8000",
						Type:            entities.EnvServedPortBindingTypePort,
						RedirectToHTTPS: false,
					},
				},

				"6000": {
					{
						Value:           "6000",
						Type:            entities.EnvServedPortBindingTypePort,
						RedirectToHTTPS: false,
					},
				},

				"4000": {
					{
						Value:           "test.test.io",
						Type:            entities.EnvServedPortBindingTypeDomain,
						RedirectToHTTPS: true,
					},
				},
			},

			expectedProtoServedPorts: map[string]*proto.EnvServedPortBindings{
				"8000": {
					Bindings: []*proto.EnvServedPortBinding{
						{
							Type:            "domain",
							Value:           "api.eleven.sh",
							RedirectToHttps: true,
						},

						{
							Type:            "port",
							Value:           "8000",
							RedirectToHttps: false,
						},
					},
				},

				"6000": {
					Bindings: []*proto.EnvServedPortBinding{
						{
							Type:            "port",
							Value:           "6000",
							RedirectToHttps: false,
						},
					},
				},

				"4000": {
					Bindings: []*proto.EnvServedPortBinding{
						{
							Type:            "domain",
							Value:           "test.test.io",
							RedirectToHttps: true,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			env := entities.NewEnv(
				"test_env",
				0,
				"test_instance_type",
				[]entities.EnvRepository{},
				entities.EnvRuntimes{},
			)

			env.ServedPorts = tc.servedPorts

			protoServedPorts := BuildProtoEnvServedPortsFromEnv(env)

			if !reflect.DeepEqual(tc.expectedProtoServedPorts, protoServedPorts) {
				t.Fatalf(
					"expected proto served ports to equal '%+v', got '%+v'",
					tc.expectedProtoServedPorts,
					protoServedPorts,
				)
			}
		})
	}
}
