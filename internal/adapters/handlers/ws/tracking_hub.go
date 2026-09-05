package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type Client struct {
	Hub    *TrackingHub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string // The user this client is tracking
}

type TrackingHub struct {
	// Registered clients mapped by the UserID they are tracking
	clients    map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

type Message struct {
	UserID  string
	Payload []byte
}

func NewTrackingHub() *TrackingHub {
	return &TrackingHub{
		clients:    make(map[string]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *TrackingHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("WS: Client connected to track user: %s", client.UserID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID][client]; ok {
				delete(h.clients[client.UserID], client)
				close(client.Send)
				if len(h.clients[client.UserID]) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			log.Printf("WS: Client disconnected from user: %s", client.UserID)
		case msg := <-h.Broadcast:
			h.mu.Lock()
			// Send to specific user channel
			if clientsToNotify, ok := h.clients[msg.UserID]; ok {
				for client := range clientsToNotify {
					select {
					case client.Send <- msg.Payload:
					default:
						close(client.Send)
						delete(h.clients[msg.UserID], client)
					}
				}
			}
			// Send to "all" channel (Admin channel)
			if clientsToNotifyAll, ok := h.clients["all"]; ok {
				for client := range clientsToNotifyAll {
					select {
					case client.Send <- msg.Payload:
					default:
						close(client.Send)
						delete(h.clients["all"], client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Write pump for the client
func (c *Client) writePump() {
	defer c.Conn.Close()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *TrackingHub) ServeWS(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Hub:    h,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	client.Hub.Register <- client

	go client.writePump()
	
	go func() {
		defer func() {
			client.Hub.Unregister <- client
			client.Conn.Close()
		}()
		for {
			_, _, err := client.Conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
