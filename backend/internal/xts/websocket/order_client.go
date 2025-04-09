package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/trade-execution-platform/backend/internal/xts/errors"
)

// OrderClient represents a WebSocket client for order updates
type OrderClient struct {
	conn          *websocket.Conn
	url           string
	token         string
	userID        string
	messageChan   chan []byte
	errorChan     chan error
	stopChan      chan struct{}
	reconnectChan chan struct{}
	mutex         sync.RWMutex
	isConnected   bool
	handlers      map[string]func([]byte)
}

// OrderMessage represents an order update message
type OrderMessage struct {
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data"`
}

// NewOrderClient creates a new order WebSocket client
func NewOrderClient(url, token, userID string) *OrderClient {
	client := &OrderClient{
		url:           url,
		token:         token,
		userID:        userID,
		messageChan:   make(chan []byte, 100),
		errorChan:     make(chan error, 10),
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
		handlers:      make(map[string]func([]byte)),
	}

	// Register default handlers
	client.registerDefaultHandlers()

	return client
}

// Connect establishes a WebSocket connection
func (c *OrderClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	// Construct connection URL with authentication parameters
	fullURL := fmt.Sprintf("%s/?token=%s&userID=%s&apiType=INTERACTIVE", 
		c.url, c.token, c.userID)

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(fullURL, nil)
	if err != nil {
		return errors.Wrap(err, "ws_connect_error", "Failed to connect to order WebSocket", 0)
	}

	c.conn = conn
	c.isConnected = true

	// Start message reader
	go c.readMessages()

	// Start reconnection handler
	go c.handleReconnection()

	log.Println("Order WebSocket connected successfully")
	return nil
}

// Disconnect closes the WebSocket connection
func (c *OrderClient) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return nil
	}

	// Signal to stop all goroutines
	close(c.stopChan)

	// Close WebSocket connection
	err := c.conn.WriteMessage(websocket.CloseMessage, 
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("Error sending close message: %v", err)
	}

	err = c.conn.Close()
	c.isConnected = false
	c.conn = nil

	log.Println("Order WebSocket disconnected")
	return err
}

// GetMessageChannel returns the channel for receiving order update messages
func (c *OrderClient) GetMessageChannel() <-chan []byte {
	return c.messageChan
}

// GetErrorChannel returns the channel for receiving errors
func (c *OrderClient) GetErrorChannel() <-chan error {
	return c.errorChan
}

// RegisterHandler registers a handler function for a specific message type
func (c *OrderClient) RegisterHandler(messageType string, handler func([]byte)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[messageType] = handler
}

// IsConnected returns whether the client is connected
func (c *OrderClient) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// readMessages reads messages from the WebSocket connection
func (c *OrderClient) readMessages() {
	defer func() {
		c.mutex.Lock()
		c.isConnected = false
		c.mutex.Unlock()
		
		// Trigger reconnection if not explicitly stopped
		select {
		case <-c.stopChan:
			// Do nothing, we're shutting down
		default:
			c.reconnectChan <- struct{}{}
		}
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			// Set read deadline
			err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				c.errorChan <- errors.Wrap(err, "ws_deadline_error", "Failed to set read deadline", 0)
				return
			}

			// Read message
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errorChan <- errors.Wrap(err, "ws_read_error", "Failed to read from WebSocket", 0)
				return
			}

			// Process message
			c.processMessage(message)
		}
	}
}

// processMessage processes a received WebSocket message
func (c *OrderClient) processMessage(message []byte) {
	// Parse message type
	var msg OrderMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.errorChan <- errors.Wrap(err, "json_error", "Failed to parse order message", 0)
		return
	}

	// Route message to appropriate handler
	c.mutex.RLock()
	handler, exists := c.handlers[msg.MessageType]
	c.mutex.RUnlock()

	if exists {
		handler(message)
	}

	// Send message to channel for external processing
	select {
	case c.messageChan <- message:
		// Message sent successfully
	default:
		// Channel is full, log warning
		log.Println("Warning: Order message channel is full, dropping message")
	}
}

// handleReconnection handles automatic reconnection
func (c *OrderClient) handleReconnection() {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-c.stopChan:
			return
		case <-c.reconnectChan:
			log.Println("Attempting to reconnect to order WebSocket...")
			
			// Wait before reconnecting
			time.Sleep(backoff)
			
			// Try to reconnect
			c.mutex.Lock()
			c.isConnected = false
			c.mutex.Unlock()
			
			err := c.Connect()
			if err != nil {
				log.Printf("Reconnection failed: %v", err)
				
				// Increase backoff time for next attempt
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				
				// Schedule another reconnection attempt
				go func() {
					time.Sleep(100 * time.Millisecond)
					c.reconnectChan <- struct{}{}
				}()
			} else {
				log.Println("Reconnected to order WebSocket")
				backoff = 1 * time.Second
			}
		}
	}
}

// registerDefaultHandlers registers default message handlers
func (c *OrderClient) registerDefaultHandlers() {
	// Handler for connection message
	c.RegisterHandler("connect", func(data []byte) {
		log.Println("Order Socket connected successfully!")
	})

	// Handler for joined message
	c.RegisterHandler("joined", func(data []byte) {
		log.Println("Order Socket joined successfully!")
	})

	// Handler for order updates
	c.RegisterHandler("order", func(data []byte) {
		log.Println("Received order update")
	})

	// Handler for trade updates
	c.RegisterHandler("trade", func(data []byte) {
		log.Println("Received trade update")
	})

	// Handler for position updates
	c.RegisterHandler("position", func(data []byte) {
		log.Println("Received position update")
	})

	// Handler for trade conversion updates
	c.RegisterHandler("tradeConversion", func(data []byte) {
		log.Println("Received trade conversion update")
	})

	// Handler for logout message
	c.RegisterHandler("logout", func(data []byte) {
		log.Println("User logged out!")
	})

	// Handler for disconnect
	c.RegisterHandler("disconnect", func(data []byte) {
		log.Println("Order Socket disconnected!")
	})

	// Handler for error message
	c.RegisterHandler("error", func(data []byte) {
		log.Println("Received error from order socket")
	})
}

// WithContext returns a new OrderClient with context support
func (c *OrderClient) WithContext(ctx context.Context) *OrderClient {
	// Create a copy of the client
	newClient := &OrderClient{
		url:           c.url,
		token:         c.token,
		userID:        c.userID,
		messageChan:   make(chan []byte, 100),
		errorChan:     make(chan error, 10),
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
		handlers:      make(map[string]func([]byte)),
	}

	// Copy handlers
	for k, v := range c.handlers {
		newClient.handlers[k] = v
	}

	// Monitor context cancellation
	go func() {
		<-ctx.Done()
		newClient.Disconnect()
	}()

	return newClient
}
