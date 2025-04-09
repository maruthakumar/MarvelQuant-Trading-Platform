package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// RedisCacheManager implements CacheManager for Redis
type RedisCacheManager struct {
	client *redis.Client
	mutex  sync.RWMutex
}

// NewRedisCacheManager creates a new Redis cache manager
func NewRedisCacheManager(addr, password string, db int) *RedisCacheManager {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCacheManager{
		client: client,
	}
}

// Get gets a value from the cache
func (m *RedisCacheManager) Get(key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	val, err := m.client.Get(key).Result()
	if err != nil {
		return nil, false
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, false
	}

	return result, true
}

// Set sets a value in the cache
func (m *RedisCacheManager) Set(key string, value interface{}, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return m.client.Set(key, val, expiration).Err()
}

// Delete deletes a value from the cache
func (m *RedisCacheManager) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.client.Del(key).Err()
}

// Clear clears the cache
func (m *RedisCacheManager) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.client.FlushDB().Err()
}

// InMemoryCacheManager implements CacheManager for in-memory cache
type InMemoryCacheManager struct {
	cache map[string]cacheEntry
	mutex sync.RWMutex
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCacheManager creates a new in-memory cache manager
func NewInMemoryCacheManager() *InMemoryCacheManager {
	manager := &InMemoryCacheManager{
		cache: make(map[string]cacheEntry),
	}

	// Start cleanup goroutine
	go manager.cleanup()

	return manager
}

// Get gets a value from the cache
func (m *InMemoryCacheManager) Get(key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, ok := m.cache[key]
	if !ok {
		return nil, false
	}

	// Check if entry has expired
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		return nil, false
	}

	return entry.value, true
}

// Set sets a value in the cache
func (m *InMemoryCacheManager) Set(key string, value interface{}, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var expirationTime time.Time
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration)
	}

	m.cache[key] = cacheEntry{
		value:      value,
		expiration: expirationTime,
	}

	return nil
}

// Delete deletes a value from the cache
func (m *InMemoryCacheManager) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.cache, key)
	return nil
}

// Clear clears the cache
func (m *InMemoryCacheManager) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.cache = make(map[string]cacheEntry)
	return nil
}

// cleanup periodically removes expired entries
func (m *InMemoryCacheManager) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		now := time.Now()
		for key, entry := range m.cache {
			if !entry.expiration.IsZero() && now.After(entry.expiration) {
				delete(m.cache, key)
			}
		}
		m.mutex.Unlock()
	}
}

// AlphaVantageConnector implements DataSourceConnector for Alpha Vantage API
type AlphaVantageConnector struct {
	apiKey      string
	isConnected bool
	baseURL     string
	httpClient  *http.Client
	mutex       sync.RWMutex
	callbacks   map[string][]MarketDataCallback
	callbacksMu sync.RWMutex
}

// NewAlphaVantageConnector creates a new Alpha Vantage connector
func NewAlphaVantageConnector(apiKey string) *AlphaVantageConnector {
	return &AlphaVantageConnector{
		apiKey:     apiKey,
		baseURL:    "https://www.alphavantage.co/query",
		httpClient: &http.Client{Timeout: 10 * time.Second},
		callbacks:  make(map[string][]MarketDataCallback),
	}
}

// Connect connects to the Alpha Vantage API
func (c *AlphaVantageConnector) Connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return nil
	}

	// Alpha Vantage API doesn't require a connection, so we just set the flag
	c.isConnected = true
	return nil
}

// Disconnect disconnects from the Alpha Vantage API
func (c *AlphaVantageConnector) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.isConnected = false
	return nil
}

// IsConnected checks if the connector is connected
func (c *AlphaVantageConnector) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// GetMarketData gets market data for the specified symbols
func (c *AlphaVantageConnector) GetMarketData(ctx context.Context, symbols []string) (map[string]MarketData, error) {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("not connected to Alpha Vantage API")
	}
	c.mutex.RUnlock()

	result := make(map[string]MarketData)

	// Alpha Vantage API doesn't support batch requests, so we need to make a request for each symbol
	for _, symbol := range symbols {
		// Create request URL
		req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create quote request: %w", err)
		}

		// Add query parameters
		q := req.URL.Query()
		q.Add("function", "GLOBAL_QUOTE")
		q.Add("symbol", symbol)
		q.Add("apikey", c.apiKey)
		req.URL.RawQuery = q.Encode()

		// Send request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send quote request: %w", err)
		}
		defer resp.Body.Close()

		// Read response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read quote response: %w", err)
		}

		// Parse response
		var quoteResp struct {
			GlobalQuote struct {
				Symbol           string `json:"01. symbol"`
				Open             string `json:"02. open"`
				High             string `json:"03. high"`
				Low              string `json:"04. low"`
				Price            string `json:"05. price"`
				Volume           string `json:"06. volume"`
				LatestTradingDay string `json:"07. latest trading day"`
				PreviousClose    string `json:"08. previous close"`
				Change           string `json:"09. change"`
				ChangePercent    string `json:"10. change percent"`
			} `json:"Global Quote"`
		}
		if err := json.Unmarshal(respBody, &quoteResp); err != nil {
			return nil, fmt.Errorf("failed to parse quote response: %w", err)
		}

		// Check if we got data
		if quoteResp.GlobalQuote.Symbol == "" {
			continue
		}

		// Parse values
		open, _ := strconv.ParseFloat(quoteResp.GlobalQuote.Open, 64)
		high, _ := strconv.ParseFloat(quoteResp.GlobalQuote.High, 64)
		low, _ := strconv.ParseFloat(quoteResp.GlobalQuote.Low, 64)
		price, _ := strconv.ParseFloat(quoteResp.GlobalQuote.Price, 64)
		volume, _ := strconv.ParseInt(quoteResp.GlobalQuote.Volume, 10, 64)
		prevClose, _ := strconv.ParseFloat(quoteResp.GlobalQuote.PreviousClose, 64)

		// Parse date
		date, err := time.Parse("2006-01-02", quoteResp.GlobalQuote.LatestTradingDay)
		if err != nil {
			date = time.Now()
		}

		// Create market data
		result[symbol] = MarketData{
			Symbol:     symbol,
			Exchange:   "US", // Alpha Vantage doesn't provide exchange info
			LastPrice:  price,
			BidPrice:   0, // Alpha Vantage doesn't provide bid/ask
			AskPrice:   0,
			BidSize:    0,
			AskSize:    0,
			Volume:     int(volume),
			OpenPrice:  open,
			HighPrice:  high,
			LowPrice:   low,
			ClosePrice: prevClose,
			Timestamp:  date,
		}
	}

	return result, nil
}

// GetHistoricalData gets historical data for the specified symbol and interval
func (c *AlphaVantageConnector) GetHistoricalData(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error) {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("not connected to Alpha Vantage API")
	}
	c.mutex.RUnlock()

	// Map interval to Alpha Vantage interval
	var function, outputsize string
	switch interval {
	case "1m":
		function = "TIME_SERIES_INTRADAY"
		outputsize = "full"
		interval = "1min"
	case "5m":
		function = "TIME_SERIES_INTRADAY"
		outputsize = "full"
		interval = "5min"
	case "15m":
		function = "TIME_SERIES_INTRADAY"
		outputsize = "full"
		interval = "15min"
	case "30m":
		function = "TIME_SERIES_INTRADAY"
		outputsize = "full"
		interval = "30min"
	case "1h":
		function = "TIME_SERIES_INTRADAY"
		outputsize = "full"
		interval = "60min"
	case "1d":
		function = "TIME_SERIES_DAILY"
		outputsize = "full"
	case "1w":
		function = "TIME_SERIES_WEEKLY"
		outputsize = "full"
	case "1mo":
		function = "TIME_SERIES_MONTHLY"
		outputsize = "full"
	default:
		return nil, fmt.Errorf("unsupported interval: %s", interval)
	}

	// Create request URL
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create historical data request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("function", function)
	q.Add("symbol", symbol)
	q.Add("apikey", c.apiKey)
	q.Add("outputsize", outputsize)
	if function == "TIME_SERIES_INTRADAY" {
		q.Add("interval", interval)
	}
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send historical data request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read historical data response: %w", err)
	}

	// Parse response based on function
	var result []OHLCV
	if function == "TIME_SERIES_INTRADAY" {
		var resp struct {
			MetaData struct {
				Symbol string `json:"2. Symbol"`
			} `json:"Meta Data"`
			TimeSeries map[string]struct {
				Open   string `json:"1. open"`
				High   string `json:"2. high"`
				Low    string `json:"3. low"`
				Close  string `json:"4. close"`
				Volume string `json:"5. volume"`
			} `json:"Time Series (` + interval + `)"`
		}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse intraday response: %w", err)
		}

		for dateStr, data := range resp.TimeSeries {
			date, err := time.Parse("2006-01-02 15:04:05", dateStr)
			if err != nil {
				continue
			}

			// Skip if outside date range
			if date.Before(from) || date.After(to) {
				continue
			}

			open, _ := strconv.ParseFloat(data.Open, 64)
			high, _ := strconv.ParseFloat(data.High, 64)
			low, _ := strconv.ParseFloat(data.Low, 64)
			close, _ := strconv.ParseFloat(data.Close, 64)
			volume, _ := strconv.ParseInt(data.Volume, 10, 64)

			result = append(result, OHLCV{
				Symbol:    symbol,
				Interval:  interval,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    int(volume),
				Timestamp: date,
			})
		}
	} else if function == "TIME_SERIES_DAILY" {
		var resp struct {
			MetaData struct {
				Symbol string `json:"2. Symbol"`
			} `json:"Meta Data"`
			TimeSeries map[string]struct {
				Open   string `json:"1. open"`
				High   string `json:"2. high"`
				Low    string `json:"3. low"`
				Close  string `json:"4. close"`
				Volume string `json:"5. volume"`
			} `json:"Time Series (Daily)"`
		}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse daily response: %w", err)
		}

		for dateStr, data := range resp.TimeSeries {
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// Skip if outside date range
			if date.Before(from) || date.After(to) {
				continue
			}

			open, _ := strconv.ParseFloat(data.Open, 64)
			high, _ := strconv.ParseFloat(data.High, 64)
			low, _ := strconv.ParseFloat(data.Low, 64)
			close, _ := strconv.ParseFloat(data.Close, 64)
			volume, _ := strconv.ParseInt(data.Volume, 10, 64)

			result = append(result, OHLCV{
				Symbol:    symbol,
				Interval:  "1d",
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    int(volume),
				Timestamp: date,
			})
		}
	} else if function == "TIME_SERIES_WEEKLY" {
		var resp struct {
			MetaData struct {
				Symbol string `json:"2. Symbol"`
			} `json:"Meta Data"`
			TimeSeries map[string]struct {
				Open   string `json:"1. open"`
				High   string `json:"2. high"`
				Low    string `json:"3. low"`
				Close  string `json:"4. close"`
				Volume string `json:"5. volume"`
			} `json:"Weekly Time Series"`
		}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse weekly response: %w", err)
		}

		for dateStr, data := range resp.TimeSeries {
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// Skip if outside date range
			if date.Before(from) || date.After(to) {
				continue
			}

			open, _ := strconv.ParseFloat(data.Open, 64)
			high, _ := strconv.ParseFloat(data.High, 64)
			low, _ := strconv.ParseFloat(data.Low, 64)
			close, _ := strconv.ParseFloat(data.Close, 64)
			volume, _ := strconv.ParseInt(data.Volume, 10, 64)

			result = append(result, OHLCV{
				Symbol:    symbol,
				Interval:  "1w",
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    int(volume),
				Timestamp: date,
			})
		}
	} else if function == "TIME_SERIES_MONTHLY" {
		var resp struct {
			MetaData struct {
				Symbol string `json:"2. Symbol"`
			} `json:"Meta Data"`
			TimeSeries map[string]struct {
				Open   string `json:"1. open"`
				High   string `json:"2. high"`
				Low    string `json:"3. low"`
				Close  string `json:"4. close"`
				Volume string `json:"5. volume"`
			} `json:"Monthly Time Series"`
		}
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse monthly response: %w", err)
		}

		for dateStr, data := range resp.TimeSeries {
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// Skip if outside date range
			if date.Before(from) || date.After(to) {
				continue
			}

			open, _ := strconv.ParseFloat(data.Open, 64)
			high, _ := strconv.ParseFloat(data.High, 64)
			low, _ := strconv.ParseFloat(data.Low, 64)
			close, _ := strconv.ParseFloat(data.Close, 64)
			volume, _ := strconv.ParseInt(data.Volume, 10, 64)

			result = append(result, OHLCV{
				Symbol:    symbol,
				Interval:  "1mo",
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    int(volume),
				Timestamp: date,
			})
		}
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result, nil
}

// SubscribeToMarketData subscribes to market data for the specified symbols
func (c *AlphaVantageConnector) SubscribeToMarketData(ctx context.Context, symbols []string, callback MarketDataCallback) error {
	c.mutex.RLock()
	if !c.isConnected {
		c.mutex.RUnlock()
		return fmt.Errorf("not connected to Alpha Vantage API")
	}
	c.mutex.RUnlock()

	// Register callback
	c.callbacksMu.Lock()
	for _, symbol := range symbols {
		c.callbacks[symbol] = append(c.callbacks[symbol], callback)
	}
	c.callbacksMu.Unlock()

	// Alpha Vantage API doesn't support real-time subscriptions
	// We'll simulate it by polling at regular intervals
	go func() {
		ticker := time.NewTicker(60 * time.Second) // Alpha Vantage has strict rate limits
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get current symbols
				c.callbacksMu.RLock()
				if len(c.callbacks) == 0 {
					c.callbacksMu.RUnlock()
					continue
				}
				
				// Get list of symbols to fetch
				currentSymbols := make([]string, 0, len(c.callbacks))
				for symbol := range c.callba
(Content truncated due to size limit. Use line ranges to read in chunks)