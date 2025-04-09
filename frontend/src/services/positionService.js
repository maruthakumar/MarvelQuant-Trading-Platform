import apiService from './apiService';

const POSITION_API_PATH = '/positions';

export const positionService = {
  // Get all positions with optional filters
  getPositions: async (filters = {}) => {
    return await apiService.get(POSITION_API_PATH, filters);
  },
  
  // Get a specific position by ID
  getPositionById: async (positionId) => {
    return await apiService.get(`${POSITION_API_PATH}/${positionId}`);
  },
  
  // Close a position
  closePosition: async (positionId, closeData = {}) => {
    return await apiService.post(`${POSITION_API_PATH}/${positionId}/close`, closeData);
  },
  
  // Get position history
  getPositionHistory: async (filters = {}) => {
    return await apiService.get(`${POSITION_API_PATH}/history`, filters);
  },
  
  // Get position statistics
  getPositionStats: async (timeframe = 'day') => {
    return await apiService.get(`${POSITION_API_PATH}/stats`, { timeframe });
  },
  
  // Get position P&L
  getPositionPnL: async (positionId) => {
    return await apiService.get(`${POSITION_API_PATH}/${positionId}/pnl`);
  },
  
  // Update position stop loss or take profit
  updatePositionLimits: async (positionId, limitData) => {
    return await apiService.put(`${POSITION_API_PATH}/${positionId}/limits`, limitData);
  }
};

export default positionService;
