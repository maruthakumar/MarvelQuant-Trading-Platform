import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  Divider,
  Button,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Switch,
  FormControlLabel,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RefreshIcon from '@mui/icons-material/Refresh';
import FilterListIcon from '@mui/icons-material/FilterList';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import LogsPanel from '../logs/LogsPanel';

const OrderBook = () => {
  const [orders, setOrders] = useState([
    { 
      id: 'order1', 
      symbol: 'NIFTY', 
      type: 'LIMIT', 
      side: 'BUY',
      quantity: 75,
      price: 22500,
      status: 'OPEN',
      time: '09:31:16',
      strategy: 'MARU'
    },
    { 
      id: 'order2', 
      symbol: 'BANKNIFTY', 
      type: 'MARKET', 
      side: 'SELL',
      quantity: 25,
      price: 48750,
      status: 'EXECUTED',
      time: '09:25:45',
      strategy: 'LUX24VR'
    },
    { 
      id: 'order3', 
      symbol: 'NIFTY', 
      type: 'SL', 
      side: 'SELL',
      quantity: 50,
      price: 22450,
      status: 'CANCELLED',
      time: '09:18:30',
      strategy: 'MARU'
    },
    { 
      id: 'order4', 
      symbol: 'RELIANCE', 
      type: 'LIMIT', 
      side: 'BUY',
      quantity: 100,
      price: 2750.50,
      status: 'REJECTED',
      time: '09:17:22',
      strategy: 'NFTTR'
    },
    { 
      id: 'order5', 
      symbol: 'BANKNIFTY', 
      type: 'MARKET', 
      side: 'BUY',
      quantity: 25,
      price: 48500,
      status: 'EXECUTED',
      time: '09:15:45',
      strategy: 'LUX24VR'
    }
  ]);
  
  const [activeTab, setActiveTab] = useState(0);
  const [openDialog, setOpenDialog] = useState(false);
  const [newOrder, setNewOrder] = useState({
    symbol: 'NIFTY',
    type: 'LIMIT',
    side: 'BUY',
    quantity: 75,
    price: 22500,
    strategy: 'MARU'
  });
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [filterStatus, setFilterStatus] = useState('ALL');
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handleAddOrder = () => {
    setOpenDialog(true);
    setNewOrder({
      symbol: 'NIFTY',
      type: 'LIMIT',
      side: 'BUY',
      quantity: 75,
      price: 22500,
      strategy: 'MARU'
    });
  };
  
  const handleCloseDialog = () => {
    setOpenDialog(false);
  };
  
  const handlePlaceOrder = () => {
    const newOrderObj = {
      id: `order${Date.now()}`,
      ...newOrder,
      status: 'OPEN',
      time: new Date().toLocaleTimeString('en-US', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' })
    };
    
    setOrders([newOrderObj, ...orders]);
    setOpenDialog(false);
  };
  
  const handleCancelOrder = (orderId) => {
    if (window.confirm('Are you sure you want to cancel this order?')) {
      const updatedOrders = orders.map(order => 
        order.id === orderId 
          ? { ...order, status: 'CANCELLED' } 
          : order
      );
      setOrders(updatedOrders);
    }
  };
  
  const handleOrderSelect = (order) => {
    setSelectedOrder(order);
  };
  
  const handleRefreshOrders = () => {
    // In a real implementation, this would fetch the latest orders from the backend
    console.log('Refreshing orders');
  };
  
  const filteredOrders = filterStatus === 'ALL' 
    ? orders 
    : orders.filter(order => order.status === filterStatus);
  
  // Tab panel component
  const TabPanel = (props) => {
    const { children, value, index, ...other } = props;
    return (
      <div
        role="tabpanel"
        hidden={value !== index}
        id={`orderbook-tabpanel-${index}`}
        aria-labelledby={`orderbook-tab-${index}`}
        {...other}
      >
        {value === index && (
          <Box sx={{ p: 2 }}>
            {children}
          </Box>
        )}
      </div>
    );
  };
  
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Typography variant="h6" sx={{ mb: 2 }}>Order Book</Typography>
      
      <Grid container spacing={2} sx={{ flexGrow: 1, mb: 2 }}>
        {/* Order Entry Panel */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>Order Entry</Typography>
            <Divider sx={{ mb: 2 }} />
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Symbol</InputLabel>
              <Select
                value="NIFTY"
                label="Symbol"
              >
                <MenuItem value="NIFTY">NIFTY</MenuItem>
                <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
                <MenuItem value="RELIANCE">RELIANCE</MenuItem>
                <MenuItem value="TCS">TCS</MenuItem>
                <MenuItem value="INFY">INFY</MenuItem>
              </Select>
            </FormControl>
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Order Type</InputLabel>
              <Select
                value="LIMIT"
                label="Order Type"
              >
                <MenuItem value="MARKET">MARKET</MenuItem>
                <MenuItem value="LIMIT">LIMIT</MenuItem>
                <MenuItem value="SL">STOP LOSS</MenuItem>
                <MenuItem value="SL-M">STOP LOSS MARKET</MenuItem>
              </Select>
            </FormControl>
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Product Type</InputLabel>
              <Select
                value="MIS"
                label="Product Type"
              >
                <MenuItem value="MIS">MIS</MenuItem>
                <MenuItem value="NRML">NRML</MenuItem>
                <MenuItem value="CNC">CNC</MenuItem>
              </Select>
            </FormControl>
            
            <TextField
              fullWidth
              label="Quantity"
              type="number"
              value="75"
              sx={{ mb: 2 }}
            />
            
            <TextField
              fullWidth
              label="Price"
              type="number"
              value="22500"
              sx={{ mb: 2 }}
            />
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Strategy</InputLabel>
              <Select
                value="MARU"
                label="Strategy"
              >
                <MenuItem value="MARU">MARU</MenuItem>
                <MenuItem value="LUX24VR">LUX24VR</MenuItem>
                <MenuItem value="NFTTR">NFTTR</MenuItem>
                <MenuItem value="MANUAL">MANUAL</MenuItem>
              </Select>
            </FormControl>
            
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Button 
                variant="contained" 
                color="primary" 
                sx={{ width: '48%' }}
              >
                BUY
              </Button>
              <Button 
                variant="contained" 
                color="error" 
                sx={{ width: '48%' }}
              >
                SELL
              </Button>
            </Box>
          </Paper>
        </Grid>
        
        {/* Order Book Panel */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="subtitle1">Orders</Typography>
              <Box>
                <FormControl sx={{ minWidth: 120, mr: 1 }}>
                  <InputLabel size="small">Status</InputLabel>
                  <Select
                    value={filterStatus}
                    label="Status"
                    size="small"
                    onChange={(e) => setFilterStatus(e.target.value)}
                  >
                    <MenuItem value="ALL">All</MenuItem>
                    <MenuItem value="OPEN">Open</MenuItem>
                    <MenuItem value="EXECUTED">Executed</MenuItem>
                    <MenuItem value="CANCELLED">Cancelled</MenuItem>
                    <MenuItem value="REJECTED">Rejected</MenuItem>
                  </Select>
                </FormControl>
                <IconButton size="small" onClick={handleRefreshOrders}>
                  <RefreshIcon />
                </IconButton>
                <IconButton size="small" onClick={handleAddOrder}>
                  <AddIcon />
                </IconButton>
              </Box>
            </Box>
            
            <Divider sx={{ mb: 2 }} />
            
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
              <Tabs value={activeTab} onChange={handleTabChange}>
                <Tab label="Today's Orders" />
                <Tab label="Order History" />
                <Tab label="Trade History" />
              </Tabs>
            </Box>
            
            {/* Tab 1: Today's Orders */}
            <TabPanel value={activeTab} index={0}>
              <TableContainer sx={{ maxHeight: 400 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Time</TableCell>
                      <TableCell>Symbol</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Side</TableCell>
                      <TableCell>Qty</TableCell>
                      <TableCell>Price</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Strategy</TableCell>
                      <TableCell>Action</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredOrders.map((order) => (
                      <TableRow 
                        key={order.id}
                        hover
                        onClick={() => handleOrderSelect(order)}
                        sx={{ 
                          cursor: 'pointer',
                          bgcolor: order.status === 'REJECTED' ? 'rgba(244, 67, 54, 0.1)' : 
                                  order.status === 'EXECUTED' ? 'rgba(76, 175, 80, 0.1)' : 
                                  order.status === 'CANCELLED' ? 'rgba(158, 158, 158, 0.1)' : 'inherit'
                        }}
                      >
                        <TableCell>{order.time}</TableCell>
                        <TableCell>{order.symbol}</TableCell>
                        <TableCell>{order.type}</TableCell>
                        <TableCell sx={{ color: order.side === 'BUY' ? 'primary.main' : 'error.main', fontWeight: 'bold' }}>
                          {order.side}
                        </TableCell>
                        <TableCell>{order.quantity}</TableCell>
                        <TableCell>{order.price}</TableCell>
                        <TableCell>
                          <Box 
                            component="span" 
                            sx={{ 
                              color: order.status === 'EXECUTED' ? 'success.main' : 
                                    order.status === 'REJECTED' ? 'error.main' : 
                                    order.status === 'CANCELLED' ? 'text.secondary' : 'primary.main',
                              fontWeight: 'bold'
                            }}
                          >
                            {order.status}
                          </Box>
                        </TableCell>
                        <TableCell>{order.strategy}</TableCell>
                        <TableCell>
                          {order.status === 'OPEN' && (
                            <IconButton 
                              size="small" 
                              onClick={(e) => {
                                e.stopPropagation();
                                handleCancelOrder(order.id);
                              }}
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </TabPanel>
            
            {/* Tab 2: Order History */}
            <TabPanel value={activeTab} index={1}>
              <TableContainer sx={{ maxHeight: 400 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Date</TableCell>
                      <TableCell>Time</TableCell>
                      <TableCell>Symbol</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Side</TableCell>
                      <TableCell>Qty</TableCell>
                      <TableCell>Price</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Strategy</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell>2025-04-06</TableCell>
                      <TableCell>15:25:45</TableCell>
                      <TableCell>NIFTY</TableCell>
                      <TableCell>LIMIT</TableCell>
                      <TableCell sx={{ color: 'primary.main', fontWeight: 'bold' }}>BUY</TableCell>
                      <TableCell>75</TableCell>
                      <TableCell>22450</TableCell>
                      <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>EXECUTED</TableCell>
                      <TableCell>MARU</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2025-04-06</TableCell>
                      <TableCell>14:30:12</TableCell>
                      <TableCell>BANKNIFTY</TableCell>
                      <TableCell>MARKET</TableCell>
                      <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>SELL</TableCell>
                      <TableCell>25</TableCell>
                      <TableCell>48650</TableCell>
                      <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>EXECUTED</TableCell>
                      <TableCell>LUX24VR</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2025-04-05</TableCell>
                      <TableCell>15:28:33</TableCell>
                      <TableCell>NIFTY</TableCell>
                      <TableCell>SL</TableCell>
                      <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>SELL</TableCell>
                      <TableCell>50</TableCell>
                      <TableCell>22350</TableCell>
                      <TableCell sx={{ color: 'text.secondary', fontWeight: 'bold' }}>CANCELLED</TableCell>
                      <TableCell>MARU</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </TabPanel>
            
            {/* Tab 3: Trade History */}
            <TabPanel value={activeTab} index={2}>
              <TableContainer sx={{ maxHeight: 400 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Date</TableCell>
                      <TableCell>Time</TableCell>
                      <TableCell>Symbol</TableCell>
                      <TableCell>Side</TableCell>
                      <TableCell>Qty</TableCell>
                      <TableCell>Price</TableCell>
                      <TableCell>Value</TableCell>
                      <TableCell>P&L</TableCell>
                      <TableCell>Strategy</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell>2025-04-06</TableCell>
                      <TableCell>15:25:45</TableCell>
                      <TableCell>NIFTY</TableCell>
                      <TableCell sx={{ color: 'primary.main', fontWeight: 'bold' }}>BUY</TableCell>
                      <TableCell>75</TableCell>
                      <TableCell>22450</TableCell>
                      <TableCell>₹1,683,750</TableCell>
                      <TableCell>-</TableCell>
                      <TableCell>MARU</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2025-04-06</TableCell>
                      <TableCell>14:30:12</TableCell>
                      <TableCell>BANKNIFTY</TableCell>
                      <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>SELL</TableCell>
                      <TableCell>25</TableCell>
                      <TableCell>48650</TableCell>
                      <TableCell>₹1,216,250</TableCell>
                      <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹3,750</TableCell>
                      <TableCell>LUX24VR</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2025-04-06</TableCell>
                      <TableCell>10:15:30</TableCell>
                      <TableCell>BANKNIFTY</TableCell>
                      <TableCell sx={{ color: 'primary.main', fontWeight: 'bold' }}>BUY</TableCell>
                      <TableCell>25</TableCell>
                      <TableCell>48500</TableCell>
                      <TableCell>₹1,212,500</TableCell>
                      <TableCell>-</TableCell>
                      <TableCell>LUX24VR</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>2025-04-05</TableCell>
                      <TableCell>15:20:45</TableCell>
                      <TableCell>NIFTY</TableCell>
                      <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>SELL</TableCell>
                      <TableCell>50</TableCell>
                      <TableCell>22550</TableCell>
                      <TableCell>₹1,127,500</TableCell>
                      <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹5,000</TableCell>
                      <TableCell>MARU</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </TabPanel>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Logs Panel */}
      <Box sx={{ mt: 2 }}>
        <LogsPanel />
      </Box>
      
      {/* Add Order Dialog */}
      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Place New Order</DialogTitle>
        <DialogContent>
          <FormControl fullWidth sx={{ mb: 2, mt: 2 }}>
            <InputLabel>Symbol</InputLabel>
            <Select
              value={newOrder.symbol}
              label="Symbol"
              onChange={(e) => setNewOrder({...newOrder, symbol: e.target.value})}
            >
              <MenuItem value="NIFTY">NIFTY</MenuItem>
              <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
              <MenuItem value="RELIANCE">RELIANCE</MenuItem>
              <MenuItem value="TCS">TCS</MenuItem>
              <MenuItem value="INFY">INFY</MenuItem>
            </Select>
          </FormControl>
          
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Order Type</InputLabel>
            <Select
              value={newOrder.type}
              label="Order Type"
              onChange={(e) => setNewOrder({...newOrder, type: e.target.value})}
            >
              <MenuItem value="MARKET">MARKET</MenuItem>
              <MenuItem value="LIMIT">LIMIT</MenuItem>
              <MenuItem value="SL">STOP LOSS</MenuItem>
              <MenuItem value="SL-M">STOP LOSS MARKET</MenuItem>
            </Select>
          </FormControl>
          
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Side</InputLabel>
            <Select
              value={newOrder.side}
              label="Side"
              onChange={(e) => setNewOrder({...newOrder, side: e.target.value})}
            >
              <MenuItem value="BUY">BUY</MenuItem>
              <MenuItem value="SELL">SELL</MenuItem>
            </Select>
          </FormControl>
          
          <TextField
            fullWidth
            label="Quantity"
            type="number"
            value={newOrder.quantity}
            onChange={(e) => setNewOrder({...newOrder, quantity: e.target.value})}
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Price"
            type="number"
            value={newOrder.price}
            onChange={(e) => setNewOrder({...newOrder, price: e.target.value})}
            sx={{ mb: 2 }}
            disabled={newOrder.type === 'MARKET'}
          />
          
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Strategy</InputLabel>
            <Select
              value={newOrder.strategy}
              label="Strategy"
              onChange={(e) => setNewOrder({...newOrder, strategy: e.target.value})}
            >
              <MenuItem value="MARU">MARU</MenuItem>
              <MenuItem value="LUX24VR">LUX24VR</MenuItem>
              <MenuItem value="NFTTR">NFTTR</MenuItem>
              <MenuItem value="MANUAL">MANUAL</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handlePlaceOrder} variant="contained">Place Order</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default OrderBook;
