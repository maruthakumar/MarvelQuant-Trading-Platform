import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  TextField, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Button,
  Grid,
  Breadcrumbs,
  Link,
  Container,
  Card,
  CardContent,
  Divider,
  Tabs,
  Tab,
  FormControlLabel,
  Switch,
  Tooltip,
  IconButton
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import InfoIcon from '@mui/icons-material/Info';

const NewPortfolioPage = () => {
  const navigate = useNavigate();
  const [portfolioName, setPortfolioName] = useState('');
  const [symbol, setSymbol] = useState('NIFTY');
  const [strategy, setStrategy] = useState('BACKENZOBUYING');
  const [activeTab, setActiveTab] = useState(0);
  
  // Basic Information
  const [description, setDescription] = useState('');
  const [initialCapital, setInitialCapital] = useState('100000');
  const [exchangeValue, setExchangeValue] = useState('NSE');
  const [expiryValue, setExpiryValue] = useState('09-Apr-25');
  const [defaultLots, setDefaultLots] = useState('10');
  const [lotSize, setLotSize] = useState('75');
  
  // Strategy Configuration
  const [strikeSelection, setStrikeSelection] = useState('Relative');
  const [underlying, setUnderlying] = useState('Spot');
  const [priceType, setPriceType] = useState('LTP');
  const [strikeStep, setStrikeStep] = useState('50');
  
  // Portfolio Behavior
  const [positionalPortfolio, setPositionalPortfolio] = useState(false);
  const [buyTradesFirst, setBuyTradesFirst] = useState(true);
  const [allowFarStrikes, setAllowFarStrikes] = useState(false);
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handleCreatePortfolio = () => {
    if (!portfolioName) {
      alert('Portfolio name is required');
      return;
    }
    
    // Here you would typically make an API call to create the portfolio
    // For now, we'll just navigate back to the multi-leg page
    navigate('/multi-leg', { 
      state: { 
        newPortfolio: {
          name: portfolioName,
          symbol,
          strategy,
          description,
          initialCapital,
          exchange: exchangeValue,
          expiry: expiryValue,
          defaultLots,
          lotSize,
          strikeSelection,
          underlying,
          priceType,
          strikeStep,
          positionalPortfolio,
          buyTradesFirst,
          allowFarStrikes
        },
        showSuccessMessage: true
      } 
    });
  };
  
  const handleCancel = () => {
    navigate('/multi-leg');
  };
  
  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={handleCancel}
          sx={{ mr: 2 }}
        >
          Back
        </Button>
        <Breadcrumbs aria-label="breadcrumb">
          <Link color="inherit" href="/multi-leg" onClick={(e) => { e.preventDefault(); navigate('/multi-leg'); }}>
            Multi-Leg
          </Link>
          <Typography color="text.primary">New Portfolio</Typography>
        </Breadcrumbs>
      </Box>
      
      <Typography variant="h4" gutterBottom>
        Create New Portfolio
      </Typography>
      
      <Box sx={{ width: '100%', mb: 3 }}>
        <Tabs 
          value={activeTab} 
          onChange={handleTabChange} 
          aria-label="portfolio configuration tabs"
          sx={{
            '& .MuiTabs-indicator': {
              backgroundColor: '#1976d2',
            },
            '& .Mui-selected': {
              color: '#1976d2',
              fontWeight: 'bold',
            },
          }}
        >
          <Tab label="Basic Information" />
          <Tab label="Strategy Configuration" />
          <Tab label="Portfolio Behavior" />
        </Tabs>
      </Box>
      
      {/* Basic Information Tab */}
      {activeTab === 0 && (
        <Paper sx={{ p: 3, mb: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6">Basic Information</Typography>
            <Tooltip title="Enter the fundamental details of your portfolio">
              <IconButton size="small" sx={{ ml: 1 }}>
                <InfoIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
          <Divider sx={{ mb: 3 }} />
          
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                required
                id="portfolioName"
                label="Portfolio Name"
                fullWidth
                value={portfolioName}
                onChange={(e) => setPortfolioName(e.target.value)}
                helperText="Enter a unique name for your portfolio"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                id="description"
                label="Description"
                fullWidth
                multiline
                rows={2}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                helperText="Optional: Add a description for this portfolio"
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel id="exchange-label">Exchange</InputLabel>
                <Select
                  labelId="exchange-label"
                  id="exchange"
                  value={exchangeValue}
                  label="Exchange"
                  onChange={(e) => setExchangeValue(e.target.value)}
                >
                  <MenuItem value="NSE">NSE</MenuItem>
                  <MenuItem value="BSE">BSE</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel id="symbol-label">Symbol</InputLabel>
                <Select
                  labelId="symbol-label"
                  id="symbol"
                  value={symbol}
                  label="Symbol"
                  onChange={(e) => setSymbol(e.target.value)}
                >
                  <MenuItem value="NIFTY">NIFTY</MenuItem>
                  <MenuItem value="BANKNIFTY">BANKNIFTY</MenuItem>
                  <MenuItem value="FINNIFTY">FINNIFTY</MenuItem>
                  <MenuItem value="SENSEX">SENSEX</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel id="expiry-label">Expiry</InputLabel>
                <Select
                  labelId="expiry-label"
                  id="expiry"
                  value={expiryValue}
                  label="Expiry"
                  onChange={(e) => setExpiryValue(e.target.value)}
                >
                  <MenuItem value="09-Apr-25">09-Apr-25</MenuItem>
                  <MenuItem value="16-Apr-25">16-Apr-25</MenuItem>
                  <MenuItem value="23-Apr-25">23-Apr-25</MenuItem>
                  <MenuItem value="30-Apr-25">30-Apr-25</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                id="defaultLots"
                label="Default Lots"
                fullWidth
                type="number"
                value={defaultLots}
                onChange={(e) => setDefaultLots(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                id="lotSize"
                label="Lot Size"
                fullWidth
                type="number"
                value={lotSize}
                onChange={(e) => setLotSize(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                id="initialCapital"
                label="Initial Capital"
                fullWidth
                type="number"
                value={initialCapital}
                onChange={(e) => setInitialCapital(e.target.value)}
                InputProps={{
                  startAdornment: <span style={{ marginRight: 8 }}>₹</span>,
                }}
              />
            </Grid>
          </Grid>
        </Paper>
      )}
      
      {/* Strategy Configuration Tab */}
      {activeTab === 1 && (
        <Paper sx={{ p: 3, mb: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6">Strategy Configuration</Typography>
            <Tooltip title="Configure the strategy parameters for your portfolio">
              <IconButton size="small" sx={{ ml: 1 }}>
                <InfoIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
          <Divider sx={{ mb: 3 }} />
          
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel id="strategy-label">Strategy</InputLabel>
                <Select
                  labelId="strategy-label"
                  id="strategy"
                  value={strategy}
                  label="Strategy"
                  onChange={(e) => setStrategy(e.target.value)}
                >
                  <MenuItem value="BACKENZOBUYING">BACKENZOBUYING</MenuItem>
                  <MenuItem value="NF-NDSTR-D">NF-NDSTR-D</MenuItem>
                  <MenuItem value="CUSTOM">CUSTOM</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel id="strike-selection-label">Strike Selection</InputLabel>
                <Select
                  labelId="strike-selection-label"
                  id="strikeSelection"
                  value={strikeSelection}
                  label="Strike Selection"
                  onChange={(e) => setStrikeSelection(e.target.value)}
                >
                  <MenuItem value="Relative">Relative</MenuItem>
                  <MenuItem value="Absolute">Absolute</MenuItem>
                  <MenuItem value="ATM">ATM</MenuItem>
                  <MenuItem value="NearestDelta">Nearest Delta</MenuItem>
                  <MenuItem value="NearestPremium">Nearest Premium</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel id="underlying-label">Underlying</InputLabel>
                <Select
                  labelId="underlying-label"
                  id="underlying"
                  value={underlying}
                  label="Underlying"
                  onChange={(e) => setUnderlying(e.target.value)}
                >
                  <MenuItem value="Spot">Spot</MenuItem>
                  <MenuItem value="Futures">Futures</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel id="price-type-label">Price Type</InputLabel>
                <Select
                  labelId="price-type-label"
                  id="priceType"
                  value={priceType}
                  label="Price Type"
                  onChange={(e) => setPriceType(e.target.value)}
                >
                  <MenuItem value="LTP">LTP</MenuItem>
                  <MenuItem value="Bid">Bid</MenuItem>
                  <MenuItem value="Ask">Ask</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                id="strikeStep"
                label="Strike Step"
                fullWidth
                type="number"
                value={strikeStep}
                onChange={(e) => setStrikeStep(e.target.value)}
              />
            </Grid>
          </Grid>
        </Paper>
      )}
      
      {/* Portfolio Behavior Tab */}
      {activeTab === 2 && (
        <Paper sx={{ p: 3, mb: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6">Portfolio Behavior</Typography>
            <Tooltip title="Configure how your portfolio behaves during trading">
              <IconButton size="small" sx={{ ml: 1 }}>
                <InfoIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
          <Divider sx={{ mb: 3 }} />
          
          <Grid container spacing={3}>
            <Grid item xs={12} md={4}>
              <FormControlLabel
                control={
                  <Switch
                    checked={positionalPortfolio}
                    onChange={(e) => setPositionalPortfolio(e.target.checked)}
                    name="positionalPortfolio"
                  />
                }
                label="Positional Portfolio"
              />
              <Tooltip title="Enable for portfolios held overnight">
                <IconButton size="small">
                  <HelpOutlineIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControlLabel
                control={
                  <Switch
                    checked={buyTradesFirst}
                    onChange={(e) => setBuyTradesFirst(e.target.checked)}
                    name="buyTradesFirst"
                  />
                }
                label="Buy Trades First"
              />
              <Tooltip title="Execute buy trades before sell trades">
                <IconButton size="small">
                  <HelpOutlineIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControlLabel
                control={
                  <Switch
                    checked={allowFarStrikes}
                    onChange={(e) => setAllowFarStrikes(e.target.checked)}
                    name="allowFarStrikes"
                  />
                }
                label="Allow Far Strikes"
              />
              <Tooltip title="Allow selection of strikes far from ATM">
                <IconButton size="small">
                  <HelpOutlineIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Paper>
      )}
      
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
        <Button 
          variant="outlined" 
          onClick={() => setActiveTab(Math.max(0, activeTab - 1))}
          disabled={activeTab === 0}
        >
          Previous
        </Button>
        
        <Box>
          <Button 
            variant="outlined" 
            onClick={handleCancel}
            sx={{ mr: 2 }}
          >
            Cancel
          </Button>
          {activeTab === 2 ? (
            <Button 
              variant="contained" 
              onClick={handleCreatePortfolio}
              disabled={!portfolioName}
            >
              Create Portfolio
            </Button>
          ) : (
            <Button 
              variant="contained" 
              onClick={() => setActiveTab(Math.min(2, activeTab + 1))}
            >
              Next
            </Button>
          )}
        </Box>
      </Box>
    </Container>
  );
};

export default NewPortfolioPage;
