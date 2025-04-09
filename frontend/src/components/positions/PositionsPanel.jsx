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
  DialogActions,
  Chip
} from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import FilterListIcon from '@mui/icons-material/FilterList';
import CloseIcon from '@mui/icons-material/Close';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import LogsPanel from '../logs/LogsPanel';

const PositionsPanel = () => {
  const [positions, setPositions] = useState([
    { 
      id: 'pos1', 
      symbol: 'NIFTY', 
      type: 'LONG',
      quantity: 75,
      entryPrice: 22500,
      currentPrice: 22650,
      pnl: 11250,
      pnlPercent: 0.67,
      strategy: 'MARU',
      openTime: '09:31:45'
    },
    { 
      id: 'pos2', 
      symbol: 'BANKNIFTY', 
      type: 'SHORT',
      quantity: 25,
      entryPrice: 48750,
      currentPrice: 48650,
      pnl: 2500,
      pnlPercent: 0.21,
      strategy: 'LUX24VR',
      openTime: '09:25:45'
    },
    { 
      id: 'pos3', 
      symbol: 'RELIANCE', 
      type: 'LONG',
      quantity: 100,
      entryPrice: 2750.50,
      currentPrice: 2735.25,
      pnl: -1525,
      pnlPercent: -0.55,
      strategy: 'NFTTR',
      openTime: '09:40:22'
    }
  ]);
  
  const [activeTab, setActiveTab] = useState(0);
  const [openDialog, setOpenDialog] = useState(false);
  const [selectedPosition, setSelectedPosition] = useState(null);
  const [filterStrategy, setFilterStrategy] = useState('ALL');
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handlePositionSelect = (position) => {
    setSelectedPosition(position);
  };
  
  const handleRefreshPositions = () => {
    // In a real implementation, this would fetch the latest positions from the backend
    console.log('Refreshing positions');
    
    // Simulate price updates
    const updatedPositions = positions.map(position => {
      const priceChange = (Math.random() * 20 - 10).toFixed(2); // Random price change between -10 and +10
      const newPrice = parseFloat(position.currentPrice) + parseFloat(priceChange);
      const priceDiff = newPrice - position.entryPrice;
      const newPnl = priceDiff * position.quantity;
      const newPnlPercent = (priceDiff / position.entryPrice * 100).toFixed(2);
      
      return {
        ...position,
        currentPrice: newPrice,
        pnl: newPnl,
        pnlPercent: newPnlPercent
      };
    });
    
    setPositions(updatedPositions);
  };
  
  const handleClosePosition = (positionId) => {
    if (window.confirm('Are you sure you want to close this position?')) {
      // In a real implementation, this would send a request to close the position
      console.log('Closing position:', positionId);
      
      // For demo purposes, remove the position from the list
      const updatedPositions = positions.filter(position => position.id !== positionId);
      setPositions(updatedPositions);
      
      if (selectedPosition && selectedPosition.id === positionId) {
        setSelectedPosition(null);
      }
    }
  };
  
  const filteredPositions = filterStrategy === 'ALL' 
    ? positions 
    : positions.filter(position => position.strategy === filterStrategy);
  
  const totalPnl = positions.reduce((sum, position) => sum + position.pnl, 0);
  const totalPnlPercent = (totalPnl / positions.reduce((sum, position) => sum + (position.entryPrice * position.quantity), 0) * 100).toFixed(2);
  
  // Tab panel component
  const TabPanel = (props) => {
    const { children, value, index, ...other } = props;
    return (
      <div
        role="tabpanel"
        hidden={value !== index}
        id={`positions-tabpanel-${index}`}
        aria-labelledby={`positions-tab-${index}`}
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
      <Typography variant="h6" sx={{ mb: 2 }}>Positions</Typography>
      
      <Grid container spacing={2} sx={{ flexGrow: 1, mb: 2 }}>
        {/* Summary Panel */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>Summary</Typography>
            <Divider sx={{ mb: 2 }} />
            
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Paper variant="outlined" sx={{ p: 2, mb: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Open Positions
                  </Typography>
                  <Typography variant="h4">
                    {positions.length}
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6}>
                <Paper variant="outlined" sx={{ p: 2, mb: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Total P&L
                  </Typography>
                  <Typography 
                    variant="h4" 
                    sx={{ 
                      color: totalPnl > 0 ? 'success.main' : totalPnl < 0 ? 'error.main' : 'text.primary',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center'
                    }}
                  >
                    {totalPnl > 0 ? <TrendingUpIcon sx={{ mr: 0.5 }} /> : totalPnl < 0 ? <TrendingDownIcon sx={{ mr: 0.5 }} /> : null}
                    ₹{Math.abs(totalPnl).toLocaleString()}
                  </Typography>
                  <Typography 
                    variant="body2" 
                    sx={{ 
                      color: totalPnl > 0 ? 'success.main' : totalPnl < 0 ? 'error.main' : 'text.primary',
                      fontWeight: 'bold'
                    }}
                  >
                    ({totalPnl > 0 ? '+' : ''}{totalPnlPercent}%)
                  </Typography>
                </Paper>
              </Grid>
            </Grid>
            
            <Typography variant="subtitle2" gutterBottom>Position Distribution</Typography>
            
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                By Strategy
              </Typography>
              
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                {Array.from(new Set(positions.map(p => p.strategy))).map(strategy => (
                  <Chip 
                    key={strategy} 
                    label={`${strategy} (${positions.filter(p => p.strategy === strategy).length})`} 
                    color="primary" 
                    variant={filterStrategy === strategy ? "filled" : "outlined"}
                    onClick={() => setFilterStrategy(filterStrategy === strategy ? 'ALL' : strategy)}
                  />
                ))}
                {positions.length > 0 && (
                  <Chip 
                    label="All" 
                    color="primary" 
                    variant={filterStrategy === 'ALL' ? "filled" : "outlined"}
                    onClick={() => setFilterStrategy('ALL')}
                  />
                )}
              </Box>
            </Box>
            
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                By Type
              </Typography>
              
              <Grid container spacing={1}>
                <Grid item xs={6}>
                  <Paper variant="outlined" sx={{ p: 1, textAlign: 'center' }}>
                    <Typography variant="body2" color="primary.main" gutterBottom>
                      LONG
                    </Typography>
                    <Typography variant="h6">
                      {positions.filter(p => p.type === 'LONG').length}
                    </Typography>
                  </Paper>
                </Grid>
                
                <Grid item xs={6}>
                  <Paper variant="outlined" sx={{ p: 1, textAlign: 'center' }}>
                    <Typography variant="body2" color="error.main" gutterBottom>
                      SHORT
                    </Typography>
                    <Typography variant="h6">
                      {positions.filter(p => p.type === 'SHORT').length}
                    </Typography>
                  </Paper>
                </Grid>
              </Grid>
            </Box>
            
            <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
              <Button 
                variant="outlined" 
                startIcon={<RefreshIcon />}
                onClick={handleRefreshPositions}
              >
                Refresh
              </Button>
            </Box>
          </Paper>
        </Grid>
        
        {/* Positions Panel */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="subtitle1">Open Positions</Typography>
              <Box>
                <IconButton size="small" onClick={handleRefreshPositions}>
                  <RefreshIcon />
                </IconButton>
              </Box>
            </Box>
            
            <Divider sx={{ mb: 2 }} />
            
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
              <Tabs value={activeTab} onChange={handleTabChange}>
                <Tab label="Open Positions" />
                <Tab label="Closed Positions" />
                <Tab label="Performance" />
              </Tabs>
            </Box>
            
            {/* Tab 1: Open Positions */}
            <TabPanel value={activeTab} index={0}>
              <TableContainer sx={{ maxHeight: 400 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Symbol</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Qty</TableCell>
                      <TableCell>Entry Price</TableCell>
                      <TableCell>Current Price</TableCell>
                      <TableCell>P&L</TableCell>
                      <TableCell>Strategy</TableCell>
                      <TableCell>Open Time</TableCell>
                      <TableCell>Action</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredPositions.map((position) => (
                      <TableRow 
                        key={position.id}
                        hover
                        onClick={() => handlePositionSelect(position)}
                        sx={{ 
                          cursor: 'pointer',
                          bgcolor: position.pnl > 0 ? 'rgba(76, 175, 80, 0.05)' : 
                                  position.pnl < 0 ? 'rgba(244, 67, 54, 0.05)' : 'inherit'
                        }}
                      >
                        <TableCell>{position.symbol}</TableCell>
                        <TableCell sx={{ color: position.type === 'LONG' ? 'primary.main' : 'error.main', fontWeight: 'bold' }}>
                          {position.type}
                        </TableCell>
                        <TableCell>{position.quantity}</TableCell>
                        <TableCell>{position.entryPrice.toLocaleString()}</TableCell>
                        <TableCell>{position.currentPrice.toLocaleString()}</TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                            <Typography 
                              variant="body2" 
                              sx={{ 
                                color: position.pnl > 0 ? 'success.main' : position.pnl < 0 ? 'error.main' : 'text.primary',
                                fontWeight: 'bold',
                                display: 'flex',
                                alignItems: 'center'
                              }}
                            >
                              {position.pnl > 0 ? <TrendingUpIcon fontSize="small" sx={{ mr: 0.5 }} /> : 
                               position.pnl < 0 ? <TrendingDownIcon fontSize="small" sx={{ mr: 0.5 }} /> : null}
                              ₹{Math.abs(position.pnl).toLocaleString()}
                            </Typography>
                            <Typography 
                              variant="caption" 
                              sx={{ 
                                color: position.pnl > 0 ? 'success.main' : position.pnl < 0 ? 'error.main' : 'text.primary'
                              }}
                            >
                              ({position.pnl > 0 ? '+' : ''}{position.pnlPercent}%)
                            </Typography>
                          </Box>
                        </TableCell>
                        <TableCell>{position.strategy}</TableCell>
                        <TableCell>{position.openTime}</TableCell>
                        <TableCell>
                          <IconButton 
                            size="small" 
                            onClick={(e) => {
                              e.stopPropagation();
                              handleClosePosition(position.id);
                            }}
                            title="Close Position"
                          >
                            <CloseIcon fontSize="small" />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))}
                    {filteredPositions.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={9} align="center" sx={{ py: 3 }}>
                          <Typography variant="body2" color="textSecondary">
                            No open positions found
                          </Typography>
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </TableContainer>
            </TabPanel>
            
            {/* Tab 2: Closed Positions */}
            <TabPanel value={activeTab} index={1}>
              <TableContainer sx={{ maxHeight: 400 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Symbol</TableCell>
                      <TableCell>Type</TableCell>
                      <TableCell>Qty</TableCell>
                      <TableCell>Entry Price</TableCell>
                      <TableCell>Exit Price</TableCell>
                      <TableCell>P&L</TableCell>
                      <TableCell>Strategy</TableCell>
                      <TableCell>Open Time</TableCell>
                      <TableCell>Close Time</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell>NIFTY</TableCell>
                      <TableCell sx={{ color: 'primary.main', fontWeight: 'bold' }}>LONG</TableCell>
                      <TableCell>50</TableCell>
                      <TableCell>22450</TableCell>
                      <TableCell>22550</TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                          <Typography 
                            variant="body2" 
                            sx={{ 
                              color: 'success.main',
                              fontWeight: 'bold',
                              display: 'flex',
                              alignItems: 'center'
                            }}
                          >
                            <TrendingUpIcon fontSize="small" sx={{ mr: 0.5 }} />
                            ₹5,000
                          </Typography>
                          <Typography variant="caption" sx={{ color: 'success.main' }}>
                            (+0.44%)
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>MARU</TableCell>
                      <TableCell>09:15:30</TableCell>
                      <TableCell>15:20:45</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>BANKNIFTY</TableCell>
                      <TableCell sx={{ color: 'primary.main', fontWeight: 'bold' }}>LONG</TableCell>
                      <TableCell>25</TableCell>
                      <TableCell>48500</TableCell>
                      <TableCell>48650</TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                          <Typography 
                            variant="body2" 
                            sx={{ 
                              color: 'success.main',
                              fontWeight: 'bold',
                              display: 'flex',
                              alignItems: 'center'
                            }}
                          >
                            <TrendingUpIcon fontSize="small" sx={{ mr: 0.5 }} />
                            ₹3,750
                          </Typography>
                          <Typography variant="caption" sx={{ color: 'success.main' }}>
                            (+0.31%)
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>LUX24VR</TableCell>
                      <TableCell>09:15:45</TableCell>
                      <TableCell>14:30:12</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>RELIANCE</TableCell>
                      <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>SHORT</TableCell>
                      <TableCell>50</TableCell>
                      <TableCell>2780.25</TableCell>
                      <TableCell>2755.50</TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                          <Typography 
                            variant="body2" 
                            sx={{ 
                              color: 'success.main',
                              fontWeight: 'bold',
                              display: 'flex',
                              alignItems: 'center'
                            }}
                          >
                            <TrendingUpIcon fontSize="small" sx={{ mr: 0.5 }} />
                            ₹1,237.50
                          </Typography>
                          <Typography variant="caption" sx={{ color: 'success.main' }}>
                            (+0.89%)
                          </Typography>
                        </Box>
                      </TableCell>
                      <TableCell>NFTTR</TableCell>
                      <TableCell>09:20:15</TableCell>
                      <TableCell>11:45:30</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </TabPanel>
            
            {/* Tab 3: Performance */}
            <TabPanel value={activeTab} index={2}>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>Daily P&L</Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                      <Typography variant="h5" sx={{ color: 'success.main', fontWeight: 'bold' }}>
                        +₹15,750
                      </Typography>
                      <Typography variant="body2" sx={{ color: 'success.main' }}>
                        +1.25%
                      </Typography>
                    </Box>
                  </Paper>
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>Weekly P&L</Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                      <Typography variant="h5" sx={{ color: 'success.main', fontWeight: 'bold' }}>
                        +₹42,500
                      </Typography>
                      <Typography variant="body2" sx={{ color: 'success.main' }}>
                        +3.40%
                      </Typography>
                    </Box>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>Performance by Strategy</Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Strategy</TableCell>
                            <TableCell>Trades</TableCell>
                            <TableCell>Win Rate</TableCell>
                            <TableCell>P&L</TableCell>
                            <TableCell>Return</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          <TableRow>
                            <TableCell>MARU</TableCell>
                            <TableCell>24</TableCell>
                            <TableCell>75%</TableCell>
                            <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹28,500</TableCell>
                            <TableCell sx={{ color: 'success.main' }}>+2.28%</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>LUX24VR</TableCell>
                            <TableCell>18</TableCell>
                            <TableCell>67%</TableCell>
                            <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹15,750</TableCell>
                            <TableCell sx={{ color: 'success.main' }}>+1.26%</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>NFTTR</TableCell>
                            <TableCell>12</TableCell>
                            <TableCell>58%</TableCell>
                            <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>-₹1,750</TableCell>
                            <TableCell sx={{ color: 'error.main' }}>-0.14%</TableCell>
                          </TableRow>
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>Performance by Instrument</Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Instrument</TableCell>
                            <TableCell>Trades</TableCell>
                            <TableCell>Win Rate</TableCell>
                            <TableCell>P&L</TableCell>
                            <TableCell>Return</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          <TableRow>
                            <TableCell>NIFTY</TableCell>
                            <TableCell>28</TableCell>
                            <TableCell>71%</TableCell>
                            <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹32,250</TableCell>
                            <TableCell sx={{ color: 'success.main' }}>+2.58%</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>BANKNIFTY</TableCell>
                            <TableCell>15</TableCell>
                            <TableCell>67%</TableCell>
                            <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>+₹18,750</TableCell>
                            <TableCell sx={{ color: 'success.main' }}>+1.50%</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>RELIANCE</TableCell>
                            <TableCell>11</TableCell>
                            <TableCell>55%</TableCell>
                            <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>-₹8,500</TableCell>
                            <TableCell sx={{ color: 'error.main' }}>-0.68%</TableCell>
                          </TableRow>
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Paper>
                </Grid>
              </Grid>
            </TabPanel>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Logs Panel */}
      <Box sx={{ mt: 2 }}>
        <LogsPanel />
      </Box>
    </Box>
  );
};

export default PositionsPanel;
