package ports

import (
	"github.com/jmoiron/sqlx"
	"worthly-tracker/model"
)

type Connection interface {
	GetDB() *sqlx.DB
	BeginTx() (*sqlx.Tx, error)
}

type AssetRepo interface {
	Get(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.AssetDetail, error)
	GetNames(isActive *bool, typeId *int, tx *sqlx.Tx) ([]model.NameDetail, error)
	Upsert(asset model.AssetDetail, tx *sqlx.Tx) error
	Delete(id int, tx *sqlx.Tx) error
	UpdateSequence(sequence model.SequenceDetail, tx *sqlx.Tx) error
}

type AssetTypeRepo interface {
	Get(isActive *bool, tx *sqlx.Tx) ([]model.AssetTypeDetail, error)
	GetNames(isActive *bool, tx *sqlx.Tx) ([]model.NameDetail, error)
	Upsert(assetType model.AssetTypeDetail, tx *sqlx.Tx) error
	Delete(id int, tx *sqlx.Tx) error
	UpdateSequence(sequences model.SequenceDetail, tx *sqlx.Tx) error
}

type BoughtValueOffsetRepo interface {
	Get(date model.Date, assetId int, tx *sqlx.Tx) (model.OffsetDetail, error)
	GetAllByAssetId(assetId int, tx *sqlx.Tx) ([]model.OffsetDetail, error)
	GetAllByDate(date model.Date, tx *sqlx.Tx) ([]model.OffsetDetail, error)
	Upsert(data model.OffsetDetail, tx *sqlx.Tx) error
	Delete(id int, tx *sqlx.Tx) error
}

type RecordRepo interface {
	GetDate(current model.Date, tx *sqlx.Tx) (*model.DateList, error)
	GetLatestDate(tx *sqlx.Tx) (*model.Date, error)
	GetRecordByDate(date model.Date, tx *sqlx.Tx) ([]model.AssetTypeRecord, error)
	GetRecordDraft(tx *sqlx.Tx) ([]model.AssetTypeRecord, error)
	UpsertRecord(record model.AssetRecord, date model.Date, tx *sqlx.Tx) error
	DeleteRecordById(id int, tx *sqlx.Tx) error
	DeleteRecordByDate(date model.Date, tx *sqlx.Tx) (int64, error)
}
