# Module 3: Positions Management Backend Documentation

## Overview
This document provides detailed information about the Positions Management Backend implemented in the trading platform. The implementation follows a clean, layered architecture with handlers, services, and repositories.

## Architecture

### Layered Design
The Positions Management Backend follows a clean, layered architecture:

1. **Handlers Layer**: Responsible for HTTP request/response handling, parameter parsing, and input validation
2. **Services Layer**: Contains business logic, validation, and orchestration of operations
3. **Repositories Layer**: Handles data persistence and retrieval from the database

This separation of concerns ensures maintainability, testability, and scalability of the codebase.

## Position Model

The Position model represents a trading position with the following key features:
- Support for different position statuses (Open, Closed, Partial)
- Support for different position directions (Long, Short)
- Support for different instrument types (Option, Future, Stock)
- Comprehensive validation logic for all fields
- Calculation methods for P&L, Greeks, and days till expiry

### Key Fields
- `ID`: Unique identifier for the position
- `UserID`: ID of the user who owns the position
- `OrderID`: ID of the order that created the position
- `Symbol`: Trading symbol
- `Exchange`: Exchange where the position is held
- `Direction`: Long or Short
- `EntryPrice`: Price at which the position was entered
- `ExitPrice`: Price at which the position was exited (if applicable)
- `Quantity`: Total quantity of the position
- `ExitQuantity`: Quantity that has been exited (if applicable)
- `Status`: Open, Closed, or Partial
- `ProductType`: MIS, NRML, CNC
- `InstrumentType`: Option, Future, Stock
- `OptionType`: Call or Put (for options)
- `StrikePrice`: Strike price (for options)
- `Expiry`: Expiry date (for options and futures)
- `UnrealizedPnL`: Unrealized profit and loss
- `RealizedPnL`: Realized profit and loss
- `Greeks`: Option Greeks (Delta, Gamma, Theta, Vega)
- `PortfolioID`: ID of the portfolio (if applicable)
- `StrategyID`: ID of the strategy (if applicable)
- `LegID`: ID of the leg (if applicable)
- `Tags`: Custom tags for the position
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

### Greeks Structure
- `Delta`: Rate of change of option price with respect to underlying price
- `Gamma`: Rate of change of delta with respect to underlying price
- `Theta`: Rate of change of option price with respect to time
- `Vega`: Rate of change of option price with respect to volatility

### AggregatedPosition Structure
- `Key`: Aggregation key (e.g., symbol, instrument type)
- `GroupBy`: Field used for grouping
- `TotalQuantity`: Sum of quantities across positions
- `NetQuantity`: Net quantity (long minus short)
- `TotalValue`: Total value of positions
- `NetValue`: Net value (long minus short)
- `PnL`: Total profit and loss
- `Greeks`: Aggregated Greeks
- `PositionCount`: Number of positions in the group

## API Endpoints

### Position Management Endpoints

#### Create Position from Order
- **Endpoint**: `POST /api/positions/create-from-order`
- **Description**: Creates a new position from an executed order
- **Request Body**: Order object with execution details
- **Response**: Created position with assigned ID and status
- **Status Codes**:
  - 201: Position created successfully
  - 400: Invalid request payload or validation error
  - 500: Internal server error

#### Get Position by ID
- **Endpoint**: `GET /api/positions/{id}`
- **Description**: Retrieves a specific position by its ID
- **URL Parameters**: `id` - The unique identifier of the position
- **Response**: Position object with complete details
- **Status Codes**:
  - 200: Position retrieved successfully
  - 404: Position not found
  - 500: Internal server error

#### Get Positions (with filtering)
- **Endpoint**: `GET /api/positions`
- **Description**: Retrieves a list of positions with optional filtering and pagination
- **Query Parameters**:
  - `userId`: Filter by user ID
  - `symbol`: Filter by symbol
  - `status`: Filter by position status
  - `direction`: Filter by position direction
  - `productType`: Filter by product type
  - `instrumentType`: Filter by instrument type
  - `portfolioId`: Filter by portfolio ID
  - `strategyId`: Filter by strategy ID
  - `orderId`: Filter by order ID
  - `fromDate`: Filter by creation date (start)
  - `toDate`: Filter by creation date (end)
  - `page`: Page number for pagination (default: 1)
  - `limit`: Number of items per page (default: 50)
- **Response**: List of positions with pagination metadata
- **Status Codes**:
  - 200: Positions retrieved successfully
  - 500: Internal server error

#### Update Position
- **Endpoint**: `PUT /api/positions/{id}`
- **Description**: Updates an existing position
- **URL Parameters**: `id` - The unique identifier of the position
- **Request Body**: Updated position object
- **Response**: Updated position with complete details
- **Status Codes**:
  - 200: Position updated successfully
  - 400: Invalid request payload or validation error
  - 404: Position not found
  - 500: Internal server error

#### Close Position
- **Endpoint**: `POST /api/positions/{id}/close`
- **Description**: Closes an existing position
- **URL Parameters**: `id` - The unique identifier of the position
- **Request Body**: Close parameters (exit price and quantity)
- **Response**: Closed position with updated details
- **Status Codes**:
  - 200: Position closed successfully
  - 400: Invalid request payload or validation error
  - 404: Position not found
  - 500: Internal server error

#### Calculate P&L
- **Endpoint**: `GET /api/positions/{id}/pnl`
- **Description**: Calculates the P&L for a specific position
- **URL Parameters**: `id` - The unique identifier of the position
- **Response**: P&L calculation result
- **Status Codes**:
  - 200: P&L calculated successfully
  - 404: Position not found
  - 500: Internal server error

#### Calculate Greeks
- **Endpoint**: `GET /api/positions/{id}/greeks`
- **Description**: Calculates the Greeks for a specific position
- **URL Parameters**: `id` - The unique identifier of the position
- **Response**: Greeks calculation result
- **Status Codes**:
  - 200: Greeks calculated successfully
  - 404: Position not found
  - 500: Internal server error

### User-Specific Position Endpoints

#### Get Positions by User
- **Endpoint**: `GET /api/users/{userId}/positions`
- **Description**: Retrieves all positions for a specific user
- **URL Parameters**: `userId` - The unique identifier of the user
- **Query Parameters**: Same filtering and pagination options as the main Get Positions endpoint
- **Response**: List of positions with pagination metadata
- **Status Codes**:
  - 200: Positions retrieved successfully
  - 500: Internal server error

#### Calculate Exposure
- **Endpoint**: `GET /api/users/{userId}/exposure`
- **Description**: Calculates the total exposure for a user's positions
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: Exposure calculation result
- **Status Codes**:
  - 200: Exposure calculated successfully
  - 500: Internal server error

### Strategy-Specific Position Endpoints

#### Get Positions by Strategy
- **Endpoint**: `GET /api/strategies/{strategyId}/positions`
- **Description**: Retrieves all positions for a specific strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Query Parameters**: Same filtering and pagination options as the main Get Positions endpoint
- **Response**: List of positions with pagination metadata
- **Status Codes**:
  - 200: Positions retrieved successfully
  - 500: Internal server error

### Portfolio-Specific Position Endpoints

#### Get Positions by Portfolio
- **Endpoint**: `GET /api/portfolios/{portfolioId}/positions`
- **Description**: Retrieves all positions for a specific portfolio
- **URL Parameters**: `portfolioId` - The unique identifier of the portfolio
- **Query Parameters**: Same filtering and pagination options as the main Get Positions endpoint
- **Response**: List of positions with pagination metadata
- **Status Codes**:
  - 200: Positions retrieved successfully
  - 500: Internal server error

### Position Aggregation Endpoint

#### Aggregate Positions
- **Endpoint**: `GET /api/positions/aggregate`
- **Description**: Aggregates positions by the specified grouping
- **Query Parameters**:
  - `userId`: Filter by user ID
  - `groupBy`: Field to group by (symbol, instrumentType, productType, strategy, portfolio)
- **Response**: List of aggregated positions
- **Status Codes**:
  - 200: Positions aggregated successfully
  - 400: Invalid groupBy parameter
  - 500: Internal server error

## Implementation Details

### Position Service
The position service is responsible for:
- Creating positions from executed orders
- Retrieving positions with filtering and pagination
- Updating positions
- Closing positions (fully or partially)
- Calculating P&L for positions
- Calculating Greeks for positions
- Calculating exposure for positions
- Aggregating positions by different criteria

Key features:
- Comprehensive validation of all operations
- Automatic calculation of P&L and Greeks
- Support for partial position closure
- Position aggregation by different criteria

### Position Repository
The position repository handles data persistence:
- MongoDB integration
- CRUD operations
- Query building for filtering
- Pagination implementation

Key features:
- MongoDB integration with proper error handling
- Filter construction based on query parameters
- Pagination with offset and limit
- Sorting options

### Position Handler
The position handler is responsible for:
- Parsing HTTP requests
- Validating input data
- Calling the appropriate service methods
- Formatting and returning HTTP responses

Key features:
- Comprehensive error handling
- Proper HTTP status codes
- JSON response formatting
- Pagination support

## Calculation Methods

### P&L Calculation
The P&L calculation takes into account:
- Position direction (long or short)
- Entry price
- Current market price (or exit price for closed positions)
- Position quantity
- Exit quantity (for partially closed positions)

For long positions:
- Unrealized P&L = (Current Price - Entry Price) * Remaining Quantity
- Realized P&L = (Exit Price - Entry Price) * Exit Quantity

For short positions:
- Unrealized P&L = (Entry Price - Current Price) * Remaining Quantity
- Realized P&L = (Entry Price - Exit Price) * Exit Quantity

### Greeks Calculation
The Greeks calculation is based on option pricing models and takes into account:
- Option type (call or put)
- Position direction (long or short)
- Strike price
- Underlying price
- Time to expiry
- Volatility
- Interest rate

The implementation provides:
- Delta: Sensitivity to underlying price changes
- Gamma: Sensitivity of delta to underlying price changes
- Theta: Sensitivity to time decay
- Vega: Sensitivity to volatility changes

### Exposure Calculation
The exposure calculation takes into account:
- Position value (entry price * quantity)
- Position direction (long or short)
- Position status (open, closed, partial)

The implementation provides:
- Total exposure: Sum of position values
- Net exposure: Sum of long position values minus sum of short position values

### Position Aggregation
The position aggregation allows grouping positions by:
- Symbol
- Instrument type
- Product type
- Strategy
- Portfolio

The aggregation provides:
- Total quantity: Sum of quantities across positions
- Net quantity: Sum of long quantities minus sum of short quantities
- Total value: Sum of position values
- Net value: Sum of long position values minus sum of short position values
- P&L: Sum of P&L across positions
- Greeks: Sum of Greeks across positions
- Position count: Number of positions in the group

## Testing
Comprehensive unit tests have been implemented for all layers:
- Handler tests with mock services
- Service tests with mock repositories
- Repository tests with mock MongoDB interfaces

The tests cover:
- Happy path scenarios
- Error handling
- Edge cases
- Validation logic

## Error Handling
The API implements consistent error handling:
- Validation errors return 400 Bad Request with descriptive messages
- Not found errors return 404 Not Found
- Server errors return 500 Internal Server Error
- Custom error responses include an "error" field with a descriptive message

## Pagination
All list endpoints support pagination:
- Default page size is 50 items
- Maximum page size is 100 items
- Response includes metadata:
  - total: Total number of items matching the filter
  - page: Current page number
  - limit: Number of items per page
  - totalPages: Total number of pages
  - hasNextPage: Boolean indicating if there are more pages

## Usage Examples

### Creating a Position from an Order
```
POST /api/positions/create-from-order
Content-Type: application/json

{
  "id": "order123",
  "userId": "user123",
  "symbol": "NIFTY",
  "exchange": "NSE",
  "orderType": "LIMIT",
  "direction": "BUY",
  "quantity": 10,
  "price": 500.50,
  "executionPrice": 500.75,
  "filledQuantity": 10,
  "status": "EXECUTED",
  "productType": "MIS",
  "instrumentType": "OPTION",
  "optionType": "CE",
  "strikePrice": 18000,
  "expiry": "2025-05-03T00:00:00Z"
}
```

### Retrieving Positions with Filtering
```
GET /api/positions?userId=user123&status=OPEN&page=1&limit=20
```

### Closing a Position
```
POST /api/positions/position123/close
Content-Type: application/json

{
  "exitPrice": 550.0,
  "exitQuantity": 10
}
```

### Calculating P&L
```
GET /api/positions/position123/pnl
```

### Calculating Greeks
```
GET /api/positions/position123/greeks
```

### Calculating Exposure
```
GET /api/users/user123/exposure
```

### Aggregating Positions
```
GET /api/positions/aggregate?userId=user123&groupBy=symbol
```

## Future Enhancements
Potential future enhancements for the Positions Management Backend include:
- Real-time market data integration for accurate P&L and Greeks calculation
- Advanced risk management features
- Position hedging automation
- Performance optimization for large position portfolios
- Integration with external analytics services
- Historical position analysis
- Position performance benchmarking

## Implementation Notes
- All components follow a clean, layered architecture
- Proper separation of concerns between layers
- Comprehensive validation at all levels
- Thorough error handling
- MongoDB integration for data persistence
- RESTful API design principles
