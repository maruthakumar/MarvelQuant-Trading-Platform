# API Documentation with Frontend Usage Examples

## Overview

This document provides comprehensive documentation for the Trading Platform API with practical frontend usage examples. It covers authentication, error handling, and specific API endpoints with TypeScript code examples for frontend integration.

## Table of Contents

1. [API Client Overview](#api-client-overview)
2. [Authentication](#authentication)
3. [Error Handling](#error-handling)
4. [Order Management API](#order-management-api)
5. [Position Management API](#position-management-api)
6. [Strategy Management API](#strategy-management-api)
7. [User Management API](#user-management-api)
8. [Market Data API](#market-data-api)
9. [WebSocket Integration](#websocket-integration)
10. [Best Practices](#best-practices)

## API Client Overview

The Trading Platform provides a robust API client for frontend integration. There are two main approaches to using the API in frontend components:

1. Using the `apiClient singleton` for direct API calls
2. Using `specialized API services` for domain-specific operations

### API Client Singleton

The `apiClient` singleton provides a unified interface for making API requests with proper error handling, request/response interceptors, and retry capabilities.

```typescript
// Import the API client
import apiClient from '../api/apiClient';

// Example usage in a component
const fetchData = async () => {
  try {
    const data = await apiClient.get('/orders');
    setOrders(data);
  } catch (error) {
    handleError(error);
  }
};
```

### Specialized API Services

The Trading Platform also provides specialized API services for domain-specific operations, which are built on top of the base API client.

```typescript
// Import specialized API services
import { orderApi, positionApi, strategyApi } from '../services/apiService';

// Example usage in a component
const fetchOrders = async () => {
  try {
    const orders = await orderApi.getOrders({ status: 'OPEN' });
    setOrders(orders);
  } catch (error) {
    handleError(error);
  }
};
```

## Authentication

### Login Flow

The authentication flow involves sending credentials to the server and storing the returned token for subsequent requests.

#### Backend API Endpoint

```
POST /api/auth/login
```

Request body:
```json
{
  "username": "string",
  "password": "string"
}
```

Response:
```json
{
  "token": "string",
  "expiresAt": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string"
  }
}
```

#### Frontend Implementation

Using the specialized auth API service:

```typescript
import { authApi } from '../services/apiService';
import { useDispatch } from 'react-redux';
import { loginSuccess, loginFailure } from '../store/slices/authSlice';

const LoginComponent = () => {
  const dispatch = useDispatch();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      const response = await authApi.login(username, password);
      
      // Store token in localStorage
      localStorage.setItem('token', response.token);
      
      // Update Redux state
      dispatch(loginSuccess({
        token: response.token,
        user: response.user
      }));
      
      // Redirect to dashboard
      navigate('/dashboard');
    } catch (error) {
      setError(
        error.response?.data?.message || 
        'Login failed. Please check your credentials.'
      );
      dispatch(loginFailure(error.response?.data?.message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleLogin}>
      {/* Form fields */}
      {error && <div className="error">{error}</div>}
      <button type="submit" disabled={loading}>
        {loading ? 'Logging in...' : 'Login'}
      </button>
    </form>
  );
};
```

### Token Validation

The token validation flow involves checking if the current token is still valid.

#### Backend API Endpoint

```
GET /api/auth/validate
```

Response:
```json
{
  "valid": true,
  "user": {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string"
  }
}
```

#### Frontend Implementation

```typescript
import { authApi } from '../services/apiService';
import { useDispatch } from 'react-redux';
import { validateTokenSuccess, validateTokenFailure } from '../store/slices/authSlice';

const validateToken = async () => {
  const dispatch = useDispatch();
  
  try {
    const response = await authApi.validateToken();
    
    if (response.valid) {
      dispatch(validateTokenSuccess({
        user: response.user
      }));
      return true;
    } else {
      // Token is invalid
      localStorage.removeItem('token');
      dispatch(validateTokenFailure());
      return false;
    }
  } catch (error) {
    // Token validation failed
    localStorage.removeItem('token');
    dispatch(validateTokenFailure());
    return false;
  }
};
```

### Logout Flow

The logout flow involves invalidating the token on the server and removing it from local storage.

#### Backend API Endpoint

```
POST /api/auth/logout
```

Response:
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### Frontend Implementation

```typescript
import { authApi } from '../services/apiService';
import { useDispatch } from 'react-redux';
import { logout } from '../store/slices/authSlice';

const LogoutButton = () => {
  const dispatch = useDispatch();
  const [loading, setLoading] = useState(false);

  const handleLogout = async () => {
    setLoading(true);
    
    try {
      await authApi.logout();
      
      // Remove token from localStorage
      localStorage.removeItem('token');
      
      // Update Redux state
      dispatch(logout());
      
      // Redirect to login
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
      
      // Even if the server request fails, we still want to log out locally
      localStorage.removeItem('token');
      dispatch(logout());
      navigate('/login');
    } finally {
      setLoading(false);
    }
  };

  return (
    <button onClick={handleLogout} disabled={loading}>
      {loading ? 'Logging out...' : 'Logout'}
    </button>
  );
};
```

## Error Handling

The Trading Platform implements comprehensive error handling for API requests. This section demonstrates how to handle different types of errors in frontend components.

### API Client Error Handling

The `apiClient` includes built-in error handling with retry capabilities for network errors and server errors (5xx).

```typescript
// From apiClient.ts
private setupInterceptors(): void {
  // Response interceptor
  this.client.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;
      
      // Handle token expiration (401 Unauthorized)
      if (error.response && error.response.status === 401) {
        // Clear token and redirect to login
        localStorage.removeItem('token');
        window.location.href = '/login';
        return Promise.reject(error);
      }
      
      // Implement retry logic for network errors or 5xx server errors
      if (
        (error.message.includes('Network Error') || 
        (error.response && error.response.status >= 500)) && 
        originalRequest && 
        !originalRequest._retry &&
        originalRequest._retryCount < this.retryCount
      ) {
        originalRequest._retry = true;
        originalRequest._retryCount = (originalRequest._retryCount || 0) + 1;
        
        // Exponential backoff
        const delay = this.retryDelay * Math.pow(2, originalRequest._retryCount - 1);
        
        return new Promise(resolve => {
          setTimeout(() => resolve(this.client(originalRequest)), delay);
        });
      }
      
      return Promise.reject(error);
    }
  );
}
```

### Component-Level Error Handling

In React components, you should implement try-catch blocks to handle API errors gracefully:

```typescript
import React, { useState, useEffect } from 'react';
import { orderApi } from '../services/apiService';

const OrderList = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await orderApi.getOrders();
        setOrders(data);
      } catch (error) {
        console.error('Error fetching orders:', error);
        
        // Handle different error types
        if (error.response) {
          // Server responded with an error status
          if (error.response.status === 404) {
            setError('No orders found.');
          } else if (error.response.status === 403) {
            setError('You do not have permission to view orders.');
          } else {
            setError(`Server error: ${error.response.data.message || 'Unknown error'}`);
          }
        } else if (error.request) {
          // Request was made but no response received
          setError('Network error. Please check your connection and try again.');
        } else {
          // Something else happened while setting up the request
          setError(`Error: ${error.message}`);
        }
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, []);

  if (loading) return <div>Loading orders...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div className="order-list">
      <h2>Your Orders</h2>
      {orders.length === 0 ? (
        <p>No orders found.</p>
      ) : (
        <ul>
          {orders.map(order => (
            <li key={order.id}>
              {order.symbol} - {order.quantity} @ {order.price} - {order.status}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
```

### Error Boundary

For React applications, implement Error Boundaries to catch JavaScript errors anywhere in the component tree:

```typescript
import React, { Component, ErrorInfo, ReactNode } from 'react';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    // Log the error to an error reporting service
    console.error('Error caught by boundary:', error, errorInfo);
    
    // You could also send this to a logging service like Sentry
    // logErrorToService(error, errorInfo);
  }

  render(): ReactNode {
    if (this.state.hasError) {
      // You can render any custom fallback UI
      return this.props.fallback || (
        <div className="error-boundary">
          <h2>Something went wrong.</h2>
          <p>{this.state.error?.message || 'Unknown error'}</p>
          <button onClick={() => window.location.reload()}>
            Reload Page
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

// Usage:
// <ErrorBoundary>
//   <YourComponent />
// </ErrorBoundary>
```

## Order Management API

The Order Management API allows you to create, retrieve, update, and cancel orders.

### Get Orders

#### Backend API Endpoint

```
GET /api/orders
```

Query parameters:
- `status` (optional): Filter by order status (e.g., 'OPEN', 'FILLED', 'CANCELLED')
- `symbol` (optional): Filter by symbol
- `from` (optional): Filter by start date (ISO format)
- `to` (optional): Filter by end date (ISO format)
- `limit` (optional): Limit the number of results (default: 100)
- `offset` (optional): Offset for pagination (default: 0)

Response:
```json
{
  "orders": [
    {
      "id": "string",
      "symbol": "string",
      "quantity": "number",
      "price": "number",
      "side": "string",
      "type": "string",
      "status": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ],
  "total": "number",
  "limit": "number",
  "offset": "number"
}
```

#### Frontend Implementation

Using the specialized order API service:

```typescript
import React, { useState, useEffect } from 'react';
import { orderApi } from '../services/apiService';

const OrderListComponent = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    status: 'OPEN',
    symbol: '',
    from: '',
    to: '',
    limit: 100,
    offset: 0
  });

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await orderApi.getOrders(filters);
        setOrders(data.orders);
      } catch (error) {
        console.error('Error fetching orders:', error);
        setError('Failed to fetch orders. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [filters]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handlePageChange = (newOffset) => {
    setFilters(prev => ({
      ...prev,
      offset: newOffset
    }));
  };

  if (loading) return <div>Loading orders...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div className="order-list">
      <h2>Your Orders</h2>
      
      {/* Filter controls */}
      <div className="filters">
        <select 
          name="status" 
          value={filters.status} 
          onChange={handleFilterChange}
        >
          <option value="">All Statuses</option>
          <option value="OPEN">Open</option>
          <option value="FILLED">Filled</option>
          <option value="CANCELLED">Cancelled</option>
        </select>
        
        <input
          type="text"
          name="symbol"
          placeholder="Symbol"
          value={filters.symbol}
          onChange={handleFilterChange}
        />
        
        <input
          type="date"
          name="from"
          value={filters.from}
          onChange={handleFilterChange}
        />
        
        <input
          type="date"
          name="to"
          value={filters.to}
          onChange={handleFilterChange}
        />
      </div>
      
      {/* Order table */}
      <table className="orders-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Symbol</th>
            <th>Side</th>
            <th>Type</th>
            <th>Quantity</th>
            <th>Price</th>
            <th>Status</th>
            <th>Created At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {orders.map(order => (
            <tr key={order.id}>
              <td>{order.id}</td>
              <td>{order.symbol}</td>
              <td>{order.side}</td>
              <td>{order.type}</td>
              <td>{order.quantity}</td>
              <td>{order.price}</td>
              <td>{order.status}</td>
              <td>{new Date(order.createdAt).toLocaleString()}</td>
              <td>
                {order.status === 'OPEN' && (
                  <button onClick={() => handleCancelOrder(order.id)}>
                    Cancel
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      
      {/* Pagination */}
      <div className="pagination">
        <button 
          disabled={filters.offset === 0} 
          onClick={() => handlePageChange(filters.offset - filters.limit)}
        >
          Previous
        </button>
        <span>Page {Math.floor(filters.offset / filters.limit) + 1}</span>
        <button 
          disabled={orders.length < filters.limit} 
          onClick={() => handlePageChange(filters.offset + filters.limit)}
        >
          Next
        </button>
      </div>
    </div>
  );
};
```

### Create Order

#### Backend API Endpoint

```
POST /api/orders
```

Request body:
```json
{
  "symbol": "string",
  "quantity": "number",
  "price": "number",
  "side": "string",
  "type": "string",
  "timeInForce": "string",
  "stopPrice": "number (optional)",
  "clientOrderId": "string (optional)"
}
```

Response:
```json
{
  "id": "string",
  "clientOrderId": "string",
  "symbol": "string",
  "quantity": "number",
  "price": "number",
  "side": "string",
  "type": "string",
  "timeInForce": "string",
  "stopPrice": "number",
  "status": "string",
  "createdAt": "string"
}
```

#### Frontend Implementation

```typescript
import React, { useState } from 'react';
import { orderApi } from '../services/apiService';

const CreateOrderForm = () => {
  const [orderData, setOrderData] = useState({
    symbol: '',
    quantity: '',
    price: '',
    side: 'BUY',
    type: 'LIMIT',
    timeInForce: 'GTC',
    stopPrice: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setOrderData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);
    
    // Validate form
    if (!orderData.symbol) {
      setError('Symbol is required');
      setLoading(false);
      return;
    }
    
    if (!orderData.quantity || parseFloat(orderData.quantity) <= 0) {
      setError('Quantity must be greater than 0');
      setLoading(false);
      return;
    }
    
    if (orderData.type === 'LIMIT' && (!orderData.price || parseFloat(orderData.price) <= 0)) {
      setError('Price must be greater than 0 for limit orders');
      setLoading(false);
      return;
    }
    
    if (orderData.type === 'STOP' && (!orderData.stopPrice || parseFloat(orderData.stopPrice) <= 0)) {
      setError('Stop price must be greater than 0 for stop orders');
      setLoading(false);
      return;
    }
    
    try {
      // Convert string values to numbers
      const formattedOrderData = {
        ...orderData,
        quantity: parseFloat(orderData.quantity),
        price: orderData.price ? parseFloat(orderData.price) : undefined,
        stopPrice: orderData.stopPrice ? parseFloat(orderData.stopPrice) : undefined
      };
      
      const response = await orderApi.createOrder(formattedOrderData);
      
      setSuccess(`Order created successfully! Order ID: ${response.id}`);
      
      // Reset form
      setOrderData({
        symbol: '',
        quantity: '',
        price: '',
        side: 'BUY',
        type: 'LIMIT',
        timeInForce: 'GTC',
        stopPrice: ''
      });
    } catch (error) {
      console.error('Error creating order:', error);
      setError(error.response?.data?.message || 'Failed to create order. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="create-order-form">
      <h2>Create New Order</h2>
      
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="symbol">Symbol</label>
          <input
            type="text"
            id="symbol"
            name="symbol"
            value={orderData.symbol}
            onChange={handleChange}
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="side">Side</label>
          <select
            id="side"
            name="side"
            value={orderData.side}
            onChange={handleChange}
            required
          >
            <option value="BUY">Buy</option>
            <option value="SELL">Sell</option>
          </select>
        </div>
        
        <div className="form-group">
          <label htmlFor="type">Order Type</label>
          <select
            id="type"
            name="type"
            value={orderData.type}
            onChange={handleChange}
            required
          >
            <option value="MARKET">Market</option>
            <option value="LIMIT">Limit</option>
            <option value="STOP">Stop</option>
            <option value="STOP_LIMIT">Stop Limit</option>
          </select>
        </div>
        
        <div className="form-group">
          <label htmlFor="quantity">Quantity</label>
          <input
            type="number"
            id="quantity"
            name="quantity"
            value={orderData.quantity}
            onChange={handleChange}
            min="0.00000001"
            step="0.00000001"
            required
          />
        </div>
        
        {(orderData.type === 'LIMIT' || orderData.type === 'STOP_LIMIT') && (
          <div className="form-group">
            <label htmlFor="price">Price</label>
            <input
              type="number"
              id="price"
              name="price"
              value={orderData.price}
              onChange={handleChange}
              min="0.00000001"
              step="0.00000001"
              required
            />
          </div>
        )}
        
        {(orderData.type === 'STOP' || orderData.type === 'STOP_LIMIT') && (
          <div className="form-group">
            <label htmlFor="stopPrice">Stop Price</label>
            <input
              type="number"
              id="stopPrice"
              name="stopPrice"
              value={orderData.stopPrice}
              onChange={handleChange}
              min="0.00000001"
              step="0.00000001"
              required
            />
          </div>
        )}
        
        <div className="form-group">
          <label htmlFor="timeInForce">Time In Force</label>
          <select
            id="timeInForce"
            name="timeInForce"
            value={orderData.timeInForce}
            onChange={handleChange}
            required
          >
            <option value="GTC">Good Till Cancelled</option>
            <option value="IOC">Immediate Or Cancel</option>
            <option value="FOK">Fill Or Kill</option>
            <option value="DAY">Day</option>
          </select>
        </div>
        
        <button type="submit" disabled={loading}>
          {loading ? 'Creating Order...' : 'Create Order'}
        </button>
      </form>
    </div>
  );
};
```

### Cancel Order

#### Backend API Endpoint

```
DELETE /api/orders/{orderId}
```

Response:
```json
{
  "id": "string",
  "status": "CANCELLED",
  "message": "Order cancelled successfully"
}
```

#### Frontend Implementation

```typescript
import React, { useState } from 'react';
import { orderApi } from '../services/apiService';

const CancelOrderButton = ({ orderId, onSuccess }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleCancelOrder = async () => {
    if (!confirm('Are you sure you want to cancel this order?')) {
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      await orderApi.cancelOrder(orderId);
      
      // Call the success callback
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error('Error cancelling order:', error);
      setError(error.response?.data?.message || 'Failed to cancel order. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button 
        onClick={handleCancelOrder} 
        disabled={loading}
        className="cancel-button"
      >
        {loading ? 'Cancelling...' : 'Cancel Order'}
      </button>
      {error && <div className="error-message">{error}</div>}
    </>
  );
};
```

## Position Management API

The Position Management API allows you to retrieve and manage trading positions.

### Get Positions

#### Backend API Endpoint

```
GET /api/positions
```

Query parameters:
- `symbol` (optional): Filter by symbol
- `status` (optional): Filter by position status (e.g., 'OPEN', 'CLOSED')

Response:
```json
{
  "positions": [
    {
      "id": "string",
      "symbol": "string",
      "quantity": "number",
      "entryPrice": "number",
      "currentPrice": "number",
      "unrealizedPnl": "number",
      "realizedPnl": "number",
      "side": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import { positionApi } from '../services/apiService';

const PositionList = () => {
  const [positions, setPositions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    symbol: '',
    status: 'OPEN'
  });

  useEffect(() => {
    const fetchPositions = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await positionApi.getPositions(filters);
        setPositions(data.positions);
      } catch (error) {
        console.error('Error fetching positions:', error);
        setError('Failed to fetch positions. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchPositions();
  }, [filters]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const calculateTotalPnL = () => {
    return positions.reduce((total, position) => {
      return total + position.unrealizedPnl + position.realizedPnl;
    }, 0);
  };

  if (loading) return <div>Loading positions...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div className="position-list">
      <h2>Your Positions</h2>
      
      {/* Filter controls */}
      <div className="filters">
        <input
          type="text"
          name="symbol"
          placeholder="Symbol"
          value={filters.symbol}
          onChange={handleFilterChange}
        />
        
        <select 
          name="status" 
          value={filters.status} 
          onChange={handleFilterChange}
        >
          <option value="">All Statuses</option>
          <option value="OPEN">Open</option>
          <option value="CLOSED">Closed</option>
        </select>
      </div>
      
      {/* Summary */}
      <div className="positions-summary">
        <div className="summary-item">
          <span>Total Positions:</span>
          <span>{positions.length}</span>
        </div>
        <div className="summary-item">
          <span>Total P&L:</span>
          <span className={calculateTotalPnL() >= 0 ? 'profit' : 'loss'}>
            ${calculateTotalPnL().toFixed(2)}
          </span>
        </div>
      </div>
      
      {/* Positions table */}
      <table className="positions-table">
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Side</th>
            <th>Quantity</th>
            <th>Entry Price</th>
            <th>Current Price</th>
            <th>Unrealized P&L</th>
            <th>Realized P&L</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {positions.length === 0 ? (
            <tr>
              <td colSpan="8">No positions found.</td>
            </tr>
          ) : (
            positions.map(position => (
              <tr key={position.id}>
                <td>{position.symbol}</td>
                <td>{position.side}</td>
                <td>{position.quantity}</td>
                <td>${position.entryPrice.toFixed(2)}</td>
                <td>${position.currentPrice.toFixed(2)}</td>
                <td className={position.unrealizedPnl >= 0 ? 'profit' : 'loss'}>
                  ${position.unrealizedPnl.toFixed(2)}
                </td>
                <td className={position.realizedPnl >= 0 ? 'profit' : 'loss'}>
                  ${position.realizedPnl.toFixed(2)}
                </td>
                <td>
                  {position.status === 'OPEN' && (
                    <button onClick={() => handleSquareOff(position.id)}>
                      Square Off
                    </button>
                  )}
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
};
```

### Square Off Position

#### Backend API Endpoint

```
POST /api/positions/{positionId}/square-off
```

Response:
```json
{
  "id": "string",
  "status": "CLOSED",
  "message": "Position squared off successfully",
  "order": {
    "id": "string",
    "symbol": "string",
    "quantity": "number",
    "price": "number",
    "side": "string",
    "type": "string",
    "status": "string"
  }
}
```

#### Frontend Implementation

```typescript
import React, { useState } from 'react';
import { positionApi } from '../services/apiService';

const SquareOffButton = ({ positionId, onSuccess }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSquareOff = async () => {
    if (!confirm('Are you sure you want to square off this position?')) {
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      const result = await positionApi.squareOffPosition(positionId);
      
      // Call the success callback
      if (onSuccess) {
        onSuccess(result);
      }
    } catch (error) {
      console.error('Error squaring off position:', error);
      setError(error.response?.data?.message || 'Failed to square off position. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button 
        onClick={handleSquareOff} 
        disabled={loading}
        className="square-off-button"
      >
        {loading ? 'Processing...' : 'Square Off'}
      </button>
      {error && <div className="error-message">{error}</div>}
    </>
  );
};
```

## Strategy Management API

The Strategy Management API allows you to create, retrieve, update, and execute trading strategies.

### Get Strategies

#### Backend API Endpoint

```
GET /api/strategies
```

Query parameters:
- `status` (optional): Filter by strategy status (e.g., 'ACTIVE', 'INACTIVE')
- `type` (optional): Filter by strategy type

Response:
```json
{
  "strategies": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "type": "string",
      "status": "string",
      "parameters": "object",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import { strategyApi } from '../services/apiService';

const StrategyList = () => {
  const [strategies, setStrategies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    status: '',
    type: ''
  });

  useEffect(() => {
    const fetchStrategies = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await strategyApi.getStrategies(filters);
        setStrategies(data.strategies);
      } catch (error) {
        console.error('Error fetching strategies:', error);
        setError('Failed to fetch strategies. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchStrategies();
  }, [filters]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  if (loading) return <div>Loading strategies...</div>;
  if (error) return <div className="error-message">{error}</div>;

  return (
    <div className="strategy-list">
      <h2>Your Strategies</h2>
      
      {/* Filter controls */}
      <div className="filters">
        <select 
          name="status" 
          value={filters.status} 
          onChange={handleFilterChange}
        >
          <option value="">All Statuses</option>
          <option value="ACTIVE">Active</option>
          <option value="INACTIVE">Inactive</option>
        </select>
        
        <select 
          name="type" 
          value={filters.type} 
          onChange={handleFilterChange}
        >
          <option value="">All Types</option>
          <option value="TREND_FOLLOWING">Trend Following</option>
          <option value="MEAN_REVERSION">Mean Reversion</option>
          <option value="BREAKOUT">Breakout</option>
          <option value="CUSTOM">Custom</option>
        </select>
      </div>
      
      {/* Strategies table */}
      <table className="strategies-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Status</th>
            <th>Created At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {strategies.length === 0 ? (
            <tr>
              <td colSpan="5">No strategies found.</td>
            </tr>
          ) : (
            strategies.map(strategy => (
              <tr key={strategy.id}>
                <td>{strategy.name}</td>
                <td>{strategy.type}</td>
                <td>{strategy.status}</td>
                <td>{new Date(strategy.createdAt).toLocaleString()}</td>
                <td>
                  <button onClick={() => handleViewStrategy(strategy.id)}>
                    View
                  </button>
                  <button onClick={() => handleEditStrategy(strategy.id)}>
                    Edit
                  </button>
                  <button onClick={() => handleExecuteStrategy(strategy.id)}>
                    Execute
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
};
```

### Execute Strategy

#### Backend API Endpoint

```
POST /api/strategies/{strategyId}/execute
```

Request body:
```json
{
  "parameters": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

Response:
```json
{
  "executionId": "string",
  "strategyId": "string",
  "status": "string",
  "message": "Strategy execution started",
  "orders": [
    {
      "id": "string",
      "symbol": "string",
      "quantity": "number",
      "price": "number",
      "side": "string",
      "type": "string",
      "status": "string"
    }
  ]
}
```

#### Frontend Implementation

```typescript
import React, { useState } from 'react';
import { strategyApi } from '../services/apiService';

const executeStrategy = async (strategyId, parameters) => {
  try {
    const response = await strategyApi.executeStrategy(strategyId, { parameters });
    return response;
  } catch (error) {
    console.error('Error executing strategy:', error);
    throw error;
  }
};
```

## User Management API

The User Management API allows you to manage user profiles, preferences, and API keys.

### Get User Profile

#### Backend API Endpoint

```
GET /api/users/profile
```

Response:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "phone": "string",
  "role": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import { userApi } from '../services/apiService';

const UserProfile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const data = await userApi.getUserProfile();
        setProfile(data);
      } catch (error) {
        console.error('Error fetching user profile:', error);
        setError('Failed to fetch user profile. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) return <div>Loading profile...</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!profile) return <div>No profile data available.</div>;

  return (
    <div className="user-profile">
      <h2>User Profile</h2>
      
      <div className="profile-details">
        <div className="profile-field">
          <span className="field-label">Username:</span>
          <span className="field-value">{profile.username}</span>
        </div>
        
        <div className="profile-field">
          <span className="field-label">Email:</span>
          <span className="field-value">{profile.email}</span>
        </div>
        
        <div className="profile-field">
          <span className="field-label">Name:</span>
          <span className="field-value">{profile.firstName} {profile.lastName}</span>
        </div>
        
        <div className="profile-field">
          <span className="field-label">Phone:</span>
          <span className="field-value">{profile.phone || 'Not provided'}</span>
        </div>
        
        <div className="profile-field">
          <span className="field-label">Role:</span>
          <span className="field-value">{profile.role}</span>
        </div>
        
        <div className="profile-field">
          <span className="field-label">Member Since:</span>
          <span className="field-value">{new Date(profile.createdAt).toLocaleDateString()}</span>
        </div>
      </div>
      
      <button onClick={() => handleEditProfile()}>Edit Profile</button>
    </div>
  );
};
```

### Update User Profile

#### Backend API Endpoint

```
PUT /api/users/profile
```

Request body:
```json
{
  "firstName": "string",
  "lastName": "string",
  "email": "string",
  "phone": "string"
}
```

Response:
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "phone": "string",
  "role": "string",
  "updatedAt": "string"
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import { userApi } from '../services/apiService';

const EditProfileForm = ({ initialProfile, onSuccess, onCancel }) => {
  const [profileData, setProfileData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (initialProfile) {
      setProfileData({
        firstName: initialProfile.firstName || '',
        lastName: initialProfile.lastName || '',
        email: initialProfile.email || '',
        phone: initialProfile.phone || ''
      });
    }
  }, [initialProfile]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setProfileData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    
    try {
      const updatedProfile = await userApi.updateUserProfile(profileData);
      
      // Call the success callback
      if (onSuccess) {
        onSuccess(updatedProfile);
      }
    } catch (error) {
      console.error('Error updating profile:', error);
      setError(error.response?.data?.message || 'Failed to update profile. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="edit-profile-form">
      <h3>Edit Profile</h3>
      
      {error && <div className="error-message">{error}</div>}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="firstName">First Name</label>
          <input
            type="text"
            id="firstName"
            name="firstName"
            value={profileData.firstName}
            onChange={handleChange}
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="lastName">Last Name</label>
          <input
            type="text"
            id="lastName"
            name="lastName"
            value={profileData.lastName}
            onChange={handleChange}
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            value={profileData.email}
            onChange={handleChange}
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="phone">Phone</label>
          <input
            type="tel"
            id="phone"
            name="phone"
            value={profileData.phone}
            onChange={handleChange}
          />
        </div>
        
        <div className="form-actions">
          <button type="submit" disabled={loading}>
            {loading ? 'Saving...' : 'Save Changes'}
          </button>
          <button type="button" onClick={onCancel} disabled={loading}>
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};
```

## Market Data API

The Market Data API provides access to market data, quotes, historical data, and technical indicators.

### Get Quote

#### Backend API Endpoint

```
GET /api/v1/market-data/quote/{symbol}
```

Response:
```json
{
  "status": "success",
  "quote": {
    "symbol": "string",
    "price": "number",
    "change": "number",
    "changePercent": "number",
    "high": "number",
    "low": "number",
    "open": "number",
    "close": "number",
    "volume": "number",
    "timestamp": "string"
  }
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import apiClient from '../api/apiClient';

const StockQuote = ({ symbol }) => {
  const [quote, setQuote] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchQuote = async () => {
      if (!symbol) return;
      
      try {
        setLoading(true);
        setError(null);
        
        const response = await apiClient.get(`/v1/market-data/quote/${symbol}`);
        setQuote(response.quote);
      } catch (error) {
        console.error(`Error fetching quote for ${symbol}:`, error);
        setError('Failed to fetch quote. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchQuote();
    
    // Set up polling for real-time updates
    const intervalId = setInterval(fetchQuote, 5000); // Poll every 5 seconds
    
    return () => clearInterval(intervalId); // Clean up on unmount
  }, [symbol]);

  if (!symbol) return <div>Please specify a symbol.</div>;
  if (loading) return <div>Loading quote for {symbol}...</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!quote) return <div>No quote data available for {symbol}.</div>;

  return (
    <div className="stock-quote">
      <h3>{symbol}</h3>
      
      <div className="quote-price">
        <span className="current-price">${quote.price.toFixed(2)}</span>
        <span className={`price-change ${quote.change >= 0 ? 'positive' : 'negative'}`}>
          {quote.change >= 0 ? '+' : ''}{quote.change.toFixed(2)} ({quote.changePercent.toFixed(2)}%)
        </span>
      </div>
      
      <div className="quote-details">
        <div className="quote-detail">
          <span className="detail-label">Open:</span>
          <span className="detail-value">${quote.open.toFixed(2)}</span>
        </div>
        
        <div className="quote-detail">
          <span className="detail-label">High:</span>
          <span className="detail-value">${quote.high.toFixed(2)}</span>
        </div>
        
        <div className="quote-detail">
          <span className="detail-label">Low:</span>
          <span className="detail-value">${quote.low.toFixed(2)}</span>
        </div>
        
        <div className="quote-detail">
          <span className="detail-label">Close:</span>
          <span className="detail-value">${quote.close.toFixed(2)}</span>
        </div>
        
        <div className="quote-detail">
          <span className="detail-label">Volume:</span>
          <span className="detail-value">{quote.volume.toLocaleString()}</span>
        </div>
        
        <div className="quote-detail">
          <span className="detail-label">Updated:</span>
          <span className="detail-value">{new Date(quote.timestamp).toLocaleString()}</span>
        </div>
      </div>
    </div>
  );
};
```

### Get Historical Data

#### Backend API Endpoint

```
GET /api/v1/market-data/historical/{symbol}
```

Query parameters:
- `interval` (optional): Data interval (e.g., '1d', '1h', '15m')
- `from` (optional): Start date (YYYY-MM-DD)
- `to` (optional): End date (YYYY-MM-DD)

Response:
```json
{
  "status": "success",
  "data": {
    "symbol": "string",
    "interval": "string",
    "from": "string",
    "to": "string",
    "candles": [
      {
        "timestamp": "string",
        "open": "number",
        "high": "number",
        "low": "number",
        "close": "number",
        "volume": "number"
      }
    ]
  }
}
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect } from 'react';
import apiClient from '../api/apiClient';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const HistoricalChart = ({ symbol, interval = '1d', days = 30 }) => {
  const [historicalData, setHistoricalData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchHistoricalData = async () => {
      if (!symbol) return;
      
      try {
        setLoading(true);
        setError(null);
        
        // Calculate from date based on days
        const to = new Date();
        const from = new Date();
        from.setDate(from.getDate() - days);
        
        const toStr = to.toISOString().split('T')[0];
        const fromStr = from.toISOString().split('T')[0];
        
        const response = await apiClient.get(`/v1/market-data/historical/${symbol}`, {
          params: {
            interval,
            from: fromStr,
            to: toStr
          }
        });
        
        setHistoricalData(response.data);
      } catch (error) {
        console.error(`Error fetching historical data for ${symbol}:`, error);
        setError('Failed to fetch historical data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchHistoricalData();
  }, [symbol, interval, days]);

  if (!symbol) return <div>Please specify a symbol.</div>;
  if (loading) return <div>Loading historical data for {symbol}...</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!historicalData || !historicalData.candles || historicalData.candles.length === 0) {
    return <div>No historical data available for {symbol}.</div>;
  }

  // Format data for chart
  const chartData = historicalData.candles.map(candle => ({
    date: new Date(candle.timestamp).toLocaleDateString(),
    close: candle.close,
    open: candle.open,
    high: candle.high,
    low: candle.low,
    volume: candle.volume
  }));

  return (
    <div className="historical-chart">
      <h3>{symbol} Historical Data ({interval})</h3>
      
      <div className="chart-container" style={{ width: '100%', height: 400 }}>
        <ResponsiveContainer>
          <LineChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis domain={['auto', 'auto']} />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="close" stroke="#8884d8" activeDot={{ r: 8 }} name="Close Price" />
          </LineChart>
        </ResponsiveContainer>
      </div>
      
      <div className="chart-controls">
        <div className="interval-selector">
          <label htmlFor="interval">Interval:</label>
          <select 
            id="interval" 
            value={interval} 
            onChange={(e) => setInterval(e.target.value)}
          >
            <option value="1d">Daily</option>
            <option value="1h">Hourly</option>
            <option value="15m">15 Minutes</option>
            <option value="5m">5 Minutes</option>
            <option value="1m">1 Minute</option>
          </select>
        </div>
        
        <div className="period-selector">
          <label htmlFor="period">Period:</label>
          <select 
            id="period" 
            value={days} 
            onChange={(e) => setDays(parseInt(e.target.value))}
          >
            <option value="7">1 Week</option>
            <option value="30">1 Month</option>
            <option value="90">3 Months</option>
            <option value="180">6 Months</option>
            <option value="365">1 Year</option>
          </select>
        </div>
      </div>
    </div>
  );
};
```

## WebSocket Integration

The Trading Platform provides WebSocket integration for real-time data updates.

### WebSocket Connection

#### Backend API Endpoint

```
WebSocket: /api/v1/market-data/stream
```

#### Frontend Implementation

```typescript
import React, { useState, useEffect, useRef } from 'react';

const WebSocketComponent = ({ symbols = [] }) => {
  const [connected, setConnected] = useState(false);
  const [quotes, setQuotes] = useState({});
  const [error, setError] = useState(null);
  const wsRef = useRef(null);

  useEffect(() => {
    // Connect to WebSocket
    const connectWebSocket = () => {
      const token = localStorage.getItem('token');
      const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const wsURL = baseURL.replace(/^http/, 'ws') + '/api/v1/market-data/stream';
      
      // Add token to URL if available
      const url = token ? `${wsURL}?token=${token}` : wsURL;
      
      const ws = new WebSocket(url);
      wsRef.current = ws;
      
      ws.onopen = () => {
        console.log('WebSocket connected');
        setConnected(true);
        setError(null);
        
        // Subscribe to symbols
        if (symbols.length > 0) {
          subscribeToSymbols(symbols);
        }
      };
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          if (data.type === 'quote') {
            // Update quote data
            setQuotes(prev => ({
              ...prev,
              [data.symbol]: data
            }));
          } else if (data.type === 'error') {
            console.error('WebSocket error:', data.message);
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
      
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setError('WebSocket connection error. Please try again.');
      };
      
      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setConnected(false);
        
        // Attempt to reconnect after a delay
        setTimeout(() => {
          if (wsRef.current === ws) { // Only reconnect if this is still the current ws
            connectWebSocket();
          }
        }, 5000);
      };
    };
    
    connectWebSocket();
    
    // Clean up on unmount
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);
  
  // Subscribe to symbols when they change
  useEffect(() => {
    if (connected && symbols.length > 0) {
      subscribeToSymbols(symbols);
    }
  }, [connected, symbols]);
  
  const subscribeToSymbols = (symbols) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      return;
    }
    
    // Unsubscribe from all symbols first
    wsRef.current.send(JSON.stringify({
      action: 'unsubscribe',
      symbols: Object.keys(quotes)
    }));
    
    // Subscribe to new symbols
    wsRef.current.send(JSON.stringify({
      action: 'subscribe',
      symbols: symbols
    }));
  };

  return (
    <div className="websocket-component">
      <div className="connection-status">
        Status: {connected ? 'Connected' : 'Disconnected'}
        {error && <div className="error-message">{error}</div>}
      </div>
      
      <div className="real-time-quotes">
        <h3>Real-Time Quotes</h3>
        
        {Object.keys(quotes).length === 0 ? (
          <p>No quotes available. Please subscribe to symbols.</p>
        ) : (
          <div className="quotes-grid">
            {Object.entries(quotes).map(([symbol, data]) => (
              <div key={symbol} className="quote-card">
                <h4>{symbol}</h4>
                <div className="quote-price">
                  <span className="current-price">${data.price.toFixed(2)}</span>
                  <span className={`price-change ${data.change >= 0 ? 'positive' : 'negative'}`}>
                    {data.change >= 0 ? '+' : ''}{data.change.toFixed(2)} ({data.changePercent.toFixed(2)}%)
                  </span>
                </div>
                <div className="quote-timestamp">
                  Updated: {new Date(data.timestamp).toLocaleTimeString()}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
```

### WebSocket Subscription Management

```typescript
import React, { useState } from 'react';

const SymbolSubscriptionManager = ({ onSubscribe }) => {
  const [symbol, setSymbol] = useState('');
  const [subscribedSymbols, setSubscribedSymbols] = useState([]);

  const handleAddSymbol = (e) => {
    e.preventDefault();
    
    if (!symbol) return;
    
    // Check if already subscribed
    if (subscribedSymbols.includes(symbol)) {
      alert(`Already subscribed to ${symbol}`);
      return;
    }
    
    // Add to subscribed symbols
    const newSymbols = [...subscribedSymbols, symbol.toUpperCase()];
    setSubscribedSymbols(newSymbols);
    
    // Call the onSubscribe callback
    if (onSubscribe) {
      onSubscribe(newSymbols);
    }
    
    // Clear input
    setSymbol('');
  };

  const handleRemoveSymbol = (symbolToRemove) => {
    const newSymbols = subscribedSymbols.filter(s => s !== symbolToRemove);
    setSubscribedSymbols(newSymbols);
    
    // Call the onSubscribe callback
    if (onSubscribe) {
      onSubscribe(newSymbols);
    }
  };

  return (
    <div className="symbol-subscription-manager">
      <h3>Manage Subscriptions</h3>
      
      <form onSubmit={handleAddSymbol} className="add-symbol-form">
        <input
          type="text"
          value={symbol}
          onChange={(e) => setSymbol(e.target.value.toUpperCase())}
          placeholder="Enter symbol (e.g., AAPL)"
        />
        <button type="submit">Subscribe</button>
      </form>
      
      <div className="subscribed-symbols">
        <h4>Subscribed Symbols</h4>
        
        {subscribedSymbols.length === 0 ? (
          <p>No symbols subscribed.</p>
        ) : (
          <ul>
            {subscribedSymbols.map(s => (
              <li key={s}>
                {s}
                <button onClick={() => handleRemoveSymbol(s)} className="remove-button">
                  Unsubscribe
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};
```

## Best Practices

### API Client Configuration

Configure API client with appropriate settings for your environment:

```typescript
// Configure API client for production
import apiClient from '../api/apiClient';

// Set retry configuration
apiClient.setRetryConfig(3, 1000); // 3 retries with 1000ms base delay

// For testing environments, you might want to disable retries
if (process.env.NODE_ENV === 'test') {
  apiClient.setRetryConfig(0, 0);
}
```

### Error Handling Strategy

Implement a consistent error handling strategy across your application:

```typescript
// Create a reusable error handler
const handleApiError = (error, setError, defaultMessage = 'An error occurred. Please try again.') => {
  console.error('API Error:', error);
  
  if (error.response) {
    // Server responded with an error status
    const status = error.response.status;
    const message = error.response.data?.message;
    
    if (status === 401) {
      // Unauthorized - redirect to login
      localStorage.removeItem('token');
      window.location.href = '/login';
      return;
    } else if (status === 403) {
      // Forbidden
      setError('You do not have permission to perform this action.');
    } else if (status === 404) {
      // Not found
      setError('The requested resource was not found.');
    } else if (status === 422) {
      // Validation error
      setError(`Validation error: ${message || 'Please check your input.'}`);
    } else if (status >= 500) {
      // Server error
      setError(`Server error: ${message || 'Please try again later.'}`);
    } else {
      // Other errors
      setError(message || defaultMessage);
    }
  } else if (error.request) {
    // Request was made but no response received
    setError('Network error. Please check your connection and try again.');
  } else {
    // Something else happened while setting up the request
    setError(error.message || defaultMessage);
  }
};

// Usage in a component
const MyComponent = () => {
  const [error, setError] = useState(null);
  
  const fetchData = async () => {
    try {
      const data = await apiClient.get('/some-endpoint');
      // Process data
    } catch (error) {
      handleApiError(error, setError, 'Failed to fetch data. Please try again.');
    }
  };
  
  // ...
};
```

### Loading States

Implement consistent loading states across your application:

```typescript
import React, { useState, useEffect } from 'react';
import { orderApi } from '../services/apiService';

const LoadingStateExample = () => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const result = await orderApi.getOrders();
        setData(result);
      } catch (error) {
        console.error('Error fetching data:', error);
        setError('Failed to fetch data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error-icon">⚠️</div>
        <p className="error-message">{error}</p>
        <button onClick={() => window.location.reload()}>Retry</button>
      </div>
    );
  }

  if (!data) {
    return <div>No data available.</div>;
  }

  return (
    <div className="data-container">
      {/* Render data here */}
    </div>
  );
};
```

### Authentication HOC

Create a Higher-Order Component (HOC) for protected routes:

```typescript
import React, { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { authApi } from '../services/apiService';

const withAuth = (Component) => {
  const AuthenticatedComponent = (props) => {
    const [loading, setLoading] = useState(true);
    const [authenticated, setAuthenticated] = useState(false);
    const location = useLocation();

    useEffect(() => {
      const checkAuth = async () => {
        const token = localStorage.getItem('token');
        
        if (!token) {
          setLoading(false);
          return;
        }
        
        try {
          const response = await authApi.validateToken();
          
          if (response.valid) {
            setAuthenticated(true);
          } else {
            // Token is invalid
            localStorage.removeItem('token');
          }
        } catch (error) {
          console.error('Auth validation error:', error);
          localStorage.removeItem('token');
        } finally {
          setLoading(false);
        }
      };

      checkAuth();
    }, []);

    if (loading) {
      return (
        <div className="auth-loading">
          <div className="loading-spinner"></div>
          <p>Authenticating...</p>
        </div>
      );
    }

    if (!authenticated) {
      // Redirect to login page and save the location they were trying to access
      return <Navigate to="/login" state={{ from: location }} replace />;
    }

    return <Component {...props} />;
  };

  return AuthenticatedComponent;
};

// Usage
const ProtectedComponent = () => {
  return <div>This is a protected component</div>;
};

export default withAuth(ProtectedComponent);
```

### API Service Provider

Create a context provider for API services:

```typescript
import React, { createContext, useContext } from 'react';
import { authApi, orderApi, positionApi, strategyApi, userApi } from '../services/apiService';

// Create context
const ApiContext = createContext(null);

// Create provider
export const ApiProvider = ({ children }) => {
  const apiServices = {
    auth: authApi,
    orders: orderApi,
    positions: positionApi,
    strategies: strategyApi,
    users: userApi
  };

  return (
    <ApiContext.Provider value={apiServices}>
      {children}
    </ApiContext.Provider>
  );
};

// Create hook for using API services
export const useApi = () => {
  const context = useContext(ApiContext);
  
  if (!context) {
    throw new Error('useApi must be used within an ApiProvider');
  }
  
  return context;
};

// Usage in a component
const MyComponent = () => {
  const api = useApi();
  
  const fetchOrders = async () => {
    try {
      const orders = await api.orders.getOrders();
      // Process orders
    } catch (error) {
      // Handle error
    }
  };
  
  // ...
};
```

### WebSocket Provider

Create a context provider for WebSocket connections:

```typescript
import React, { createContext, useContext, useEffect, useRef, useState } from 'react';

// Create context
const WebSocketContext = createContext(null);

// Create provider
export const WebSocketProvider = ({ children }) => {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState(null);
  const [subscriptions, setSubscriptions] = useState({});
  const wsRef = useRef(null);

  // Connect to WebSocket
  useEffect(() => {
    const connectWebSocket = () => {
      const token = localStorage.getItem('token');
      const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const wsURL = baseURL.replace(/^http/, 'ws') + '/api/v1/market-data/stream';
      
      // Add token to URL if available
      const url = token ? `${wsURL}?token=${token}` : wsURL;
      
      const ws = new WebSocket(url);
      wsRef.current = ws;
      
      ws.onopen = () => {
        console.log('WebSocket connected');
        setConnected(true);
        setError(null);
        
        // Resubscribe to previous subscriptions
        if (Object.keys(subscriptions).length > 0) {
          ws.send(JSON.stringify({
            action: 'subscribe',
            symbols: Object.keys(subscriptions)
          }));
        }
      };
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          if (data.type === 'quote' && data.symbol) {
            // Update subscription data
            setSubscriptions(prev => ({
              ...prev,
              [data.symbol]: {
                ...prev[data.symbol],
                data: data
              }
            }));
            
            // Call callback if exists
            if (subscriptions[data.symbol]?.callback) {
              subscriptions[data.symbol].callback(data);
            }
          } else if (data.type === 'error') {
            console.error('WebSocket error:', data.message);
            setError(data.message);
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
      
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setError('WebSocket connection error. Please try again.');
      };
      
      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setConnected(false);
        
        // Attempt to reconnect after a delay
        setTimeout(() => {
          if (wsRef.current === ws) { // Only reconnect if this is still the current ws
            connectWebSocket();
          }
        }, 5000);
      };
    };
    
    connectWebSocket();
    
    // Clean up on unmount
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);
  
  // Subscribe to a symbol
  const subscribe = (symbol, callback) => {
    if (!symbol) return;
    
    // Add to subscriptions
    setSubscriptions(prev => ({
      ...prev,
      [symbol]: {
        ...prev[symbol],
        callback
      }
    }));
    
    // Send subscribe message if connected
    if (connected && wsRef.current) {
      wsRef.current.send(JSON.stringify({
        action: 'subscribe',
        symbols: [symbol]
      }));
    }
  };
  
  // Unsubscribe from a symbol
  const unsubscribe = (symbol) => {
    if (!symbol) return;
    
    // Remove from subscriptions
    setSubscriptions(prev => {
      const newSubscriptions = { ...prev };
      delete newSubscriptions[symbol];
      return newSubscriptions;
    });
    
    // Send unsubscribe message if connected
    if (connected && wsRef.current) {
      wsRef.current.send(JSON.stringify({
        action: 'unsubscribe',
        symbols: [symbol]
      }));
    }
  };
  
  // Get data for a symbol
  const getData = (symbol) => {
    return subscriptions[symbol]?.data || null;
  };
  
  // WebSocket context value
  const value = {
    connected,
    error,
    subscribe,
    unsubscribe,
    getData,
    subscriptions: Object.keys(subscriptions)
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
};

// Create hook for using WebSocket
export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  
  return context;
};

// Usage in a component
const RealTimeQuote = ({ symbol }) => {
  const [quote, setQuote] = useState(null);
  const ws = useWebSocket();
  
  useEffect(() => {
    // Subscribe to symbol
    ws.subscribe(symbol, (data) => {
      setQuote(data);
    });
    
    // Unsubscribe on unmount
    return () => {
      ws.unsubscribe(symbol);
    };
  }, [symbol, ws]);
  
  if (!ws.connected) {
    return <div>Connecting to WebSocket...</div>;
  }
  
  if (!quote) {
    return <div>Waiting for quote data...</div>;
  }
  
  return (
    <div className="real-time-quote">
      <h3>{symbol}</h3>
      <div className="quote-price">${quote.price.toFixed(2)}</div>
      <div className="quote-change">
        {quote.change >= 0 ? '+' : ''}{quote.change.toFixed(2)} ({quote.changePercent.toFixed(2)}%)
      </div>
    </div>
  );
};
```

## Conclusion

This documentation provides comprehensive guidance for integrating frontend components with the Trading Platform API. By following these examples and best practices, developers can create robust, maintainable, and user-friendly trading applications.

For any questions or issues, please contact the platform support team.
