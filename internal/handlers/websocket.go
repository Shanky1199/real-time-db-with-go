package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-database/internal/cache"
	"real-time-database/internal/database"
	"real-time-database/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Item)
var clientsMutex sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}

		// Process message and update cache/database as needed
		var item models.Item
		if err := json.Unmarshal(msg, &item); err == nil {
			if err := database.AddItem(item); err == nil {
				// Notify all connected clients
				broadcast <- item
			} else {
				log.Printf("Failed to add item: %v", err)
			}
		}
	}
}

func broadcastMessage() {
	for {
		item := <-broadcast
		clientsMutex.Lock()
		for client := range clients {
			err := client.WriteJSON(item)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}

func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleWebSocket)
	go broadcastMessage()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	cache.InitCache()
	database.InitDB()
	StartWebSocketServer()
}
