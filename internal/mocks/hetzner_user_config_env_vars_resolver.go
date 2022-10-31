// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eleven-sh/cli/internal/cloudproviders/hetzner (interfaces: UserConfigEnvVarsResolver)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	userconfig "github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	gomock "github.com/golang/mock/gomock"
)

// HetznerUserConfigEnvVarsResolver is a mock of UserConfigEnvVarsResolver interface.
type HetznerUserConfigEnvVarsResolver struct {
	ctrl     *gomock.Controller
	recorder *HetznerUserConfigEnvVarsResolverMockRecorder
}

// HetznerUserConfigEnvVarsResolverMockRecorder is the mock recorder for HetznerUserConfigEnvVarsResolver.
type HetznerUserConfigEnvVarsResolverMockRecorder struct {
	mock *HetznerUserConfigEnvVarsResolver
}

// NewHetznerUserConfigEnvVarsResolver creates a new mock instance.
func NewHetznerUserConfigEnvVarsResolver(ctrl *gomock.Controller) *HetznerUserConfigEnvVarsResolver {
	mock := &HetznerUserConfigEnvVarsResolver{ctrl: ctrl}
	mock.recorder = &HetznerUserConfigEnvVarsResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *HetznerUserConfigEnvVarsResolver) EXPECT() *HetznerUserConfigEnvVarsResolverMockRecorder {
	return m.recorder
}

// Resolve mocks base method.
func (m *HetznerUserConfigEnvVarsResolver) Resolve() (*userconfig.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve")
	ret0, _ := ret[0].(*userconfig.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve.
func (mr *HetznerUserConfigEnvVarsResolverMockRecorder) Resolve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*HetznerUserConfigEnvVarsResolver)(nil).Resolve))
}
