# Module 16: SIM User Management

## Overview
Module 16 implements a comprehensive SIM User Management system that restricts simulation and backtesting activities to dedicated SIM users, ensuring complete isolation between simulation/backtesting and live trading environments. The implementation follows a phased approach to ensure stability and proper isolation.

## Implementation Phases

### Phase 1: Core SIM User Infrastructure
- Added user type system with STANDARD, ADMIN, and SIM user types
- Implemented environment context in authentication and authorization
- Created environment-aware middleware for enforcing isolation
- Developed frontend components for environment indication
- Implemented comprehensive test suite for SIM user infrastructure

### Phase 2: Execution Platform Testing with SIM Capabilities
- Created test scenarios for order creation with SIM users
- Implemented tests for position tracking with SIM users
- Developed tests for strategy execution with SIM users
- Verified environment isolation in an integrated environment
- Validated proper environment context preservation

### Phase 3: Isolated Simulation/Backtesting System
- Created simulation account data models
- Implemented virtual balance management
- Developed simulation order processing engine
- Created market simulation engine with realistic conditions
- Implemented backtesting service for historical data analysis

### Phase 4: Controlled Interfaces Between Systems
- Defined interface contracts between execution platform and simulation system
- Implemented API Gateway for cross-system communication
- Created data synchronization mechanisms
- Developed error handling across system boundaries
- Implemented security controls for cross-system access

## Architecture

### User Type System
The user type system is implemented as an enumeration with three values:
- `STANDARD`: Regular users with access to live trading
- `ADMIN`: Administrative users with full system access
- `SIM`: Simulation users restricted to simulation/backtesting environments

```go
// UserType defines the type of user in the system
type UserType string

const (
    // UserTypeStandard represents a standard user with access to live trading
    UserTypeStandard UserType = "STANDARD"
    
    // UserTypeAdmin represents an administrative user with full system access
    UserTypeAdmin UserType = "ADMIN"
    
    // UserTypeSim represents a simulation user restricted to simulation/backtesting
    UserTypeSim UserType = "SIM"
)
```

### Environment Context
The environment context is implemented as an enumeration with two values:
- `LIVE`: Live trading environment with real money
- `SIM`: Simulation environment with virtual money

```go
// Environment defines the execution environment
type Environment string

const (
    // EnvironmentLive represents the live trading environment
    EnvironmentLive Environment = "LIVE"
    
    // EnvironmentSim represents the simulation environment
    EnvironmentSim Environment = "SIM"
)
```

### Authentication and Authorization
The JWT token structure has been extended to include user type and environment information:

```go
// Claims represents the JWT claims structure
type Claims struct {
    UserID      string      `json:"userId"`
    Username    string      `json:"username"`
    UserType    string      `json:"userType"`
    Environment string      `json:"environment"`
    Permissions []string    `json:"permissions"`
    ExpiresAt   int64       `json:"exp"`
    IssuedAt    int64       `json:"iat"`
}
```

The authentication middleware validates the user type and environment, ensuring that:
- SIM users can only access the SIM environment
- STANDARD users can access both LIVE and SIM environments
- ADMIN users can access both environments

### Simulation Account Model
The simulation account model extends the standard account model with simulation-specific fields:

```go
// SimulationAccount represents a simulation/paper trading account
type SimulationAccount struct {
    ID                string      `json:"id"`
    UserID            string      `json:"userId"`
    Name              string      `json:"name"`
    InitialBalance    float64     `json:"initialBalance"`
    CurrentBalance    float64     `json:"currentBalance"`
    Currency          string      `json:"currency"`
    SimulationType    string      `json:"simulationType"` // "PAPER" or "BACKTEST"
    CreatedAt         time.Time   `json:"createdAt"`
    UpdatedAt         time.Time   `json:"updatedAt"`
    IsActive          bool        `json:"isActive"`
    RiskSettings      RiskSettings `json:"riskSettings"`
    MarketSettings    MarketSettings `json:"marketSettings"`
}
```

### API Gateway
The API Gateway serves as the central communication point between the execution platform and simulation system, enforcing security, handling cross-system errors, managing rate limiting, and ensuring proper data synchronization.

```go
// APIGateway implements the interfaces.ExecutionSimulationInterface
type APIGateway struct {
    simulationService       *simulation.SimulationAccountService
    virtualBalanceService   *simulation.VirtualBalanceService
    simulationOrderService  *simulation.SimulationOrderService
    marketSimulationService *simulation.MarketSimulationService
    backtestService         *simulation.BacktestService
    
    // Execution platform interface
    executionPlatform       interfaces.ExecutionPlatformInterface
    
    // Security and rate limiting
    accessControlList       map[string][]string // userID -> permissions
    rateLimits              map[string]RateLimit
    
    // Data synchronization
    lastSyncTime            map[string]time.Time // resource -> last sync time
    
    // Error handling
    errorHandlers           map[string]ErrorHandler
}
```

### Interface Contracts
The interface contracts define the communication boundaries between the execution platform and simulation system:

```go
// ExecutionSimulationInterface defines the contract between the execution platform
// and the simulation system
type ExecutionSimulationInterface interface {
    // Account Management
    CreateSimulationAccount(ctx context.Context, userID string, account models.SimulationAccount) (*models.SimulationAccount, error)
    GetSimulationAccount(ctx context.Context, accountID string) (*models.SimulationAccount, error)
    // ... additional methods
}

// ExecutionPlatformInterface defines the contract for the execution platform
// to expose functionality to the simulation system
type ExecutionPlatformInterface interface {
    // Market Data Access
    GetRealTimeMarketData(ctx context.Context, symbol string) (*models.MarketDataSnapshot, error)
    GetHistoricalMarketData(ctx context.Context, symbol string, startDate, endDate time.Time, timeframe string) ([]*models.MarketDataSnapshot, error)
    // ... additional methods
}
```

## Security Features

### Permission System
The permission system uses a fine-grained approach with resource:action format:
- `simulation:account:read`: Permission to read simulation accounts
- `simulation:order:create`: Permission to create simulation orders
- `backtest:session:run`: Permission to run backtest sessions

### Rate Limiting
Rate limiting is implemented per user and per API category:
- Market data: 300 requests per minute
- Order management: 100 requests per minute
- Account management: 60 requests per minute
- Backtesting: 30 requests per minute

### Error Handling
Cross-system error handling ensures that:
- Authentication errors return standardized messages without exposing details
- Authorization errors return standardized messages without exposing details
- Validation errors return detailed messages to help users fix issues
- System errors return generic messages without exposing internal details

### Data Synchronization
Data synchronization ensures that:
- Market data is synchronized between execution platform and simulation system
- Synchronization is cached for 5 minutes to reduce load
- Critical operations like order creation always trigger synchronization

## Frontend Components

### Environment Indicator
The Environment Indicator component provides a visual indication of the current environment:
- SIM environment: Yellow indicator with "SIMULATION" text
- LIVE environment: Green indicator with "LIVE TRADING" text

```tsx
const EnvironmentIndicator: React.FC<EnvironmentIndicatorProps> = ({ environment }) => {
  const getIndicatorClass = () => {
    return environment === 'SIM' ? 'simulation-indicator' : 'live-indicator';
  };

  const getIndicatorText = () => {
    return environment === 'SIM' ? 'SIMULATION' : 'LIVE TRADING';
  };

  return (
    <div className={`environment-indicator ${getIndicatorClass()}`}>
      {getIndicatorText()}
    </div>
  );
};
```

### User Context
The User Context provides environment state management across the application:

```tsx
export const UserContext = createContext<UserContextType>({
  user: null,
  environment: 'LIVE',
  setEnvironment: () => {},
  isLoading: false,
  error: null,
});

export const UserProvider: React.FC<UserProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [environment, setEnvironment] = useState<Environment>('LIVE');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // ... implementation details

  return (
    <UserContext.Provider value={{ user, environment, setEnvironment, isLoading, error }}>
      {children}
    </UserContext.Provider>
  );
};
```

## Testing

### Unit Tests
Unit tests verify the functionality of individual components:
- User type and environment models
- Authentication and authorization middleware
- Environment service
- API Gateway components

### Integration Tests
Integration tests verify the interaction between components:
- API Gateway with HTTP handlers
- Cross-system error handling
- Data synchronization between systems

### Security Tests
Security tests verify the isolation between environments:
- SIM users cannot access LIVE resources
- LIVE users cannot access SIM resources
- Admin users can access all resources
- Permission checking works correctly
- Rate limiting prevents abuse

## Usage Examples

### Creating a Simulation Account
```go
// Create a simulation account
account := models.SimulationAccount{
    Name:           "Test Paper Trading",
    InitialBalance: 100000.0,
    Currency:       "USD",
    SimulationType: "PAPER",
    RiskSettings: models.RiskSettings{
        MaxPositionSize:      0.05, // 5% of account
        MaxDrawdown:          0.20, // 20% drawdown
        DailyLossLimit:       0.05, // 5% daily loss limit
    },
    MarketSettings: models.MarketSettings{
        Slippage:             0.001, // 0.1% slippage
        Latency:              200,   // 200ms latency
        CommissionRate:       0.0025, // 0.25% commission
    },
}

// Create the account through the API Gateway
result, err := gateway.CreateSimulationAccount(ctx, "user123", account)
```

### Creating a Simulation Order
```go
// Create a simulation order
order := models.SimulationOrder{
    Order: models.Order{
        Symbol:    "AAPL",
        Quantity:  100,
        Side:      "BUY",
        OrderType: "MARKET",
    },
}

// Create the order through the API Gateway
result, err := gateway.CreateOrder(ctx, "sim123", order)
```

### Running a Backtest
```go
// Create a backtest session
session := models.BacktestSession{
    Name:           "AAPL Strategy Backtest",
    StartDate:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
    EndDate:        time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
    Symbols:        []string{"AAPL"},
    InitialBalance: 100000.0,
    Strategy:       "moving_average_crossover",
    Parameters: map[string]interface{}{
        "short_period": 10,
        "long_period":  30,
    },
}

// Create the backtest session
result, err := gateway.CreateBacktestSession(ctx, "sim123", session)

// Run the backtest
err = gateway.RunBacktest(ctx, result.ID)

// Get backtest results
results, err := gateway.GetBacktestResults(ctx, result.ID)

// Get performance metrics
metrics, err := gateway.GetBacktestPerformanceMetrics(ctx, result.ID)
```

## Conclusion
Module 16: SIM User Management provides a comprehensive solution for simulation and backtesting activities within the trading platform. By implementing a robust user type system, environment context, and controlled interfaces between systems, it ensures complete isolation between simulation/backtesting and live trading environments while providing a seamless user experience.

The phased implementation approach ensures stability and proper isolation, with each phase building on the previous one to create a complete solution. The comprehensive testing suite verifies that all components work correctly and that the isolation between environments is maintained.
