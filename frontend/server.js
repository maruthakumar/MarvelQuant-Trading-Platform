// Server configuration for MarvelQuant Trading Platform
const express = require('express');
const path = require('path');
const cors = require('cors');
const compression = require('compression');
const helmet = require('helmet');
const db = require('./src/db/connection');

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(compression()); // Compress responses
app.use(helmet({ contentSecurityPolicy: false })); // Security headers
app.use(cors()); // Enable CORS
app.use(express.json()); // Parse JSON bodies

// Serve static files from the React app
app.use(express.static(path.join(__dirname, 'build')));

// API Routes
app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', timestamp: new Date() });
});

// Database connection test endpoint
app.get('/api/db-test', async (req, res) => {
  try {
    // Test database connection
    const result = await db.raw('SELECT 1+1 as result');
    res.json({ 
      status: 'ok', 
      message: 'Database connection successful',
      result: result[0][0].result
    });
  } catch (error) {
    console.error('Database connection error:', error);
    res.status(500).json({ 
      status: 'error', 
      message: 'Database connection failed',
      error: error.message
    });
  }
});

// Catch-all handler to serve React app
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, 'build', 'index.html'));
});

// Start server
const server = app.listen(PORT, '0.0.0.0', () => {
  console.log(`MarvelQuant Trading Platform server running on port ${PORT} in ${process.env.NODE_ENV || 'production'} mode`);
});

module.exports = server;
