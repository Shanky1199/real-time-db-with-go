package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetDB() *sql.DB {
	return db
}
