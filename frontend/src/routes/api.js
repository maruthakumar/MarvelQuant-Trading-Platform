// API routes for MarvelQuant Trading Platform
const express = require('express');
const router = express.Router();
const db = require('../db/connection');

// Get all trading strategies
router.get('/strategies', async (req, res) => {
  try {
    const strategies = await db('strategies').select('*');
    res.json(strategies);
  } catch (error) {
    console.error('Error fetching strategies:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get all positions
router.get('/positions', async (req, res) => {
  try {
    const positions = await db('positions').select('*');
    res.json(positions);
  } catch (error) {
    console.error('Error fetching positions:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get order book entries
router.get('/orderbook', async (req, res) => {
  try {
    const orders = await db('orders').select('*').orderBy('timestamp', 'desc');
    res.json(orders);
  } catch (error) {
    console.error('Error fetching orders:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get user settings
router.get('/user-settings/:userId', async (req, res) => {
  try {
    const { userId } = req.params;
    const settings = await db('user_settings').where({ user_id: userId }).first();
    res.json(settings || {});
  } catch (error) {
    console.error('Error fetching user settings:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get multi-leg configurations
router.get('/multi-leg', async (req, res) => {
  try {
    const configs = await db('multi_leg_configs').select('*');
    res.json(configs);
  } catch (error) {
    console.error('Error fetching multi-leg configurations:', error);
    res.status(500).json({ error: error.message });
  }
});

// Get system logs
router.get('/logs', async (req, res) => {
  try {
    const logs = await db('system_logs')
      .select('*')
      .orderBy('timestamp', 'desc')
      .limit(100);
    res.json(logs);
  } catch (error) {
    console.error('Error fetching system logs:', error);
    res.status(500).json({ error: error.message });
  }
});

module.exports = router;
