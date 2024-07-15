package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"real-time-database/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestHandleWebSocket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandleWebSocket))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	item := models.Item{ID: 1, Name: "Test Item"}

	msg, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal item: %v", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	var receivedItem models.Item
	if err := conn.ReadJSON(&receivedItem); err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	if receivedItem.ID != item.ID || receivedItem.Name != item.Name {
		t.Fatalf("Received item does not match sent item: got %v, want %v", receivedItem, item)
	}
}

func TestBroadcastMessage(t *testing.T) {
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Item)

	server := httptest.NewServer(http.HandlerFunc(HandleWebSocket))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]

	conn1, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn2.Close()

	go broadcastMessage()

	item := models.Item{ID: 1, Name: "Broadcast Test Item"}
	broadcast <- item

	var receivedItem1 models.Item
	if err := conn1.ReadJSON(&receivedItem1); err != nil {
		t.Fatalf("Failed to read message from conn1: %v", err)
	}

	var receivedItem2 models.Item
	if err := conn2.ReadJSON(&receivedItem2); err != nil {
		t.Fatalf("Failed to read message from conn2: %v", err)
	}

	if receivedItem1.ID != item.ID || receivedItem1.Name != item.Name {
		t.Fatalf("conn1 received item does not match sent item: got %v, want %v", receivedItem1, item)
	}

	if receivedItem2.ID != item.ID || receivedItem2.Name != item.Name {
		t.Fatalf("conn2 received item does not match sent item: got %v, want %v", receivedItem2, item)
	}
}
