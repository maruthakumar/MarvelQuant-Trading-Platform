import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import App from './App';
import LoginPage from './components/auth/LoginPage';
import Dashboard from './components/dashboard/Dashboard';
import OrderBook from './components/orderbook/OrderBook';
import PositionsPanel from './components/positions/PositionsPanel';
import UserSettings from './components/user/UserSettings';
import StrategiesComponent from './components/strategies/StrategiesComponent';
import MultiLegComponent from './components/multileg/MultiLegComponent';
import LogsPanel from './components/logs/LogsPanel';

const ProtectedRoute = ({ children }) => {
  const isAuthenticated = localStorage.getItem('isAuthenticated') === 'true';
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

const AppRouter = () => {
  // Force check authentication on initial load
  useEffect(() => {
    const isAuthenticated = localStorage.getItem('isAuthenticated') === 'true';
    const path = window.location.pathname;
    
    // If not authenticated and not already on login page, redirect to login
    if (!isAuthenticated && path !== '/login') {
      window.location.href = '/login';
    }
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={
          <ProtectedRoute>
            <App />
          </ProtectedRoute>
        }>
          <Route index element={<Dashboard />} />
          <Route path="orderbook" element={<OrderBook />} />
          <Route path="positions" element={<PositionsPanel />} />
          <Route path="user-settings" element={<UserSettings />} />
          <Route path="strategies" element={<StrategiesComponent />} />
          <Route path="multi-leg" element={<MultiLegComponent />} />
          <Route path="logs" element={<LogsPanel />} />
        </Route>
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default AppRouter;
