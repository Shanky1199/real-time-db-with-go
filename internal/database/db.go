package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"real-time-database/internal/cache"
	"real-time-database/internal/models"
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

func AddItem(item models.Item) error {
	_, err := db.ExecContext(context.Background(), "INSERT INTO items (name, data) VALUES ($1, $2)", item.Name, item.Data)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := cache.SetItem("items", "", 0); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func UpdateItem(item models.Item) error {
	_, err := db.ExecContext(context.Background(), "UPDATE items SET name = $1, data = $2 WHERE id = $3", item.Name, item.Data, item.ID)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := cache.SetItem("items", "", 0); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	return nil
}

func GetItems() ([]models.Item, error) {
	const cacheKey = "items"

	cachedItems, err := cache.GetItem(cacheKey)
	if err == nil {
		var items []models.Item
		if err := json.Unmarshal([]byte(cachedItems), &items); err == nil {
			return items, nil
		}
	}

	// If cache miss or unmarshalling failed, query the database
	rows, err := db.QueryContext(context.Background(), "SELECT id, name, data FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Data); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if data, err := json.Marshal(items); err == nil {
		cache.SetItem(cacheKey, data, 5*time.Minute)
	}

	return items, nil
}
