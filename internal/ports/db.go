package ports

import (
	"context"
	"database/sql"
	"github.com/go-sqlx/sqlx"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/rickb777/date/v2"
)

type TX interface {
	Commit() error
	Rollback() error
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

type Conn interface {
	DB
	Beginx() (*sqlx.Tx, error)
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// AssetTypesRepository defines the interface for asset type repository operations
type AssetTypesRepository interface {
	BeginTx() (AssetTypesRepository, TX, error)
	Get(isActive *bool) ([]model.AssetTypeDetail, error)
	GetNames(isActive *bool) ([]model.AssetTypeNameDetail, error)
	Upsert(assetType model.AssetTypeDetail) error
	Delete(id int) error
	UpdateSequence(sequences model.AssetTypeSequenceDetail) error
}

// AssetsRepository defines the interface for asset repository operations
type AssetsRepository interface {
	BeginTx() (AssetsRepository, TX, error)
	Get(isActive *bool, typeId *int) ([]model.AssetDetail, error)
	GetNames(isActive *bool, typeId *int) ([]model.AssetNameDetail, error)
	Upsert(asset model.AssetDetail) error
	Delete(id int) error
	UpdateSequence(sequence model.AssetSequenceDetail) error
}

// RecordsRepository defines the interface for record repository operations
type RecordsRepository interface {
	BeginTx() (RecordsRepository, TX, error)
	GetDate(current date.Date) (*model.DateList, error)
	GetLatestDate() (date.Date, error)
	GetRecordByDate(d date.Date) ([]model.AssetTypeRecord, error)
	GetRecordDraft() ([]model.AssetTypeRecord, error)
	UpsertRecord(record model.AssetRecord, d date.Date) error
	DeleteRecordById(id int) error
	DeleteRecordByDate(d date.Date) (int64, error)
}
