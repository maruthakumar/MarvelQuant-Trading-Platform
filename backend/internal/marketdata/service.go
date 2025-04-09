package marketdata

import (
	"context"
	"sync"
	"time"
)

// MarketData represents real-time market data
type MarketData struct {
	Symbol     string    `json:"symbol"`
	Exchange   string    `json:"exchange"`
	LastPrice  float64   `json:"lastPrice"`
	BidPrice   float64   `json:"bidPrice"`
	AskPrice   float64   `json:"askPrice"`
	BidSize    int       `json:"bidSize"`
	AskSize    int       `json:"askSize"`
	Volume     int       `json:"volume"`
	OpenPrice  float64   `json:"openPrice"`
	HighPrice  float64   `json:"highPrice"`
	LowPrice   float64   `json:"lowPrice"`
	ClosePrice float64   `json:"closePrice"`
	Timestamp  time.Time `json:"timestamp"`
}

// OHLCV represents Open, High, Low, Close, Volume data
type OHLCV struct {
	Symbol    string    `json:"symbol"`
	Interval  string    `json:"interval"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int       `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// IndicatorValue represents a calculated indicator value
type IndicatorValue struct {
	Symbol    string                 `json:"symbol"`
	Indicator string                 `json:"indicator"`
	Value     float64                `json:"value"`
	Values    map[string]float64     `json:"values,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MarketDataCallback is a function that is called when market data is received
type MarketDataCallback func(data MarketData)

// DataSourceConnector interface for connecting to data sources
type DataSourceConnector interface {
	Connect(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
	GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error)
	GetHistoricalData(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error)
	SubscribeToMarketData(ctx context.Context, symbols []string, callback MarketDataCallback) error
	UnsubscribeFromMarketData(ctx context.Context, symbols []string) error
}

// DataSourceManager manages multiple data sources
type DataSourceManager struct {
	primarySource   DataSourceConnector
	backupSources   []DataSourceConnector
	activeSource    DataSourceConnector
	mutex           sync.RWMutex
	subscriptions   map[string][]MarketDataCallback
	subscriptionsMu sync.RWMutex
}

// NewDataSourceManager creates a new data source manager
func NewDataSourceManager(primarySource DataSourceConnector, backupSources ...DataSourceConnector) *DataSourceManager {
	return &DataSourceManager{
		primarySource: primarySource,
		backupSources: backupSources,
		activeSource:  primarySource,
		subscriptions: make(map[string][]MarketDataCallback),
	}
}

// Connect connects to the data sources
func (m *DataSourceManager) Connect(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Try to connect to primary source
	err := m.primarySource.Connect(ctx)
	if err == nil {
		m.activeSource = m.primarySource
		return nil
	}

	// If primary source fails, try backup sources
	for _, source := range m.backupSources {
		err = source.Connect(ctx)
		if err == nil {
			m.activeSource = source
			return nil
		}
	}

	return err
}

// Disconnect disconnects from the data sources
func (m *DataSourceManager) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Disconnect from primary source
	if err := m.primarySource.Disconnect(); err != nil {
		return err
	}

	// Disconnect from backup sources
	for _, source := range m.backupSources {
		if err := source.Disconnect(); err != nil {
			return err
		}
	}

	return nil
}

// GetMarketData gets market data for the specified symbols
func (m *DataSourceManager) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	m.mutex.RLock()
	activeSource := m.activeSource
	m.mutex.RUnlock()

	// Try to get market data from active source
	data, err := activeSource.GetMarketData(ctx, symbols)
	if err == nil {
		return data, nil
	}

	// If active source fails, try other sources
	if activeSource == m.primarySource {
		for _, source := range m.backupSources {
			data, err = source.GetMarketData(ctx, symbols)
			if err == nil {
				m.mutex.Lock()
				m.activeSource = source
				m.mutex.Unlock()
				return data, nil
			}
		}
	} else {
		// Try primary source
		data, err = m.primarySource.GetMarketData(ctx, symbols)
		if err == nil {
			m.mutex.Lock()
			m.activeSource = m.primarySource
			m.mutex.Unlock()
			return data, nil
		}

		// Try other backup sources
		for _, source := range m.backupSources {
			if source == activeSource {
				continue
			}
			data, err = source.GetMarketData(ctx, symbols)
			if err == nil {
				m.mutex.Lock()
				m.activeSource = source
				m.mutex.Unlock()
				return data, nil
			}
		}
	}

	return nil, err
}

// GetHistoricalData gets historical data for the specified symbol and interval
func (m *DataSourceManager) GetHistoricalData(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error) {
	m.mutex.RLock()
	activeSource := m.activeSource
	m.mutex.RUnlock()

	// Try to get historical data from active source
	data, err := activeSource.GetHistoricalData(ctx, symbol, interval, from, to)
	if err == nil {
		return data, nil
	}

	// If active source fails, try other sources
	if activeSource == m.primarySource {
		for _, source := range m.backupSources {
			data, err = source.GetHistoricalData(ctx, symbol, interval, from, to)
			if err == nil {
				m.mutex.Lock()
				m.activeSource = source
				m.mutex.Unlock()
				return data, nil
			}
		}
	} else {
		// Try primary source
		data, err = m.primarySource.GetHistoricalData(ctx, symbol, interval, from, to)
		if err == nil {
			m.mutex.Lock()
			m.activeSource = m.primarySource
			m.mutex.Unlock()
			return data, nil
		}

		// Try other backup sources
		for _, source := range m.backupSources {
			if source == activeSource {
				continue
			}
			data, err = source.GetHistoricalData(ctx, symbol, interval, from, to)
			if err == nil {
				m.mutex.Lock()
				m.activeSource = source
				m.mutex.Unlock()
				return data, nil
			}
		}
	}

	return nil, err
}

// SubscribeToMarketData subscribes to market data for the specified symbols
func (m *DataSourceManager) SubscribeToMarketData(ctx context.Context, symbols []string, callback MarketDataCallback) error {
	m.mutex.RLock()
	activeSource := m.activeSource
	m.mutex.RUnlock()

	// Register callback
	m.subscriptionsMu.Lock()
	for _, symbol := range symbols {
		m.subscriptions[symbol] = append(m.subscriptions[symbol], callback)
	}
	m.subscriptionsMu.Unlock()

	// Create a wrapper callback that distributes the data to all registered callbacks
	wrapperCallback := func(data MarketData) {
		m.subscriptionsMu.RLock()
		callbacks := m.subscriptions[data.Symbol]
		m.subscriptionsMu.RUnlock()

		for _, cb := range callbacks {
			cb(data)
		}
	}

	// Subscribe to market data from active source
	return activeSource.SubscribeToMarketData(ctx, symbols, wrapperCallback)
}

// UnsubscribeFromMarketData unsubscribes from market data for the specified symbols
func (m *DataSourceManager) UnsubscribeFromMarketData(ctx context.Context, symbols []string) error {
	m.mutex.RLock()
	activeSource := m.activeSource
	m.mutex.RUnlock()

	// Unregister callbacks
	m.subscriptionsMu.Lock()
	for _, symbol := range symbols {
		delete(m.subscriptions, symbol)
	}
	m.subscriptionsMu.Unlock()

	// Unsubscribe from market data from active source
	return activeSource.UnsubscribeFromMarketData(ctx, symbols)
}

// DataProcessor interface for processing data
type DataProcessor interface {
	Process(data interface{}) (interface{}, error)
}

// DataNormalizer normalizes data from different sources
type DataNormalizer struct {
	sourceFormat string
	targetFormat string
}

// NewDataNormalizer creates a new data normalizer
func NewDataNormalizer(sourceFormat, targetFormat string) *DataNormalizer {
	return &DataNormalizer{
		sourceFormat: sourceFormat,
		targetFormat: targetFormat,
	}
}

// Process normalizes the data
func (n *DataNormalizer) Process(data interface{}) (interface{}, error) {
	// Implementation depends on the specific source and target formats
	// For now, we'll just return the data as is
	return data, nil
}

// DataEnricher enriches data with additional information
type DataEnricher struct {
	enrichmentRules []EnrichmentRule
}

// EnrichmentRule defines a rule for enriching data
type EnrichmentRule interface {
	Apply(data interface{}) (interface{}, error)
}

// NewDataEnricher creates a new data enricher
func NewDataEnricher(rules ...EnrichmentRule) *DataEnricher {
	return &DataEnricher{
		enrichmentRules: rules,
	}
}

// Process enriches the data
func (e *DataEnricher) Process(data interface{}) (interface{}, error) {
	result := data

	// Apply each enrichment rule
	for _, rule := range e.enrichmentRules {
		var err error
		result, err = rule.Apply(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// DataAggregator aggregates data into different time frames
type DataAggregator struct {
	timeFrame string
}

// NewDataAggregator creates a new data aggregator
func NewDataAggregator(timeFrame string) *DataAggregator {
	return &DataAggregator{
		timeFrame: timeFrame,
	}
}

// Process aggregates the data
func (a *DataAggregator) Process(data interface{}) (interface{}, error) {
	// Implementation depends on the specific time frame and data format
	// For now, we'll just return the data as is
	return data, nil
}

// MarketDataService is the main service for market data
type MarketDataService struct {
	dataSourceManager *DataSourceManager
	dataStorage       DataStorage
	cacheManager      CacheManager
	processors        []DataProcessor
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(
	dataSourceManager *DataSourceManager,
	dataStorage DataStorage,
	cacheManager CacheManager,
	processors ...DataProcessor,
) *MarketDataService {
	return &MarketDataService{
		dataSourceManager: dataSourceManager,
		dataStorage:       dataStorage,
		cacheManager:      cacheManager,
		processors:        processors,
	}
}

// Start starts the market data service
func (s *MarketDataService) Start(ctx context.Context) error {
	// Connect to data sources
	if err := s.dataSourceManager.Connect(ctx); err != nil {
		return err
	}

	// TODO: Start any background processes

	return nil
}

// Stop stops the market data service
func (s *MarketDataService) Stop() error {
	// Disconnect from data sources
	if err := s.dataSourceManager.Disconnect(); err != nil {
		return err
	}

	// TODO: Stop any background processes

	return nil
}

// GetMarketData gets market data for the specified symbols
func (s *MarketDataService) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	// Check cache first
	cachedData := make(map[string]MarketData)
	missingSymbols := make([]string, 0, len(symbols))

	for _, symbol := range symbols {
		key := "market_data:" + symbol
		data, found := s.cacheManager.Get(key)
		if found {
			if md, ok := data.(MarketData); ok {
				cachedData[symbol] = md
				continue
			}
		}
		missingSymbols = append(missingSymbols, symbol)
	}

	// If all symbols were in cache, return cached data
	if len(missingSymbols) == 0 {
		return cachedData, nil
	}

	// Get data for missing symbols from data source
	data, err := s.dataSourceManager.GetMarketData(ctx, missingSymbols)
	if err != nil {
		return nil, err
	}

	// Process the data
	for symbol, md := range data {
		processedData := md

		// Apply each processor
		for _, processor := range s.processors {
			result, err := processor.Process(processedData)
			if err != nil {
				return nil, err
			}
			if processed, ok := result.(MarketData); ok {
				processedData = processed
			}
		}

		// Store processed data in cache
		key := "market_data:" + symbol
		s.cacheManager.Set(key, processedData, 5*time.Second)

		// Store processed data in storage
		if err := s.dataStorage.StoreMarketData(ctx, processedData); err != nil {
			return nil, err
		}

		// Add to result
		data[symbol] = processedData
	}

	// Merge cached data and new data
	for symbol, md := range cachedData {
		data[symbol] = md
	}

	return data, nil
}

// GetHistoricalData gets historical data for the specified symbol and interval
func (s *MarketDataService) GetHistoricalData(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error) {
	// Check cache first
	key := "ohlcv:" + symbol + ":" + interval + ":" + from.Format(time.RFC3339) + ":" + to.Format(time.RFC3339)
	data, found := s.cacheManager.Get(key)
	if found {
		if ohlcv, ok := data.([]OHLCV); ok {
			return ohlcv, nil
		}
	}

	// Get data from storage
	ohlcv, err := s.dataStorage.GetOHLCV(ctx, symbol, interval, from, to)
	if err == nil && len(ohlcv) > 0 {
		// Store in cache
		s.cacheManager.Set(key, ohlcv, 5*time.Minute)
		return ohlcv, nil
	}

	// If not in storage or error, get from data source
	ohlcv, err = s.dataSourceManager.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}

	// Process the data
	processedData := ohlcv

	// Apply each processor
	for _, processor := range s.processors {
		result, err := processor.Process(processedData)
		if err != nil {
			return nil, err
		}
		if processed, ok := result.([]OHLCV); ok {
			processedData = processed
		}
	}

	// Store processed data in cache
	s.cacheManager.Set(key, processedData, 5*time.Minute)

	// Store processed data in storage
	if err := s.dataStorage.StoreOHLCV(ctx, symbol, interval, processedData); err != nil {
		return nil, err
	}

	return processedData, nil
}

// SubscribeToMarketData subscribes to market data for the specified symbols
func (s *MarketDataService) SubscribeToMarketData(ctx context.Context, symbols []string, callback MarketDataCallback) error {
	// Create a wrapper callback that processes the data before calling the user callback
	wrapperCallback := func(data MarketData) {
		processedData := data

		// Apply each processor
		for _, processor := range s.processors {
			result, err := processor.Process(processedData)
			if err != nil {
				// Log error and continue with unprocessed data
				continue
			}
			if processed, ok := result.(MarketData); ok {
				processedData = processed
			}
		}

		// Store processed data in cache
		key := "market_data:" + data.Symbol
		s.cacheManager.Set(key, processedData, 5*time.Second)

		// Store processed data in storage (async)
		go func() {
			if err := s.dataStorage.StoreMarketData(context.Background(), processedData); err != nil {
				// Log error
			}
		}()

		// Call user callback
		callback(processedData)
	}

	return s.dataSourceManager.SubscribeToMarketData(ctx, symbols, wrapperCallback)
}

// UnsubscribeFromMarketData unsubscribes from market data for the specified symbols
func (s *MarketDataService) UnsubscribeFromMarketData(ctx context.Context, symbols []string) error {
	return s.dataSourceManager.UnsubscribeFromMarketData(ctx, symbols)
}

// DataStorage interface for storing and retrieving data
type DataStorage interface {
	StoreMarketData(ctx context.Context, data MarketData) error
	StoreOHLCV(ctx context.Context, symbol string, interval string, data []OHLCV) error
	GetMarketData(ctx context.Context, symbol string) (MarketData, error)
	GetOHLCV(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error)
	GetLatestOHLCV(ctx context.Context, symbol string, interval string, limit int) ([]OHLCV, error)
}

// CacheManager interface for managing cache
type CacheManager interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	Clear() error
}
