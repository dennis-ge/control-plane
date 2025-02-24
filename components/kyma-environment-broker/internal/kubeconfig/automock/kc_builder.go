// Code generated by mockery v2.14.0. DO NOT EDIT.

package automock

import (
	internal "github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"

	mock "github.com/stretchr/testify/mock"
)

// KcBuilder is an autogenerated mock type for the KcBuilder type
type KcBuilder struct {
	mock.Mock
}

// Build provides a mock function with given fields: _a0
func (_m *KcBuilder) Build(_a0 *internal.Instance) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(*internal.Instance) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*internal.Instance) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildFromAdminKubeconfig provides a mock function with given fields: instance, adminKubeconfig
func (_m *KcBuilder) BuildFromAdminKubeconfig(instance *internal.Instance, adminKubeconfig string) (string, error) {
	ret := _m.Called(instance, adminKubeconfig)

	var r0 string
	if rf, ok := ret.Get(0).(func(*internal.Instance, string) string); ok {
		r0 = rf(instance, adminKubeconfig)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*internal.Instance, string) error); ok {
		r1 = rf(instance, adminKubeconfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewKcBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewKcBuilder creates a new instance of KcBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKcBuilder(t mockConstructorTestingTNewKcBuilder) *KcBuilder {
	mock := &KcBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
