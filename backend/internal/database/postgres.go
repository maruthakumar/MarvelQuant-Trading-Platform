package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config Config) (*PostgresDB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &PostgresDB{pool: pool}, nil
}

// Close closes the database connection
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// GetPool returns the connection pool
func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.pool
}

// InitSchema initializes the database schema
func (db *PostgresDB) InitSchema(ctx context.Context) error {
	// Create TimescaleDB extension if not exists
	_, err := db.pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;")
	if err != nil {
		return fmt.Errorf("failed to create TimescaleDB extension: %w", err)
	}

	// Create tables
	err = db.createTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Create hypertables
	err = db.createHypertables(ctx)
	if err != nil {
		return fmt.Errorf("failed to create hypertables: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// createTables creates the necessary tables
func (db *PostgresDB) createTables(ctx context.Context) error {
	// Create users table
	_, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create strategies table
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS strategies (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			name VARCHAR(100) NOT NULL,
			description TEXT,
			is_active BOOLEAN DEFAULT false,
			max_profit_value NUMERIC(15, 5),
			max_profit_type VARCHAR(20),
			target_action VARCHAR(20),
			target_monitoring VARCHAR(20),
			max_loss_value NUMERIC(15, 5),
			max_loss_type VARCHAR(20),
			sl_action VARCHAR(20),
			sl_monitoring VARCHAR(20),
			profit_locking_enabled BOOLEAN DEFAULT false,
			profit_locking_threshold NUMERIC(15, 5),
			profit_locking_value NUMERIC(15, 5),
			profit_trailing_enabled BOOLEAN DEFAULT false,
			profit_trailing_step NUMERIC(15, 5),
			profit_trailing_value NUMERIC(15, 5),
			sl_trailing_enabled BOOLEAN DEFAULT false,
			sl_trailing_step NUMERIC(15, 5),
			sl_trailing_value NUMERIC(15, 5),
			scheduling_enabled BOOLEAN DEFAULT false,
			scheduling_type VARCHAR(20),
			scheduling_config JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create portfolios table
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS portfolios (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			strategy_id INTEGER REFERENCES strategies(id),
			name VARCHAR(100) NOT NULL,
			symbol VARCHAR(20) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			expiry VARCHAR(20),
			default_lots INTEGER DEFAULT 1,
			is_active BOOLEAN DEFAULT false,
			is_positional BOOLEAN DEFAULT false,
			buy_trades_first BOOLEAN DEFAULT true,
			allow_far_strikes BOOLEAN DEFAULT false,
			use_implied_futures BOOLEAN DEFAULT false,
			product VARCHAR(10) DEFAULT 'NRML',
			leg_failure_action VARCHAR(20) DEFAULT 'KEEP_PLACED_LEGS',
			legs_execution VARCHAR(20) DEFAULT 'PARALLEL',
			max_lots INTEGER DEFAULT 1,
			premium_gap NUMERIC(15, 5) DEFAULT 0,
			run_on_days VARCHAR(50) DEFAULT 'Mon,Tue,Wed,Thu,Fri',
			start_time TIME DEFAULT '09:15:00',
			end_time TIME DEFAULT '15:15:00',
			sq_off_time TIME DEFAULT '15:20:00',
			execution_mode VARCHAR(20) DEFAULT 'START_TIME',
			entry_order_type VARCHAR(20) DEFAULT 'MARKET',
			range_breakout_enabled BOOLEAN DEFAULT false,
			range_end_time TIME,
			high_buffer NUMERIC(15, 5) DEFAULT 0,
			low_buffer NUMERIC(15, 5) DEFAULT 0,
			use_opposite_side_for_sl BOOLEAN DEFAULT false,
			range_buffer NUMERIC(15, 5) DEFAULT 0,
			dynamic_hedge_enabled BOOLEAN DEFAULT false,
			hedge_type VARCHAR(20) DEFAULT 'DELTA_NEUTRAL',
			hedge_interval INTEGER DEFAULT 300,
			hedge_threshold NUMERIC(15, 5) DEFAULT 0.1,
			max_profit_value NUMERIC(15, 5),
			max_profit_type VARCHAR(20) DEFAULT 'ABSOLUTE',
			max_loss_value NUMERIC(15, 5),
			max_loss_type VARCHAR(20) DEFAULT 'ABSOLUTE',
			monitoring_frequency INTEGER DEFAULT 5,
			monitoring_type VARCHAR(20) DEFAULT 'REAL_TIME',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create portfolio_legs table
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS portfolio_legs (
			id SERIAL PRIMARY KEY,
			portfolio_id INTEGER REFERENCES portfolios(id),
			leg_id INTEGER NOT NULL,
			buy_sell VARCHAR(4) NOT NULL,
			option_type VARCHAR(3) NOT NULL,
			strike VARCHAR(20) NOT NULL,
			lots INTEGER DEFAULT 1,
			expiry VARCHAR(20),
			is_idle BOOLEAN DEFAULT false,
			ltp NUMERIC(15, 5) DEFAULT 0,
			hedge_required BOOLEAN DEFAULT false,
			wait_and_trade VARCHAR(20),
			target_type VARCHAR(20) DEFAULT 'NONE',
			target_value NUMERIC(15, 5) DEFAULT 0,
			trail_target BOOLEAN DEFAULT false,
			trail_target_value NUMERIC(15, 5) DEFAULT 0,
			sl_type VARCHAR(20) DEFAULT 'NONE',
			sl_value NUMERIC(15, 5) DEFAULT 0,
			trail_sl BOOLEAN DEFAULT false,
			trail_sl_value NUMERIC(15, 5) DEFAULT 0,
			on_target_action VARCHAR(30) DEFAULT 'NONE',
			on_target_action_params JSONB,
			on_stoploss_action VARCHAR(30) DEFAULT 'NONE',
			on_stoploss_action_params JSONB,
			on_start_action VARCHAR(30) DEFAULT 'NONE',
			on_start_action_params JSONB,
			start_time TIME,
			spread_limit NUMERIC(15, 5) DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create orders table
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			portfolio_id INTEGER REFERENCES portfolios(id),
			leg_id INTEGER,
			broker_order_id VARCHAR(50),
			symbol VARCHAR(20) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			order_type VARCHAR(20) NOT NULL,
			transaction_type VARCHAR(4) NOT NULL,
			product_type VARCHAR(10) NOT NULL,
			quantity INTEGER NOT NULL,
			price NUMERIC(15, 5),
			trigger_price NUMERIC(15, 5),
			status VARCHAR(20) NOT NULL,
			message TEXT,
			order_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create market_data table for time-series data
	_, err = db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS market_data (
			time TIMESTAMP WITH TIME ZONE NOT NULL,
			symbol VARCHAR(20) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			last_price NUMERIC(15, 5) NOT NULL,
			bid_price NUMERIC(15, 5),
			ask_price NUMERIC(15, 5),
			volume INTEGER,
			open_interest INTEGER,
			PRIMARY KEY (time, symbol, exchange)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

// createHypertables creates TimescaleDB hypertables
func (db *PostgresDB) createHypertables(ctx context.Context) error {
	// Convert market_data to a hypertable
	_, err := db.pool.Exec(ctx, `
		SELECT create_hypertable('market_data', 'time', if_not_exists => TRUE);
	`)
	if err != nil {
		return err
	}

	return nil
}
