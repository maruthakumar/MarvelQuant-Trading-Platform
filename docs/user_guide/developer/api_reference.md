# API Reference

## Introduction

This API Reference provides comprehensive documentation for developers integrating with the Trading Platform. The platform offers a robust set of RESTful APIs and WebSocket endpoints that enable programmatic access to all trading functionality, market data, account information, and more.

## API Overview

The Trading Platform API is organized into logical service groups, each handling specific functionality domains. All API endpoints follow RESTful principles and use standard HTTP methods.

### Base URLs

- **Production Environment**: `https://api.tradingplatform.example.com/v1`
- **Staging Environment**: `https://api-staging.tradingplatform.example.com/v1`
- **Sandbox Environment**: `https://api-sandbox.tradingplatform.example.com/v1`

### API Versioning

The API uses versioning in the URL path to ensure backward compatibility. The current version is `v1`. When breaking changes are introduced, a new version will be released, and the previous version will be maintained for a deprecation period of at least 6 months.

### Authentication

All API requests require authentication using one of the following methods:

#### OAuth 2.0 (Recommended)

For user-context operations, the OAuth 2.0 authorization framework is used:

1. Redirect users to: `https://auth.tradingplatform.example.com/oauth/authorize`
2. Users authenticate and authorize your application
3. Users are redirected back to your application with an authorization code
4. Exchange the code for an access token using: `https://auth.tradingplatform.example.com/oauth/token`
5. Include the access token in the `Authorization` header of API requests:

```
Authorization: Bearer {access_token}
```

#### API Keys

For server-to-server integration, API keys can be used:

1. Generate API keys in the platform's developer portal
2. Include the API key and secret in the request headers:

```
X-API-Key: {api_key}
X-API-Secret: {api_secret}
```

### Rate Limiting

To ensure fair usage and system stability, rate limits are applied to API requests:

- **OAuth 2.0**: 120 requests per minute per user
- **API Keys**: 300 requests per minute per API key
- **WebSocket**: 10 connections per user, with subscription limits per connection

Rate limit headers are included in API responses:

```
X-RateLimit-Limit: 120
X-RateLimit-Remaining: 115
X-RateLimit-Reset: 1617981234
```

When rate limits are exceeded, the API returns a `429 Too Many Requests` status code.

## Common Patterns

### Request Format

- **Content-Type**: `application/json` for request bodies
- **Accept**: `application/json` for response format
- **Charset**: UTF-8

### Response Format

All API responses follow a consistent structure:

```json
{
  "status": "success",
  "data": {
    // Response data specific to the endpoint
  },
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 25,
      "total": 100,
      "total_pages": 4
    }
  }
}
```

For error responses:

```json
{
  "status": "error",
  "error": {
    "code": "invalid_parameter",
    "message": "The provided parameter is invalid",
    "details": {
      "parameter": "symbol",
      "reason": "Symbol AAPL123 does not exist"
    }
  }
}
```

### Pagination

For endpoints that return collections of resources, pagination is supported:

- **page**: Page number (starting from 1)
- **per_page**: Number of items per page (default: 25, max: 100)

Example:
```
GET /orders?page=2&per_page=50
```

Response includes pagination metadata:

```json
"meta": {
  "pagination": {
    "page": 2,
    "per_page": 50,
    "total": 320,
    "total_pages": 7
  }
}
```

### Filtering

Many endpoints support filtering using query parameters:

```
GET /orders?status=open&symbol=AAPL&side=buy
```

### Sorting

Sorting is supported using the `sort` parameter:

```
GET /orders?sort=created_at:desc
```

Multiple sort fields can be specified:

```
GET /orders?sort=status:asc,created_at:desc
```

### Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests:

- **2xx**: Success
- **4xx**: Client error (invalid request)
- **5xx**: Server error

Common error codes:

- **400 Bad Request**: Invalid request parameters
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **422 Unprocessable Entity**: Validation error
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server-side error

## API Services

### Authentication Service

#### Endpoints

##### `POST /auth/token`

Obtain an access token using OAuth 2.0 client credentials flow.

**Request:**
```json
{
  "grant_type": "client_credentials",
  "client_id": "your_client_id",
  "client_secret": "your_client_secret"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "scope": "read write"
  }
}
```

##### `POST /auth/refresh`

Refresh an expired access token.

**Request:**
```json
{
  "grant_type": "refresh_token",
  "refresh_token": "your_refresh_token"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "refresh_token": "new_refresh_token",
    "scope": "read write"
  }
}
```

##### `POST /auth/revoke`

Revoke an access token.

**Request:**
```json
{
  "token": "your_access_token",
  "token_type_hint": "access_token"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "message": "Token revoked successfully"
  }
}
```

### User Service

#### Endpoints

##### `GET /users/me`

Get the current user's profile information.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "usr_123456789",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "created_at": "2023-01-15T08:30:00Z",
    "updated_at": "2023-03-20T14:15:30Z",
    "preferences": {
      "theme": "dark",
      "default_order_type": "limit",
      "notification_settings": {
        "email": true,
        "push": true,
        "sms": false
      }
    }
  }
}
```

##### `PATCH /users/me`

Update the current user's profile information.

**Request:**
```json
{
  "first_name": "Johnny",
  "phone": "+1987654321",
  "preferences": {
    "theme": "light"
  }
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "usr_123456789",
    "email": "user@example.com",
    "first_name": "Johnny",
    "last_name": "Doe",
    "phone": "+1987654321",
    "created_at": "2023-01-15T08:30:00Z",
    "updated_at": "2023-04-05T09:45:12Z",
    "preferences": {
      "theme": "light",
      "default_order_type": "limit",
      "notification_settings": {
        "email": true,
        "push": true,
        "sms": false
      }
    }
  }
}
```

##### `GET /users/me/api-keys`

List all API keys for the current user.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "key_abcdef123456",
      "name": "Trading Bot",
      "created_at": "2023-02-10T12:00:00Z",
      "last_used_at": "2023-04-04T18:30:45Z",
      "permissions": ["read", "trade"],
      "expires_at": "2024-02-10T12:00:00Z"
    },
    {
      "id": "key_ghijkl789012",
      "name": "Data Analysis",
      "created_at": "2023-03-15T09:30:00Z",
      "last_used_at": "2023-04-05T08:15:20Z",
      "permissions": ["read"],
      "expires_at": "2024-03-15T09:30:00Z"
    }
  ]
}
```

##### `POST /users/me/api-keys`

Create a new API key.

**Request:**
```json
{
  "name": "Automated Trading",
  "permissions": ["read", "trade"],
  "expires_in": 31536000
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "key_mnopqr345678",
    "name": "Automated Trading",
    "api_key": "pk_live_abcdefghijklmnopqrstuvwxyz",
    "api_secret": "sk_live_abcdefghijklmnopqrstuvwxyz123456789",
    "created_at": "2023-04-05T10:00:00Z",
    "permissions": ["read", "trade"],
    "expires_at": "2024-04-04T10:00:00Z"
  },
  "meta": {
    "warning": "This is the only time the API secret will be displayed. Please store it securely."
  }
}
```

##### `DELETE /users/me/api-keys/{key_id}`

Delete an API key.

**Response:**
```json
{
  "status": "success",
  "data": {
    "message": "API key deleted successfully"
  }
}
```

### Order Service

#### Endpoints

##### `GET /orders`

List orders with optional filtering.

**Parameters:**
- `status` (string, optional): Filter by order status (open, filled, canceled, rejected)
- `symbol` (string, optional): Filter by instrument symbol
- `side` (string, optional): Filter by order side (buy, sell)
- `type` (string, optional): Filter by order type (market, limit, stop, stop_limit)
- `from_date` (string, optional): Filter by creation date (ISO 8601 format)
- `to_date` (string, optional): Filter by creation date (ISO 8601 format)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "ord_123456789",
      "user_id": "usr_123456789",
      "symbol": "AAPL",
      "side": "buy",
      "type": "limit",
      "status": "open",
      "quantity": 100,
      "filled_quantity": 0,
      "price": 150.00,
      "stop_price": null,
      "time_in_force": "gtc",
      "created_at": "2023-04-05T09:30:00Z",
      "updated_at": "2023-04-05T09:30:00Z",
      "expires_at": null,
      "client_order_id": "my_order_123"
    },
    {
      "id": "ord_987654321",
      "user_id": "usr_123456789",
      "symbol": "MSFT",
      "side": "sell",
      "type": "market",
      "status": "filled",
      "quantity": 50,
      "filled_quantity": 50,
      "price": null,
      "stop_price": null,
      "time_in_force": "day",
      "created_at": "2023-04-04T14:15:00Z",
      "updated_at": "2023-04-04T14:15:05Z",
      "expires_at": "2023-04-04T20:00:00Z",
      "client_order_id": "my_order_456",
      "executions": [
        {
          "id": "exec_123456",
          "price": 287.50,
          "quantity": 50,
          "timestamp": "2023-04-04T14:15:05Z"
        }
      ]
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 25,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

##### `GET /orders/{order_id}`

Get a specific order by ID.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "ord_123456789",
    "user_id": "usr_123456789",
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "status": "open",
    "quantity": 100,
    "filled_quantity": 0,
    "price": 150.00,
    "stop_price": null,
    "time_in_force": "gtc",
    "created_at": "2023-04-05T09:30:00Z",
    "updated_at": "2023-04-05T09:30:00Z",
    "expires_at": null,
    "client_order_id": "my_order_123",
    "executions": []
  }
}
```

##### `POST /orders`

Create a new order.

**Request:**
```json
{
  "symbol": "AAPL",
  "side": "buy",
  "type": "limit",
  "quantity": 100,
  "price": 150.00,
  "time_in_force": "gtc",
  "client_order_id": "my_order_789"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "ord_abcdef123456",
    "user_id": "usr_123456789",
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "status": "open",
    "quantity": 100,
    "filled_quantity": 0,
    "price": 150.00,
    "stop_price": null,
    "time_in_force": "gtc",
    "created_at": "2023-04-05T10:15:00Z",
    "updated_at": "2023-04-05T10:15:00Z",
    "expires_at": null,
    "client_order_id": "my_order_789",
    "executions": []
  }
}
```

##### `PATCH /orders/{order_id}`

Modify an existing order.

**Request:**
```json
{
  "quantity": 50,
  "price": 152.00
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "ord_123456789",
    "user_id": "usr_123456789",
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "status": "open",
    "quantity": 50,
    "filled_quantity": 0,
    "price": 152.00,
    "stop_price": null,
    "time_in_force": "gtc",
    "created_at": "2023-04-05T09:30:00Z",
    "updated_at": "2023-04-05T10:20:00Z",
    "expires_at": null,
    "client_order_id": "my_order_123",
    "executions": []
  }
}
```

##### `DELETE /orders/{order_id}`

Cancel an order.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "ord_123456789",
    "user_id": "usr_123456789",
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "status": "canceled",
    "quantity": 50,
    "filled_quantity": 0,
    "price": 152.00,
    "stop_price": null,
    "time_in_force": "gtc",
    "created_at": "2023-04-05T09:30:00Z",
    "updated_at": "2023-04-05T10:25:00Z",
    "expires_at": null,
    "client_order_id": "my_order_123",
    "executions": []
  }
}
```

### Portfolio Service

#### Endpoints

##### `GET /portfolio/positions`

Get current positions.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "symbol": "AAPL",
      "quantity": 150,
      "average_price": 145.75,
      "current_price": 160.25,
      "market_value": 24037.50,
      "unrealized_pl": 2175.00,
      "unrealized_pl_percent": 9.95,
      "day_pl": 375.00,
      "day_pl_percent": 1.58,
      "cost_basis": 21862.50,
      "updated_at": "2023-04-05T10:30:00Z"
    },
    {
      "symbol": "MSFT",
      "quantity": -50,
      "average_price": 287.50,
      "current_price": 280.75,
      "market_value": -14037.50,
      "unrealized_pl": 337.50,
      "unrealized_pl_percent": 2.35,
      "day_pl": 125.00,
      "day_pl_percent": 0.88,
      "cost_basis": -14375.00,
      "updated_at": "2023-04-05T10:30:00Z"
    }
  ]
}
```

##### `GET /portfolio/positions/{symbol}`

Get position details for a specific symbol.

**Response:**
```json
{
  "status": "success",
  "data": {
    "symbol": "AAPL",
    "quantity": 150,
    "average_price": 145.75,
    "current_price": 160.25,
    "market_value": 24037.50,
    "unrealized_pl": 2175.00,
    "unrealized_pl_percent": 9.95,
    "day_pl": 375.00,
    "day_pl_percent": 1.58,
    "cost_basis": 21862.50,
    "updated_at": "2023-04-05T10:30:00Z",
    "trades": [
      {
        "order_id": "ord_123456789",
        "side": "buy",
        "quantity": 100,
        "price": 142.50,
        "timestamp": "2023-03-15T14:30:00Z"
      },
      {
        "order_id": "ord_abcdefghij",
        "side": "buy",
        "quantity": 50,
        "price": 152.25,
        "timestamp": "2023-03-28T10:15:00Z"
      }
    ],
    "risk_metrics": {
      "beta": 1.15,
      "volatility": 0.28,
      "var_95": 1250.00,
      "sharpe_ratio": 1.8
    }
  }
}
```

##### `GET /portfolio/performance`

Get portfolio performance metrics.

**Parameters:**
- `period` (string, optional): Time period for performance calculation (day, week, month, year, all)
- `from_date` (string, optional): Start date for custom period (ISO 8601 format)
- `to_date` (string, optional): End date for custom period (ISO 8601 format)

**Response:**
```json
{
  "status": "success",
  "data": {
    "total_value": 125000.00,
    "cash_balance": 35000.00,
    "invested_value": 90000.00,
    "day_pl": 1250.00,
    "day_pl_percent": 1.01,
    "total_pl": 15000.00,
    "total_pl_percent": 13.64,
    "annualized_return": 12.5,
    "risk_metrics": {
      "volatility": 0.18,
      "sharpe_ratio": 1.4,
      "sortino_ratio": 2.1,
      "max_drawdown": 8.5,
      "beta": 0.95,
      "alpha": 2.3
    },
    "allocation": {
      "by_asset_class": {
        "equity": 65.0,
        "fixed_income": 20.0,
        "cash": 15.0
      },
      "by_sector": {
        "technology": 40.0,
        "healthcare": 15.0,
        "financials": 10.0,
        "consumer_discretionary": 10.0,
        "other": 25.0
      },
      "by_geography": {
        "north_america": 70.0,
        "europe": 15.0,
        "asia_pacific": 10.0,
        "other": 5.0
      }
    },
    "performance_series": {
      "timestamps": [
        "2023-03-06T00:00:00Z",
        "2023-03-13T00:00:00Z",
        "2023-03-20T00:00:00Z",
        "2023-03-27T00:00:00Z",
        "2023-04-03T00:00:00Z"
      ],
      "values": [
        115000.00,
        118500.00,
        117000.00,
        121000.00,
        125000.00
      ],
      "returns": [
        null,
        3.04,
        -1.27,
        3.42,
        3.31
      ]
    }
  }
}
```

### Market Data Service

#### Endpoints

##### `GET /market-data/quotes/{symbol}`

Get real-time quote for a symbol.

**Response:**
```json
{
  "status": "success",
  "data": {
    "symbol": "AAPL",
    "bid": 160.10,
    "ask": 160.15,
    "last": 160.12,
    "volume": 12500000,
    "open": 158.50,
    "high": 161.20,
    "low": 158.25,
    "prev_close": 159.30,
    "change": 0.82,
    "change_percent": 0.51,
    "timestamp": "2023-04-05T10:35:45Z"
  }
}
```

##### `GET /market-data/quotes`

Get real-time quotes for multiple symbols.

**Parameters:**
- `symbols` (string, required): Comma-separated list of symbols

**Response:**
```json
{
  "status": "success",
  "data": {
    "AAPL": {
      "bid": 160.10,
      "ask": 160.15,
      "last": 160.12,
      "volume": 12500000,
      "open": 158.50,
      "high": 161.20,
      "low": 158.25,
      "prev_close": 159.30,
      "change": 0.82,
      "change_percent": 0.51,
      "timestamp": "2023-04-05T10:35:45Z"
    },
    "MSFT": {
      "bid": 280.60,
      "ask": 280.75,
      "last": 280.70,
      "volume": 8750000,
      "open": 279.00,
      "high": 282.50,
      "low": 278.75,
      "prev_close": 281.50,
      "change": -0.80,
      "change_percent": -0.28,
      "timestamp": "2023-04-05T10:35:45Z"
    }
  }
}
```

##### `GET /market-data/historical/{symbol}`

Get historical price data for a symbol.

**Parameters:**
- `interval` (string, required): Data interval (1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1mo)
- `from` (string, required): Start date/time (ISO 8601 format)
- `to` (string, required): End date/time (ISO 8601 format)
- `adjusted` (boolean, optional): Whether to return adjusted prices (default: true)

**Response:**
```json
{
  "status": "success",
  "data": {
    "symbol": "AAPL",
    "interval": "1d",
    "adjusted": true,
    "candles": [
      {
        "timestamp": "2023-04-03T00:00:00Z",
        "open": 157.80,
        "high": 159.45,
        "low": 157.25,
        "close": 159.30,
        "volume": 52500000
      },
      {
        "timestamp": "2023-04-04T00:00:00Z",
        "open": 159.40,
        "high": 160.75,
        "low": 158.80,
        "close": 160.10,
        "volume": 48750000
      },
      {
        "timestamp": "2023-04-05T00:00:00Z",
        "open": 158.50,
        "high": 161.20,
        "low": 158.25,
        "close": 160.12,
        "volume": 12500000
      }
    ]
  }
}
```

##### `GET /market-data/indicators/{symbol}`

Calculate technical indicators for a symbol.

**Parameters:**
- `indicator` (string, required): Indicator type (sma, ema, rsi, macd, bollinger, etc.)
- `params` (object, required): Indicator-specific parameters
- `interval` (string, required): Data interval (1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1mo)
- `from` (string, required): Start date/time (ISO 8601 format)
- `to` (string, required): End date/time (ISO 8601 format)

**Example Request:**
```
GET /market-data/indicators/AAPL?indicator=sma&params[period]=20&interval=1d&from=2023-03-01T00:00:00Z&to=2023-04-05T00:00:00Z
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "symbol": "AAPL",
    "indicator": "sma",
    "params": {
      "period": 20
    },
    "interval": "1d",
    "values": [
      {
        "timestamp": "2023-03-20T00:00:00Z",
        "value": 152.35
      },
      {
        "timestamp": "2023-03-21T00:00:00Z",
        "value": 153.10
      },
      // ... more data points
      {
        "timestamp": "2023-04-05T00:00:00Z",
        "value": 158.75
      }
    ]
  }
}
```

### Strategy Service

#### Endpoints

##### `GET /strategies`

List all strategies.

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "strat_123456789",
      "name": "Moving Average Crossover",
      "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
      "status": "active",
      "created_at": "2023-03-10T09:00:00Z",
      "updated_at": "2023-04-01T14:30:00Z",
      "symbols": ["AAPL", "MSFT", "GOOGL"],
      "performance": {
        "total_return": 8.5,
        "annualized_return": 12.3,
        "sharpe_ratio": 1.5,
        "max_drawdown": 5.2
      }
    },
    {
      "id": "strat_987654321",
      "name": "RSI Mean Reversion",
      "description": "Buy when RSI is below 30, sell when RSI is above 70",
      "status": "paused",
      "created_at": "2023-02-15T11:30:00Z",
      "updated_at": "2023-03-20T16:45:00Z",
      "symbols": ["AMZN", "NFLX", "META"],
      "performance": {
        "total_return": 5.2,
        "annualized_return": 9.8,
        "sharpe_ratio": 1.2,
        "max_drawdown": 7.5
      }
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 25,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

##### `GET /strategies/{strategy_id}`

Get a specific strategy.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "strat_123456789",
    "name": "Moving Average Crossover",
    "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
    "status": "active",
    "created_at": "2023-03-10T09:00:00Z",
    "updated_at": "2023-04-01T14:30:00Z",
    "symbols": ["AAPL", "MSFT", "GOOGL"],
    "parameters": {
      "fast_period": 50,
      "slow_period": 200,
      "position_size_percent": 5,
      "max_positions": 10
    },
    "conditions": {
      "entry": {
        "type": "crossover",
        "indicators": [
          {
            "type": "sma",
            "params": { "period": 50 }
          },
          {
            "type": "sma",
            "params": { "period": 200 }
          }
        ]
      },
      "exit": {
        "type": "crossunder",
        "indicators": [
          {
            "type": "sma",
            "params": { "period": 50 }
          },
          {
            "type": "sma",
            "params": { "period": 200 }
          }
        ]
      }
    },
    "performance": {
      "total_return": 8.5,
      "annualized_return": 12.3,
      "sharpe_ratio": 1.5,
      "max_drawdown": 5.2,
      "win_rate": 65.0,
      "avg_win": 2.8,
      "avg_loss": 1.2,
      "profit_factor": 2.3
    },
    "trades": [
      {
        "symbol": "AAPL",
        "entry_date": "2023-03-15T10:30:00Z",
        "entry_price": 150.25,
        "exit_date": "2023-03-28T14:45:00Z",
        "exit_price": 158.75,
        "quantity": 100,
        "profit_loss": 850.00,
        "profit_loss_percent": 5.66
      },
      {
        "symbol": "MSFT",
        "entry_date": "2023-03-18T11:15:00Z",
        "entry_price": 275.50,
        "exit_date": null,
        "exit_price": null,
        "quantity": 50,
        "profit_loss": null,
        "profit_loss_percent": null
      }
    ]
  }
}
```

##### `POST /strategies`

Create a new strategy.

**Request:**
```json
{
  "name": "Bollinger Band Breakout",
  "description": "Buy on upper band breakout, sell on lower band breakout",
  "symbols": ["SPY", "QQQ", "IWM"],
  "parameters": {
    "period": 20,
    "std_dev": 2,
    "position_size_percent": 3,
    "max_positions": 5
  },
  "conditions": {
    "entry": {
      "type": "breakout",
      "indicators": [
        {
          "type": "bollinger_bands",
          "params": { "period": 20, "std_dev": 2 }
        }
      ]
    },
    "exit": {
      "type": "stop_loss",
      "params": { "percent": 5 }
    }
  }
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "strat_abcdef123456",
    "name": "Bollinger Band Breakout",
    "description": "Buy on upper band breakout, sell on lower band breakout",
    "status": "inactive",
    "created_at": "2023-04-05T11:00:00Z",
    "updated_at": "2023-04-05T11:00:00Z",
    "symbols": ["SPY", "QQQ", "IWM"],
    "parameters": {
      "period": 20,
      "std_dev": 2,
      "position_size_percent": 3,
      "max_positions": 5
    },
    "conditions": {
      "entry": {
        "type": "breakout",
        "indicators": [
          {
            "type": "bollinger_bands",
            "params": { "period": 20, "std_dev": 2 }
          }
        ]
      },
      "exit": {
        "type": "stop_loss",
        "params": { "percent": 5 }
      }
    }
  }
}
```

##### `PATCH /strategies/{strategy_id}`

Update a strategy.

**Request:**
```json
{
  "status": "active",
  "parameters": {
    "std_dev": 2.5,
    "position_size_percent": 2
  }
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "strat_abcdef123456",
    "name": "Bollinger Band Breakout",
    "description": "Buy on upper band breakout, sell on lower band breakout",
    "status": "active",
    "created_at": "2023-04-05T11:00:00Z",
    "updated_at": "2023-04-05T11:15:00Z",
    "symbols": ["SPY", "QQQ", "IWM"],
    "parameters": {
      "period": 20,
      "std_dev": 2.5,
      "position_size_percent": 2,
      "max_positions": 5
    },
    "conditions": {
      "entry": {
        "type": "breakout",
        "indicators": [
          {
            "type": "bollinger_bands",
            "params": { "period": 20, "std_dev": 2 }
          }
        ]
      },
      "exit": {
        "type": "stop_loss",
        "params": { "percent": 5 }
      }
    }
  }
}
```

##### `DELETE /strategies/{strategy_id}`

Delete a strategy.

**Response:**
```json
{
  "status": "success",
  "data": {
    "message": "Strategy deleted successfully"
  }
}
```

##### `POST /strategies/{strategy_id}/backtest`

Run a backtest for a strategy.

**Request:**
```json
{
  "start_date": "2022-01-01T00:00:00Z",
  "end_date": "2023-01-01T00:00:00Z",
  "initial_capital": 100000,
  "parameters": {
    "period": 20,
    "std_dev": 2.5
  }
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "backtest_id": "bt_123456789",
    "strategy_id": "strat_abcdef123456",
    "status": "processing",
    "estimated_completion_time": "2023-04-05T11:30:00Z"
  }
}
```

##### `GET /strategies/{strategy_id}/backtest/{backtest_id}`

Get backtest results.

**Response:**
```json
{
  "status": "success",
  "data": {
    "backtest_id": "bt_123456789",
    "strategy_id": "strat_abcdef123456",
    "status": "completed",
    "start_date": "2022-01-01T00:00:00Z",
    "end_date": "2023-01-01T00:00:00Z",
    "initial_capital": 100000,
    "final_capital": 112500,
    "total_return": 12.5,
    "annualized_return": 12.5,
    "sharpe_ratio": 1.3,
    "sortino_ratio": 1.8,
    "max_drawdown": 8.2,
    "win_rate": 58.0,
    "profit_factor": 1.9,
    "trades_count": 45,
    "parameters": {
      "period": 20,
      "std_dev": 2.5,
      "position_size_percent": 2,
      "max_positions": 5
    },
    "equity_curve": {
      "timestamps": [
        "2022-01-01T00:00:00Z",
        "2022-02-01T00:00:00Z",
        // ... more timestamps
        "2023-01-01T00:00:00Z"
      ],
      "equity": [
        100000,
        102500,
        // ... more equity values
        112500
      ],
      "drawdown": [
        0,
        0,
        // ... more drawdown values
        0
      ]
    },
    "trades": [
      {
        "symbol": "SPY",
        "entry_date": "2022-01-15T10:30:00Z",
        "entry_price": 450.25,
        "exit_date": "2022-02-10T14:45:00Z",
        "exit_price": 458.75,
        "quantity": 50,
        "profit_loss": 425.00,
        "profit_loss_percent": 1.89
      },
      // ... more trades
    ]
  }
}
```

## WebSocket API

The WebSocket API provides real-time data streams for market data, order updates, and other events.

### Connection

Connect to the WebSocket API at:

```
wss://ws.tradingplatform.example.com/v1
```

Authentication is required using one of the following methods:

1. **JWT Token**: Include the token in the connection URL:
   ```
   wss://ws.tradingplatform.example.com/v1?token=your_jwt_token
   ```

2. **API Key**: Include the API key in the connection URL:
   ```
   wss://ws.tradingplatform.example.com/v1?api_key=your_api_key&api_secret=your_api_secret
   ```

### Message Format

All WebSocket messages use JSON format:

```json
{
  "type": "message_type",
  "data": {
    // Message-specific data
  }
}
```

### Subscription

To subscribe to data streams, send a subscription message:

```json
{
  "type": "subscribe",
  "channels": [
    {
      "name": "quotes",
      "symbols": ["AAPL", "MSFT", "GOOGL"]
    },
    {
      "name": "orders"
    }
  ]
}
```

The server will respond with a confirmation:

```json
{
  "type": "subscription_success",
  "channels": [
    {
      "name": "quotes",
      "symbols": ["AAPL", "MSFT", "GOOGL"]
    },
    {
      "name": "orders"
    }
  ]
}
```

### Available Channels

#### Quotes Channel

Real-time market quotes for specified symbols.

**Subscription:**
```json
{
  "type": "subscribe",
  "channels": [
    {
      "name": "quotes",
      "symbols": ["AAPL", "MSFT"]
    }
  ]
}
```

**Message Format:**
```json
{
  "type": "quote",
  "data": {
    "symbol": "AAPL",
    "bid": 160.10,
    "ask": 160.15,
    "last": 160.12,
    "volume": 12500000,
    "timestamp": "2023-04-05T11:30:45.123Z"
  }
}
```

#### Orders Channel

Real-time updates for your orders.

**Subscription:**
```json
{
  "type": "subscribe",
  "channels": [
    {
      "name": "orders"
    }
  ]
}
```

**Message Format:**
```json
{
  "type": "order_update",
  "data": {
    "id": "ord_123456789",
    "user_id": "usr_123456789",
    "symbol": "AAPL",
    "side": "buy",
    "type": "limit",
    "status": "filled",
    "quantity": 100,
    "filled_quantity": 100,
    "price": 150.00,
    "time_in_force": "gtc",
    "created_at": "2023-04-05T09:30:00Z",
    "updated_at": "2023-04-05T11:45:00Z",
    "client_order_id": "my_order_123",
    "execution": {
      "id": "exec_abcdef",
      "price": 150.00,
      "quantity": 100,
      "timestamp": "2023-04-05T11:45:00Z"
    }
  }
}
```

#### Positions Channel

Real-time updates for your positions.

**Subscription:**
```json
{
  "type": "subscribe",
  "channels": [
    {
      "name": "positions"
    }
  ]
}
```

**Message Format:**
```json
{
  "type": "position_update",
  "data": {
    "symbol": "AAPL",
    "quantity": 250,
    "average_price": 147.50,
    "current_price": 160.25,
    "market_value": 40062.50,
    "unrealized_pl": 3187.50,
    "unrealized_pl_percent": 8.65,
    "updated_at": "2023-04-05T11:45:00Z"
  }
}
```

#### Strategy Channel

Real-time updates for your strategies.

**Subscription:**
```json
{
  "type": "subscribe",
  "channels": [
    {
      "name": "strategies"
    }
  ]
}
```

**Message Format:**
```json
{
  "type": "strategy_signal",
  "data": {
    "strategy_id": "strat_123456789",
    "symbol": "AAPL",
    "signal": "buy",
    "price": 160.25,
    "timestamp": "2023-04-05T11:50:00Z",
    "indicators": {
      "sma_50": 155.75,
      "sma_200": 145.50
    },
    "notes": "50-day MA crossed above 200-day MA"
  }
}
```

### Unsubscription

To unsubscribe from data streams, send an unsubscription message:

```json
{
  "type": "unsubscribe",
  "channels": [
    {
      "name": "quotes",
      "symbols": ["AAPL"]
    }
  ]
}
```

The server will respond with a confirmation:

```json
{
  "type": "unsubscription_success",
  "channels": [
    {
      "name": "quotes",
      "symbols": ["AAPL"]
    }
  ]
}
```

### Heartbeat

The server sends heartbeat messages every 30 seconds to keep the connection alive:

```json
{
  "type": "heartbeat",
  "timestamp": "2023-04-05T11:55:00Z"
}
```

Clients should respond with a heartbeat acknowledgment:

```json
{
  "type": "heartbeat_ack",
  "timestamp": "2023-04-05T11:55:00Z"
}
```

If no heartbeat acknowledgment is received for 90 seconds, the server will close the connection.

## Error Codes

The API uses standardized error codes to provide detailed information about failures:

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `authentication_required` | Authentication is required | 401 |
| `invalid_credentials` | Invalid API key or token | 401 |
| `permission_denied` | Insufficient permissions | 403 |
| `resource_not_found` | Requested resource not found | 404 |
| `validation_error` | Request validation failed | 422 |
| `rate_limit_exceeded` | Rate limit exceeded | 429 |
| `internal_error` | Internal server error | 500 |
| `service_unavailable` | Service temporarily unavailable | 503 |

## SDK Libraries

The Trading Platform provides official SDK libraries for popular programming languages:

- [Python SDK](https://github.com/tradingplatform/python-sdk)
- [JavaScript SDK](https://github.com/tradingplatform/js-sdk)
- [Java SDK](https://github.com/tradingplatform/java-sdk)
- [C# SDK](https://github.com/tradingplatform/csharp-sdk)

## Best Practices

### Authentication
- Store API keys securely and never expose them in client-side code
- Implement token refresh logic to handle expiring access tokens
- Use the minimum required permissions for API keys

### Rate Limiting
- Implement exponential backoff for rate limit errors
- Cache frequently accessed data to reduce API calls
- Batch operations when possible (e.g., getting quotes for multiple symbols)

### WebSocket Usage
- Implement reconnection logic with exponential backoff
- Handle heartbeat messages to keep connections alive
- Subscribe only to the data you need

### Error Handling
- Implement proper error handling for all API calls
- Log detailed error information for troubleshooting
- Provide user-friendly error messages in your application

## Support and Resources

- [API Documentation Portal](https://developers.tradingplatform.example.com)
- [API Status Page](https://status.tradingplatform.example.com)
- [Developer Forum](https://forum.tradingplatform.example.com)
- [GitHub Repositories](https://github.com/tradingplatform)
- [Support Email](mailto:api-support@tradingplatform.example.com)
