import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  List, 
  ListItem, 
  ListItemText,
  Divider,
  Button,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Checkbox,
  Tooltip,
  Chip,
  Snackbar,
  Alert,
  Menu,
  Badge
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import SettingsIcon from '@mui/icons-material/Settings';
import RefreshIcon from '@mui/icons-material/Refresh';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import StopIcon from '@mui/icons-material/Stop';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import RestartAltIcon from '@mui/icons-material/RestartAlt';
import ShowChartIcon from '@mui/icons-material/ShowChart';
import ReplayIcon from '@mui/icons-material/Replay';
import SplitscreenIcon from '@mui/icons-material/Splitscreen';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import FilterListIcon from '@mui/icons-material/FilterList';
import PortfolioComponent from './PortfolioComponent';
import LogsPanel from '../logs/LogsPanel';

const MultiLegComponent = () => {
  const [selectedPortfolio, setSelectedPortfolio] = useState(null);
  const [portfolios, setPortfolios] = useState([
    { 
      id: 'portfolio1', 
      name: 'S8 1005-1100 CP35 3:1', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio2', 
      name: 'S8 1051 CP45 15/20/25', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio3', 
      name: 'SSNIF', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio4', 
      name: 'S10', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio5', 
      name: 'S4 935 CP45 15 LT 40%', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio6', 
      name: 'SENIF', 
      symbol: 'NIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio7', 
      name: '1051 CP120 35-40 15%', 
      symbol: 'BANKNIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio8', 
      name: 'MORB31 V8 9:16 TO 9:17', 
      symbol: 'BANKNIFTY',
      status: 'Completed',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: -131820.00,
      currentValue: 22100.00,
      valuePerLot: 110.50,
      underlyingPrice: 51625.15,
      underlyingLTP: 51680.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: true,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio9', 
      name: 'S 1', 
      symbol: 'BANKNIFTY',
      status: 'Monitoring',
      strategy: 'BACKENZOBUYING',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    },
    { 
      id: 'portfolio10', 
      name: 'NF_NDSTR', 
      symbol: 'NIFTY',
      status: 'UnderExecution',
      strategy: 'NF-NDSTR-D',
      enabled: true,
      pnl: -1612.50,
      currentValue: -1996.50,
      valuePerLot: -199.65,
      underlyingPrice: 23690.25,
      underlyingLTP: 23690.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    }
  ]);
  
  // Group portfolios by symbol
  const symbols = {};
  portfolios.forEach(portfolio => {
    if (!symbols[portfolio.symbol]) {
      symbols[portfolio.symbol] = [];
    }
    symbols[portfolio.symbol].push(portfolio);
  });
  
  const [expandedSymbols, setExpandedSymbols] = useState({
    NIFTY: true,
    BANKNIFTY: true,
    SENSEX: false
  });
  
  const [openDialog, setOpenDialog] = useState(false);
  const [newPortfolioName, setNewPortfolioName] = useState('');
  const [newPortfolioSymbol, setNewPortfolioSymbol] = useState('NIFTY');
  const [newPortfolioStrategy, setNewPortfolioStrategy] = useState('BACKENZOBUYING');
  
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success'
  });
  
  const [filterMenuAnchor, setFilterMenuAnchor] = useState(null);
  const [statusFilter, setStatusFilter] = useState('All');
  
  const toggleSymbol = (symbol) => {
    setExpandedSymbols({
      ...expandedSymbols,
      [symbol]: !expandedSymbols[symbol]
    });
  };
  
  const handlePortfolioSelect = (portfolio) => {
    setSelectedPortfolio(portfolio);
  };
  
  const handleAddPortfolio = () => {
    setOpenDialog(true);
    setNewPortfolioName('');
    setNewPortfolioSymbol('NIFTY');
    setNewPortfolioStrategy('BACKENZOBUYING');
  };
  
  const handleCloseDialog = () => {
    setOpenDialog(false);
  };
  
  const handleCreatePortfolio = () => {
    if (!newPortfolioName) {
      setSnackbar({
        open: true,
        message: 'Portfolio name is required',
        severity: 'error'
      });
      return;
    }
    
    const newPortfolio = {
      id: `portfolio${Date.now()}`,
      name: newPortfolioName,
      symbol: newPortfolioSymbol,
      strategy: newPortfolioStrategy,
      status: 'Monitoring',
      enabled: true,
      pnl: 0.00,
      currentValue: 0.00,
      valuePerLot: 0.00,
      underlyingPrice: 0.00,
      underlyingLTP: 0.00,
      executeSquareOff: false,
      edit: false,
      makeCopy: false,
      clone: false,
      delete: false,
      markAsCompleted: false,
      reset: false,
      payoff: false,
      chart: false,
      reexecute: false,
      partEntryExit: false
    };
    
    setPortfolios([...portfolios, newPortfolio]);
    
    // Ensure the symbol is expanded
    if (!expandedSymbols[newPortfolioSymbol]) {
      setExpandedSymbols({
        ...expandedSymbols,
        [newPortfolioSymbol]: true
      });
    }
    
    setOpenDialog(false);
    
    // Select the new portfolio
    setSelectedPortfolio(newPortfolio);
    
    setSnackbar({
      open: true,
      message: 'Portfolio created successfully',
      severity: 'success'
    });
  };
  
  const handleDeletePortfolio = (portfolioId) => {
    if (window.confirm('Are you sure you want to delete this portfolio?')) {
      const updatedPortfolios = portfolios.filter(p => p.id !== portfolioId);
      setPortfolios(updatedPortfolios);
      
      if (selectedPortfolio && selectedPortfolio.id === portfolioId) {
        setSelectedPortfolio(null);
      }
      
      setSnackbar({
        open: true,
        message: 'Portfolio deleted successfully',
        severity: 'success'
      });
    }
  };
  
  const handleToggleEnabled = (portfolioId) => {
    const updatedPortfolios = portfolios.map(p => 
      p.id === portfolioId ? { ...p, enabled: !p.enabled } : p
    );
    setPortfolios(updatedPortfolios);
    
    if (selectedPortfolio && selectedPortfolio.id === portfolioId) {
      setSelectedPortfolio({
        ...selectedPortfolio,
        enabled: !selectedPortfolio.enabled
      });
    }
  };
  
  const handleMarkAsCompleted = (portfolioId) => {
    const updatedPortfolios = portfolios.map(p => 
      p.id === portfolioId ? { 
        ...p, 
        markAsCompleted: !p.markAsCompleted,
        status: p.markAsCompleted ? 'Monitoring' : 'Completed'
      } : p
    );
    setPortfolios(updatedPortfolios);
    
    if (selectedPortfolio && selectedPortfolio.id === portfolioId) {
      setSelectedPortfolio({
        ...selectedPortfolio,
        markAsCompleted: !selectedPortfolio.markAsCompleted,
        status: selectedPortfolio.markAsCompleted ? 'Monitoring' : 'Completed'
      });
    }
  };
  
  const handleCloseSnackbar = () => {
    setSnackbar({
      ...snackbar,
      open: false
    });
  };
  
  const handleOpenFilterMenu = (event) => {
    setFilterMenuAnchor(event.currentTarget);
  };
  
  const handleCloseFilterMenu = () => {
    setFilterMenuAnchor(null);
  };
  
  const handleStatusFilterChange = (status) => {
    setStatusFilter(status);
    handleCloseFilterMenu();
  };
  
  const filteredPortfolios = statusFilter === 'All' 
    ? portfolios 
    : portfolios.filter(p => p.status === statusFilter);
  
  const getStatusColor = (status) => {
    switch (status) {
      case 'Monitoring':
        return '#4caf50';
      case 'Completed':
        return '#2196f3';
      case 'UnderExecution':
        return '#ff9800';
      case 'Disabled':
        return '#f44336';
      default:
        return 'inherit';
    }
  };
  
  const handleExecuteSquareOff = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Execute/Square Off action triggered',
      severity: 'info'
    });
  };
  
  const handleEdit = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Edit action triggered',
      severity: 'info'
    });
  };
  
  const handleMakeCopy = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Make Copy action triggered',
      severity: 'info'
    });
  };
  
  const handleClone = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Clone action triggered',
      severity: 'info'
    });
  };
  
  const handleReset = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Reset action triggered',
      severity: 'info'
    });
  };
  
  const handlePayoff = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Payoff action triggered',
      severity: 'info'
    });
  };
  
  const handleChart = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Chart action triggered',
      severity: 'info'
    });
  };
  
  const handleReexecute = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Reexecute action triggered',
      severity: 'info'
    });
  };
  
  const handlePartEntryExit = (portfolioId) => {
    setSnackbar({
      open: true,
      message: 'Part Entry/Exit action triggered',
      severity: 'info'
    });
  };
  
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">Multi-Leg Portfolio Management</Typography>
        <Box>
          <Button 
            variant="contained" 
            size="small" 
            startIcon={<AddIcon />}
            onClick={handleAddPortfolio}
            sx={{ mr: 1 }}
          >
            New Portfolio
          </Button>
          <IconButton size="small" onClick={handleOpenFilterMenu}>
            <Badge color="primary" variant="dot" invisible={statusFilter === 'All'}>
              <FilterListIcon />
            </Badge>
          </IconButton>
        </Box>
      </Box>
      
      <Menu
        anchorEl={filterMenuAnchor}
        open={Boolean(filterMenuAnchor)}
        onClose={handleCloseFilterMenu}
      >
        <MenuItem 
          onClick={() => handleStatusFilterChange('All')}
          selected={statusFilter === 'All'}
        >
          All Statuses
        </MenuItem>
        <MenuItem 
          onClick={() => handleStatusFilterChange('Monitoring')}
          selected={statusFilter === 'Monitoring'}
        >
          Monitoring
        </MenuItem>
        <MenuItem 
          onClick={() => handleStatusFilterChange('Completed')}
          selected={statusFilter === 'Completed'}
        >
          Completed
        </MenuItem>
        <MenuItem 
          onClick={() => handleStatusFilterChange('UnderExecution')}
          selected={statusFilter === 'UnderExecution'}
        >
          Under Execution
        </MenuItem>
        <MenuItem 
          onClick={() => handleStatusFilterChange('Disabled')}
          selected={statusFilter === 'Disabled'}
        >
          Disabled
        </MenuItem>
      </Menu>
      
      <Typography variant="caption" color="textSecondary" sx={{ mb: 1 }}>
        Drag a column to group
      </Typography>
      
      <Paper sx={{ flexGrow: 1, mb: 2, overflow: 'hidden' }}>
        <TableContainer sx={{ maxHeight: 'calc(100vh - 300px)' }}>
          <Table stickyHeader size="small">
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox">
                  <Typography variant="subtitle2">Enabled</Typography>
                </TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Portfolio Name</TableCell>
                <TableCell>Symbol</TableCell>
                <TableCell>Execute / SqOff</TableCell>
                <TableCell>Edit</TableCell>
                <TableCell>Make Copy</TableCell>
                <TableCell>Clone</TableCell>
                <TableCell>Delete</TableCell>
                <TableCell>Mark As Completed</TableCell>
                <TableCell>Strategy Tag</TableCell>
                <TableCell>Reset</TableCell>
                <TableCell>PayOff</TableCell>
                <TableCell>Chart</TableCell>
                <TableCell>Reexecute</TableCell>
                <TableCell>Part Entry / Exit</TableCell>
                <TableCell>PNL</TableCell>
                <TableCell>Current Value</TableCell>
                <TableCell>Value Per Lot</TableCell>
                <TableCell>Underlying Price</TableCell>
                <TableCell>Underlying LTP</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredPortfolios.map((portfolio) => (
                <TableRow 
                  key={portfolio.id}
                  hover
                  selected={selectedPortfolio && selectedPortfolio.id === portfolio.id}
                  onClick={() => handlePortfolioSelect(portfolio)}
                  sx={{ 
                    cursor: 'pointer',
                    bgcolor: portfolio.status === 'Completed' ? 'rgba(33, 150, 243, 0.05)' : 
                            portfolio.status === 'UnderExecution' ? 'rgba(255, 152, 0, 0.05)' :
                            portfolio.status === 'Disabled' ? 'rgba(244, 67, 54, 0.05)' :
                            'rgba(76, 175, 80, 0.05)'
                  }}
                >
                  <TableCell padding="checkbox">
                    <Checkbox
                      checked={portfolio.enabled}
                      onChange={() => handleToggleEnabled(portfolio.id)}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </TableCell>
                  <TableCell>
                    <Typography 
                      variant="body2" 
                      sx={{ 
                        color: getStatusColor(portfolio.status),
                        fontWeight: 'medium'
                      }}
                    >
                      {portfolio.status}
                    </Typography>
                  </TableCell>
                  <TableCell>{portfolio.name}</TableCell>
                  <TableCell>{portfolio.symbol}</TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleExecuteSquareOff(portfolio.id);
                      }}
                    >
                      <PlayArrowIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleEdit(portfolio.id);
                      }}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleMakeCopy(portfolio.id);
                      }}
                    >
                      <ContentCopyIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleClone(portfolio.id);
                      }}
                    >
                      <ContentCopyIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="error"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeletePortfolio(portfolio.id);
                      }}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <Checkbox
                      checked={portfolio.markAsCompleted || portfolio.status === 'Completed'}
                      onChange={() => handleMarkAsCompleted(portfolio.id)}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={portfolio.strategy} 
                      size="small" 
                      color="primary" 
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleReset(portfolio.id);
                      }}
                    >
                      <RestartAltIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handlePayoff(portfolio.id);
                      }}
                    >
                      <AttachMoneyIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleChart(portfolio.id);
                      }}
                    >
                      <ShowChartIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleReexecute(portfolio.id);
                      }}
                    >
                      <ReplayIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={(e) => {
                        e.stopPropagation();
                        handlePartEntryExit(portfolio.id);
                      }}
                    >
                      <SplitscreenIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                  <TableCell sx={{ 
                    color: portfolio.pnl > 0 ? 'success.main' : portfolio.pnl < 0 ? 'error.main' : 'text.primary',
                    fontWeight: portfolio.pnl !== 0 ? 'bold' : 'normal'
                  }}>
                    {portfolio.pnl.toFixed(2)}
                  </TableCell>
                  <TableCell>{portfolio.currentValue.toFixed(2)}</TableCell>
                  <TableCell>{portfolio.valuePerLot.toFixed(2)}</TableCell>
                  <TableCell>{portfolio.underlyingPrice.toFixed(2)}</TableCell>
                  <TableCell>{portfolio.underlyingLTP.toFixed(2)}</TableCell>
                </TableRow>
              ))}
              {filteredPortfolios.length === 0 && (
                <TableRow>
                  <TableCell colSpan={21} align="center" sx={{ py: 3 }}>
                    <Typography variant="body2" color="textSecondary">
                      No portfolios found
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
      
      {/* Main content area for portfolio details */}
      <Grid container spacing={2} sx={{ flexGrow: 1 }}>
        <Grid item xs={12}>
          {selectedPortfolio ? (
            <PortfolioComponent portfolio={selectedPortfolio} />
          ) : (
            <Paper sx={{ p: 4, height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center', flexDirection: 'column' }}>
              <Typography variant="h6" color="textSecondary" sx={{ mb: 2 }}>
                No Portfolio Selected
              </Typography>
              <Typography variant="body2" color="textSecondary" sx={{ mb: 3, textAlign: 'center' }}>
                Please select a portfolio from the table above or create a new one.
              </Typography>
              <Button 
                variant="contained" 
                startIcon={<AddIcon />}
                onClick={handleAddPortfolio}
              >
                Create New Portfolio
              </Button>
            </Paper>
          )}
        </Grid>
      </Grid>
      
      {/* Add Portfolio Dialog */}
      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>Create New Portfolio</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Portfolio Name"
            fullWidth
            value={newPortfolioName}
            onChange={(e) => setNewPortfolioName(e.target.value)}
            sx={{ mb: 2, mt: 1 }}
          />
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Symbol</InputLabel>
            <Select
              value={newPortfolioSymbol}
              label="Symbol"
              onChange={(e) => setNewPortfolioSymbol(e.target.value)}
            >
              <MenuItem value="NIFTY">NIFTY</MenuItem>
              <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
              <MenuItem value="FINNIFTY">FINNIFTY</MenuItem>
              <MenuItem value="SENSEX">SENSEX</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel>Strategy</InputLabel>
            <Select
              value={newPortfolioStrategy}
              label="Strategy"
              onChange={(e) => setNewPortfolioStrategy(e.target.value)}
            >
              <MenuItem value="BACKENZOBUYING">BACKENZOBUYING</MenuItem>
              <MenuItem value="NF-NDSTR-D">NF-NDSTR-D</MenuItem>
              <MenuItem value="MARU">MARU</MenuItem>
              <MenuItem value="LUX24VR">LUX24VR</MenuItem>
              <MenuItem value="NFTTR">NFTTR</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handleCreatePortfolio} variant="contained">Create</Button>
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
    </Box>
  );
};

export default MultiLegComponent;
