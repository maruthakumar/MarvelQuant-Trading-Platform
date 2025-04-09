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

// MarketDataService provides market data functionality
type MarketDataService struct {
	restClient   *rest.Client
	wsClient     *websocket.MarketDataClient
	wsURL        string
	mutex        sync.RWMutex
	subscribers  map[string][]chan *models.Quote
	isRunning    bool
	stopChan     chan struct{}
	config       *config.XTSConfig
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(restClient *rest.Client, config *config.XTSConfig) *MarketDataService {
	return &MarketDataService{
		restClient:  restClient,
		wsURL:       fmt.Sprintf("%s/apimarketdata/socket.io", config.BaseURL),
		subscribers: make(map[string][]chan *models.Quote),
		stopChan:    make(chan struct{}),
		config:      config,
	}
}

// Start initializes and starts the market data service
func (s *MarketDataService) Start(token, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return nil
	}

	// Create WebSocket client
	s.wsClient = websocket.NewMarketDataClient(s.wsURL, token, userID)

	// Connect to WebSocket
	if err := s.wsClient.Connect(); err != nil {
		return err
	}

	// Start message processor
	go s.processMessages()

	s.isRunning = true
	return nil
}

// Stop stops the market data service
func (s *MarketDataService) Stop() error {
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

// GetQuotes gets quotes for instruments
func (s *MarketDataService) GetQuotes(instruments []string) ([]*models.Quote, error) {
	// Implementation for REST API call to get quotes
	// This would typically call the restClient to fetch quotes
	return nil, fmt.Errorf("not implemented")
}

// SubscribeToQuotes subscribes to real-time quotes for instruments
func (s *MarketDataService) SubscribeToQuotes(instruments []models.Instrument) (chan *models.Quote, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil, errors.ErrWebSocketClosed
	}

	// Create a channel for this subscriber
	quoteChan := make(chan *models.Quote, 100)

	// Register subscriber for each instrument
	for _, instrument := range instruments {
		key := fmt.Sprintf("%s:%s", instrument.ExchangeSegment, instrument.ExchangeInstrumentID)
		s.subscribers[key] = append(s.subscribers[key], quoteChan)
	}

	// Subscribe to instruments via WebSocket
	if err := s.wsClient.Subscribe(instruments); err != nil {
		return nil, err
	}

	return quoteChan, nil
}

// UnsubscribeFromQuotes unsubscribes from real-time quotes
func (s *MarketDataService) UnsubscribeFromQuotes(instruments []models.Instrument, quoteChan chan *models.Quote) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return errors.ErrWebSocketClosed
	}

	// Unregister subscriber for each instrument
	for _, instrument := range instruments {
		key := fmt.Sprintf("%s:%s", instrument.ExchangeSegment, instrument.ExchangeInstrumentID)
		
		// Find and remove the channel from subscribers
		subscribers := s.subscribers[key]
		for i, ch := range subscribers {
			if ch == quoteChan {
				// Remove this channel
				s.subscribers[key] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}

		// If no more subscribers for this instrument, unsubscribe from WebSocket
		if len(s.subscribers[key]) == 0 {
			delete(s.subscribers, key)
		}
	}

	// Unsubscribe from instruments via WebSocket if no more subscribers
	return s.wsClient.Unsubscribe(instruments)
}

// GetHistoricalData gets historical OHLC data for an instrument
func (s *MarketDataService) GetHistoricalData(
	exchangeSegment, exchangeInstrumentID string,
	startTime, endTime time.Time,
	interval string) ([]*models.OHLC, error) {
	
	// Implementation for REST API call to get historical data
	// This would typically call the restClient to fetch historical data
	return nil, fmt.Errorf("not implemented")
}

// SearchInstruments searches for instruments by string
func (s *MarketDataService) SearchInstruments(searchString string) ([]models.Instrument, error) {
	// Implementation for REST API call to search instruments
	// This would typically call the restClient to search instruments
	return nil, fmt.Errorf("not implemented")
}

// processMessages processes messages from the WebSocket
func (s *MarketDataService) processMessages() {
	messageChan := s.wsClient.GetMessageChannel()
	errorChan := s.wsClient.GetErrorChannel()

	for {
		select {
		case <-s.stopChan:
			return
		case msg := <-messageChan:
			// Process market data message
			s.handleMarketDataMessage(msg)
		case err := <-errorChan:
			// Log error
			log.Printf("Market data WebSocket error: %v", err)
		}
	}
}

// handleMarketDataMessage handles a market data message
func (s *MarketDataService) handleMarketDataMessage(msg []byte) {
	// Parse message to extract quote data
	quote, err := parseQuoteFromMessage(msg)
	if err != nil {
		log.Printf("Error parsing quote: %v", err)
		return
	}

	if quote == nil {
		return
	}

	// Distribute quote to subscribers
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", quote.ExchangeSegment, quote.ExchangeInstrumentID)
	subscribers := s.subscribers[key]

	for _, ch := range subscribers {
		select {
		case ch <- quote:
			// Quote sent successfully
		default:
			// Channel is full, log warning
			log.Printf("Warning: Quote channel is full, dropping quote for %s", key)
		}
	}
}

// parseQuoteFromMessage parses a quote from a WebSocket message
func parseQuoteFromMessage(msg []byte) (*models.Quote, error) {
	// Implementation to parse the WebSocket message into a Quote
	// This would typically involve JSON unmarshaling and data transformation
	return nil, fmt.Errorf("not implemented")
}

// WithContext returns a new MarketDataService with context support
func (s *MarketDataService) WithContext(ctx context.Context) *MarketDataService {
	// Create a copy of the service
	newService := &MarketDataService{
		restClient:  s.restClient,
		wsURL:       s.wsURL,
		subscribers: make(map[string][]chan *models.Quote),
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
