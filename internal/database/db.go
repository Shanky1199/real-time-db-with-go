package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	// Set connection pooling configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return err
	}
	log.Println("Database connected successfully")
	return nil
}

func GetDB() *sql.DB {
	return db
}

func InsertItem(item map[string]interface{}) error {
	query := "INSERT INTO items (data) VALUES ($1)"
	_, err := db.Exec(query, item)
	return err
}

func GetItems() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT data FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var item map[string]interface{}
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
