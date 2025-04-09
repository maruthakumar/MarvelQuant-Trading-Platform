# Paper Trading Workflow Documentation

## Overview

This document outlines the complete workflow for the paper trading environment with SIM user functionality in the trading platform. Paper trading provides a risk-free environment for users to test trading strategies, learn platform features, and practice trading without using real capital.

## User Workflows

### 1. SIM User Account Creation

#### Process Flow

1. **User Registration**
   - User navigates to registration page
   - User selects "Simulation Account" option
   - User completes registration form with required information
   - System validates information and creates SIM user account
   - System generates welcome email with simulation account details

2. **Simulation Account Setup**
   - User logs in with SIM user credentials
   - System displays simulation account setup form
   - User configures initial account parameters:
     - Initial balance (default: $100,000)
     - Commission model (none, realistic, custom)
     - Slippage model (none, realistic, custom)
   - System creates simulation account with specified parameters
   - System redirects to simulation dashboard

#### UI Components

- Registration form with account type selection
- Simulation account setup form
- Simulation mode indicator (displayed throughout the UI)

#### Database Operations

- Create user record with `user_type = 'SIM'`
- Create simulation account record with initial parameters
- Initialize account balance

### 2. Switching Between Real and Simulation Modes

#### Process Flow

1. **For Users with Both Account Types**
   - User clicks on account switcher in navigation bar
   - System displays available account types (REAL, SIM)
   - User selects desired account type
   - System switches context to selected account type
   - UI updates to reflect selected account type (including visual indicators for simulation mode)

2. **Session Management**
   - System maintains separate session data for each account type
   - System preserves state when switching between account types
   - System applies appropriate permissions and restrictions based on account type

#### UI Components

- Account type switcher in navigation bar
- Visual indicators for simulation mode
- Confirmation dialog when switching account types

#### Security Considerations

- Clear visual distinction between real and simulation modes
- Confirmation prompts when switching modes
- Session isolation between real and simulation accounts

### 3. Paper Trading Order Execution

#### Process Flow

1. **Order Creation**
   - User navigates to order entry form
   - System displays simulation mode indicator
   - User enters order details:
     - Instrument
     - Order type (market, limit, stop, etc.)
     - Side (buy, sell)
     - Quantity
     - Price (for limit orders)
     - Additional parameters
   - User submits order
   - System validates order against simulation account parameters (balance, position limits)

2. **Order Simulation**
   - System routes order to paper trading adapter
   - Paper trading adapter processes order:
     - Applies configured execution delay
     - Calculates execution price based on slippage model
     - Updates simulation account balance
     - Creates simulated position or updates existing position
     - Records transaction in simulation history

3. **Order Confirmation**
   - System displays order confirmation with simulation indicator
   - System updates orderbook with simulated order
   - System sends notification of order execution (clearly marked as simulation)

#### UI Components

- Order entry form with simulation indicator
- Order confirmation dialog with simulation indicator
- Orderbook with simulation order highlighting
- Notifications with simulation indicators

#### Simulation Parameters

- **Execution Delay**: Configurable delay to simulate real-world latency
- **Slippage Model**: Calculation of price deviation from requested price
- **Commission Model**: Calculation of transaction costs

### 4. Position Management in Simulation

#### Process Flow

1. **Position Creation**
   - System creates position when order is executed
   - System calculates average price for position
   - System updates simulation account balance

2. **Position Monitoring**
   - System updates position P&L in real-time based on market data
   - System displays position details in positions view
   - System calculates and displays risk metrics (Greeks for options)

3. **Position Closing**
   - User initiates square-off action for position
   - System creates opposite order to close position
   - System simulates order execution
   - System updates position status and P&L
   - System records transaction in simulation history

#### UI Components

- Positions view with simulation indicators
- Position details panel
- Square-off controls
- P&L visualization

#### Calculations

- **Unrealized P&L**: `(Current Price - Average Price) * Quantity * Multiplier`
- **Realized P&L**: Calculated when position is closed
- **Risk Metrics**: Greeks calculations for options positions

### 5. Performance Tracking and Reporting

#### Process Flow

1. **Real-time Performance Monitoring**
   - System calculates and displays key performance metrics:
     - Current balance
     - Unrealized P&L
     - Realized P&L
     - Win/loss ratio
     - Sharpe ratio
   - System updates metrics in real-time as market prices change

2. **Historical Performance Analysis**
   - User selects date range for analysis
   - System retrieves historical simulation data
   - System calculates performance metrics for selected period
   - System generates performance charts and visualizations

3. **Report Generation**
   - User requests performance report
   - System compiles simulation data
   - System generates report with performance metrics, trade history, and analysis
   - System provides export options (PDF, CSV)

#### UI Components

- Performance dashboard
- Equity curve chart
- Trade history table
- Performance metrics cards
- Report configuration form

#### Metrics Calculated

- **Total P&L**: Sum of realized and unrealized P&L
- **Win Rate**: Percentage of profitable trades
- **Sharpe Ratio**: Risk-adjusted return metric
- **Maximum Drawdown**: Largest peak-to-trough decline
- **Return on Investment**: Percentage return on initial capital

### 6. Strategy Testing in Simulation

#### Process Flow

1. **Strategy Configuration**
   - User navigates to strategy creation interface
   - User configures strategy parameters
   - User selects simulation account for strategy execution
   - System validates strategy configuration

2. **Strategy Execution**
   - System executes strategy in simulation environment
   - System generates orders based on strategy rules
   - Paper trading adapter processes orders
   - System tracks strategy performance

3. **Strategy Monitoring**
   - User monitors strategy execution in real-time
   - System displays strategy metrics and performance
   - User can pause, resume, or stop strategy execution
   - System provides alerts for significant events

#### UI Components

- Strategy configuration interface
- Strategy monitoring dashboard
- Strategy performance metrics
- Strategy control panel

#### Strategy Types

- **Single-leg Strategies**: Simple directional strategies
- **Multi-leg Strategies**: Complex option strategies
- **Portfolio Strategies**: Multiple instruments with allocation rules
- **MAXOI/MAXCOI Strategies**: Open interest based strategies

## Technical Workflows

### 1. Simulation Engine Architecture

#### Components

1. **Paper Trading Adapter**
   - Implements broker adapter interface
   - Routes simulation requests to appropriate components
   - Maintains isolation from real trading systems

2. **Market Data Simulator**
   - Provides market data for simulation
   - Uses real market data with optional adjustments
   - Supports real-time and historical data modes

3. **Order Matching Engine**
   - Simulates order execution
   - Applies slippage and execution delay models
   - Processes different order types (market, limit, stop)

4. **Position Keeper**
   - Tracks simulated positions
   - Calculates position metrics
   - Manages position lifecycle

5. **Account Manager**
   - Manages simulation account balances
   - Processes transactions
   - Enforces account limits and rules

6. **Performance Tracker**
   - Calculates performance metrics
   - Generates performance reports
   - Provides historical analysis

#### Data Flow

```
User Interface → API Gateway → Paper Trading Adapter → Order Matching Engine → Position Keeper → Account Manager
                                      ↑                        ↓
                                      |                        |
                            Market Data Simulator  ←→  Performance Tracker
```

### 2. Database Schema

#### Tables

1. **users**
   - Contains user records with `user_type` field to distinguish SIM users
   - Stores authentication and profile information

2. **sim_accounts**
   - Stores simulation account configuration
   - Tracks initial and current balance
   - Links to user record

3. **sim_positions**
   - Tracks simulated positions
   - Stores quantity, average price, and P&L information
   - Links to simulation account

4. **sim_orders**
   - Records all simulated orders
   - Stores order parameters and execution details
   - Tracks order status and history

5. **sim_transactions**
   - Records all financial transactions in simulation accounts
   - Tracks balance changes, commissions, and fees
   - Provides audit trail for simulation activity

#### Relationships

- One user can have multiple simulation accounts
- One simulation account can have multiple positions
- One simulation account can have multiple orders
- One simulation account can have multiple transactions
- Orders and positions have a many-to-many relationship

### 3. Integration Points

#### Authentication System

- Extended to support SIM user type
- Provides session management with simulation flag
- Enforces appropriate permissions for simulation accounts

#### Order Execution Engine

- Routes orders based on user type
- Uses paper trading adapter for simulation orders
- Maintains complete isolation between real and simulated orders

#### Market Data Services

- Provides real market data to simulation engine
- Supports optional adjustments for testing scenarios
- Ensures consistent data between real and simulation environments

#### UI Framework

- Displays clear visual indicators for simulation mode
- Provides simulation-specific controls and forms
- Ensures consistent user experience across real and simulation modes

## Security and Isolation

### 1. Data Isolation

- Separate database tables for simulation data
- Logical separation of real and simulated orders
- Distinct API endpoints for simulation functionality

### 2. Visual Indicators

- Prominent simulation mode indicator throughout UI
- Color-coding of simulation elements (blue for simulation)
- Clear labeling of all simulation-related components

### 3. Permission Controls

- Specific permissions for simulation functionality
- Restricted access to real trading for SIM-only users
- Controlled access to simulation configuration

### 4. Validation Checks

- Multiple validation layers to prevent crossover between real and simulation
- Confirmation prompts for critical actions in simulation
- System checks to verify appropriate routing of orders

## Error Handling

### 1. Simulation-Specific Errors

- Insufficient simulation balance
- Invalid simulation configuration
- Simulation engine failures

### 2. Error Responses

- Clear error messages indicating simulation context
- Appropriate error codes for simulation issues
- Detailed logging for troubleshooting

### 3. Recovery Procedures

- Automatic retry for non-critical simulation failures
- Manual intervention options for critical failures
- Simulation state recovery mechanisms

## Conclusion

The paper trading environment provides a comprehensive, risk-free platform for users to test trading strategies, learn platform features, and practice trading without using real capital. The workflows outlined in this document ensure a seamless, intuitive user experience while maintaining complete isolation between real and simulated trading activities.

The implementation follows a modular approach with clear separation of concerns, allowing for future enhancements and extensions to the simulation capabilities. The system is designed to provide a realistic trading experience while clearly indicating the simulation nature of all activities to prevent confusion with real trading.
