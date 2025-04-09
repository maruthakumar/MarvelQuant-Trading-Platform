# Platform Integration Implementation Plan

## Overview

This document outlines the implementation plan for integrating all components of the Trading Platform into a cohesive system. The integration will bring together the WebSocket implementation, infrastructure components, order execution engine, risk management system, and position and portfolio management system to create a complete trading platform.

## Architecture

The integrated platform follows a layered architecture:

1. **Core Services Layer**: Provides fundamental services used by all components
2. **Integration Layer**: Connects different components and manages data flow
3. **API Layer**: Exposes unified APIs for client applications
4. **WebSocket Layer**: Provides real-time updates to clients
5. **Monitoring Layer**: Monitors system health and performance

## Implementation Components

### 1. Core Services Layer

```go
// core_services.go
package core

import (
    "github.com/trading-platform/backend/internal/database"
    "github.com/trading-platform/backend/internal/cache"
    "github.com/trading-platform/backend/internal/messaging"
    "github.com/trading-platform/backend/internal/logging"
)

// CoreServices provides fundamental services used by all components
type CoreServices struct {
    DB          *database.Database
    Cache       *cache.Cache
    MessageBus  *messaging.MessageBus
    Logger      *logging.Logger
    Config      *Config
}

// NewCoreServices creates a new core services instance
func NewCoreServices() (*CoreServices, error) {
    // Load configuration
    config, err := LoadConfig()
    if err != nil {
        return nil, err
    }
    
    // Initialize logger
    logger := logging.NewLogger("core_services")
    
    // Initialize database
    db, err := database.NewDatabase(config.Database)
    if err != nil {
        logger.Error("Failed to initialize database", "error", err)
        return nil, err
    }
    
    // Initialize cache
    cache, err := cache.NewCache(config.Cache)
    if err != nil {
        logger.Error("Failed to initialize cache", "error", err)
        return nil, err
    }
    
    // Initialize message bus
    messageBus, err := messaging.NewMessageBus(config.Messaging)
    if err != nil {
        logger.Error("Failed to initialize message bus", "error", err)
        return nil, err
    }
    
    return &CoreServices{
        DB:         db,
        Cache:      cache,
        MessageBus: messageBus,
        Logger:     logger,
        Config:     config,
    }, nil
}

// Config represents the configuration for core services
type Config struct {
    Database  *database.Config
    Cache     *cache.Config
    Messaging *messaging.Config
}

// LoadConfig loads the configuration from environment variables or config file
func LoadConfig() (*Config, error) {
    // Implementation
    return &Config{}, nil
}
```

### 2. Integration Layer

```go
// integration_layer.go
package integration

import (
    "github.com/trading-platform/backend/internal/core"
    "github.com/trading-platform/backend/internal/broker"
    "github.com/trading-platform/backend/internal/websocket"
    "github.com/trading-platform/backend/internal/execution"
    "github.com/trading-platform/backend/internal/risk"
    "github.com/trading-platform/backend/internal/portfolio"
)

// IntegrationLayer connects different components and manages data flow
type IntegrationLayer struct {
    coreServices      *core.CoreServices
    brokerManager     *broker.BrokerManager
    websocketManager  *websocket.WebSocketManager
    orderProcessor    *execution.OrderProcessor
    riskEngine        *risk.PreTradeRiskEngine
    positionTracker   *portfolio.PositionTracker
    portfolioManager  *portfolio.PortfolioManager
    eventHandlers     map[string][]EventHandler
    logger            *core.Logger
}

// NewIntegrationLayer creates a new integration layer
func NewIntegrationLayer(
    coreServices *core.CoreServices,
    brokerManager *broker.BrokerManager,
    websocketManager *websocket.WebSocketManager,
    orderProcessor *execution.OrderProcessor,
    riskEngine *risk.PreTradeRiskEngine,
    positionTracker *portfolio.PositionTracker,
    portfolioManager *portfolio.PortfolioManager,
) *IntegrationLayer {
    return &IntegrationLayer{
        coreServices:     coreServices,
        brokerManager:    brokerManager,
        websocketManager: websocketManager,
        orderProcessor:   orderProcessor,
        riskEngine:       riskEngine,
        positionTracker:  positionTracker,
        portfolioManager: portfolioManager,
        eventHandlers:    make(map[string][]EventHandler),
        logger:           coreServices.Logger.WithPrefix("integration_layer"),
    }
}

// Initialize initializes the integration layer
func (l *IntegrationLayer) Initialize() error {
    // Register event handlers
    l.registerEventHandlers()
    
    // Subscribe to message bus topics
    if err := l.subscribeToTopics(); err != nil {
        l.logger.Error("Failed to subscribe to topics", "error", err)
        return err
    }
    
    l.logger.Info("Integration layer initialized")
    return nil
}

// registerEventHandlers registers event handlers for different event types
func (l *IntegrationLayer) registerEventHandlers() {
    // Order events
    l.RegisterEventHandler("order.placed", l.handleOrderPlaced)
    l.RegisterEventHandler("order.executed", l.handleOrderExecuted)
    l.RegisterEventHandler("order.canceled", l.handleOrderCanceled)
    
    // Trade events
    l.RegisterEventHandler("trade.executed", l.handleTradeExecuted)
    
    // Position events
    l.RegisterEventHandler("position.updated", l.handlePositionUpdated)
    
    // Market data events
    l.RegisterEventHandler("market.data.updated", l.handleMarketDataUpdated)
    
    // Risk events
    l.RegisterEventHandler("risk.limit.exceeded", l.handleRiskLimitExceeded)
    
    // Portfolio events
    l.RegisterEventHandler("portfolio.updated", l.handlePortfolioUpdated)
}

// subscribeToTopics subscribes to message bus topics
func (l *IntegrationLayer) subscribeToTopics() error {
    topics := []string{
        "order.placed",
        "order.executed",
        "order.canceled",
        "trade.executed",
        "position.updated",
        "market.data.updated",
        "risk.limit.exceeded",
        "portfolio.updated",
    }
    
    for _, topic := range topics {
        if err := l.coreServices.MessageBus.Subscribe(topic, l.handleMessage); err != nil {
            return err
        }
    }
    
    return nil
}

// handleMessage handles messages from the message bus
func (l *IntegrationLayer) handleMessage(topic string, message []byte) {
    // Parse message
    event, err := ParseEvent(topic, message)
    if err != nil {
        l.logger.Error("Failed to parse event", "error", err, "topic", topic)
        return
    }
    
    // Handle event
    l.HandleEvent(event)
}

// Event represents an event in the system
type Event struct {
    Type    string
    Payload interface{}
}

// EventHandler represents an event handler function
type EventHandler func(event *Event)

// RegisterEventHandler registers an event handler for an event type
func (l *IntegrationLayer) RegisterEventHandler(eventType string, handler EventHandler) {
    l.eventHandlers[eventType] = append(l.eventHandlers[eventType], handler)
}

// HandleEvent handles an event
func (l *IntegrationLayer) HandleEvent(event *Event) {
    handlers, ok := l.eventHandlers[event.Type]
    if !ok {
        l.logger.Warn("No handlers registered for event type", "type", event.Type)
        return
    }
    
    for _, handler := range handlers {
        go handler(event)
    }
}

// ParseEvent parses an event from a message
func ParseEvent(topic string, message []byte) (*Event, error) {
    // Implementation
    return &Event{}, nil
}

// Event handlers

func (l *IntegrationLayer) handleOrderPlaced(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handleOrderExecuted(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handleOrderCanceled(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handleTradeExecuted(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handlePositionUpdated(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handleMarketDataUpdated(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handleRiskLimitExceeded(event *Event) {
    // Implementation
}

func (l *IntegrationLayer) handlePortfolioUpdated(event *Event) {
    // Implementation
}
```

### 3. API Layer

```go
// api_layer.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/trading-platform/backend/internal/core"
    "github.com/trading-platform/backend/internal/broker"
    "github.com/trading-platform/backend/internal/execution"
    "github.com/trading-platform/backend/internal/risk"
    "github.com/trading-platform/backend/internal/portfolio"
)

// APILayer exposes unified APIs for client applications
type APILayer struct {
    coreServices     *core.CoreServices
    brokerController *BrokerController
    executionController *ExecutionController
    riskController   *RiskController
    portfolioController *PortfolioController
    router           *gin.Engine
    logger           *core.Logger
}

// NewAPILayer creates a new API layer
func NewAPILayer(
    coreServices *core.CoreServices,
    brokerManager *broker.BrokerManager,
    orderProcessor *execution.OrderProcessor,
    riskEngine *risk.PreTradeRiskEngine,
    positionTracker *portfolio.PositionTracker,
    portfolioManager *portfolio.PortfolioManager,
) *APILayer {
    // Create router
    router := gin.Default()
    
    // Create controllers
    brokerController := NewBrokerController(brokerManager)
    executionController := NewExecutionController(orderProcessor)
    riskController := NewRiskController(riskEngine)
    portfolioController := NewPortfolioController(positionTracker, portfolioManager)
    
    return &APILayer{
        coreServices:       coreServices,
        brokerController:    brokerController,
        executionController: executionController,
        riskController:      riskController,
        portfolioController: portfolioController,
        router:              router,
        logger:              coreServices.Logger.WithPrefix("api_layer"),
    }
}

// Initialize initializes the API layer
func (l *APILayer) Initialize() error {
    // Set up middleware
    l.setupMiddleware()
    
    // Register routes
    l.registerRoutes()
    
    l.logger.Info("API layer initialized")
    return nil
}

// setupMiddleware sets up middleware for the router
func (l *APILayer) setupMiddleware() {
    // CORS middleware
    l.router.Use(cors())
    
    // Authentication middleware
    l.router.Use(authenticate(l.coreServices))
    
    // Logging middleware
    l.router.Use(logRequest(l.logger))
    
    // Recovery middleware
    l.router.Use(gin.Recovery())
}

// registerRoutes registers routes for the router
func (l *APILayer) registerRoutes() {
    // API version group
    v1 := l.router.Group("/api/v1")
    
    // Register controller routes
    l.brokerController.RegisterRoutes(v1)
    l.executionController.RegisterRoutes(v1)
    l.riskController.RegisterRoutes(v1)
    l.portfolioController.RegisterRoutes(v1)
    
    // Health check route
    l.router.GET("/health", l.healthCheck)
}

// Start starts the API server
func (l *APILayer) Start() error {
    return l.router.Run(l.coreServices.Config.API.Address)
}

// healthCheck handles health check requests
func (l *APILayer) healthCheck(ctx *gin.Context) {
    ctx.JSON(200, gin.H{
        "status": "ok",
    })
}

// Middleware functions

func cors() gin.HandlerFunc {
    // Implementation
    return func(ctx *gin.Context) {
        // Implementation
    }
}

func authenticate(coreServices *core.CoreServices) gin.HandlerFunc {
    // Implementation
    return func(ctx *gin.Context) {
        // Implementation
    }
}

func logRequest(logger *core.Logger) gin.HandlerFunc {
    // Implementation
    return func(ctx *gin.Context) {
        // Implementation
    }
}
```

### 4. WebSocket Layer

```go
// websocket_layer.go
package websocket

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/trading-platform/backend/internal/core"
    "github.com/trading-platform/backend/internal/broker"
    "github.com/trading-platform/backend/internal/execution"
    "github.com/trading-platform/backend/internal/portfolio"
)

// WebSocketLayer provides real-time updates to clients
type WebSocketLayer struct {
    coreServices     *core.CoreServices
    brokerManager    *broker.BrokerManager
    orderProcessor   *execution.OrderProcessor
    positionTracker  *portfolio.PositionTracker
    portfolioManager *portfolio.PortfolioManager
    router           *gin.Engine
    upgrader         websocket.Upgrader
    clients          map[string]*Client
    logger           *core.Logger
}

// NewWebSocketLayer creates a new WebSocket layer
func NewWebSocketLayer(
    coreServices *core.CoreServices,
    brokerManager *broker.BrokerManager,
    orderProcessor *execution.OrderProcessor,
    positionTracker *portfolio.PositionTracker,
    portfolioManager *portfolio.PortfolioManager,
    router *gin.Engine,
) *WebSocketLayer {
    return &WebSocketLayer{
        coreServices:     coreServices,
        brokerManager:    brokerManager,
        orderProcessor:   orderProcessor,
        positionTracker:  positionTracker,
        portfolioManager: portfolioManager,
        router:           router,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
        clients: make(map[string]*Client),
        logger:  coreServices.Logger.WithPrefix("websocket_layer"),
    }
}

// Initialize initializes the WebSocket layer
func (l *WebSocketLayer) Initialize() error {
    // Register WebSocket routes
    l.registerRoutes()
    
    // Subscribe to message bus topics
    if err := l.subscribeToTopics(); err != nil {
        l.logger.Error("Failed to subscribe to topics", "error", err)
        return err
    }
    
    l.logger.Info("WebSocket layer initialized")
    return nil
}

// registerRoutes registers WebSocket routes
func (l *WebSocketLayer) registerRoutes() {
    l.router.GET("/ws", l.handleConnection)
}

// subscribeToTopics subscribes to message bus topics
func (l *WebSocketLayer) subscribeToTopics() error {
    topics := []string{
        "order.placed",
        "order.executed",
        "order.canceled",
        "trade.executed",
        "position.updated",
        "market.data.updated",
        "portfolio.updated",
    }
    
    for _, topic := range topics {
        if err := l.coreServices.MessageBus.Subscribe(topic, l.handleMessage); err != nil {
            return err
        }
    }
    
    return nil
}

// handleConnection handles WebSocket connections
func (l *WebSocketLayer) handleConnection(ctx *gin.Context) {
    // Upgrade HTTP connection to WebSocket
    conn, err := l.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        l.logger.Error("Failed to upgrade connection", "error", err)
        return
    }
    
    // Get user ID from context
    userID := ctx.GetString("userID")
    if userID == "" {
        l.logger.Error("User ID not found in context")
        conn.Close()
        return
    }
    
    // Create client
    client := NewClient(userID, conn, l.coreServices.Logger)
    l.clients[userID] = client
    
    // Start client
    go client.Start()
    
    l.logger.Info("Client connected", "user_id", userID)
}

// handleMessage handles messages from the message bus
func (l *WebSocketLayer) handleMessage(topic string, message []byte) {
    // Parse message
    event, err := ParseEvent(topic, message)
    if err != nil {
        l.logger.Error("Failed to parse event", "error", err, "topic", topic)
        return
    }
    
    // Get user ID from event
    userID, ok := getUserIDFromEvent(event)
    if !ok {
        l.logger.Error("User ID not found in event", "topic", topic)
        return
    }
    
    // Get client
    client, ok := l.clients[userID]
    if !ok {
        // Client not connected
        return
    }
    
    // Send event to client
    client.Send(event)
}

// Client represents a WebSocket client
type Client struct {
    UserID string
    conn   *websocket.Conn
    send   chan interface{}
    logger *core.Logger
}

// NewClient creates a new WebSocket client
func NewClient(userID string, conn *websocket.Conn, logger *core.Logger) *Client {
    return &Client{
        UserID: userID,
        conn:   conn,
        send:   make(chan interface{}, 256),
        logger: logger.WithPrefix("websocket_client"),
    }
}

// Start starts the client
func (c *Client) Start() {
    // Start writer goroutine
    go c.writer()
    
    // Start reader goroutine
    go c.reader()
}

// Send sends a message to the client
func (c *Client) Send(message interface{}) {
    c.send <- message
}

// writer writes messages to the WebSocket connection
func (c *Client) writer() {
    // Implementation
}

// reader reads messages from the WebSocket connection
func (c *Client) reader() {
    // Implementation
}

// ParseEvent parses an event from a message
func ParseEvent(topic string, message []byte) (interface{}, error) {
    // Implementation
    return nil, nil
}

// getUserIDFromEvent gets the user ID from an event
func getUserIDFromEvent(event interface{}) (string, bool) {
    // Implementation
    return "", false
}
```

### 5. Monitoring Layer

```go
// monitoring_layer.go
package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gin-gonic/gin"
    "github.com/trading-platform/backend/internal/core"
)

// MonitoringLayer monitors system health and performance
type MonitoringLayer struct {
    coreServices *core.CoreServices
    router       *gin.Engine
    registry     *prometheus.Registry
    metrics      *Metrics
    logger       *core.Logger
}

// NewMonitoringLayer creates a new monitoring layer
func NewMonitoringLayer(coreServices *core.CoreServices, router *gin.Engine) *MonitoringLayer {
    // Create registry
    registry := prometheus.NewRegistry()
    
    // Create metrics
    metrics := NewMetrics(registry)
    
    return &MonitoringLayer{
        coreServices: coreServices,
        router:       router,
        registry:     registry,
        metrics:      metrics,
        logger:       coreServices.Logger.WithPrefix("monitoring_layer"),
    }
}

// Initialize initializes the monitoring layer
func (l *MonitoringLayer) Initialize() error {
    // Register metrics handler
    l.router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(l.registry, promhttp.HandlerOpts{})))
    
    // Register health check handler
    l.router.GET("/health", l.healthCheck)
    
    // Start metrics collection
    go l.collectMetrics()
    
    l.logger.Info("Monitoring layer initialized")
    return nil
}

// healthCheck handles health check requests
func (l *MonitoringLayer) healthCheck(ctx *gin.Context) {
    // Check database health
    dbHealth := l.checkDatabaseHealth()
    
    // Check cache health
    cacheHealth := l.checkCacheHealth()
    
    // Check message bus health
    messageBusHealth := l.checkMessageBusHealth()
    
    // Overall health
    overallHealth := dbHealth && cacheHealth && messageBusHealth
    
    if overallHealth {
        ctx.JSON(200, gin.H{
            "status": "ok",
            "components": gin.H{
                "database":    dbHealth,
                "cache":       cacheHealth,
                "message_bus": messageBusHealth,
            },
        })
    } else {
        ctx.JSON(500, gin.H{
            "status": "error",
            "components": gin.H{
                "database":    dbHealth,
                "cache":       cacheHealth,
                "message_bus": messageBusHealth,
            },
        })
    }
}

// checkDatabaseHealth checks the health of the database
func (l *MonitoringLayer) checkDatabaseHealth() bool {
    // Implementation
    return true
}

// checkCacheHealth checks the health of the cache
func (l *MonitoringLayer) checkCacheHealth() bool {
    // Implementation
    return true
}

// checkMessageBusHealth checks the health of the message bus
func (l *MonitoringLayer) checkMessageBusHealth() bool {
    // Implementation
    return true
}

// collectMetrics collects metrics
func (l *MonitoringLayer) collectMetrics() {
    // Implementation
}

// Metrics represents the metrics collected by the monitoring layer
type Metrics struct {
    // HTTP metrics
    HTTPRequestsTotal      *prometheus.CounterVec
    HTTPRequestDuration    *prometheus.HistogramVec
    
    // Database metrics
    DatabaseQueryTotal     *prometheus.CounterVec
    DatabaseQueryDuration  *prometheus.HistogramVec
    
    // Cache metrics
    CacheHitTotal          *prometheus.CounterVec
    CacheMissTotal         *prometheus.CounterVec
    
    // Message bus metrics
    MessageBusPublishTotal *prometheus.CounterVec
    MessageBusConsumeTotal *prometheus.CounterVec
    
    // Order metrics
    OrderPlacedTotal       *prometheus.CounterVec
    OrderExecutedTotal     *prometheus.CounterVec
    OrderCanceledTotal     *prometheus.CounterVec
    
    // Trade metrics
    TradeExecutedTotal     *prometheus.CounterVec
    
    // Position metrics
    PositionUpdatedTotal   *prometheus.CounterVec
    
    // Portfolio metrics
    PortfolioUpdatedTotal  *prometheus.CounterVec
    
    // WebSocket metrics
    WebSocketConnectionsTotal *prometheus.CounterVec
    WebSocketMessagesSentTotal *prometheus.CounterVec
    WebSocketMessagesReceivedTotal *prometheus.CounterVec
}

// NewMetrics creates a new metrics instance
func NewMetrics(registry *prometheus.Registry) *Metrics {
    // Implementation
    return &Metrics{}
}
```

### 6. Main Application

```go
// main.go
package main

import (
    "github.com/trading-platform/backend/internal/core"
    "github.com/trading-platform/backend/internal/broker"
    "github.com/trading-platform/backend/internal/websocket"
    "github.com/trading-platform/backend/internal/execution"
    "github.com/trading-platform/backend/internal/risk"
    "github.com/trading-platform/backend/internal/portfolio"
    "github.com/trading-platform/backend/internal/integration"
    "github.com/trading-platform/backend/internal/api"
    "github.com/trading-platform/backend/internal/monitoring"
)

func main() {
    // Initialize core services
    coreServices, err := core.NewCoreServices()
    if err != nil {
        panic(err)
    }
    
    // Initialize broker manager
    brokerManager, err := broker.NewBrokerManager(coreServices)
    if err != nil {
        panic(err)
    }
    
    // Initialize WebSocket manager
    websocketManager, err := websocket.NewWebSocketManager(coreServices)
    if err != nil {
        panic(err)
    }
    
    // Initialize position tracker
    positionTracker, err := portfolio.NewPositionTracker(coreServices, brokerManager)
    if err != nil {
        panic(err)
    }
    
    // Initialize portfolio manager
    portfolioManager, err := portfolio.NewPortfolioManager(coreServices, positionTracker)
    if err != nil {
        panic(err)
    }
    
    // Initialize risk engine
    riskEngine, err := risk.NewPreTradeRiskEngine(coreServices, positionTracker)
    if err != nil {
        panic(err)
    }
    
    // Initialize order processor
    orderProcessor, err := execution.NewOrderProcessor(coreServices, brokerManager, riskEngine)
    if err != nil {
        panic(err)
    }
    
    // Initialize integration layer
    integrationLayer, err := integration.NewIntegrationLayer(
        coreServices,
        brokerManager,
        websocketManager,
        orderProcessor,
        riskEngine,
        positionTracker,
        portfolioManager,
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize API layer
    apiLayer, err := api.NewAPILayer(
        coreServices,
        brokerManager,
        orderProcessor,
        riskEngine,
        positionTracker,
        portfolioManager,
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize WebSocket layer
    websocketLayer, err := websocket.NewWebSocketLayer(
        coreServices,
        brokerManager,
        orderProcessor,
        positionTracker,
        portfolioManager,
        apiLayer.GetRouter(),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize monitoring layer
    monitoringLayer, err := monitoring.NewMonitoringLayer(
        coreServices,
        apiLayer.GetRouter(),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize all layers
    if err := integrationLayer.Initialize(); err != nil {
        panic(err)
    }
    
    if err := apiLayer.Initialize(); err != nil {
        panic(err)
    }
    
    if err := websocketLayer.Initialize(); err != nil {
        panic(err)
    }
    
    if err := monitoringLayer.Initialize(); err != nil {
        panic(err)
    }
    
    // Start API server
    if err := apiLayer.Start(); err != nil {
        panic(err)
    }
}
```

## Database Schema Integration

The integrated platform will use a unified database schema that supports all components:

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

-- Risk Parameters
CREATE TABLE risk_parameters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    max_order_size INTEGER NOT NULL,
    max_position_size INTEGER NOT NULL,
    max_order_value DECIMAL(18, 2) NOT NULL,
    max_position_value DECIMAL(18, 2) NOT NULL,
    max_loss DECIMAL(18, 2) NOT NULL,
    price_range_percent DECIMAL(5, 2) NOT NULL,
    max_daily_loss DECIMAL(18, 2) NOT NULL,
    max_daily_trades INTEGER NOT NULL,
    max_daily_volume INTEGER NOT NULL,
    enable_circuit_breaker BOOLEAN DEFAULT TRUE,
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

-- Portfolios
CREATE TABLE portfolios (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    total_value DECIMAL(18, 2),
    realized_profit DECIMAL(18, 2),
    unrealized_profit DECIMAL(18, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Portfolio Positions
CREATE TABLE portfolio_positions (
    id SERIAL PRIMARY KEY,
    portfolio_id INTEGER REFERENCES portfolios(id),
    position_id INTEGER REFERENCES positions(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

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

-- Reports
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    content JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_broker_id ON orders(broker_id);
CREATE INDEX idx_orders_order_status ON orders(order_status);
CREATE INDEX idx_trades_user_id ON trades(user_id);
CREATE INDEX idx_trades_broker_id ON trades(broker_id);
CREATE INDEX idx_positions_user_id ON positions(user_id);
CREATE INDEX idx_positions_broker_id ON positions(broker_id);
CREATE INDEX idx_market_data_symbol ON market_data(exchange_segment, trading_symbol);
CREATE INDEX idx_portfolios_user_id ON portfolios(user_id);
CREATE INDEX idx_portfolio_positions_portfolio_id ON portfolio_positions(portfolio_id);
CREATE INDEX idx_reports_user_id ON reports(user_id);
CREATE INDEX idx_reports_type ON reports(type);
```

## Message Bus Topics

The integrated platform will use a unified message bus with the following topics:

| Topic | Description |
|-------|-------------|
| `order.placed` | Order placed |
| `order.executed` | Order executed |
| `order.canceled` | Order canceled |
| `trade.executed` | Trade executed |
| `position.updated` | Position updated |
| `market.data.updated` | Market data updated |
| `risk.limit.exceeded` | Risk limit exceeded |
| `portfolio.updated` | Portfolio updated |
| `report.generated` | Report generated |

## API Endpoints

The integrated platform will expose the following API endpoints:

### Authentication

- `POST /api/v1/auth/login`: Login
- `POST /api/v1/auth/logout`: Logout
- `POST /api/v1/auth/register`: Register
- `GET /api/v1/auth/profile`: Get user profile
- `PUT /api/v1/auth/profile`: Update user profile

### Broker Management

- `GET /api/v1/brokers`: Get all brokers
- `POST /api/v1/brokers`: Register a broker
- `GET /api/v1/brokers/:brokerID`: Get broker details
- `PUT /api/v1/brokers/:brokerID`: Update broker details
- `DELETE /api/v1/brokers/:brokerID`: Delete broker

### Order Execution

- `POST /api/v1/orders`: Place an order
- `PUT /api/v1/orders/:orderID`: Modify an order
- `DELETE /api/v1/orders/:orderID`: Cancel an order
- `GET /api/v1/orders`: Get order book
- `GET /api/v1/orders/:orderID`: Get order details
- `POST /api/v1/strategies`: Execute a strategy

### Risk Management

- `GET /api/v1/risk/parameters`: Get risk parameters
- `PUT /api/v1/risk/parameters`: Update risk parameters
- `GET /api/v1/risk/circuit-breaker`: Get circuit breaker state
- `POST /api/v1/risk/circuit-breaker/reset`: Reset circuit breaker

### Position and Portfolio Management

- `GET /api/v1/positions`: Get all positions
- `GET /api/v1/positions/:exchangeSegment/:tradingSymbol`: Get position details
- `GET /api/v1/positions/:exchangeSegment/:tradingSymbol/performance`: Get position performance
- `GET /api/v1/portfolios`: Get all portfolios
- `POST /api/v1/portfolios`: Create a portfolio
- `GET /api/v1/portfolios/:portfolioID`: Get portfolio details
- `PUT /api/v1/portfolios/:portfolioID`: Update portfolio details
- `DELETE /api/v1/portfolios/:portfolioID`: Delete portfolio
- `POST /api/v1/portfolios/:portfolioID/positions`: Add position to portfolio
- `DELETE /api/v1/portfolios/:portfolioID/positions/:exchangeSegment/:tradingSymbol`: Remove position from portfolio
- `GET /api/v1/portfolios/:portfolioID/performance`: Get portfolio performance
- `GET /api/v1/portfolios/:portfolioID/analytics`: Get portfolio analytics

### Reporting

- `GET /api/v1/reports`: Get all reports
- `GET /api/v1/reports/:reportID`: Get report details
- `POST /api/v1/reports/positions/:exchangeSegment/:tradingSymbol`: Generate position report
- `POST /api/v1/reports/portfolios/:portfolioID`: Generate portfolio report
- `DELETE /api/v1/reports/:reportID`: Delete report

### Market Data

- `GET /api/v1/market-data/:exchangeSegment/:tradingSymbol`: Get market data
- `GET /api/v1/market-data/:exchangeSegment/:tradingSymbol/history`: Get historical market data

### System

- `GET /health`: Health check
- `GET /metrics`: Prometheus metrics

## WebSocket Events

The integrated platform will provide the following WebSocket events:

| Event | Description |
|-------|-------------|
| `order.update` | Order status update |
| `trade.update` | Trade update |
| `position.update` | Position update |
| `market.data.update` | Market data update |
| `portfolio.update` | Portfolio update |
| `risk.alert` | Risk alert |

## Testing Strategy

1. **Unit Tests**: Test individual components in isolation
   - Test core services
   - Test integration layer
   - Test API layer
   - Test WebSocket layer
   - Test monitoring layer

2. **Integration Tests**: Test the integration between components
   - Test order flow from API to broker
   - Test market data flow from broker to WebSocket
   - Test position updates from trades to portfolio

3. **System Tests**: Test the entire system
   - Test end-to-end order placement and execution
   - Test end-to-end portfolio management
   - Test end-to-end risk management

4. **Performance Tests**: Test system performance
   - Test order throughput
   - Test market data throughput
   - Test WebSocket message throughput

## Deployment

The integrated platform will be deployed using Docker and Kubernetes:

```yaml
# docker-compose.yml for development
version: '3'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=trading
      - DB_PASSWORD=trading123
      - DB_NAME=trading_platform
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - RABBITMQ_HOST=rabbitmq
      - RABBITMQ_PORT=5672
      - RABBITMQ_USER=trading
      - RABBITMQ_PASSWORD=trading123
    depends_on:
      - postgres
      - redis
      - rabbitmq

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

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

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

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: admin123
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  postgres_data:
  redis_data:
  rabbitmq_data:
  prometheus_data:
  grafana_data:
```

```yaml
# Kubernetes deployment for production
# Kubernetes deployment configurations will be provided separately
```

## Implementation Timeline

1. **Week 1**: Implement Core Services Layer and Integration Layer
2. **Week 2**: Implement API Layer and WebSocket Layer
3. **Week 3**: Implement Monitoring Layer and Database Schema
4. **Week 4**: Integrate all components and test
5. **Week 5**: Deploy and optimize

## Conclusion

This Platform Integration Implementation Plan provides a comprehensive approach to integrating all components of the Trading Platform into a cohesive system. The integration follows a layered architecture that allows for flexibility and extensibility, with support for real-time updates, monitoring, and reporting.
