package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"trading-platform/backend/internal/portfolioanalytics"
	"trading-platform/backend/internal/orderexecution"
)

// Handler handles WebSocket connections
type Handler struct {
	portfolioService portfolioanalytics.Service
	orderService     orderexecution.Service
	clients          map[*Client]bool
	register         chan *Client
	unregister       chan *Client
	broadcast        chan []byte
	mutex            sync.Mutex
}

// Client represents a WebSocket client
type Client struct {
	conn       *websocket.Conn
	handler    *Handler
	send       chan []byte
	userID     string
	subscriptions map[string]string // Map of subscription type to ID
}

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Subscription represents a subscription request
type Subscription struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Upgrader upgrades HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, this should check the origin
		return true
	},
}

// NewHandler creates a new WebSocket handler
func NewHandler(portfolioService portfolioanalytics.Service, orderService orderexecution.Service) *Handler {
	return &Handler{
		portfolioService: portfolioService,
		orderService:     orderService,
		clients:          make(map[*Client]bool),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan []byte),
	}
}

// HandleConnection handles a WebSocket connection
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Get user ID from request (in a real implementation, this would come from the JWT token)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	// Create new client
	client := &Client{
		conn:         conn,
		handler:      h,
		send:         make(chan []byte, 256),
		userID:       userID,
		subscriptions: make(map[string]string),
	}

	// Register client
	h.register <- client

	// Start client goroutines
	go client.readPump()
	go client.writePump()
}

// Run starts the WebSocket handler
func (h *Handler) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Unsubscribe from all subscriptions
				for subType, subID := range client.subscriptions {
					switch subType {
					case "portfolio":
						h.portfolioService.UnsubscribeFromUpdates(subID)
					}
				}
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// readPump pumps messages from the WebSocket connection to the handler
func (c *Client) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// Handle message based on type
		switch msg.Type {
		case "subscribe":
			var sub Subscription
			if err := json.Unmarshal(msg.Payload, &sub); err != nil {
				log.Printf("Failed to parse subscription: %v", err)
				continue
			}
			c.handleSubscription(sub)
		case "unsubscribe":
			var sub Subscription
			if err := json.Unmarshal(msg.Payload, &sub); err != nil {
				log.Printf("Failed to parse subscription: %v", err)
				continue
			}
			c.handleUnsubscription(sub)
		}
	}
}

// writePump pumps messages from the handler to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The handler closed the channel
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
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleSubscription handles a subscription request
func (c *Client) handleSubscription(sub Subscription) {
	switch sub.Type {
	case "portfolio":
		// Subscribe to portfolio updates
		subID, err := c.handler.portfolioService.SubscribeToUpdates(sub.ID, func(data interface{}) {
			// Send update to client
			message, err := json.Marshal(Message{
				Type:    "portfolio_update",
				Payload: data.(json.RawMessage),
			})
			if err != nil {
				log.Printf("Failed to marshal portfolio update: %v", err)
				return
			}
			c.send <- message
		})
		if err != nil {
			log.Printf("Failed to subscribe to portfolio updates: %v", err)
			return
		}
		c.subscriptions[sub.Type] = subID
	}
}

// handleUnsubscription handles an unsubscription request
func (c *Client) handleUnsubscription(sub Subscription) {
	subID, ok := c.subscriptions[sub.Type]
	if !ok {
		return
	}

	switch sub.Type {
	case "portfolio":
		// Unsubscribe from portfolio updates
		err := c.handler.portfolioService.UnsubscribeFromUpdates(subID)
		if err != nil {
			log.Printf("Failed to unsubscribe from portfolio updates: %v", err)
			return
		}
		delete(c.subscriptions, sub.Type)
	}
}
