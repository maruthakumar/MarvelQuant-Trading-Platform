package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/trading-platform/backend/internal/models"
)

// OrderUpdateService handles real-time order updates
type OrderUpdateService struct {
	hub *Hub
}

// NewOrderUpdateService creates a new OrderUpdateService
func NewOrderUpdateService(hub *Hub) *OrderUpdateService {
	return &OrderUpdateService{
		hub: hub,
	}
}

// BroadcastOrderUpdate sends an order update to all subscribed clients
func (s *OrderUpdateService) BroadcastOrderUpdate(order *models.Order) error {
	// Create order update payload
	orderUpdate := struct {
		OrderID    string           `json:"orderId"`
		UserID     string           `json:"userId"`
		Symbol     string           `json:"symbol"`
		Side       models.OrderSide `json:"side"`
		Quantity   int              `json:"quantity"`
		Price      float64          `json:"price"`
		Status     models.OrderStatus `json:"status"`
		UpdatedAt  time.Time        `json:"updatedAt"`
		StrategyID string           `json:"strategyId,omitempty"`
	}{
		OrderID:    order.ID,
		UserID:     order.UserID,
		Symbol:     order.Symbol,
		Side:       order.Side,
		Quantity:   order.Quantity,
		Price:      order.Price,
		Status:     order.Status,
		UpdatedAt:  order.UpdatedAt,
		StrategyID: order.StrategyID,
	}

	// Marshal the order update
	payload, err := json.Marshal(orderUpdate)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := WebSocketMessage{
		Type:      MessageTypeOrderUpdate,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	// Marshal the message
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Broadcast to order topic
	s.hub.BroadcastToTopic("orders", messageJSON)

	// Broadcast to user-specific topic
	s.hub.BroadcastToTopic("user:"+order.UserID+":orders", messageJSON)

	// Broadcast to strategy-specific topic if applicable
	if order.StrategyID != "" {
		s.hub.BroadcastToTopic("strategy:"+order.StrategyID+":orders", messageJSON)
	}

	return nil
}

// PositionUpdateService handles real-time position updates
type PositionUpdateService struct {
	hub *Hub
}

// NewPositionUpdateService creates a new PositionUpdateService
func NewPositionUpdateService(hub *Hub) *PositionUpdateService {
	return &PositionUpdateService{
		hub: hub,
	}
}

// BroadcastPositionUpdate sends a position update to all subscribed clients
func (s *PositionUpdateService) BroadcastPositionUpdate(position *models.Position) error {
	// Create position update payload
	positionUpdate := struct {
		PositionID   string              `json:"positionId"`
		UserID       string              `json:"userId"`
		Symbol       string              `json:"symbol"`
		Quantity     int                 `json:"quantity"`
		EntryPrice   float64             `json:"entryPrice"`
		CurrentPrice float64             `json:"currentPrice,omitempty"`
		UnrealizedPnL float64            `json:"unrealizedPnL"`
		RealizedPnL  float64             `json:"realizedPnL"`
		Status       models.PositionStatus `json:"status"`
		UpdatedAt    time.Time           `json:"updatedAt"`
		StrategyID   string              `json:"strategyId,omitempty"`
	}{
		PositionID:   position.ID,
		UserID:       position.UserID,
		Symbol:       position.Symbol,
		Quantity:     position.Quantity,
		EntryPrice:   position.EntryPrice,
		CurrentPrice: position.CurrentPrice,
		UnrealizedPnL: position.UnrealizedPnL,
		RealizedPnL:  position.RealizedPnL,
		Status:       position.Status,
		UpdatedAt:    position.UpdatedAt,
		StrategyID:   position.StrategyID,
	}

	// Marshal the position update
	payload, err := json.Marshal(positionUpdate)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := WebSocketMessage{
		Type:      MessageTypePositionUpdate,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	// Marshal the message
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Broadcast to positions topic
	s.hub.BroadcastToTopic("positions", messageJSON)

	// Broadcast to user-specific topic
	s.hub.BroadcastToTopic("user:"+position.UserID+":positions", messageJSON)

	// Broadcast to strategy-specific topic if applicable
	if position.StrategyID != "" {
		s.hub.BroadcastToTopic("strategy:"+position.StrategyID+":positions", messageJSON)
	}

	return nil
}

// StrategyMonitorService handles real-time strategy monitoring
type StrategyMonitorService struct {
	hub *Hub
}

// NewStrategyMonitorService creates a new StrategyMonitorService
func NewStrategyMonitorService(hub *Hub) *StrategyMonitorService {
	return &StrategyMonitorService{
		hub: hub,
	}
}

// BroadcastStrategyUpdate sends a strategy update to all subscribed clients
func (s *StrategyMonitorService) BroadcastStrategyUpdate(strategy *models.Strategy) error {
	// Create strategy update payload
	strategyUpdate := struct {
		StrategyID  string    `json:"strategyId"`
		UserID      string    `json:"userId"`
		Name        string    `json:"name"`
		Status      string    `json:"status"`
		Performance struct {
			TotalPnL    float64 `json:"totalPnL"`
			WinRate     float64 `json:"winRate"`
			TotalTrades int     `json:"totalTrades"`
		} `json:"performance"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}{
		StrategyID:  strategy.ID,
		UserID:      strategy.UserID,
		Name:        strategy.Name,
		Status:      strategy.Status,
		UpdatedAt:   strategy.UpdatedAt,
	}

	// Set performance metrics if available
	if strategy.Performance != nil {
		strategyUpdate.Performance.TotalPnL = strategy.Performance.TotalPnL
		strategyUpdate.Performance.WinRate = strategy.Performance.WinRate
		strategyUpdate.Performance.TotalTrades = strategy.Performance.TotalTrades
	}

	// Marshal the strategy update
	payload, err := json.Marshal(strategyUpdate)
	if err != nil {
		return err
	}

	// Create WebSocket message
	message := WebSocketMessage{
		Type:      MessageTypeStrategyUpdate,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	// Marshal the message
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Broadcast to strategies topic
	s.hub.BroadcastToTopic("strategies", messageJSON)

	// Broadcast to user-specific topic
	s.hub.BroadcastToTopic("user:"+strategy.UserID+":strategies", messageJSON)

	// Broadcast to strategy-specific topic
	s.hub.BroadcastToTopic("strategy:"+strategy.ID, messageJSON)

	return nil
}

// ConnectionManager handles WebSocket connection management
type ConnectionManager struct {
	hub *Hub
}

// NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(hub *Hub) *ConnectionManager {
	return &ConnectionManager{
		hub: hub,
	}
}

// GetActiveConnections returns the number of active connections
func (m *ConnectionManager) GetActiveConnections() int {
	return len(m.hub.clients)
}

// GetTopicSubscriptions returns the number of subscriptions for a topic
func (m *ConnectionManager) GetTopicSubscriptions(topic string) int {
	m.hub.mu.Lock()
	defer m.hub.mu.Unlock()

	if clients, ok := m.hub.topics[topic]; ok {
		return len(clients)
	}
	return 0
}

// GetUserConnections returns the number of connections for a user
func (m *ConnectionManager) GetUserConnections(userID string) int {
	count := 0
	for client := range m.hub.clients {
		if client.userID == userID {
			count++
		}
	}
	return count
}

// DisconnectUser disconnects all connections for a user
func (m *ConnectionManager) DisconnectUser(userID string) {
	for client := range m.hub.clients {
		if client.userID == userID {
			client.conn.Close()
		}
	}
}

// AuthenticationMiddleware provides authentication for WebSocket connections
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from query parameter or Authorization header
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
			// Remove "Bearer " prefix if present
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		// Validate token
		// In a real implementation, this would involve checking the token against a database
		// For now, we'll just check if it's not empty
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token
		// In a real implementation, this would involve decoding the token
		// For now, we'll just use a placeholder
		userID := "user123" // Placeholder

		// Set user ID in request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
