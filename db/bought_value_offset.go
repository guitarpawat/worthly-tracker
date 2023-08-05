package db

import (
	"github.com/jmoiron/sqlx"
	"worthly-tracker/model"
)

var boughtValueOffsetRepo = &SqliteOffsetRepo{}

func GetBoughtValueOffsetRepo() BoughtValueOffsetRepo {
	return boughtValueOffsetRepo
}

type BoughtValueOffsetRepo interface {
	Get(date *model.Date, assetIds []int, tx *sqlx.Tx) ([]model.OffsetDetail, error)
	Upsert(data []model.OffsetDetail, tx *sqlx.Tx) error
}

type SqliteOffsetRepo struct {
}

func (s *SqliteOffsetRepo) Get(date *model.Date, assetIds []int, tx *sqlx.Tx) ([]model.OffsetDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SqliteOffsetRepo) Upsert(data []model.OffsetDetail, tx *sqlx.Tx) error {
	//TODO implement me
	panic("implement me")
}
