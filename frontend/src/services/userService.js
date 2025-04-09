import apiService from './apiService';

const USER_API_PATH = '/users';

export const userService = {
  // Get current user profile
  getCurrentUser: async () => {
    return await apiService.get(`${USER_API_PATH}/me`);
  },
  
  // Update user profile
  updateProfile: async (profileData) => {
    return await apiService.put(`${USER_API_PATH}/profile`, profileData);
  },
  
  // Update user password
  updatePassword: async (passwordData) => {
    return await apiService.put(`${USER_API_PATH}/password`, passwordData);
  },
  
  // Get user broker settings
  getBrokerSettings: async () => {
    return await apiService.get(`${USER_API_PATH}/broker-settings`);
  },
  
  // Update user broker settings
  updateBrokerSettings: async (brokerSettings) => {
    return await apiService.put(`${USER_API_PATH}/broker-settings`, brokerSettings);
  },
  
  // Get user notification settings
  getNotificationSettings: async () => {
    return await apiService.get(`${USER_API_PATH}/notification-settings`);
  },
  
  // Update user notification settings
  updateNotificationSettings: async (notificationSettings) => {
    return await apiService.put(`${USER_API_PATH}/notification-settings`, notificationSettings);
  },
  
  // Get user theme settings
  getThemeSettings: async () => {
    return await apiService.get(`${USER_API_PATH}/theme-settings`);
  },
  
  // Update user theme settings
  updateThemeSettings: async (themeSettings) => {
    return await apiService.put(`${USER_API_PATH}/theme-settings`, themeSettings);
  },
  
  // Get user activity logs
  getActivityLogs: async (filters = {}) => {
    return await apiService.get(`${USER_API_PATH}/activity-logs`, filters);
  },
  
  // Test broker connection
  testBrokerConnection: async (brokerData) => {
    return await apiService.post(`${USER_API_PATH}/test-broker-connection`, brokerData);
  }
};

export default userService;
