package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketServer manages WebSocket connections for market data streaming
type WebSocketServer struct {
	upgrader        websocket.Upgrader
	connections     map[*websocket.Conn]*WebSocketConnection
	connectionsMu   sync.RWMutex
	marketDataSvc   *MarketDataService
	realTimeManager *RealTimeUpdateManager
}

// WebSocketConnection represents a single WebSocket connection
type WebSocketConnection struct {
	conn           *websocket.Conn
	subscriptions  map[string]bool
	subscriptionMu sync.RWMutex
	send           chan []byte
	server         *WebSocketServer
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(marketDataSvc *MarketDataService, realTimeManager *RealTimeUpdateManager) *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
		connections:     make(map[*websocket.Conn]*WebSocketConnection),
		marketDataSvc:   marketDataSvc,
		realTimeManager: realTimeManager,
	}
}

// ServeWS handles WebSocket requests from clients
func (s *WebSocketServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Create new connection
	wsConn := &WebSocketConnection{
		conn:          conn,
		subscriptions: make(map[string]bool),
		send:          make(chan []byte, 256),
		server:        s,
	}

	// Register connection
	s.connectionsMu.Lock()
	s.connections[conn] = wsConn
	s.connectionsMu.Unlock()

	// Start goroutines for reading and writing
	go wsConn.writePump()
	go wsConn.readPump()

	// Send welcome message
	welcome := map[string]interface{}{
		"type":    "welcome",
		"message": "Connected to MarverQuant Market Data Service",
		"time":    time.Now().Format(time.RFC3339),
	}
	welcomeJSON, _ := json.Marshal(welcome)
	wsConn.send <- welcomeJSON
}

// BroadcastToAll sends a message to all connected clients
func (s *WebSocketServer) BroadcastToAll(message []byte) {
	s.connectionsMu.RLock()
	defer s.connectionsMu.RUnlock()

	for _, conn := range s.connections {
		select {
		case conn.send <- message:
		default:
			// Channel is full, close connection
			s.closeConnection(conn)
		}
	}
}

// BroadcastToSubscribers sends a message to clients subscribed to a symbol
func (s *WebSocketServer) BroadcastToSubscribers(symbol string, message []byte) {
	s.connectionsMu.RLock()
	defer s.connectionsMu.RUnlock()

	for _, conn := range s.connections {
		conn.subscriptionMu.RLock()
		subscribed := conn.subscriptions[symbol]
		conn.subscriptionMu.RUnlock()

		if subscribed {
			select {
			case conn.send <- message:
			default:
				// Channel is full, close connection
				s.closeConnection(conn)
			}
		}
	}
}

// closeConnection closes a WebSocket connection
func (s *WebSocketServer) closeConnection(conn *WebSocketConnection) {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()

	// Unsubscribe from all symbols
	symbols := make([]string, 0, len(conn.subscriptions))
	conn.subscriptionMu.Lock()
	for symbol := range conn.subscriptions {
		symbols = append(symbols, symbol)
	}
	conn.subscriptionMu.Unlock()

	if len(symbols) > 0 {
		s.realTimeManager.Unsubscribe(context.Background(), symbols)
	}

	// Close connection
	if _, ok := s.connections[conn.conn]; ok {
		close(conn.send)
		conn.conn.Close()
		delete(s.connections, conn.conn)
	}
}

// GetConnectionCount returns the number of active connections
func (s *WebSocketServer) GetConnectionCount() int {
	s.connectionsMu.RLock()
	defer s.connectionsMu.RUnlock()
	return len(s.connections)
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *WebSocketConnection) readPump() {
	defer func() {
		c.server.closeConnection(c)
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

		// Process message
		c.processMessage(message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *WebSocketConnection) writePump() {
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
				// Channel was closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
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

// processMessage processes a message from the client
func (c *WebSocketConnection) processMessage(message []byte) {
	// Parse message
	var msg struct {
		Type    string   `json:"type"`
		Symbols []string `json:"symbols,omitempty"`
		Symbol  string   `json:"symbol,omitempty"`
		Action  string   `json:"action,omitempty"`
	}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error parsing message: %v", err)
		c.sendError("Invalid message format")
		return
	}

	// Process based on message type
	switch msg.Type {
	case "subscribe":
		c.handleSubscribe(msg.Symbols)
	case "unsubscribe":
		c.handleUnsubscribe(msg.Symbols)
	case "quote":
		c.handleQuoteRequest(msg.Symbol)
	case "ping":
		c.handlePing()
	default:
		c.sendError(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleSubscribe handles a subscribe request
func (c *WebSocketConnection) handleSubscribe(symbols []string) {
	if len(symbols) == 0 {
		c.sendError("No symbols specified")
		return
	}

	// Filter out already subscribed symbols
	var newSymbols []string
	c.subscriptionMu.Lock()
	for _, symbol := range symbols {
		if !c.subscriptions[symbol] {
			newSymbols = append(newSymbols, symbol)
			c.subscriptions[symbol] = true
		}
	}
	c.subscriptionMu.Unlock()

	if len(newSymbols) == 0 {
		c.sendSuccess("subscribe", "Already subscribed to all specified symbols")
		return
	}

	// Create callback for real-time updates
	callback := func(data MarketData) {
		// Convert to JSON
		message, err := json.Marshal(map[string]interface{}{
			"type": "update",
			"data": data,
		})
		if err != nil {
			log.Printf("Error marshaling market data: %v", err)
			return
		}

		// Send to client
		select {
		case c.send <- message:
		default:
			// Channel is full, close connection
			c.server.closeConnection(c)
		}
	}

	// Subscribe to real-time updates
	if err := c.server.realTimeManager.Subscribe(context.Background(), newSymbols, callback); err != nil {
		log.Printf("Error subscribing to symbols: %v", err)
		c.sendError(fmt.Sprintf("Error subscribing to symbols: %v", err))
		return
	}

	c.sendSuccess("subscribe", fmt.Sprintf("Subscribed to %d symbols", len(newSymbols)))
}

// handleUnsubscribe handles an unsubscribe request
func (c *WebSocketConnection) handleUnsubscribe(symbols []string) {
	if len(symbols) == 0 {
		c.sendError("No symbols specified")
		return
	}

	// Filter out symbols that aren't subscribed
	var unsubSymbols []string
	c.subscriptionMu.Lock()
	for _, symbol := range symbols {
		if c.subscriptions[symbol] {
			unsubSymbols = append(unsubSymbols, symbol)
			delete(c.subscriptions, symbol)
		}
	}
	c.subscriptionMu.Unlock()

	if len(unsubSymbols) == 0 {
		c.sendSuccess("unsubscribe", "Not subscribed to any of the specified symbols")
		return
	}

	// Unsubscribe from real-time updates
	if err := c.server.realTimeManager.Unsubscribe(context.Background(), unsubSymbols); err != nil {
		log.Printf("Error unsubscribing from symbols: %v", err)
		c.sendError(fmt.Sprintf("Error unsubscribing from symbols: %v", err))
		return
	}

	c.sendSuccess("unsubscribe", fmt.Sprintf("Unsubscribed from %d symbols", len(unsubSymbols)))
}

// handleQuoteRequest handles a quote request
func (c *WebSocketConnection) handleQuoteRequest(symbol string) {
	if symbol == "" {
		c.sendError("No symbol specified")
		return
	}

	// Get market data
	data, err := c.server.marketDataSvc.GetMarketData(context.Background(), []string{symbol})
	if err != nil {
		log.Printf("Error getting market data: %v", err)
		c.sendError(fmt.Sprintf("Error getting market data: %v", err))
		return
	}

	// Check if we got data for the symbol
	quote, ok := data[symbol]
	if !ok {
		c.sendError(fmt.Sprintf("No data found for symbol: %s", symbol))
		return
	}

	// Send quote
	message, err := json.Marshal(map[string]interface{}{
		"type":  "quote",
		"quote": quote,
	})
	if err != nil {
		log.Printf("Error marshaling quote: %v", err)
		c.sendError(fmt.Sprintf("Error marshaling quote: %v", err))
		return
	}

	c.send <- message
}

// handlePing handles a ping request
func (c *WebSocketConnection) handlePing() {
	c.send <- []byte(`{"type":"pong","time":"` + time.Now().Format(time.RFC3339) + `"}`)
}

// sendError sends an error message to the client
func (c *WebSocketConnection) sendError(message string) {
	errorMsg, _ := json.Marshal(map[string]interface{}{
		"type":    "error",
		"message": message,
	})
	c.send <- errorMsg
}

// sendSuccess sends a success message to the client
func (c *WebSocketConnection) sendSuccess(action, message string) {
	successMsg, _ := json.Marshal(map[string]interface{}{
		"type":    "success",
		"action":  action,
		"message": message,
	})
	c.send <- successMsg
}

// MarketDataStreamingService manages streaming of market data
type MarketDataStreamingService struct {
	wsServer        *WebSocketServer
	marketDataSvc   *MarketDataService
	realTimeManager *RealTimeUpdateManager
	eventBus        *MarketDataEventBus
}

// NewMarketDataStreamingService creates a new market data streaming service
func NewMarketDataStreamingService(
	marketDataSvc *MarketDataService,
	realTimeManager *RealTimeUpdateManager,
	eventBus *MarketDataEventBus,
) *MarketDataStreamingService {
	wsServer := NewWebSocketServer(marketDataSvc, realTimeManager)
	
	return &MarketDataStreamingService{
		wsServer:        wsServer,
		marketDataSvc:   marketDataSvc,
		realTimeManager: realTimeManager,
		eventBus:        eventBus,
	}
}

// Start starts the streaming service
func (s *MarketDataStreamingService) Start() {
	// Subscribe to market data events
	s.eventBus.Subscribe("market_data", func(data MarketData) {
		// Convert to JSON
		message, err := json.Marshal(map[string]interface{}{
			"type": "update",
			"data": data,
		})
		if err != nil {
			log.Printf("Error marshaling market data: %v", err)
			return
		}

		// Broadcast to subscribers
		s.wsServer.BroadcastToSubscribers(data.Symbol, message)
	})

	log.Println("Market data streaming service started")
}

// Stop stops the streaming service
func (s *MarketDataStreamingService) Stop() {
	// Nothing to do for now
	log.Println("Market data streaming service stopped")
}

// GetWSServer returns the WebSocket server
func (s *MarketDataStreamingService) GetWSServer() *WebSocketServer {
	return s.wsServer
}

// HandleWebSocket handles a WebSocket connection
func (s *MarketDataStreamingService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	s.wsServer.ServeWS(w, r)
}

// GetConnectionCount returns the number of active connections
func (s *MarketDataStreamingService) GetConnectionCount() int {
	return s.wsServer.GetConnectionCount()
}
