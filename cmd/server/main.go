package main

import (
	"log"
	"net/http"
	"os"
	"real-time-database/internal/database"
	"real-time-database/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://username:password@localhost/real_time_db?sslmode=disable" // Update with your actual DB details
	}
	if err := database.InitDB(dbURL); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/ws", handlers.HandleWebSocket)
	r.HandleFunc("/api/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/api/items", handlers.GetItems).Methods("GET")

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
