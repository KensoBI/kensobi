// Code generated by mockery v2.31.4. DO NOT EDIT.

package utils

import (
	mock "github.com/stretchr/testify/mock"
	cli "github.com/urfave/cli/v2"
)

// MockCommandLine is an autogenerated mock type for the CommandLine type
type MockCommandLine struct {
	mock.Mock
}

// Application provides a mock function with given fields:
func (_m *MockCommandLine) Application() *cli.App {
	ret := _m.Called()

	var r0 *cli.App
	if rf, ok := ret.Get(0).(func() *cli.App); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cli.App)
		}
	}

	return r0
}

// Args provides a mock function with given fields:
func (_m *MockCommandLine) Args() cli.Args {
	ret := _m.Called()

	var r0 cli.Args
	if rf, ok := ret.Get(0).(func() cli.Args); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cli.Args)
		}
	}

	return r0
}

// Bool provides a mock function with given fields: name
func (_m *MockCommandLine) Bool(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FlagNames provides a mock function with given fields:
func (_m *MockCommandLine) FlagNames() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Generic provides a mock function with given fields: name
func (_m *MockCommandLine) Generic(name string) interface{} {
	ret := _m.Called(name)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Int provides a mock function with given fields: name
func (_m *MockCommandLine) Int(name string) int {
	ret := _m.Called(name)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// PluginDirectory provides a mock function with given fields:
func (_m *MockCommandLine) PluginDirectory() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// PluginRepoURL provides a mock function with given fields:
func (_m *MockCommandLine) PluginRepoURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// PluginURL provides a mock function with given fields:
func (_m *MockCommandLine) PluginURL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ShowHelp provides a mock function with given fields:
func (_m *MockCommandLine) ShowHelp() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShowVersion provides a mock function with given fields:
func (_m *MockCommandLine) ShowVersion() {
	_m.Called()
}

// String provides a mock function with given fields: name
func (_m *MockCommandLine) String(name string) string {
	ret := _m.Called(name)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// StringSlice provides a mock function with given fields: name
func (_m *MockCommandLine) StringSlice(name string) []string {
	ret := _m.Called(name)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// NewMockCommandLine creates a new instance of MockCommandLine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommandLine(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommandLine {
	mock := &MockCommandLine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
