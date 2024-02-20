package db

import (
	"errors"
	"github.com/go-sqlx/sqlx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"worthly-tracker/logs"
	"worthly-tracker/ports"
	"worthly-tracker/resource"
)

var db ports.Connection

type SqlConn struct {
	conn *sqlx.DB
}

func (p *SqlConn) GetDB() *sqlx.DB {
	return p.conn
}

func (p *SqlConn) BeginTx() (*sqlx.Tx, error) {
	return p.conn.Beginx()
}

var cacheStr = ""

const datasourceUriKey = "datasource.uri"

func Init() {
	if viper.GetBool("datasource.cache") {
		cacheStr = "&mode=memory&cache=shared"
		conn, _ := sqlx.Open("sqlite3", viper.GetString(datasourceUriKey)+"?_foreign_keys=true"+cacheStr)
		_ = conn.Ping()
	}
	migrateDB()
	connectDB()
}

func connectDB() {
	connDb, err := sqlx.Open("sqlite3", viper.GetString(datasourceUriKey)+"?_foreign_keys=true"+cacheStr)
	if err != nil {
		logs.Log().Panicf("Unable to connect to database: %v\n", err)
	}
	if err = connDb.Ping(); err != nil {
		logs.Log().Panicf("Unable to connect to database: %v\n", err)
	}
	db = &SqlConn{connDb}
	logs.Log().Info("Database connected successfully")
}

func migrateDB() {
	connMigrate, err := sqlx.Open("sqlite3", viper.GetString(datasourceUriKey)+"?_foreign_keys=false"+cacheStr)
	if err != nil {
		logs.Log().Panicf("Unable to connect to database: %v\n", err)
	}
	defer connMigrate.Close()

	migrationFs, err := iofs.New(resource.Loader(), "migration")
	if err != nil {
		logs.Log().Panicf("Unable to create iofs for db/migration: %v\n", err)
	}

	driver, err := sqlite.WithInstance(connMigrate.DB, &sqlite.Config{})
	if err != nil {
		logs.Log().Panicf("Unable to create driver instance: %v\n", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrationFs, "worthly_tracker", driver)
	if err != nil {
		logs.Log().Panicf("Unable to create go-migrate instance: %v\n", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logs.Log().Debug("No changes in migration")
		} else {
			logs.Log().Panicf("Unable to migrate the database: %v\n", err)
		}
	}

	logs.Log().Info("Database migration run successfully")
}

func GetDB() ports.Connection {
	return db
}
