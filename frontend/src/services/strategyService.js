import apiService from './apiService';

const STRATEGY_API_PATH = '/strategies';

export const strategyService = {
  // Get all strategies
  getStrategies: async () => {
    return await apiService.get(STRATEGY_API_PATH);
  },
  
  // Get a specific strategy by ID
  getStrategyById: async (strategyId) => {
    return await apiService.get(`${STRATEGY_API_PATH}/${strategyId}`);
  },
  
  // Create a new strategy
  createStrategy: async (strategyData) => {
    return await apiService.post(STRATEGY_API_PATH, strategyData);
  },
  
  // Update an existing strategy
  updateStrategy: async (strategyId, strategyData) => {
    return await apiService.put(`${STRATEGY_API_PATH}/${strategyId}`, strategyData);
  },
  
  // Delete a strategy
  deleteStrategy: async (strategyId) => {
    return await apiService.delete(`${STRATEGY_API_PATH}/${strategyId}`);
  },
  
  // Get all portfolios for a strategy
  getPortfolios: async (strategyId) => {
    return await apiService.get(`${STRATEGY_API_PATH}/${strategyId}/portfolios`);
  },
  
  // Get a specific portfolio by ID
  getPortfolioById: async (strategyId, portfolioId) => {
    return await apiService.get(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}`);
  },
  
  // Create a new portfolio for a strategy
  createPortfolio: async (strategyId, portfolioData) => {
    return await apiService.post(`${STRATEGY_API_PATH}/${strategyId}/portfolios`, portfolioData);
  },
  
  // Update an existing portfolio
  updatePortfolio: async (strategyId, portfolioId, portfolioData) => {
    return await apiService.put(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}`, portfolioData);
  },
  
  // Delete a portfolio
  deletePortfolio: async (strategyId, portfolioId) => {
    return await apiService.delete(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}`);
  },
  
  // Get portfolio configuration
  getPortfolioConfig: async (strategyId, portfolioId) => {
    return await apiService.get(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}/config`);
  },
  
  // Update portfolio configuration
  updatePortfolioConfig: async (strategyId, portfolioId, configData) => {
    return await apiService.put(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}/config`, configData);
  },
  
  // Start a portfolio
  startPortfolio: async (strategyId, portfolioId) => {
    return await apiService.post(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}/start`);
  },
  
  // Stop a portfolio
  stopPortfolio: async (strategyId, portfolioId) => {
    return await apiService.post(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}/stop`);
  },
  
  // Get portfolio performance
  getPortfolioPerformance: async (strategyId, portfolioId, timeframe = 'day') => {
    return await apiService.get(`${STRATEGY_API_PATH}/${strategyId}/portfolios/${portfolioId}/performance`, { timeframe });
  }
};

export default strategyService;
