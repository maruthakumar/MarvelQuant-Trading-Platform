package websocket

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/models"
)

// TestWebSocketHub tests the WebSocket hub functionality
func TestWebSocketHub(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a test client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "test-user",
		topics: make(map[string]bool),
	}
	
	// Register the client
	hub.register <- client
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check that the client is registered
	assert.Contains(t, hub.clients, client)
	
	// Subscribe the client to a topic
	hub.Subscribe(client, "test-topic")
	
	// Check that the client is subscribed to the topic
	assert.Contains(t, hub.topics, "test-topic")
	assert.Contains(t, hub.topics["test-topic"], client)
	assert.Contains(t, client.topics, "test-topic")
	
	// Broadcast a message to the topic
	message := []byte("test message")
	hub.BroadcastToTopic("test-topic", message)
	
	// Wait for the message to be sent
	time.Sleep(100 * time.Millisecond)
	
	// Check that the message was received
	select {
	case received := <-client.send:
		assert.Equal(t, message, received)
	default:
		t.Fatal("Message not received")
	}
	
	// Unsubscribe the client from the topic
	hub.Unsubscribe(client, "test-topic")
	
	// Check that the client is unsubscribed from the topic
	assert.NotContains(t, hub.topics["test-topic"], client)
	assert.NotContains(t, client.topics, "test-topic")
	
	// Unregister the client
	hub.unregister <- client
	
	// Wait for unregistration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check that the client is unregistered
	assert.NotContains(t, hub.clients, client)
}

// TestOrderUpdateService tests the OrderUpdateService
func TestOrderUpdateService(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a test client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "test-user",
		topics: make(map[string]bool),
	}
	
	// Register the client
	hub.register <- client
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Subscribe the client to the orders topic
	hub.Subscribe(client, "orders")
	
	// Create an order update service
	service := NewOrderUpdateService(hub)
	
	// Create a test order
	order := &models.Order{
		ID:        "order123",
		UserID:    "test-user",
		Symbol:    "AAPL",
		Side:      models.OrderSideBuy,
		Quantity:  10,
		Price:     150.0,
		Status:    models.OrderStatusFilled,
		UpdatedAt: time.Now(),
	}
	
	// Broadcast an order update
	err := service.BroadcastOrderUpdate(order)
	assert.NoError(t, err)
	
	// Wait for the message to be sent
	time.Sleep(100 * time.Millisecond)
	
	// Check that the message was received
	select {
	case received := <-client.send:
		// Parse the message
		var message WebSocketMessage
		err := json.Unmarshal(received, &message)
		assert.NoError(t, err)
		
		// Check the message type
		assert.Equal(t, MessageTypeOrderUpdate, message.Type)
		
		// Parse the payload
		var orderUpdate struct {
			OrderID  string `json:"orderId"`
			UserID   string `json:"userId"`
			Symbol   string `json:"symbol"`
			Quantity int    `json:"quantity"`
		}
		err = json.Unmarshal(message.Payload, &orderUpdate)
		assert.NoError(t, err)
		
		// Check the order details
		assert.Equal(t, order.ID, orderUpdate.OrderID)
		assert.Equal(t, order.UserID, orderUpdate.UserID)
		assert.Equal(t, order.Symbol, orderUpdate.Symbol)
		assert.Equal(t, order.Quantity, orderUpdate.Quantity)
	default:
		t.Fatal("Message not received")
	}
}

// TestPositionUpdateService tests the PositionUpdateService
func TestPositionUpdateService(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a test client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "test-user",
		topics: make(map[string]bool),
	}
	
	// Register the client
	hub.register <- client
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Subscribe the client to the positions topic
	hub.Subscribe(client, "positions")
	
	// Create a position update service
	service := NewPositionUpdateService(hub)
	
	// Create a test position
	position := &models.Position{
		ID:           "position123",
		UserID:       "test-user",
		Symbol:       "AAPL",
		Quantity:     10,
		EntryPrice:   150.0,
		CurrentPrice: 160.0,
		UnrealizedPnL: 100.0,
		RealizedPnL:  0.0,
		Status:       models.PositionStatusOpen,
		UpdatedAt:    time.Now(),
	}
	
	// Broadcast a position update
	err := service.BroadcastPositionUpdate(position)
	assert.NoError(t, err)
	
	// Wait for the message to be sent
	time.Sleep(100 * time.Millisecond)
	
	// Check that the message was received
	select {
	case received := <-client.send:
		// Parse the message
		var message WebSocketMessage
		err := json.Unmarshal(received, &message)
		assert.NoError(t, err)
		
		// Check the message type
		assert.Equal(t, MessageTypePositionUpdate, message.Type)
		
		// Parse the payload
		var positionUpdate struct {
			PositionID   string  `json:"positionId"`
			UserID       string  `json:"userId"`
			Symbol       string  `json:"symbol"`
			Quantity     int     `json:"quantity"`
			EntryPrice   float64 `json:"entryPrice"`
			CurrentPrice float64 `json:"currentPrice"`
			UnrealizedPnL float64 `json:"unrealizedPnL"`
		}
		err = json.Unmarshal(message.Payload, &positionUpdate)
		assert.NoError(t, err)
		
		// Check the position details
		assert.Equal(t, position.ID, positionUpdate.PositionID)
		assert.Equal(t, position.UserID, positionUpdate.UserID)
		assert.Equal(t, position.Symbol, positionUpdate.Symbol)
		assert.Equal(t, position.Quantity, positionUpdate.Quantity)
		assert.Equal(t, position.EntryPrice, positionUpdate.EntryPrice)
		assert.Equal(t, position.CurrentPrice, positionUpdate.CurrentPrice)
		assert.Equal(t, position.UnrealizedPnL, positionUpdate.UnrealizedPnL)
	default:
		t.Fatal("Message not received")
	}
}

// TestStrategyMonitorService tests the StrategyMonitorService
func TestStrategyMonitorService(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a test client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "test-user",
		topics: make(map[string]bool),
	}
	
	// Register the client
	hub.register <- client
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Subscribe the client to the strategies topic
	hub.Subscribe(client, "strategies")
	
	// Create a strategy monitor service
	service := NewStrategyMonitorService(hub)
	
	// Create a test strategy
	strategy := &models.Strategy{
		ID:        "strategy123",
		UserID:    "test-user",
		Name:      "Test Strategy",
		Status:    "ACTIVE",
		UpdatedAt: time.Now(),
		Performance: &models.StrategyPerformance{
			TotalPnL:    100.0,
			WinRate:     75.0,
			TotalTrades: 4,
		},
	}
	
	// Broadcast a strategy update
	err := service.BroadcastStrategyUpdate(strategy)
	assert.NoError(t, err)
	
	// Wait for the message to be sent
	time.Sleep(100 * time.Millisecond)
	
	// Check that the message was received
	select {
	case received := <-client.send:
		// Parse the message
		var message WebSocketMessage
		err := json.Unmarshal(received, &message)
		assert.NoError(t, err)
		
		// Check the message type
		assert.Equal(t, MessageTypeStrategyUpdate, message.Type)
		
		// Parse the payload
		var strategyUpdate struct {
			StrategyID  string `json:"strategyId"`
			UserID      string `json:"userId"`
			Name        string `json:"name"`
			Status      string `json:"status"`
			Performance struct {
				TotalPnL    float64 `json:"totalPnL"`
				WinRate     float64 `json:"winRate"`
				TotalTrades int     `json:"totalTrades"`
			} `json:"performance"`
		}
		err = json.Unmarshal(message.Payload, &strategyUpdate)
		assert.NoError(t, err)
		
		// Check the strategy details
		assert.Equal(t, strategy.ID, strategyUpdate.StrategyID)
		assert.Equal(t, strategy.UserID, strategyUpdate.UserID)
		assert.Equal(t, strategy.Name, strategyUpdate.Name)
		assert.Equal(t, strategy.Status, strategyUpdate.Status)
		assert.Equal(t, strategy.Performance.TotalPnL, strategyUpdate.Performance.TotalPnL)
		assert.Equal(t, strategy.Performance.WinRate, strategyUpdate.Performance.WinRate)
		assert.Equal(t, strategy.Performance.TotalTrades, strategyUpdate.Performance.TotalTrades)
	default:
		t.Fatal("Message not received")
	}
}

// TestConnectionManager tests the ConnectionManager
func TestConnectionManager(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a connection manager
	manager := NewConnectionManager(hub)
	
	// Check initial connection count
	assert.Equal(t, 0, manager.GetActiveConnections())
	
	// Create test clients
	client1 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "user1",
		topics: make(map[string]bool),
	}
	
	client2 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "user2",
		topics: make(map[string]bool),
	}
	
	client3 := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "user1",
		topics: make(map[string]bool),
	}
	
	// Register the clients
	hub.register <- client1
	hub.register <- client2
	hub.register <- client3
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check connection count
	assert.Equal(t, 3, manager.GetActiveConnections())
	
	// Check user connection count
	assert.Equal(t, 2, manager.GetUserConnections("user1"))
	assert.Equal(t, 1, manager.GetUserConnections("user2"))
	
	// Subscribe clients to topics
	hub.Subscribe(client1, "topic1")
	hub.Subscribe(client2, "topic1")
	hub.Subscribe(client3, "topic2")
	
	// Wait for subscriptions to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check topic subscription count
	assert.Equal(t, 2, manager.GetTopicSubscriptions("topic1"))
	assert.Equal(t, 1, manager.GetTopicSubscriptions("topic2"))
}

// TestAuthenticationMiddleware tests the AuthenticationMiddleware
func TestAuthenticationMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value("userID").(string)
		assert.True(t, ok)
		assert.Equal(t, "user123", userID)
		
		w.WriteHeader(http.StatusOK)
	})
	
	// Create the middleware
	middleware := AuthenticationMiddleware(testHandler)
	
	// Create a test server
	server := httptest.NewServer(middleware)
	defer server.Close()
	
	// Test with token in query parameter
	req, err := http.NewRequest("GET", server.URL+"?token=valid-token", nil)
	assert.NoError(t, err)
	
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test with token in Authorization header
	req, err = http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer valid-token")
	
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test with no token
	req, err = http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)
	
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestWebSocketHandler tests the WebSocketHandler
func TestWebSocketHandler(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a WebSocket handler
	handler := NewWebSocketHandler(hub)
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set user ID in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", "test-user")
		r = r.WithContext(ctx)
		
		// Handle WebSocket
		handler.HandleWebSocket(w, r)
	}))
	defer server.Close()
	
	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	
	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()
	
	// Send a subscription message
	subscriptionMsg := WebSocketMessage{
		Type:      MessageTypeSubscription,
		Timestamp: time.Now(),
		Payload:   json.RawMessage(`{"action":"subscribe","topics":["orders"]}`),
	}
	
	subscriptionJSON, err := json.Marshal(subscriptionMsg)
	assert.NoError(t, err)
	
	err = ws.WriteMessage(websocket.TextMessage, subscriptionJSON)
	assert.NoError(t, err)
	
	// Wait for subscription to complete
	time.Sleep(100 * time.Millisecond)
	
	// Send a heartbeat message
	heartbeatMsg := WebSocketMessage{
		Type:      MessageTypeHeartbeat,
		Timestamp: time.Now(),
		Payload:   json.RawMessage(`{}`),
	}
	
	heartbeatJSON, err := json.Marshal(heartbeatMsg)
	assert.NoError(t, err)
	
	err = ws.WriteMessage(websocket.TextMessage, heartbeatJSON)
	assert.NoError(t, err)
	
	// Wait for response
	_, response, err := ws.ReadMessage()
	assert.NoError(t, err)
	
	// Parse the response
	var responseMsg WebSocketMessage
	err = json.Unmarshal(response, &responseMsg)
	assert.NoError(t, err)
	
	// Check the response type
	assert.Equal(t, MessageTypeHeartbeat, responseMsg.Type)
}

// TestWebSocketIntegration tests the WebSocketIntegration
func TestWebSocketIntegration(t *testing.T) {
	// Create a new hub
	hub := NewHub()
	
	// Start the hub
	go hub.Run()
	
	// Create a WebSocket integration
	integration := NewWebSocketIntegration(hub)
	
	// Create a test client
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		userID: "test-user",
		topics: make(map[string]bool),
	}
	
	// Register the client
	hub.register <- client
	
	// Wait for registration to complete
	time.Sleep(100 * time.Millisecond)
	
	// Subscribe the client to topics
	hub.Subscribe(client, "orders")
	hub.Subscribe(client, "positions")
	hub.Subscribe(client, "strategies")
	
	// Create test data
	order := &models.Order{
		ID:        "order123",
		UserID:    "test-user",
		Symbol:    "AAPL",
		Side:      models.OrderSideBuy,
		Quantity:  10,
		Price:     150.0,
		Status:    models.OrderStatusFilled,
		UpdatedAt: time.Now(),
	}
	
	position := &models.Position{
		ID:           "position123",
		UserID:       "test-user",
		Symbol:       "AAPL",
		Quantity:     10,
		EntryPrice:   150.0,
		CurrentPrice: 160.0,
		UnrealizedPnL: 100.0,
		RealizedPnL:  0.0,
		Status:       models.PositionStatusOpen,
		UpdatedAt:    time.Now(),
	}
	
	strategy := &models.Strategy{
		ID:        "strategy123",
		UserID:    "test-user",
		Name:      "Test Strategy",
		Status:    "ACTIVE",
		UpdatedAt: time.Now(),
		Performance: &models.StrategyPerformance{
			TotalPnL:    100.0,
			WinRate:     75.0,
			TotalTrades: 4,
		},
	}
	
	// Send notifications
	integration.NotifyOrderUpdate(order)
	integration.NotifyPositionUpdate(position)
	integration.NotifyStrategyUpdate(strategy)
	
	// Wait for messages to be sent
	time.Sleep(100 * time.Millisecond)
	
	// Check that messages were received
	messageCount := 0
	for i := 0; i < 3; i++ {
		select {
		case <-client.send:
			messageCount++
		default:
			// No more messages
		}
	}
	
	assert.Equal(t, 3, messageCount)
}
