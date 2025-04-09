package marketdata

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

// TestMarketDataService tests the market data service
func TestMarketDataService(t *testing.T) {
	// Create test context
	ctx := context.Background()

	// Initialize components
	config := DefaultServiceConfig()
	dataSourceManager := createTestDataSourceManager()
	dataStorage := createTestDataStorage()
	cacheManager := createTestCacheManager()
	
	// Create service
	service := NewMarketDataService(config, dataSourceManager, dataStorage, cacheManager)
	
	// Test getting market data
	t.Run("GetMarketData", func(t *testing.T) {
		symbols := []string{"AAPL", "MSFT", "GOOG"}
		data, err := service.GetMarketData(ctx, symbols)
		
		if err != nil {
			t.Errorf("Error getting market data: %v", err)
		}
		
		if len(data) != len(symbols) {
			t.Errorf("Expected %d symbols, got %d", len(symbols), len(data))
		}
		
		for _, symbol := range symbols {
			if _, ok := data[symbol]; !ok {
				t.Errorf("Missing data for symbol: %s", symbol)
			}
		}
	})
	
	// Test getting historical data
	t.Run("GetHistoricalData", func(t *testing.T) {
		symbol := "AAPL"
		interval := "1d"
		from := time.Now().AddDate(0, -1, 0) // 1 month ago
		to := time.Now()
		
		data, err := service.GetHistoricalData(ctx, symbol, interval, from, to)
		
		if err != nil {
			t.Errorf("Error getting historical data: %v", err)
		}
		
		if len(data) == 0 {
			t.Errorf("No historical data returned")
		}
		
		// Check data properties
		for _, candle := range data {
			if candle.Symbol != symbol {
				t.Errorf("Expected symbol %s, got %s", symbol, candle.Symbol)
			}
			
			if candle.Interval != interval {
				t.Errorf("Expected interval %s, got %s", interval, candle.Interval)
			}
			
			if candle.Timestamp.Before(from) || candle.Timestamp.After(to) {
				t.Errorf("Candle timestamp %v outside range %v to %v", candle.Timestamp, from, to)
			}
		}
	})
	
	// Test calculating indicators
	t.Run("CalculateIndicator", func(t *testing.T) {
		symbol := "AAPL"
		interval := "1d"
		from := time.Now().AddDate(0, -3, 0) // 3 months ago
		to := time.Now()
		
		// Test SMA indicator
		t.Run("SMA", func(t *testing.T) {
			params := map[string]interface{}{
				"period": float64(14),
				"price":  "close",
			}
			
			result, err := service.CalculateIndicator(ctx, "SMA", symbol, interval, from, to, params)
			
			if err != nil {
				t.Errorf("Error calculating SMA: %v", err)
			}
			
			if len(result) == 0 {
				t.Errorf("No SMA results returned")
			}
		})
		
		// Test RSI indicator
		t.Run("RSI", func(t *testing.T) {
			params := map[string]interface{}{
				"period": float64(14),
				"price":  "close",
			}
			
			result, err := service.CalculateIndicator(ctx, "RSI", symbol, interval, from, to, params)
			
			if err != nil {
				t.Errorf("Error calculating RSI: %v", err)
			}
			
			if len(result) == 0 {
				t.Errorf("No RSI results returned")
			}
			
			// Check RSI values are in range 0-100
			for _, value := range result {
				rsi := value["value"].(float64)
				if rsi < 0 || rsi > 100 {
					t.Errorf("RSI value %f outside valid range 0-100", rsi)
				}
			}
		})
		
		// Test MACD indicator
		t.Run("MACD", func(t *testing.T) {
			params := map[string]interface{}{
				"fastPeriod":   float64(12),
				"slowPeriod":   float64(26),
				"signalPeriod": float64(9),
				"price":        "close",
			}
			
			result, err := service.CalculateIndicator(ctx, "MACD", symbol, interval, from, to, params)
			
			if err != nil {
				t.Errorf("Error calculating MACD: %v", err)
			}
			
			if len(result) == 0 {
				t.Errorf("No MACD results returned")
			}
			
			// Check MACD values contain required fields
			for _, value := range result {
				if _, ok := value["macd"]; !ok {
					t.Errorf("Missing MACD value")
				}
				if _, ok := value["signal"]; !ok {
					t.Errorf("Missing signal value")
				}
				if _, ok := value["histogram"]; !ok {
					t.Errorf("Missing histogram value")
				}
			}
		})
	})
}

// TestHistoricalDataManager tests the historical data manager
func TestHistoricalDataManager(t *testing.T) {
	// Create test context
	ctx := context.Background()

	// Initialize components
	dataSourceManager := createTestDataSourceManager()
	dataStorage := createTestDataStorage()
	cacheConfig := DefaultCacheConfig()
	cache := NewMarketDataCache(cacheConfig)
	
	// Create manager
	manager := NewHistoricalDataManager(dataSourceManager, dataStorage, cache)
	
	// Test getting historical data
	t.Run("GetHistoricalData", func(t *testing.T) {
		symbol := "AAPL"
		interval := "1d"
		from := time.Now().AddDate(0, -1, 0) // 1 month ago
		to := time.Now()
		
		data, err := manager.GetHistoricalData(ctx, symbol, interval, from, to)
		
		if err != nil {
			t.Errorf("Error getting historical data: %v", err)
		}
		
		if len(data) == 0 {
			t.Errorf("No historical data returned")
		}
	})
	
	// Test getting historical data batch
	t.Run("GetHistoricalDataBatch", func(t *testing.T) {
		symbols := []string{"AAPL", "MSFT", "GOOG"}
		interval := "1d"
		from := time.Now().AddDate(0, -1, 0) // 1 month ago
		to := time.Now()
		
		data, err := manager.GetHistoricalDataBatch(ctx, symbols, interval, from, to)
		
		if err != nil {
			t.Errorf("Error getting historical data batch: %v", err)
		}
		
		if len(data) != len(symbols) {
			t.Errorf("Expected %d symbols, got %d", len(symbols), len(data))
		}
		
		for _, symbol := range symbols {
			if _, ok := data[symbol]; !ok {
				t.Errorf("Missing data for symbol: %s", symbol)
			}
		}
	})
	
	// Test getting historical data range
	t.Run("GetHistoricalDataRange", func(t *testing.T) {
		symbol := "AAPL"
		interval := "1d"
		from := time.Now().AddDate(0, -1, 0) // 1 month ago
		to := time.Now()
		limit := 10
		
		data, err := manager.GetHistoricalDataRange(ctx, symbol, interval, from, to, limit)
		
		if err != nil {
			t.Errorf("Error getting historical data range: %v", err)
		}
		
		if len(data) > limit {
			t.Errorf("Expected at most %d candles, got %d", limit, len(data))
		}
	})
}

// TestIndicatorLibrary tests the technical indicator library
func TestIndicatorLibrary(t *testing.T) {
	// Create library
	library := NewIndicatorLibrary()
	
	// Test available indicators
	t.Run("GetAvailableIndicators", func(t *testing.T) {
		indicators := library.GetAvailableIndicators()
		
		if len(indicators) == 0 {
			t.Errorf("No indicators available")
		}
		
		// Check for required indicators
		requiredIndicators := []IndicatorType{SMA, EMA, RSI, MACD, BOLL}
		for _, indicator := range requiredIndicators {
			found := false
			for _, available := range indicators {
				if available == indicator {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Required indicator %s not available", indicator)
			}
		}
	})
	
	// Test SMA indicator
	t.Run("SMA", func(t *testing.T) {
		// Create test data
		data := createTestOHLCVData("AAPL", "1d", 50)
		
		// Get default parameters
		params, err := library.GetDefaultParameters(SMA)
		if err != nil {
			t.Errorf("Error getting default parameters: %v", err)
		}
		
		// Calculate indicator
		result, err := library.Calculate(SMA, data, params)
		if err != nil {
			t.Errorf("Error calculating SMA: %v", err)
		}
		
		// Check result
		period := int(params["period"].(float64))
		expectedLength := len(data) - period + 1
		if len(result.Values) != expectedLength {
			t.Errorf("Expected %d values, got %d", expectedLength, len(result.Values))
		}
	})
	
	// Test RSI indicator
	t.Run("RSI", func(t *testing.T) {
		// Create test data
		data := createTestOHLCVData("AAPL", "1d", 50)
		
		// Get default parameters
		params, err := library.GetDefaultParameters(RSI)
		if err != nil {
			t.Errorf("Error getting default parameters: %v", err)
		}
		
		// Calculate indicator
		result, err := library.Calculate(RSI, data, params)
		if err != nil {
			t.Errorf("Error calculating RSI: %v", err)
		}
		
		// Check result
		period := int(params["period"].(float64))
		expectedLength := len(data) - period
		if len(result.Values) != expectedLength {
			t.Errorf("Expected %d values, got %d", expectedLength, len(result.Values))
		}
		
		// Check RSI values are in range 0-100
		for _, value := range result.Values {
			rsi := value["value"].(float64)
			if rsi < 0 || rsi > 100 {
				t.Errorf("RSI value %f outside valid range 0-100", rsi)
			}
		}
	})
	
	// Test MACD indicator
	t.Run("MACD", func(t *testing.T) {
		// Create test data
		data := createTestOHLCVData("AAPL", "1d", 100)
		
		// Get default parameters
		params, err := library.GetDefaultParameters(MACD)
		if err != nil {
			t.Errorf("Error getting default parameters: %v", err)
		}
		
		// Calculate indicator
		result, err := library.Calculate(MACD, data, params)
		if err != nil {
			t.Errorf("Error calculating MACD: %v", err)
		}
		
		// Check result
		fastPeriod := int(params["fastPeriod"].(float64))
		slowPeriod := int(params["slowPeriod"].(float64))
		signalPeriod := int(params["signalPeriod"].(float64))
		expectedLength := len(data) - (slowPeriod + signalPeriod - 1)
		if len(result.Values) != expectedLength {
			t.Errorf("Expected %d values, got %d", expectedLength, len(result.Values))
		}
		
		// Check MACD values contain required fields
		for _, value := range result.Values {
			if _, ok := value["macd"]; !ok {
				t.Errorf("Missing MACD value")
			}
			if _, ok := value["signal"]; !ok {
				t.Errorf("Missing signal value")
			}
			if _, ok := value["histogram"]; !ok {
				t.Errorf("Missing histogram value")
			}
		}
	})
	
	// Test Bollinger Bands indicator
	t.Run("BOLL", func(t *testing.T) {
		// Create test data
		data := createTestOHLCVData("AAPL", "1d", 50)
		
		// Get default parameters
		params, err := library.GetDefaultParameters(BOLL)
		if err != nil {
			t.Errorf("Error getting default parameters: %v", err)
		}
		
		// Calculate indicator
		result, err := library.Calculate(BOLL, data, params)
		if err != nil {
			t.Errorf("Error calculating Bollinger Bands: %v", err)
		}
		
		// Check result
		period := int(params["period"].(float64))
		expectedLength := len(data) - period + 1
		if len(result.Values) != expectedLength {
			t.Errorf("Expected %d values, got %d", expectedLength, len(result.Values))
		}
		
		// Check Bollinger Bands values contain required fields
		for _, value := range result.Values {
			if _, ok := value["middle"]; !ok {
				t.Errorf("Missing middle band value")
			}
			if _, ok := value["upper"]; !ok {
				t.Errorf("Missing upper band value")
			}
			if _, ok := value["lower"]; !ok {
				t.Errorf("Missing lower band value")
			}
			
			// Check upper band is greater than middle band
			if value["upper"].(float64) <= value["middle"].(float64) {
				t.Errorf("Upper band %f not greater than middle band %f", value["upper"].(float64), value["middle"].(float64))
			}
			
			// Check lower band is less than middle band
			if value["lower"].(float64) >= value["middle"].(float64) {
				t.Errorf("Lower band %f not less than middle band %f", value["lower"].(float64), value["middle"].(float64))
			}
		}
	})
}

// TestCacheManager tests the cache manager
func TestCacheManager(t *testing.T) {
	// Create cache manager
	config := DefaultCacheConfig()
	config.DefaultTTL = 1 * time.Second // Short TTL for testing
	cacheManager := NewCacheManager(config)
	
	// Test setting and getting values
	t.Run("SetGet", func(t *testing.T) {
		key := "test_key"
		value := map[string]interface{}{
			"name":  "Test Value",
			"value": 123.45,
		}
		
		// Set value
		err := cacheManager.Set(key, value, 0)
		if err != nil {
			t.Errorf("Error setting cache value: %v", err)
		}
		
		// Get value
		var result map[string]interface{}
		found := cacheManager.Get(key, &result)
		
		if !found {
			t.Errorf("Value not found in cache")
		}
		
		if result["name"] != value["name"] || result["value"] != value["value"] {
			t.Errorf("Expected %v, got %v", value, result)
		}
	})
	
	// Test TTL expiration
	t.Run("TTLExpiration", func(t *testing.T) {
		key := "expiring_key"
		value := "This will expire"
		
		// Set value with short TTL
		err := cacheManager.Set(key, value, 100*time.Millisecond)
		if err != nil {
			t.Errorf("Error setting cache value: %v", err)
		}
		
		// Get value immediately (should exist)
		var result string
		found := cacheManager.Get(key, &result)
		
		if !found {
			t.Errorf("Value not found in cache immediately after setting")
		}
		
		// Wait for expiration
		time.Sleep(200 * time.Millisecond)
		
		// Get value after expiration (should not exist)
		found = cacheManager.Get(key, &result)
		
		if found {
			t.Errorf("Value found in cache after expiration")
		}
	})
	
	// Test deleting values
	t.Run("Delete", func(t *testing.T) {
		key := "delete_key"
		value := "This will be deleted"
		
		// Set value
		err := cacheManager.Set(key, value, 0)
		if err != nil {
			t.Errorf("Error setting cache value: %v", err)
		}
		
		// Delete value
		cacheManager.Delete(key)
		
		// Get value after deletion (should not exist)
		var result string
		found := cacheManager.Get(key, &result)
		
		if found {
			t.Errorf("Value found in cache after deletion")
		}
	})
	
	// Test clearing cache
	t.Run("Clear", func(t *testing.T) {
		// Set multiple values
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("key_%d", i)
			value := fmt.Sprintf("value_%d", i)
			
			err := cacheManager.Set(key, value, 0)
			if err != nil {
				t.Errorf("Error setting cache value: %v", err)
			}
		}
		
		// Check size
		size := cacheManager.Size()
		if size != 10 {
			t.Errorf("Expected cache size 10, got %d", size)
		}
		
		// Clear cache
		cacheManager.Clear()
		
		// Check size after clearing
		size = cacheManager.Size()
		if size != 0 {
			t.Errorf("Expected cache size 0 after clearing, got %d", size)
		}
	})
	
	// Clean up
	cacheManager.Stop()
}

// TestRealTimeUpdateManager tests the real-time update manager
func TestRealTimeUpdateManager(t *testing.T) {
	// Create test context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize components
	dataSourceManager := createTestDataSourceManager()
	dataStorage := createTestDataStorage()
	cacheConfig := DefaultCacheConfig()
	cacheManager := NewCacheManager(cacheConfig)
	
	// Create manager
	manager := NewRealTimeUpdateManager(dataSourceManager, dataStorage, cacheManager)
	
	// Test subscribing to market data
	t.Run("Subscribe", func(t *testing.T) {
		symbols := []string{"AAPL", "MSFT", "GOOG"}
		
		// Create callback channel
		callbackCh := make(chan MarketData, 10)
		callback := func(data MarketData) {
			callbackCh <- data
		}
		
		// Subscribe to market data
		err := manager.Subscribe(ctx, symbols, callback)
		if err != nil {
			t.Errorf("Error subscribing to market data: %v", err)
		}
		
		// Wait for callbacks
		timeout := time.After(5 * time.Second)
		callbackCount := 0
		
		for callbackCount < len(symbols) {
			select {
			case data := <-callbackCh:
				log.Printf("Received market data for %s", data.Symbol)
				callbackCount++
			case <-timeout:
				t.Errorf("Timeout waiting for callbacks, received %d/%d", callbackCount, len(symbols))
				return
			}
		}
		
		// Check subscribed symbols
		subscribedSymbols := manager.GetSubscribedSymbols()
		if len(subscribedSymbols) != len(symbols) {
			t.Errorf("Expected %d subscribed symbols, got %d", len(symbols), len(subscribedSymbols))
		}
		
		// Unsubscribe from market data
		err = manager.Unsubscribe(ctx, symbols)
		if err != nil {
			t.Errorf("Error unsubscribing from market data: %v", err)
		}
		
		// Check subscribed symbols after unsubscribing
		subscribedSymbols = manager.GetSubscribedSymbols()
		if len(subscribedSymbols) != 0 {
			t.Errorf("Expected 0 subscribed symbols after unsubscribing, got %d", len(subscribedSymbols))
		}
	})
}

// Helper functions for testing

// createTestDataSourceManager creates a test data source manager
func createTestDataSourceManager() *DataSourceManager {
	// Create mock connectors
	mockConnector := &MockDataSourceConnector{
		connected: true,
		marketData: map[string]MarketData{
			"AAPL": {
				Symbol
(Content truncated due to size limit. Use line ranges to read in chunks)