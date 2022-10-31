// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eleven-sh/cli/internal/cloudproviders/hetzner (interfaces: UserConfigFilesResolver)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	userconfig "github.com/eleven-sh/hetzner-cloud-provider/userconfig"
	gomock "github.com/golang/mock/gomock"
)

// HetznerUserConfigFilesResolver is a mock of UserConfigFilesResolver interface.
type HetznerUserConfigFilesResolver struct {
	ctrl     *gomock.Controller
	recorder *HetznerUserConfigFilesResolverMockRecorder
}

// HetznerUserConfigFilesResolverMockRecorder is the mock recorder for HetznerUserConfigFilesResolver.
type HetznerUserConfigFilesResolverMockRecorder struct {
	mock *HetznerUserConfigFilesResolver
}

// NewHetznerUserConfigFilesResolver creates a new mock instance.
func NewHetznerUserConfigFilesResolver(ctrl *gomock.Controller) *HetznerUserConfigFilesResolver {
	mock := &HetznerUserConfigFilesResolver{ctrl: ctrl}
	mock.recorder = &HetznerUserConfigFilesResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *HetznerUserConfigFilesResolver) EXPECT() *HetznerUserConfigFilesResolverMockRecorder {
	return m.recorder
}

// Resolve mocks base method.
func (m *HetznerUserConfigFilesResolver) Resolve() (*userconfig.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve")
	ret0, _ := ret[0].(*userconfig.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve.
func (mr *HetznerUserConfigFilesResolverMockRecorder) Resolve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*HetznerUserConfigFilesResolver)(nil).Resolve))
}