//go:build test || unit || integration

package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sqlx/sqlx"
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
