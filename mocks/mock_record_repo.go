// Code generated by mockery v2.42.2. DO NOT EDIT.

//go:build test || unit || integration

package mocks

import (
	model "worthly-tracker/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/go-sqlx/sqlx"
)

// MockRecordRepo is an autogenerated mock type for the RecordRepo type
type MockRecordRepo struct {
	mock.Mock
}

type MockRecordRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRecordRepo) EXPECT() *MockRecordRepo_Expecter {
	return &MockRecordRepo_Expecter{mock: &_m.Mock}
}

// DeleteRecordByDate provides a mock function with given fields: date, tx
func (_m *MockRecordRepo) DeleteRecordByDate(date model.Date, tx *sqlx.Tx) (int64, error) {
	ret := _m.Called(date, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRecordByDate")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) (int64, error)); ok {
		return rf(date, tx)
	}
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) int64); ok {
		r0 = rf(date, tx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(model.Date, *sqlx.Tx) error); ok {
		r1 = rf(date, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRecordRepo_DeleteRecordByDate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRecordByDate'
type MockRecordRepo_DeleteRecordByDate_Call struct {
	*mock.Call
}

// DeleteRecordByDate is a helper method to define mock.On call
//   - date model.Date
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) DeleteRecordByDate(date interface{}, tx interface{}) *MockRecordRepo_DeleteRecordByDate_Call {
	return &MockRecordRepo_DeleteRecordByDate_Call{Call: _e.mock.On("DeleteRecordByDate", date, tx)}
}

func (_c *MockRecordRepo_DeleteRecordByDate_Call) Run(run func(date model.Date, tx *sqlx.Tx)) *MockRecordRepo_DeleteRecordByDate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.Date), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_DeleteRecordByDate_Call) Return(_a0 int64, _a1 error) *MockRecordRepo_DeleteRecordByDate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRecordRepo_DeleteRecordByDate_Call) RunAndReturn(run func(model.Date, *sqlx.Tx) (int64, error)) *MockRecordRepo_DeleteRecordByDate_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteRecordById provides a mock function with given fields: id, tx
func (_m *MockRecordRepo) DeleteRecordById(id int, tx *sqlx.Tx) error {
	ret := _m.Called(id, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRecordById")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *sqlx.Tx) error); ok {
		r0 = rf(id, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRecordRepo_DeleteRecordById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRecordById'
type MockRecordRepo_DeleteRecordById_Call struct {
	*mock.Call
}

// DeleteRecordById is a helper method to define mock.On call
//   - id int
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) DeleteRecordById(id interface{}, tx interface{}) *MockRecordRepo_DeleteRecordById_Call {
	return &MockRecordRepo_DeleteRecordById_Call{Call: _e.mock.On("DeleteRecordById", id, tx)}
}

func (_c *MockRecordRepo_DeleteRecordById_Call) Run(run func(id int, tx *sqlx.Tx)) *MockRecordRepo_DeleteRecordById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_DeleteRecordById_Call) Return(_a0 error) *MockRecordRepo_DeleteRecordById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRecordRepo_DeleteRecordById_Call) RunAndReturn(run func(int, *sqlx.Tx) error) *MockRecordRepo_DeleteRecordById_Call {
	_c.Call.Return(run)
	return _c
}

// GetDate provides a mock function with given fields: current, tx
func (_m *MockRecordRepo) GetDate(current model.Date, tx *sqlx.Tx) (*model.DateList, error) {
	ret := _m.Called(current, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetDate")
	}

	var r0 *model.DateList
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) (*model.DateList, error)); ok {
		return rf(current, tx)
	}
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) *model.DateList); ok {
		r0 = rf(current, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DateList)
		}
	}

	if rf, ok := ret.Get(1).(func(model.Date, *sqlx.Tx) error); ok {
		r1 = rf(current, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRecordRepo_GetDate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDate'
type MockRecordRepo_GetDate_Call struct {
	*mock.Call
}

// GetDate is a helper method to define mock.On call
//   - current model.Date
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) GetDate(current interface{}, tx interface{}) *MockRecordRepo_GetDate_Call {
	return &MockRecordRepo_GetDate_Call{Call: _e.mock.On("GetDate", current, tx)}
}

func (_c *MockRecordRepo_GetDate_Call) Run(run func(current model.Date, tx *sqlx.Tx)) *MockRecordRepo_GetDate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.Date), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_GetDate_Call) Return(_a0 *model.DateList, _a1 error) *MockRecordRepo_GetDate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRecordRepo_GetDate_Call) RunAndReturn(run func(model.Date, *sqlx.Tx) (*model.DateList, error)) *MockRecordRepo_GetDate_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestDate provides a mock function with given fields: tx
func (_m *MockRecordRepo) GetLatestDate(tx *sqlx.Tx) (*model.Date, error) {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestDate")
	}

	var r0 *model.Date
	var r1 error
	if rf, ok := ret.Get(0).(func(*sqlx.Tx) (*model.Date, error)); ok {
		return rf(tx)
	}
	if rf, ok := ret.Get(0).(func(*sqlx.Tx) *model.Date); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Date)
		}
	}

	if rf, ok := ret.Get(1).(func(*sqlx.Tx) error); ok {
		r1 = rf(tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRecordRepo_GetLatestDate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestDate'
type MockRecordRepo_GetLatestDate_Call struct {
	*mock.Call
}

// GetLatestDate is a helper method to define mock.On call
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) GetLatestDate(tx interface{}) *MockRecordRepo_GetLatestDate_Call {
	return &MockRecordRepo_GetLatestDate_Call{Call: _e.mock.On("GetLatestDate", tx)}
}

func (_c *MockRecordRepo_GetLatestDate_Call) Run(run func(tx *sqlx.Tx)) *MockRecordRepo_GetLatestDate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_GetLatestDate_Call) Return(_a0 *model.Date, _a1 error) *MockRecordRepo_GetLatestDate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRecordRepo_GetLatestDate_Call) RunAndReturn(run func(*sqlx.Tx) (*model.Date, error)) *MockRecordRepo_GetLatestDate_Call {
	_c.Call.Return(run)
	return _c
}

// GetRecordByDate provides a mock function with given fields: date, tx
func (_m *MockRecordRepo) GetRecordByDate(date model.Date, tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	ret := _m.Called(date, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetRecordByDate")
	}

	var r0 []model.AssetTypeRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) ([]model.AssetTypeRecord, error)); ok {
		return rf(date, tx)
	}
	if rf, ok := ret.Get(0).(func(model.Date, *sqlx.Tx) []model.AssetTypeRecord); ok {
		r0 = rf(date, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.AssetTypeRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(model.Date, *sqlx.Tx) error); ok {
		r1 = rf(date, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRecordRepo_GetRecordByDate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRecordByDate'
type MockRecordRepo_GetRecordByDate_Call struct {
	*mock.Call
}

// GetRecordByDate is a helper method to define mock.On call
//   - date model.Date
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) GetRecordByDate(date interface{}, tx interface{}) *MockRecordRepo_GetRecordByDate_Call {
	return &MockRecordRepo_GetRecordByDate_Call{Call: _e.mock.On("GetRecordByDate", date, tx)}
}

func (_c *MockRecordRepo_GetRecordByDate_Call) Run(run func(date model.Date, tx *sqlx.Tx)) *MockRecordRepo_GetRecordByDate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.Date), args[1].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_GetRecordByDate_Call) Return(_a0 []model.AssetTypeRecord, _a1 error) *MockRecordRepo_GetRecordByDate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRecordRepo_GetRecordByDate_Call) RunAndReturn(run func(model.Date, *sqlx.Tx) ([]model.AssetTypeRecord, error)) *MockRecordRepo_GetRecordByDate_Call {
	_c.Call.Return(run)
	return _c
}

// GetRecordDraft provides a mock function with given fields: tx
func (_m *MockRecordRepo) GetRecordDraft(tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for GetRecordDraft")
	}

	var r0 []model.AssetTypeRecord
	var r1 error
	if rf, ok := ret.Get(0).(func(*sqlx.Tx) ([]model.AssetTypeRecord, error)); ok {
		return rf(tx)
	}
	if rf, ok := ret.Get(0).(func(*sqlx.Tx) []model.AssetTypeRecord); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.AssetTypeRecord)
		}
	}

	if rf, ok := ret.Get(1).(func(*sqlx.Tx) error); ok {
		r1 = rf(tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRecordRepo_GetRecordDraft_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRecordDraft'
type MockRecordRepo_GetRecordDraft_Call struct {
	*mock.Call
}

// GetRecordDraft is a helper method to define mock.On call
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) GetRecordDraft(tx interface{}) *MockRecordRepo_GetRecordDraft_Call {
	return &MockRecordRepo_GetRecordDraft_Call{Call: _e.mock.On("GetRecordDraft", tx)}
}

func (_c *MockRecordRepo_GetRecordDraft_Call) Run(run func(tx *sqlx.Tx)) *MockRecordRepo_GetRecordDraft_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_GetRecordDraft_Call) Return(_a0 []model.AssetTypeRecord, _a1 error) *MockRecordRepo_GetRecordDraft_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRecordRepo_GetRecordDraft_Call) RunAndReturn(run func(*sqlx.Tx) ([]model.AssetTypeRecord, error)) *MockRecordRepo_GetRecordDraft_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertRecord provides a mock function with given fields: record, date, tx
func (_m *MockRecordRepo) UpsertRecord(record model.AssetRecord, date model.Date, tx *sqlx.Tx) error {
	ret := _m.Called(record, date, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpsertRecord")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.AssetRecord, model.Date, *sqlx.Tx) error); ok {
		r0 = rf(record, date, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRecordRepo_UpsertRecord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertRecord'
type MockRecordRepo_UpsertRecord_Call struct {
	*mock.Call
}

// UpsertRecord is a helper method to define mock.On call
//   - record model.AssetRecord
//   - date model.Date
//   - tx *sqlx.Tx
func (_e *MockRecordRepo_Expecter) UpsertRecord(record interface{}, date interface{}, tx interface{}) *MockRecordRepo_UpsertRecord_Call {
	return &MockRecordRepo_UpsertRecord_Call{Call: _e.mock.On("UpsertRecord", record, date, tx)}
}

func (_c *MockRecordRepo_UpsertRecord_Call) Run(run func(record model.AssetRecord, date model.Date, tx *sqlx.Tx)) *MockRecordRepo_UpsertRecord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.AssetRecord), args[1].(model.Date), args[2].(*sqlx.Tx))
	})
	return _c
}

func (_c *MockRecordRepo_UpsertRecord_Call) Return(_a0 error) *MockRecordRepo_UpsertRecord_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRecordRepo_UpsertRecord_Call) RunAndReturn(run func(model.AssetRecord, model.Date, *sqlx.Tx) error) *MockRecordRepo_UpsertRecord_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRecordRepo creates a new instance of MockRecordRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRecordRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRecordRepo {
	mock := &MockRecordRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
