package db

import (
	"context"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/guitarpawat/worthly-tracker/internal/model"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	"gorm.io/gorm"
)

type SqliteConfig struct {
	Uri         string `validate:"required"`
	AutoMigrate bool
}

func NewSqlite(ctx context.Context, cfg *SqliteConfig) (*gorm.DB, error) {
	uri := fmt.Sprintf("%s?_pragma=foreign_keys(1)", cfg.Uri)
	logs.Log().Debugf("sqlite uri: %s", uri)
	db, err := gorm.Open(sqlite.Open(uri), &gorm.Config{
		Logger: logs.Log().ToGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	db = db.WithContext(ctx)

	if cfg.AutoMigrate {
		logs.Log().Info("auto migrate database enabled")
		err = db.AutoMigrate(&model.Asset{}, &model.AssetType{}, &model.Record{})
		if err != nil {
			return nil, err
		}
	}

	return db, err
}
