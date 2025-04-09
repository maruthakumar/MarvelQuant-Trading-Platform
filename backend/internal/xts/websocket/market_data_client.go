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
	"github.com/trade-execution-platform/backend/internal/xts/models"
)

// MarketDataClient represents a WebSocket client for market data
type MarketDataClient struct {
	conn          *websocket.Conn
	url           string
	token         string
	userID        string
	subscriptions map[string]bool
	messageChan   chan []byte
	errorChan     chan error
	stopChan      chan struct{}
	reconnectChan chan struct{}
	mutex         sync.RWMutex
	isConnected   bool
	handlers      map[string]func([]byte)
}

// MarketDataMessage represents a market data message
type MarketDataMessage struct {
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data"`
}

// NewMarketDataClient creates a new market data WebSocket client
func NewMarketDataClient(url, token, userID string) *MarketDataClient {
	client := &MarketDataClient{
		url:           url,
		token:         token,
		userID:        userID,
		subscriptions: make(map[string]bool),
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
func (c *MarketDataClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	// Construct connection URL with authentication parameters
	fullURL := fmt.Sprintf("%s/?token=%s&userID=%s&publishFormat=JSON", 
		c.url, c.token, c.userID)

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(fullURL, nil)
	if err != nil {
		return errors.Wrap(err, "ws_connect_error", "Failed to connect to market data WebSocket", 0)
	}

	c.conn = conn
	c.isConnected = true

	// Start message reader
	go c.readMessages()

	// Start reconnection handler
	go c.handleReconnection()

	log.Println("Market Data WebSocket connected successfully")
	return nil
}

// Disconnect closes the WebSocket connection
func (c *MarketDataClient) Disconnect() error {
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

	log.Println("Market Data WebSocket disconnected")
	return err
}

// Subscribe subscribes to market data for instruments
func (c *MarketDataClient) Subscribe(instruments []models.Instrument) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return errors.ErrWebSocketClosed
	}

	// Prepare subscription request
	instrumentList := make([]map[string]interface{}, 0, len(instruments))
	for _, instrument := range instruments {
		instrumentMap := map[string]interface{}{
			"exchangeSegment":     instrument.ExchangeSegment,
			"exchangeInstrumentID": instrument.ExchangeInstrumentID,
		}
		instrumentList = append(instrumentList, instrumentMap)
		
		// Track subscriptions
		key := fmt.Sprintf("%s:%s", instrument.ExchangeSegment, instrument.ExchangeInstrumentID)
		c.subscriptions[key] = true
	}

	request := map[string]interface{}{
		"instruments": instrumentList,
		"xtsMessageCode": 1501, // Touchline data
	}

	// Send subscription request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "json_error", "Failed to marshal subscription request", 0)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return errors.Wrap(err, "ws_write_error", "Failed to send subscription request", 0)
	}

	log.Printf("Subscribed to %d instruments", len(instruments))
	return nil
}

// Unsubscribe unsubscribes from market data for instruments
func (c *MarketDataClient) Unsubscribe(instruments []models.Instrument) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return errors.ErrWebSocketClosed
	}

	// Prepare unsubscription request
	instrumentList := make([]map[string]interface{}, 0, len(instruments))
	for _, instrument := range instruments {
		instrumentMap := map[string]interface{}{
			"exchangeSegment":     instrument.ExchangeSegment,
			"exchangeInstrumentID": instrument.ExchangeInstrumentID,
		}
		instrumentList = append(instrumentList, instrumentMap)
		
		// Remove from subscriptions
		key := fmt.Sprintf("%s:%s", instrument.ExchangeSegment, instrument.ExchangeInstrumentID)
		delete(c.subscriptions, key)
	}

	request := map[string]interface{}{
		"instruments": instrumentList,
		"xtsMessageCode": 1501, // Touchline data
	}

	// Send unsubscription request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "json_error", "Failed to marshal unsubscription request", 0)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return errors.Wrap(err, "ws_write_error", "Failed to send unsubscription request", 0)
	}

	log.Printf("Unsubscribed from %d instruments", len(instruments))
	return nil
}

// GetMessageChannel returns the channel for receiving market data messages
func (c *MarketDataClient) GetMessageChannel() <-chan []byte {
	return c.messageChan
}

// GetErrorChannel returns the channel for receiving errors
func (c *MarketDataClient) GetErrorChannel() <-chan error {
	return c.errorChan
}

// RegisterHandler registers a handler function for a specific message type
func (c *MarketDataClient) RegisterHandler(messageType string, handler func([]byte)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[messageType] = handler
}

// IsConnected returns whether the client is connected
func (c *MarketDataClient) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// readMessages reads messages from the WebSocket connection
func (c *MarketDataClient) readMessages() {
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
func (c *MarketDataClient) processMessage(message []byte) {
	// Parse message type
	var msg MarketDataMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.errorChan <- errors.Wrap(err, "json_error", "Failed to parse market data message", 0)
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
		log.Println("Warning: Market data message channel is full, dropping message")
	}
}

// handleReconnection handles automatic reconnection
func (c *MarketDataClient) handleReconnection() {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-c.stopChan:
			return
		case <-c.reconnectChan:
			log.Println("Attempting to reconnect to market data WebSocket...")
			
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
				log.Println("Reconnected to market data WebSocket")
				backoff = 1 * time.Second
				
				// Resubscribe to previous subscriptions
				c.resubscribe()
			}
		}
	}
}

// resubscribe resubscribes to previously subscribed instruments
func (c *MarketDataClient) resubscribe() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.subscriptions) == 0 {
		return
	}

	// Convert subscriptions map to instruments slice
	instruments := make([]models.Instrument, 0, len(c.subscriptions))
	for key := range c.subscriptions {
		var segment, instrumentID string
		fmt.Sscanf(key, "%s:%s", &segment, &instrumentID)
		
		instrument := models.Instrument{
			ExchangeSegment:     segment,
			ExchangeInstrumentID: instrumentID,
		}
		instruments = append(instruments, instrument)
	}

	// Resubscribe
	go func() {
		err := c.Subscribe(instruments)
		if err != nil {
			log.Printf("Failed to resubscribe: %v", err)
		} else {
			log.Printf("Resubscribed to %d instruments", len(instruments))
		}
	}()
}

// registerDefaultHandlers registers default message handlers
func (c *MarketDataClient) registerDefaultHandlers() {
	// Handler for connection message
	c.RegisterHandler("connect", func(data []byte) {
		log.Println("Market Data Socket connected successfully!")
	})

	// Handler for touchline data (1501)
	c.RegisterHandler("1501-json-full", func(data []byte) {
		log.Println("Received touchline data (full)")
	})

	// Handler for touchline data partial (1501)
	c.RegisterHandler("1501-json-partial", func(data []byte) {
		log.Println("Received touchline data (partial)")
	})

	// Handler for market depth (1502)
	c.RegisterHandler("1502-json-full", func(data []byte) {
		log.Println("Received market depth data (full)")
	})

	// Handler for market depth partial (1502)
	c.RegisterHandler("1502-json-partial", func(data []byte) {
		log.Println("Received market depth data (partial)")
	})

	// Handler for candle data (1505)
	c.RegisterHandler("1505-json-full", func(data []byte) {
		log.Println("Received candle data (full)")
	})

	// Handler for open interest (1510)
	c.RegisterHandler("1510-json-full", func(data []byte) {
		log.Println("Received open interest data (full)")
	})

	// Handler for LTP data (1512)
	c.RegisterHandler("1512-json-full", func(data []byte) {
		log.Println("Received LTP data (full)")
	})

	// Handler for disconnect
	c.RegisterHandler("disconnect", func(data []byte) {
		log.Println("Market Data Socket disconnected!")
	})
}

// WithContext returns a new MarketDataClient with context support
func (c *MarketDataClient) WithContext(ctx context.Context) *MarketDataClient {
	// Create a copy of the client
	newClient := &MarketDataClient{
		url:           c.url,
		token:         c.token,
		userID:        c.userID,
		subscriptions: make(map[string]bool),
		messageChan:   make(chan []byte, 100),
		errorChan:     make(chan error, 10),
		stopChan:      make(chan struct{}),
		reconnectChan: make(chan struct{}, 1),
		handlers:      make(map[string]func([]byte)),
	}

	// Copy subscriptions
	for k, v := range c.subscriptions {
		newClient.subscriptions[k] = v
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
