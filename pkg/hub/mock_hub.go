// Code generated by mockery v1.0.0. DO NOT EDIT.

package hub

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockHub is an autogenerated mock type for the Hub type
type MockHub struct {
	mock.Mock
}

// AddTopic provides a mock function with given fields: ctx, topic
func (_m *MockHub) AddTopic(ctx context.Context, topic Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields: ctx
func (_m *MockHub) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTopic provides a mock function with given fields: ctx, topic
func (_m *MockHub) DeleteTopic(ctx context.Context, topic Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProcessor provides a mock function with given fields: ctx, processorID
func (_m *MockHub) GetProcessor(ctx context.Context, processorID string) (Processor, error) {
	ret := _m.Called(ctx, processorID)

	var r0 Processor
	if rf, ok := ret.Get(0).(func(context.Context, string) Processor); ok {
		r0 = rf(ctx, processorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Processor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, processorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publish provides a mock function with given fields: ctx, msg
func (_m *MockHub) Publish(ctx context.Context, msg Message) error {
	ret := _m.Called(ctx, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Message) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterProcessor provides a mock function with given fields: ctx, p
func (_m *MockHub) RegisterProcessor(ctx context.Context, p Processor) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Processor) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: ctx, interrupt
func (_m *MockHub) Run(ctx context.Context, interrupt chan struct{}) error {
	ret := _m.Called(ctx, interrupt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, chan struct{}) error); ok {
		r0 = rf(ctx, interrupt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: ctx, p, topic
func (_m *MockHub) Subscribe(ctx context.Context, p Processor, topic Topic) error {
	ret := _m.Called(ctx, p, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Processor, Topic) error); ok {
		r0 = rf(ctx, p, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnregisterProcessor provides a mock function with given fields: ctx, p
func (_m *MockHub) UnregisterProcessor(ctx context.Context, p Processor) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Processor) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with given fields: ctx, p, topic
func (_m *MockHub) Unsubscribe(ctx context.Context, p Processor, topic Topic) error {
	ret := _m.Called(ctx, p, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Processor, Topic) error); ok {
		r0 = rf(ctx, p, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}