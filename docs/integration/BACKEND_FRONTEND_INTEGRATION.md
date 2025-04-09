# Backend-Frontend Integration

## Overview

This document outlines the integration between the backend and frontend components of the trading platform. The integration connects the React frontend with the Go backend through RESTful APIs and WebSocket connections for real-time data.

## API Configuration

First, let's create a configuration file for API endpoints:

```typescript
// frontend/src/config/api.config.ts
export const API_CONFIG = {
  BASE_URL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api',
  WS_URL: process.env.REACT_APP_WS_URL || 'ws://localhost:8080/ws',
  TIMEOUT: 30000, // 30 seconds
  RETRY_COUNT: 3,
};

export const ENDPOINTS = {
  // Auth endpoints
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    FORGOT_PASSWORD: '/auth/forgot-password',
    RESET_PASSWORD: '/auth/reset-password',
    REFRESH_TOKEN: '/auth/refresh-token',
    LOGOUT: '/auth/logout',
  },
  
  // User endpoints
  USER: {
    PROFILE: '/users/profile',
    UPDATE_PROFILE: '/users/profile',
    CHANGE_PASSWORD: '/users/change-password',
  },
  
  // Order endpoints
  ORDER: {
    LIST: '/orders',
    CREATE: '/orders',
    DETAILS: (id: string) => `/orders/${id}`,
    UPDATE: (id: string) => `/orders/${id}`,
    CANCEL: (id: string) => `/orders/${id}`,
  },
  
  // Position endpoints
  POSITION: {
    LIST: '/positions',
    DETAILS: (id: string) => `/positions/${id}`,
    CLOSE: (id: string) => `/positions/${id}/close`,
  },
  
  // Portfolio endpoints
  PORTFOLIO: {
    SUMMARY: '/portfolio/summary',
    HOLDINGS: '/portfolio/holdings',
    PERFORMANCE: '/portfolio/performance',
    ANALYTICS: '/portfolio/analytics',
  },
  
  // Strategy endpoints
  STRATEGY: {
    LIST: '/strategies',
    CREATE: '/strategies',
    DETAILS: (id: string) => `/strategies/${id}`,
    UPDATE: (id: string) => `/strategies/${id}`,
    DELETE: (id: string) => `/strategies/${id}`,
    ACTIVATE: (id: string) => `/strategies/${id}/activate`,
    DEACTIVATE: (id: string) => `/strategies/${id}/deactivate`,
  },
  
  // Market data endpoints
  MARKET: {
    QUOTES: '/market/quotes',
    QUOTE: (symbol: string) => `/market/quotes/${symbol}`,
    WATCHLISTS: '/market/watchlists',
    WATCHLIST: (id: string) => `/market/watchlists/${id}`,
    WATCHLIST_ADD_SYMBOL: (id: string) => `/market/watchlists/${id}/symbols`,
    WATCHLIST_REMOVE_SYMBOL: (id: string, symbol: string) => 
      `/market/watchlists/${id}/symbols/${symbol}`,
    MARKET_SUMMARY: '/market/summary',
  },
};
```

## API Client

Next, let's create a reusable API client with error handling and authentication:

```typescript
// frontend/src/services/api.client.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { API_CONFIG } from '../config/api.config';

class ApiClient {
  private instance: AxiosInstance;
  private authToken: string | null = null;
  
  constructor() {
    this.instance = axios.create({
      baseURL: API_CONFIG.BASE_URL,
      timeout: API_CONFIG.TIMEOUT,
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    // Initialize auth token from localStorage
    this.authToken = localStorage.getItem('token');
    
    // Request interceptor for adding auth token
    this.instance.interceptors.request.use(
      (config) => {
        if (this.authToken) {
          config.headers.Authorization = `Bearer ${this.authToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );
    
    // Response interceptor for handling errors
    this.instance.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };
        
        // Handle 401 Unauthorized error (token expired)
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          
          try {
            // Try to refresh the token
            const refreshToken = localStorage.getItem('refreshToken');
            if (refreshToken) {
              const response = await this.instance.post('/auth/refresh-token', {
                refreshToken,
              });
              
              const { token } = response.data;
              localStorage.setItem('token', token);
              this.authToken = token;
              
              // Retry the original request with new token
              originalRequest.headers = {
                ...originalRequest.headers,
                Authorization: `Bearer ${token}`,
              };
              
              return this.instance(originalRequest);
            }
          } catch (refreshError) {
            // If refresh token fails, redirect to login
            localStorage.removeItem('token');
            localStorage.removeItem('refreshToken');
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }
        
        // Handle other errors
        return Promise.reject(error);
      }
    );
  }
  
  // Set auth token
  public setAuthToken(token: string): void {
    this.authToken = token;
    localStorage.setItem('token', token);
  }
  
  // Clear auth token
  public clearAuthToken(): void {
    this.authToken = null;
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
  }
  
  // Generic request method
  public async request<T>(config: AxiosRequestConfig): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.instance.request(config);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        // Handle specific error cases
        const status = error.response?.status;
        const errorData = error.response?.data;
        
        // Throw a more informative error
        throw {
          status,
          message: errorData?.message || error.message,
          data: errorData,
          originalError: error,
        };
      }
      throw error;
    }
  }
  
  // Convenience methods for common HTTP methods
  public async get<T>(url: string, params?: any): Promise<T> {
    return this.request<T>({ method: 'GET', url, params });
  }
  
  public async post<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'POST', url, data });
  }
  
  public async put<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'PUT', url, data });
  }
  
  public async patch<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'PATCH', url, data });
  }
  
  public async delete<T>(url: string): Promise<T> {
    return this.request<T>({ method: 'DELETE', url });
  }
}

// Create and export a singleton instance
export const apiClient = new ApiClient();
```

## Service Implementations

Now, let's implement the service modules that will use the API client to communicate with the backend:

### Authentication Service

```typescript
// frontend/src/services/auth.service.ts
import { apiClient } from './api.client';
import { ENDPOINTS } from '../config/api.config';

export interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  refreshToken: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
  };
}

export const authService = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>(ENDPOINTS.AUTH.LOGIN, data);
    
    // Store tokens
    apiClient.setAuthToken(response.token);
    localStorage.setItem('refreshToken', response.refreshToken);
    
    return response;
  },
  
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>(ENDPOINTS.AUTH.REGISTER, data);
    
    // Store tokens
    apiClient.setAuthToken(response.token);
    localStorage.setItem('refreshToken', response.refreshToken);
    
    return response;
  },
  
  forgotPassword: async (data: ForgotPasswordRequest): Promise<{ message: string }> => {
    return apiClient.post<{ message: string }>(ENDPOINTS.AUTH.FORGOT_PASSWORD, data);
  },
  
  resetPassword: async (data: ResetPasswordRequest): Promise<{ message: string }> => {
    return apiClient.post<{ message: string }>(ENDPOINTS.AUTH.RESET_PASSWORD, data);
  },
  
  logout: async (): Promise<void> => {
    try {
      await apiClient.post(ENDPOINTS.AUTH.LOGOUT);
    } finally {
      // Clear tokens regardless of API response
      apiClient.clearAuthToken();
    }
  },
  
  refreshToken: async (): Promise<AuthResponse> => {
    const refreshToken = localStorage.getItem('refreshToken');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }
    
    const response = await apiClient.post<AuthResponse>(
      ENDPOINTS.AUTH.REFRESH_TOKEN,
      { refreshToken }
    );
    
    // Update tokens
    apiClient.setAuthToken(response.token);
    localStorage.setItem('refreshToken', response.refreshToken);
    
    return response;
  },
};
```

### Order Service

```typescript
// frontend/src/services/order.service.ts
import { apiClient } from './api.client';
import { ENDPOINTS } from '../config/api.config';

export interface OrderFormData {
  symbol: string;
  side: 'BUY' | 'SELL';
  type: 'MARKET' | 'LIMIT' | 'STOP' | 'STOP_LIMIT';
  quantity: number;
  price?: number;
  stopPrice?: number;
  timeInForce?: 'DAY' | 'GTC' | 'IOC' | 'FOK';
  strategyId?: string;
}

export interface Order {
  id: string;
  symbol: string;
  side: 'BUY' | 'SELL';
  type: 'MARKET' | 'LIMIT' | 'STOP' | 'STOP_LIMIT';
  quantity: number;
  price?: number;
  stopPrice?: number;
  timeInForce: 'DAY' | 'GTC' | 'IOC' | 'FOK';
  status: 'OPEN' | 'FILLED' | 'PARTIALLY_FILLED' | 'CANCELLED' | 'REJECTED';
  filledQuantity: number;
  averagePrice?: number;
  strategyId?: string;
  createdAt: string;
  updatedAt: string;
  filledAt?: string;
}

export const orderService = {
  getOrders: async (): Promise<Order[]> => {
    return apiClient.get<Order[]>(ENDPOINTS.ORDER.LIST);
  },
  
  getOrderById: async (id: string): Promise<Order> => {
    return apiClient.get<Order>(ENDPOINTS.ORDER.DETAILS(id));
  },
  
  submitOrder: async (orderData: OrderFormData): Promise<Order> => {
    return apiClient.post<Order>(ENDPOINTS.ORDER.CREATE, orderData);
  },
  
  updateOrder: async (id: string, orderData: Partial<OrderFormData>): Promise<Order> => {
    return apiClient.put<Order>(ENDPOINTS.ORDER.UPDATE(id), orderData);
  },
  
  cancelOrder: async (id: string): Promise<Order> => {
    return apiClient.delete<Order>(ENDPOINTS.ORDER.CANCEL(id));
  },
};
```

### Portfolio Service

```typescript
// frontend/src/services/portfolio.service.ts
import { apiClient } from './api.client';
import { ENDPOINTS } from '../config/api.config';

export interface PortfolioSummary {
  totalValue: number;
  cashBalance: number;
  investedValue: number;
  dayChange: number;
  dayChangePercent: number;
  totalReturn: number;
  totalReturnPercent: number;
}

export interface Holding {
  symbol: string;
  name: string;
  quantity: number;
  averagePrice: number;
  currentPrice: number;
  marketValue: number;
  unrealizedPL: number;
  unrealizedPLPercent: number;
  weight: number;
}

export interface PerformanceData {
  daily: { date: string; value: number }[];
  weekly: { date: string; value: number }[];
  monthly: { date: string; value: number }[];
  yearly: { date: string; value: number }[];
}

export interface PortfolioSettings {
  currency: string;
  benchmark: string;
}

export interface RiskMetrics {
  sharpeRatio: number;
  sortinoRatio: number;
  maxDrawdown: number;
  volatility: number;
  beta: number;
  alpha: number;
}

export const portfolioService = {
  getSummary: async (): Promise<PortfolioSummary> => {
    return apiClient.get<PortfolioSummary>(ENDPOINTS.PORTFOLIO.SUMMARY);
  },
  
  getHoldings: async (): Promise<Holding[]> => {
    return apiClient.get<Holding[]>(ENDPOINTS.PORTFOLIO.HOLDINGS);
  },
  
  getPerformance: async (timeframe?: string): Promise<PerformanceData> => {
    return apiClient.get<PerformanceData>(ENDPOINTS.PORTFOLIO.PERFORMANCE, { timeframe });
  },
  
  updateSettings: async (settings: PortfolioSettings): Promise<PortfolioSettings> => {
    return apiClient.put<PortfolioSettings>(ENDPOINTS.PORTFOLIO.SUMMARY, settings);
  },
  
  getAnalytics: async (): Promise<{
    riskMetrics: RiskMetrics;
    correlationMatrix: Record<string, Record<string, number>>;
    attribution: { symbol: string; contribution: number; weight: number }[];
  }> => {
    return apiClient.get(ENDPOINTS.PORTFOLIO.ANALYTICS);
  },
};
```

### Strategy Service

```typescript
// frontend/src/services/strategy.service.ts
import { apiClient } from './api.client';
import { ENDPOINTS } from '../config/api.config';

export interface StrategyFormData {
  name: string;
  description: string;
  type: string;
  parameters: Record<string, any>;
}

export interface Strategy {
  id: string;
  name: string;
  description: string;
  status: 'ACTIVE' | 'INACTIVE';
  type: string;
  parameters: Record<string, any>;
  performance: {
    totalReturn: number;
    sharpeRatio: number;
    maxDrawdown: number;
    winRate: number;
  };
  createdAt: string;
  lastModified: string;
}

export const strategyService = {
  getStrategies: async (): Promise<Strategy[]> => {
    return apiClient.get<Strategy[]>(ENDPOINTS.STRATEGY.LIST);
  },
  
  getStrategyById: async (id: string): Promise<Strategy> => {
    return apiClient.get<Strategy>(ENDPOINTS.STRATEGY.DETAILS(id));
  },
  
  createStrategy: async (strategyData: StrategyFormData): Promise<Strategy> => {
    return apiClient.post<Strategy>(ENDPOINTS.STRATEGY.CREATE, strategyData);
  },
  
  updateStrategy: async (id: string, strategyData: Partial<StrategyFormData>): Promise<Strategy> => {
    return apiClient.put<Strategy>(ENDPOINTS.STRATEGY.UPDATE(id), strategyData);
  },
  
  deleteStrategy: async (id: string): Promise<void> => {
    return apiClient.delete(ENDPOINTS.STRATEGY.DELETE(id));
  },
  
  activateStrategy: async (id: string): Promise<Strategy> => {
    return apiClient.post<Strategy>(ENDPOINTS.STRATEGY.ACTIVATE(id));
  },
  
  deactivateStrategy: async (id: string): Promise<Strategy> => {
    return apiClient.post<Strategy>(ENDPOINTS.STRATEGY.DEACTIVATE(id));
  },
};
```

### Market Service

```typescript
// frontend/src/services/market.service.ts
import { apiClient } from './api.client';
import { ENDPOINTS } from '../config/api.config';

export interface Quote {
  symbol: string;
  name: string;
  price: number;
  change: number;
  changePercent: number;
  volume: number;
  marketCap: number;
  bid?: number;
  ask?: number;
  high?: number;
  low?: number;
  open?: number;
  close?: number;
}

export interface Watchlist {
  id: string;
  name: string;
  symbols: string[];
}

export interface MarketSummary {
  indices: {
    name: string;
    value: number;
    change: number;
    changePercent: number;
  }[];
  sectors: {
    name: string;
    performance: number;
  }[];
}

export const marketService = {
  getQuotes: async (symbols: string[]): Promise<Record<string, Quote>> => {
    return apiClient.get<Record<string, Quote>>(ENDPOINTS.MARKET.QUOTES, { symbols: symbols.join(',') });
  },
  
  getQuote: async (symbol: string): Promise<Quote> => {
    return apiClient.get<Quote>(ENDPOINTS.MARKET.QUOTE(symbol));
  },
  
  getWatchlists: async (): Promise<Watchlist[]> => {
    return apiClient.get<Watchlist[]>(ENDPOINTS.MARKET.WATCHLISTS);
  },
  
  getWatchlist: async (id: string): Promise<Watchlist> => {
    return apiClient.get<Watchlist>(ENDPOINTS.MARKET.WATCHLIST(id));
  },
  
  createWatchlist: async (name: string, symbols: string[] = []): Promise<Watchlist> => {
    return apiClient.post<Watchlist>(ENDPOINTS.MARKET.WATCHLISTS, { name, symbols });
  },
  
  updateWatchlist: async (id: string, name: string): Promise<Watchlist> => {
    return apiClient.put<Watchlist>(ENDPOINTS.MARKET.WATCHLIST(id), { name });
  },
  
  deleteWatchlist: async (id: string): Promise<void> => {
    return apiClient.delete(ENDPOINTS.MARKET.WATCHLIST(id));
  },
  
  addSymbolToWatchlist: async (id: string, symbol: string): Promise<Watchlist> => {
    return apiClient.post<Watchlist>(ENDPOINTS.MARKET.WATCHLIST_ADD_SYMBOL(id), { symbol });
  },
  
  removeSymbolFromWatchlist: async (id: string, symbol: string): Promise<Watchlist> => {
    return apiClient.delete<Watchlist>(ENDPOINTS.MARKET.WATCHLIST_REMOVE_SYMBOL(id, symbol));
  },
  
  getMarketSummary: async (): Promise<MarketSummary> => {
    return apiClient.get<MarketSummary>(ENDPOINTS.MARKET.MARKET_SUMMARY);
  },
};
```

## WebSocket Integration

For real-time updates, we'll implement a WebSocket service:

```typescript
// frontend/src/services/websocket.service.ts
import { API_CONFIG } from '../config/api.config';

export type MessageHandler = (data: any) => void;

export class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 1000; // Start with 1 second
  private messageHandlers: Record<string, MessageHandler[]> = {};
  private isConnecting = false;
  
  // Connect to WebSocket server
  public connect(): Promise<void> {
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      return Promise.resolve();
    }
    
    if (this.isConnecting) {
      return new Promise((resolve) => {
        const checkInterval = setInterval(() => {
          if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            clearInterval(checkInterval);
            resolve();
          }
        }, 100);
      });
    }
    
    this.isConnecting = true;
    
    return new Promise((resolve, reject) => {
      try {
        const token = localStorage.getItem('token');
        const url = `${API_CONFIG.WS_URL}?token=${token}`;
        
        this.socket = new WebSocket(url);
        
        this.socket.onopen = () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          this.isConnecting = false;
          resolve();
        };
        
        this.socket.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.handleMessage(data);
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
          }
        };
        
        this.socket.onclose = () => {
          console.log('WebSocket disconnected');
          this.socket = null;
          this.isConnecting = false;
          
          // Attempt to reconnect
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const timeout = this.reconnectTimeout * Math.pow(2, this.reconnectAttempts - 1);
            
            setTimeout(() => {
              this.connect().catch(console.error);
            }, timeout);
          }
        };
        
        this.socket.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnecting = false;
          reject(error);
        };
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }
  
  // Disconnect from WebSocket server
  public disconnect(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
  
  // Send message to WebSocket server
  public send(type: string, payload: any): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      this.connect().then(() => {
        this.sendMessage(type, payload);
      }).catch(console.error);
      return;
    }
    
    this.sendMessage(type, payload);
  }
  
  private sendMessage(type: string, payload: any): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({ type, payload });
      this.socket.send(message);
    }
  }
  
  // Register message handler
  public on(type: string, handler: MessageHandler): void {
    if (!this.messageHandlers[type]) {
      this.messageHandlers[type] = [];
    }
    
    this.messageHandlers[type].push(handler);
  }
  
  // Remove message handler
  public off(type: string, handler: MessageHandler): void {
    if (!this.messageHandlers[type]) {
      return;
    }
    
    this.messageHandlers[type] = this.messageHandlers[type].filter(h => h !== handler);
  }
  
  // Handle incoming message
  private handleMessage(data: any): void {
    const { type, payload } = data;
    
    if (!type) {
      console.error('Received WebSocket message without type:', data);
      return;
    }
    
    const handlers = this.messageHandlers[type] || [];
    
    handlers.forEach(handler => {
      try {
        handler(payload);
      } catch (error) {
        console.error(`Error in WebSocket handler for type ${type}:`, error);
      }
    });
  }
}

// Create and export a singleton instance
export const websocketService = new WebSocketService();
```

## Redux Integration

Now, let's update the Redux slices to use the API services:

### Auth Slice

```typescript
// frontend/src/store/slices/authSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { authService, LoginRequest, RegisterRequest, AuthResponse } from '../../services/auth.service';

interface AuthState {
  user: AuthResponse['user'] | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  user: null,
  isAuthenticated: !!localStorage.getItem('token'),
  loading: false,
  error: null,
};

export const login = createAsyncThunk(
  'auth/login',
  async (credentials: LoginRequest, { rejectWithValue }) => {
    try {
      return await authService.login(credentials);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Login failed');
    }
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async (userData: RegisterRequest, { rejectWithValue }) => {
    try {
      return await authService.register(userData);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Registration failed');
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      await authService.logout();
      return null;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Logout failed');
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    // Login
    builder
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
    
    // Register
    builder
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.user = action.payload.user;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
    
    // Logout
    builder
      .addCase(logout.pending, (state) => {
        state.loading = true;
      })
      .addCase(logout.fulfilled, (state) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.user = null;
      })
      .addCase(logout.rejected, (state) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.user = null;
      });
  },
});

export const { clearError } = authSlice.actions;
export default authSlice.reducer;
```

### Order Slice

```typescript
// frontend/src/store/slices/orderSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { orderService, OrderFormData, Order } from '../../services/order.service';
import { websocketService } from '../../services/websocket.service';

interface OrderState {
  orders: Order[];
  loading: boolean;
  error: string | null;
}

const initialState: OrderState = {
  orders: [],
  loading: false,
  error: null,
};

export const fetchOrders = createAsyncThunk(
  'order/fetchOrders',
  async (_, { rejectWithValue }) => {
    try {
      return await orderService.getOrders();
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch orders');
    }
  }
);

export const submitOrder = createAsyncThunk(
  'order/submitOrder',
  async (orderData: OrderFormData, { rejectWithValue }) => {
    try {
      return await orderService.submitOrder(orderData);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to submit order');
    }
  }
);

export const cancelOrder = createAsyncThunk(
  'order/cancelOrder',
  async (orderId: string, { rejectWithValue }) => {
    try {
      return await orderService.cancelOrder(orderId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to cancel order');
    }
  }
);

export const modifyOrder = createAsyncThunk(
  'order/modifyOrder',
  async ({ id, ...orderData }: { id: string } & Partial<OrderFormData>, { rejectWithValue }) => {
    try {
      return await orderService.updateOrder(id, orderData);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to modify order');
    }
  }
);

const orderSlice = createSlice({
  name: 'order',
  initialState,
  reducers: {
    updateOrders: (state, action) => {
      // Handle WebSocket updates
      const updatedOrder = action.payload;
      const index = state.orders.findIndex(order => order.id === updatedOrder.id);
      
      if (index !== -1) {
        // Update existing order
        state.orders[index] = updatedOrder;
      } else {
        // Add new order
        state.orders.push(updatedOrder);
      }
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    // Fetch orders
    builder
      .addCase(fetchOrders.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchOrders.fulfilled, (state, action) => {
        state.loading = false;
        state.orders = action.payload;
      })
      .addCase(fetchOrders.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
    
    // Submit order
    builder
      .addCase(submitOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(submitOrder.fulfilled, (state, action) => {
        state.loading = false;
        state.orders.push(action.payload);
      })
      .addCase(submitOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
    
    // Cancel order
    builder
      .addCase(cancelOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(cancelOrder.fulfilled, (state, action) => {
        state.loading = false;
        const index = state.orders.findIndex(order => order.id === action.payload.id);
        if (index !== -1) {
          state.orders[index] = action.payload;
        }
      })
      .addCase(cancelOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
    
    // Modify order
    builder
      .addCase(modifyOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(modifyOrder.fulfilled, (state, action) => {
        state.loading = false;
        const index = state.orders.findIndex(order => order.id === action.payload.id);
        if (index !== -1) {
          state.orders[index] = action.payload;
        }
      })
      .addCase(modifyOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

// Set up WebSocket handler for order updates
websocketService.on('ORDER_UPDATE', (payload) => {
  orderSlice.actions.updateOrders(payload);
});

export const { updateOrders, clearError } = orderSlice.actions;
export default orderSlice.reducer;
```

## WebSocket Provider

To make WebSocket integration available throughout the application, let's create a WebSocket provider:

```tsx
// frontend/src/providers/WebSocketProvider.tsx
import React, { createContext, useContext, useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import { websocketService } from '../services/websocket.service';

interface WebSocketContextType {
  connected: boolean;
  connect: () => Promise<void>;
  disconnect: () => void;
  send: (type: string, payload: any) => void;
}

const WebSocketContext = createContext<WebSocketContextType>({
  connected: false,
  connect: () => Promise.resolve(),
  disconnect: () => {},
  send: () => {},
});

export const useWebSocket = () => useContext(WebSocketContext);

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [connected, setConnected] = useState(false);
  const { isAuthenticated } = useSelector((state: RootState) => state.auth);
  
  // Connect to WebSocket when authenticated
  useEffect(() => {
    if (isAuthenticated) {
      const connectWebSocket = async () => {
        try {
          await websocketService.connect();
          setConnected(true);
        } catch (error) {
          console.error('Failed to connect to WebSocket:', error);
          setConnected(false);
        }
      };
      
      connectWebSocket();
      
      // Set up connection status handler
      const checkConnection = setInterval(() => {
        const socket = (websocketService as any).socket;
        if (socket) {
          setConnected(socket.readyState === WebSocket.OPEN);
        } else {
          setConnected(false);
        }
      }, 5000);
      
      return () => {
        clearInterval(checkConnection);
        websocketService.disconnect();
      };
    } else {
      websocketService.disconnect();
      setConnected(false);
    }
  }, [isAuthenticated]);
  
  const connect = async () => {
    try {
      await websocketService.connect();
      setConnected(true);
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      setConnected(false);
      throw error;
    }
  };
  
  const disconnect = () => {
    websocketService.disconnect();
    setConnected(false);
  };
  
  const send = (type: string, payload: any) => {
    websocketService.send(type, payload);
  };
  
  return (
    <WebSocketContext.Provider value={{ connected, connect, disconnect, send }}>
      {children}
    </WebSocketContext.Provider>
  );
};
```

## App Integration

Finally, let's update the main App component to use our providers:

```tsx
// frontend/src/App.tsx
import React from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { store } from './store';
import { WebSocketProvider } from './providers/WebSocketProvider';
import AppRoutes from './routes';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <BrowserRouter>
          <WebSocketProvider>
            <AppRoutes />
          </WebSocketProvider>
        </BrowserRouter>
      </ThemeProvider>
    </Provider>
  );
};

export default App;
```

## Backend CORS Configuration

To allow the frontend to communicate with the backend, we need to configure CORS on the backend:

```go
// backend/internal/api/server.go
package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	
	router.Use(cors.New(config))

	// Setup routes
	setupRoutes(router)

	return router
}
```

## WebSocket Server Implementation

Let's implement the WebSocket server on the backend:

```go
// backend/internal/websocket/handler.go
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/trading-platform/backend/internal/auth"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}
	
	// Map of client connections
	clients = make(map[*websocket.Conn]string) // conn -> userID
	
	// Mutex for thread-safe operations on clients map
	clientsMutex = &sync.Mutex{}
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}
	
	// Validate token
	claims, err := auth.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	
	// Register client
	clientsMutex.Lock()
	clients[conn] = claims.UserID
	clientsMutex.Unlock()
	
	// Send welcome message
	welcomeMsg := Message{
		Type: "CONNECTED",
		Payload: map[string]interface{}{
			"message": "Connected to WebSocket server",
			"userID":  claims.UserID,
		},
	}
	
	err = conn.WriteJSON(welcomeMsg)
	if err != nil {
		log.Printf("Failed to send welcome message: %v", err)
		conn.Close()
		return
	}
	
	// Handle incoming messages
	go handleMessages(conn, claims.UserID)
}

// handleMessages processes incoming messages from a client
func handleMessages(conn *websocket.Conn, userID string) {
	defer func() {
		// Unregister client on disconnect
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		
		conn.Close()
	}()
	
	for {
		// Read message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		// Parse message
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}
		
		// Handle message based on type
		switch message.Type {
		case "PING":
			// Respond with pong
			pongMsg := Message{
				Type:    "PONG",
				Payload: nil,
			}
			
			if err := conn.WriteJSON(pongMsg); err != nil {
				log.Printf("Failed to send pong: %v", err)
			}
		
		// Add more message type handlers as needed
		default:
			log.Printf("Unknown message type: %s", message.Type)
		}
	}
}

// BroadcastToAll sends a message to all connected clients
func BroadcastToAll(messageType string, payload interface{}) {
	message := Message{
		Type:    messageType,
		Payload: payload,
	}
	
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	
	for conn := range clients {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Failed to broadcast to client: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

// BroadcastToUser sends a message to a specific user
func BroadcastToUser(userID string, messageType string, payload interface{}) {
	message := Message{
		Type:    messageType,
		Payload: payload,
	}
	
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	
	for conn, id := range clients {
		if id == userID {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Failed to broadcast to user %s: %v", userID, err)
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}
```

## Integration Testing

Let's create a simple integration test to verify the backend-frontend communication:

```typescript
// frontend/src/tests/integration/api.integration.test.ts
import { apiClient } from '../../services/api.client';
import { authService } from '../../services/auth.service';
import { orderService } from '../../services/order.service';
import { portfolioService } from '../../services/portfolio.service';
import { strategyService } from '../../services/strategy.service';
import { marketService } from '../../services/market.service';

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString();
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

// Integration tests
describe('API Integration Tests', () => {
  // Test credentials
  const testUser = {
    email: 'test@example.com',
    password: 'Password123!',
    firstName: 'Test',
    lastName: 'User',
  };
  
  // Clear localStorage before each test
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });
  
  // Mock fetch
  global.fetch = jest.fn();
  
  test('Authentication flow', async () => {
    // Mock successful login response
    (fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          token: 'test-token',
          refreshToken: 'test-refresh-token',
          user: {
            id: '123',
            email: testUser.email,
            firstName: testUser.firstName,
            lastName: testUser.lastName,
          },
        }),
      })
    );
    
    // Login
    const loginResponse = await authService.login({
      email: testUser.email,
      password: testUser.password,
    });
    
    // Check response
    expect(loginResponse).toEqual({
      token: 'test-token',
      refreshToken: 'test-refresh-token',
      user: {
        id: '123',
        email: testUser.email,
        firstName: testUser.firstName,
        lastName: testUser.lastName,
      },
    });
    
    // Check localStorage
    expect(localStorage.getItem('token')).toBe('test-token');
    expect(localStorage.getItem('refreshToken')).toBe('test-refresh-token');
    
    // Mock successful logout
    (fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({}),
      })
    );
    
    // Logout
    await authService.logout();
    
    // Check localStorage after logout
    expect(localStorage.getItem('token')).toBeNull();
    expect(localStorage.getItem('refreshToken')).toBeNull();
  });
  
  test('Order submission and management', async () => {
    // Set token for authenticated requests
    localStorage.setItem('token', 'test-token');
    
    // Mock order submission response
    (fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          id: 'order-123',
          symbol: 'AAPL',
          side: 'BUY',
          type: 'MARKET',
          quantity: 100,
          status: 'OPEN',
          filledQuantity: 0,
          createdAt: '2023-04-01T10:30:00Z',
          updatedAt: '2023-04-01T10:30:00Z',
        }),
      })
    );
    
    // Submit order
    const orderResponse = await orderService.submitOrder({
      symbol: 'AAPL',
      side: 'BUY',
      type: 'MARKET',
      quantity: 100,
    });
    
    // Check response
    expect(orderResponse).toEqual({
      id: 'order-123',
      symbol: 'AAPL',
      side: 'BUY',
      type: 'MARKET',
      quantity: 100,
      status: 'OPEN',
      filledQuantity: 0,
      createdAt: '2023-04-01T10:30:00Z',
      updatedAt: '2023-04-01T10:30:00Z',
    });
    
    // Mock get orders response
    (fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve([
          {
            id: 'order-123',
            symbol: 'AAPL',
            side: 'BUY',
            type: 'MARKET',
            quantity: 100,
            status: 'OPEN',
            filledQuantity: 0,
            createdAt: '2023-04-01T10:30:00Z',
            updatedAt: '2023-04-01T10:30:00Z',
          },
        ]),
      })
    );
    
    // Get orders
    const orders = await orderService.getOrders();
    
    // Check response
    expect(orders).toHaveLength(1);
    expect(orders[0].id).toBe('order-123');
    
    // Mock cancel order response
    (fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          id: 'order-123',
          symbol: 'AAPL',
          side: 'BUY',
          type: 'MARKET',
          quantity: 100,
          status: 'CANCELLED',
          filledQuantity: 0,
          createdAt: '2023-04-01T10:30:00Z',
          updatedAt: '2023-04-01T10:35:00Z',
        }),
      })
    );
    
    // Cancel order
    const cancelledOrder = await orderService.cancelOrder('order-123');
    
    // Check response
    expect(cancelledOrder.status).toBe('CANCELLED');
  });
});
```

## Conclusion

This integration connects the frontend React components with the backend Go services through RESTful APIs and WebSocket for real-time updates. The key components of the integration include:

1. **API Configuration**: Centralized configuration of API endpoints
2. **API Client**: Reusable client with error handling and authentication
3. **Service Modules**: Typed service modules for each feature area
4. **WebSocket Integration**: Real-time updates with reconnection logic
5. **Redux Integration**: Updated Redux slices to use the API services
6. **CORS Configuration**: Backend configuration to allow frontend requests
7. **Integration Testing**: Tests to verify the backend-frontend communication

This integration provides a solid foundation for the trading platform, enabling seamless communication between the frontend and backend components.
