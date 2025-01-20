package db

import (
	"context"
	"github.com/glebarez/sqlite"
	"github.com/guitarpawat/worthly-tracker/internal/entity"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"gorm.io/gorm"
)

type SqliteConfig struct {
	Uri         string `validate:"required"`
	AutoMigrate bool
}

func NewSqlite(ctx context.Context, cfg *SqliteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Uri), &gorm.Config{
		Logger: logs.Log().ToGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	db = db.WithContext(ctx)

	if cfg.AutoMigrate {
		err = db.AutoMigrate(&entity.Asset{}, &entity.AssetType{}, &entity.Record{}, &entity.ValueOffset{})
		if err != nil {
			return nil, err
		}
	}

	return db, err
}
