# Trading Platform Integration Troubleshooting Guide

## Overview

This guide provides solutions for common integration issues encountered when working with the Trading Platform. It covers issues related to broker integration, WebSocket connections, C++ execution engine integration, and general system integration problems.

## Table of Contents

1. [Authentication Issues](#authentication-issues)
2. [Network and Connection Issues](#network-and-connection-issues)
3. [Order Execution Issues](#order-execution-issues)
4. [WebSocket Integration Issues](#websocket-integration-issues)
5. [C++ Execution Engine Issues](#c-execution-engine-issues)
6. [Multi-Broker Integration Issues](#multi-broker-integration-issues)
7. [Performance Issues](#performance-issues)
8. [Environment Isolation Issues](#environment-isolation-issues)
9. [Debugging Techniques](#debugging-techniques)
10. [Common Error Codes](#common-error-codes)

## Authentication Issues

### Token Expiration

**Problem**: Authentication fails with `ErrTokenExpired` or HTTP 401 Unauthorized.

**Solution**:
1. Implement automatic token refresh:
```go
if errors.IsAuthError(err) && errors.Is(err, errors.ErrTokenExpired) {
    // Refresh token
    session, err := client.Login()
    if err != nil {
        return err
    }
    
    // Retry the original operation with new token
    return client.ExecuteWithToken(session.Token)
}
```

2. Check token expiration before operations:
```go
if time.Now().After(session.ExpiryTime) {
    // Refresh token before proceeding
    session, err = client.Login()
    if err != nil {
        return err
    }
}
```

### Invalid Credentials

**Problem**: Authentication fails with `ErrInvalidCredentials`.

**Solution**:
1. Verify API key and secret key are correct
2. Check if API key has been revoked or disabled
3. Ensure the API key has the necessary permissions
4. Verify the correct environment (UAT vs Production) is being used

```go
// Configuration validation
if config.APIKey == "" {
    return nil, errors.ErrEmptyAPIKey
}
if config.SecretKey == "" {
    return nil, errors.ErrEmptySecretKey
}
```

### Session Invalidation

**Problem**: Session becomes invalid during operations with `ErrSessionInvalid`.

**Solution**:
1. Implement session monitoring:
```go
// Check session validity periodically
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := client.ValidateSession(); err != nil {
                // Refresh session
                client.RefreshSession()
            }
        case <-ctx.Done():
            return
        }
    }
}()
```

2. Handle session errors with retry logic:
```go
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    if errors.IsAuthError(err) {
        // Refresh session before retry
        session, err = client.Login()
        if err != nil {
            return err
        }
    }
    return client.ExecuteOperation()
})
```

## Network and Connection Issues

### Connection Failures

**Problem**: Unable to connect to broker API with `ErrConnectionFailed`.

**Solution**:
1. Check network connectivity:
```bash
# Test connectivity to API endpoint
curl -v https://api.xts.com/health
```

2. Verify firewall and proxy settings:
```go
// Configure HTTP client with proxy if needed
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}
```

3. Implement connection retry with backoff:
```go
// Retry with exponential backoff
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    return client.Connect()
})
```

### Request Timeouts

**Problem**: API requests time out with `ErrRequestTimeout`.

**Solution**:
1. Increase timeout duration:
```go
// Configure longer timeout for specific operations
client.SetTimeout(operation, 30*time.Second)
```

2. Check for network congestion or high latency:
```bash
# Test network latency
ping api.xts.com
```

3. Implement request timeout handling:
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := client.ExecuteWithContext(ctx)
if errors.Is(err, context.DeadlineExceeded) {
    // Handle timeout specifically
    log.Warn("Operation timed out, will retry with longer timeout")
    
    // Retry with longer timeout
    ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    err = client.ExecuteWithContext(ctx)
}
```

### API Rate Limiting

**Problem**: Requests are rejected due to rate limiting with `ErrAPIRateLimited`.

**Solution**:
1. Implement rate limiting on client side:
```go
// Create rate limiter
limiter := rate.NewLimiter(rate.Limit(5), 10) // 5 requests per second, burst of 10

// Use rate limiter before making requests
if err := limiter.Wait(ctx); err != nil {
    return err
}
// Proceed with request
```

2. Implement adaptive rate limiting:
```go
// Adjust rate limit based on response headers
if resp.StatusCode == http.StatusTooManyRequests {
    // Extract rate limit info from headers
    remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
    reset, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
    
    // Adjust limiter based on remaining quota
    if remaining < 10 {
        // Slow down requests
        time.Sleep(time.Duration(reset) * time.Second)
    }
}
```

## Order Execution Issues

### Order Rejection

**Problem**: Orders are rejected with `ErrOrderRejected`.

**Solution**:
1. Check order parameters:
```go
// Validate order parameters before submission
if order.Quantity <= 0 {
    return nil, errors.Wrap(errors.ErrInvalidOrderParams, "ORDER_VALIDATION", "Quantity must be greater than zero")
}
if order.Price <= 0 && order.OrderType == models.OrderTypeLimit {
    return nil, errors.Wrap(errors.ErrInvalidOrderParams, "ORDER_VALIDATION", "Price must be greater than zero for limit orders")
}
```

2. Verify account has sufficient margin:
```go
// Check margin before placing order
margin, err := client.GetMarginRequired(order)
if err != nil {
    return nil, err
}

available, err := client.GetAvailableMargin()
if err != nil {
    return nil, err
}

if available < margin {
    return nil, errors.Wrap(errors.ErrOrderRejected, "INSUFFICIENT_MARGIN", 
        fmt.Sprintf("Required: %f, Available: %f", margin, available))
}
```

3. Check for position limits:
```go
// Check position limits before placing order
currentExposure, err := client.GetCurrentExposure(order.Symbol)
if err != nil {
    return nil, err
}

maxExposure := client.GetMaxExposure(order.Symbol)
if currentExposure + order.Quantity > maxExposure {
    return nil, errors.Wrap(errors.ErrOrderRejected, "POSITION_LIMIT_EXCEEDED", 
        fmt.Sprintf("Max allowed: %f", maxExposure))
}
```

### Order Not Found

**Problem**: Unable to find order with `ErrOrderNotFound`.

**Solution**:
1. Verify order ID is correct:
```go
// Validate order ID format
if !isValidOrderID(orderID) {
    return nil, errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_ORDER_ID", 
        "Order ID format is invalid")
}
```

2. Check if order was actually submitted:
```go
// Check order submission status
if order.Status == models.OrderStatusPending {
    // Order is still being processed
    return nil, errors.Wrap(errors.ErrOrderNotFound, "ORDER_PENDING", 
        "Order is still being processed")
}
```

3. Implement order tracking:
```go
// Track orders in local cache
orderCache := make(map[string]*models.Order)

// When placing order
orderCache[order.OrderID] = order

// When checking order
cachedOrder, exists := orderCache[orderID]
if !exists {
    // Order not in cache, try to fetch from API
    order, err := client.GetOrder(orderID)
    if err != nil {
        if errors.IsOrderError(err) && errors.Is(err, errors.ErrOrderNotFound) {
            // Handle specifically
            return nil, errors.Wrap(err, "ORDER_TRACKING", 
                "Order not found in cache or API")
        }
        return nil, err
    }
    // Update cache
    orderCache[orderID] = order
    return order, nil
}
```

### Order Modification Failures

**Problem**: Unable to modify orders with `ErrOrderModifyFailed`.

**Solution**:
1. Check if order is in modifiable state:
```go
// Check order status before modification
order, err := client.GetOrder(orderID)
if err != nil {
    return err
}

if order.Status != models.OrderStatusOpen && order.Status != models.OrderStatusPending {
    return errors.Wrap(errors.ErrOrderModifyFailed, "INVALID_ORDER_STATE", 
        fmt.Sprintf("Cannot modify order in %s state", order.Status))
}
```

2. Verify modification parameters:
```go
// Validate modification parameters
if modifyOrder.Price <= 0 && order.OrderType == models.OrderTypeLimit {
    return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_PRICE", 
        "Price must be greater than zero for limit orders")
}
```

3. Implement modification retry logic:
```go
// Retry modification with backoff
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    return client.ModifyOrder(modifyOrder)
})
```

## WebSocket Integration Issues

### Connection Failures

**Problem**: WebSocket connection fails with `ErrWebSocketConnFailed`.

**Solution**:
1. Check WebSocket endpoint URL:
```go
// Validate WebSocket URL
if !strings.HasPrefix(wsURL, "wss://") && !strings.HasPrefix(wsURL, "ws://") {
    return nil, errors.Wrap(errors.ErrWebSocketConnFailed, "INVALID_URL", 
        "WebSocket URL must start with ws:// or wss://")
}
```

2. Implement connection retry with backoff:
```go
// Retry WebSocket connection with backoff
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        return errors.Wrap(errors.ErrWebSocketConnFailed, "CONNECTION_FAILED", err.Error())
    }
    client.conn = conn
    return nil
})
```

3. Handle authentication for WebSocket:
```go
// Authenticate WebSocket connection
authMsg := map[string]interface{}{
    "type": "auth",
    "token": session.Token,
}
err = client.conn.WriteJSON(authMsg)
if err != nil {
    return errors.Wrap(errors.ErrWebSocketConnFailed, "AUTH_FAILED", err.Error())
}
```

### Connection Closed Unexpectedly

**Problem**: WebSocket connection closes unexpectedly with `ErrWebSocketClosed`.

**Solution**:
1. Implement automatic reconnection:
```go
// Reconnect when connection closes
go func() {
    for {
        _, _, err := client.conn.ReadMessage()
        if err != nil {
            log.Warn("WebSocket connection closed, reconnecting...", "error", err)
            
            // Reconnect with backoff
            backoff := 1 * time.Second
            maxBackoff := 30 * time.Second
            
            for {
                time.Sleep(backoff)
                
                err := client.Connect()
                if err == nil {
                    // Resubscribe to previous subscriptions
                    client.Resubscribe()
                    break
                }
                
                log.Warn("Reconnection failed, retrying...", "error", err)
                backoff = backoff * 2
                if backoff > maxBackoff {
                    backoff = maxBackoff
                }
            }
        }
    }
}()
```

2. Implement heartbeat to keep connection alive:
```go
// Send heartbeat every 30 seconds
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            err := client.conn.WriteMessage(websocket.PingMessage, []byte{})
            if err != nil {
                log.Warn("Failed to send heartbeat", "error", err)
            }
        case <-client.done:
            return
        }
    }
}()
```

3. Handle connection close gracefully:
```go
// Close connection gracefully
func (client *WebSocketClient) Close() error {
    // Send close message
    err := client.conn.WriteMessage(websocket.CloseMessage, 
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
    if err != nil {
        log.Warn("Failed to send close message", "error", err)
    }
    
    // Close channel to stop goroutines
    close(client.done)
    
    // Close connection
    return client.conn.Close()
}
```

### Subscription Failures

**Problem**: Unable to subscribe to market data with `ErrSubscriptionFailed`.

**Solution**:
1. Verify subscription parameters:
```go
// Validate subscription parameters
if len(symbols) == 0 {
    return nil, errors.Wrap(errors.ErrSubscriptionFailed, "INVALID_SYMBOLS", 
        "Symbol list cannot be empty")
}
```

2. Implement subscription retry logic:
```go
// Retry subscription with backoff
err = recovery.RetryWithBackoff(ctx, recovery.DefaultRetryConfig(), func() error {
    subMsg := map[string]interface{}{
        "type": "subscribe",
        "symbols": symbols,
    }
    return client.conn.WriteJSON(subMsg)
})
```

3. Handle subscription confirmation:
```go
// Wait for subscription confirmation
select {
case resp := <-client.responseChan:
    if resp.Type == "subscribed" {
        // Subscription successful
        return nil
    } else if resp.Type == "error" {
        // Subscription failed
        return errors.Wrap(errors.ErrSubscriptionFailed, "SUBSCRIPTION_REJECTED", 
            resp.Message)
    }
case <-time.After(5 * time.Second):
    // Subscription timed out
    return errors.Wrap(errors.ErrSubscriptionFailed, "SUBSCRIPTION_TIMEOUT", 
        "No confirmation received within timeout")
}
```

## C++ Execution Engine Issues

### Integration Failures

**Problem**: C++ execution engine fails to integrate with Go backend.

**Solution**:
1. Verify shared library loading:
```go
// Load C++ shared library
lib := C.dlopen(C.CString("./libexecution.so"), C.RTLD_LAZY)
if lib == nil {
    return fmt.Errorf("failed to load execution engine library: %s", C.GoString(C.dlerror()))
}
```

2. Check function signatures:
```go
// Get function pointer
initFunc := C.dlsym(lib, C.CString("InitExecutionEngine"))
if initFunc == nil {
    return fmt.Errorf("failed to find InitExecutionEngine function: %s", C.GoString(C.dlerror()))
}
```

3. Implement proper error handling for C++ calls:
```go
// Call C++ function and handle errors
result := C.CallExecutionEngine(C.CString(orderJSON))
if result.status != 0 {
    return fmt.Errorf("execution engine error: %s (code: %d)", 
        C.GoString(result.error_message), result.status)
}
```

### Memory Management Issues

**Problem**: Memory leaks or corruption when using C++ execution engine.

**Solution**:
1. Implement proper memory cleanup:
```go
// Free memory allocated by C++
defer C.free(unsafe.Pointer(result.data))
defer C.free(unsafe.Pointer(result.error_message))
```

2. Use CGO memory allocation for data passed to C++:
```go
// Allocate memory that C++ can free
cData := C.CString(data)
defer C.free(unsafe.Pointer(cData))

// Pass to C++ function
C.ProcessData(cData, C.int(len(data)))
```

3. Implement resource tracking:
```go
// Track resources for cleanup
type Resource struct {
    ptr  unsafe.Pointer
    size int
}

var resources = make(map[string]Resource)
var resourcesMutex sync.Mutex

// Allocate and track resource
func allocateResource(name string, size int) unsafe.Pointer {
    ptr := C.malloc(C.size_t(size))
    
    resourcesMutex.Lock()
    resources[name] = Resource{ptr, size}
    resourcesMutex.Unlock()
    
    return ptr
}

// Free resource
func freeResource(name string) {
    resourcesMutex.Lock()
    resource, exists := resources[name]
    if exists {
        C.free(resource.ptr)
        delete(resources, name)
    }
    resourcesMutex.Unlock()
}

// Free all resources
func freeAllResources() {
    resourcesMutex.Lock()
    for name, resource := range resources {
        C.free(resource.ptr)
        delete(resources, name)
    }
    resourcesMutex.Unlock()
}
```

### Performance Issues

**Problem**: C++ execution engine performance is slower than expected.

**Solution**:
1. Minimize data serialization/deserialization:
```go
// Use shared memory instead of JSON serialization
shm, err := shm.Create("execution_engine_shm", 1024*1024)
if err != nil {
    return err
}
defer shm.Close()

// Write data to shared memory
copy(shm.Data(), []byte(data))

// Call C++ function with shared memory reference
C.ProcessDataFromSharedMemory(C.CString("execution_engine_shm"), C.int(len(data)))
```

2. Batch operations when possible:
```go
// Batch multiple orders into a single call
orders := []*Order{order1, order2, order3}
ordersJSON, err := json.Marshal(orders)
if err != nil {
    return err
}

// Process batch
result := C.ProcessOrderBatch(C.CString(string(ordersJSON)))
```

3. Profile and optimize critical paths:
```go
// Profile C++ execution time
start := time.Now()
result := C.ExecuteOrder(C.CString(orderJSON))
duration := time.Since(start)

// Log performance metrics
log.Info("Order execution completed", 
    "duration_ms", duration.Milliseconds(),
    "order_id", order.OrderID)

// Adjust thread pool size based on performance
if duration > 100*time.Millisecond {
    // Increase thread pool size
    C.SetThreadPoolSize(C.int(runtime.NumCPU()))
}
```

## Multi-Broker Integration Issues

### Broker-Specific Errors

**Problem**: Different brokers return different error formats.

**Solution**:
1. Implement broker-specific error handling:
```go
// Handle broker-specific errors
switch client.BrokerType {
case BrokerTypeXTSPro:
    return handleXTSProError(err)
case BrokerTypeXTSClient:
    return handleXTSClientError(err)
case BrokerTypeZerodha:
    return handleZerodhaError(err)
default:
    return err
}
```

2. Normalize error codes across brokers:
```go
// Map broker-specific error codes to common codes
func mapErrorCode(brokerType BrokerType, brokerErrorCode string) string {
    errorMappings := map[BrokerType]map[string]string{
        BrokerTypeXTSPro: {
            "E001": ErrCodeAuthenticationFailed,
            "E002": ErrCodeConnectionFailed,
            // ...
        },
        BrokerTypeZerodha: {
            "NE001": ErrCodeAuthenticationFailed,
            "NE002": ErrCodeConnectionFailed,
            // ...
        },
    }
    
    if mappings, exists := errorMappings[brokerType]; exists {
        if commonCode, mapped := mappings[brokerErrorCode]; mapped {
            return commonCode
        }
    }
    
    return ErrCodeUnknown
}
```

3. Create unified error reporting:
```go
// Create unified error report
func createErrorReport(err error) map[string]interface{} {
    report := map[string]interface{}{
        "timestamp": time.Now(),
        "error": err.Error(),
    }
    
    var xtsErr *errors.XTSError
    if errors.As(err, &xtsErr) {
        report["code"] = xtsErr.Code
        report["message"] = xtsErr.Message
        report["description"] = xtsErr.Description
        report["http_status"] = xtsErr.HTTPStatus
    }
    
    var execErr *orderexecution.ExecutionError
    if errors.As(err, &execErr) {
        report["type"] = execErr.Type
        report["severity"] = execErr.Severity
        report["code"] = execErr.Code
        report["order_id"] = execErr.OrderID
        report["retry_count"] = execErr.RetryCount
    }
    
    return report
}
```

### Authentication Differences

**Problem**: Different brokers have different authentication mechanisms.

**Solution**:
1. Implement broker-specific authentication:
```go
// Authenticate based on broker type
func (manager *BrokerManager) Login(brokerID string, credentials *Credentials) (*Session, error) {
    client, err := manager.GetBrokerClient(brokerID)
    if err != nil {
        return nil, err
    }
    
    switch client.BrokerType {
    case BrokerTypeXTSPro, BrokerTypeXTSClient:
        // XTS authentication (API key + secret)
        return client.Login(credentials)
    case BrokerTypeZerodha:
        // Zerodha two-step authentication
        requestToken, err := client.GenerateRequestToken(credentials)
        if err != nil {
            return nil, err
        }
        
        return client.GenerateSession(requestToken)
    default:
        return nil, fmt.Errorf("unsupported broker type: %s", client.BrokerType)
    }
}
```

2. Handle token refresh differences:
```go
// Refresh token based on broker type
func (manager *BrokerManager) RefreshToken(brokerID string, session *Session) (*Session, error) {
    client, err := manager.GetBrokerClient(brokerID)
    if err != nil {
        return nil, err
    }
    
    switch client.BrokerType {
    case BrokerTypeXTSPro:
        // XTS Pro token refresh
        return client.RefreshToken(session.Token)
    case BrokerTypeXTSClient:
        // XTS Client token refresh
        return client.RefreshToken(session.Token)
    case BrokerTypeZerodha:
        // Zerodha requires re-authentication
        credentials, err := manager.GetStoredCredentials(brokerID)
        if err != nil {
            return nil, err
        }
        
        return manager.Login(brokerID, credentials)
    default:
        return nil, fmt.Errorf("unsupported broker type: %s", client.BrokerType)
    }
}
```

3. Implement session management:
```go
// Session manager for different broker types
type SessionManager struct {
    sessions     map[string]*Session
    sessionMutex sync.RWMutex
    refreshers   map[string]*time.Timer
}

// Get session with automatic refresh
func (sm *SessionManager) GetSession(brokerID string) (*Session, error) {
    sm.sessionMutex.RLock()
    session, exists := sm.sessions[brokerID]
    sm.sessionMutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("no session found for broker: %s", brokerID)
    }
    
    // Check if session is about to expire
    if time.Until(session.ExpiryTime) < 5*time.Minute {
        // Refresh session
        newSession, err := manager.RefreshToken(brokerID, session)
        if err != nil {
            return nil, err
        }
        
        sm.sessionMutex.Lock()
        sm.sessions[brokerID] = newSession
        sm.sessionMutex.Unlock()
        
        return newSession, nil
    }
    
    return session, nil
}
```

### Order Parameter Differences

**Problem**: Different brokers require different order parameters.

**Solution**:
1. Create broker-specific order adapters:
```go
// Adapt common order to broker-specific format
func adaptOrder(order *Order, brokerType BrokerType) (interface{}, error) {
    switch brokerType {
    case BrokerTypeXTSPro:
        return &xts.Order{
            ExchangeSegment:      order.Exchange,
            ExchangeInstrumentID: order.Symbol,
            OrderSide:            mapOrderSide(order.Side, brokerType),
            OrderType:            mapOrderType(order.Type, brokerType),
            Quantity:             order.Quantity,
            Price:                order.Price,
            TriggerPrice:         order.StopPrice,
            ProductType:          mapProductType(order.ProductType, brokerType),
            TimeInForce:          mapTimeInForce(order.TimeInForce, brokerType),
        }, nil
    case BrokerTypeZerodha:
        return &zerodha.Order{
            Exchange:      mapExchange(order.Exchange, brokerType),
            TradingSymbol: order.Symbol,
            TransactionType: mapOrderSide(order.Side, brokerType),
            OrderType:     mapOrderType(order.Type, brokerType),
            Quantity:      int(order.Quantity),
            Price:         float64(order.Price),
            TriggerPrice:  float64(order.StopPrice),
            Product:       mapProductType(order.ProductType, brokerType),
            Validity:      mapTimeInForce(order.TimeInForce, brokerType),
        }, nil
    default:
        return nil, fmt.Errorf("unsupported broker type: %s", brokerType)
    }
}
```

2. Implement parameter mapping functions:
```go
// Map order side to broker-specific value
func mapOrderSide(side OrderSide, brokerType BrokerType) string {
    switch brokerType {
    case BrokerTypeXTSPro, BrokerTypeXTSClient:
        if side == OrderSideBuy {
            return "BUY"
        }
        return "SELL"
    case BrokerTypeZerodha:
        if side == OrderSideBuy {
            return "BUY"
        }
        return "SELL"
    default:
        return string(side)
    }
}

// Map order type to broker-specific value
func mapOrderType(orderType OrderType, brokerType BrokerType) string {
    switch brokerType {
    case BrokerTypeXTSPro, BrokerTypeXTSClient:
        switch orderType {
        case OrderTypeMarket:
            return "MARKET"
        case OrderTypeLimit:
            return "LIMIT"
        case OrderTypeStopLimit:
            return "STOPLIMIT"
        default:
            return "MARKET"
        }
    case BrokerTypeZerodha:
        switch orderType {
        case OrderTypeMarket:
            return "MARKET"
        case OrderTypeLimit:
            return "LIMIT"
        case OrderTypeStopLimit:
            return "SL"
        default:
            return "MARKET"
        }
    default:
        return string(orderType)
    }
}
```

3. Validate broker-specific constraints:
```go
// Validate order based on broker-specific constraints
func validateOrder(order *Order, brokerType BrokerType) error {
    switch brokerType {
    case BrokerTypeXTSPro:
        // XTS Pro specific validations
        if order.Quantity <= 0 {
            return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_QUANTITY", 
                "Quantity must be greater than zero")
        }
        if order.Type == OrderTypeLimit && order.Price <= 0 {
            return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_PRICE", 
                "Price must be greater than zero for limit orders")
        }
    case BrokerTypeZerodha:
        // Zerodha specific validations
        if order.Quantity <= 0 || int(order.Quantity)%1 != 0 {
            return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_QUANTITY", 
                "Quantity must be a positive integer")
        }
        if order.Type == OrderTypeLimit && order.Price <= 0 {
            return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_PRICE", 
                "Price must be greater than zero for limit orders")
        }
        // Zerodha requires trigger price for stop orders
        if order.Type == OrderTypeStopLimit && order.StopPrice <= 0 {
            return errors.Wrap(errors.ErrInvalidOrderParams, "INVALID_TRIGGER_PRICE", 
                "Trigger price must be greater than zero for stop orders")
        }
    }
    
    return nil
}
```

## Performance Issues

### High Latency

**Problem**: System experiences high latency during order execution.

**Solution**:
1. Implement performance monitoring:
```go
// Monitor operation latency
func monitorLatency(operation string, fn func() error) error {
    start := time.Now()
    err := fn()
    duration := time.Since(start)
    
    // Log latency
    log.Info("Operation completed", 
        "operation", operation,
        "duration_ms", duration.Milliseconds(),
        "error", err != nil)
    
    // Record metrics
    metrics.RecordLatency(operation, duration)
    
    return err
}
```

2. Optimize critical paths:
```go
// Use connection pooling
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 10 * time.Second,
}
```

3. Implement caching for frequently accessed data:
```go
// Create cache with expiration
cache := cache.New(5*time.Minute, 10*time.Minute)

// Get data with caching
func getDataWithCache(key string, fetchFn func() (interface{}, error)) (interface{}, error) {
    // Check cache first
    if data, found := cache.Get(key); found {
        return data, nil
    }
    
    // Fetch data
    data, err := fetchFn()
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    cache.Set(key, data, cache.DefaultExpiration)
    
    return data, nil
}
```

### Memory Leaks

**Problem**: System experiences memory growth over time.

**Solution**:
1. Implement resource cleanup:
```go
// Ensure resources are properly closed
func withResource(resource io.Closer, fn func() error) error {
    defer func() {
        err := resource.Close()
        if err != nil {
            log.Error("Failed to close resource", "error", err)
        }
    }()
    
    return fn()
}
```

2. Use context for cancellation:
```go
// Create context with cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Use context for operations
client.ExecuteWithContext(ctx)
```

3. Implement memory usage monitoring:
```go
// Monitor memory usage
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            log.Info("Memory stats", 
                "alloc_mb", m.Alloc/1024/1024,
                "sys_mb", m.Sys/1024/1024,
                "num_gc", m.NumGC)
            
            // Alert if memory usage is too high
            if m.Alloc > 1024*1024*1024 { // 1GB
                log.Warn("High memory usage detected")
            }
        case <-ctx.Done():
            return
        }
    }
}()
```

### Concurrency Issues

**Problem**: System experiences race conditions or deadlocks.

**Solution**:
1. Use proper synchronization:
```go
// Use mutex for shared resources
var (
    cache      = make(map[string]interface{})
    cacheMutex sync.RWMutex
)

// Thread-safe cache operations
func getFromCache(key string) (interface{}, bool) {
    cacheMutex.RLock()
    defer cacheMutex.RUnlock()
    
    value, exists := cache[key]
    return value, exists
}

func setInCache(key string, value interface{}) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    
    cache[key] = value
}
```

2. Implement deadlock detection:
```go
// Detect potential deadlocks
func withTimeout(timeout time.Duration, fn func() error) error {
    done := make(chan error, 1)
    
    go func() {
        done <- fn()
    }()
    
    select {
    case err := <-done:
        return err
    case <-time.After(timeout):
        return fmt.Errorf("operation timed out after %v, possible deadlock", timeout)
    }
}
```

3. Use channels for coordination:
```go
// Coordinate multiple goroutines
func processItems(items []Item) error {
    numWorkers := runtime.NumCPU()
    itemChan := make(chan Item, len(items))
    errChan := make(chan error, 1)
    doneChan := make(chan struct{})
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        go func() {
            for item := range itemChan {
                if err := processItem(item); err != nil {
                    select {
                    case errChan <- err:
                    default:
                    }
                }
            }
        }()
    }
    
    // Send items to workers
    go func() {
        for _, item := range items {
            itemChan <- item
        }
        close(itemChan)
        doneChan <- struct{}{}
    }()
    
    // Wait for completion or error
    select {
    case <-doneChan:
        return nil
    case err := <-errChan:
        return err
    }
}
```

## Environment Isolation Issues

### Cross-Environment Data Leakage

**Problem**: Data from one environment (e.g., live) leaks into another (e.g., simulation).

**Solution**:
1. Implement strict environment checking:
```go
// Check environment before operations
func (service *OrderService) PlaceOrder(ctx context.Context, order *Order) (*OrderResponse, error) {
    // Get environment from context
    env, err := GetEnvironmentFromContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Validate environment
    if order.Environment != env {
        return nil, errors.Wrap(errors.ErrInvalidOrderParams, "ENVIRONMENT_MISMATCH", 
            fmt.Sprintf("Order environment (%s) does not match current environment (%s)", 
                order.Environment, env))
    }
    
    // Proceed with order placement
    return service.orderProcessor.ProcessOrder(ctx, order)
}
```

2. Use environment-specific database schemas:
```go
// Get database connection for specific environment
func GetDBForEnvironment(env Environment) (*sql.DB, error) {
    switch env {
    case EnvironmentLive:
        return sql.Open("postgres", "host=db-live user=app password=secret dbname=trading_live")
    case EnvironmentSimulation:
        return sql.Open("postgres", "host=db-sim user=app password=secret dbname=trading_sim")
    default:
        return nil, fmt.Errorf("unknown environment: %s", env)
    }
}
```

3. Implement environment middleware:
```go
// Middleware to enforce environment isolation
func EnvironmentMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract environment from request (e.g., from JWT token)
        env, err := ExtractEnvironmentFromRequest(r)
        if err != nil {
            http.Error(w, "Invalid environment", http.StatusBadRequest)
            return
        }
        
        // Add environment to request context
        ctx := context.WithValue(r.Context(), EnvironmentKey, env)
        
        // Call next handler with updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Authentication Confusion

**Problem**: Users authenticated in one environment try to access another.

**Solution**:
1. Include environment in authentication tokens:
```go
// Generate JWT token with environment
func GenerateToken(userID string, env Environment) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "env": string(env),
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

2. Validate environment during authentication:
```go
// Validate token and environment
func ValidateToken(tokenString string, expectedEnv Environment) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })
    if err != nil {
        return "", err
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", errors.New("invalid token")
    }
    
    // Check environment
    tokenEnv, ok := claims["env"].(string)
    if !ok || Environment(tokenEnv) != expectedEnv {
        return "", errors.New("environment mismatch")
    }
    
    // Return user ID
    return claims["sub"].(string), nil
}
```

3. Implement environment switching:
```go
// Switch user environment
func (service *UserService) SwitchEnvironment(userID string, newEnv Environment) (string, error) {
    // Validate environment
    if newEnv != EnvironmentLive && newEnv != EnvironmentSimulation {
        return "", errors.New("invalid environment")
    }
    
    // Update user preferences
    err := service.userRepo.UpdateUserEnvironment(userID, newEnv)
    if err != nil {
        return "", err
    }
    
    // Generate new token with updated environment
    return GenerateToken(userID, newEnv)
}
```

### UI Environment Indicators

**Problem**: Users are unaware of which environment they are currently using.

**Solution**:
1. Implement clear visual indicators:
```typescript
// React component for environment indicator
const EnvironmentIndicator: React.FC<{ environment: string }> = ({ environment }) => {
    const getIndicatorStyle = () => {
        switch (environment) {
            case 'LIVE':
                return {
                    backgroundColor: '#d9534f',
                    color: 'white',
                };
            case 'SIMULATION':
                return {
                    backgroundColor: '#5bc0de',
                    color: 'white',
                };
            default:
                return {
                    backgroundColor: '#f0ad4e',
                    color: 'white',
                };
        }
    };
    
    return (
        <div className="environment-indicator" style={getIndicatorStyle()}>
            {environment}
        </div>
    );
};
```

2. Add environment to page title:
```typescript
// Update page title based on environment
useEffect(() => {
    document.title = `Trading Platform - ${environment.toUpperCase()}`;
}, [environment]);
```

3. Implement confirmation dialogs for critical actions:
```typescript
// Confirm order placement with environment check
const confirmOrderPlacement = (order) => {
    return new Promise((resolve, reject) => {
        if (order.environment === 'LIVE') {
            if (window.confirm(`You are about to place a LIVE order for ${order.quantity} ${order.symbol} at ${order.price}. Proceed?`)) {
                resolve();
            } else {
                reject(new Error('Order placement cancelled by user'));
            }
        } else {
            resolve();
        }
    });
};
```

## Debugging Techniques

### Logging

Implement comprehensive logging to diagnose integration issues:

```go
// Configure structured logging
logger := log.NewLogger(log.Config{
    Level:      log.LevelInfo,
    Format:     log.FormatJSON,
    OutputPath: "logs/integration.log",
})

// Log with context
logger.Info("Processing order",
    "order_id", order.OrderID,
    "symbol", order.Symbol,
    "quantity", order.Quantity,
    "price", order.Price,
    "environment", order.Environment,
)

// Log errors with context
if err != nil {
    logger.Error("Order processing failed",
        "order_id", order.OrderID,
        "error", err.Error(),
        "stack_trace", debug.Stack(),
    )
}
```

### Request/Response Tracing

Implement request/response tracing for API calls:

```go
// Trace HTTP requests
func TraceHTTPRequest(req *http.Request) (string, error) {
    traceID := uuid.New().String()
    
    // Add trace ID to request headers
    req.Header.Set("X-Trace-ID", traceID)
    
    // Log request
    body, err := io.ReadAll(req.Body)
    if err != nil {
        return traceID, err
    }
    
    // Restore request body
    req.Body = io.NopCloser(bytes.NewBuffer(body))
    
    log.Info("API Request",
        "trace_id", traceID,
        "method", req.Method,
        "url", req.URL.String(),
        "headers", req.Header,
        "body", string(body),
    )
    
    return traceID, nil
}

// Trace HTTP responses
func TraceHTTPResponse(traceID string, resp *http.Response) error {
    // Log response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    // Restore response body
    resp.Body = io.NopCloser(bytes.NewBuffer(body))
    
    log.Info("API Response",
        "trace_id", traceID,
        "status", resp.Status,
        "headers", resp.Header,
        "body", string(body),
    )
    
    return nil
}
```

### Correlation IDs

Implement correlation IDs to track operations across components:

```go
// Generate correlation ID
func GenerateCorrelationID() string {
    return uuid.New().String()
}

// Add correlation ID to context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
    return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// Get correlation ID from context
func GetCorrelationID(ctx context.Context) string {
    correlationID, ok := ctx.Value(CorrelationIDKey).(string)
    if !ok {
        return "unknown"
    }
    return correlationID
}

// Use correlation ID in logs
func LogWithCorrelation(ctx context.Context, msg string, fields ...interface{}) {
    correlationID := GetCorrelationID(ctx)
    
    // Add correlation ID to fields
    fields = append(fields, "correlation_id", correlationID)
    
    log.Info(msg, fields...)
}
```

## Common Error Codes

| Error Code | Description | Troubleshooting Steps |
|------------|-------------|----------------------|
| `ERR_AUTHENTICATION_FAILED` | Authentication with broker API failed | Check API keys, verify credentials, ensure API is enabled |
| `ERR_TOKEN_EXPIRED` | Authentication token has expired | Implement automatic token refresh, re-authenticate |
| `ERR_CONNECTION_FAILED` | Connection to broker API failed | Check network connectivity, verify API endpoint, check firewall settings |
| `ERR_REQUEST_TIMEOUT` | Request to broker API timed out | Increase timeout duration, check network latency, retry with backoff |
| `ERR_RATE_LIMIT_EXCEEDED` | API rate limit exceeded | Implement rate limiting, add delays between requests, batch operations |
| `ERR_INVALID_ORDER_PARAMS` | Invalid order parameters | Validate order parameters before submission, check broker-specific constraints |
| `ERR_ORDER_REJECTED` | Order rejected by broker | Check margin requirements, verify position limits, validate order parameters |
| `ERR_ORDER_NOT_FOUND` | Order not found | Verify order ID, check if order was actually submitted, implement order tracking |
| `ERR_WEBSOCKET_CONN_FAILED` | WebSocket connection failed | Check WebSocket endpoint, verify authentication, implement reconnection logic |
| `ERR_WEBSOCKET_CLOSED` | WebSocket connection closed unexpectedly | Implement automatic reconnection, add heartbeat, handle connection close gracefully |
| `ERR_SUBSCRIPTION_FAILED` | Subscription to market data failed | Verify subscription parameters, implement subscription retry logic, check permissions |
| `ERR_ENVIRONMENT_MISMATCH` | Environment mismatch | Implement strict environment checking, use environment-specific connections |
| `ERR_INSUFFICIENT_MARGIN` | Insufficient margin for order | Check account balance, verify margin requirements, implement pre-trade checks |
| `ERR_POSITION_LIMIT_EXCEEDED` | Position limit exceeded | Check position limits, implement pre-trade risk checks |
| `ERR_INTERNAL_ERROR` | Internal system error | Check logs for details, implement comprehensive error handling |

## Conclusion

This troubleshooting guide provides solutions for common integration issues in the Trading Platform. By following these guidelines, developers can identify and resolve issues more efficiently, ensuring a smooth integration experience.

For additional assistance, please contact the platform support team or refer to the API documentation for specific broker integrations.
