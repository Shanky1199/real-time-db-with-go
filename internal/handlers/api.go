package handlers

import (
	"encoding/json"
	"net/http"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Insert item into the database
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	// Retrieve items from the database
	items := []map[string]interface{}{}
	json.NewEncoder(w).Encode(items)
}
