import apiService from './apiService';

const LOG_API_PATH = '/logs';

export const logService = {
  // Get all logs with optional filters
  getLogs: async (filters = {}) => {
    return await apiService.get(LOG_API_PATH, filters);
  },
  
  // Get logs by level
  getLogsByLevel: async (level, filters = {}) => {
    return await apiService.get(`${LOG_API_PATH}/level/${level}`, filters);
  },
  
  // Get logs by source
  getLogsBySource: async (source, filters = {}) => {
    return await apiService.get(`${LOG_API_PATH}/source/${source}`, filters);
  },
  
  // Get logs by date range
  getLogsByDateRange: async (startDate, endDate, filters = {}) => {
    return await apiService.get(`${LOG_API_PATH}/date-range`, {
      ...filters,
      startDate,
      endDate
    });
  },
  
  // Add a new log entry
  addLog: async (logData) => {
    return await apiService.post(LOG_API_PATH, logData);
  },
  
  // Clear all logs
  clearLogs: async () => {
    return await apiService.delete(`${LOG_API_PATH}/clear`);
  },
  
  // Export logs to file
  exportLogs: async (format = 'json', filters = {}) => {
    return await apiService.get(`${LOG_API_PATH}/export`, {
      ...filters,
      format
    });
  }
};

export default logService;
