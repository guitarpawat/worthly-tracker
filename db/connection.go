package db

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"worthly-tracker/logs"
	"worthly-tracker/resource"
)

var db Connection

type Connection interface {
	GetDB() *sqlx.DB
	BeginTx() (*sqlx.Tx, error)
}

type SqlConn struct {
	conn *sqlx.DB
}

func (p *SqlConn) GetDB() *sqlx.DB {
	return p.conn
}

func (p *SqlConn) BeginTx() (*sqlx.Tx, error) {
	return p.conn.Beginx()
}

func Init() {
	connectDB()
	migrateDB()
}

func connectDB() {
	connDb, err := sqlx.Open("sqlite3", viper.GetString("datasource.uri")+"?_foreign_keys=true")
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
	migrationFs, err := iofs.New(resource.Loader(), "migration")
	if err != nil {
		logs.Log().Panicf("Unable to create iofs for db/migration: %v\n", err)
	}

	driver, err := sqlite.WithInstance(db.GetDB().DB, &sqlite.Config{})
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

func GetDB() Connection {
	return db
}
