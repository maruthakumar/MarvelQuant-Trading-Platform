# Market Data Service Architecture Design

## Overview
This document outlines the architecture design for the Market Data Service component of the MarverQuant Trading Platform. The Market Data Service is responsible for providing real-time and historical market data, technical indicators, and data streaming capabilities to other components of the platform.

## Architecture Principles
- **Modularity**: Separate concerns into distinct modules for easier maintenance and testing
- **Scalability**: Design for horizontal scaling to handle increasing data volumes
- **Resilience**: Implement fault tolerance and graceful degradation
- **Performance**: Optimize for low latency and high throughput
- **Extensibility**: Make it easy to add new data sources and indicators

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Market Data Service                             │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│  Data       │  Data       │  Technical  │  Data       │  API        │
│  Ingestion  │  Processing │  Indicators │  Storage    │  Layer      │
│  Layer      │  Layer      │  Layer      │  Layer      │             │
└─────┬───────┴──────┬──────┴──────┬──────┴──────┬──────┴──────┬──────┘
      │              │              │             │             │
┌─────▼───────┐┌─────▼───────┐┌─────▼───────┐┌────▼────────┐┌───▼───────┐
│  External   ││  Stream     ││  Indicator  ││  Time-Series││  REST &   │
│  Data       ││  Processing ││  Library    ││  Database   ││  WebSocket│
│  Sources    ││  Engine     ││             ││  & Cache    ││  APIs     │
└─────────────┘└─────────────┘└─────────────┘└─────────────┘└───────────┘
```

## Component Details

### 1. Data Ingestion Layer
The Data Ingestion Layer is responsible for connecting to external data sources and retrieving market data.

#### Components:
- **Data Source Connectors**: Adapters for different data sources (XTS, Yahoo Finance, etc.)
- **Data Source Manager**: Manages connections to data sources and handles failover
- **Rate Limiter**: Ensures compliance with API rate limits
- **Data Validator**: Validates incoming data for accuracy and completeness

#### Interfaces:
- `DataSourceConnector`: Interface for connecting to data sources
- `DataSourceManager`: Interface for managing data sources
- `DataValidator`: Interface for validating data

### 2. Data Processing Layer
The Data Processing Layer is responsible for processing, normalizing, and enriching market data.

#### Components:
- **Data Normalizer**: Normalizes data from different sources into a consistent format
- **Data Enricher**: Enriches data with additional information
- **Data Aggregator**: Aggregates data into different time frames
- **Stream Processor**: Processes real-time data streams

#### Interfaces:
- `DataProcessor`: Interface for processing data
- `DataNormalizer`: Interface for normalizing data
- `DataEnricher`: Interface for enriching data
- `DataAggregator`: Interface for aggregating data

### 3. Technical Indicators Layer
The Technical Indicators Layer is responsible for calculating technical indicators.

#### Components:
- **Indicator Calculator**: Calculates technical indicators
- **Indicator Registry**: Registry of available indicators
- **Indicator Factory**: Factory for creating indicator calculators

#### Interfaces:
- `IndicatorCalculator`: Interface for calculating indicators
- `IndicatorRegistry`: Interface for registering and retrieving indicators
- `IndicatorFactory`: Interface for creating indicator calculators

### 4. Data Storage Layer
The Data Storage Layer is responsible for storing and retrieving market data.

#### Components:
- **Time-Series Database**: Stores historical market data
- **Cache Manager**: Manages data caching
- **Data Access Object**: Provides access to stored data
- **Query Builder**: Builds queries for data retrieval

#### Interfaces:
- `DataStorage`: Interface for storing and retrieving data
- `CacheManager`: Interface for managing cache
- `DataAccessObject`: Interface for accessing data
- `QueryBuilder`: Interface for building queries

### 5. API Layer
The API Layer is responsible for exposing market data to other components.

#### Components:
- **REST API**: Provides REST endpoints for data retrieval
- **WebSocket Server**: Provides WebSocket endpoints for real-time data
- **API Gateway**: Routes API requests to appropriate handlers
- **Authentication/Authorization**: Secures API access

#### Interfaces:
- `ApiHandler`: Interface for handling API requests
- `WebSocketHandler`: Interface for handling WebSocket connections
- `AuthenticationProvider`: Interface for authenticating requests

## Detailed Component Design

### Data Source Connectors

```go
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

// XTSConnector implements DataSourceConnector for XTS API
type XTSConnector struct {
    apiKey      string
    secretKey   string
    source      string
    userID      string
    isConnected bool
    mutex       sync.RWMutex
}

// YahooFinanceConnector implements DataSourceConnector for Yahoo Finance API
type YahooFinanceConnector struct {
    apiKey      string
    isConnected bool
    mutex       sync.RWMutex
}
```

### Data Processing Components

```go
// DataProcessor interface for processing data
type DataProcessor interface {
    Process(data interface{}) (interface{}, error)
}

// DataNormalizer implements DataProcessor for normalizing data
type DataNormalizer struct {
    sourceFormat string
    targetFormat string
}

// DataEnricher implements DataProcessor for enriching data
type DataEnricher struct {
    enrichmentRules []EnrichmentRule
}

// DataAggregator implements DataProcessor for aggregating data
type DataAggregator struct {
    timeFrame string
}
```

### Technical Indicators

```go
// IndicatorCalculator interface for calculating indicators
type IndicatorCalculator interface {
    Calculate(data []OHLCV) ([]IndicatorValue, error)
    GetParameters() map[string]interface{}
    SetParameters(params map[string]interface{}) error
}

// MovingAverageCalculator implements IndicatorCalculator for moving averages
type MovingAverageCalculator struct {
    period int
    type   string // "simple", "exponential", "weighted"
}

// RSICalculator implements IndicatorCalculator for RSI
type RSICalculator struct {
    period int
}

// MACDCalculator implements IndicatorCalculator for MACD
type MACDCalculator struct {
    fastPeriod   int
    slowPeriod   int
    signalPeriod int
}
```

### Data Storage Components

```go
// DataStorage interface for storing and retrieving data
type DataStorage interface {
    StoreMarketData(ctx context.Context, data MarketData) error
    StoreOHLCV(ctx context.Context, symbol string, interval string, data []OHLCV) error
    GetMarketData(ctx context.Context, symbol string) (MarketData, error)
    GetOHLCV(ctx context.Context, symbol string, interval string, from, to time.Time) ([]OHLCV, error)
    GetLatestOHLCV(ctx context.Context, symbol string, interval string, limit int) ([]OHLCV, error)
}

// TimescaleDBStorage implements DataStorage for TimescaleDB
type TimescaleDBStorage struct {
    db *sql.DB
}

// CacheManager interface for managing cache
type CacheManager interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, expiration time.Duration) error
    Delete(key string) error
    Clear() error
}

// RedisCacheManager implements CacheManager for Redis
type RedisCacheManager struct {
    client *redis.Client
}
```

### API Components

```go
// ApiHandler interface for handling API requests
type ApiHandler interface {
    HandleRequest(w http.ResponseWriter, r *http.Request)
}

// MarketDataHandler implements ApiHandler for market data
type MarketDataHandler struct {
    dataStorage DataStorage
    cacheManager CacheManager
}

// HistoricalDataHandler implements ApiHandler for historical data
type HistoricalDataHandler struct {
    dataStorage DataStorage
    cacheManager CacheManager
}

// WebSocketHandler interface for handling WebSocket connections
type WebSocketHandler interface {
    HandleConnection(conn *websocket.Conn)
}

// MarketDataWebSocketHandler implements WebSocketHandler for market data
type MarketDataWebSocketHandler struct {
    dataSourceManager DataSourceManager
    subscriptionManager SubscriptionManager
}
```

## Data Models

### Market Data

```go
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
```

## API Endpoints

### REST API

```
GET /api/v1/market-data/symbols
GET /api/v1/market-data/quote/{symbol}
GET /api/v1/market-data/quotes?symbols=AAPL,MSFT,GOOG
GET /api/v1/market-data/historical/{symbol}?interval=1d&from=2023-01-01&to=2023-12-31
GET /api/v1/market-data/indicators/{indicator}/{symbol}?interval=1d&from=2023-01-01&to=2023-12-31&period=14
```

### WebSocket API

```
ws://host/api/v1/market-data/stream
```

Message format:
```json
{
  "action": "subscribe",
  "symbols": ["AAPL", "MSFT", "GOOG"],
  "fields": ["lastPrice", "bidPrice", "askPrice"]
}
```

## Database Schema

### TimescaleDB Tables

```sql
-- Market data table
CREATE TABLE market_data (
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
);

-- Convert to hypertable
SELECT create_hypertable('market_data', 'timestamp');

-- OHLCV table
CREATE TABLE ohlcv (
    symbol TEXT NOT NULL,
    interval TEXT NOT NULL,
    open DOUBLE PRECISION NOT NULL,
    high DOUBLE PRECISION NOT NULL,
    low DOUBLE PRECISION NOT NULL,
    close DOUBLE PRECISION NOT NULL,
    volume INTEGER NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (symbol, interval, timestamp)
);

-- Convert to hypertable
SELECT create_hypertable('ohlcv', 'timestamp');

-- Indicator values table
CREATE TABLE indicator_values (
    symbol TEXT NOT NULL,
    indicator TEXT NOT NULL,
    value DOUBLE PRECISION,
    values JSONB,
    timestamp TIMESTAMPTZ NOT NULL,
    metadata JSONB,
    PRIMARY KEY (symbol, indicator, timestamp)
);

-- Convert to hypertable
SELECT create_hypertable('indicator_values', 'timestamp');
```

## Caching Strategy

1. **Market Data Cache**:
   - Cache latest market data for frequently accessed symbols
   - TTL: 5 seconds
   - Key format: `market_data:{symbol}`

2. **OHLCV Cache**:
   - Cache historical OHLCV data for frequently accessed symbols and intervals
   - TTL: 5 minutes
   - Key format: `ohlcv:{symbol}:{interval}:{from}:{to}`

3. **Indicator Cache**:
   - Cache calculated indicator values
   - TTL: 5 minutes
   - Key format: `indicator:{indicator}:{symbol}:{interval}:{from}:{to}:{params}`

## Error Handling

1. **Data Source Errors**:
   - Retry with exponential backoff
   - Failover to backup data sources
   - Log errors and notify administrators

2. **Processing Errors**:
   - Log errors and continue processing
   - Return partial results when possible
   - Implement circuit breaker pattern for repeated errors

3. **Storage Errors**:
   - Retry with exponential backoff
   - Use write-ahead logging for durability
   - Implement data recovery mechanisms

4. **API Errors**:
   - Return appropriate HTTP status codes
   - Provide detailed error messages
   - Log errors for troubleshooting

## Monitoring and Metrics

1. **System Metrics**:
   - CPU, memory, disk usage
   - Network I/O
   - Database connections and queries

2. **Application Metrics**:
   - Request rate and latency
   - Error rate
   - Cache hit/miss ratio
   - Data processing throughput

3. **Business Metrics**:
   - Number of active symbols
   - Data freshness
   - Subscription count
   - API usage by endpoint

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Kubernetes Cluster                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │  Market     │  │  Market     │  │  Market     │  │  Market     │ │
│  │  Data API   │  │  Data       │  │  Data       │  │  Data       │ │
│  │  Service    │  │  Ingestion  │  │  Processing │  │  WebSocket  │ │
│  │  (3+ pods)  │  │  Service    │  │  Service    │  │  Service    │ │
│  │             │  │  (2+ pods)  │  │  (2+ pods)  │  │  (2+ pods)  │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │
│                                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │  TimescaleDB│  │  Redis      │  │  NATS       │  │  Prometheus │ │
│  │  Cluster    │  │  Cluster    │  │  Streaming  │  │  &          │ │
│  │  (3+ nodes) │  │  (3+ nodes) │  │  (3+ nodes) │  │  Grafana    │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Scalability Considerations

1. **Horizontal Scaling**:
   - Scale API and processing services based on load
   - Use Kubernetes Horizontal Pod Autoscaler
   - Implement sharding for database if needed

2. **Vertical Scaling**:
   - Optimize code for performance
   - Use efficient data structures and algorithms
   - Implement caching at multiple levels

3. **Load Balancing**:
   - Use Kubernetes Service for load balancing
   - Implement client-side load balancing for internal services
   - Use connection pooling for database connections

## Security Considerations

1. **Authentication and Authorization**:
   - Use JWT for API authentication
   - Implement role-b
(Content truncated due to size limit. Use line ranges to read in chunks)