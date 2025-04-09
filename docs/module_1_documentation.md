# Module 1: Data Models and Database Schema Documentation

## Overview
This document provides detailed information about the data models and database schema implemented in the trading platform. The models are designed to support a comprehensive trading system with options, futures, and stock trading capabilities.

## Models

### Order Model
The Order model represents a trading order in the system with the following key features:
- Support for different order types (Market, Limit, Stop-Loss Limit)
- Support for different order directions (Buy, Sell)
- Support for different product types (MIS, NRML, CNC)
- Support for different instrument types (Option, Future, Stock)
- Comprehensive validation logic for all fields
- Calculation methods for slippage, completion status, and remaining quantity

### Position Model
The Position model represents a trading position with the following key features:
- Support for different position statuses (Open, Closed, Partial)
- Support for different instrument types (Option, Future, Stock)
- Calculation methods for P&L, Greeks, and days till expiry
- Support for partial position closure
- Comprehensive validation logic for all fields

### User Model
The User model represents a user in the trading system with the following key features:
- Different user roles (Admin, Trader, Viewer)
- Secure password handling
- Two-factor authentication support
- User preferences for trading settings
- API key management for broker integration
- Comprehensive validation logic for all fields

### Strategy Model
The Strategy model represents a trading strategy with the following key features:
- Different strategy types (Directional, Market Neutral, Volatility, etc.)
- Different execution modes (Time, Signal, Combined Premium, etc.)
- Target and stop-loss settings
- Risk management parameters
- Performance metrics tracking
- Comprehensive validation logic for all fields

### Portfolio Model
The Portfolio model represents a multi-leg options portfolio with the following key features:
- Support for different portfolio statuses (Pending, Active, Completed, Failed)
- Strike selection modes for options
- Execution parameters for entry and exit
- Target and stop-loss settings
- Performance metrics tracking
- Comprehensive validation logic for all fields

### Leg Model
The Leg model represents a single leg in a multi-leg options portfolio with the following key features:
- Support for different leg types (Option, Future, Stock)
- Entry and exit parameters
- Individual target and stop-loss settings
- Greeks calculation
- Performance metrics tracking
- Comprehensive validation logic for all fields

## Validation Logic
All models include comprehensive validation logic to ensure data integrity:
- Required field validation
- Format validation for fields like email, phone, time formats
- Range validation for numeric fields
- Consistency validation between related fields
- Type validation for enumerated types

## Calculation Methods
The models include various calculation methods:
- P&L calculations (unrealized, realized, total, percentage)
- Greeks calculations (delta, gamma, theta, vega)
- Slippage calculations
- Status determination methods
- Time-based calculations (days till expiry, should execute now, etc.)

## Database Schema
The models are designed to be stored in a MongoDB database with the following schema considerations:
- Appropriate field tagging for JSON and BSON serialization
- Indexed fields for efficient querying
- Embedded documents for related data
- References to related documents using ID fields

## Testing
Comprehensive unit tests have been implemented for all models, covering:
- Validation logic for all fields
- Calculation methods
- Edge cases and error conditions

## Usage Examples
```go
// Create a new order
order := &models.Order{
    UserID:         "user123",
    Symbol:         "NIFTY",
    Exchange:       "NSE",
    OrderType:      models.OrderTypeLimit,
    Direction:      models.OrderDirectionBuy,
    Quantity:       10,
    Price:          500.50,
    Status:         models.OrderStatusPending,
    ProductType:    models.ProductTypeMIS,
    InstrumentType: models.InstrumentTypeOption,
    OptionType:     models.OptionTypeCall,
    StrikePrice:    18000,
    Expiry:         time.Now().AddDate(0, 1, 0),
}

// Validate the order
if err := order.Validate(); err != nil {
    // Handle validation error
}

// Calculate slippage
slippage := order.CalculateSlippage()

// Check if order is complete
isComplete := order.IsComplete()

// Get remaining quantity
remaining := order.RemainingQuantity()
```

## Future Enhancements
Potential future enhancements for the data models include:
- Support for more complex order types (OCO, bracket orders)
- Enhanced risk management parameters
- Additional performance metrics
- Support for different asset classes
- Integration with more brokers

## Implementation Notes
- All models have been implemented with Go struct tags for JSON and BSON serialization
- Validation logic is implemented in the Validate() method of each model
- Calculation methods are implemented as methods on the model structs
- Constants are used for enumerated types to ensure type safety
