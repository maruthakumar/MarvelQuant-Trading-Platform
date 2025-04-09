import apiService from './apiService';

const ORDER_API_PATH = '/orders';

export const orderService = {
  // Get all orders with optional filters
  getOrders: async (filters = {}) => {
    return await apiService.get(ORDER_API_PATH, filters);
  },
  
  // Get a specific order by ID
  getOrderById: async (orderId) => {
    return await apiService.get(`${ORDER_API_PATH}/${orderId}`);
  },
  
  // Create a new order
  createOrder: async (orderData) => {
    return await apiService.post(ORDER_API_PATH, orderData);
  },
  
  // Update an existing order
  updateOrder: async (orderId, orderData) => {
    return await apiService.put(`${ORDER_API_PATH}/${orderId}`, orderData);
  },
  
  // Cancel an order
  cancelOrder: async (orderId) => {
    return await apiService.post(`${ORDER_API_PATH}/${orderId}/cancel`);
  },
  
  // Get order history
  getOrderHistory: async (filters = {}) => {
    return await apiService.get(`${ORDER_API_PATH}/history`, filters);
  },
  
  // Get order statistics
  getOrderStats: async (timeframe = 'day') => {
    return await apiService.get(`${ORDER_API_PATH}/stats`, { timeframe });
  }
};

export default orderService;
