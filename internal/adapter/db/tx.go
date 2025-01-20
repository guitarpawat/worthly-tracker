package db

import (
	"context"
	"gorm.io/gorm"
)

type TxRepository[T any] interface {
	BeginTx(ctx context.Context) (T, *Tx)
}

type Tx struct {
	tx *gorm.DB
}

func (t *Tx) Commit() error {
	return t.tx.Commit().Error
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback().Error
}
