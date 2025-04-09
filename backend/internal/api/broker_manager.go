// Package api provides a unified API layer for broker interactions
package api

import (
	"errors"
	"fmt"
	"sync"

	"github.com/trading-platform/backend/internal/broker/common"
	"github.com/trading-platform/backend/internal/broker/factory"
)

// BrokerManager manages broker clients and provides a unified API
type BrokerManager struct {
	clients     map[string]common.BrokerClient
	configs     map[string]*common.BrokerConfig
	activeUsers map[string]string // Maps userID to clientID
	mu          sync.RWMutex
}

// NewBrokerManager creates a new broker manager
func NewBrokerManager() *BrokerManager {
	return &BrokerManager{
		clients:     make(map[string]common.BrokerClient),
		configs:     make(map[string]*common.BrokerConfig),
		activeUsers: make(map[string]string),
	}
}

// RegisterBroker registers a broker configuration with the manager
func (m *BrokerManager) RegisterBroker(clientID string, config *common.BrokerConfig) error {
	if clientID == "" {
		return errors.New("client ID is required")
	}

	if config == nil {
		return errors.New("broker configuration is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Store the configuration
	m.configs[clientID] = config

	return nil
}

// GetBrokerClient gets or creates a broker client for the specified client ID
func (m *BrokerManager) GetBrokerClient(clientID string) (common.BrokerClient, error) {
	if clientID == "" {
		return nil, errors.New("client ID is required")
	}

	m.mu.RLock()
	client, exists := m.clients[clientID]
	m.mu.RUnlock()

	if exists {
		return client, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check again in case another goroutine created the client
	client, exists = m.clients[clientID]
	if exists {
		return client, nil
	}

	// Get the configuration
	config, exists := m.configs[clientID]
	if !exists {
		return nil, fmt.Errorf("no configuration found for client ID: %s", clientID)
	}

	// Create the client
	client, err := factory.NewBrokerClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create broker client: %w", err)
	}

	// Store the client
	m.clients[clientID] = client

	return client, nil
}

// Login authenticates a user with the specified broker
func (m *BrokerManager) Login(clientID string, credentials *common.Credentials) (*common.Session, error) {
	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	session, err := client.Login(credentials)
	if err != nil {
		return nil, err
	}

	// Store the active user
	m.mu.Lock()
	m.activeUsers[session.UserID] = clientID
	m.mu.Unlock()

	return session, nil
}

// Logout invalidates a user's session with the specified broker
func (m *BrokerManager) Logout(clientID string) error {
	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return err
	}

	err = client.Logout()
	if err != nil {
		return err
	}

	// Remove the active user
	m.mu.Lock()
	for userID, cID := range m.activeUsers {
		if cID == clientID {
			delete(m.activeUsers, userID)
			break
		}
	}
	m.mu.Unlock()

	return nil
}

// GetClientIDForUser gets the client ID for the specified user ID
func (m *BrokerManager) GetClientIDForUser(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}

	m.mu.RLock()
	clientID, exists := m.activeUsers[userID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("no active session found for user ID: %s", userID)
	}

	return clientID, nil
}

// PlaceOrder places an order for the specified user
func (m *BrokerManager) PlaceOrder(userID string, order *common.Order) (*common.OrderResponse, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.PlaceOrder(order)
}

// ModifyOrder modifies an order for the specified user
func (m *BrokerManager) ModifyOrder(userID string, order *common.ModifyOrder) (*common.OrderResponse, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.ModifyOrder(order)
}

// CancelOrder cancels an order for the specified user
func (m *BrokerManager) CancelOrder(userID string, orderID string) (*common.OrderResponse, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.CancelOrder(orderID, "")
}

// GetOrderBook gets the order book for the specified user
func (m *BrokerManager) GetOrderBook(userID string) (*common.OrderBook, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.GetOrderBook("")
}

// GetPositions gets the positions for the specified user
func (m *BrokerManager) GetPositions(userID string) ([]common.Position, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.GetPositions("")
}

// GetHoldings gets the holdings for the specified user
func (m *BrokerManager) GetHoldings(userID string) ([]common.Holding, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.GetHoldings("")
}

// GetQuote gets quotes for the specified symbols
func (m *BrokerManager) GetQuote(userID string, symbols []string) (map[string]common.Quote, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.GetQuote(symbols)
}

// SubscribeToQuotes subscribes to real-time quotes for the specified symbols
func (m *BrokerManager) SubscribeToQuotes(userID string, symbols []string) (chan common.Quote, error) {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return nil, err
	}

	return client.SubscribeToQuotes(symbols)
}

// UnsubscribeFromQuotes unsubscribes from real-time quotes for the specified symbols
func (m *BrokerManager) UnsubscribeFromQuotes(userID string, symbols []string) error {
	clientID, err := m.GetClientIDForUser(userID)
	if err != nil {
		return err
	}

	client, err := m.GetBrokerClient(clientID)
	if err != nil {
		return err
	}

	return client.UnsubscribeFromQuotes(symbols)
}

// For dealer-specific operations, we need to add methods that take a dealer client ID and a target client ID

// PlaceDealerOrder places an order on behalf of a client as a dealer
func (m *BrokerManager) PlaceDealerOrder(dealerUserID string, targetClientID string, order *common.Order) (*common.OrderResponse, error) {
	dealerClientID, err := m.GetClientIDForUser(dealerUserID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(dealerClientID)
	if err != nil {
		return nil, err
	}

	// Set the client ID in the order
	order.ClientID = targetClientID

	// Check if the client supports dealer operations
	if xtsClient, ok := client.(interface {
		PlaceDealerOrder(order *common.Order) (*common.OrderResponse, error)
	}); ok {
		return xtsClient.PlaceDealerOrder(order)
	}

	// Fall back to regular order placement
	return client.PlaceOrder(order)
}

// GetDealerOrderBook gets the dealer order book for the specified client
func (m *BrokerManager) GetDealerOrderBook(dealerUserID string, targetClientID string) (*common.OrderBook, error) {
	dealerClientID, err := m.GetClientIDForUser(dealerUserID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(dealerClientID)
	if err != nil {
		return nil, err
	}

	// Check if the client supports dealer operations
	if xtsClient, ok := client.(interface {
		GetDealerOrderBook(clientID string) (*common.OrderBook, error)
	}); ok {
		return xtsClient.GetDealerOrderBook(targetClientID)
	}

	return nil, errors.New("dealer operations not supported by this broker")
}

// GetDealerPositions gets the dealer positions for the specified client
func (m *BrokerManager) GetDealerPositions(dealerUserID string, targetClientID string) ([]common.Position, error) {
	dealerClientID, err := m.GetClientIDForUser(dealerUserID)
	if err != nil {
		return nil, err
	}

	client, err := m.GetBrokerClient(dealerClientID)
	if err != nil {
		return nil, err
	}

	// Check if the client supports dealer operations
	if xtsClient, ok := client.(interface {
		GetDealerPositions(clientID string) ([]common.Position, error)
	}); ok {
		return xtsClient.GetDealerPositions(targetClientID)
	}

	return nil, errors.New("dealer operations not supported by this broker")
}
