import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Tabs, 
  Tab, 
  Accordion,
  AccordionSummary,
  AccordionDetails,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Button,
  Grid,
  Divider,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Switch,
  FormControlLabel,
  Chip
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FolderIcon from '@mui/icons-material/Folder';
import BarChartIcon from '@mui/icons-material/BarChart';
import { strategyService } from '../../services/strategyService';

// Mock data for initial development
const mockStrategies = [
  {
    id: 'strat-001',
    name: 'MARU',
    description: 'Mean reversion strategy for index futures',
    status: 'ACTIVE',
    portfolios: [
      {
        id: 'port-001',
        name: 'Portfolio 1',
        description: 'Nifty futures with 5-minute timeframe',
        instruments: ['NIFTY 50'],
        status: 'ACTIVE'
      },
      {
        id: 'port-002',
        name: 'Portfolio 2',
        description: 'BankNifty futures with 5-minute timeframe',
        instruments: ['BANKNIFTY'],
        status: 'ACTIVE'
      }
    ]
  },
  {
    id: 'strat-002',
    name: 'LUX24VR',
    description: 'Volatility breakout strategy for equity futures',
    status: 'ACTIVE',
    portfolios: [
      {
        id: 'port-003',
        name: 'Portfolio 1',
        description: 'Large cap stocks with 15-minute timeframe',
        instruments: ['RELIANCE', 'HDFCBANK', 'TCS'],
        status: 'ACTIVE'
      }
    ]
  },
  {
    id: 'strat-003',
    name: 'NFTTR',
    description: 'Trend following strategy for commodity futures',
    status: 'INACTIVE',
    portfolios: [
      {
        id: 'port-004',
        name: 'Portfolio 1',
        description: 'Gold and Silver futures with 30-minute timeframe',
        instruments: ['GOLD', 'SILVER'],
        status: 'INACTIVE'
      },
      {
        id: 'port-005',
        name: 'Portfolio 2',
        description: 'Crude Oil futures with 30-minute timeframe',
        instruments: ['CRUDEOIL'],
        status: 'INACTIVE'
      }
    ]
  }
];

const StrategyPanel = () => {
  const [strategies, setStrategies] = useState(mockStrategies);
  const [selectedStrategy, setSelectedStrategy] = useState(null);
  const [selectedPortfolio, setSelectedPortfolio] = useState(null);
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(false);
  
  // Portfolio configuration settings
  const [portfolioConfig, setPortfolioConfig] = useState({
    // Execution Parameters
    executionMode: 'LIVE',
    orderType: 'LIMIT',
    slippagePercent: 0.1,
    maxOrdersPerDay: 10,
    
    // Range Breakout
    rangeBreakoutEnabled: true,
    rangeBreakoutPeriod: 20,
    rangeBreakoutMultiplier: 1.5,
    
    // Extra Conditions
    volumeFilterEnabled: true,
    volumeThreshold: 1.5,
    rsiFilterEnabled: true,
    rsiLowerThreshold: 30,
    rsiUpperThreshold: 70,
    
    // Other Settings
    maxPositions: 5,
    maxRiskPerTrade: 1.0,
    maxDrawdownPercent: 5.0,
    
    // Monitoring
    emailAlerts: true,
    smsAlerts: false,
    autoShutdown: true,
    maxLossAmount: 10000,
    
    // Dynamic Hedge
    dynamicHedgeEnabled: false,
    hedgeInstrument: 'NIFTY 50',
    hedgeRatio: 0.5,
    
    // Target Settings
    targetType: 'PERCENT',
    targetValue: 1.5,
    trailingStopEnabled: true,
    trailingStopPercent: 0.5,
    
    // Stoploss Settings
    stopLossType: 'PERCENT',
    stopLossValue: 1.0,
    timeBasedStopLossEnabled: true,
    timeBasedStopLossMinutes: 60,
    
    // Exit Settings
    exitTimeEnabled: true,
    exitTime: '15:15',
    partialExitEnabled: true,
    partialExitPercent: 50,
    
    // At Broker
    brokerName: 'Zerodha',
    exchangeSegment: 'NSE',
    productType: 'MIS',
    disclosedQuantity: 0
  });

  // Fetch strategies from API when component mounts
  useEffect(() => {
    const fetchStrategies = async () => {
      try {
        setLoading(true);
        // In a real implementation, this would call the API
        // const response = await strategyService.getStrategies();
        // setStrategies(response.data);
        
        // Using mock data for now
        setStrategies(mockStrategies);
      } catch (error) {
        console.error('Error fetching strategies:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStrategies();
  }, []);

  // Set default selected strategy and portfolio when data is loaded
  useEffect(() => {
    if (strategies.length > 0 && !selectedStrategy) {
      setSelectedStrategy(strategies[0]);
      if (strategies[0].portfolios.length > 0) {
        setSelectedPortfolio(strategies[0].portfolios[0]);
      }
    }
  }, [strategies, selectedStrategy]);

  // Fetch portfolio configuration when a portfolio is selected
  useEffect(() => {
    if (selectedPortfolio) {
      const fetchPortfolioConfig = async () => {
        try {
          setLoading(true);
          // In a real implementation, this would call the API
          // const response = await strategyService.getPortfolioConfig(selectedPortfolio.id);
          // setPortfolioConfig(response.data);
          
          // Using mock data for now
          // In a real implementation, we would fetch the actual config
          setLoading(false);
        } catch (error) {
          console.error('Error fetching portfolio configuration:', error);
          setLoading(false);
        }
      };

      fetchPortfolioConfig();
    }
  }, [selectedPortfolio]);

  const handleStrategySelect = (strategy) => {
    setSelectedStrategy(strategy);
    if (strategy.portfolios.length > 0) {
      setSelectedPortfolio(strategy.portfolios[0]);
    } else {
      setSelectedPortfolio(null);
    }
    setTabValue(0);
  };

  const handlePortfolioSelect = (portfolio) => {
    setSelectedPortfolio(portfolio);
    setTabValue(0);
  };

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleConfigChange = (section, field, value) => {
    setPortfolioConfig(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleSaveConfig = async () => {
    try {
      setLoading(true);
      // In a real implementation, this would call the API
      // await strategyService.updatePortfolioConfig(selectedPortfolio.id, portfolioConfig);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500));
      
      console.log('Portfolio configuration saved:', portfolioConfig);
      setLoading(false);
    } catch (error) {
      console.error('Error saving portfolio configuration:', error);
      setLoading(false);
    }
  };

  const handleStartPortfolio = async () => {
    try {
      // In a real implementation, this would call the API
      // await strategyService.startPortfolio(selectedPortfolio.id);
      
      // Update local state
      setStrategies(prev => 
        prev.map(strategy => 
          strategy.id === selectedStrategy.id
            ? {
                ...strategy,
                portfolios: strategy.portfolios.map(portfolio => 
                  portfolio.id === selectedPortfolio.id
                    ? { ...portfolio, status: 'ACTIVE' }
                    : portfolio
                )
              }
            : strategy
        )
      );
      
      // Update selected portfolio
      setSelectedPortfolio(prev => ({ ...prev, status: 'ACTIVE' }));
    } catch (error) {
      console.error('Error starting portfolio:', error);
    }
  };

  const handleStopPortfolio = async () => {
    try {
      // In a real implementation, this would call the API
      // await strategyService.stopPortfolio(selectedPortfolio.id);
      
      // Update local state
      setStrategies(prev => 
        prev.map(strategy => 
          strategy.id === selectedStrategy.id
            ? {
                ...strategy,
                portfolios: strategy.portfolios.map(portfolio => 
                  portfolio.id === selectedPortfolio.id
                    ? { ...portfolio, status: 'INACTIVE' }
                    : portfolio
                )
              }
            : strategy
        )
      );
      
      // Update selected portfolio
      setSelectedPortfolio(prev => ({ ...prev, status: 'INACTIVE' }));
    } catch (error) {
      console.error('Error stopping portfolio:', error);
    }
  };

  // Strategy list sidebar
  const renderStrategySidebar = () => (
    <Paper elevation={2} sx={{ width: '100%', height: '100%' }}>
      <Typography variant="h6" sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
        Strategies
      </Typography>
      
      <List component="nav">
        {strategies.map((strategy) => (
          <React.Fragment key={strategy.id}>
            <ListItem 
              button 
              onClick={() => handleStrategySelect(strategy)}
              selected={selectedStrategy && selectedStrategy.id === strategy.id}
            >
              <ListItemIcon>
                <BarChartIcon />
              </ListItemIcon>
              <ListItemText 
                primary={strategy.name} 
                secondary={strategy.description}
              />
              <Chip 
                label={strategy.status} 
                size="small" 
                color={strategy.status === 'ACTIVE' ? 'success' : 'default'} 
              />
            </ListItem>
            
            {selectedStrategy && selectedStrategy.id === strategy.id && (
              <List component="div" disablePadding>
                {strategy.portfolios.map((portfolio) => (
                  <ListItem 
                    key={portfolio.id}
                    button 
                    onClick={() => handlePortfolioSelect(portfolio)}
                    selected={selectedPortfolio && selectedPortfolio.id === portfolio.id}
                    sx={{ pl: 4 }}
                  >
                    <ListItemIcon>
                      <FolderIcon />
                    </ListItemIcon>
                    <ListItemText 
                      primary={portfolio.name} 
                      secondary={portfolio.description}
                    />
                    <Chip 
                      label={portfolio.status} 
                      size="small" 
                      color={portfolio.status === 'ACTIVE' ? 'success' : 'default'} 
                    />
                  </ListItem>
                ))}
              </List>
            )}
          </React.Fragment>
        ))}
      </List>
    </Paper>
  );

  // Portfolio configuration tabs
  const renderPortfolioTabs = () => (
    <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
      <Tabs 
        value={tabValue} 
        onChange={handleTabChange} 
        variant="scrollable"
        scrollButtons="auto"
        aria-label="portfolio configuration tabs"
      >
        <Tab label="Execution Parameters" />
        <Tab label="Range Breakout" />
        <Tab label="Extra Conditions" />
        <Tab label="Other Settings" />
        <Tab label="Monitoring" />
        <Tab label="Dynamic Hedge" />
        <Tab label="Target Settings" />
        <Tab label="Stoploss Settings" />
        <Tab label="Exit Settings" />
        <Tab label="At Broker" />
      </Tabs>
    </Box>
  );

  // Execution Parameters tab
  const renderExecutionParameters = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Execution Parameters
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Execution Mode</InputLabel>
            <Select
              value={portfolioConfig.executionMode}
              label="Execution Mode"
              onChange={(e) => handleConfigChange('execution', 'executionMode', e.target.value)}
            >
              <MenuItem value="LIVE">Live Trading</MenuItem>
              <MenuItem value="PAPER">Paper Trading</MenuItem>
              <MenuItem value="BACKTEST">Backtest</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Order Type</InputLabel>
            <Select
              value={portfolioConfig.orderType}
              label="Order Type"
              onChange={(e) => handleConfigChange('execution', 'orderType', e.target.value)}
            >
              <MenuItem value="MARKET">Market</MenuItem>
              <MenuItem value="LIMIT">Limit</MenuItem>
              <MenuItem value="SL">Stop Loss</MenuItem>
              <MenuItem value="SL-M">Stop Loss Market</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Slippage Percent"
            type="number"
            value={portfolioConfig.slippagePercent}
            onChange={(e) => handleConfigChange('execution', 'slippagePercent', parseFloat(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0, step: 0.1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Max Orders Per Day"
            type="number"
            value={portfolioConfig.maxOrdersPerDay}
            onChange={(e) => handleConfigChange('execution', 'maxOrdersPerDay', parseInt(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Range Breakout tab
  const renderRangeBreakout = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Range Breakout Settings
      </Typography>
      
      <FormControlLabel
        control={
          <Switch
            checked={portfolioConfig.rangeBreakoutEnabled}
            onChange={(e) => handleConfigChange('rangeBreakout', 'rangeBreakoutEnabled', e.target.checked)}
            color="primary"
          />
        }
        label="Enable Range Breakout"
        sx={{ mb: 2 }}
      />
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Range Period (Bars)"
            type="number"
            value={portfolioConfig.rangeBreakoutPeriod}
            onChange={(e) => handleConfigChange('rangeBreakout', 'rangeBreakoutPeriod', parseInt(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.rangeBreakoutEnabled}
            InputProps={{
              inputProps: { min: 1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Breakout Multiplier"
            type="number"
            value={portfolioConfig.rangeBreakoutMultiplier}
            onChange={(e) => handleConfigChange('rangeBreakout', 'rangeBreakoutMultiplier', parseFloat(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.rangeBreakoutEnabled}
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Extra Conditions tab
  const renderExtraConditions = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Extra Conditions
      </Typography>
      
      <Accordion defaultExpanded>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>Volume Filter</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.volumeFilterEnabled}
                onChange={(e) => handleConfigChange('extraConditions', 'volumeFilterEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable Volume Filter"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Volume Threshold (x Average)"
            type="number"
            value={portfolioConfig.volumeThreshold}
            onChange={(e) => handleConfigChange('extraConditions', 'volumeThreshold', parseFloat(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.volumeFilterEnabled}
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </AccordionDetails>
      </Accordion>
      
      <Accordion defaultExpanded>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>RSI Filter</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.rsiFilterEnabled}
                onChange={(e) => handleConfigChange('extraConditions', 'rsiFilterEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable RSI Filter"
            sx={{ mb: 2 }}
          />
          
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="RSI Lower Threshold"
                type="number"
                value={portfolioConfig.rsiLowerThreshold}
                onChange={(e) => handleConfigChange('extraConditions', 'rsiLowerThreshold', parseInt(e.target.value))}
                margin="normal"
                disabled={!portfolioConfig.rsiFilterEnabled}
                InputProps={{
                  inputProps: { min: 0, max: 100 }
                }}
              />
            </Grid>
            
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="RSI Upper Threshold"
                type="number"
                value={portfolioConfig.rsiUpperThreshold}
                onChange={(e) => handleConfigChange('extraConditions', 'rsiUpperThreshold', parseInt(e.target.value))}
                margin="normal"
                disabled={!portfolioConfig.rsiFilterEnabled}
                InputProps={{
                  inputProps: { min: 0, max: 100 }
                }}
              />
            </Grid>
          </Grid>
        </AccordionDetails>
      </Accordion>
    </Box>
  );

  // Other Settings tab
  const renderOtherSettings = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Other Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Max Positions"
            type="number"
            value={portfolioConfig.maxPositions}
            onChange={(e) => handleConfigChange('otherSettings', 'maxPositions', parseInt(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Max Risk Per Trade (%)"
            type="number"
            value={portfolioConfig.maxRiskPerTrade}
            onChange={(e) => handleConfigChange('otherSettings', 'maxRiskPerTrade', parseFloat(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Max Drawdown (%)"
            type="number"
            value={portfolioConfig.maxDrawdownPercent}
            onChange={(e) => handleConfigChange('otherSettings', 'maxDrawdownPercent', parseFloat(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Monitoring tab
  const renderMonitoring = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Monitoring Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.emailAlerts}
                onChange={(e) => handleConfigChange('monitoring', 'emailAlerts', e.target.checked)}
                color="primary"
              />
            }
            label="Email Alerts"
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.smsAlerts}
                onChange={(e) => handleConfigChange('monitoring', 'smsAlerts', e.target.checked)}
                color="primary"
              />
            }
            label="SMS Alerts"
          />
        </Grid>
        
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.autoShutdown}
                onChange={(e) => handleConfigChange('monitoring', 'autoShutdown', e.target.checked)}
                color="primary"
              />
            }
            label="Auto Shutdown on Max Loss"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Max Loss Amount"
            type="number"
            value={portfolioConfig.maxLossAmount}
            onChange={(e) => handleConfigChange('monitoring', 'maxLossAmount', parseInt(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.autoShutdown}
            InputProps={{
              inputProps: { min: 1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Dynamic Hedge tab
  const renderDynamicHedge = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Dynamic Hedge Settings
      </Typography>
      
      <FormControlLabel
        control={
          <Switch
            checked={portfolioConfig.dynamicHedgeEnabled}
            onChange={(e) => handleConfigChange('dynamicHedge', 'dynamicHedgeEnabled', e.target.checked)}
            color="primary"
          />
        }
        label="Enable Dynamic Hedging"
        sx={{ mb: 2 }}
      />
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal" disabled={!portfolioConfig.dynamicHedgeEnabled}>
            <InputLabel>Hedge Instrument</InputLabel>
            <Select
              value={portfolioConfig.hedgeInstrument}
              label="Hedge Instrument"
              onChange={(e) => handleConfigChange('dynamicHedge', 'hedgeInstrument', e.target.value)}
            >
              <MenuItem value="NIFTY 50">NIFTY 50</MenuItem>
              <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
              <MenuItem value="FINNIFTY">FINNIFTY</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Hedge Ratio"
            type="number"
            value={portfolioConfig.hedgeRatio}
            onChange={(e) => handleConfigChange('dynamicHedge', 'hedgeRatio', parseFloat(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.dynamicHedgeEnabled}
            InputProps={{
              inputProps: { min: 0.1, max: 1.0, step: 0.1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Target Settings tab
  const renderTargetSettings = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Target Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Target Type</InputLabel>
            <Select
              value={portfolioConfig.targetType}
              label="Target Type"
              onChange={(e) => handleConfigChange('targetSettings', 'targetType', e.target.value)}
            >
              <MenuItem value="PERCENT">Percent</MenuItem>
              <MenuItem value="POINTS">Points</MenuItem>
              <MenuItem value="AMOUNT">Amount</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Target Value"
            type="number"
            value={portfolioConfig.targetValue}
            onChange={(e) => handleConfigChange('targetSettings', 'targetValue', parseFloat(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.trailingStopEnabled}
                onChange={(e) => handleConfigChange('targetSettings', 'trailingStopEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable Trailing Stop"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Trailing Stop Percent"
            type="number"
            value={portfolioConfig.trailingStopPercent}
            onChange={(e) => handleConfigChange('targetSettings', 'trailingStopPercent', parseFloat(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.trailingStopEnabled}
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Stoploss Settings tab
  const renderStoplossSettings = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Stoploss Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Stoploss Type</InputLabel>
            <Select
              value={portfolioConfig.stopLossType}
              label="Stoploss Type"
              onChange={(e) => handleConfigChange('stoplossSettings', 'stopLossType', e.target.value)}
            >
              <MenuItem value="PERCENT">Percent</MenuItem>
              <MenuItem value="POINTS">Points</MenuItem>
              <MenuItem value="AMOUNT">Amount</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Stoploss Value"
            type="number"
            value={portfolioConfig.stopLossValue}
            onChange={(e) => handleConfigChange('stoplossSettings', 'stopLossValue', parseFloat(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0.1, step: 0.1 }
            }}
          />
        </Grid>
        
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.timeBasedStopLossEnabled}
                onChange={(e) => handleConfigChange('stoplossSettings', 'timeBasedStopLossEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable Time-Based Stoploss"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Time-Based Stoploss (Minutes)"
            type="number"
            value={portfolioConfig.timeBasedStopLossMinutes}
            onChange={(e) => handleConfigChange('stoplossSettings', 'timeBasedStopLossMinutes', parseInt(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.timeBasedStopLossEnabled}
            InputProps={{
              inputProps: { min: 1 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Exit Settings tab
  const renderExitSettings = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Exit Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.exitTimeEnabled}
                onChange={(e) => handleConfigChange('exitSettings', 'exitTimeEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable Exit Time"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Exit Time"
            type="time"
            value={portfolioConfig.exitTime}
            onChange={(e) => handleConfigChange('exitSettings', 'exitTime', e.target.value)}
            margin="normal"
            disabled={!portfolioConfig.exitTimeEnabled}
            InputLabelProps={{
              shrink: true,
            }}
            inputProps={{
              step: 300, // 5 min
            }}
          />
        </Grid>
        
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Switch
                checked={portfolioConfig.partialExitEnabled}
                onChange={(e) => handleConfigChange('exitSettings', 'partialExitEnabled', e.target.checked)}
                color="primary"
              />
            }
            label="Enable Partial Exit"
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            label="Partial Exit Percent"
            type="number"
            value={portfolioConfig.partialExitPercent}
            onChange={(e) => handleConfigChange('exitSettings', 'partialExitPercent', parseInt(e.target.value))}
            margin="normal"
            disabled={!portfolioConfig.partialExitEnabled}
            InputProps={{
              inputProps: { min: 1, max: 99 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // At Broker tab
  const renderAtBroker = () => (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        At Broker Settings
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Broker</InputLabel>
            <Select
              value={portfolioConfig.brokerName}
              label="Broker"
              onChange={(e) => handleConfigChange('atBroker', 'brokerName', e.target.value)}
            >
              <MenuItem value="Zerodha">Zerodha</MenuItem>
              <MenuItem value="XTS">XTS</MenuItem>
              <MenuItem value="SIM1">Simulator 1</MenuItem>
              <MenuItem value="SIM2">Simulator 2</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Exchange Segment</InputLabel>
            <Select
              value={portfolioConfig.exchangeSegment}
              label="Exchange Segment"
              onChange={(e) => handleConfigChange('atBroker', 'exchangeSegment', e.target.value)}
            >
              <MenuItem value="NSE">NSE</MenuItem>
              <MenuItem value="BSE">BSE</MenuItem>
              <MenuItem value="NFO">NFO</MenuItem>
              <MenuItem value="MCX">MCX</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Product Type</InputLabel>
            <Select
              value={portfolioConfig.productType}
              label="Product Type"
              onChange={(e) => handleConfigChange('atBroker', 'productType', e.target.value)}
            >
              <MenuItem value="MIS">MIS (Intraday)</MenuItem>
              <MenuItem value="NRML">NRML (Carry Forward)</MenuItem>
              <MenuItem value="CNC">CNC (Delivery)</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Disclosed Quantity"
            type="number"
            value={portfolioConfig.disclosedQuantity}
            onChange={(e) => handleConfigChange('atBroker', 'disclosedQuantity', parseInt(e.target.value))}
            margin="normal"
            InputProps={{
              inputProps: { min: 0 }
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );

  // Render the appropriate tab content
  const renderTabContent = () => {
    switch (tabValue) {
      case 0:
        return renderExecutionParameters();
      case 1:
        return renderRangeBreakout();
      case 2:
        return renderExtraConditions();
      case 3:
        return renderOtherSettings();
      case 4:
        return renderMonitoring();
      case 5:
        return renderDynamicHedge();
      case 6:
        return renderTargetSettings();
      case 7:
        return renderStoplossSettings();
      case 8:
        return renderExitSettings();
      case 9:
        return renderAtBroker();
      default:
        return renderExecutionParameters();
    }
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Grid container spacing={2}>
        <Grid item xs={12} md={3}>
          {renderStrategySidebar()}
        </Grid>
        
        <Grid item xs={12} md={9}>
          {selectedPortfolio ? (
            <Paper elevation={2}>
              <Box sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: 1, borderColor: 'divider' }}>
                <Box>
                  <Typography variant="h6">
                    {selectedStrategy.name} / {selectedPortfolio.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {selectedPortfolio.description}
                  </Typography>
                </Box>
                
                <Box>
                  {selectedPortfolio.status === 'ACTIVE' ? (
                    <Button 
                      variant="outlined" 
                      color="error" 
                      onClick={handleStopPortfolio}
                      disabled={loading}
                    >
                      Stop Portfolio
                    </Button>
                  ) : (
                    <Button 
                      variant="contained" 
                      color="success" 
                      onClick={handleStartPortfolio}
                      disabled={loading}
                    >
                      Start Portfolio
                    </Button>
                  )}
                </Box>
              </Box>
              
              {renderPortfolioTabs()}
              {renderTabContent()}
              
              <Box sx={{ p: 2, display: 'flex', justifyContent: 'flex-end', borderTop: 1, borderColor: 'divider' }}>
                <Button 
                  variant="contained" 
                  color="primary" 
                  onClick={handleSaveConfig}
                  disabled={loading}
                >
                  Save Configuration
                </Button>
              </Box>
            </Paper>
          ) : (
            <Paper elevation={2} sx={{ p: 3, textAlign: 'center' }}>
              <Typography variant="h6">
                Select a portfolio to configure
              </Typography>
            </Paper>
          )}
        </Grid>
      </Grid>
    </Box>
  );
};

export default StrategyPanel;
