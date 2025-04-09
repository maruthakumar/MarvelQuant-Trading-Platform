# Module 5: Strategies Management Backend Documentation

## Overview
This document provides detailed information about the Strategies Management Backend implemented in the trading platform. The implementation follows a clean, layered architecture with handlers, services, and repositories.

## Architecture

### Layered Design
The Strategies Management Backend follows a clean, layered architecture:

1. **Handlers Layer**: Responsible for HTTP request/response handling, parameter parsing, and input validation
2. **Services Layer**: Contains business logic, validation, and orchestration of operations
3. **Repositories Layer**: Handles data persistence and retrieval from the database

This separation of concerns ensures maintainability, testability, and scalability of the codebase.

## Strategy Models

The Strategies Management Backend includes several models to represent different aspects of trading strategies:

### Strategy
Represents a trading strategy with entry and exit conditions:
- Basic information (name, description)
- User association
- Strategy type (manual, automated, algo)
- Status (draft, active, paused, stopped)
- Entry and exit conditions
- Risk parameters
- Instruments
- Tags for categorization

### Condition
Represents a trading condition used in strategies:
- Condition type
- Parameter
- Operator
- Value

### RiskParameters
Represents risk management parameters for a strategy:
- Maximum position size
- Maximum loss
- Maximum daily loss
- Trailing stop percentage
- Take profit percentage

### StrategySchedule
Represents a schedule for strategy execution:
- Frequency (once, daily, weekly, monthly)
- Start time
- End time
- Days of week
- Enabled status

### StrategyPerformance
Represents performance metrics for a strategy:
- Total P&L
- Win/loss counts
- Win rate
- Maximum drawdown
- Order and position counts
- Time period

## Core Components

### Strategy Service
The strategy service is responsible for:
- Strategy CRUD operations
- Strategy execution operations
- Strategy monitoring operations
- Strategy scheduling operations
- Strategy tagging operations

Key features:
- Comprehensive validation of all operations
- Status management for strategies
- Performance calculation
- Tag management

### Strategy Execution Engine
The strategy execution engine is responsible for:
- Executing strategies based on their conditions
- Processing entry and exit conditions
- Monitoring active strategies
- Checking risk parameters

Key features:
- Concurrent execution of multiple strategies
- Safe management of active strategies
- Integration with order service for trade execution

### Strategy Scheduler
The strategy scheduler is responsible for:
- Scheduling strategies for execution
- Managing execution schedules
- Converting schedule specifications to cron expressions

Key features:
- Support for various schedule frequencies
- Cron-based scheduling
- Schedule management

### Strategy Monitoring Service
The monitoring service is responsible for:
- Monitoring stop losses
- Monitoring take profits
- Monitoring risk parameters

Key features:
- Concurrent monitoring of multiple strategies
- Real-time monitoring of positions
- Automatic strategy stopping when risk parameters are exceeded

### Strategy Repository
The strategy repository handles data persistence:
- MongoDB integration
- CRUD operations for strategies
- Schedule management
- Query building for filtering

Key features:
- MongoDB integration with proper error handling
- Efficient querying for strategies by various criteria
- Proper error handling for not found cases

## API Endpoints

### Strategy Management Endpoints

#### Create Strategy
- **Endpoint**: `POST /api/strategies`
- **Description**: Creates a new strategy
- **Request Body**: Strategy object
- **Response**: Created strategy object
- **Status Codes**:
  - 201: Strategy created successfully
  - 400: Invalid request payload or validation error
  - 500: Internal server error

#### Get Strategy
- **Endpoint**: `GET /api/strategies/{strategyId}`
- **Description**: Retrieves a strategy by ID
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Strategy object
- **Status Codes**:
  - 200: Strategy retrieved successfully
  - 404: Strategy not found
  - 500: Internal server error

#### Get User Strategies
- **Endpoint**: `GET /api/users/{userId}/strategies`
- **Description**: Retrieves all strategies for a user
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: Array of strategy objects
- **Status Codes**:
  - 200: Strategies retrieved successfully
  - 500: Internal server error

#### Update Strategy
- **Endpoint**: `PUT /api/strategies/{strategyId}`
- **Description**: Updates a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Request Body**: Updated strategy object
- **Response**: Updated strategy object
- **Status Codes**:
  - 200: Strategy updated successfully
  - 400: Invalid request payload or validation error
  - 404: Strategy not found
  - 500: Internal server error

#### Delete Strategy
- **Endpoint**: `DELETE /api/strategies/{strategyId}`
- **Description**: Deletes a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy deleted successfully
  - 400: Invalid request or strategy is active
  - 404: Strategy not found
  - 500: Internal server error

### Strategy Execution Endpoints

#### Execute Strategy
- **Endpoint**: `POST /api/strategies/{strategyId}/execute`
- **Description**: Executes a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy execution started
  - 400: Invalid request or strategy is already active
  - 404: Strategy not found
  - 500: Internal server error

#### Pause Strategy
- **Endpoint**: `POST /api/strategies/{strategyId}/pause`
- **Description**: Pauses a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy paused successfully
  - 400: Invalid request or strategy is not active
  - 404: Strategy not found
  - 500: Internal server error

#### Resume Strategy
- **Endpoint**: `POST /api/strategies/{strategyId}/resume`
- **Description**: Resumes a paused strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy resumed successfully
  - 400: Invalid request or strategy is not paused
  - 404: Strategy not found
  - 500: Internal server error

#### Stop Strategy
- **Endpoint**: `POST /api/strategies/{strategyId}/stop`
- **Description**: Stops a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy stopped successfully
  - 400: Invalid request or strategy is not active/paused
  - 404: Strategy not found
  - 500: Internal server error

### Strategy Monitoring Endpoints

#### Get Strategy Status
- **Endpoint**: `GET /api/strategies/{strategyId}/status`
- **Description**: Retrieves the status of a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Status object
- **Status Codes**:
  - 200: Status retrieved successfully
  - 404: Strategy not found
  - 500: Internal server error

#### Get Strategy Performance
- **Endpoint**: `GET /api/strategies/{strategyId}/performance`
- **Description**: Retrieves the performance of a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Performance object
- **Status Codes**:
  - 200: Performance retrieved successfully
  - 404: Strategy not found
  - 500: Internal server error

### Strategy Scheduling Endpoints

#### Schedule Strategy
- **Endpoint**: `POST /api/strategies/{strategyId}/schedule`
- **Description**: Schedules a strategy for execution
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Request Body**: Schedule object
- **Response**: Success message
- **Status Codes**:
  - 200: Strategy scheduled successfully
  - 400: Invalid request payload or validation error
  - 404: Strategy not found
  - 500: Internal server error

#### Get Strategy Schedule
- **Endpoint**: `GET /api/strategies/{strategyId}/schedule`
- **Description**: Retrieves the schedule for a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Schedule object
- **Status Codes**:
  - 200: Schedule retrieved successfully
  - 404: Strategy or schedule not found
  - 500: Internal server error

#### Update Strategy Schedule
- **Endpoint**: `PUT /api/strategies/{strategyId}/schedule`
- **Description**: Updates the schedule for a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Request Body**: Updated schedule object
- **Response**: Success message
- **Status Codes**:
  - 200: Schedule updated successfully
  - 400: Invalid request payload or validation error
  - 404: Strategy or schedule not found
  - 500: Internal server error

#### Delete Strategy Schedule
- **Endpoint**: `DELETE /api/strategies/{strategyId}/schedule`
- **Description**: Deletes the schedule for a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Response**: Success message
- **Status Codes**:
  - 200: Schedule deleted successfully
  - 404: Strategy or schedule not found
  - 500: Internal server error

### Strategy Tagging Endpoints

#### Add Strategy Tag
- **Endpoint**: `POST /api/strategies/{strategyId}/tags`
- **Description**: Adds a tag to a strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Request Body**: Tag object
- **Response**: Success message
- **Status Codes**:
  - 200: Tag added successfully
  - 400: Invalid request payload or tag already exists
  - 404: Strategy not found
  - 500: Internal server error

#### Remove Strategy Tag
- **Endpoint**: `DELETE /api/strategies/{strategyId}/tags/{tag}`
- **Description**: Removes a tag from a strategy
- **URL Parameters**: 
  - `strategyId` - The unique identifier of the strategy
  - `tag` - The tag to remove
- **Response**: Success message
- **Status Codes**:
  - 200: Tag removed successfully
  - 404: Strategy or tag not found
  - 500: Internal server error

#### Get Strategies By Tag
- **Endpoint**: `GET /api/strategies/tags/{tag}`
- **Description**: Retrieves all strategies with a specific tag
- **URL Parameters**: `tag` - The tag to search for
- **Response**: Array of strategy objects
- **Status Codes**:
  - 200: Strategies retrieved successfully
  - 500: Internal server error

## Implementation Details

### Strategy Creation and Validation
When creating a strategy, the system:
1. Validates required fields (name, user ID, type, instruments)
2. Validates risk parameters
3. Sets initial status to DRAFT
4. Sets creation and update timestamps
5. Persists the strategy in the database

### Strategy Execution
When executing a strategy, the system:
1. Checks if the strategy exists and is not already active
2. Updates the strategy status to ACTIVE
3. Records the execution time
4. Processes entry conditions to create new orders
5. Processes exit conditions for existing positions
6. Monitors the strategy for risk parameter violations

### Strategy Scheduling
When scheduling a strategy, the system:
1. Validates the schedule parameters
2. Creates a cron expression based on the schedule frequency
3. Adds the schedule to the cron scheduler
4. Persists the schedule in the database

### Strategy Monitoring
The monitoring service:
1. Continuously monitors active strategies
2. Checks stop loss conditions for positions
3. Checks take profit conditions for positions
4. Checks risk parameters for the overall strategy
5. Automatically stops strategies that exceed risk parameters

### Strategy Performance Calculation
When calculating strategy performance, the system:
1. Retrieves all orders and positions for the strategy
2. Calculates total P&L
3. Counts winning and losing trades
4. Calculates win rate
5. Determines maximum drawdown
6. Compiles all metrics into a performance object

## Data Models

### Strategy
```json
{
  "id": "strategy123",
  "name": "Moving Average Crossover",
  "description": "A strategy that trades based on moving average crossovers",
  "userId": "user123",
  "type": "AUTOMATED",
  "status": "ACTIVE",
  "entryConditions": [
    {
      "type": "INDICATOR",
      "parameter": "MA_CROSSOVER",
      "operator": "CROSSES_ABOVE",
      "value": {
        "fastPeriod": 9,
        "slowPeriod": 21
      }
    }
  ],
  "exitConditions": [
    {
      "type": "INDICATOR",
      "parameter": "MA_CROSSOVER",
      "operator": "CROSSES_BELOW",
      "value": {
        "fastPeriod": 9,
        "slowPeriod": 21
      }
    }
  ],
  "riskParameters": {
    "maxPositionSize": 100,
    "maxLoss": 1000,
    "maxDailyLoss": 2000,
    "trailingStopPercent": 2.0,
    "takeProfitPercent": 5.0
  },
  "instruments": ["AAPL", "MSFT", "GOOGL"],
  "tags": ["momentum", "technical", "stocks"],
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T10:15:00Z",
  "lastExecutedAt": "2025-04-03T10:15:00Z"
}
```

### StrategySchedule
```json
{
  "strategyId": "strategy123",
  "frequency": "DAILY",
  "startTime": "2025-04-03T09:30:00Z",
  "endTime": "2025-04-03T16:00:00Z",
  "daysOfWeek": [1, 2, 3, 4, 5],
  "enabled": true,
  "createdAt": "2025-04-03T09:00:00Z",
  "updatedAt": "2025-04-03T09:00:00Z"
}
```

### StrategyPerformance
```json
{
  "strategyId": "strategy123",
  "totalPnL": 1250.75,
  "winCount": 15,
  "lossCount": 5,
  "totalTrades": 20,
  "winRate": 75.0,
  "maxDrawdown": -350.25,
  "orderCount": 42,
  "positionCount": 20,
  "startDate": "2025-03-01T00:00:00Z",
  "endDate": "2025-04-03T00:00:00Z"
}
```

## Validation

Each model includes validation logic to ensure data integrity:

### Strategy Validation
- Name is required
- User ID is required
- Strategy type is required
- At least one instrument is required
- Risk parameters must be valid

### StrategySchedule Validation
- Strategy ID is required
- Frequency is required
- Start time is required
- For weekly frequency, days of week are required
- Days of week must be valid (0-6)

## Error Handling
The API implements consistent error handling:
- Validation errors return 400 Bad Request with descriptive messages
- Not found errors return 404 Not Found
- Server errors return 500 Internal Server Error
- Custom error responses include an "error" field with a descriptive message

## Security Considerations
The Strategies Management Backend implements several security measures:
- User ID validation to prevent unauthorized access
- Input validation to prevent injection attacks
- Proper error handling to avoid information leakage
- Status checks to prevent invalid operations

## Usage Examples

### Creating a Strategy
```
POST /api/strategies
Content-Type: application/json

{
  "name": "Moving Average Crossover",
  "description": "A strategy that trades based on moving average crossovers",
  "userId": "user123",
  "type": "AUTOMATED",
  "entryConditions": [
    {
      "type": "INDICATOR",
      "parameter": "MA_CROSSOVER",
      "operator": "CROSSES_ABOVE",
      "value": {
        "fastPeriod": 9,
        "slowPeriod": 21
      }
    }
  ],
  "exitConditions": [
    {
      "type": "INDICATOR",
      "parameter": "MA_CROSSOVER",
      "operator": "CROSSES_BELOW",
      "value": {
        "fastPeriod": 9,
        "slowPeriod": 21
      }
    }
  ],
  "riskParameters": {
    "maxPositionSize": 100,
    "maxLoss": 1000,
    "maxDailyLoss": 2000,
    "trailingStopPercent": 2.0,
    "takeProfitPercent": 5.0
  },
  "instruments": ["AAPL", "MSFT", "GOOGL"],
  "tags": ["momentum", "technical", "stocks"]
}
```

### Executing a Strategy
```
POST /api/strategies/strategy123/execute
```

### Scheduling a Strategy
```
POST /api/strategies/strategy123/schedule
Content-Type: application/json

{
  "frequency": "DAILY",
  "startTime": "2025-04-03T09:30:00Z",
  "endTime": "2025-04-03T16:00:00Z",
  "daysOfWeek": [1, 2, 3, 4, 5],
  "enabled": true
}
```

### Getting Strategy Performance
```
GET /api/strategies/strategy123/performance
```

### Adding a Tag to a Strategy
```
POST /api/strategies/strategy123/tags
Content-Type: application/json

{
  "tag": "high-frequency"
}
```

## Testing
Comprehensive unit tests have been implemented for all components:
- Service tests with mock repositories
- Repository tests with mock MongoDB interfaces
- Handler tests with mock services

The tests cover:
- Happy path scenarios
- Error handling
- Edge cases
- Validation logic

## Future Enhancements
Potential future enhancements for the Strategies Management Backend include:
- Advanced condition types for more complex strategies
- Machine learning integration for strategy optimization
- Strategy backtesting capabilities
- Strategy sharing and collaboration features
- Strategy templates and cloning
- Advanced performance analytics
- Real-time strategy monitoring dashboard
- Strategy comparison tools
- Integration with external data sources for more sophisticated conditions

## Implementation Notes
- All components follow a clean, layered architecture
- Proper separation of concerns between layers
- Comprehensive validation at all levels
- Thorough error handling
- MongoDB integration for data persistence
- RESTful API design principles
- Concurrent execution and monitoring of strategies
