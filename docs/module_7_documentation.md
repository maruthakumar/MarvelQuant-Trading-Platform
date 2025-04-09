# Module 7: WebSocket Integration

## Overview
The WebSocket Integration module provides real-time communication capabilities for the trading platform, enabling instant updates for orders, positions, and strategies. This module establishes a bidirectional communication channel between the server and clients, allowing for efficient and timely delivery of trading information.

## Architecture
The module follows a clean, layered architecture with proper separation of concerns:

1. **Core WebSocket Framework**: Manages connections, message handling, and client tracking
2. **Specialized Services**: Handles specific types of real-time updates
3. **Integration Layer**: Connects WebSocket functionality with other platform components
4. **API Layer**: Exposes WebSocket endpoints and status information

## Components

### Hub
The central component that manages WebSocket connections and message distribution:

```go
type Hub struct {
    // Registered clients
    clients map[*Client]bool

    // Register requests from clients
    register chan *Client

    // Unregister requests from clients
    unregister chan *Client

    // Inbound messages from clients
    broadcast chan []byte
    
    // Topic subscriptions
    topics map[string]map[*Client]bool
}
```

Key features:
- Client registration and management
- Topic-based subscription system
- Efficient message broadcasting
- Thread-safe operations

### Client
Represents an individual WebSocket connection:

```go
type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    userID   string
    topics   map[string]bool
}
```

Key features:
- Message reading and writing
- Connection lifecycle management
- Topic subscription tracking
- User association

### WebSocket Message
Standardized message format for WebSocket communication:

```go
type WebSocketMessage struct {
    Type      MessageType     // Message type (ORDER_UPDATE, POSITION_UPDATE, etc.)
    Timestamp time.Time       // Message timestamp
    Payload   json.RawMessage // Message payload
}
```

### Specialized Services

#### OrderUpdateService
Handles real-time order updates:
- Broadcasts order status changes
- Formats order information for WebSocket transmission
- Targets updates to relevant topics and users

#### PositionUpdateService
Handles real-time position updates:
- Broadcasts position changes
- Formats position information for WebSocket transmission
- Includes P&L and risk metrics

#### StrategyMonitorService
Handles real-time strategy monitoring:
- Broadcasts strategy status changes
- Provides performance metrics
- Notifies of strategy execution events

#### ConnectionManager
Manages WebSocket connections:
- Tracks active connections
- Monitors topic subscriptions
- Provides connection statistics
- Handles user disconnection

### WebSocket Handler
Exposes WebSocket endpoints:
- `/ws`: Main WebSocket connection endpoint
- `/ws/status`: Connection status information endpoint

### Authentication Middleware
Secures WebSocket connections:
- Token-based authentication
- User identification
- Connection authorization

## Message Types
The module supports various message types:

- `ORDER_UPDATE`: Order status and details updates
- `POSITION_UPDATE`: Position changes and P&L updates
- `STRATEGY_UPDATE`: Strategy status and performance updates
- `MARKET_DATA`: Real-time market data
- `AUTHENTICATION`: Authentication messages
- `SUBSCRIPTION`: Topic subscription management
- `ERROR`: Error notifications
- `HEARTBEAT`: Connection health checks

## Topic-Based Subscription
The module implements a topic-based subscription system:

- Global topics: `orders`, `positions`, `strategies`
- User-specific topics: `user:{userId}:orders`, `user:{userId}:positions`
- Strategy-specific topics: `strategy:{strategyId}:orders`, `strategy:{strategyId}`

Clients can subscribe to multiple topics to receive only relevant updates.

## Implementation Details

### Connection Lifecycle
1. **Connection Establishment**: Client connects to WebSocket endpoint
2. **Authentication**: Client is authenticated via token
3. **Registration**: Client is registered with the hub
4. **Subscription**: Client subscribes to relevant topics
5. **Message Exchange**: Bidirectional communication
6. **Disconnection**: Client disconnects or connection times out

### Message Handling
1. **Message Reception**: Messages received from clients
2. **Message Parsing**: JSON messages parsed into structured format
3. **Message Routing**: Messages routed based on type
4. **Message Processing**: Business logic applied to messages
5. **Message Broadcasting**: Updates broadcast to relevant clients

### Error Handling
- Connection errors are logged and handled gracefully
- Invalid messages trigger error responses
- Authentication failures prevent connection establishment
- Unexpected disconnections are detected and managed

### Performance Considerations
- Efficient message broadcasting to minimize latency
- Connection pooling to manage resources
- Heartbeat mechanism to detect stale connections
- Buffer management to prevent memory issues

## Integration with Other Modules
The WebSocket module integrates with:

- **Order Module**: Receives order updates for broadcasting
- **Position Module**: Receives position updates for broadcasting
- **Strategy Module**: Receives strategy updates for broadcasting
- **Authentication Module**: Verifies user identity

## API Endpoints

### WebSocket Connection
```
GET /ws
```
Query parameters:
- `token`: Authentication token

Headers:
- `Authorization`: Bearer token (alternative to query parameter)

### WebSocket Status
```
GET /ws/status
```
Response:
```json
{
  "activeConnections": 42,
  "topicSubscriptions": {
    "orders": 35,
    "positions": 28,
    "strategies": 15
  },
  "timestamp": "2025-04-03T10:15:30Z"
}
```

## Client Usage Examples

### Connecting to WebSocket
```javascript
const socket = new WebSocket('wss://api.trading-platform.com/ws?token=YOUR_TOKEN');

socket.onopen = () => {
  console.log('Connected to WebSocket');
};

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
};

socket.onclose = () => {
  console.log('Disconnected from WebSocket');
};
```

### Subscribing to Topics
```javascript
// Subscribe to topics
const subscriptionMessage = {
  type: 'SUBSCRIPTION',
  timestamp: new Date().toISOString(),
  payload: {
    action: 'subscribe',
    topics: ['orders', 'positions', 'strategies']
  }
};
socket.send(JSON.stringify(subscriptionMessage));
```

### Handling Order Updates
```javascript
socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'ORDER_UPDATE') {
    const order = JSON.parse(message.payload);
    console.log('Order update:', order);
    // Update UI with order information
  }
};
```

## Testing
The module includes comprehensive unit tests:

- **Hub Tests**: Test connection management and message broadcasting
- **Service Tests**: Test specialized services for different update types
- **Handler Tests**: Test WebSocket endpoints and authentication
- **Integration Tests**: Test interaction with other modules

## Security Considerations
- All WebSocket connections require authentication
- User-specific topics prevent unauthorized access to data
- Token validation on connection establishment
- Connection monitoring for suspicious activity
- Rate limiting to prevent abuse

## Future Enhancements
- **Compression**: Message compression for bandwidth optimization
- **Reconnection Logic**: Automatic reconnection with exponential backoff
- **Message Queuing**: Persistent message queuing for offline clients
- **Presence Awareness**: User online/offline status tracking
- **Custom Subscriptions**: More granular subscription options
- **Binary Protocol**: Binary message format for improved performance

## Dependencies
- **Gorilla WebSocket**: WebSocket implementation for Go
- **Standard Library**: JSON encoding/decoding, HTTP handling
- **Context Package**: Request context management
- **Sync Package**: Thread-safe operations
