# Module 6: Multileg Layout Backend

## Overview
The Multileg Layout Backend module provides comprehensive functionality for creating, managing, executing, and monitoring complex multileg trading strategies. This module enables traders to define sophisticated strategies with multiple legs, execution parameters, risk management rules, and dynamic hedging capabilities.

## Architecture
The module follows a clean, layered architecture with proper separation of concerns:

1. **Models Layer**: Defines the data structures and validation logic
2. **Repository Layer**: Handles data persistence with MongoDB
3. **Service Layer**: Implements business logic and operations
4. **API Layer**: Exposes RESTful endpoints for client interaction
5. **Specialized Components**: Provides execution, monitoring, and hedging capabilities

## Data Models

### MultilegStrategy
The core model representing a complex trading strategy with multiple legs:

```go
type MultilegStrategy struct {
    ID              string           // Unique identifier
    Name            string           // Strategy name
    Description     string           // Strategy description
    UserID          string           // Owner user ID
    PortfolioID     string           // Associated portfolio
    Legs            []Leg            // Strategy legs
    ExecutionParams ExecutionParams  // Execution parameters
    RiskParams      RiskParameters   // Risk management parameters
    HedgeParams     HedgeParameters  // Hedging parameters
    Status          string           // Strategy status
    Tags            []string         // Strategy tags
    CreatedAt       time.Time        // Creation timestamp
    UpdatedAt       time.Time        // Last update timestamp
    LastExecutedAt  time.Time        // Last execution timestamp
}
```

### Leg
Represents a single leg in a multileg strategy:

```go
type Leg struct {
    ID              string        // Unique identifier
    StrategyID      string        // Parent strategy ID
    Symbol          string        // Trading symbol
    Type            LegType       // Leg type (BUY, SELL, etc.)
    Quantity        int           // Quantity to trade
    ExecutionType   ExecutionType // Execution type (MARKET, LIMIT, etc.)
    Price           float64       // Limit price (if applicable)
    StopPrice       float64       // Stop price (if applicable)
    TrailingAmount  float64       // Trailing amount (if applicable)
    TrailingPercent float64       // Trailing percentage (if applicable)
    Status          LegStatus     // Leg status
    OrderID         string        // Associated order ID
    ExecutionTime   time.Time     // Execution timestamp
    ExecutedPrice   float64       // Executed price
    Sequence        int           // Execution sequence
    DependsOn       []string      // Dependency leg IDs
    CreatedAt       time.Time     // Creation timestamp
    UpdatedAt       time.Time     // Last update timestamp
}
```

### ExecutionParams
Defines parameters for strategy execution:

```go
type ExecutionParams struct {
    Sequential      bool           // Execute legs sequentially
    SimultaneousLegs bool          // Execute legs simultaneously
    TimeWindow      int            // Execution time window in seconds
    MaxSlippage     float64        // Maximum allowed slippage
    EntryConditions []Condition    // Entry conditions
    ExitConditions  []Condition    // Exit conditions
    RangeBreakout   RangeBreakout  // Range breakout parameters
}
```

### RangeBreakout
Defines parameters for range breakout execution:

```go
type RangeBreakout struct {
    Enabled      bool    // Range breakout enabled
    UpperBound   float64 // Upper bound price
    LowerBound   float64 // Lower bound price
    Confirmation int     // Confirmation ticks
    Symbol       string  // Symbol to monitor
}
```

### HedgeParameters
Defines parameters for strategy hedging:

```go
type HedgeParameters struct {
    Type             HedgeType // Hedge type (DELTA, GAMMA, etc.)
    Instrument       string    // Hedging instrument
    Ratio            float64   // Hedge ratio
    RebalanceFreq    int       // Rebalance frequency in minutes
    DynamicThreshold float64   // Dynamic threshold for rebalancing
    Enabled          bool      // Hedging enabled
}
```

## Components

### MultilegService
The core service providing operations for multileg strategies:

- **CRUD Operations**: Create, retrieve, update, and delete strategies
- **Leg Operations**: Add, update, remove, and retrieve legs
- **Execution Operations**: Execute, pause, resume, and cancel strategies
- **Monitoring Operations**: Get status and performance metrics

### MultilegRepository
Handles data persistence for multileg strategies:

- **Strategy Operations**: Store, retrieve, update, and delete strategies
- **MongoDB Integration**: Efficient data storage and retrieval

### MultilegHandler
Exposes RESTful API endpoints for multileg strategies:

- **Strategy Endpoints**: CRUD operations for strategies
- **Leg Endpoints**: Operations for managing legs
- **Execution Endpoints**: Control strategy execution
- **Monitoring Endpoints**: Retrieve status and performance

### ExecutionEngine
Handles the execution of multileg strategies:

- **Sequential Execution**: Execute legs in sequence
- **Simultaneous Execution**: Execute legs simultaneously
- **Order Creation**: Create orders for each leg
- **Execution Tracking**: Track execution status

### RangeBreakoutMonitor
Monitors for range breakouts:

- **Price Monitoring**: Monitor price movements
- **Breakout Detection**: Detect range breakouts
- **Strategy Execution**: Trigger strategy execution on breakout

### DynamicHedgeService
Provides dynamic hedging capabilities:

- **Hedge Creation**: Create hedges for strategies
- **Hedge Adjustment**: Dynamically adjust hedges
- **Rebalancing**: Periodic hedge rebalancing

## API Endpoints

### Strategy Endpoints
- `POST /api/multileg`: Create a new multileg strategy
- `GET /api/multileg/{strategyId}`: Get a strategy by ID
- `GET /api/users/{userId}/multileg`: Get all strategies for a user
- `GET /api/portfolios/{portfolioId}/multileg`: Get all strategies for a portfolio
- `PUT /api/multileg/{strategyId}`: Update a strategy
- `DELETE /api/multileg/{strategyId}`: Delete a strategy

### Leg Endpoints
- `POST /api/multileg/{strategyId}/legs`: Add a leg to a strategy
- `GET /api/multileg/{strategyId}/legs`: Get all legs for a strategy
- `PUT /api/multileg/{strategyId}/legs/{legId}`: Update a leg
- `DELETE /api/multileg/{strategyId}/legs/{legId}`: Remove a leg

### Execution Endpoints
- `POST /api/multileg/{strategyId}/execute`: Execute a strategy
- `POST /api/multileg/{strategyId}/pause`: Pause a strategy
- `POST /api/multileg/{strategyId}/resume`: Resume a strategy
- `POST /api/multileg/{strategyId}/cancel`: Cancel a strategy

### Monitoring Endpoints
- `GET /api/multileg/{strategyId}/status`: Get strategy status
- `GET /api/multileg/{strategyId}/performance`: Get strategy performance

## Implementation Details

### Validation
The module implements comprehensive validation for all models:

- **Strategy Validation**: Name, user ID, portfolio ID, legs
- **Leg Validation**: Symbol, type, quantity, execution type
- **Execution Parameter Validation**: Slippage, time window
- **Risk Parameter Validation**: Max loss, max daily loss
- **Hedge Parameter Validation**: Type, ratio, rebalance frequency

### Error Handling
The module provides robust error handling:

- **Input Validation Errors**: Invalid input parameters
- **Not Found Errors**: Strategy or leg not found
- **Execution Errors**: Errors during strategy execution
- **Monitoring Errors**: Errors during status or performance retrieval

### Security
The module implements security measures:

- **User Authorization**: Strategies are associated with users
- **Portfolio Authorization**: Strategies are associated with portfolios
- **Input Validation**: Prevent injection attacks
- **Error Handling**: Prevent information leakage

## Usage Examples

### Creating a Multileg Strategy
```json
POST /api/multileg
{
  "name": "Bull Call Spread",
  "description": "A bullish options strategy",
  "userId": "user123",
  "portfolioId": "portfolio123",
  "legs": [
    {
      "symbol": "AAPL220121C00150000",
      "type": "BUY_TO_OPEN",
      "quantity": 10,
      "executionType": "LIMIT",
      "price": 5.0,
      "sequence": 1
    },
    {
      "symbol": "AAPL220121C00160000",
      "type": "SELL_TO_OPEN",
      "quantity": 10,
      "executionType": "LIMIT",
      "price": 2.0,
      "sequence": 2
    }
  ],
  "executionParams": {
    "sequential": true,
    "simultaneousLegs": false,
    "maxSlippage": 0.1
  },
  "riskParams": {
    "maxLoss": 3000,
    "maxDailyLoss": 5000
  },
  "hedgeParams": {
    "type": "DELTA",
    "instrument": "AAPL",
    "ratio": 0.5,
    "rebalanceFreq": 60,
    "enabled": true
  }
}
```

### Executing a Strategy
```
POST /api/multileg/strategy123/execute
```

### Getting Strategy Performance
```
GET /api/multileg/strategy123/performance
```

## Testing
The module includes comprehensive unit tests:

- **Service Tests**: Test business logic with mock repositories
- **Repository Tests**: Test data persistence with mock MongoDB
- **Handler Tests**: Test API endpoints with mock services
- **Execution Tests**: Test strategy execution
- **Monitoring Tests**: Test status and performance retrieval

## Future Enhancements
Potential future enhancements for the module:

- **Advanced Execution Algorithms**: TWAP, VWAP, etc.
- **Machine Learning Integration**: Predictive analytics for strategy optimization
- **Real-time Monitoring**: WebSocket integration for real-time updates
- **Strategy Templates**: Pre-defined strategy templates
- **Strategy Backtesting**: Historical performance analysis
- **Strategy Optimization**: Parameter optimization

## Dependencies
The module depends on:

- **Order Module**: For order creation and management
- **Position Module**: For position tracking
- **Portfolio Module**: For portfolio association
- **User Module**: For user association
- **MongoDB**: For data persistence
