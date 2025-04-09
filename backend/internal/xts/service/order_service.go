package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/trade-execution-platform/backend/internal/xts/config"
	"github.com/trade-execution-platform/backend/internal/xts/errors"
	"github.com/trade-execution-platform/backend/internal/xts/models"
	"github.com/trade-execution-platform/backend/internal/xts/rest"
	"github.com/trade-execution-platform/backend/internal/xts/websocket"
)

// OrderService provides order management functionality
type OrderService struct {
	restClient   *rest.Client
	wsClient     *websocket.OrderClient
	wsURL        string
	mutex        sync.RWMutex
	subscribers  map[string][]chan interface{}
	isRunning    bool
	stopChan     chan struct{}
	config       *config.XTSConfig
}

// NewOrderService creates a new order service
func NewOrderService(restClient *rest.Client, config *config.XTSConfig) *OrderService {
	return &OrderService{
		restClient:  restClient,
		wsURL:       fmt.Sprintf("%s/interactive/socket.io", config.BaseURL),
		subscribers: make(map[string][]chan interface{}),
		stopChan:    make(chan struct{}),
		config:      config,
	}
}

// Start initializes and starts the order service
func (s *OrderService) Start(token, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return nil
	}

	// Create WebSocket client
	s.wsClient = websocket.NewOrderClient(s.wsURL, token, userID)

	// Connect to WebSocket
	if err := s.wsClient.Connect(); err != nil {
		return err
	}

	// Start message processor
	go s.processMessages()

	s.isRunning = true
	return nil
}

// Stop stops the order service
func (s *OrderService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil
	}

	// Signal to stop message processor
	close(s.stopChan)

	// Disconnect WebSocket
	if s.wsClient != nil {
		if err := s.wsClient.Disconnect(); err != nil {
			return err
		}
	}

	s.isRunning = false
	return nil
}

// PlaceOrder places an order
func (s *OrderService) PlaceOrder(order *models.Order) (*models.OrderResponse, error) {
	if !s.isRunning {
		return nil, errors.ErrSessionInvalid
	}

	return s.restClient.PlaceOrder(order)
}

// ModifyOrder modifies an existing order
func (s *OrderService) ModifyOrder(modifyOrder *models.ModifyOrder) (*models.OrderResponse, error) {
	if !s.isRunning {
		return nil, errors.ErrSessionInvalid
	}

	return s.restClient.ModifyOrder(modifyOrder)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(appOrderID int, clientID string) (*models.OrderResponse, error) {
	if !s.isRunning {
		return nil, errors.ErrSessionInvalid
	}

	return s.restClient.CancelOrder(appOrderID, clientID)
}

// GetOrderBook gets the order book
func (s *OrderService) GetOrderBook(clientID string) (*models.OrderBook, error) {
	if !s.isRunning {
		return nil, errors.ErrSessionInvalid
	}

	return s.restClient.GetOrderBook(clientID)
}

// GetPositions gets current positions
func (s *OrderService) GetPositions(clientID string) ([]models.Position, error) {
	if !s.isRunning {
		return nil, errors.ErrSessionInvalid
	}

	return s.restClient.GetPositions(clientID)
}

// SubscribeToOrderUpdates subscribes to real-time order updates
func (s *OrderService) SubscribeToOrderUpdates() (chan interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil, errors.ErrWebSocketClosed
	}

	// Create a channel for this subscriber
	updateChan := make(chan interface{}, 100)

	// Register subscriber for order updates
	s.subscribers["order"] = append(s.subscribers["order"], updateChan)

	return updateChan, nil
}

// SubscribeToTradeUpdates subscribes to real-time trade updates
func (s *OrderService) SubscribeToTradeUpdates() (chan interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil, errors.ErrWebSocketClosed
	}

	// Create a channel for this subscriber
	updateChan := make(chan interface{}, 100)

	// Register subscriber for trade updates
	s.subscribers["trade"] = append(s.subscribers["trade"], updateChan)

	return updateChan, nil
}

// SubscribeToPositionUpdates subscribes to real-time position updates
func (s *OrderService) SubscribeToPositionUpdates() (chan interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil, errors.ErrWebSocketClosed
	}

	// Create a channel for this subscriber
	updateChan := make(chan interface{}, 100)

	// Register subscriber for position updates
	s.subscribers["position"] = append(s.subscribers["position"], updateChan)

	return updateChan, nil
}

// Unsubscribe unsubscribes from updates
func (s *OrderService) Unsubscribe(updateType string, updateChan chan interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return errors.ErrWebSocketClosed
	}

	// Find and remove the channel from subscribers
	subscribers := s.subscribers[updateType]
	for i, ch := range subscribers {
		if ch == updateChan {
			// Remove this channel
			s.subscribers[updateType] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}

	return nil
}

// processMessages processes messages from the WebSocket
func (s *OrderService) processMessages() {
	messageChan := s.wsClient.GetMessageChannel()
	errorChan := s.wsClient.GetErrorChannel()

	for {
		select {
		case <-s.stopChan:
			return
		case msg := <-messageChan:
			// Process order message
			s.handleOrderMessage(msg)
		case err := <-errorChan:
			// Log error
			log.Printf("Order WebSocket error: %v", err)
		}
	}
}

// handleOrderMessage handles an order message
func (s *OrderService) handleOrderMessage(msg []byte) {
	// Parse message to determine type and extract data
	messageType, data, err := parseOrderMessage(msg)
	if err != nil {
		log.Printf("Error parsing order message: %v", err)
		return
	}

	if data == nil {
		return
	}

	// Distribute message to subscribers
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	subscribers := s.subscribers[messageType]
	for _, ch := range subscribers {
		select {
		case ch <- data:
			// Message sent successfully
		default:
			// Channel is full, log warning
			log.Printf("Warning: Update channel is full, dropping %s update", messageType)
		}
	}
}

// parseOrderMessage parses an order message from a WebSocket message
func parseOrderMessage(msg []byte) (string, interface{}, error) {
	// Implementation to parse the WebSocket message
	// This would typically involve JSON unmarshaling and data transformation
	return "", nil, fmt.Errorf("not implemented")
}

// WithContext returns a new OrderService with context support
func (s *OrderService) WithContext(ctx context.Context) *OrderService {
	// Create a copy of the service
	newService := &OrderService{
		restClient:  s.restClient,
		wsURL:       s.wsURL,
		subscribers: make(map[string][]chan interface{}),
		stopChan:    make(chan struct{}),
		config:      s.config,
	}

	// Monitor context cancellation
	go func() {
		<-ctx.Done()
		newService.Stop()
	}()

	return newService
}
