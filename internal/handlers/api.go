package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-database/internal/cache"
	"real-time-database/internal/database"
	"real-time-database/internal/models"
)

var itemCache = cache.NewCache()

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the item
	if err := item.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert item into the database
	if err := database.InsertItem(item.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add item to cache
	itemCache.Set(item.Data["id"].(string), item.Data)

	// Broadcast item to WebSocket clients
	broadcast <- item

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	cachedItems, found := itemCache.Get("allItems")
	if found {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedItems)
		return
	}

	items, err := database.GetItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemCache.Set("allItems", items)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
