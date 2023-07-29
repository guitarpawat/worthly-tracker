//go:build test || unit || integration

package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"worthly-tracker/model"
)

type MockConn struct {
	conn *sqlx.DB
	mock sqlmock.Sqlmock
}

func (m *MockConn) GetDB() *sqlx.DB {
	return m.conn
}

func (m *MockConn) BeginTx() (*sqlx.Tx, error) {
	return m.conn.Beginx()
}

func (m *MockConn) GetMock() sqlmock.Sqlmock {
	return m.mock
}

func InitMock() sqlmock.Sqlmock {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	db = &MockConn{sqlxDb, mock}

	return mock
}

type MockRecordRepo struct {
	mock.Mock
}

func (r *MockRecordRepo) GetDate(current model.Date, tx *sqlx.Tx) (*model.DateList, error) {
	args := r.Called(current, tx)
	return args.Get(0).(*model.DateList), args.Error(1)
}

func (r *MockRecordRepo) GetLatestDate(tx *sqlx.Tx) (*model.Date, error) {
	args := r.Called(tx)
	a0 := args.Get(0)
	if a0 == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Date), args.Error(1)
}

func (r *MockRecordRepo) GetRecordByDate(date model.Date, tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	args := r.Called(date, tx)
	return args.Get(0).([]model.AssetTypeRecord), args.Error(1)
}

func (r *MockRecordRepo) GetRecordDraft(tx *sqlx.Tx) ([]model.AssetTypeRecord, error) {
	args := r.Called(tx)
	return args.Get(0).([]model.AssetTypeRecord), args.Error(1)
}

func (r *MockRecordRepo) UpsertRecord(record model.AssetRecord, date model.Date, tx *sqlx.Tx) error {
	args := r.Called(record, date, tx)
	return args.Error(0)
}

func (r *MockRecordRepo) DeleteRecordById(id int, tx *sqlx.Tx) error {
	args := r.Called(id, tx)
	return args.Error(0)
}

func (r *MockRecordRepo) DeleteRecordByDate(date model.Date, tx *sqlx.Tx) (int64, error) {
	args := r.Called(date, tx)
	return int64(args.Int(0)), args.Error(1)
}

type MockAssetTypeRepo struct {
	mock.Mock
}

func (r *MockAssetTypeRepo) Get(isActive *bool, tx *sqlx.Tx) ([]model.AssetTypeDetail, error) {
	args := r.Called(isActive, tx)
	return args.Get(0).([]model.AssetTypeDetail), args.Error(1)
}

func (r *MockAssetTypeRepo) GetNames(isActive *bool, tx *sqlx.Tx) ([]model.NameDetail, error) {
	args := r.Called(isActive, tx)
	return args.Get(0).([]model.NameDetail), args.Error(1)
}

func (r *MockAssetTypeRepo) Upsert(assetType model.AssetTypeDetail, tx *sqlx.Tx) error {
	args := r.Called(assetType, tx)
	return args.Error(0)
}

func (r *MockAssetTypeRepo) Delete(id int, tx *sqlx.Tx) error {
	args := r.Called(id, tx)
	return args.Error(0)
}

func (r *MockAssetTypeRepo) UpdateSequence(sequences model.SequenceDetail, tx *sqlx.Tx) error {
	args := r.Called(sequences, tx)
	return args.Error(0)
}

type MockAssetRepo struct {
	mock.Mock
}

func (r *MockAssetRepo) Get(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.AssetDetail, error) {
	args := r.Called(isActive, typeId, tx)
	return args.Get(0).([]model.AssetDetail), args.Error(1)
}

func (r *MockAssetRepo) GetNames(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.NameDetail, error) {
	args := r.Called(isActive, typeId, tx)
	return args.Get(0).([]model.NameDetail), args.Error(1)
}

func (r *MockAssetRepo) Upsert(asset model.AssetDetail, tx *sqlx.Tx) error {
	args := r.Called(asset, tx)
	return args.Error(0)
}

func (r *MockAssetRepo) Delete(id int, tx *sqlx.Tx) error {
	args := r.Called(id, tx)
	return args.Error(0)
}

func (r *MockAssetRepo) UpdateSequence(sequence model.SequenceDetail, tx *sqlx.Tx) error {
	args := r.Called(sequence, tx)
	return args.Error(0)
}
