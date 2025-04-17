package db

import (
	"context"
	"database/sql"
	"github.com/go-sqlx/sqlx"
)

type Conn interface {
	DB
	Beginx() (*sqlx.Tx, error)
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DB interface {
	sqlx.Preparer
	sqlx.PreparerContext
	sqlx.Execer
	sqlx.ExecerContext
	sqlx.Queryer
	sqlx.QueryerContext
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
}

type TX interface {
	Commit() error
	Rollback() error
}

type tx interface {
	DB
	TX
}

type TxRepository[T any] interface {
	BeginTx() (T, TX, error)
	beginTx() (T, tx, error)
}
