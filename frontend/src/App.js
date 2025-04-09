import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './components/auth/LoginPage';
import Dashboard from './components/dashboard/Dashboard';
import OrderBook from './components/orderbook/OrderBook';
import Positions from './components/positions/Positions';
import UserSettings from './components/settings/UserSettings';
import Strategies from './components/strategies/Strategies';
import MultiLegComponent from './components/multileg/MultiLegComponent';
import NewPortfolioPage from './components/multileg/NewPortfolioPage';
import ProtectedRoute from './components/auth/ProtectedRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="/orderbook" element={<ProtectedRoute><OrderBook /></ProtectedRoute>} />
        <Route path="/positions" element={<ProtectedRoute><Positions /></ProtectedRoute>} />
        <Route path="/user-settings" element={<ProtectedRoute><UserSettings /></ProtectedRoute>} />
        <Route path="/strategies" element={<ProtectedRoute><Strategies /></ProtectedRoute>} />
        <Route path="/multi-leg" element={<ProtectedRoute><MultiLegComponent /></ProtectedRoute>} />
        <Route path="/multi-leg/new-portfolio" element={<ProtectedRoute><NewPortfolioPage /></ProtectedRoute>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
