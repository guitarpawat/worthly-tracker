// Code generated by mockery v2.42.2. DO NOT EDIT.

//go:build test || unit || integration

package mocks

import (
	model "worthly-tracker/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/go-sqlx/sqlx"
)

// MockAssetTypeRepo is an autogenerated mock type for the AssetTypeRepo type
type MockAssetTypeRepo struct {
	mock.Mock
}

type MockAssetTypeRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAssetTypeRepo) EXPECT() *MockAssetTypeRepo_Expecter {
	return &MockAssetTypeRepo_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: id, tx
func (_m *MockAssetTypeRepo) Delete(id int, tx *sqlx.Tx) error {
	ret := _m.Called(id, tx)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *sqlx.Tx) error); ok {
		r0 = rf(id, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAssetTypeRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockAssetTypeRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id int
//   - tx *sqlx.Tx
func (_e *MockAssetTypeRepo_Expecter) Delete(id interface{}, tx interface{}) *MockAssetTypeRepo_Delete_Call {
	return &MockAssetTypeRepo_Delete_Call{Call: _e.mock.On("Delete", id, tx)}
}

func (_c *MockAssetTypeRepo_Delete_Call) Run(run func(id int, tx *sqlx.Tx)) *MockAssetTypeRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockAssetTypeRepo_Delete_Call) Return(_a0 error) *MockAssetTypeRepo_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAssetTypeRepo_Delete_Call) RunAndReturn(run func(int, *sqlx.Tx) error) *MockAssetTypeRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: isActive, tx
func (_m *MockAssetTypeRepo) Get(isActive *bool, tx *sqlx.Tx) ([]model.AssetTypeDetail, error) {
	ret := _m.Called(isActive, tx)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []model.AssetTypeDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(*bool, *sqlx.Tx) ([]model.AssetTypeDetail, error)); ok {
		return rf(isActive, tx)
	}
	if rf, ok := ret.Get(0).(func(*bool, *sqlx.Tx) []model.AssetTypeDetail); ok {
		r0 = rf(isActive, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.AssetTypeDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(*bool, *sqlx.Tx) error); ok {
		r1 = rf(isActive, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAssetTypeRepo_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockAssetTypeRepo_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - isActive *bool
//   - tx *sqlx.Tx
func (_e *MockAssetTypeRepo_Expecter) Get(isActive interface{}, tx interface{}) *MockAssetTypeRepo_Get_Call {
	return &MockAssetTypeRepo_Get_Call{Call: _e.mock.On("Get", isActive, tx)}
}

func (_c *MockAssetTypeRepo_Get_Call) Run(run func(isActive *bool, tx *sqlx.Tx)) *MockAssetTypeRepo_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bool), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockAssetTypeRepo_Get_Call) Return(_a0 []model.AssetTypeDetail, _a1 error) *MockAssetTypeRepo_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAssetTypeRepo_Get_Call) RunAndReturn(run func(*bool, *sqlx.Tx) ([]model.AssetTypeDetail, error)) *MockAssetTypeRepo_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetNames provides a mock function with given fields: isActive, tx
func (_m *MockAssetTypeRepo) GetNames(isActive *bool, tx *sqlx.Tx) ([]model.NameDetail, error) {
	ret := _m.Called(isActive, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetNames")
	}

	var r0 []model.NameDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(*bool, *sqlx.Tx) ([]model.NameDetail, error)); ok {
		return rf(isActive, tx)
	}
	if rf, ok := ret.Get(0).(func(*bool, *sqlx.Tx) []model.NameDetail); ok {
		r0 = rf(isActive, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.NameDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(*bool, *sqlx.Tx) error); ok {
		r1 = rf(isActive, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAssetTypeRepo_GetNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNames'
type MockAssetTypeRepo_GetNames_Call struct {
	*mock.Call
}

// GetNames is a helper method to define mock.On call
//   - isActive *bool
//   - tx *sqlx.Tx
func (_e *MockAssetTypeRepo_Expecter) GetNames(isActive interface{}, tx interface{}) *MockAssetTypeRepo_GetNames_Call {
	return &MockAssetTypeRepo_GetNames_Call{Call: _e.mock.On("GetNames", isActive, tx)}
}

func (_c *MockAssetTypeRepo_GetNames_Call) Run(run func(isActive *bool, tx *sqlx.Tx)) *MockAssetTypeRepo_GetNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bool), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockAssetTypeRepo_GetNames_Call) Return(_a0 []model.NameDetail, _a1 error) *MockAssetTypeRepo_GetNames_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAssetTypeRepo_GetNames_Call) RunAndReturn(run func(*bool, *sqlx.Tx) ([]model.NameDetail, error)) *MockAssetTypeRepo_GetNames_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSequence provides a mock function with given fields: sequences, tx
func (_m *MockAssetTypeRepo) UpdateSequence(sequences model.SequenceDetail, tx *sqlx.Tx) error {
	ret := _m.Called(sequences, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSequence")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.SequenceDetail, *sqlx.Tx) error); ok {
		r0 = rf(sequences, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAssetTypeRepo_UpdateSequence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSequence'
type MockAssetTypeRepo_UpdateSequence_Call struct {
	*mock.Call
}

// UpdateSequence is a helper method to define mock.On call
//   - sequences model.SequenceDetail
//   - tx *sqlx.Tx
func (_e *MockAssetTypeRepo_Expecter) UpdateSequence(sequences interface{}, tx interface{}) *MockAssetTypeRepo_UpdateSequence_Call {
	return &MockAssetTypeRepo_UpdateSequence_Call{Call: _e.mock.On("UpdateSequence", sequences, tx)}
}

func (_c *MockAssetTypeRepo_UpdateSequence_Call) Run(run func(sequences model.SequenceDetail, tx *sqlx.Tx)) *MockAssetTypeRepo_UpdateSequence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.SequenceDetail), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockAssetTypeRepo_UpdateSequence_Call) Return(_a0 error) *MockAssetTypeRepo_UpdateSequence_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAssetTypeRepo_UpdateSequence_Call) RunAndReturn(run func(model.SequenceDetail, *sqlx.Tx) error) *MockAssetTypeRepo_UpdateSequence_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: assetType, tx
func (_m *MockAssetTypeRepo) Upsert(assetType model.AssetTypeDetail, tx *sqlx.Tx) error {
	ret := _m.Called(assetType, tx)

	if len(ret) == 0 {
		panic("no return value specified for Upsert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.AssetTypeDetail, *sqlx.Tx) error); ok {
		r0 = rf(assetType, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAssetTypeRepo_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type MockAssetTypeRepo_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - assetType model.AssetTypeDetail
//   - tx *sqlx.Tx
func (_e *MockAssetTypeRepo_Expecter) Upsert(assetType interface{}, tx interface{}) *MockAssetTypeRepo_Upsert_Call {
	return &MockAssetTypeRepo_Upsert_Call{Call: _e.mock.On("Upsert", assetType, tx)}
}

func (_c *MockAssetTypeRepo_Upsert_Call) Run(run func(assetType model.AssetTypeDetail, tx *sqlx.Tx)) *MockAssetTypeRepo_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.AssetTypeDetail), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockAssetTypeRepo_Upsert_Call) Return(_a0 error) *MockAssetTypeRepo_Upsert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAssetTypeRepo_Upsert_Call) RunAndReturn(run func(model.AssetTypeDetail, *sqlx.Tx) error) *MockAssetTypeRepo_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAssetTypeRepo creates a new instance of MockAssetTypeRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAssetTypeRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAssetTypeRepo {
	mock := &MockAssetTypeRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
