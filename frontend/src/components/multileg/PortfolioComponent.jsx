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
  Checkbox,
  FormControlLabel,
  IconButton,
  Link,
  Tooltip,
  InputAdornment,
  Card,
  CardContent,
  Snackbar,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import AddIcon from '@mui/icons-material/Add';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import SaveIcon from '@mui/icons-material/Save';
import InfoIcon from '@mui/icons-material/Info';
import Panel3Tabs from './Panel3Tabs';

// Portfolio Component with 4-panel layout
const PortfolioComponent = ({ portfolio }) => {
  const [activeTab, setActiveTab] = useState(0);
  
  // Panel 1: Default Portfolio Settings state
  const [exchangeValue, setExchangeValue] = useState('NSE');
  const [symbolValue, setSymbolValue] = useState('NIFTY');
  const [niftyLiveValue, setNiftyLiveValue] = useState('NIFTY 50: 23250.10');
  const [expiryValue, setExpiryValue] = useState('09-Apr-25');
  const [defaultLots, setDefaultLots] = useState('10');
  const [lotSize, setLotSize] = useState('75');
  const [payPremium, setPayPremium] = useState('116100.0');
  
  // Strategy Configuration state
  const [predefinedStrategy, setPredefinedStrategy] = useState('Custom');
  const [strikeSelection, setStrikeSelection] = useState('Relative');
  const [underlying, setUnderlying] = useState('Spot');
  const [priceType, setPriceType] = useState('LTP');
  const [strikeStep, setStrikeStep] = useState('50');
  
  // Portfolio Behavior Settings state
  const [positionalPortfolio, setPositionalPortfolio] = useState(false);
  const [buyTradesFirst, setBuyTradesFirst] = useState(true);
  const [allowFarStrikes, setAllowFarStrikes] = useState(false);
  
  // Value Calculation state
  const [valueAllLots, setValueAllLots] = useState('₹116,100.00');
  const [valuePerLot, setValuePerLot] = useState('₹11,610.00');
  const [premiumGreekLog, setPremiumGreekLog] = useState(false);
  
  // Event handlers
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handleExchangeChange = (event) => {
    setExchangeValue(event.target.value);
  };
  
  const handleSymbolChange = (event) => {
    setSymbolValue(event.target.value);
  };
  
  const handleExpiryChange = (event) => {
    setExpiryValue(event.target.value);
  };
  
  const handleDefaultLotsChange = (event) => {
    setDefaultLots(event.target.value);
  };
  
  const handleLotSizeChange = (event) => {
    setLotSize(event.target.value);
  };
  
  const handlePayPremiumChange = (event) => {
    setPayPremium(event.target.value);
  };
  
  const handlePredefinedStrategyChange = (event) => {
    setPredefinedStrategy(event.target.value);
  };
  
  const handleStrikeSelectionChange = (event) => {
    setStrikeSelection(event.target.value);
  };
  
  const handleUnderlyingChange = (event) => {
    setUnderlying(event.target.value);
  };
  
  const handlePriceTypeChange = (event) => {
    setPriceType(event.target.value);
  };
  
  const handleStrikeStepChange = (event) => {
    setStrikeStep(event.target.value);
  };
  
  const handlePositionalPortfolioChange = (event) => {
    setPositionalPortfolio(event.target.checked);
  };
  
  const handleBuyTradesFirstChange = (event) => {
    setBuyTradesFirst(event.target.checked);
  };
  
  const handleAllowFarStrikesChange = (event) => {
    setAllowFarStrikes(event.target.checked);
  };
  
  const handlePremiumGreekLogChange = (event) => {
    setPremiumGreekLog(event.target.checked);
  };
  
  const handleRefreshClick = () => {
    console.log('Refresh clicked');
  };
  
  const handleAddLegClick = () => {
    console.log('Add Leg clicked');
  };
  
  return (
    <Box sx={{ width: '100%', height: '100%' }}>
      <Paper sx={{ width: '100%', height: '100%', p: 1, borderRadius: 1 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
          <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
            <Grid container spacing={1}>
              {/* Panel 1: Default Portfolio Settings */}
              <Grid item xs={12} md={4}>
                <Card variant="outlined" sx={{ height: '100%', bgcolor: '#fff' }}>
                  <CardContent sx={{ p: 0.5, '&:last-child': { pb: 0.5 } }}>
                    <Typography variant="caption" sx={{ mb: 0.3, color: '#2c3e50', fontWeight: 'bold', display: 'block', fontSize: '0.65rem' }}>
                      Default Portfolio Settings
                    </Typography>
                    
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                      <Grid container spacing={0.5}>
                        <Grid item xs={6}>
                          <FormControl fullWidth size="small">
                            <InputLabel id="exchange-label" sx={{ fontSize: '0.7rem' }}>Exchange</InputLabel>
                            <Select
                              labelId="exchange-label"
                              value={exchangeValue}
                              label="Exchange"
                              onChange={handleExchangeChange}
                              sx={{ fontSize: '0.7rem', height: 24 }}
                            >
                              <MenuItem value="NSE" sx={{ fontSize: '0.7rem' }}>NSE</MenuItem>
                              <MenuItem value="BSE" sx={{ fontSize: '0.7rem' }}>BSE</MenuItem>
                            </Select>
                          </FormControl>
                        </Grid>
                        <Grid item xs={6}>
                          <FormControl fullWidth size="small">
                            <InputLabel id="symbol-label" sx={{ fontSize: '0.7rem' }}>Symbol</InputLabel>
                            <Select
                              labelId="symbol-label"
                              value={symbolValue}
                              label="Symbol"
                              onChange={handleSymbolChange}
                              sx={{ fontSize: '0.7rem', height: 24 }}
                            >
                              <MenuItem value="NIFTY" sx={{ fontSize: '0.7rem' }}>NIFTY</MenuItem>
                              <MenuItem value="BANKNIFTY" sx={{ fontSize: '0.7rem' }}>BANKNIFTY</MenuItem>
                              <MenuItem value="FINNIFTY" sx={{ fontSize: '0.7rem' }}>FINNIFTY</MenuItem>
                            </Select>
                          </FormControl>
                        </Grid>
                      </Grid>
                      
                      <Box sx={{ 
                        bgcolor: '#f5f5f5', 
                        border: '1px solid #e0e0e0', 
                        borderRadius: 1, 
                        p: 0.3,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between'
                      }}>
                        <Typography variant="caption" sx={{ fontWeight: 'medium', color: '#555', fontSize: '0.65rem' }}>
                          LIVE
                        </Typography>
                        <Typography variant="caption" sx={{ fontWeight: 'bold', color: '#333', fontSize: '0.65rem' }}>
                          {niftyLiveValue}
                        </Typography>
                      </Box>
                      
                      <FormControl fullWidth size="small">
                        <InputLabel id="expiry-label" sx={{ fontSize: '0.7rem' }}>Expiry</InputLabel>
                        <Select
                          labelId="expiry-label"
                          value={expiryValue}
                          label="Expiry"
                          onChange={handleExpiryChange}
                          sx={{ fontSize: '0.7rem', height: 24 }}
                        >
                          <MenuItem value="09-Apr-25" sx={{ fontSize: '0.7rem' }}>09-Apr-25</MenuItem>
                          <MenuItem value="16-Apr-25" sx={{ fontSize: '0.7rem' }}>16-Apr-25</MenuItem>
                          <MenuItem value="23-Apr-25" sx={{ fontSize: '0.7rem' }}>23-Apr-25</MenuItem>
                          <MenuItem value="30-Apr-25" sx={{ fontSize: '0.7rem' }}>30-Apr-25</MenuItem>
                        </Select>
                      </FormControl>
                      
                      <Grid container spacing={0.5}>
                        <Grid item xs={6}>
                          <TextField
                            size="small"
                            label="Default Lots"
                            value={defaultLots}
                            onChange={handleDefaultLotsChange}
                            sx={{ '& .MuiInputBase-root': { height: 24, fontSize: '0.7rem' } }}
                          />
                        </Grid>
                        <Grid item xs={6}>
                          <TextField
                            size="small"
                            label="Lot Size"
                            value={lotSize}
                            onChange={handleLotSizeChange}
                            sx={{ '& .MuiInputBase-root': { height: 24, fontSize: '0.7rem' } }}
                          />
                        </Grid>
                      </Grid>
                      
                      <TextField
                        size="small"
                        label="Pay Premium"
                        value={payPremium}
                        onChange={handlePayPremiumChange}
                        sx={{ '& .MuiInputBase-root': { height: 24, fontSize: '0.7rem' } }}
                      />
                      
                      <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
                        <Tooltip title="When checked, portfolio will be treated as positional (overnight) rather than intraday" placement="top">
                          <FormControlLabel
                            control={
                              <Checkbox 
                                checked={positionalPortfolio}
                                onChange={handlePositionalPortfolioChange}
                                size="small"
                                sx={{ '& .MuiSvgIcon-root': { fontSize: 16 } }}
                              />
                            }
                            label={<Typography variant="caption" sx={{ fontSize: '0.65rem' }}>Positional Portfolio</Typography>}
                            sx={{ m: 0, width: '100%', height: 20 }}
                          />
                        </Tooltip>
                        <Tooltip title="When checked, buy trades will be executed before sell trades" placement="top">
                          <FormControlLabel
                            control={
                              <Checkbox 
                                checked={buyTradesFirst}
                                onChange={handleBuyTradesFirstChange}
                                size="small"
                                sx={{ '& .MuiSvgIcon-root': { fontSize: 16 } }}
                              />
                            }
                            label={<Typography variant="caption" sx={{ fontSize: '0.65rem' }}>Buy Trades First</Typography>}
                            sx={{ m: 0, width: '50%', height: 20 }}
                          />
                        </Tooltip>
                        <Tooltip title="When checked, allows selection of strikes far from ATM. Useful for extreme tail risk strategies." placement="top">
                          <FormControlLabel
                            control={
                              <Checkbox 
                                checked={allowFarStrikes}
                                onChange={handleAllowFarStrikesChange}
                                size="small"
                                sx={{ '& .MuiSvgIcon-root': { fontSize: 16 } }}
                              />
                            }
                            label={<Typography variant="caption" sx={{ fontSize: '0.65rem' }}>Allow Far Strikes</Typography>}
                            sx={{ m: 0, width: '50%', height: 20 }}
                          />
                        </Tooltip>
                      </Box>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              {/* Value Calculation Display */}
              <Grid item xs={12} md={4}>
                <Card variant="outlined" sx={{ height: '100%', bgcolor: '#fff' }}>
                  <CardContent sx={{ p: 0.5, '&:last-child': { pb: 0.5 } }}>
                    <Typography variant="caption" sx={{ mb: 0.3, color: '#2c3e50', fontWeight: 'bold', display: 'block', fontSize: '0.65rem' }}>
                      Value Calculation
                    </Typography>
                    
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                      <Box sx={{ 
                        bgcolor: '#f5f5f5', 
                        border: '1px solid #e0e0e0', 
                        borderRadius: 1, 
                        p: 0.3,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between'
                      }}>
                        <Typography variant="caption" sx={{ fontWeight: 'medium', color: '#555', fontSize: '0.65rem' }}>
                          VALUE ALL LOTS
                        </Typography>
                        <Typography variant="caption" sx={{ fontWeight: 'bold', color: '#333', fontSize: '0.65rem' }}>
                          {valueAllLots}
                        </Typography>
                      </Box>
                      
                      <Box sx={{ 
                        bgcolor: '#f5f5f5', 
                        border: '1px solid #e0e0e0', 
                        borderRadius: 1, 
                        p: 0.3,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between'
                      }}>
                        <Typography variant="caption" sx={{ fontWeight: 'medium', color: '#555', fontSize: '0.65rem' }}>
                          VALUE PER LOT
                        </Typography>
                        <Typography variant="caption" sx={{ fontWeight: 'bold', color: '#333', fontSize: '0.65rem' }}>
                          {valuePerLot}
                        </Typography>
                      </Box>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              
              {/* Action Buttons */}
              <Grid item xs={12} md={4}>
                <Card variant="outlined" sx={{ height: '100%', bgcolor: '#fff' }}>
                  <CardContent sx={{ p: 0.5, '&:last-child': { pb: 0.5 } }}>
                    <Typography variant="caption" sx={{ mb: 0.3, color: '#2c3e50', fontWeight: 'bold', display: 'block', fontSize: '0.65rem' }}>
                      Actions
                    </Typography>
                    
                    <Box sx={{ display: 'flex', flexDirection: 'row', gap: 0.5, alignItems: 'center' }}>
                      <Button 
                        variant="outlined" 
                        size="small" 
                        startIcon={<RefreshIcon sx={{ fontSize: '0.7rem' }} />}
                        sx={{ fontSize: '0.65rem', py: 0.3, minHeight: 0, height: 22 }}
                        onClick={handleRefreshClick}
                      >
                        Refresh
                      </Button>
                      <Button 
                        variant="outlined" 
                        size="small" 
                        startIcon={<AddIcon sx={{ fontSize: '0.7rem' }} />}
                        sx={{ fontSize: '0.65rem', py: 0.3, minHeight: 0, height: 22 }}
                        onClick={handleAddLegClick}
                      >
                        Add Leg
                      </Button>
                      <FormControlLabel
                        control={
                          <Checkbox 
                            checked={premiumGreekLog}
                            onChange={handlePremiumGreekLogChange}
                            size="small"
                            sx={{ '& .MuiSvgIcon-root': { fontSize: 16 } }}
                          />
                        }
                        label={<Typography variant="caption" sx={{ fontSize: '0.65rem' }}>Premium/Greek Log</Typography>}
                        sx={{ m: 0, height: 20 }}
                      />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </Box>
        </Box>
      </Paper>
    </Box>
  );
};

export default PortfolioComponent;
