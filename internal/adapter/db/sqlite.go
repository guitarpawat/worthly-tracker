package db

import (
	"errors"
	"fmt"
	"github.com/go-sqlx/sqlx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/guitarpawat/worthly-tracker/internal/ports"
	"github.com/guitarpawat/worthly-tracker/resource"
	"github.com/guitarpawat/worthly-tracker/utility/logs"
	_ "modernc.org/sqlite"
)

type SqliteConfig struct {
	Uri         string `validate:"required"`
	AutoMigrate bool
	Cache       bool
}

func NewSqlite(cfg *SqliteConfig, log *logs.Logger) (ports.Conn, error) {
	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}
	if cfg.AutoMigrate {
		err = migrateDB(cfg, log)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func NewInMemorySqlite(log *logs.Logger) (ports.Conn, error) {
	return NewSqlite(&SqliteConfig{
		Uri:         "file::memory:",
		Cache:       true,
		AutoMigrate: true,
	}, log)
}

func connectDB(cfg *SqliteConfig) (ports.Conn, error) {
	var uri string
	if cfg.Cache {
		uri = cfg.Uri + "?_foreign_keys=true&mode=memory&cache=shared"
	} else {
		uri = cfg.Uri + "?_foreign_keys=true"
	}
	conn, err := sqlx.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}

func migrateDB(cfg *SqliteConfig, log *logs.Logger) error {
	var uri string
	if cfg.Cache {
		uri = cfg.Uri + "?_foreign_keys=false&mode=memory&cache=shared"
	} else {
		uri = cfg.Uri + "?_foreign_keys=false"
	}
	conn, err := sqlx.Open("sqlite", uri)
	if err != nil {
		return err
	}
	defer conn.Close()

	migrationFs, err := iofs.New(resource.DbMigration, "migration")
	if err != nil {
		return fmt.Errorf("unable to create iofs for db/migration: %w", err)
	}

	driver, err := sqlite.WithInstance(conn.DB, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("unable to create driver instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrationFs, "worthly_tracker", driver)
	if err != nil {
		return fmt.Errorf("unable to create go-migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Debug("No changes in migration")
		} else {
			return fmt.Errorf("unable to migrate the database: %w", err)
		}
	}

	log.Info("Database migration run successfully")
	return nil
}
