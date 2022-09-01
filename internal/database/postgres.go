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
	err = db.Ping()

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
