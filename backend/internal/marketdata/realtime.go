package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RealTimeUpdateManager manages real-time updates for market data
type RealTimeUpdateManager struct {
	dataSourceManager *DataSourceManager
	dataStorage       DataStorage
	cacheManager      CacheManager
	processors        []DataProcessor
	subscribers       map[string][]MarketDataCallback
	subscribersMu     sync.RWMutex
	wsConnections     map[string]*websocket.Conn
	wsConnectionsMu   sync.RWMutex
}

// NewRealTimeUpdateManager creates a new real-time update manager
func NewRealTimeUpdateManager(
	dataSourceManager *DataSourceManager,
	dataStorage DataStorage,
	cacheManager CacheManager,
	processors ...DataProcessor,
) *RealTimeUpdateManager {
	return &RealTimeUpdateManager{
		dataSourceManager: dataSourceManager,
		dataStorage:       dataStorage,
		cacheManager:      cacheManager,
		processors:        processors,
		subscribers:       make(map[string][]MarketDataCallback),
		wsConnections:     make(map[string]*websocket.Conn),
	}
}

// Subscribe subscribes to real-time updates for the specified symbols
func (m *RealTimeUpdateManager) Subscribe(ctx context.Context, symbols []string, callback MarketDataCallback) error {
	// Register callback
	m.subscribersMu.Lock()
	for _, symbol := range symbols {
		m.subscribers[symbol] = append(m.subscribers[symbol], callback)
	}
	m.subscribersMu.Unlock()

	// Create a wrapper callback that processes the data before distributing it to subscribers
	wrapperCallback := func(data MarketData) {
		processedData := data

		// Apply each processor
		for _, processor := range m.processors {
			result, err := processor.Process(processedData)
			if err != nil {
				// Log error and continue with unprocessed data
				log.Printf("Error processing market data: %v", err)
				continue
			}
			if processed, ok := result.(MarketData); ok {
				processedData = processed
			}
		}

		// Store processed data in cache
		key := "market_data:" + data.Symbol
		m.cacheManager.Set(key, processedData, 5*time.Second)

		// Store processed data in storage (async)
		go func() {
			if err := m.dataStorage.StoreMarketData(context.Background(), processedData); err != nil {
				log.Printf("Error storing market data: %v", err)
			}
		}()

		// Distribute to subscribers
		m.subscribersMu.RLock()
		subscribers := m.subscribers[data.Symbol]
		m.subscribersMu.RUnlock()

		for _, subscriber := range subscribers {
			subscriber(processedData)
		}
	}

	// Subscribe to market data from data source
	return m.dataSourceManager.SubscribeToMarketData(ctx, symbols, wrapperCallback)
}

// Unsubscribe unsubscribes from real-time updates for the specified symbols
func (m *RealTimeUpdateManager) Unsubscribe(ctx context.Context, symbols []string) error {
	// Unregister callbacks
	m.subscribersMu.Lock()
	for _, symbol := range symbols {
		delete(m.subscribers, symbol)
	}
	m.subscribersMu.Unlock()

	// Unsubscribe from market data from data source
	return m.dataSourceManager.UnsubscribeFromMarketData(ctx, symbols)
}

// GetSubscribedSymbols gets the list of symbols with active subscriptions
func (m *RealTimeUpdateManager) GetSubscribedSymbols() []string {
	m.subscribersMu.RLock()
	defer m.subscribersMu.RUnlock()

	symbols := make([]string, 0, len(m.subscribers))
	for symbol := range m.subscribers {
		symbols = append(symbols, symbol)
	}

	return symbols
}

// WebSocketHandler handles WebSocket connections for real-time updates
type WebSocketHandler struct {
	updateManager *RealTimeUpdateManager
	connections   map[*websocket.Conn]bool
	connectionsMu sync.RWMutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(updateManager *RealTimeUpdateManager) *WebSocketHandler {
	return &WebSocketHandler{
		updateManager: updateManager,
		connections:   make(map[*websocket.Conn]bool),
	}
}

// HandleConnection handles a WebSocket connection
func (h *WebSocketHandler) HandleConnection(conn *websocket.Conn) {
	// Register connection
	h.connectionsMu.Lock()
	h.connections[conn] = true
	h.connectionsMu.Unlock()

	// Clean up on disconnect
	defer func() {
		h.connectionsMu.Lock()
		delete(h.connections, conn)
		h.connectionsMu.Unlock()
		conn.Close()
	}()

	// Create a context for this connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Parse message
		var msg struct {
			Action  string   `json:"action"`
			Symbols []string `json:"symbols"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing WebSocket message: %v", err)
			continue
		}

		// Handle message based on action
		switch msg.Action {
		case "subscribe":
			// Create a callback that sends data to this connection
			callback := func(data MarketData) {
				// Convert data to JSON
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshaling market data: %v", err)
					return
				}

				// Send data to connection
				if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("Error writing to WebSocket: %v", err)
					return
				}
			}

			// Subscribe to symbols
			if err := h.updateManager.Subscribe(ctx, msg.Symbols, callback); err != nil {
				log.Printf("Error subscribing to symbols: %v", err)
				continue
			}

			// Send confirmation
			response := struct {
				Action  string   `json:"action"`
				Status  string   `json:"status"`
				Symbols []string `json:"symbols"`
			}{
				Action:  "subscribe",
				Status:  "success",
				Symbols: msg.Symbols,
			}
			jsonResponse, _ := json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				continue
			}

		case "unsubscribe":
			// Unsubscribe from symbols
			if err := h.updateManager.Unsubscribe(ctx, msg.Symbols); err != nil {
				log.Printf("Error unsubscribing from symbols: %v", err)
				continue
			}

			// Send confirmation
			response := struct {
				Action  string   `json:"action"`
				Status  string   `json:"status"`
				Symbols []string `json:"symbols"`
			}{
				Action:  "unsubscribe",
				Status:  "success",
				Symbols: msg.Symbols,
			}
			jsonResponse, _ := json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				continue
			}

		default:
			// Send error for unknown action
			response := struct {
				Action  string `json:"action"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Action:  msg.Action,
				Status:  "error",
				Message: fmt.Sprintf("Unknown action: %s", msg.Action),
			}
			jsonResponse, _ := json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				continue
			}
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients
func (h *WebSocketHandler) BroadcastMessage(message []byte) {
	h.connectionsMu.RLock()
	defer h.connectionsMu.RUnlock()

	for conn := range h.connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}

// GetConnectionCount gets the number of active connections
func (h *WebSocketHandler) GetConnectionCount() int {
	h.connectionsMu.RLock()
	defer h.connectionsMu.RUnlock()
	return len(h.connections)
}

// MarketDataEventBus implements a publish-subscribe pattern for market data events
type MarketDataEventBus struct {
	subscribers map[string][]MarketDataCallback
	mutex       sync.RWMutex
}

// NewMarketDataEventBus creates a new market data event bus
func NewMarketDataEventBus() *MarketDataEventBus {
	return &MarketDataEventBus{
		subscribers: make(map[string][]MarketDataCallback),
	}
}

// Subscribe subscribes to market data events for the specified topic
func (b *MarketDataEventBus) Subscribe(topic string, callback MarketDataCallback) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], callback)
}

// Unsubscribe unsubscribes from market data events for the specified topic
func (b *MarketDataEventBus) Unsubscribe(topic string, callback MarketDataCallback) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscribers := b.subscribers[topic]
	for i, sub := range subscribers {
		if fmt.Sprintf("%p", sub) == fmt.Sprintf("%p", callback) {
			// Remove subscriber
			b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}

// Publish publishes a market data event to the specified topic
func (b *MarketDataEventBus) Publish(topic string, data MarketData) {
	b.mutex.RLock()
	subscribers := b.subscribers[topic]
	b.mutex.RUnlock()

	for _, subscriber := range subscribers {
		subscriber(data)
	}
}

// GetSubscriberCount gets the number of subscribers for the specified topic
func (b *MarketDataEventBus) GetSubscriberCount(topic string) int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return len(b.subscribers[topic])
}

// MarketDataStreamProcessor processes streams of market data
type MarketDataStreamProcessor struct {
	processors []DataProcessor
	eventBus   *MarketDataEventBus
}

// NewMarketDataStreamProcessor creates a new market data stream processor
func NewMarketDataStreamProcessor(eventBus *MarketDataEventBus, processors ...DataProcessor) *MarketDataStreamProcessor {
	return &MarketDataStreamProcessor{
		processors: processors,
		eventBus:   eventBus,
	}
}

// ProcessStream processes a stream of market data
func (p *MarketDataStreamProcessor) ProcessStream(ctx context.Context, inputChan <-chan MarketData) (<-chan MarketData, error) {
	outputChan := make(chan MarketData, 100)

	go func() {
		defer close(outputChan)

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-inputChan:
				if !ok {
					return
				}

				// Process data
				processedData := data
				for _, processor := range p.processors {
					result, err := processor.Process(processedData)
					if err != nil {
						log.Printf("Error processing market data: %v", err)
						continue
					}
					if processed, ok := result.(MarketData); ok {
						processedData = processed
					}
				}

				// Publish to event bus
				p.eventBus.Publish(data.Symbol, processedData)

				// Send to output channel
				select {
				case outputChan <- processedData:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return outputChan, nil
}

// AddProcessor adds a processor to the stream processor
func (p *MarketDataStreamProcessor) AddProcessor(processor DataProcessor) {
	p.processors = append(p.processors, processor)
}
