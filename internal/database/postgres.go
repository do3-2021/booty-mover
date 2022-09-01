package database

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

var ErrPostgresNotFound = errors.New("POSTGRES environment variable not found")

var globalDB *sql.DB = nil

func ConnectPostgres() (db *sql.DB, err error) {
	connStr, found := os.LookupEnv("POSTGRES")
	if !found {
		return nil, ErrPostgresNotFound
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = initDB(db)

	if err == nil {
		globalDB = db
	}
	return
}

func GetDB() (db *sql.DB, err error) {
	if globalDB != nil {
		return globalDB, nil
	}
	return ConnectPostgres()

}

func initDB(db *sql.DB) (err error) {
	_, err = db.Query("CREATE TABLE IF NOT EXISTS guilds (id VARCHAR(24) PRIMARY KEY, group_channel VARCHAR(24))")
	return
}
