# API Documentation

## Overview

This document provides comprehensive documentation for the Trading Platform API. The API allows developers to programmatically access trading functionality, market data, account information, and more. This documentation covers authentication, endpoints, request/response formats, error handling, and includes examples for common use cases.

## Base URL

All API endpoints are relative to the base URL:

```
https://api.tradingplatform.com/v1
```

## Authentication

### JWT Authentication

The API uses JSON Web Tokens (JWT) for authentication. To authenticate, you must first obtain an access token by calling the login endpoint.

#### Obtaining a Token

```
POST /auth/login
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 900,
  "user": {
    "id": "user123",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

#### Using the Token

Include the token in the Authorization header for all authenticated requests:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Token Refresh

Access tokens expire after 15 minutes. Use the refresh token to obtain a new access token:

```
POST /auth/refresh-token
```

**Request Body:**

```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 900
}
```

### API Key Authentication

For server-to-server integrations, you can use API key authentication:

1. Generate an API key in the platform settings
2. Include the API key in the `X-API-Key` header:

```
X-API-Key: your_api_key
```

## Rate Limiting

The API implements rate limiting to protect against abuse. Limits are applied per API key or user account:

- 100 requests per minute for standard endpoints
- 10 requests per minute for resource-intensive endpoints

When a rate limit is exceeded, the API returns a 429 Too Many Requests response with a Retry-After header indicating when you can resume making requests.

## Endpoints

### User Management

#### Get User Profile

```
GET /users/profile
```

**Response:**

```json
{
  "id": "user123",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "phone": "+1234567890",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "preferences": {
    "theme": "dark",
    "notifications": {
      "email": true,
      "sms": false,
      "push": true
    }
  },
  "createdAt": "2023-01-15T08:30:00Z",
  "updatedAt": "2023-03-20T14:15:30Z"
}
```

#### Update User Profile

```
PUT /users/profile
```

**Request Body:**

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "phone": "+1234567890",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "preferences": {
    "theme": "dark",
    "notifications": {
      "email": true,
      "sms": false,
      "push": true
    }
  }
}
```

**Response:**

```json
{
  "id": "user123",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "phone": "+1234567890",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "preferences": {
    "theme": "dark",
    "notifications": {
      "email": true,
      "sms": false,
      "push": true
    }
  },
  "updatedAt": "2023-04-02T10:25:30Z"
}
```

### Market Data

#### Get Quote

```
GET /market/quote?symbol=AAPL
```

**Parameters:**

| Parameter | Type   | Required | Description                                      |
|-----------|--------|----------|--------------------------------------------------|
| symbol    | string | Yes      | Symbol of the instrument (e.g., AAPL, MSFT, BTC) |

**Response:**

```json
{
  "symbol": "AAPL",
  "exchange": "NASDAQ",
  "lastPrice": 175.25,
  "change": 2.35,
  "changePercent": 1.36,
  "bidPrice": 175.20,
  "bidSize": 500,
  "askPrice": 175.30,
  "askSize": 300,
  "volume": 15234567,
  "openPrice": 173.50,
  "highPrice": 176.40,
  "lowPrice": 173.10,
  "previousClose": 172.90,
  "timestamp": "2023-04-02T16:30:00Z"
}
```

#### Get Historical Data

```
GET /market/history?symbol=AAPL&interval=1d&from=2023-03-01&to=2023-04-01
```

**Parameters:**

| Parameter | Type   | Required | Description                                                                      |
|-----------|--------|----------|----------------------------------------------------------------------------------|
| symbol    | string | Yes      | Symbol of the instrument                                                         |
| interval  | string | Yes      | Time interval (1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1M)                            |
| from      | string | Yes      | Start date (YYYY-MM-DD) or timestamp                                             |
| to        | string | Yes      | End date (YYYY-MM-DD) or timestamp                                               |
| fields    | string | No       | Comma-separated list of fields to include (default: o,h,l,c,v)                   |

**Response:**

```json
{
  "symbol": "AAPL",
  "interval": "1d",
  "data": [
    {
      "timestamp": "2023-03-01T00:00:00Z",
      "open": 165.30,
      "high": 167.40,
      "low": 164.80,
      "close": 166.90,
      "volume": 12345678
    },
    {
      "timestamp": "2023-03-02T00:00:00Z",
      "open": 167.10,
      "high": 168.20,
      "low": 166.50,
      "close": 167.80,
      "volume": 10987654
    },
    // Additional data points...
  ]
}
```

#### Search Instruments

```
GET /market/search?query=apple&type=stock&limit=10
```

**Parameters:**

| Parameter | Type   | Required | Description                                                        |
|-----------|--------|----------|--------------------------------------------------------------------|
| query     | string | Yes      | Search query                                                       |
| type      | string | No       | Instrument type (stock, option, future, forex, crypto)             |
| exchange  | string | No       | Filter by exchange                                                 |
| limit     | number | No       | Maximum number of results to return (default: 10, max: 100)        |

**Response:**

```json
{
  "results": [
    {
      "symbol": "AAPL",
      "name": "Apple Inc.",
      "type": "stock",
      "exchange": "NASDAQ",
      "currency": "USD"
    },
    {
      "symbol": "AAPL230421C00170000",
      "name": "AAPL Apr 21 2023 170 Call",
      "type": "option",
      "exchange": "OPRA",
      "currency": "USD",
      "underlying": "AAPL",
      "strikePrice": 170.00,
      "expirationDate": "2023-04-21",
      "optionType": "call"
    },
    // Additional results...
  ],
  "count": 2,
  "total": 2
}
```

### Order Management

#### Place Order

```
POST /orders
```

**Request Body:**

```json
{
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "quantity": 10,
  "price": 175.50,
  "timeInForce": "GTC",
  "strategyId": "strategy123"
}
```

**Response:**

```json
{
  "id": "order123",
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "quantity": 10,
  "price": 175.50,
  "timeInForce": "GTC",
  "status": "OPEN",
  "filledQuantity": 0,
  "averagePrice": null,
  "strategyId": "strategy123",
  "createdAt": "2023-04-02T16:45:30Z",
  "updatedAt": "2023-04-02T16:45:30Z"
}
```

#### Place Bracket Order

```
POST /orders/bracket
```

**Request Body:**

```json
{
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "quantity": 10,
  "price": 175.50,
  "timeInForce": "GTC",
  "takeProfitPrice": 180.00,
  "stopLossPrice": 170.00,
  "strategyId": "strategy123"
}
```

**Response:**

```json
{
  "mainOrder": {
    "id": "order123",
    "symbol": "AAPL",
    "side": "BUY",
    "type": "LIMIT",
    "quantity": 10,
    "price": 175.50,
    "timeInForce": "GTC",
    "status": "OPEN",
    "filledQuantity": 0,
    "averagePrice": null,
    "strategyId": "strategy123",
    "createdAt": "2023-04-02T16:45:30Z",
    "updatedAt": "2023-04-02T16:45:30Z"
  },
  "takeProfitOrder": {
    "id": "order124",
    "symbol": "AAPL",
    "side": "SELL",
    "type": "LIMIT",
    "quantity": 10,
    "price": 180.00,
    "timeInForce": "GTC",
    "status": "PENDING",
    "parentId": "order123",
    "strategyId": "strategy123",
    "createdAt": "2023-04-02T16:45:30Z",
    "updatedAt": "2023-04-02T16:45:30Z"
  },
  "stopLossOrder": {
    "id": "order125",
    "symbol": "AAPL",
    "side": "SELL",
    "type": "STOP",
    "quantity": 10,
    "stopPrice": 170.00,
    "timeInForce": "GTC",
    "status": "PENDING",
    "parentId": "order123",
    "strategyId": "strategy123",
    "createdAt": "2023-04-02T16:45:30Z",
    "updatedAt": "2023-04-02T16:45:30Z"
  }
}
```

#### Get Order

```
GET /orders/{orderId}
```

**Parameters:**

| Parameter | Type   | Required | Description    |
|-----------|--------|----------|----------------|
| orderId   | string | Yes      | Order ID       |

**Response:**

```json
{
  "id": "order123",
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "quantity": 10,
  "price": 175.50,
  "timeInForce": "GTC",
  "status": "PARTIALLY_FILLED",
  "filledQuantity": 5,
  "averagePrice": 175.45,
  "strategyId": "strategy123",
  "createdAt": "2023-04-02T16:45:30Z",
  "updatedAt": "2023-04-02T16:50:15Z",
  "fills": [
    {
      "id": "fill123",
      "price": 175.45,
      "quantity": 5,
      "timestamp": "2023-04-02T16:50:15Z"
    }
  ]
}
```

#### Get Orders

```
GET /orders?status=OPEN&symbol=AAPL&limit=10&page=1
```

**Parameters:**

| Parameter | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| status    | string | No       | Filter by status (OPEN, FILLED, PARTIALLY_FILLED, CANCELLED, REJECTED)     |
| symbol    | string | No       | Filter by symbol                                                           |
| side      | string | No       | Filter by side (BUY, SELL)                                                 |
| from      | string | No       | Filter by start date/time                                                  |
| to        | string | No       | Filter by end date/time                                                    |
| limit     | number | No       | Maximum number of results to return (default: 10, max: 100)                |
| page      | number | No       | Page number for pagination (default: 1)                                    |

**Response:**

```json
{
  "orders": [
    {
      "id": "order123",
      "symbol": "AAPL",
      "side": "BUY",
      "type": "LIMIT",
      "quantity": 10,
      "price": 175.50,
      "timeInForce": "GTC",
      "status": "PARTIALLY_FILLED",
      "filledQuantity": 5,
      "averagePrice": 175.45,
      "strategyId": "strategy123",
      "createdAt": "2023-04-02T16:45:30Z",
      "updatedAt": "2023-04-02T16:50:15Z"
    },
    {
      "id": "order126",
      "symbol": "AAPL",
      "side": "BUY",
      "type": "LIMIT",
      "quantity": 5,
      "price": 174.80,
      "timeInForce": "GTC",
      "status": "OPEN",
      "filledQuantity": 0,
      "averagePrice": null,
      "strategyId": "strategy123",
      "createdAt": "2023-04-02T17:15:30Z",
      "updatedAt": "2023-04-02T17:15:30Z"
    }
    // Additional orders...
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalItems": 2,
    "totalPages": 1
  }
}
```

#### Cancel Order

```
DELETE /orders/{orderId}
```

**Parameters:**

| Parameter | Type   | Required | Description    |
|-----------|--------|----------|----------------|
| orderId   | string | Yes      | Order ID       |

**Response:**

```json
{
  "id": "order123",
  "status": "CANCELLED",
  "updatedAt": "2023-04-02T17:30:45Z"
}
```

#### Modify Order

```
PUT /orders/{orderId}
```

**Parameters:**

| Parameter | Type   | Required | Description    |
|-----------|--------|----------|----------------|
| orderId   | string | Yes      | Order ID       |

**Request Body:**

```json
{
  "quantity": 15,
  "price": 175.75
}
```

**Response:**

```json
{
  "id": "order123",
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "quantity": 15,
  "price": 175.75,
  "timeInForce": "GTC",
  "status": "OPEN",
  "filledQuantity": 0,
  "averagePrice": null,
  "strategyId": "strategy123",
  "createdAt": "2023-04-02T16:45:30Z",
  "updatedAt": "2023-04-02T17:35:20Z"
}
```

### Portfolio Management

#### Get Portfolio Summary

```
GET /portfolio/summary
```

**Response:**

```json
{
  "totalValue": 125750.45,
  "cashBalance": 25000.75,
  "investedValue": 100749.70,
  "dayPnL": 1250.30,
  "dayPnLPercent": 1.01,
  "totalPnL": 15750.45,
  "totalPnLPercent": 14.32,
  "allocation": {
    "stocks": 65.5,
    "options": 15.2,
    "etfs": 19.3
  },
  "timestamp": "2023-04-02T18:00:00Z"
}
```

#### Get Positions

```
GET /portfolio/positions
```

**Response:**

```json
{
  "positions": [
    {
      "symbol": "AAPL",
      "quantity": 50,
      "averageCost": 165.30,
      "currentPrice": 175.25,
      "marketValue": 8762.50,
      "unrealizedPnL": 497.50,
      "unrealizedPnLPercent": 6.02,
      "dayPnL": 117.50,
      "dayPnLPercent": 1.36,
      "allocation": 8.70,
      "lastUpdated": "2023-04-02T18:00:00Z"
    },
    {
      "symbol": "MSFT",
      "quantity": 30,
      "averageCost": 280.50,
      "currentPrice": 305.75,
      "marketValue": 9172.50,
      "unrealizedPnL": 757.50,
      "unrealizedPnLPercent": 9.00,
      "dayPnL": 135.00,
      "dayPnLPercent": 1.49,
      "allocation": 9.10,
      "lastUpdated": "2023-04-02T18:00:00Z"
    }
    // Additional positions...
  ]
}
```

#### Get Position

```
GET /portfolio/positions/{symbol}
```

**Parameters:**

| Parameter | Type   | Required | Description                |
|-----------|--------|----------|----------------------------|
| symbol    | string | Yes      | Symbol of the instrument   |

**Response:**

```json
{
  "symbol": "AAPL",
  "quantity": 50,
  "averageCost": 165.30,
  "currentPrice": 175.25,
  "marketValue": 8762.50,
  "unrealizedPnL": 497.50,
  "unrealizedPnLPercent": 6.02,
  "dayPnL": 117.50,
  "dayPnLPercent": 1.36,
  "allocation": 8.70,
  "trades": [
    {
      "id": "trade123",
      "orderId": "order123",
      "side": "BUY",
      "quantity": 30,
      "price": 160.25,
      "timestamp": "2023-03-15T14:30:45Z"
    },
    {
      "id": "trade124",
      "orderId": "order125",
      "side": "BUY",
      "quantity": 20,
      "price": 172.88,
      "timestamp": "2023-03-28T10:15:30Z"
    }
  ],
  "lastUpdated": "2023-04-02T18:00:00Z"
}
```

#### Get Portfolio Performance

```
GET /portfolio/performance?period=1m&interval=1d
```

**Parameters:**

| Parameter | Type   | Required | Description                                                        |
|-----------|--------|----------|--------------------------------------------------------------------|
| period    | string | No       | Time period (1d, 1w, 1m, 3m, 6m, 1y, ytd, all) (default: 1m)       |
| interval  | string | No       | Time interval for data points (1h, 1d, 1w, 1m) (default: 1d)       |
| benchmark | string | No       | Symbol of benchmark to compare against (e.g., SPY)                 |

**Response:**

```json
{
  "period": "1m",
  "interval": "1d",
  "startDate": "2023-03-02T00:00:00Z",
  "endDate": "2023-04-02T00:00:00Z",
  "startValue": 115000.25,
  "endValue": 125750.45,
  "absoluteReturn": 10750.20,
  "percentReturn": 9.35,
  "benchmark": {
    "symbol": "SPY",
    "startValue": 100,
    "endValue": 105.75,
    "absoluteReturn": 5.75,
    "percentReturn": 5.75
  },
  "alpha": 3.60,
  "beta": 1.15,
  "sharpeRatio": 1.85,
  "maxDrawdown": -3.25,
  "volatility": 12.50,
  "dataPoints": [
    {
      "date": "2023-03-02T00:00:00Z",
      "value": 115000.25,
      "benchmarkValue": 100
    },
    {
      "date": "2023-03-03T00:00:00Z",
      "value": 116250.75,
      "benchmarkValue": 101.25
    },
    // Additional data points...
    {
      "date": "2023-04-02T00:00:00Z",
      "value": 125750.45,
      "benchmarkValue": 105.75
    }
  ]
}
```

### Strategy Management

#### Create Strategy

```
POST /strategies
```

**Request Body:**

```json
{
  "name": "Moving Average Crossover",
  "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
  "symbols": ["AAPL", "MSFT", "GOOG"],
  "type": "TREND_FOLLOWING",
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.01,
    "exitThreshold": -0.01
  },
  "riskManagement": {
    "maxPositionSize": 10,
    "stopLossPercent": 5,
    "takeProfitPercent": 15
  },
  "active": false
}
```

**Response:**

```json
{
  "id": "strategy123",
  "name": "Moving Average Crossover",
  "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
  "symbols": ["AAPL", "MSFT", "GOOG"],
  "type": "TREND_FOLLOWING",
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.01,
    "exitThreshold": -0.01
  },
  "riskManagement": {
    "maxPositionSize": 10,
    "stopLossPercent": 5,
    "takeProfitPercent": 15
  },
  "active": false,
  "createdAt": "2023-04-02T19:15:30Z",
  "updatedAt": "2023-04-02T19:15:30Z"
}
```

#### Get Strategies

```
GET /strategies
```

**Response:**

```json
{
  "strategies": [
    {
      "id": "strategy123",
      "name": "Moving Average Crossover",
      "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
      "symbols": ["AAPL", "MSFT", "GOOG"],
      "type": "TREND_FOLLOWING",
      "active": false,
      "createdAt": "2023-04-02T19:15:30Z",
      "updatedAt": "2023-04-02T19:15:30Z"
    },
    {
      "id": "strategy124",
      "name": "RSI Oversold",
      "description": "Buy when RSI goes below 30, sell when it goes above 70",
      "symbols": ["AMZN", "NFLX", "META"],
      "type": "MEAN_REVERSION",
      "active": true,
      "createdAt": "2023-04-01T15:30:45Z",
      "updatedAt": "2023-04-02T10:20:15Z"
    }
    // Additional strategies...
  ]
}
```

#### Get Strategy

```
GET /strategies/{strategyId}
```

**Parameters:**

| Parameter  | Type   | Required | Description    |
|------------|--------|----------|----------------|
| strategyId | string | Yes      | Strategy ID    |

**Response:**

```json
{
  "id": "strategy123",
  "name": "Moving Average Crossover",
  "description": "Buy when 50-day MA crosses above 200-day MA, sell when it crosses below",
  "symbols": ["AAPL", "MSFT", "GOOG"],
  "type": "TREND_FOLLOWING",
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.01,
    "exitThreshold": -0.01
  },
  "riskManagement": {
    "maxPositionSize": 10,
    "stopLossPercent": 5,
    "takeProfitPercent": 15
  },
  "active": false,
  "performance": {
    "totalTrades": 25,
    "winningTrades": 15,
    "losingTrades": 10,
    "winRate": 60,
    "averageWin": 8.5,
    "averageLoss": -3.2,
    "profitFactor": 2.65,
    "totalReturn": 75.25,
    "annualizedReturn": 15.5,
    "maxDrawdown": -12.3
  },
  "createdAt": "2023-04-02T19:15:30Z",
  "updatedAt": "2023-04-02T19:15:30Z"
}
```

#### Update Strategy

```
PUT /strategies/{strategyId}
```

**Parameters:**

| Parameter  | Type   | Required | Description    |
|------------|--------|----------|----------------|
| strategyId | string | Yes      | Strategy ID    |

**Request Body:**

```json
{
  "name": "Enhanced Moving Average Crossover",
  "description": "Buy when 50-day MA crosses above 200-day MA with volume confirmation",
  "symbols": ["AAPL", "MSFT", "GOOG", "AMZN"],
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.02,
    "exitThreshold": -0.02,
    "volumeThreshold": 1.5
  },
  "riskManagement": {
    "maxPositionSize": 5,
    "stopLossPercent": 7,
    "takeProfitPercent": 20
  },
  "active": true
}
```

**Response:**

```json
{
  "id": "strategy123",
  "name": "Enhanced Moving Average Crossover",
  "description": "Buy when 50-day MA crosses above 200-day MA with volume confirmation",
  "symbols": ["AAPL", "MSFT", "GOOG", "AMZN"],
  "type": "TREND_FOLLOWING",
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.02,
    "exitThreshold": -0.02,
    "volumeThreshold": 1.5
  },
  "riskManagement": {
    "maxPositionSize": 5,
    "stopLossPercent": 7,
    "takeProfitPercent": 20
  },
  "active": true,
  "createdAt": "2023-04-02T19:15:30Z",
  "updatedAt": "2023-04-02T20:30:45Z"
}
```

#### Delete Strategy

```
DELETE /strategies/{strategyId}
```

**Parameters:**

| Parameter  | Type   | Required | Description    |
|------------|--------|----------|----------------|
| strategyId | string | Yes      | Strategy ID    |

**Response:**

```json
{
  "id": "strategy123",
  "deleted": true
}
```

#### Backtest Strategy

```
POST /strategies/{strategyId}/backtest
```

**Parameters:**

| Parameter  | Type   | Required | Description    |
|------------|--------|----------|----------------|
| strategyId | string | Yes      | Strategy ID    |

**Request Body:**

```json
{
  "startDate": "2022-01-01",
  "endDate": "2023-01-01",
  "initialCapital": 100000,
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.01,
    "exitThreshold": -0.01
  }
}
```

**Response:**

```json
{
  "id": "backtest123",
  "strategyId": "strategy123",
  "startDate": "2022-01-01T00:00:00Z",
  "endDate": "2023-01-01T00:00:00Z",
  "initialCapital": 100000,
  "finalCapital": 125750.45,
  "totalReturn": 25.75,
  "annualizedReturn": 25.75,
  "sharpeRatio": 1.85,
  "maxDrawdown": -15.30,
  "trades": [
    {
      "symbol": "AAPL",
      "entryDate": "2022-02-15T00:00:00Z",
      "entryPrice": 160.25,
      "quantity": 50,
      "exitDate": "2022-04-20T00:00:00Z",
      "exitPrice": 175.50,
      "pnl": 762.50,
      "pnlPercent": 9.52
    },
    // Additional trades...
  ],
  "equityCurve": [
    {
      "date": "2022-01-01T00:00:00Z",
      "equity": 100000
    },
    // Additional equity curve points...
    {
      "date": "2023-01-01T00:00:00Z",
      "equity": 125750.45
    }
  ],
  "parameters": {
    "shortPeriod": 50,
    "longPeriod": 200,
    "entryThreshold": 0.01,
    "exitThreshold": -0.01
  },
  "createdAt": "2023-04-02T21:00:15Z"
}
```

### WebSocket API

The WebSocket API provides real-time updates for market data, orders, and account information.

#### Connection

Connect to the WebSocket server:

```
wss://api.tradingplatform.com/v1/ws
```

Authentication is required via a token query parameter:

```
wss://api.tradingplatform.com/v1/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Message Format

All messages follow this format:

```json
{
  "type": "MESSAGE_TYPE",
  "data": {
    // Message-specific data
  }
}
```

#### Subscription

Subscribe to specific data channels:

```json
{
  "type": "SUBSCRIBE",
  "channels": [
    {
      "name": "QUOTES",
      "symbols": ["AAPL", "MSFT", "GOOG"]
    },
    {
      "name": "ORDERS"
    },
    {
      "name": "POSITIONS"
    }
  ]
}
```

#### Quote Updates

Real-time quote updates:

```json
{
  "type": "QUOTE",
  "data": {
    "symbol": "AAPL",
    "lastPrice": 175.25,
    "change": 2.35,
    "changePercent": 1.36,
    "bidPrice": 175.20,
    "bidSize": 500,
    "askPrice": 175.30,
    "askSize": 300,
    "volume": 15234567,
    "timestamp": "2023-04-02T16:30:00.123Z"
  }
}
```

#### Order Updates

Real-time order status updates:

```json
{
  "type": "ORDER_UPDATE",
  "data": {
    "id": "order123",
    "symbol": "AAPL",
    "side": "BUY",
    "type": "LIMIT",
    "quantity": 10,
    "price": 175.50,
    "timeInForce": "GTC",
    "status": "PARTIALLY_FILLED",
    "filledQuantity": 5,
    "averagePrice": 175.45,
    "strategyId": "strategy123",
    "updatedAt": "2023-04-02T16:50:15.456Z"
  }
}
```

#### Position Updates

Real-time position updates:

```json
{
  "type": "POSITION_UPDATE",
  "data": {
    "symbol": "AAPL",
    "quantity": 50,
    "averageCost": 165.30,
    "currentPrice": 175.25,
    "marketValue": 8762.50,
    "unrealizedPnL": 497.50,
    "unrealizedPnLPercent": 6.02,
    "dayPnL": 117.50,
    "dayPnLPercent": 1.36,
    "lastUpdated": "2023-04-02T18:00:00.789Z"
  }
}
```

## Error Handling

### Error Format

All API errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      // Additional error details (optional)
    }
  }
}
```

### Common Error Codes

| Code                  | HTTP Status | Description                                           |
|-----------------------|-------------|-------------------------------------------------------|
| AUTHENTICATION_FAILED | 401         | Invalid or expired authentication token               |
| AUTHORIZATION_FAILED  | 403         | Insufficient permissions for the requested operation  |
| RESOURCE_NOT_FOUND    | 404         | The requested resource does not exist                 |
| VALIDATION_ERROR      | 400         | Invalid request parameters or body                    |
| RATE_LIMIT_EXCEEDED   | 429         | Too many requests, rate limit exceeded                |
| INTERNAL_ERROR        | 500         | Internal server error                                 |

## Pagination

For endpoints that return collections of resources, pagination is supported with the following parameters:

| Parameter | Type   | Description                                                |
|-----------|--------|------------------------------------------------------------|
| limit     | number | Maximum number of results to return (default: 10, max: 100)|
| page      | number | Page number (default: 1)                                   |

Pagination information is included in the response:

```json
{
  "data": [
    // Resource items...
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalItems": 42,
    "totalPages": 5
  }
}
```

## Versioning

The API uses versioning in the URL path (e.g., `/v1/orders`). When breaking changes are introduced, a new version will be released. The current version is v1.

## Code Examples

### Authentication (JavaScript)

```javascript
async function login(email, password) {
  try {
    const response = await fetch('https://api.tradingplatform.com/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, password })
    });
    
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error.message);
    }
    
    const data = await response.json();
    
    // Store tokens
    localStorage.setItem('token', data.token);
    localStorage.setItem('refreshToken', data.refreshToken);
    
    return data;
  } catch (error) {
    console.error('Login failed:', error);
    throw error;
  }
}
```

### Placing an Order (Python)

```python
import requests
import json

def place_order(token, symbol, side, order_type, quantity, price=None):
    url = "https://api.tradingplatform.com/v1/orders"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    }
    
    payload = {
        "symbol": symbol,
        "side": side,
        "type": order_type,
        "quantity": quantity,
        "timeInForce": "GTC"
    }
    
    if price is not None and order_type != "MARKET":
        payload["price"] = price
    
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    
    if response.status_code != 200:
        error_data = response.json()
        raise Exception(f"Order placement failed: {error_data['error']['message']}")
    
    return response.json()

# Example usage
try:
    order = place_order(
        token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        symbol="AAPL",
        side="BUY",
        order_type="LIMIT",
        quantity=10,
        price=175.50
    )
    print(f"Order placed successfully: {order['id']}")
except Exception as e:
    print(e)
```

### WebSocket Connection (JavaScript)

```javascript
class TradingWebSocket {
  constructor(token) {
    this.token = token;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
    this.handlers = {};
  }
  
  connect() {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(`wss://api.tradingplatform.com/v1/ws?token=${this.token}`);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        resolve();
      };
      
      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          if (this.handlers[message.type]) {
            this.handlers[message.type].forEach(handler => handler(message.data));
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
      
      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        this.attemptReconnect();
      };
      
      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };
    });
  }
  
  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
      
      console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`);
      
      setTimeout(() => {
        this.connect().catch(console.error);
      }, delay);
    } else {
      console.error('Maximum reconnect attempts reached');
    }
  }
  
  subscribe(channels) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        type: 'SUBSCRIBE',
        channels
      }));
    } else {
      throw new Error('WebSocket not connected');
    }
  }
  
  on(messageType, handler) {
    if (!this.handlers[messageType]) {
      this.handlers[messageType] = [];
    }
    this.handlers[messageType].push(handler);
  }
  
  off(messageType, handler) {
    if (this.handlers[messageType]) {
      this.handlers[messageType] = this.handlers[messageType].filter(h => h !== handler);
    }
  }
  
  close() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Example usage
async function connectToWebSocket() {
  const token = localStorage.getItem('token');
  const ws = new TradingWebSocket(token);
  
  try {
    await ws.connect();
    
    // Subscribe to channels
    ws.subscribe([
      {
        name: 'QUOTES',
        symbols: ['AAPL', 'MSFT', 'GOOG']
      },
      {
        name: 'ORDERS'
      }
    ]);
    
    // Register handlers
    ws.on('QUOTE', (data) => {
      console.log(`Quote update for ${data.symbol}: ${data.lastPrice}`);
      // Update UI with quote data
    });
    
    ws.on('ORDER_UPDATE', (data) => {
      console.log(`Order ${data.id} updated: ${data.status}`);
      // Update UI with order data
    });
    
  } catch (error) {
    console.error('Failed to connect to WebSocket:', error);
  }
  
  return ws;
}
```

## Conclusion

This API documentation provides a comprehensive reference for integrating with the Trading Platform API. For additional support or questions, please contact our developer support team at api-support@tradingplatform.com.

Remember to keep your API keys and authentication tokens secure, and follow best practices for error handling and reconnection logic in your applications.
