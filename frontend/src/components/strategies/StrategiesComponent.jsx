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
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Checkbox,
  Tooltip,
  Alert,
  Snackbar,
  Chip
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import SaveIcon from '@mui/icons-material/Save';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import StopIcon from '@mui/icons-material/Stop';
import RefreshIcon from '@mui/icons-material/Refresh';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import WarningIcon from '@mui/icons-material/Warning';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TrendingDownIcon from '@mui/icons-material/TrendingDown';
import LogsPanel from '../logs/LogsPanel';

const StrategiesComponent = () => {
  const [strategies, setStrategies] = useState([
    { 
      id: 'strategy1', 
      enabled: true,
      name: 'LUX24VR', 
      description: 'Liquidity Utilization Extended 24h Volatility Response',
      type: 'Mean Reversion',
      instruments: ['NIFTY', 'BANKNIFTY'],
      status: 'Allowed',
      pnl: -130882.50,
      tradeValue: 1977130.00,
      marketOrders: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-06',
      manualSquareOff: true,
      markAsCompleted: true
    },
    { 
      id: 'strategy2', 
      enabled: true,
      name: 'NFTTR', 
      description: 'Non-Farm Trend Trading Response',
      type: 'Trend Following',
      instruments: ['NIFTY', 'BANKNIFTY'],
      status: 'Allowed',
      pnl: 0.00,
      tradeValue: 0.00,
      marketOrders: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-06',
      manualSquareOff: true,
      markAsCompleted: true
    },
    { 
      id: 'strategy3', 
      enabled: true,
      name: 'BNTTR', 
      description: 'BankNifty Trend Trading Response',
      type: 'Trend Following',
      instruments: ['BANKNIFTY'],
      status: 'Allowed',
      pnl: 0.00,
      tradeValue: 0.00,
      marketOrders: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-05',
      manualSquareOff: true,
      markAsCompleted: true
    },
    { 
      id: 'strategy4', 
      enabled: true,
      name: 'LUX24VRPE', 
      description: 'Liquidity Utilization Extended 24h Volatility Response with Price Extension',
      type: 'Mean Reversion',
      instruments: ['NIFTY', 'BANKNIFTY'],
      status: 'Allowed',
      pnl: 0.00,
      tradeValue: 0.00,
      marketOrders: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-06',
      manualSquareOff: true,
      markAsCompleted: true
    },
    { 
      id: 'strategy5', 
      enabled: true,
      name: 'NIFEMA3', 
      description: 'NIFTY Exponential Moving Average 3-period Strategy',
      type: 'Trend Following',
      instruments: ['NIFTY'],
      status: 'Allowed',
      pnl: -3393.75,
      tradeValue: 22187.25,
      marketOrders: 'Allowed',
      squareOffTime: '15:15:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-06',
      manualSquareOff: true,
      markAsCompleted: true
    },
    { 
      id: 'strategy6', 
      enabled: true,
      name: 'NF3STM2', 
      description: 'NIFTY 3-period Short-Term Momentum Strategy 2',
      type: 'Momentum',
      instruments: ['NIFTY'],
      status: 'Allowed',
      pnl: 0.00,
      tradeValue: 0.00,
      marketOrders: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      maxProfit: 0,
      maxLoss: 0,
      delayBetweenUsers: 0.00,
      uniqueIdReq: false,
      cancelPreviousSignal: true,
      lastRun: '2025-04-05',
      manualSquareOff: true,
      markAsCompleted: true
    }
  ]);
  
  const [selectedStrategy, setSelectedStrategy] = useState(null);
  const [activeTab, setActiveTab] = useState(0);
  const [openDialog, setOpenDialog] = useState(false);
  const [newStrategy, setNewStrategy] = useState({
    name: '',
    description: '',
    type: 'Range Breakout',
    instruments: ['NIFTY'],
    status: 'Allowed',
    squareOffTime: '15:25:00',
    userAccount: 'SIM1',
    enabled: true,
    manualSquareOff: true,
    markAsCompleted: true,
    uniqueIdReq: false,
    cancelPreviousSignal: true
  });
  const [isEditing, setIsEditing] = useState(false);
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success'
  });
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handleStrategySelect = (strategy) => {
    setSelectedStrategy(strategy);
    setActiveTab(0);
  };
  
  const handleAddStrategy = () => {
    setIsEditing(false);
    setNewStrategy({
      name: '',
      description: '',
      type: 'Range Breakout',
      instruments: ['NIFTY'],
      status: 'Allowed',
      squareOffTime: '15:25:00',
      userAccount: 'SIM1',
      enabled: true,
      manualSquareOff: true,
      markAsCompleted: true,
      uniqueIdReq: false,
      cancelPreviousSignal: true
    });
    setOpenDialog(true);
  };
  
  const handleEditStrategy = () => {
    if (!selectedStrategy) return;
    
    setIsEditing(true);
    setNewStrategy({
      name: selectedStrategy.name,
      description: selectedStrategy.description,
      type: selectedStrategy.type,
      instruments: [...selectedStrategy.instruments],
      status: selectedStrategy.status,
      squareOffTime: selectedStrategy.squareOffTime,
      userAccount: selectedStrategy.userAccount,
      enabled: selectedStrategy.enabled,
      manualSquareOff: selectedStrategy.manualSquareOff,
      markAsCompleted: selectedStrategy.markAsCompleted,
      uniqueIdReq: selectedStrategy.uniqueIdReq,
      cancelPreviousSignal: selectedStrategy.cancelPreviousSignal
    });
    setOpenDialog(true);
  };
  
  const handleCloseDialog = () => {
    setOpenDialog(false);
  };
  
  const handleSaveStrategy = () => {
    if (!newStrategy.name) {
      setSnackbar({
        open: true,
        message: 'Strategy name is required',
        severity: 'error'
      });
      return;
    }
    
    if (isEditing && selectedStrategy) {
      // Update existing strategy
      const updatedStrategies = strategies.map(s => 
        s.id === selectedStrategy.id 
          ? { 
              ...s, 
              name: newStrategy.name,
              description: newStrategy.description,
              type: newStrategy.type,
              instruments: [...newStrategy.instruments],
              status: newStrategy.status,
              squareOffTime: newStrategy.squareOffTime,
              userAccount: newStrategy.userAccount,
              enabled: newStrategy.enabled,
              manualSquareOff: newStrategy.manualSquareOff,
              markAsCompleted: newStrategy.markAsCompleted,
              uniqueIdReq: newStrategy.uniqueIdReq,
              cancelPreviousSignal: newStrategy.cancelPreviousSignal
            } 
          : s
      );
      setStrategies(updatedStrategies);
      setSelectedStrategy({
        ...selectedStrategy,
        name: newStrategy.name,
        description: newStrategy.description,
        type: newStrategy.type,
        instruments: [...newStrategy.instruments],
        status: newStrategy.status,
        squareOffTime: newStrategy.squareOffTime,
        userAccount: newStrategy.userAccount,
        enabled: newStrategy.enabled,
        manualSquareOff: newStrategy.manualSquareOff,
        markAsCompleted: newStrategy.markAsCompleted,
        uniqueIdReq: newStrategy.uniqueIdReq,
        cancelPreviousSignal: newStrategy.cancelPreviousSignal
      });
      
      setSnackbar({
        open: true,
        message: 'Strategy updated successfully',
        severity: 'success'
      });
    } else {
      // Create new strategy
      const newStrategyObj = {
        id: `strategy${Date.now()}`,
        name: newStrategy.name,
        description: newStrategy.description,
        type: newStrategy.type,
        instruments: [...newStrategy.instruments],
        status: newStrategy.status,
        pnl: 0.00,
        tradeValue: 0.00,
        marketOrders: 'Allowed',
        squareOffTime: newStrategy.squareOffTime,
        userAccount: newStrategy.userAccount,
        maxProfit: 0,
        maxLoss: 0,
        delayBetweenUsers: 0.00,
        uniqueIdReq: newStrategy.uniqueIdReq,
        cancelPreviousSignal: newStrategy.cancelPreviousSignal,
        lastRun: 'Never',
        enabled: newStrategy.enabled,
        manualSquareOff: newStrategy.manualSquareOff,
        markAsCompleted: newStrategy.markAsCompleted
      };
      
      setStrategies([...strategies, newStrategyObj]);
      setSelectedStrategy(newStrategyObj);
      
      setSnackbar({
        open: true,
        message: 'Strategy created successfully',
        severity: 'success'
      });
    }
    
    setOpenDialog(false);
  };
  
  const handleDeleteStrategy = (strategyId) => {
    const strategyToDelete = strategies.find(s => s.id === strategyId);
    if (!strategyToDelete) return;
    
    if (window.confirm(`Are you sure you want to delete the strategy "${strategyToDelete.name}"?`)) {
      const updatedStrategies = strategies.filter(s => s.id !== strategyId);
      setStrategies(updatedStrategies);
      
      if (selectedStrategy && selectedStrategy.id === strategyId) {
        setSelectedStrategy(null);
      }
      
      setSnackbar({
        open: true,
        message: 'Strategy deleted successfully',
        severity: 'success'
      });
    }
  };
  
  const handleToggleStatus = (strategyId) => {
    const updatedStrategies = strategies.map(s => 
      s.id === strategyId 
        ? { ...s, enabled: !s.enabled } 
        : s
    );
    
    setStrategies(updatedStrategies);
    
    if (selectedStrategy && selectedStrategy.id === strategyId) {
      setSelectedStrategy({
        ...selectedStrategy,
        enabled: !selectedStrategy.enabled
      });
    }
  };
  
  const handleToggleManualSquareOff = (strategyId) => {
    const updatedStrategies = strategies.map(s => 
      s.id === strategyId 
        ? { ...s, manualSquareOff: !s.manualSquareOff } 
        : s
    );
    
    setStrategies(updatedStrategies);
    
    if (selectedStrategy && selectedStrategy.id === strategyId) {
      setSelectedStrategy({
        ...selectedStrategy,
        manualSquareOff: !selectedStrategy.manualSquareOff
      });
    }
  };
  
  const handleToggleMarkAsCompleted = (strategyId) => {
    const updatedStrategies = strategies.map(s => 
      s.id === strategyId 
        ? { ...s, markAsCompleted: !s.markAsCompleted } 
        : s
    );
    
    setStrategies(updatedStrategies);
    
    if (selectedStrategy && selectedStrategy.id === strategyId) {
      setSelectedStrategy({
        ...selectedStrategy,
        markAsCompleted: !selectedStrategy.markAsCompleted
      });
    }
  };
  
  const handleCloseSnackbar = () => {
    setSnackbar({
      ...snackbar,
      open: false
    });
  };
  
  const handleRefreshStrategies = () => {
    // In a real implementation, this would fetch the latest strategy data from the backend
    console.log('Refreshing strategies');
    setSnackbar({
      open: true,
      message: 'Strategies refreshed',
      severity: 'info'
    });
  };
  
  // Tab panel component
  const TabPanel = (props) => {
    const { children, value, index, ...other } = props;
    return (
      <div
        role="tabpanel"
        hidden={value !== index}
        id={`strategy-tabpanel-${index}`}
        aria-labelledby={`strategy-tab-${index}`}
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
      <Typography variant="h6" sx={{ mb: 2 }}>Strategy Management</Typography>
      
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box>
          <Button 
            variant="contained" 
            size="small" 
            startIcon={<AddIcon />}
            onClick={handleAddStrategy}
            sx={{ mr: 1 }}
          >
            New Strategy
          </Button>
          <IconButton size="small" onClick={handleRefreshStrategies}>
            <RefreshIcon />
          </IconButton>
        </Box>
        <Typography variant="caption" color="textSecondary">
          Drag a column to group
        </Typography>
      </Box>
      
      <Paper sx={{ flexGrow: 1, mb: 2, overflow: 'hidden' }}>
        <TableContainer sx={{ maxHeight: 'calc(100vh - 200px)' }}>
          <Table stickyHeader size="small">
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox">
                  <Typography variant="subtitle2">Enabled</Typography>
                </TableCell>
                <TableCell>Delete</TableCell>
                <TableCell>Manual Square Off</TableCell>
                <TableCell>Mark As Completed</TableCell>
                <TableCell>Strategy Tag</TableCell>
                <TableCell>P&L</TableCell>
                <TableCell>Trade Value</TableCell>
                <TableCell>Market Orders</TableCell>
                <TableCell>SqOff Time</TableCell>
                <TableCell>User Account</TableCell>
                <TableCell>Max Profit</TableCell>
                <TableCell>Max Loss</TableCell>
                <TableCell>Delay Between Users</TableCell>
                <TableCell>Unique ID Req for Order</TableCell>
                <TableCell>Cancel Previous Open Signal</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {strategies.map((strategy) => (
                <TableRow 
                  key={strategy.id}
                  hover
                  selected={selectedStrategy && selectedStrategy.id === strategy.id}
                  onClick={() => handleStrategySelect(strategy)}
                  sx={{ 
                    cursor: 'pointer',
                    bgcolor: strategy.enabled ? 'rgba(76, 175, 80, 0.05)' : 'rgba(244, 67, 54, 0.05)'
                  }}
                >
                  <TableCell padding="checkbox">
                    <Checkbox
                      checked={strategy.enabled}
                      onChange={() => handleToggleStatus(strategy.id)}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="error"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeleteStrategy(strategy.id);
                      }}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <Checkbox
                      checked={strategy.manualSquareOff}
                      onChange={() => handleToggleManualSquareOff(strategy.id)}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </TableCell>
                  <TableCell>
                    <Checkbox
                      checked={strategy.markAsCompleted}
                      onChange={() => handleToggleMarkAsCompleted(strategy.id)}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </TableCell>
                  <TableCell>{strategy.name}</TableCell>
                  <TableCell sx={{ 
                    color: strategy.pnl > 0 ? 'success.main' : strategy.pnl < 0 ? 'error.main' : 'text.primary',
                    fontWeight: strategy.pnl !== 0 ? 'bold' : 'normal'
                  }}>
                    {strategy.pnl.toFixed(2)}
                  </TableCell>
                  <TableCell>{strategy.tradeValue.toFixed(2)}</TableCell>
                  <TableCell>{strategy.marketOrders}</TableCell>
                  <TableCell>{strategy.squareOffTime}</TableCell>
                  <TableCell>{strategy.userAccount}</TableCell>
                  <TableCell>{strategy.maxProfit}</TableCell>
                  <TableCell>{strategy.maxLoss}</TableCell>
                  <TableCell>{strategy.delayBetweenUsers.toFixed(2)}</TableCell>
                  <TableCell>
                    <Checkbox
                      checked={strategy.uniqueIdReq}
                      disabled
                    />
                  </TableCell>
                  <TableCell>
                    <Checkbox
                      checked={strategy.cancelPreviousSignal}
                      disabled
                    />
                  </TableCell>
                </TableRow>
              ))}
              {strategies.length === 0 && (
                <TableRow>
                  <TableCell colSpan={15} align="center" sx={{ py: 3 }}>
                    <Typography variant="body2" color="textSecondary">
                      No strategies found
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
      
      {selectedStrategy && (
        <Paper sx={{ p: 2, mb: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6">{selectedStrategy.name}</Typography>
            <Box>
              <Button 
                variant={selectedStrategy.enabled ? "contained" : "outlined"} 
                color={selectedStrategy.enabled ? "error" : "success"}
                startIcon={selectedStrategy.enabled ? <StopIcon /> : <PlayArrowIcon />}
                onClick={() => handleToggleStatus(selectedStrategy.id)}
                sx={{ mr: 1 }}
              >
                {selectedStrategy.enabled ? "Disable" : "Enable"}
              </Button>
              <Button 
                variant="outlined" 
                startIcon={<EditIcon />}
                onClick={handleEditStrategy}
                sx={{ mr: 1 }}
              >
                Edit
              </Button>
              <Button 
                variant="outlined" 
                color="error"
                startIcon={<DeleteIcon />}
                onClick={() => handleDeleteStrategy(selectedStrategy.id)}
              >
                Delete
              </Button>
            </Box>
          </Box>
          
          <Divider sx={{ mb: 2 }} />
          
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs value={activeTab} onChange={handleTabChange}>
              <Tab label="Overview" />
              <Tab label="Parameters" />
              <Tab label="Backtest Results" />
              <Tab label="Performance" />
              <Tab label="Logs" />
            </Tabs>
          </Box>
          
          {/* Tab 1: Overview */}
          <TabPanel value={activeTab} index={0}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" gutterBottom>Description</Typography>
                <Typography variant="body2" paragraph>
                  {selectedStrategy.description}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>Strategy Type</Typography>
                <Typography variant="body2" paragraph>
                  {selectedStrategy.type}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>Instruments</Typography>
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                  {selectedStrategy.instruments.map(instrument => (
                    <Chip key={instrument} label={instrument} color="primary" variant="outlined" />
                  ))}
                </Box>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" gutterBottom>Status</Typography>
                <Typography 
                  variant="body2" 
                  paragraph
                  sx={{ 
                    color: selectedStrategy.enabled ? 'success.main' : 'error.main',
                    fontWeight: 'bold',
                    display: 'flex',
                    alignItems: 'center'
                  }}
                >
                  {selectedStrategy.enabled ? (
                    <>
                      <CheckCircleIcon fontSize="small" sx={{ mr: 0.5 }} />
                      Enabled
                    </>
                  ) : (
                    <>
                      <ErrorIcon fontSize="small" sx={{ mr: 0.5 }} />
                      Disabled
                    </>
                  )}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>P&L</Typography>
                <Typography 
                  variant="body2" 
                  paragraph
                  sx={{ 
                    color: selectedStrategy.pnl > 0 ? 'success.main' : selectedStrategy.pnl < 0 ? 'error.main' : 'text.primary',
                    fontWeight: 'bold',
                    display: 'flex',
                    alignItems: 'center'
                  }}
                >
                  {selectedStrategy.pnl > 0 ? (
                    <>
                      <TrendingUpIcon fontSize="small" sx={{ mr: 0.5 }} />
                      ₹{selectedStrategy.pnl.toFixed(2)}
                    </>
                  ) : selectedStrategy.pnl < 0 ? (
                    <>
                      <TrendingDownIcon fontSize="small" sx={{ mr: 0.5 }} />
                      ₹{Math.abs(selectedStrategy.pnl).toFixed(2)}
                    </>
                  ) : (
                    <>₹{selectedStrategy.pnl.toFixed(2)}</>
                  )}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>Last Run</Typography>
                <Typography variant="body2" paragraph>
                  {selectedStrategy.lastRun}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>Square Off Time</Typography>
                <Typography variant="body2" paragraph>
                  {selectedStrategy.squareOffTime}
                </Typography>
                
                <Typography variant="subtitle2" gutterBottom>User Account</Typography>
                <Typography variant="body2" paragraph>
                  {selectedStrategy.userAccount}
                </Typography>
              </Grid>
            </Grid>
          </TabPanel>
          
          {/* Tab 2: Parameters */}
          <TabPanel value={activeTab} index={1}>
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Strategy Parameters</Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Entry Type</InputLabel>
                  <Select
                    value="breakout"
                    label="Entry Type"
                    disabled={selectedStrategy.enabled}
                  >
                    <MenuItem value="breakout">Breakout</MenuItem>
                    <MenuItem value="reversion">Mean Reversion</MenuItem>
                    <MenuItem value="trend">Trend Following</MenuItem>
                  </Select>
                </FormControl>
                
                <TextField
                  fullWidth
                  label="Entry Threshold"
                  type="number"
                  value="0.5"
                  sx={{ mb: 2 }}
                  disabled={selectedStrategy.enabled}
                />
                
                <FormControlLabel
                  control={<Switch checked={true} disabled={selectedStrategy.enabled} />}
                  label="Use ATR for Entry"
                  sx={{ mb: 1, display: 'block' }}
                />
                
                <TextField
                  fullWidth
                  label="ATR Multiplier"
                  type="number"
                  value="1.5"
                  sx={{ mb: 2 }}
                  disabled={selectedStrategy.enabled}
                />
              </Grid>
              
              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Exit Type</InputLabel>
                  <Select
                    value="target"
                    label="Exit Type"
                    disabled={selectedStrategy.enabled}
                  >
                    <MenuItem value="target">Target & Stoploss</MenuItem>
                    <MenuItem value="trailing">Trailing Stoploss</MenuItem>
                    <MenuItem value="time">Time-based</MenuItem>
                  </Select>
                </FormControl>
                
                <TextField
                  fullWidth
                  label="Profit Target (%)"
                  type="number"
                  value="1.5"
                  sx={{ mb: 2 }}
                  disabled={selectedStrategy.enabled}
                />
                
                <TextField
                  fullWidth
                  label="Stoploss (%)"
                  type="number"
                  value="0.8"
                  sx={{ mb: 2 }}
                  disabled={selectedStrategy.enabled}
                />
                
                <FormControlLabel
                  control={<Switch checked={true} disabled={selectedStrategy.enabled} />}
                  label="Enable Trailing Stoploss"
                  sx={{ mb: 1, display: 'block' }}
                />
              </Grid>
              
              <Grid item xs={12}>
                <Divider sx={{ my: 2 }} />
                <Typography variant="subtitle2" sx={{ mb: 2 }}>Advanced Settings</Typography>
                
                <Grid container spacing={2}>
                  <Grid item xs={12} md={6}>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={selectedStrategy.manualSquareOff}
                          onChange={() => handleToggleManualSquareOff(selectedStrategy.id)}
                          disabled={selectedStrategy.enabled}
                        />
                      }
                      label="Enable Manual Square Off"
                      sx={{ mb: 1, display: 'block' }}
                    />
                    
                    <FormControlLabel
                      control={
                        <Switch
                          checked={selectedStrategy.markAsCompleted}
                          onChange={() => handleToggleMarkAsCompleted(selectedStrategy.id)}
                          disabled={selectedStrategy.enabled}
                        />
                      }
                      label="Mark As Completed"
                      sx={{ mb: 1, display: 'block' }}
                    />
                    
                    <FormControlLabel
                      control={
                        <Switch
                          checked={selectedStrategy.uniqueIdReq}
                          onChange={(e) => {
                            if (selectedStrategy) {
                              const updatedStrategies = strategies.map(s => 
                                s.id === selectedStrategy.id 
                                  ? { ...s, uniqueIdReq: e.target.checked } 
                                  : s
                              );
                              setStrategies(updatedStrategies);
                              setSelectedStrategy({
                                ...selectedStrategy,
                                uniqueIdReq: e.target.checked
                              });
                            }
                          }}
                          disabled={selectedStrategy.enabled}
                        />
                      }
                      label="Require Unique ID for Order"
                      sx={{ mb: 1, display: 'block' }}
                    />
                  </Grid>
                  
                  <Grid item xs={12} md={6}>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={selectedStrategy.cancelPreviousSignal}
                          onChange={(e) => {
                            if (selectedStrategy) {
                              const updatedStrategies = strategies.map(s => 
                                s.id === selectedStrategy.id 
                                  ? { ...s, cancelPreviousSignal: e.target.checked } 
                                  : s
                              );
                              setStrategies(updatedStrategies);
                              setSelectedStrategy({
                                ...selectedStrategy,
                                cancelPreviousSignal: e.target.checked
                              });
                            }
                          }}
                          disabled={selectedStrategy.enabled}
                        />
                      }
                      label="Cancel Previous Open Signal"
                      sx={{ mb: 1, display: 'block' }}
                    />
                    
                    <TextField
                      fullWidth
                      label="Delay Between Users (seconds)"
                      type="number"
                      value={selectedStrategy.delayBetweenUsers}
                      onChange={(e) => {
                        if (selectedStrategy) {
                          const updatedStrategies = strategies.map(s => 
                            s.id === selectedStrategy.id 
                              ? { ...s, delayBetweenUsers: parseFloat(e.target.value) } 
                              : s
                          );
                          setStrategies(updatedStrategies);
                          setSelectedStrategy({
                            ...selectedStrategy,
                            delayBetweenUsers: parseFloat(e.target.value)
                          });
                        }
                      }}
                      sx={{ mb: 2 }}
                      disabled={selectedStrategy.enabled}
                    />
                    
                    <TextField
                      fullWidth
                      label="Square Off Time"
                      type="time"
                      value={selectedStrategy.squareOffTime}
                      onChange={(e) => {
                        if (selectedStrategy) {
                          const updatedStrategies = strategies.map(s => 
                            s.id === selectedStrategy.id 
                              ? { ...s, squareOffTime: e.target.value } 
                              : s
                          );
                          setStrategies(updatedStrategies);
                          setSelectedStrategy({
                            ...selectedStrategy,
                            squareOffTime: e.target.value
                          });
                        }
                      }}
                      sx={{ mb: 2 }}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      inputProps={{
                        step: 300, // 5 min
                      }}
                      disabled={selectedStrategy.enabled}
                    />
                  </Grid>
                </Grid>
              </Grid>
              
              <Grid item xs={12}>
                <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                  <Button 
                    variant="contained" 
                    startIcon={<SaveIcon />}
                    disabled={selectedStrategy.enabled}
                    onClick={handleEditStrategy}
                  >
                    Save Parameters
                  </Button>
                </Box>
              </Grid>
            </Grid>
          </TabPanel>
          
          {/* Tab 3: Backtest Results */}
          <TabPanel value={activeTab} index={2}>
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Backtest Results</Typography>
            
            <TableContainer sx={{ maxHeight: 300, mb: 2 }}>
              <Table stickyHeader size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Period</TableCell>
                    <TableCell>Net P&L</TableCell>
                    <TableCell>Win Rate</TableCell>
                    <TableCell>Profit Factor</TableCell>
                    <TableCell>Sharpe Ratio</TableCell>
                    <TableCell>Max Drawdown</TableCell>
                    <TableCell>Total Trades</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow>
                    <TableCell>Last 1 Month</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹15,750.00</TableCell>
                    <TableCell>68%</TableCell>
                    <TableCell>2.3</TableCell>
                    <TableCell>1.8</TableCell>
                    <TableCell>8.5%</TableCell>
                    <TableCell>42</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Last 3 Months</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹42,250.00</TableCell>
                    <TableCell>65%</TableCell>
                    <TableCell>2.1</TableCell>
                    <TableCell>1.7</TableCell>
                    <TableCell>12.3%</TableCell>
                    <TableCell>126</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Last 6 Months</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹78,500.00</TableCell>
                    <TableCell>62%</TableCell>
                    <TableCell>1.9</TableCell>
                    <TableCell>1.6</TableCell>
                    <TableCell>15.8%</TableCell>
                    <TableCell>245</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Last 1 Year</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹156,750.00</TableCell>
                    <TableCell>60%</TableCell>
                    <TableCell>1.8</TableCell>
                    <TableCell>1.5</TableCell>
                    <TableCell>18.2%</TableCell>
                    <TableCell>487</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
            
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Monthly Performance</Typography>
            
            <TableContainer sx={{ maxHeight: 300 }}>
              <Table stickyHeader size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Month</TableCell>
                    <TableCell>Net P&L</TableCell>
                    <TableCell>Win Rate</TableCell>
                    <TableCell>Total Trades</TableCell>
                    <TableCell>Avg. Trade</TableCell>
                    <TableCell>Max Profit</TableCell>
                    <TableCell>Max Loss</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow>
                    <TableCell>Apr 2025</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹4,250.00</TableCell>
                    <TableCell>70%</TableCell>
                    <TableCell>12</TableCell>
                    <TableCell>₹354.17</TableCell>
                    <TableCell>₹1,250.00</TableCell>
                    <TableCell>₹-450.00</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Mar 2025</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹5,750.00</TableCell>
                    <TableCell>65%</TableCell>
                    <TableCell>15</TableCell>
                    <TableCell>₹383.33</TableCell>
                    <TableCell>₹1,500.00</TableCell>
                    <TableCell>₹-600.00</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Feb 2025</TableCell>
                    <TableCell sx={{ color: 'error.main', fontWeight: 'bold' }}>₹-1,250.00</TableCell>
                    <TableCell>45%</TableCell>
                    <TableCell>14</TableCell>
                    <TableCell>₹-89.29</TableCell>
                    <TableCell>₹950.00</TableCell>
                    <TableCell>₹-1,200.00</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell>Jan 2025</TableCell>
                    <TableCell sx={{ color: 'success.main', fontWeight: 'bold' }}>₹7,000.00</TableCell>
                    <TableCell>72%</TableCell>
                    <TableCell>18</TableCell>
                    <TableCell>₹388.89</TableCell>
                    <TableCell>₹1,750.00</TableCell>
                    <TableCell>₹-500.00</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
          </TabPanel>
          
          {/* Tab 4: Performance */}
          <TabPanel value={activeTab} index={3}>
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Performance Metrics</Typography>
            
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Total P&L
                  </Typography>
                  <Typography 
                    variant="h6" 
                    sx={{ 
                      color: selectedStrategy.pnl > 0 ? 'success.main' : selectedStrategy.pnl < 0 ? 'error.main' : 'text.primary',
                      fontWeight: 'bold'
                    }}
                  >
                    ₹{selectedStrategy.pnl.toFixed(2)}
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Win Rate
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 'bold' }}>
                    65%
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Profit Factor
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 'bold' }}>
                    2.1
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Max Drawdown
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'error.main', fontWeight: 'bold' }}>
                    12.3%
                  </Typography>
                </Paper>
              </Grid>
            </Grid>
            
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Trade Statistics</Typography>
            
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Total Trades
                  </Typography>
                  <Typography variant="h6">
                    126
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Winning Trades
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'success.main' }}>
                    82
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Losing Trades
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'error.main' }}>
                    44
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Avg. Trade P&L
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'success.main' }}>
                    ₹335.32
                  </Typography>
                </Paper>
              </Grid>
            </Grid>
            
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Risk Metrics</Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Sharpe Ratio
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main' }}>
                    1.7
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Sortino Ratio
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main' }}>
                    2.3
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Calmar Ratio
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main' }}>
                    1.2
                  </Typography>
                </Paper>
              </Grid>
              
              <Grid item xs={6} md={3}>
                <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>
                  <Typography variant="body2" color="textSecondary" gutterBottom>
                    Recovery Factor
                  </Typography>
                  <Typography variant="h6" sx={{ color: 'primary.main' }}>
                    3.5
                  </Typography>
                </Paper>
              </Grid>
            </Grid>
          </TabPanel>
          
          {/* Tab 5: Logs */}
          <TabPanel value={activeTab} index={4}>
            <Typography variant="subtitle2" sx={{ mb: 2 }}>Strategy Execution Logs</Typography>
            
            <LogsPanel 
              category="strategy" 
              filter={selectedStrategy.name}
              height={400}
            />
          </TabPanel>
        </Paper>
      )}
      
      {/* Strategy Dialog */}
      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="md" fullWidth>
        <DialogTitle>{isEditing ? 'Edit Strategy' : 'Add New Strategy'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Strategy Name"
                value={newStrategy.name}
                onChange={(e) => setNewStrategy({ ...newStrategy, name: e.target.value })}
                sx={{ mb: 2 }}
              />
              
              <TextField
                fullWidth
                label="Description"
                value={newStrategy.description}
                onChange={(e) => setNewStrategy({ ...newStrategy, description: e.target.value })}
                multiline
                rows={3}
                sx={{ mb: 2 }}
              />
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Strategy Type</InputLabel>
                <Select
                  value={newStrategy.type}
                  label="Strategy Type"
                  onChange={(e) => setNewStrategy({ ...newStrategy, type: e.target.value })}
                >
                  <MenuItem value="Range Breakout">Range Breakout</MenuItem>
                  <MenuItem value="Mean Reversion">Mean Reversion</MenuItem>
                  <MenuItem value="Trend Following">Trend Following</MenuItem>
                  <MenuItem value="Momentum">Momentum</MenuItem>
                  <MenuItem value="Volatility">Volatility</MenuItem>
                  <MenuItem value="Arbitrage">Arbitrage</MenuItem>
                </Select>
              </FormControl>
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Instruments</InputLabel>
                <Select
                  multiple
                  value={newStrategy.instruments}
                  label="Instruments"
                  onChange={(e) => setNewStrategy({ ...newStrategy, instruments: e.target.value })}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {selected.map((value) => (
                        <Chip key={value} label={value} />
                      ))}
                    </Box>
                  )}
                >
                  <MenuItem value="NIFTY">NIFTY</MenuItem>
                  <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
                  <MenuItem value="FINNIFTY">FINNIFTY</MenuItem>
                  <MenuItem value="SENSEX">SENSEX</MenuItem>
                  <MenuItem value="RELIANCE">RELIANCE</MenuItem>
                  <MenuItem value="TCS">TCS</MenuItem>
                  <MenuItem value="INFY">INFY</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Status</InputLabel>
                <Select
                  value={newStrategy.status}
                  label="Status"
                  onChange={(e) => setNewStrategy({ ...newStrategy, status: e.target.value })}
                >
                  <MenuItem value="Allowed">Allowed</MenuItem>
                  <MenuItem value="Not Allowed">Not Allowed</MenuItem>
                </Select>
              </FormControl>
              
              <TextField
                fullWidth
                label="Square Off Time"
                type="time"
                value={newStrategy.squareOffTime}
                onChange={(e) => setNewStrategy({ ...newStrategy, squareOffTime: e.target.value })}
                sx={{ mb: 2 }}
                InputLabelProps={{
                  shrink: true,
                }}
              />
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>User Account</InputLabel>
                <Select
                  value={newStrategy.userAccount}
                  label="User Account"
                  onChange={(e) => setNewStrategy({ ...newStrategy, userAccount: e.target.value })}
                >
                  <MenuItem value="SIM1">SIM1</MenuItem>
                  <MenuItem value="SIM2">SIM2</MenuItem>
                  <MenuItem value="ZM9343">ZM9343</MenuItem>
                  <MenuItem value="DLT1182">DLT1182</MenuItem>
                  <MenuItem value="FA161611">FA161611</MenuItem>
                </Select>
              </FormControl>
              
              <Box sx={{ mb: 2 }}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={newStrategy.enabled}
                      onChange={(e) => setNewStrategy({ ...newStrategy, enabled: e.target.checked })}
                    />
                  }
                  label="Enable Strategy"
                />
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={newStrategy.manualSquareOff}
                      onChange={(e) => setNewStrategy({ ...newStrategy, manualSquareOff: e.target.checked })}
                    />
                  }
                  label="Manual Square Off"
                />
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={newStrategy.markAsCompleted}
                      onChange={(e) => setNewStrategy({ ...newStrategy, markAsCompleted: e.target.checked })}
                    />
                  }
                  label="Mark As Completed"
                />
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={newStrategy.uniqueIdReq}
                      onChange={(e) => setNewStrategy({ ...newStrategy, uniqueIdReq: e.target.checked })}
                    />
                  }
                  label="Unique ID Required"
                />
                
                <FormControlLabel
                  control={
                    <Switch
                      checked={newStrategy.cancelPreviousSignal}
                      onChange={(e) => setNewStrategy({ ...newStrategy, cancelPreviousSignal: e.target.checked })}
                    />
                  }
                  label="Cancel Previous Signal"
                />
              </Box>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handleSaveStrategy} variant="contained">Save</Button>
        </DialogActions>
      </Dialog>
      
      {/* Snackbar for notifications */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={handleCloseSnackbar} 
          severity={snackbar.severity} 
          sx={{ width: '100%' }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
      
      <Box sx={{ mt: 2 }}>
        <Typography variant="caption" color="textSecondary">
          1. If user's SqOff time hits before the Strategy Sq Off time, then all positions would be sqoff for the user and no new order will be placed.
        </Typography>
        <Typography variant="caption" color="textSecondary" display="block">
          2. Password and Pin is only required if you have selected for Auto Login. Auto login internally fills user details in browser for easy login. It is totally optional feature.
        </Typography>
        <Typography variant="caption" color="textSecondary" display="block">
          3. If you are facing Login issue with Zerodha, AliceBlue, Upstox then just Un-Tick the Auto Login and then proceed with Manual Login.
        </Typography>
      </Box>
    </Box>
  );
};

export default StrategiesComponent;
