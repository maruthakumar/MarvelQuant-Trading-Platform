package marketdata

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// TimescaleDBStorage implements DataStorage for TimescaleDB
type TimescaleDBStorage struct {
	db *sql.DB
}

// NewTimescaleDBStorage creates a new TimescaleDB storage
func NewTimescaleDBStorage(db *sql.DB) *TimescaleDBStorage {
	return &TimescaleDBStorage{
		db: db,
	}
}

// StoreMarketData stores market data in the database
func (s *TimescaleDBStorage) StoreMarketData(ctx context.Context, data MarketData) error {
	query := `
		INSERT INTO market_data (
			symbol, exchange, last_price, bid_price, ask_price, 
			bid_size, ask_size, volume, open_price, high_price, 
			low_price, close_price, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (symbol, exchange, timestamp) DO UPDATE SET
			last_price = EXCLUDED.last_price,
			bid_price = EXCLUDED.bid_price,
			ask_price = EXCLUDED.ask_price,
			bid_size = EXCLUDED.bid_size,
			ask_size = EXCLUDED.ask_size,
			volume = EXCLUDED.volume,
			open_price = EXCLUDED.open_price,
			high_price = EXCLUDED.high_price,
			low_price = EXCLUDED.low_price,
			close_price = EXCLUDED.close_price
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		data.Symbol,
		data.Exchange,
		data.LastPrice,
		data.BidPrice,
		data.AskPrice,
		data.BidSize,
		data.AskSize,
		data.Volume,
		data.OpenPrice,
		data.HighPrice,
		data.LowPrice,
		data.ClosePrice,
		data.Timestamp,
	)

	return err
}

// StoreOHLCV stores OHLCV data in the database
func (s *TimescaleDBStorage) StoreOHLCV(ctx context.Context, symbol string, interval string, data []OHLCV) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ohlcv (
			symbol, interval, open, high, low, close, volume, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (symbol, interval, timestamp) DO UPDATE SET
			open = EXCLUDED.open,
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			close = EXCLUDED.close,
			volume = EXCLUDED.volume
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert each OHLCV record
	for _, ohlcv := range data {
		_, err := stmt.ExecContext(
			ctx,
			symbol,
			interval,
			ohlcv.Open,
			ohlcv.High,
			ohlcv.Low,
			ohlcv.Close,
			ohlcv.Volume,
			ohlcv.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// GetMarketData gets market data from the database
func (s *TimescaleDBStorage) GetMarketData(ctx context.Context, symbol string) (MarketData, error) {
	query := `
		SELECT 
			symbol, exchange, last_price, bid_price, ask_price, 
			bid_size, ask_size, volume, open_price, high_price, 
			low_price, close_price, timestamp
		FROM market_data
		WHERE symbol = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var data MarketData
	err := s.db.QueryRowContext(ctx, query, symbol).Scan(
		&data.Symbol,
		&data.Exchange,
		&data.LastPrice,
		&data.BidPrice,
		&data.AskPrice,
		&data.BidSize,
		&data.AskSize,
		&data.Volume,
		&data.OpenPrice,
		&data.HighPrice,
		&data.LowPrice,
		&data.ClosePrice,
		&data.Timestamp,
	)

	if err != nil {
		return MarketData{}, err
	}

	return data, nil
}

// GetOHLCV gets OHLCV data from the database
func (s *TimescaleDBStorage) GetOHLCV(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error) {
	query := `
		SELECT 
			symbol, interval, open, high, low, close, volume, timestamp
		FROM ohlcv
		WHERE symbol = $1 AND interval = $2 AND timestamp >= $3 AND timestamp <= $4
		ORDER BY timestamp ASC
	`

	rows, err := s.db.QueryContext(ctx, query, symbol, interval, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []OHLCV
	for rows.Next() {
		var ohlcv OHLCV
		err := rows.Scan(
			&ohlcv.Symbol,
			&ohlcv.Interval,
			&ohlcv.Open,
			&ohlcv.High,
			&ohlcv.Low,
			&ohlcv.Close,
			&ohlcv.Volume,
			&ohlcv.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ohlcv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetLatestOHLCV gets the latest OHLCV data from the database
func (s *TimescaleDBStorage) GetLatestOHLCV(ctx context.Context, symbol string, interval string, limit int) ([]OHLCV, error) {
	query := `
		SELECT 
			symbol, interval, open, high, low, close, volume, timestamp
		FROM ohlcv
		WHERE symbol = $1 AND interval = $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	rows, err := s.db.QueryContext(ctx, query, symbol, interval, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []OHLCV
	for rows.Next() {
		var ohlcv OHLCV
		err := rows.Scan(
			&ohlcv.Symbol,
			&ohlcv.Interval,
			&ohlcv.Open,
			&ohlcv.High,
			&ohlcv.Low,
			&ohlcv.Close,
			&ohlcv.Volume,
			&ohlcv.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, ohlcv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse the order to get ascending order by timestamp
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// StoreIndicatorValue stores an indicator value in the database
func (s *TimescaleDBStorage) StoreIndicatorValue(ctx context.Context, value IndicatorValue) error {
	query := `
		INSERT INTO indicator_values (
			symbol, indicator, value, values, timestamp, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		ON CONFLICT (symbol, indicator, timestamp) DO UPDATE SET
			value = EXCLUDED.value,
			values = EXCLUDED.values,
			metadata = EXCLUDED.metadata
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		value.Symbol,
		value.Indicator,
		value.Value,
		pq.Array(value.Values),
		value.Timestamp,
		pq.Array(value.Metadata),
	)

	return err
}

// GetIndicatorValues gets indicator values from the database
func (s *TimescaleDBStorage) GetIndicatorValues(ctx context.Context, symbol string, indicator string, from, to time.Time) ([]IndicatorValue, error) {
	query := `
		SELECT 
			symbol, indicator, value, values, timestamp, metadata
		FROM indicator_values
		WHERE symbol = $1 AND indicator = $2 AND timestamp >= $3 AND timestamp <= $4
		ORDER BY timestamp ASC
	`

	rows, err := s.db.QueryContext(ctx, query, symbol, indicator, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []IndicatorValue
	for rows.Next() {
		var value IndicatorValue
		var valuesArray, metadataArray []byte
		err := rows.Scan(
			&value.Symbol,
			&value.Indicator,
			&value.Value,
			&valuesArray,
			&value.Timestamp,
			&metadataArray,
		)
		if err != nil {
			return nil, err
		}

		// Parse values and metadata
		if err := pq.Array(&value.Values).Scan(valuesArray); err != nil {
			return nil, err
		}
		if err := pq.Array(&value.Metadata).Scan(metadataArray); err != nil {
			return nil, err
		}

		result = append(result, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// InitializeSchema initializes the database schema
func (s *TimescaleDBStorage) InitializeSchema(ctx context.Context) error {
	// Create market_data table
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS market_data (
			symbol TEXT NOT NULL,
			exchange TEXT NOT NULL,
			last_price DOUBLE PRECISION,
			bid_price DOUBLE PRECISION,
			ask_price DOUBLE PRECISION,
			bid_size INTEGER,
			ask_size INTEGER,
			volume INTEGER,
			open_price DOUBLE PRECISION,
			high_price DOUBLE PRECISION,
			low_price DOUBLE PRECISION,
			close_price DOUBLE PRECISION,
			timestamp TIMESTAMPTZ NOT NULL,
			PRIMARY KEY (symbol, exchange, timestamp)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create market_data table: %w", err)
	}

	// Convert market_data to hypertable
	_, err = s.db.ExecContext(ctx, `
		SELECT create_hypertable('market_data', 'timestamp', if_not_exists => TRUE)
	`)
	if err != nil {
		return fmt.Errorf("failed to convert market_data to hypertable: %w", err)
	}

	// Create ohlcv table
	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS ohlcv (
			symbol TEXT NOT NULL,
			interval TEXT NOT NULL,
			open DOUBLE PRECISION NOT NULL,
			high DOUBLE PRECISION NOT NULL,
			low DOUBLE PRECISION NOT NULL,
			close DOUBLE PRECISION NOT NULL,
			volume INTEGER NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			PRIMARY KEY (symbol, interval, timestamp)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create ohlcv table: %w", err)
	}

	// Convert ohlcv to hypertable
	_, err = s.db.ExecContext(ctx, `
		SELECT create_hypertable('ohlcv', 'timestamp', if_not_exists => TRUE)
	`)
	if err != nil {
		return fmt.Errorf("failed to convert ohlcv to hypertable: %w", err)
	}

	// Create indicator_values table
	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS indicator_values (
			symbol TEXT NOT NULL,
			indicator TEXT NOT NULL,
			value DOUBLE PRECISION,
			values JSONB,
			timestamp TIMESTAMPTZ NOT NULL,
			metadata JSONB,
			PRIMARY KEY (symbol, indicator, timestamp)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create indicator_values table: %w", err)
	}

	// Convert indicator_values to hypertable
	_, err = s.db.ExecContext(ctx, `
		SELECT create_hypertable('indicator_values', 'timestamp', if_not_exists => TRUE)
	`)
	if err != nil {
		return fmt.Errorf("failed to convert indicator_values to hypertable: %w", err)
	}

	return nil
}
