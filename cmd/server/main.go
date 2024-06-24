package main

import (
	"net/http"
	"real-time-database/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handlers.HandleWebSocket)
	r.HandleFunc("/api/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/api/items", handlers.GetItems).Methods("GET")

	http.ListenAndServe(":8080", r)
}
