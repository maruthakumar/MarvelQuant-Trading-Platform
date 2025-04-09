package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	// Message types
	MessageTypeOrderUpdate     MessageType = "ORDER_UPDATE"
	MessageTypePositionUpdate  MessageType = "POSITION_UPDATE"
	MessageTypeStrategyUpdate  MessageType = "STRATEGY_UPDATE"
	MessageTypeMarketData      MessageType = "MARKET_DATA"
	MessageTypeAuthentication  MessageType = "AUTHENTICATION"
	MessageTypeSubscription    MessageType = "SUBSCRIPTION"
	MessageTypeError           MessageType = "ERROR"
	MessageTypeHeartbeat       MessageType = "HEARTBEAT"
	
	// WebSocket configuration
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      MessageType     `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Client represents a WebSocket client connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	topics   map[string]bool
	mu       sync.Mutex
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Inbound messages from clients
	broadcast chan []byte
	
	// Topic subscriptions
	topics map[string]map[*Client]bool
	
	// Mutex for thread safety
	mu sync.Mutex
}

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for now
		// In production, this should be restricted to known origins
		return true
	},
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		topics:     make(map[string]map[*Client]bool),
	}
}

// Run starts the Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				// Remove client from all topics
				h.mu.Lock()
				for topic, clients := range h.topics {
					if _, ok := clients[client]; ok {
						delete(h.topics[topic], client)
					}
				}
				h.mu.Unlock()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Subscribe adds a client to a topic
func (h *Hub) Subscribe(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, ok := h.topics[topic]; !ok {
		h.topics[topic] = make(map[*Client]bool)
	}
	h.topics[topic][client] = true
	
	client.mu.Lock()
	client.topics[topic] = true
	client.mu.Unlock()
}

// Unsubscribe removes a client from a topic
func (h *Hub) Unsubscribe(client *Client, topic string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if clients, ok := h.topics[topic]; ok {
		delete(clients, client)
	}
	
	client.mu.Lock()
	delete(client.topics, topic)
	client.mu.Unlock()
}

// BroadcastToTopic sends a message to all clients subscribed to a topic
func (h *Hub) BroadcastToTopic(topic string, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if clients, ok := h.topics[topic]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(clients, client)
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	for client := range h.clients {
		if client.userID == userID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// ServeWs handles WebSocket requests from clients
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		topics: make(map[string]bool),
	}
	client.hub.register <- client
	
	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})
	
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		// Handle client messages
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from clients
func (c *Client) handleMessage(message []byte) {
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		// Send error message back to client
		errorMsg := WebSocketMessage{
			Type:      MessageTypeError,
			Timestamp: time.Now(),
			Payload:   json.RawMessage(`{"error": "Invalid message format"}`),
		}
		if errorJSON, err := json.Marshal(errorMsg); err == nil {
			c.send <- errorJSON
		}
		return
	}
	
	// Handle different message types
	switch wsMessage.Type {
	case MessageTypeSubscription:
		var subscription struct {
			Action string   `json:"action"`
			Topics []string `json:"topics"`
		}
		if err := json.Unmarshal(wsMessage.Payload, &subscription); err != nil {
			return
		}
		
		if subscription.Action == "subscribe" {
			for _, topic := range subscription.Topics {
				c.hub.Subscribe(c, topic)
			}
		} else if subscription.Action == "unsubscribe" {
			for _, topic := range subscription.Topics {
				c.hub.Unsubscribe(c, topic)
			}
		}
		
	case MessageTypeAuthentication:
		// Authentication is already handled by the HTTP middleware
		// This is just for re-authentication if needed
		
	case MessageTypeHeartbeat:
		// Respond with a heartbeat
		heartbeat := WebSocketMessage{
			Type:      MessageTypeHeartbeat,
			Timestamp: time.Now(),
			Payload:   json.RawMessage(`{}`),
		}
		if heartbeatJSON, err := json.Marshal(heartbeat); err == nil {
			c.send <- heartbeatJSON
		}
		
	default:
		// Forward other messages to the hub for processing
		c.hub.broadcast <- message
	}
}
