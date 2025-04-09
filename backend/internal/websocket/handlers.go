package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
)

// WebSocketHandler handles HTTP requests for WebSocket connections
type WebSocketHandler struct {
	hub *Hub
	orderUpdateService *OrderUpdateService
	positionUpdateService *PositionUpdateService
	strategyMonitorService *StrategyMonitorService
	connectionManager *ConnectionManager
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		orderUpdateService: NewOrderUpdateService(hub),
		positionUpdateService: NewPositionUpdateService(hub),
		strategyMonitorService: NewStrategyMonitorService(hub),
		connectionManager: NewConnectionManager(hub),
	}
}

// RegisterRoutes registers WebSocket routes with the router
func (h *WebSocketHandler) RegisterRoutes(router *mux.Router) {
	// WebSocket endpoint
	router.HandleFunc("/ws", h.HandleWebSocket)
	
	// WebSocket status endpoint
	router.HandleFunc("/ws/status", h.HandleStatus)
}

// HandleWebSocket handles WebSocket connection requests
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthenticationMiddleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Serve WebSocket
	ServeWs(h.hub, w, r, userID)
}

// HandleStatus returns the status of WebSocket connections
func (h *WebSocketHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	// Get connection statistics
	status := struct {
		ActiveConnections int            `json:"activeConnections"`
		TopicSubscriptions map[string]int `json:"topicSubscriptions"`
		Timestamp         time.Time      `json:"timestamp"`
	}{
		ActiveConnections: h.connectionManager.GetActiveConnections(),
		TopicSubscriptions: make(map[string]int),
		Timestamp:         time.Now(),
	}
	
	// Get subscription counts for common topics
	commonTopics := []string{"orders", "positions", "strategies"}
	for _, topic := range commonTopics {
		status.TopicSubscriptions[topic] = h.connectionManager.GetTopicSubscriptions(topic)
	}
	
	// Return status as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// WebSocketServer represents the WebSocket server
type WebSocketServer struct {
	hub *Hub
	handler *WebSocketHandler
}

// NewWebSocketServer creates a new WebSocketServer
func NewWebSocketServer() *WebSocketServer {
	hub := NewHub()
	handler := NewWebSocketHandler(hub)
	
	return &WebSocketServer{
		hub: hub,
		handler: handler,
	}
}

// Start starts the WebSocket server
func (s *WebSocketServer) Start() {
	// Start the hub in a goroutine
	go s.hub.Run()
}

// GetHandler returns the WebSocket handler
func (s *WebSocketServer) GetHandler() *WebSocketHandler {
	return s.handler
}

// GetHub returns the WebSocket hub
func (s *WebSocketServer) GetHub() *Hub {
	return s.hub
}

// WebSocketIntegration provides integration with other services
type WebSocketIntegration struct {
	orderUpdateService *OrderUpdateService
	positionUpdateService *PositionUpdateService
	strategyMonitorService *StrategyMonitorService
}

// NewWebSocketIntegration creates a new WebSocketIntegration
func NewWebSocketIntegration(hub *Hub) *WebSocketIntegration {
	return &WebSocketIntegration{
		orderUpdateService: NewOrderUpdateService(hub),
		positionUpdateService: NewPositionUpdateService(hub),
		strategyMonitorService: NewStrategyMonitorService(hub),
	}
}

// NotifyOrderUpdate notifies clients of an order update
func (i *WebSocketIntegration) NotifyOrderUpdate(order *models.Order) {
	i.orderUpdateService.BroadcastOrderUpdate(order)
}

// NotifyPositionUpdate notifies clients of a position update
func (i *WebSocketIntegration) NotifyPositionUpdate(position *models.Position) {
	i.positionUpdateService.BroadcastPositionUpdate(position)
}

// NotifyStrategyUpdate notifies clients of a strategy update
func (i *WebSocketIntegration) NotifyStrategyUpdate(strategy *models.Strategy) {
	i.strategyMonitorService.BroadcastStrategyUpdate(strategy)
}
