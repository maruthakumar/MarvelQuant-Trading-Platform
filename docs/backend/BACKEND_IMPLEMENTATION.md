# Backend Implementation Documentation

## Overview

This document provides comprehensive documentation for the backend implementation of the trading platform. The backend is built using Go with a RESTful API architecture and MongoDB for data storage. It implements core trading functionality including user authentication, order management, strategy management, and portfolio management.

## Architecture

The backend follows a clean architecture pattern with the following layers:

1. **API Layer**: Handles HTTP requests and responses
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Manages data persistence
4. **Model Layer**: Defines data structures

### Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/                  # API handlers
│   │   ├── order_handler.go
│   │   ├── position_handler.go
│   │   ├── portfolio_handler.go
│   │   ├── strategy_handler.go
│   │   ├── user_handler.go
│   │   └── routes.go
│   ├── auth/                 # Authentication
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── config/               # Configuration
│   │   └── config.go
│   ├── database/             # Database connection
│   │   └── mongodb.go
│   ├── models/               # Data models
│   │   ├── order.go
│   │   ├── position.go
│   │   ├── portfolio.go
│   │   ├── strategy.go
│   │   ├── user.go
│   │   └── leg.go
│   ├── repository/           # Data repositories
│   │   ├── order_repository.go
│   │   ├── position_repository.go
│   │   ├── portfolio_repository.go
│   │   ├── strategy_repository.go
│   │   └── user_repository.go
│   └── utils/                # Utilities
│       └── http.go
└── tests/                    # Unit tests
    ├── order_handler_test.go
    ├── portfolio_handler_test.go
    ├── strategy_handler_test.go
    └── user_handler_test.go
```

## Data Models

### User Model

```go
type User struct {
    ID             string    `json:"id" bson:"_id,omitempty"`
    Email          string    `json:"email" bson:"email"`
    PasswordHash   string    `json:"-" bson:"password_hash"`
    FirstName      string    `json:"firstName" bson:"first_name"`
    LastName       string    `json:"lastName" bson:"last_name"`
    Role           string    `json:"role" bson:"role"`
    Status         string    `json:"status" bson:"status"`
    EmailVerified  bool      `json:"emailVerified" bson:"email_verified"`
    Preferences    UserPreferences `json:"preferences" bson:"preferences"`
    CreatedAt      time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt      time.Time `json:"updatedAt" bson:"updated_at"`
    LastLoginAt    time.Time `json:"lastLoginAt,omitempty" bson:"last_login_at,omitempty"`
}

type UserPreferences struct {
    Theme          string    `json:"theme" bson:"theme"`
    Language       string    `json:"language" bson:"language"`
    Notifications  bool      `json:"notifications" bson:"notifications"`
    TwoFactorAuth  bool      `json:"twoFactorAuth" bson:"two_factor_auth"`
}
```

### Order Model

```go
type Order struct {
    ID            string    `json:"id" bson:"_id,omitempty"`
    UserID        string    `json:"userId" bson:"user_id"`
    PortfolioID   string    `json:"portfolioId,omitempty" bson:"portfolio_id,omitempty"`
    StrategyID    string    `json:"strategyId,omitempty" bson:"strategy_id,omitempty"`
    Symbol        string    `json:"symbol" bson:"symbol"`
    OrderType     OrderType `json:"orderType" bson:"order_type"`
    Side          Side      `json:"side" bson:"side"`
    Quantity      float64   `json:"quantity" bson:"quantity"`
    Price         float64   `json:"price,omitempty" bson:"price,omitempty"`
    StopPrice     float64   `json:"stopPrice,omitempty" bson:"stop_price,omitempty"`
    TimeInForce   TimeInForce `json:"timeInForce" bson:"time_in_force"`
    Status        OrderStatus `json:"status" bson:"status"`
    FilledQty     float64   `json:"filledQty" bson:"filled_qty"`
    AvgFillPrice  float64   `json:"avgFillPrice,omitempty" bson:"avg_fill_price,omitempty"`
    RejectReason  string    `json:"rejectReason,omitempty" bson:"reject_reason,omitempty"`
    Notes         string    `json:"notes,omitempty" bson:"notes,omitempty"`
    CreatedAt     time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt     time.Time `json:"updatedAt" bson:"updated_at"`
    ExecutedAt    time.Time `json:"executedAt,omitempty" bson:"executed_at,omitempty"`
}

type OrderType string
const (
    OrderTypeMarket OrderType = "market"
    OrderTypeLimit  OrderType = "limit"
    OrderTypeStop   OrderType = "stop"
    OrderTypeStopLimit OrderType = "stop_limit"
)

type Side string
const (
    SideBuy  Side = "buy"
    SideSell Side = "sell"
)

type TimeInForce string
const (
    TimeInForceGTC TimeInForce = "gtc" // Good Till Cancelled
    TimeInForceIOC TimeInForce = "ioc" // Immediate or Cancel
    TimeInForceFOK TimeInForce = "fok" // Fill or Kill
    TimeInForceDAY TimeInForce = "day" // Day Order
)

type OrderStatus string
const (
    OrderStatusNew       OrderStatus = "new"
    OrderStatusPartiallyFilled OrderStatus = "partially_filled"
    OrderStatusFilled    OrderStatus = "filled"
    OrderStatusCancelled OrderStatus = "cancelled"
    OrderStatusRejected  OrderStatus = "rejected"
    OrderStatusExpired   OrderStatus = "expired"
)

type OrderFilter struct {
    UserID      string      `json:"userId,omitempty"`
    PortfolioID string      `json:"portfolioId,omitempty"`
    StrategyID  string      `json:"strategyId,omitempty"`
    Symbol      string      `json:"symbol,omitempty"`
    Status      OrderStatus `json:"status,omitempty"`
    Side        Side        `json:"side,omitempty"`
    StartDate   time.Time   `json:"startDate,omitempty"`
    EndDate     time.Time   `json:"endDate,omitempty"`
}
```

### Position Model

```go
type Position struct {
    ID            string    `json:"id" bson:"_id,omitempty"`
    UserID        string    `json:"userId" bson:"user_id"`
    PortfolioID   string    `json:"portfolioId,omitempty" bson:"portfolio_id,omitempty"`
    Symbol        string    `json:"symbol" bson:"symbol"`
    Direction     Direction `json:"direction" bson:"direction"`
    Quantity      float64   `json:"quantity" bson:"quantity"`
    EntryPrice    float64   `json:"entryPrice" bson:"entry_price"`
    CurrentPrice  float64   `json:"currentPrice" bson:"current_price"`
    UnrealizedPnL float64   `json:"unrealizedPnL" bson:"unrealized_pnl"`
    RealizedPnL   float64   `json:"realizedPnL" bson:"realized_pnl"`
    CreatedAt     time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt     time.Time `json:"updatedAt" bson:"updated_at"`
}

type Direction string
const (
    DirectionLong  Direction = "long"
    DirectionShort Direction = "short"
)

type PositionFilter struct {
    UserID      string    `json:"userId,omitempty"`
    PortfolioID string    `json:"portfolioId,omitempty"`
    Symbol      string    `json:"symbol,omitempty"`
    Direction   Direction `json:"direction,omitempty"`
}
```

### Strategy Model

```go
type Strategy struct {
    ID          string                 `json:"id" bson:"_id,omitempty"`
    UserID      string                 `json:"userId" bson:"user_id"`
    Name        string                 `json:"name" bson:"name"`
    Description string                 `json:"description" bson:"description"`
    Type        string                 `json:"type" bson:"type"`
    Symbol      string                 `json:"symbol" bson:"symbol"`
    ProductType string                 `json:"productType" bson:"product_type"`
    Parameters  map[string]interface{} `json:"parameters" bson:"parameters"`
    Active      bool                   `json:"active" bson:"active"`
    Tags        []string               `json:"tags" bson:"tags"`
    CreatedAt   time.Time              `json:"createdAt" bson:"created_at"`
    UpdatedAt   time.Time              `json:"updatedAt" bson:"updated_at"`
}

type StrategyFilter struct {
    UserID      string `json:"userId,omitempty"`
    Type        string `json:"type,omitempty"`
    Symbol      string `json:"symbol,omitempty"`
    ProductType string `json:"productType,omitempty"`
    Active      *bool  `json:"active,omitempty"`
    Tag         string `json:"tag,omitempty"`
}
```

### Portfolio Model

```go
type Portfolio struct {
    ID          string           `json:"id" bson:"_id,omitempty"`
    UserID      string           `json:"userId" bson:"user_id"`
    Name        string           `json:"name" bson:"name"`
    Description string           `json:"description" bson:"description"`
    StrategyID  string           `json:"strategyId" bson:"strategy_id"`
    Status      PortfolioStatus  `json:"status" bson:"status"`
    Capital     float64          `json:"capital" bson:"capital"`
    Currency    string           `json:"currency" bson:"currency"`
    RiskLevel   string           `json:"riskLevel" bson:"risk_level"`
    Legs        []Leg            `json:"legs" bson:"legs"`
    CreatedAt   time.Time        `json:"createdAt" bson:"created_at"`
    UpdatedAt   time.Time        `json:"updatedAt" bson:"updated_at"`
}

type PortfolioStatus string
const (
    PortfolioStatusDraft    PortfolioStatus = "draft"
    PortfolioStatusActive   PortfolioStatus = "active"
    PortfolioStatusInactive PortfolioStatus = "inactive"
    PortfolioStatusClosed   PortfolioStatus = "closed"
)

type PortfolioFilter struct {
    UserID     string          `json:"userId,omitempty"`
    StrategyID string          `json:"strategyId,omitempty"`
    Status     PortfolioStatus `json:"status,omitempty"`
    RiskLevel  string          `json:"riskLevel,omitempty"`
}
```

### Leg Model

```go
type Leg struct {
    ID          int       `json:"id" bson:"id"`
    PortfolioID string    `json:"portfolioId" bson:"portfolio_id"`
    Symbol      string    `json:"symbol" bson:"symbol"`
    Direction   Direction `json:"direction" bson:"direction"`
    Allocation  float64   `json:"allocation" bson:"allocation"`
    EntryPrice  float64   `json:"entryPrice" bson:"entry_price"`
    StopLoss    float64   `json:"stopLoss,omitempty" bson:"stop_loss,omitempty"`
    TakeProfit  float64   `json:"takeProfit,omitempty" bson:"take_profit,omitempty"`
    Description string    `json:"description,omitempty" bson:"description,omitempty"`
    CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" bson:"updated_at"`
}
```

## API Endpoints

### Authentication Endpoints

#### Register User

- **URL**: `/api/auth/register`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123",
    "firstName": "John",
    "lastName": "Doe"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "user123",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "user",
    "status": "active",
    "emailVerified": false,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `201 Created`: User successfully registered
  - `400 Bad Request`: Invalid input
  - `409 Conflict`: Email already exists

#### Login

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123"
  }
  ```
- **Response**: 
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user123",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "user"
    }
  }
  ```
- **Status Codes**:
  - `200 OK`: Login successful
  - `401 Unauthorized`: Invalid credentials

#### Refresh Token

- **URL**: `/api/auth/refresh`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **Status Codes**:
  - `200 OK`: Token refreshed
  - `401 Unauthorized`: Invalid or expired token

#### Forgot Password

- **URL**: `/api/auth/forgot-password`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response**: 
  ```json
  {
    "message": "Password reset email sent"
  }
  ```
- **Status Codes**:
  - `200 OK`: Reset email sent
  - `404 Not Found`: Email not found

#### Reset Password

- **URL**: `/api/auth/reset-password`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "token": "reset-token-123",
    "password": "newSecurePassword123"
  }
  ```
- **Response**: 
  ```json
  {
    "message": "Password reset successful"
  }
  ```
- **Status Codes**:
  - `200 OK`: Password reset successful
  - `400 Bad Request`: Invalid token or password

### User Endpoints

#### Get User Profile

- **URL**: `/api/users/profile`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "user123",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "user",
    "status": "active",
    "emailVerified": true,
    "preferences": {
      "theme": "dark",
      "language": "en",
      "notifications": true,
      "twoFactorAuth": false
    },
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z",
    "lastLoginAt": "2023-01-02T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Profile retrieved
  - `401 Unauthorized`: Not authenticated

#### Update User Profile

- **URL**: `/api/users/profile`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "firstName": "Johnny",
    "lastName": "Doe",
    "preferences": {
      "theme": "light",
      "language": "fr"
    }
  }
  ```
- **Response**: 
  ```json
  {
    "id": "user123",
    "email": "user@example.com",
    "firstName": "Johnny",
    "lastName": "Doe",
    "preferences": {
      "theme": "light",
      "language": "fr",
      "notifications": true,
      "twoFactorAuth": false
    },
    "updatedAt": "2023-01-03T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Profile updated
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated

#### Change Password

- **URL**: `/api/users/change-password`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "currentPassword": "securePassword123",
    "newPassword": "evenMoreSecurePassword456"
  }
  ```
- **Response**: 
  ```json
  {
    "message": "Password changed successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Password changed
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Incorrect current password

### Order Endpoints

#### Create Order

- **URL**: `/api/orders`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "symbol": "AAPL",
    "orderType": "limit",
    "side": "buy",
    "quantity": 100,
    "price": 150.50,
    "timeInForce": "gtc",
    "portfolioId": "portfolio123"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "order123",
    "userId": "user123",
    "portfolioId": "portfolio123",
    "symbol": "AAPL",
    "orderType": "limit",
    "side": "buy",
    "quantity": 100,
    "price": 150.50,
    "timeInForce": "gtc",
    "status": "new",
    "filledQty": 0,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `201 Created`: Order created
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated

#### Get Order

- **URL**: `/api/orders/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "order123",
    "userId": "user123",
    "portfolioId": "portfolio123",
    "symbol": "AAPL",
    "orderType": "limit",
    "side": "buy",
    "quantity": 100,
    "price": 150.50,
    "timeInForce": "gtc",
    "status": "filled",
    "filledQty": 100,
    "avgFillPrice": 150.25,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T01:00:00Z",
    "executedAt": "2023-01-01T01:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Order retrieved
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to view this order
  - `404 Not Found`: Order not found

#### Update Order

- **URL**: `/api/orders/{id}`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "quantity": 200,
    "price": 155.75
  }
  ```
- **Response**: 
  ```json
  {
    "id": "order123",
    "userId": "user123",
    "portfolioId": "portfolio123",
    "symbol": "AAPL",
    "orderType": "limit",
    "side": "buy",
    "quantity": 200,
    "price": 155.75,
    "timeInForce": "gtc",
    "status": "new",
    "filledQty": 0,
    "updatedAt": "2023-01-01T02:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Order updated
  - `400 Bad Request`: Invalid input or order cannot be modified
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this order
  - `404 Not Found`: Order not found

#### Cancel Order

- **URL**: `/api/orders/{id}/cancel`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "order123",
    "userId": "user123",
    "status": "cancelled",
    "updatedAt": "2023-01-01T03:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Order cancelled
  - `400 Bad Request`: Order cannot be cancelled
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to cancel this order
  - `404 Not Found`: Order not found

#### Get Orders

- **URL**: `/api/orders?symbol=AAPL&status=filled&page=1&limit=20`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Query Parameters**:
  - `symbol`: Filter by symbol
  - `status`: Filter by status
  - `side`: Filter by side
  - `portfolioId`: Filter by portfolio ID
  - `startDate`: Filter by start date
  - `endDate`: Filter by end date
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 20)
- **Response**: 
  ```json
  {
    "data": [
      {
        "id": "order1",
        "userId": "user123",
        "symbol": "AAPL",
        "orderType": "limit",
        "side": "buy",
        "quantity": 100,
        "price": 150.50,
        "status": "filled"
      },
      {
        "id": "order2",
        "userId": "user123",
        "symbol": "MSFT",
        "orderType": "market",
        "side": "sell",
        "quantity": 50,
        "status": "filled"
      }
    ],
    "page": 1,
    "limit": 20,
    "total": 2,
    "totalPages": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Orders retrieved
  - `401 Unauthorized`: Not authenticated

### Strategy Endpoints

#### Create Strategy

- **URL**: `/api/strategies`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "name": "Test Strategy",
    "description": "A test strategy for automated trading",
    "type": "momentum",
    "symbol": "AAPL",
    "productType": "equity",
    "parameters": {
      "lookbackPeriod": 14,
      "threshold": 0.05
    },
    "active": true,
    "tags": ["test", "momentum", "equity"]
  }
  ```
- **Response**: 
  ```json
  {
    "id": "strategy123",
    "userId": "user123",
    "name": "Test Strategy",
    "description": "A test strategy for automated trading",
    "type": "momentum",
    "symbol": "AAPL",
    "productType": "equity",
    "parameters": {
      "lookbackPeriod": 14,
      "threshold": 0.05
    },
    "active": true,
    "tags": ["test", "momentum", "equity"],
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `201 Created`: Strategy created
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated

#### Get Strategy

- **URL**: `/api/strategies/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "strategy123",
    "userId": "user123",
    "name": "Test Strategy",
    "description": "A test strategy for automated trading",
    "type": "momentum",
    "symbol": "AAPL",
    "productType": "equity",
    "parameters": {
      "lookbackPeriod": 14,
      "threshold": 0.05
    },
    "active": true,
    "tags": ["test", "momentum", "equity"],
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategy retrieved
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to view this strategy
  - `404 Not Found`: Strategy not found

#### Update Strategy

- **URL**: `/api/strategies/{id}`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "name": "Updated Strategy",
    "description": "An updated test strategy",
    "parameters": {
      "lookbackPeriod": 21,
      "threshold": 0.03
    },
    "tags": ["test", "updated", "equity"]
  }
  ```
- **Response**: 
  ```json
  {
    "id": "strategy123",
    "userId": "user123",
    "name": "Updated Strategy",
    "description": "An updated test strategy",
    "type": "momentum",
    "symbol": "AAPL",
    "productType": "equity",
    "parameters": {
      "lookbackPeriod": 21,
      "threshold": 0.03
    },
    "active": true,
    "tags": ["test", "updated", "equity"],
    "updatedAt": "2023-01-01T01:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategy updated
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this strategy
  - `404 Not Found`: Strategy not found

#### Delete Strategy

- **URL**: `/api/strategies/{id}`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "message": "Strategy deleted successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategy deleted
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to delete this strategy
  - `404 Not Found`: Strategy not found

#### Activate Strategy

- **URL**: `/api/strategies/{id}/activate`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "strategy123",
    "userId": "user123",
    "active": true,
    "updatedAt": "2023-01-01T02:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategy activated
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to activate this strategy
  - `404 Not Found`: Strategy not found

#### Deactivate Strategy

- **URL**: `/api/strategies/{id}/deactivate`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "strategy123",
    "userId": "user123",
    "active": false,
    "updatedAt": "2023-01-01T03:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategy deactivated
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to deactivate this strategy
  - `404 Not Found`: Strategy not found

#### Get Strategies

- **URL**: `/api/strategies?type=momentum&active=true&page=1&limit=20`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Query Parameters**:
  - `type`: Filter by type
  - `symbol`: Filter by symbol
  - `productType`: Filter by product type
  - `active`: Filter by active status
  - `tag`: Filter by tag
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 20)
- **Response**: 
  ```json
  {
    "data": [
      {
        "id": "strategy1",
        "userId": "user123",
        "name": "Momentum Strategy",
        "type": "momentum",
        "symbol": "AAPL",
        "active": true,
        "tags": ["momentum", "equity"]
      },
      {
        "id": "strategy2",
        "userId": "user123",
        "name": "Mean Reversion Strategy",
        "type": "mean-reversion",
        "symbol": "MSFT",
        "active": true,
        "tags": ["mean-reversion", "equity"]
      }
    ],
    "page": 1,
    "limit": 20,
    "total": 2,
    "totalPages": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Strategies retrieved
  - `401 Unauthorized`: Not authenticated

### Portfolio Endpoints

#### Create Portfolio

- **URL**: `/api/portfolios`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "name": "Test Portfolio",
    "description": "A test portfolio for automated trading",
    "strategyId": "strategy123",
    "status": "draft",
    "capital": 10000.0,
    "currency": "USD",
    "riskLevel": "medium"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "name": "Test Portfolio",
    "description": "A test portfolio for automated trading",
    "strategyId": "strategy123",
    "status": "draft",
    "capital": 10000.0,
    "currency": "USD",
    "riskLevel": "medium",
    "legs": [],
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
  ```
- **Status Codes**:
  - `201 Created`: Portfolio created
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated
  - `404 Not Found`: Strategy not found

#### Get Portfolio

- **URL**: `/api/portfolios/{id}`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "name": "Test Portfolio",
    "description": "A test portfolio for automated trading",
    "strategyId": "strategy123",
    "status": "active",
    "capital": 10000.0,
    "currency": "USD",
    "riskLevel": "medium",
    "legs": [
      {
        "id": 1,
        "portfolioId": "portfolio123",
        "symbol": "AAPL",
        "direction": "long",
        "allocation": 0.25,
        "entryPrice": 150.0,
        "stopLoss": 140.0,
        "takeProfit": 170.0,
        "description": "Apple Inc. long position",
        "createdAt": "2023-01-01T01:00:00Z",
        "updatedAt": "2023-01-01T01:00:00Z"
      }
    ],
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T01:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolio retrieved
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to view this portfolio
  - `404 Not Found`: Portfolio not found

#### Update Portfolio

- **URL**: `/api/portfolios/{id}`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "name": "Updated Portfolio",
    "description": "An updated test portfolio",
    "capital": 15000.0,
    "riskLevel": "high"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "name": "Updated Portfolio",
    "description": "An updated test portfolio",
    "strategyId": "strategy123",
    "status": "draft",
    "capital": 15000.0,
    "currency": "USD",
    "riskLevel": "high",
    "updatedAt": "2023-01-01T02:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolio updated
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this portfolio
  - `404 Not Found`: Portfolio not found

#### Delete Portfolio

- **URL**: `/api/portfolios/{id}`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "message": "Portfolio deleted successfully"
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolio deleted
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to delete this portfolio
  - `404 Not Found`: Portfolio not found

#### Activate Portfolio

- **URL**: `/api/portfolios/{id}/activate`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "status": "active",
    "updatedAt": "2023-01-01T03:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolio activated
  - `400 Bad Request`: Portfolio cannot be activated
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to activate this portfolio
  - `404 Not Found`: Portfolio not found

#### Deactivate Portfolio

- **URL**: `/api/portfolios/{id}/deactivate`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "status": "inactive",
    "updatedAt": "2023-01-01T04:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolio deactivated
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to deactivate this portfolio
  - `404 Not Found`: Portfolio not found

#### Add Leg to Portfolio

- **URL**: `/api/portfolios/{id}/legs`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "symbol": "AAPL",
    "direction": "long",
    "allocation": 0.25,
    "entryPrice": 150.0,
    "stopLoss": 140.0,
    "takeProfit": 170.0,
    "description": "Apple Inc. long position"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "legs": [
      {
        "id": 1,
        "portfolioId": "portfolio123",
        "symbol": "AAPL",
        "direction": "long",
        "allocation": 0.25,
        "entryPrice": 150.0,
        "stopLoss": 140.0,
        "takeProfit": 170.0,
        "description": "Apple Inc. long position",
        "createdAt": "2023-01-01T05:00:00Z",
        "updatedAt": "2023-01-01T05:00:00Z"
      }
    ],
    "updatedAt": "2023-01-01T05:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Leg added
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this portfolio
  - `404 Not Found`: Portfolio not found

#### Update Leg in Portfolio

- **URL**: `/api/portfolios/{id}/legs/{legId}`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer {token}`
- **Request Body**:
  ```json
  {
    "allocation": 0.30,
    "entryPrice": 155.0,
    "stopLoss": 145.0,
    "takeProfit": 175.0,
    "description": "Updated Apple Inc. long position"
  }
  ```
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "legs": [
      {
        "id": 1,
        "portfolioId": "portfolio123",
        "symbol": "AAPL",
        "direction": "long",
        "allocation": 0.30,
        "entryPrice": 155.0,
        "stopLoss": 145.0,
        "takeProfit": 175.0,
        "description": "Updated Apple Inc. long position",
        "updatedAt": "2023-01-01T06:00:00Z"
      }
    ],
    "updatedAt": "2023-01-01T06:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Leg updated
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this portfolio
  - `404 Not Found`: Portfolio or leg not found

#### Remove Leg from Portfolio

- **URL**: `/api/portfolios/{id}/legs/{legId}`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer {token}`
- **Response**: 
  ```json
  {
    "id": "portfolio123",
    "userId": "user123",
    "legs": [],
    "updatedAt": "2023-01-01T07:00:00Z"
  }
  ```
- **Status Codes**:
  - `200 OK`: Leg removed
  - `401 Unauthorized`: Not authenticated
  - `403 Forbidden`: Not authorized to modify this portfolio
  - `404 Not Found`: Portfolio or leg not found

#### Get Portfolios

- **URL**: `/api/portfolios?status=active&riskLevel=high&page=1&limit=20`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer {token}`
- **Query Parameters**:
  - `status`: Filter by status
  - `strategyId`: Filter by strategy ID
  - `riskLevel`: Filter by risk level
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 20)
- **Response**: 
  ```json
  {
    "data": [
      {
        "id": "portfolio1",
        "userId": "user123",
        "name": "Growth Portfolio",
        "strategyId": "strategy1",
        "status": "active",
        "capital": 10000.0,
        "riskLevel": "high"
      },
      {
        "id": "portfolio2",
        "userId": "user123",
        "name": "Income Portfolio",
        "strategyId": "strategy2",
        "status": "active",
        "capital": 20000.0,
        "riskLevel": "medium"
      }
    ],
    "page": 1,
    "limit": 20,
    "total": 2,
    "totalPages": 1
  }
  ```
- **Status Codes**:
  - `200 OK`: Portfolios retrieved
  - `401 Unauthorized`: Not authenticated

## Authentication Flow

The authentication system uses JWT (JSON Web Tokens) for secure authentication. Here's the flow:

1. **Registration**:
   - User submits registration form with email, password, and personal details
   - Server validates input and checks if email already exists
   - If valid, server creates a new user with hashed password
   - Server returns user details (without password)

2. **Login**:
   - User submits email and password
   - Server validates credentials
   - If valid, server generates a JWT token
   - Server returns token and user details

3. **Authentication**:
   - Client includes JWT token in Authorization header for protected routes
   - Server validates token using JWT middleware
   - If valid, server extracts user ID and adds it to request context
   - Handlers can access authenticated user ID from context

4. **Token Refresh**:
   - Client can refresh token before expiration
   - Server validates current token and issues a new one
   - This extends the session without requiring re-login

5. **Password Reset**:
   - User requests password reset with email
   - Server generates a reset token and sends email
   - User submits new password with reset token
   - Server validates token and updates password

## Error Handling

The backend implements consistent error handling with standardized error responses:

```json
{
  "error": {
    "code": "invalid_input",
    "message": "Invalid input parameters",
    "details": {
      "email": "Email is required",
      "password": "Password must be at least 8 characters"
    }
  }
}
```

Common error codes:

- `invalid_input`: Input validation failed
- `not_found`: Resource not found
- `unauthorized`: Authentication required
- `forbidden`: Not authorized to access resource
- `internal_error`: Server error

## Database Schema

The backend uses MongoDB with the following collections:

1. **users**: Stores user information
2. **orders**: Stores order information
3. **positions**: Stores position information
4. **strategies**: Stores strategy information
5. **portfolios**: Stores portfolio information

## Testing

The backend includes comprehensive unit tests for all handlers using mock repositories. Tests cover:

- Successful operations
- Error cases
- Authentication and authorization
- Input validation

To run tests:

```bash
cd backend
go test ./tests/...
```

## Deployment

The backend can be deployed as a standalone service or as part of a containerized application. See the DEPLOYMENT.md file for detailed deployment instructions.
