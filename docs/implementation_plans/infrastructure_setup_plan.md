# Infrastructure Setup Plan for Trading Platform

## Overview

This document outlines the infrastructure setup required to support the Trading Platform. The infrastructure will provide the foundation for high-performance, scalable, and reliable trading operations with a focus on data persistence, caching, message queuing, and monitoring.

## Components

The infrastructure setup consists of the following key components:

1. **Database**: PostgreSQL with TimescaleDB extension for time-series data
2. **Caching**: Redis for high-performance data caching
3. **Message Queue**: RabbitMQ for reliable message processing
4. **Monitoring**: Prometheus and Grafana for system monitoring
5. **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana) for centralized logging

## Database Setup

### PostgreSQL with TimescaleDB

PostgreSQL will serve as the primary database with TimescaleDB extension for efficient storage and querying of time-series data such as market data, order history, and trade executions.

#### Schema Design

```sql
-- Users and Authentication
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Broker Configurations
CREATE TABLE broker_configs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    broker_type VARCHAR(20) NOT NULL,
    api_key VARCHAR(100) NOT NULL,
    api_secret VARCHAR(255) NOT NULL,
    additional_params JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    broker_id INTEGER REFERENCES broker_configs(id),
    order_id VARCHAR(50) NOT NULL,
    exchange_order_id VARCHAR(50),
    exchange_segment VARCHAR(20) NOT NULL,
    trading_symbol VARCHAR(50) NOT NULL,
    order_side VARCHAR(10) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    order_quantity INTEGER NOT NULL,
    filled_quantity INTEGER DEFAULT 0,
    remaining_quantity INTEGER,
    limit_price DECIMAL(18, 2),
    stop_price DECIMAL(18, 2),
    order_status VARCHAR(20) NOT NULL,
    rejection_reason TEXT,
    order_timestamp TIMESTAMP WITH TIME ZONE,
    last_update_timestamp TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for orders
SELECT create_hypertable('orders', 'order_timestamp');

-- Trades
CREATE TABLE trades (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    user_id INTEGER REFERENCES users(id),
    broker_id INTEGER REFERENCES broker_configs(id),
    trade_id VARCHAR(50) NOT NULL,
    exchange_trade_id VARCHAR(50),
    exchange_segment VARCHAR(20) NOT NULL,
    trading_symbol VARCHAR(50) NOT NULL,
    trade_side VARCHAR(10) NOT NULL,
    trade_quantity INTEGER NOT NULL,
    trade_price DECIMAL(18, 2) NOT NULL,
    trade_timestamp TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for trades
SELECT create_hypertable('trades', 'trade_timestamp');

-- Positions
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    broker_id INTEGER REFERENCES broker_configs(id),
    exchange_segment VARCHAR(20) NOT NULL,
    trading_symbol VARCHAR(50) NOT NULL,
    product_type VARCHAR(20) NOT NULL,
    quantity INTEGER NOT NULL,
    buy_quantity INTEGER NOT NULL,
    sell_quantity INTEGER NOT NULL,
    net_quantity INTEGER NOT NULL,
    average_price DECIMAL(18, 2) NOT NULL,
    last_price DECIMAL(18, 2),
    realized_profit DECIMAL(18, 2),
    unrealized_profit DECIMAL(18, 2),
    position_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for positions
SELECT create_hypertable('positions', 'position_date');

-- Market Data
CREATE TABLE market_data (
    id SERIAL PRIMARY KEY,
    exchange_segment VARCHAR(20) NOT NULL,
    trading_symbol VARCHAR(50) NOT NULL,
    last_price DECIMAL(18, 2) NOT NULL,
    open_price DECIMAL(18, 2),
    high_price DECIMAL(18, 2),
    low_price DECIMAL(18, 2),
    close_price DECIMAL(18, 2),
    volume BIGINT,
    bid_price DECIMAL(18, 2),
    bid_size INTEGER,
    ask_price DECIMAL(18, 2),
    ask_size INTEGER,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for market_data
SELECT create_hypertable('market_data', 'timestamp');

-- Indexes for performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_broker_id ON orders(broker_id);
CREATE INDEX idx_orders_order_status ON orders(order_status);
CREATE INDEX idx_trades_user_id ON trades(user_id);
CREATE INDEX idx_trades_broker_id ON trades(broker_id);
CREATE INDEX idx_positions_user_id ON positions(user_id);
CREATE INDEX idx_positions_broker_id ON positions(broker_id);
CREATE INDEX idx_market_data_symbol ON market_data(exchange_segment, trading_symbol);
```

#### Database Configuration

```ini
# postgresql.conf
max_connections = 200
shared_buffers = 4GB
effective_cache_size = 12GB
maintenance_work_mem = 1GB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 20MB
min_wal_size = 1GB
max_wal_size = 4GB
max_worker_processes = 8
max_parallel_workers_per_gather = 4
max_parallel_workers = 8
```

## Caching Setup

### Redis

Redis will be used for high-performance caching of frequently accessed data such as user sessions, market data, and order status.

#### Redis Configuration

```conf
# redis.conf
maxmemory 4gb
maxmemory-policy allkeys-lru
appendonly yes
appendfsync everysec
```

#### Cache Keys and TTL

| Key Pattern | Description | TTL |
|-------------|-------------|-----|
| `session:{user_id}` | User session data | 24 hours |
| `market_data:{symbol}` | Latest market data for symbol | 5 seconds |
| `order_book:{user_id}` | User's order book | 30 seconds |
| `positions:{user_id}` | User's positions | 30 seconds |
| `holdings:{user_id}` | User's holdings | 30 seconds |

## Message Queue Setup

### RabbitMQ

RabbitMQ will be used for reliable message processing, ensuring that critical operations such as order placement, modification, and cancellation are processed reliably.

#### Queue Configuration

| Queue Name | Description | Durability | Auto-Delete |
|------------|-------------|-----------|-------------|
| `order_placement` | Order placement requests | Yes | No |
| `order_modification` | Order modification requests | Yes | No |
| `order_cancellation` | Order cancellation requests | Yes | No |
| `market_data` | Market data updates | No | No |
| `position_updates` | Position updates | Yes | No |

#### Exchange Configuration

| Exchange Name | Type | Description |
|---------------|------|-------------|
| `trading` | Direct | Trading operations |
| `market_data` | Topic | Market data distribution |

#### Binding Configuration

| Exchange | Queue | Routing Key |
|----------|-------|------------|
| `trading` | `order_placement` | `order.place` |
| `trading` | `order_modification` | `order.modify` |
| `trading` | `order_cancellation` | `order.cancel` |
| `market_data` | `market_data` | `market.#` |

## Monitoring Setup

### Prometheus and Grafana

Prometheus will be used for metrics collection and Grafana for visualization, providing real-time monitoring of system performance and health.

#### Metrics to Monitor

| Metric | Description |
|--------|-------------|
| `http_request_duration_seconds` | HTTP request duration |
| `http_requests_total` | Total HTTP requests |
| `database_query_duration_seconds` | Database query duration |
| `cache_hit_ratio` | Cache hit ratio |
| `message_queue_length` | Message queue length |
| `order_placement_duration_seconds` | Order placement duration |
| `websocket_connections` | Active WebSocket connections |
| `websocket_messages_sent` | WebSocket messages sent |
| `websocket_messages_received` | WebSocket messages received |

#### Alerting Rules

```yaml
groups:
  - name: trading-platform
    rules:
      - alert: HighLatency
        expr: http_request_duration_seconds{quantile="0.95"} > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is above 1 second for 5 minutes"

      - alert: DatabaseErrors
        expr: rate(database_errors_total[5m]) > 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Database errors detected"
          description: "Database errors have been detected in the last 5 minutes"

      - alert: MessageQueueBacklog
        expr: message_queue_length > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Message queue backlog detected"
          description: "Message queue has more than 1000 messages for 5 minutes"
```

## Logging Setup

### ELK Stack

The ELK Stack (Elasticsearch, Logstash, Kibana) will be used for centralized logging, providing a unified view of system logs for troubleshooting and analysis.

#### Log Format

```json
{
  "timestamp": "2023-04-02T10:30:00Z",
  "level": "INFO",
  "service": "order-service",
  "message": "Order placed successfully",
  "order_id": "order123",
  "user_id": "user456",
  "broker_id": "broker789",
  "latency_ms": 50
}
```

#### Log Retention

- Hot data: 7 days
- Warm data: 30 days
- Cold data: 90 days

## Deployment

### Docker Compose

For development and testing, Docker Compose will be used to deploy the infrastructure components.

```yaml
version: '3'

services:
  postgres:
    image: timescale/timescaledb:latest-pg14
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: trading
      POSTGRES_PASSWORD: trading123
      POSTGRES_DB: trading_platform
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: trading
      RABBITMQ_DEFAULT_PASS: trading123
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    restart: unless-stopped

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: admin123
    volumes:
      - grafana_data:/var/lib/grafana
    restart: unless-stopped

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.14.0
    ports:
      - "9200:9200"
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data
    restart: unless-stopped

  logstash:
    image: docker.elastic.co/logstash/logstash:7.14.0
    ports:
      - "5000:5000"
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    restart: unless-stopped

  kibana:
    image: docker.elastic.co/kibana/kibana:7.14.0
    ports:
      - "5601:5601"
    environment:
      ELASTICSEARCH_URL: http://elasticsearch:9200
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  rabbitmq_data:
  prometheus_data:
  grafana_data:
  elasticsearch_data:
```

### Kubernetes

For production, Kubernetes will be used to deploy the infrastructure components with high availability and scalability.

```yaml
# Kubernetes deployment configurations will be provided separately
```

## Implementation Timeline

1. **Week 1**: Set up PostgreSQL with TimescaleDB and create schema
2. **Week 2**: Set up Redis caching and RabbitMQ message queue
3. **Week 3**: Set up Prometheus, Grafana, and ELK Stack
4. **Week 4**: Integrate infrastructure with application code
5. **Week 5**: Performance testing and optimization

## Conclusion

This infrastructure setup plan provides a comprehensive approach to supporting the Trading Platform with high-performance, scalable, and reliable infrastructure components. The setup includes data persistence, caching, message queuing, monitoring, and logging, providing a solid foundation for the trading operations.
