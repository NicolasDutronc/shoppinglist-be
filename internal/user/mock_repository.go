// Code generated by mockery v2.5.1. DO NOT EDIT.

package user

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// AddPermissions provides a mock function with given fields: ctx, userID, permissions
func (_m *MockRepository) AddPermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	_va := make([]interface{}, len(permissions))
	for _i := range permissions {
		_va[_i] = permissions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, ...*Permission) int64); ok {
		r0 = rf(ctx, userID, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...*Permission) error); ok {
		r1 = rf(ctx, userID, permissions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, userID
func (_m *MockRepository) Delete(ctx context.Context, userID string) (int64, error) {
	ret := _m.Called(ctx, userID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: ctx, userID
func (_m *MockRepository) FindByID(ctx context.Context, userID string) (*User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *User
	if rf, ok := ret.Get(0).(func(context.Context, string) *User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: ctx, userName
func (_m *MockRepository) FindByName(ctx context.Context, userName string) (*User, error) {
	ret := _m.Called(ctx, userName)

	var r0 *User
	if rf, ok := ret.Get(0).(func(context.Context, string) *User); ok {
		r0 = rf(ctx, userName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemovePermissions provides a mock function with given fields: ctx, userID, permissions
func (_m *MockRepository) RemovePermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	_va := make([]interface{}, len(permissions))
	for _i := range permissions {
		_va[_i] = permissions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, ...*Permission) int64); ok {
		r0 = rf(ctx, userID, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...*Permission) error); ok {
		r1 = rf(ctx, userID, permissions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, name, password, permissions
func (_m *MockRepository) Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error) {
	_va := make([]interface{}, len(permissions))
	for _i := range permissions {
		_va[_i] = permissions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, password)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *User
	if rf, ok := ret.Get(0).(func(context.Context, string, string, ...*Permission) *User); ok {
		r0 = rf(ctx, name, password, permissions...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, ...*Permission) error); ok {
		r1 = rf(ctx, name, password, permissions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateName provides a mock function with given fields: ctx, userID, newName
func (_m *MockRepository) UpdateName(ctx context.Context, userID string, newName string) (int64, error) {
	ret := _m.Called(ctx, userID, newName)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int64); ok {
		r0 = rf(ctx, userID, newName)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, newName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: ctx, userID, newPassword
func (_m *MockRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) (int64, error) {
	ret := _m.Called(ctx, userID, newPassword)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int64); ok {
		r0 = rf(ctx, userID, newPassword)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, newPassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}