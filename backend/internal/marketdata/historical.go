package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

// HistoricalDataManager manages retrieval and storage of historical market data
type HistoricalDataManager struct {
	dataSourceManager *DataSourceManager
	dataStorage       DataStorage
	cache             *MarketDataCache
	mutex             sync.RWMutex
}

// NewHistoricalDataManager creates a new historical data manager
func NewHistoricalDataManager(
	dataSourceManager *DataSourceManager,
	dataStorage DataStorage,
	cache *MarketDataCache,
) *HistoricalDataManager {
	return &HistoricalDataManager{
		dataSourceManager: dataSourceManager,
		dataStorage:       dataStorage,
		cache:             cache,
	}
}

// GetHistoricalData retrieves historical data for a symbol and interval
func (m *HistoricalDataManager) GetHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
) ([]OHLCV, error) {
	// Check cache first
	if data, found := m.cache.GetHistoricalData(symbol, interval, from, to); found {
		return data, nil
	}

	// Try to get from storage
	data, err := m.dataStorage.GetHistoricalData(ctx, symbol, interval, from, to)
	if err == nil && len(data) > 0 {
		// Cache the data
		m.cache.SetHistoricalData(symbol, interval, from, to, data, 1*time.Hour)
		return data, nil
	}

	// Get from data source
	data, err = m.dataSourceManager.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data: %w", err)
	}

	// Store in storage asynchronously
	go func() {
		if err := m.dataStorage.StoreHistoricalData(context.Background(), data); err != nil {
			log.Printf("Error storing historical data: %v", err)
		}
	}()

	// Cache the data
	m.cache.SetHistoricalData(symbol, interval, from, to, data, 1*time.Hour)

	return data, nil
}

// GetHistoricalDataBatch retrieves historical data for multiple symbols
func (m *HistoricalDataManager) GetHistoricalDataBatch(
	ctx context.Context,
	symbols []string,
	interval string,
	from, to time.Time,
) (map[string][]OHLCV, error) {
	result := make(map[string][]OHLCV)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(symbols))

	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()

			data, err := m.GetHistoricalData(ctx, sym, interval, from, to)
			if err != nil {
				errChan <- fmt.Errorf("error getting historical data for %s: %w", sym, err)
				return
			}

			mu.Lock()
			result[sym] = data
			mu.Unlock()
		}(symbol)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return result, fmt.Errorf("errors retrieving historical data: %v", errs)
	}

	return result, nil
}

// GetHistoricalDataRange retrieves historical data for a specific date range
func (m *HistoricalDataManager) GetHistoricalDataRange(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
	limit int,
) ([]OHLCV, error) {
	// Get historical data
	data, err := m.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}

	// Sort by timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Before(data[j].Timestamp)
	})

	// Apply limit if specified
	if limit > 0 && limit < len(data) {
		data = data[len(data)-limit:]
	}

	return data, nil
}

// GetHistoricalDataWithAdjustments retrieves historical data with adjustments for splits and dividends
func (m *HistoricalDataManager) GetHistoricalDataWithAdjustments(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
	adjustments []string,
) ([]OHLCV, error) {
	// Get historical data
	data, err := m.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}

	// Apply adjustments
	adjustedData := make([]OHLCV, len(data))
	copy(adjustedData, data)

	// In a real implementation, this would apply adjustments based on corporate actions
	// For now, we'll just return the original data
	return adjustedData, nil
}

// BackfillHistoricalData backfills historical data for a symbol
func (m *HistoricalDataManager) BackfillHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
) error {
	// Get historical data
	data, err := m.dataSourceManager.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return fmt.Errorf("failed to get historical data: %w", err)
	}

	// Store in storage
	if err := m.dataStorage.StoreHistoricalData(ctx, data); err != nil {
		return fmt.Errorf("failed to store historical data: %w", err)
	}

	return nil
}

// HistoricalDataProcessor processes historical data
type HistoricalDataProcessor struct {
	processors []DataProcessor
}

// NewHistoricalDataProcessor creates a new historical data processor
func NewHistoricalDataProcessor(processors ...DataProcessor) *HistoricalDataProcessor {
	return &HistoricalDataProcessor{
		processors: processors,
	}
}

// Process processes historical data
func (p *HistoricalDataProcessor) Process(data []OHLCV) ([]OHLCV, error) {
	result := make([]OHLCV, len(data))
	copy(result, data)

	// Apply each processor
	for _, processor := range p.processors {
		for i, candle := range result {
			processed, err := processor.Process(candle)
			if err != nil {
				return nil, fmt.Errorf("error processing candle: %w", err)
			}
			if processedCandle, ok := processed.(OHLCV); ok {
				result[i] = processedCandle
			}
		}
	}

	return result, nil
}

// HistoricalDataExporter exports historical data to various formats
type HistoricalDataExporter struct {
	// Add exporter-specific fields if needed
}

// NewHistoricalDataExporter creates a new historical data exporter
func NewHistoricalDataExporter() *HistoricalDataExporter {
	return &HistoricalDataExporter{}
}

// ExportToCSV exports historical data to CSV format
func (e *HistoricalDataExporter) ExportToCSV(data []OHLCV) ([]byte, error) {
	var result strings.Builder

	// Write header
	result.WriteString("Timestamp,Open,High,Low,Close,Volume\n")

	// Write data
	for _, candle := range data {
		result.WriteString(fmt.Sprintf(
			"%s,%.2f,%.2f,%.2f,%.2f,%d\n",
			candle.Timestamp.Format("2006-01-02 15:04:05"),
			candle.Open,
			candle.High,
			candle.Low,
			candle.Close,
			candle.Volume,
		))
	}

	return []byte(result.String()), nil
}

// ExportToJSON exports historical data to JSON format
func (e *HistoricalDataExporter) ExportToJSON(data []OHLCV) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// HistoricalDataImporter imports historical data from various formats
type HistoricalDataImporter struct {
	// Add importer-specific fields if needed
}

// NewHistoricalDataImporter creates a new historical data importer
func NewHistoricalDataImporter() *HistoricalDataImporter {
	return &HistoricalDataImporter{}
}

// ImportFromCSV imports historical data from CSV format
func (i *HistoricalDataImporter) ImportFromCSV(symbol, interval string, data []byte) ([]OHLCV, error) {
	var result []OHLCV

	// Split into lines
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid CSV data: not enough lines")
	}

	// Skip header
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		// Split line
		fields := strings.Split(line, ",")
		if len(fields) < 6 {
			return nil, fmt.Errorf("invalid CSV data: not enough fields")
		}

		// Parse timestamp
		timestamp, err := time.Parse("2006-01-02 15:04:05", fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}

		// Parse values
		open, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid open: %w", err)
		}
		high, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid high: %w", err)
		}
		low, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid low: %w", err)
		}
		close, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid close: %w", err)
		}
		volume, err := strconv.ParseInt(fields[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid volume: %w", err)
		}

		// Create candle
		candle := OHLCV{
			Symbol:    symbol,
			Interval:  interval,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    int(volume),
			Timestamp: timestamp,
		}

		result = append(result, candle)
	}

	return result, nil
}

// ImportFromJSON imports historical data from JSON format
func (i *HistoricalDataImporter) ImportFromJSON(data []byte) ([]OHLCV, error) {
	var result []OHLCV
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %w", err)
	}
	return result, nil
}

// HistoricalDataService provides a service for historical data operations
type HistoricalDataService struct {
	manager  *HistoricalDataManager
	exporter *HistoricalDataExporter
	importer *HistoricalDataImporter
}

// NewHistoricalDataService creates a new historical data service
func NewHistoricalDataService(
	manager *HistoricalDataManager,
) *HistoricalDataService {
	return &HistoricalDataService{
		manager:  manager,
		exporter: NewHistoricalDataExporter(),
		importer: NewHistoricalDataImporter(),
	}
}

// GetHistoricalData retrieves historical data
func (s *HistoricalDataService) GetHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
) ([]OHLCV, error) {
	return s.manager.GetHistoricalData(ctx, symbol, interval, from, to)
}

// GetHistoricalDataBatch retrieves historical data for multiple symbols
func (s *HistoricalDataService) GetHistoricalDataBatch(
	ctx context.Context,
	symbols []string,
	interval string,
	from, to time.Time,
) (map[string][]OHLCV, error) {
	return s.manager.GetHistoricalDataBatch(ctx, symbols, interval, from, to)
}

// ExportHistoricalData exports historical data to the specified format
func (s *HistoricalDataService) ExportHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
	format string,
) ([]byte, error) {
	// Get historical data
	data, err := s.manager.GetHistoricalData(ctx, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}

	// Export to specified format
	switch format {
	case "csv":
		return s.exporter.ExportToCSV(data)
	case "json":
		return s.exporter.ExportToJSON(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ImportHistoricalData imports historical data from the specified format
func (s *HistoricalDataService) ImportHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	data []byte,
	format string,
) error {
	// Import from specified format
	var ohlcvData []OHLCV
	var err error

	switch format {
	case "csv":
		ohlcvData, err = s.importer.ImportFromCSV(symbol, interval, data)
	case "json":
		ohlcvData, err = s.importer.ImportFromJSON(data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return err
	}

	// Store in storage
	return s.manager.dataStorage.StoreHistoricalData(ctx, ohlcvData)
}

// BackfillHistoricalData backfills historical data
func (s *HistoricalDataService) BackfillHistoricalData(
	ctx context.Context,
	symbol string,
	interval string,
	from, to time.Time,
) error {
	return s.manager.BackfillHistoricalData(ctx, symbol, interval, from, to)
}
