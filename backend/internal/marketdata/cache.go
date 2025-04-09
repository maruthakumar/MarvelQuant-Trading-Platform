package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// CacheConfig defines configuration for the cache
type CacheConfig struct {
	DefaultTTL        time.Duration
	CleanupInterval   time.Duration
	MaxSize           int
	EnableCompression bool
}

// DefaultCacheConfig returns the default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL:        5 * time.Minute,
		CleanupInterval:   10 * time.Minute,
		MaxSize:           10000,
		EnableCompression: true,
	}
}

// CacheManager manages caching of market data
type CacheManager struct {
	config         CacheConfig
	cache          map[string]*CacheEntry
	mutex          sync.RWMutex
	stopCleanup    chan struct{}
	compressionMgr *CompressionManager
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Key       string
	Value     []byte
	Timestamp time.Time
	Expiry    time.Time
	Size      int
	Compressed bool
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config CacheConfig) *CacheManager {
	cm := &CacheManager{
		config:         config,
		cache:          make(map[string]*CacheEntry),
		stopCleanup:    make(chan struct{}),
		compressionMgr: NewCompressionManager(),
	}

	// Start cleanup goroutine
	go cm.startCleanup()

	return cm
}

// Get gets a value from the cache
func (cm *CacheManager) Get(key string, result interface{}) bool {
	cm.mutex.RLock()
	entry, found := cm.cache[key]
	cm.mutex.RUnlock()

	if !found {
		return false
	}

	// Check if entry has expired
	if time.Now().After(entry.Expiry) {
		cm.mutex.Lock()
		delete(cm.cache, key)
		cm.mutex.Unlock()
		return false
	}

	// Decompress if needed
	var data []byte
	var err error
	if entry.Compressed {
		data, err = cm.compressionMgr.Decompress(entry.Value)
		if err != nil {
			log.Printf("Error decompressing cache entry: %v", err)
			return false
		}
	} else {
		data = entry.Value
	}

	// Unmarshal data
	if err := json.Unmarshal(data, result); err != nil {
		log.Printf("Error unmarshaling cache entry: %v", err)
		return false
	}

	return true
}

// Set sets a value in the cache
func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	// Compress if enabled and size is large enough
	compressed := false
	if cm.config.EnableCompression && len(data) > 1024 {
		compressedData, err := cm.compressionMgr.Compress(data)
		if err != nil {
			log.Printf("Error compressing data: %v", err)
		} else if len(compressedData) < len(data) {
			data = compressedData
			compressed = true
		}
	}

	// Set TTL
	if ttl == 0 {
		ttl = cm.config.DefaultTTL
	}

	// Create cache entry
	entry := &CacheEntry{
		Key:       key,
		Value:     data,
		Timestamp: time.Now(),
		Expiry:    time.Now().Add(ttl),
		Size:      len(data),
		Compressed: compressed,
	}

	// Add to cache
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if cache is full
	if len(cm.cache) >= cm.config.MaxSize {
		cm.evictOldest()
	}

	cm.cache[key] = entry
	return nil
}

// Delete deletes a value from the cache
func (cm *CacheManager) Delete(key string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.cache, key)
}

// Clear clears the cache
func (cm *CacheManager) Clear() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cache = make(map[string]*CacheEntry)
}

// Size returns the number of items in the cache
func (cm *CacheManager) Size() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.cache)
}

// MemoryUsage returns the approximate memory usage of the cache in bytes
func (cm *CacheManager) MemoryUsage() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var total int
	for _, entry := range cm.cache {
		total += entry.Size
	}
	return total
}

// Stop stops the cache manager
func (cm *CacheManager) Stop() {
	close(cm.stopCleanup)
}

// startCleanup starts the cleanup goroutine
func (cm *CacheManager) startCleanup() {
	ticker := time.NewTicker(cm.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.cleanup()
		case <-cm.stopCleanup:
			return
		}
	}
}

// cleanup removes expired entries
func (cm *CacheManager) cleanup() {
	now := time.Now()
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for key, entry := range cm.cache {
		if now.After(entry.Expiry) {
			delete(cm.cache, key)
		}
	}
}

// evictOldest evicts the oldest entry from the cache
func (cm *CacheManager) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	// Find oldest entry
	for key, entry := range cm.cache {
		if oldestTime.IsZero() || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
		}
	}

	// Delete oldest entry
	if oldestKey != "" {
		delete(cm.cache, oldestKey)
	}
}

// CompressionManager manages compression of data
type CompressionManager struct {
	// Add compression-specific fields if needed
}

// NewCompressionManager creates a new compression manager
func NewCompressionManager() *CompressionManager {
	return &CompressionManager{}
}

// Compress compresses data
func (cm *CompressionManager) Compress(data []byte) ([]byte, error) {
	// In a real implementation, this would use a compression algorithm like gzip
	// For simplicity, we'll just return the original data
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress decompresses data
func (cm *CompressionManager) Decompress(data []byte) ([]byte, error) {
	// In a real implementation, this would use a decompression algorithm like gzip
	// For simplicity, we'll just return the original data
	buf := bytes.NewBuffer(data)
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	
	result, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MarketDataCache is a specialized cache for market data
type MarketDataCache struct {
	cacheManager *CacheManager
}

// NewMarketDataCache creates a new market data cache
func NewMarketDataCache(config CacheConfig) *MarketDataCache {
	return &MarketDataCache{
		cacheManager: NewCacheManager(config),
	}
}

// GetMarketData gets market data from the cache
func (c *MarketDataCache) GetMarketData(symbol string) (MarketData, bool) {
	var data MarketData
	key := fmt.Sprintf("market_data:%s", symbol)
	if c.cacheManager.Get(key, &data) {
		return data, true
	}
	return MarketData{}, false
}

// SetMarketData sets market data in the cache
func (c *MarketDataCache) SetMarketData(data MarketData, ttl time.Duration) error {
	key := fmt.Sprintf("market_data:%s", data.Symbol)
	return c.cacheManager.Set(key, data, ttl)
}

// GetHistoricalData gets historical data from the cache
func (c *MarketDataCache) GetHistoricalData(symbol, interval string, from, to time.Time) ([]OHLCV, bool) {
	var data []OHLCV
	key := fmt.Sprintf("historical_data:%s:%s:%d:%d", symbol, interval, from.Unix(), to.Unix())
	if c.cacheManager.Get(key, &data) {
		return data, true
	}
	return nil, false
}

// SetHistoricalData sets historical data in the cache
func (c *MarketDataCache) SetHistoricalData(symbol, interval string, from, to time.Time, data []OHLCV, ttl time.Duration) error {
	key := fmt.Sprintf("historical_data:%s:%s:%d:%d", symbol, interval, from.Unix(), to.Unix())
	return c.cacheManager.Set(key, data, ttl)
}

// GetIndicatorData gets indicator data from the cache
func (c *MarketDataCache) GetIndicatorData(indicator, symbol, interval string, from, to time.Time, params map[string]interface{}) ([]map[string]interface{}, bool) {
	var data []map[string]interface{}
	
	// Create a deterministic key from params
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, false
	}
	
	key := fmt.Sprintf("indicator_data:%s:%s:%s:%d:%d:%s", indicator, symbol, interval, from.Unix(), to.Unix(), string(paramsJSON))
	if c.cacheManager.Get(key, &data) {
		return data, true
	}
	return nil, false
}

// SetIndicatorData sets indicator data in the cache
func (c *MarketDataCache) SetIndicatorData(indicator, symbol, interval string, from, to time.Time, params map[string]interface{}, data []map[string]interface{}, ttl time.Duration) error {
	// Create a deterministic key from params
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}
	
	key := fmt.Sprintf("indicator_data:%s:%s:%s:%d:%d:%s", indicator, symbol, interval, from.Unix(), to.Unix(), string(paramsJSON))
	return c.cacheManager.Set(key, data, ttl)
}

// CachedMarketDataService wraps a market data service with caching
type CachedMarketDataService struct {
	service *MarketDataService
	cache   *MarketDataCache
}

// NewCachedMarketDataService creates a new cached market data service
func NewCachedMarketDataService(service *MarketDataService, cacheConfig CacheConfig) *CachedMarketDataService {
	return &CachedMarketDataService{
		service: service,
		cache:   NewMarketDataCache(cacheConfig),
	}
}

// GetMarketData gets market data with caching
func (s *CachedMarketDataService) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	result := make(map[string]MarketData)
	var uncachedSymbols []string

	// Check cache first
	for _, symbol := range symbols {
		if data, found := s.cache.GetMarketData(symbol); found {
			result[symbol] = data
		} else {
			uncachedSymbols = append(uncachedSymbols, symbol)
		}
	}

	// If all symbols were in cache, return immediately
	if len(uncachedSymbols) == 0 {
		return result, nil
	}

	// Get uncached symbols from service
	uncachedData, err := s.service.GetMarketData(ctx, uncachedSymbols)
	if err != nil {
		return result, err
	}

	// Cache results and add to result
	for symbol, data := range uncachedData {
		s.cache.SetMarketData(data, 5*time.Second) // Short TTL for market data
		result[symbol] = data
	}

	return result, nil
}

// GetHistoricalData gets historical data with caching
func (s *CachedMarketDataService) GetHistoricalData(ctx context.Context, symbol, interval string, from, to time.Time) ([]OHLCV, error) {
	// Check cache first
	if data, found := s.cache.GetHistoricalData(symbol, interval, from, to); found {
		return data, nil
	}

	// Get from service
	data, err := s.service.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}

	// Cache results
	s.cache.SetHistoricalData(symbol, interval, from, to, data, 1*time.Hour) // Longer TTL for historical data

	return data, nil
}

// CalculateIndicator calculates a technical indicator with caching
func (s *CachedMarketDataService) CalculateIndicator(ctx context.Context, indicator, symbol, interval string, from, to time.Time, params map[string]interface{}) ([]map[string]interface{}, error) {
	// Check cache first
	if data, found := s.cache.GetIndicatorData(indicator, symbol, interval, from, to, params); found {
		return data, nil
	}

	// Calculate indicator
	data, err := s.service.CalculateIndicator(ctx, indicator, symbol, interval, from, to, params)
	if err != nil {
		return nil, err
	}

	// Cache results
	s.cache.SetIndicatorData(indicator, symbol, interval, from, to, params, data, 1*time.Hour) // Longer TTL for indicator data

	return data, nil
}
