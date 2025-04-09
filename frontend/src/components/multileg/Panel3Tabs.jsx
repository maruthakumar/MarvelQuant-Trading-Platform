import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Tabs, 
  Tab, 
  Paper, 
  Grid, 
  TextField, 
  Select, 
  MenuItem, 
  FormControl, 
  InputLabel, 
  Checkbox, 
  FormControlLabel, 
  FormGroup, 
  Radio, 
  RadioGroup, 
  Button, 
  IconButton, 
  InputAdornment, 
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Snackbar,
  Alert
} from '@mui/material';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import RefreshIcon from '@mui/icons-material/Refresh';
import InfoIcon from '@mui/icons-material/Info';

const Panel3Tabs = ({ portfolio }) => {
  const [activeTab, setActiveTab] = useState(0);
  
  // General Settings Tab State
  const [runOnDays, setRunOnDays] = useState('Monday,Tuesday,Wednesday,Thursday,Friday');
  const [startTime, setStartTime] = useState('09:15:00');
  const [endTime, setEndTime] = useState('15:15:00');
  const [sqOffTime, setSqOffTime] = useState('15:25:00');
  const [estimatedMargin, setEstimatedMargin] = useState('₹ 0.00');
  
  // Range Breakout Tab State
  const [underlyingBreakout, setUnderlyingBreakout] = useState(false);
  const [legRangeBreakout, setLegRangeBreakout] = useState(false);
  const [strategyRangeBreakout, setStrategyRangeBreakout] = useState(false);
  
  // Extra Conditions Tab State
  const [gapUpGapDownConditions, setGapUpGapDownConditions] = useState(false);
  const [minGapUp, setMinGapUp] = useState('');
  const [maxGapUp, setMaxGapUp] = useState('');
  const [minGapDown, setMinGapDown] = useState('');
  const [maxGapDown, setMaxGapDown] = useState('');
  const [dayOpenCondition, setDayOpenCondition] = useState('None');
  const [legDayOpenCondition, setLegDayOpenCondition] = useState('None');
  const [changeType, setChangeType] = useState('Positive');
  const [minChange, setMinChange] = useState('');
  const [maxChange, setMaxChange] = useState('');
  const [waitTradeValue, setWaitTradeValue] = useState('');
  const [waitTradeMonitoring, setWaitTradeMonitoring] = useState('Realtime');
  const [executeOtherLegs, setExecuteOtherLegs] = useState(false);
  const [allCELegsMaxProfit, setAllCELegsMaxProfit] = useState('');
  const [allCELegsMaxLoss, setAllCELegsMaxLoss] = useState('');
  const [allPELegsMaxProfit, setAllPELegsMaxProfit] = useState('');
  const [allPELegsMaxLoss, setAllPELegsMaxLoss] = useState('');
  
  // Other Settings Tab State
  const [keepAllUsersInSync, setKeepAllUsersInSync] = useState(false);
  const [trailWaitTrade, setTrailWaitTrade] = useState(false);
  const [executeDelaysInSeconds, setExecuteDelaysInSeconds] = useState('');
  const [executeDelaysInMinutes, setExecuteDelaysInMinutes] = useState('');
  const [reExecuteDelaysInSeconds, setReExecuteDelaysInSeconds] = useState('');
  const [reExecuteDelaysInMinutes, setReExecuteDelaysInMinutes] = useState('');
  const [straddleWidthMultiplier, setStraddleWidthMultiplier] = useState('');
  const [delayBetweenLegsInSec, setDelayBetweenLegsInSec] = useState('');
  const [onTargetAction, setOnTargetAction] = useState('OnTarget_N_');
  const [onSLActionOn, setOnSLActionOn] = useState('OnSL_N_Tra_');
  const [legReExecution, setLegReExecution] = useState('');
  const [portfolioReExecutionSafetySeconds, setPortfolioReExecutionSafetySeconds] = useState('');
  const [portfolioReExecuteCount, setPortfolioReExecuteCount] = useState('');
  const [orderSlicingSettings, setOrderSlicingSettings] = useState('NotReqd');
  const [slicingType, setSlicingType] = useState('');
  const [ifOneSideActivated, setIfOneSideActivated] = useState(false);
  const [noReEntryReExecute, setNoReEntryReExecute] = useState(false);
  const [maxLegsSupported, setMaxLegsSupported] = useState('');
  
  // Monitoring Tab State
  const [positionalPortfolioTargetMonitoringStartTime, setPositionalPortfolioTargetMonitoringStartTime] = useState('');
  const [positionalPortfolioTargetMonitoringEndTime, setPositionalPortfolioTargetMonitoringEndTime] = useState('');
  const [combinedTargetSLMonitoring, setCombinedTargetSLMonitoring] = useState('Realtime');
  const [slMonitoring, setSlMonitoring] = useState('Realtime');
  const [minuteCloseTradeSeconds, setMinuteCloseTradeSeconds] = useState('');
  const [legReExecuteSettingsMonitoring, setLegReExecuteSettingsMonitoring] = useState('Realtime');
  const [noReExecuteIfMovedSLToCost, setNoReExecuteIfMovedSLToCost] = useState(false);
  const [reEntryMonitoring, setReEntryMonitoring] = useState('Realtime');
  const [reEntryTriggerOn, setReEntryTriggerOn] = useState('None');
  const [orderType, setOrderType] = useState('MARKET');
  const [minDelayInReEntry, setMinDelayInReEntry] = useState('');
  const [reEntryAtOriginalEntryPrice, setReEntryAtOriginalEntryPrice] = useState(false);
  
  // Dynamic Hedge Tab State
  const [hedgeDistanceFromATMMin, setHedgeDistanceFromATMMin] = useState('');
  const [hedgeDistanceFromATMMax, setHedgeDistanceFromATMMax] = useState('');
  const [minPremium, setMinPremium] = useState('');
  const [maxPremium, setMaxPremium] = useState('');
  const [sqOffLegOn, setSqOffLegOn] = useState('');
  const [unsatisfiedConditionAction, setUnsatisfiedConditionAction] = useState('IgnoreField');
  const [selectOnly500Strikes, setSelectOnly500Strikes] = useState(false);
  
  // Target Settings Tab State
  const [targetType, setTargetType] = useState('None');
  const [ifProfitReaches, setIfProfitReaches] = useState('');
  const [lockMinimumProfitAt, setLockMinimumProfitAt] = useState('');
  const [forEveryIncreaseInProfitBy, setForEveryIncreaseInProfitBy] = useState('');
  const [trailProfitBy, setTrailProfitBy] = useState('');
  
  // Stoploss Settings Tab State
  const [stoplossType, setStoplossType] = useState('None');
  const [moveSLToCost, setMoveSLToCost] = useState(false);
  
  // Exit Settings Tab State
  const [exitSettings, setExitSettings] = useState({});
  const [exitType, setExitType] = useState('Time');
  const [exitTime, setExitTime] = useState('15:25:00');
  const [exitOnProfit, setExitOnProfit] = useState(false);
  const [profitAmount, setProfitAmount] = useState('');
  const [exitOnLoss, setExitOnLoss] = useState(false);
  const [lossAmount, setLossAmount] = useState('');
  const [partialExitEnabled, setPartialExitEnabled] = useState(false);
  const [partialExitPercentage, setPartialExitPercentage] = useState('50');
  const [partialExitTrigger, setPartialExitTrigger] = useState('Profit');
  const [partialExitValue, setPartialExitValue] = useState('');
  const [exitStrategy, setExitStrategy] = useState('CloseAll');
  const [exitOrder, setExitOrder] = useState('MARKET');
  
  // At Broker Tab State
  const [logSLAtBroker, setLogSLAtBroker] = useState(false);
  const [logTargetAtBroker, setLogTargetAtBroker] = useState(false);
  const [legReEntryAtBroker, setLegReEntryAtBroker] = useState(false);
  const [legWaitTradeAtBroker, setLegWaitTradeAtBroker] = useState(false);

  // Advanced Stop Loss/Target Action Configuration Dialog State
  const [slActionDialogOpen, setSlActionDialogOpen] = useState(false);
  const [targetActionDialogOpen, setTargetActionDialogOpen] = useState(false);
  const [advancedSLConfig, setAdvancedSLConfig] = useState({
    actionType: 'SqOff',
    trailType: 'Percentage',
    trailValue: '10',
    reentryDelay: '60',
    reentryEnabled: false,
    notifyUser: true,
    executeImmediately: true,
    reversePosition: false
  });
  const [advancedTargetConfig, setAdvancedTargetConfig] = useState({
    actionType: 'SqOff',
    trailType: 'Percentage',
    trailValue: '10',
    reentryDelay: '60',
    reentryEnabled: false,
    notifyUser: true,
    executeImmediately: true,
    reversePosition: false
  });

  // Feedback Snackbar State
  const [feedbackOpen, setFeedbackOpen] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState('');
  const [feedbackSeverity, setFeedbackSeverity] = useState('info');

  // Handle tab change
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };

  // Open SL Action Configuration Dialog
  const handleOpenSLActionDialog = () => {
    setSlActionDialogOpen(true);
  };

  // Close SL Action Configuration Dialog
  const handleCloseSLActionDialog = () => {
    setSlActionDialogOpen(false);
  };

  // Open Target Action Configuration Dialog
  const handleOpenTargetActionDialog = () => {
    setTargetActionDialogOpen(true);
  };

  // Close Target Action Configuration Dialog
  const handleCloseTargetActionDialog = () => {
    setTargetActionDialogOpen(false);
  };

  // Handle SL Action Configuration Change
  const handleSLConfigChange = (field, value) => {
    setAdvancedSLConfig({
      ...advancedSLConfig,
      [field]: value
    });
  };

  // Handle Target Action Configuration Change
  const handleTargetConfigChange = (field, value) => {
    setAdvancedTargetConfig({
      ...advancedTargetConfig,
      [field]: value
    });
  };

  // Save SL Action Configuration
  const handleSaveSLConfig = () => {
    // Update onSLActionOn based on configuration
    let newSLAction = 'OnSL_';
    
    if (advancedSLConfig.actionType === 'SqOff') {
      newSLAction += 'SqOff';
    } else if (advancedSLConfig.actionType === 'Trail') {
      newSLAction += 'Trail';
    } else {
      newSLAction += 'N_Tra_';
    }
    
    setOnSLActionOn(newSLAction);
    
    // Show feedback
    setFeedbackMessage('Stop Loss action configuration saved successfully!');
    setFeedbackSeverity('success');
    setFeedbackOpen(true);
    
    // Close dialog
    setSlActionDialogOpen(false);
  };

  // Save Target Action Configuration
  const handleSaveTargetConfig = () => {
    // Update onTargetAction based on configuration
    let newTargetAction = 'OnTarget_';
    
    if (advancedTargetConfig.actionType === 'SqOff') {
      newTargetAction += 'SqOff';
    } else if (advancedTargetConfig.actionType === 'Trail') {
      newTargetAction += 'Trail';
    } else {
      newTargetAction += 'N_';
    }
    
    setOnTargetAction(newTargetAction);
    
    // Show feedback
    setFeedbackMessage('Target action configuration saved successfully!');
    setFeedbackSeverity('success');
    setFeedbackOpen(true);
    
    // Close dialog
    setTargetActionDialogOpen(false);
  };

  // Handle feedback close
  const handleFeedbackClose = () => {
    setFeedbackOpen(false);
  };
  
  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs 
          value={activeTab} 
          onChange={handleTabChange} 
          variant="scrollable"
          scrollButtons="auto"
          sx={{ 
            '& .MuiTab-root': { 
              fontSize: '0.75rem',
              minWidth: 'auto',
              py: 1,
              px: 1.5
            } 
          }}
        >
          <Tab label="General Settings" />
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
      
      {/* General Settings Tab */}
      <TabPanel value={activeTab} index={0}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>Timings</Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6} md={3}>
                <FormControl fullWidth size="small">
                  <InputLabel>Run On Days</InputLabel>
                  <Select
                    value={runOnDays}
                    label="Run On Days"
                    onChange={(e) => setRunOnDays(e.target.value)}
                  >
                    <MenuItem value="Monday,Tuesday,Wednesday,Thursday,Friday">Monday,Tuesday,Wednesday,Thursday,Friday</MenuItem>
                    <MenuItem value="Monday">Monday</MenuItem>
                    <MenuItem value="Tuesday">Tuesday</MenuItem>
                    <MenuItem value="Wednesday">Wednesday</MenuItem>
                    <MenuItem value="Thursday">Thursday</MenuItem>
                    <MenuItem value="Friday">Friday</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} sm={6} md={2}>
                <TextField
                  fullWidth
                  size="small"
                  label="Start Time"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton size="small">
                          <HelpOutlineIcon fontSize="small" />
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6} md={2}>
                <TextField
                  fullWidth
                  size="small"
                  label="End Time"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton size="small">
                          <HelpOutlineIcon fontSize="small" />
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6} md={2}>
                <TextField
                  fullWidth
                  size="small"
                  label="SqOff Time"
                  value={sqOffTime}
                  onChange={(e) => setSqOffTime(e.target.value)}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton size="small">
                          <HelpOutlineIcon fontSize="small" />
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Box sx={{ display: 'flex', alignItems: 'center', height: '100%' }}>
                  <Typography variant="body2" sx={{ mr: 1 }}>Estimated Margin:</Typography>
                  <Typography variant="body2" fontWeight="bold">{estimatedMargin}</Typography>
                  <IconButton size="small" sx={{ ml: 1 }}>
                    <RefreshIcon fontSize="small" />
                  </IconButton>
                </Box>
              </Grid>
            </Grid>
          </Grid>
          
          <Grid item xs={12}>
            <Typography variant="caption" color="textSecondary">
              Note: Settings like User Accounts, Refines, Hold Sell Sec etc will be taken from Strategy Tab. However Start and End Time etc will be taken from here.
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="caption" color="textSecondary">
              1. Target, SL, Trailing, Spread, Multi Tgt & Qty, Max Slippage etc can be filled in Points or in percent. Eg 10 will be points and 10% will be percent. For absolute values like for BUY at 200 then SL at 180 select absolute.
            </Typography>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Range Breakout Tab */}
      <TabPanel value={activeTab} index={1}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={4}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Underlying Range Breakout</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={underlyingBreakout}
                    onChange={(e) => setUnderlyingBreakout(e.target.checked)}
                  />
                }
                label="Underlying Breakout"
              />
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Legs Range Breakout</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={legRangeBreakout}
                    onChange={(e) => setLegRangeBreakout(e.target.checked)}
                  />
                }
                label="Leg Range Breakout"
              />
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Strategy Range Breakout</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={strategyRangeBreakout}
                    onChange={(e) => setStrategyRangeBreakout(e.target.checked)}
                  />
                }
                label="Strategy Range Breakout"
              />
            </Paper>
          </Grid>
          
          <Grid item xs={12}>
            <Typography variant="caption" color="textSecondary">
              Note: If Underlying & Leg Range Breakout both specified then Underlying range will be monitored first and the Leg range after selection of the Strikes.
            </Typography>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Extra Conditions Tab */}
      <TabPanel value={activeTab} index={2}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Underlying Gap-Up / Gap-Down</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={gapUpGapDownConditions}
                    onChange={(e) => setGapUpGapDownConditions(e.target.checked)}
                  />
                }
                label="Gap-Up / Gap-Down Conditions"
              />
              
              <Grid container spacing={2} sx={{ mt: 1 }}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Min Gap-Up"
                    value={minGapUp}
                    onChange={(e) => setMinGapUp(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Gap-Up"
                    value={maxGapUp}
                    onChange={(e) => setMaxGapUp(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Min Gap-Down"
                    value={minGapDown}
                    onChange={(e) => setMinGapDown(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Gap-Down"
                    value={maxGapDown}
                    onChange={(e) => setMaxGapDown(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Underlying Day Open Settings</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Day Open Condition *</InputLabel>
                <Select
                  value={dayOpenCondition}
                  label="Day Open Condition *"
                  onChange={(e) => setDayOpenCondition(e.target.value)}
                >
                  <MenuItem value="None">None</MenuItem>
                  <MenuItem value="Above">Above</MenuItem>
                  <MenuItem value="Below">Below</MenuItem>
                  <MenuItem value="Between">Between</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>Leg's Settings</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Leg's Day Open Condition *</InputLabel>
                <Select
                  value={legDayOpenCondition}
                  label="Leg's Day Open Condition *"
                  onChange={(e) => setLegDayOpenCondition(e.target.value)}
                >
                  <MenuItem value="None">None</MenuItem>
                  <MenuItem value="Above">Above</MenuItem>
                  <MenuItem value="Below">Below</MenuItem>
                  <MenuItem value="Between">Between</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>Leg's Current Day Change</Typography>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <FormControl component="fieldset">
                    <RadioGroup
                      row
                      value={changeType}
                      onChange={(e) => setChangeType(e.target.value)}
                    >
                      <FormControlLabel value="Positive" control={<Radio />} label="Positive" />
                      <FormControlLabel value="Negative" control={<Radio />} label="Negative" />
                    </RadioGroup>
                  </FormControl>
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Min Change"
                    value={minChange}
                    onChange={(e) => setMinChange(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Change"
                    value={maxChange}
                    onChange={(e) => setMaxChange(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Combined Wait & Trade Settings</Typography>
              <TextField
                fullWidth
                size="small"
                label="Wait & Trade Value"
                value={waitTradeValue}
                onChange={(e) => setWaitTradeValue(e.target.value)}
                sx={{ mb: 2 }}
              />
              
              <Typography variant="subtitle2" gutterBottom>Leg's Wait & Trade Settings</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Wait Trade Monitoring</InputLabel>
                <Select
                  value={waitTradeMonitoring}
                  label="Wait Trade Monitoring"
                  onChange={(e) => setWaitTradeMonitoring(e.target.value)}
                >
                  <MenuItem value="Realtime">Realtime</MenuItem>
                  <MenuItem value="1Min">1Min</MenuItem>
                  <MenuItem value="5Min">5Min</MenuItem>
                  <MenuItem value="15Min">15Min</MenuItem>
                </Select>
              </FormControl>
              
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={executeOtherLegs}
                    onChange={(e) => setExecuteOtherLegs(e.target.checked)}
                  />
                }
                label="If any leg hits W&T, then execute other legs at W&T Price or Execute if the price goes below the Adjusted Price or Time hit."
              />
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>All CE Legs Combined P&L (Rupees)</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Profit"
                    value={allCELegsMaxProfit}
                    onChange={(e) => setAllCELegsMaxProfit(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Loss"
                    value={allCELegsMaxLoss}
                    onChange={(e) => setAllCELegsMaxLoss(e.target.value)}
                  />
                </Grid>
              </Grid>
              
              <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>All PE Legs Combined P&L (Rupees)</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Profit"
                    value={allPELegsMaxProfit}
                    onChange={(e) => setAllPELegsMaxProfit(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Loss"
                    value={allPELegsMaxLoss}
                    onChange={(e) => setAllPELegsMaxLoss(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Other Settings Tab */}
      <TabPanel value={activeTab} index={3}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={keepAllUsersInSync}
                      onChange={(e) => setKeepAllUsersInSync(e.target.checked)}
                    />
                  }
                  label="Keep All Users In Sync"
                />
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={trailWaitTrade}
                      onChange={(e) => setTrailWaitTrade(e.target.checked)}
                    />
                  }
                  label="Trail Wait & Trade"
                />
              </FormGroup>
              
              <Grid container spacing={2} sx={{ mt: 1 }}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Execute / ReExecute Delays for Legs</Typography>
                  <TextField
                    fullWidth
                    size="small"
                    label="In Seconds"
                    value={executeDelaysInSeconds}
                    onChange={(e) => setExecuteDelaysInSeconds(e.target.value)}
                    sx={{ mb: 1 }}
                  />
                  <TextField
                    fullWidth
                    size="small"
                    label="In Candle Minutes"
                    value={executeDelaysInMinutes}
                    onChange={(e) => setExecuteDelaysInMinutes(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Execute / ReExecute Delays for Portfolio</Typography>
                  <TextField
                    fullWidth
                    size="small"
                    label="In Seconds"
                    value={reExecuteDelaysInSeconds}
                    onChange={(e) => setReExecuteDelaysInSeconds(e.target.value)}
                    sx={{ mb: 1 }}
                  />
                  <TextField
                    fullWidth
                    size="small"
                    label="In Candle Minutes"
                    value={reExecuteDelaysInMinutes}
                    onChange={(e) => setReExecuteDelaysInMinutes(e.target.value)}
                  />
                </Grid>
              </Grid>
              
              <Grid container spacing={2} sx={{ mt: 1 }}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Straddle Width Multiplier</Typography>
                  <TextField
                    fullWidth
                    size="small"
                    value={straddleWidthMultiplier}
                    onChange={(e) => setStraddleWidthMultiplier(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Delay Between Legs in Sec.</Typography>
                  <TextField
                    fullWidth
                    size="small"
                    value={delayBetweenLegsInSec}
                    onChange={(e) => setDelayBetweenLegsInSec(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>On Target Action</Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <FormControl fullWidth size="small" sx={{ mr: 1 }}>
                  <Select
                    value={onTargetAction}
                    onChange={(e) => setOnTargetAction(e.target.value)}
                  >
                    <MenuItem value="OnTarget_N_">OnTarget_N_</MenuItem>
                    <MenuItem value="OnTarget_SqOff">OnTarget_SqOff</MenuItem>
                    <MenuItem value="OnTarget_Trail">OnTarget_Trail</MenuItem>
                  </Select>
                </FormControl>
                <Button 
                  variant="outlined" 
                  size="small"
                  onClick={handleOpenTargetActionDialog}
                >
                  Configure
                </Button>
              </Box>
              
              <Typography variant="subtitle2" gutterBottom>On SL Action On</Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <FormControl fullWidth size="small" sx={{ mr: 1 }}>
                  <Select
                    value={onSLActionOn}
                    onChange={(e) => setOnSLActionOn(e.target.value)}
                  >
                    <MenuItem value="OnSL_N_Tra_">OnSL_N_Tra_</MenuItem>
                    <MenuItem value="OnSL_SqOff">OnSL_SqOff</MenuItem>
                    <MenuItem value="OnSL_Trail">OnSL_Trail</MenuItem>
                  </Select>
                </FormControl>
                <Button 
                  variant="outlined" 
                  size="small"
                  onClick={handleOpenSLActionDialog}
                >
                  Configure
                </Button>
              </Box>
              
              <Typography variant="subtitle2" gutterBottom>Leg ReExecution</Typography>
              <TextField
                fullWidth
                size="small"
                label="Safety Seconds"
                value={legReExecution}
                onChange={(e) => setLegReExecution(e.target.value)}
                sx={{ mb: 2 }}
              />
              
              <Typography variant="subtitle2" gutterBottom>Portfolio ReExecution Safety Seconds</Typography>
              <TextField
                fullWidth
                size="small"
                value={portfolioReExecutionSafetySeconds}
                onChange={(e) => setPortfolioReExecutionSafetySeconds(e.target.value)}
                sx={{ mb: 2 }}
              />
              
              <Typography variant="subtitle2" gutterBottom>Portfolio ReExecute Count</Typography>
              <Typography variant="caption" display="block" gutterBottom>
                Applicable only on Portfolio ReExecution.
              </Typography>
              <TextField
                fullWidth
                size="small"
                value={portfolioReExecuteCount}
                onChange={(e) => setPortfolioReExecuteCount(e.target.value)}
                sx={{ mb: 2 }}
              />
              
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Order Slicing Settings</Typography>
                  <FormControl fullWidth size="small">
                    <Select
                      value={orderSlicingSettings}
                      onChange={(e) => setOrderSlicingSettings(e.target.value)}
                    >
                      <MenuItem value="NotReqd">NotReqd</MenuItem>
                      <MenuItem value="Required">Required</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Slicing Type</Typography>
                  <FormControl fullWidth size="small">
                    <Select
                      value={slicingType}
                      onChange={(e) => setSlicingType(e.target.value)}
                    >
                      <MenuItem value="">Select Type</MenuItem>
                      <MenuItem value="Equal">Equal</MenuItem>
                      <MenuItem value="Random">Random</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
              
              <FormGroup sx={{ mt: 2 }}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={ifOneSideActivated}
                      onChange={(e) => setIfOneSideActivated(e.target.checked)}
                    />
                  }
                  label="If One Side Activated, then Cancel Other Side. (Applicable with W&T and Range Breakout Only)"
                />
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={noReEntryReExecute}
                      onChange={(e) => setNoReEntryReExecute(e.target.checked)}
                    />
                  }
                  label="NO ReEntry / ReExecute after Portfolio End Time."
                />
              </FormGroup>
              
              <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>Max Legs Supported</Typography>
              <TextField
                fullWidth
                size="small"
                value={maxLegsSupported}
                onChange={(e) => setMaxLegsSupported(e.target.value)}
              />
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Monitoring Tab */}
      <TabPanel value={activeTab} index={4}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Positional Portfolio Target Monitoring</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Start Time"
                    value={positionalPortfolioTargetMonitoringStartTime}
                    onChange={(e) => setPositionalPortfolioTargetMonitoringStartTime(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="End Time"
                    value={positionalPortfolioTargetMonitoringEndTime}
                    onChange={(e) => setPositionalPortfolioTargetMonitoringEndTime(e.target.value)}
                  />
                </Grid>
              </Grid>
              
              <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>Combined Target & SL Monitoring</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Target Monitoring</InputLabel>
                <Select
                  value={combinedTargetSLMonitoring}
                  label="Target Monitoring"
                  onChange={(e) => setCombinedTargetSLMonitoring(e.target.value)}
                >
                  <MenuItem value="Realtime">Realtime</MenuItem>
                  <MenuItem value="1Min">1Min</MenuItem>
                  <MenuItem value="5Min">5Min</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>SL Monitoring</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>SL Monitoring</InputLabel>
                <Select
                  value={slMonitoring}
                  label="SL Monitoring"
                  onChange={(e) => setSlMonitoring(e.target.value)}
                >
                  <MenuItem value="Realtime">Realtime</MenuItem>
                  <MenuItem value="1Min">1Min</MenuItem>
                  <MenuItem value="5Min">5Min</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>Minute Close Trade Seconds</Typography>
              <TextField
                fullWidth
                size="small"
                value={minuteCloseTradeSeconds}
                onChange={(e) => setMinuteCloseTradeSeconds(e.target.value)}
                sx={{ mb: 2 }}
              />
              
              <Typography variant="subtitle2" gutterBottom>Leg ReExecute Settings Monitoring</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Monitoring</InputLabel>
                <Select
                  value={legReExecuteSettingsMonitoring}
                  label="Monitoring"
                  onChange={(e) => setLegReExecuteSettingsMonitoring(e.target.value)}
                >
                  <MenuItem value="Realtime">Realtime</MenuItem>
                  <MenuItem value="1Min">1Min</MenuItem>
                  <MenuItem value="5Min">5Min</MenuItem>
                </Select>
              </FormControl>
              
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={noReExecuteIfMovedSLToCost}
                    onChange={(e) => setNoReExecuteIfMovedSLToCost(e.target.checked)}
                  />
                }
                label="NO ReExecute if Moved SL To Cost"
              />
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>ReEntry Monitoring</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>ReEntry Monitoring</InputLabel>
                <Select
                  value={reEntryMonitoring}
                  label="ReEntry Monitoring"
                  onChange={(e) => setReEntryMonitoring(e.target.value)}
                >
                  <MenuItem value="Realtime">Realtime</MenuItem>
                  <MenuItem value="1Min">1Min</MenuItem>
                  <MenuItem value="5Min">5Min</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>ReEntry Trigger On</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Trigger On</InputLabel>
                <Select
                  value={reEntryTriggerOn}
                  label="Trigger On"
                  onChange={(e) => setReEntryTriggerOn(e.target.value)}
                >
                  <MenuItem value="None">None</MenuItem>
                  <MenuItem value="OnSL">OnSL</MenuItem>
                  <MenuItem value="OnTarget">OnTarget</MenuItem>
                </Select>
              </FormControl>
              
              <Typography variant="subtitle2" gutterBottom>Order Type</Typography>
              <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                <InputLabel>Order Type</InputLabel>
                <Select
                  value={orderType}
                  label="Order Type"
                  onChange={(e) => setOrderType(e.target.value)}
                >
                  <MenuItem value="MARKET">MARKET</MenuItem>
                  <MenuItem value="LIMIT">LIMIT</MenuItem>
                  <MenuItem value="SL">SL</MenuItem>
                  <MenuItem value="SL-M">SL-M</MenuItem>
                </Select>
              </FormControl>
              
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="subtitle2" gutterBottom>Min Delay In ReEntry (Sec.)</Typography>
                  <TextField
                    fullWidth
                    size="small"
                    value={minDelayInReEntry}
                    onChange={(e) => setMinDelayInReEntry(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Box sx={{ display: 'flex', alignItems: 'center', height: '100%', pt: 3 }}>
                    <FormControlLabel
                      control={
                        <Checkbox 
                          checked={reEntryAtOriginalEntryPrice}
                          onChange={(e) => setReEntryAtOriginalEntryPrice(e.target.checked)}
                        />
                      }
                      label="ReEntry at Original Entry Price"
                    />
                  </Box>
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Dynamic Hedge Tab */}
      <TabPanel value={activeTab} index={5}>
        <Box sx={{ p: 2, bgcolor: '#f5f5f5', borderRadius: 1, mb: 2 }}>
          <Typography variant="body2" color="textSecondary">
            This setting will only be applicable to the Legs where 'Hedge Req' tick is ticked. Here you can set the Dynamic Hedge Settings which will be used to select Hedge Legs for those legs.
            Fill only those Fields which are really required, if the system is unable to find the Hedge Legs, it may stop the execution of Portfolio.
          </Typography>
        </Box>
        
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Mandatory Fields</Typography>
              <Typography variant="body2" gutterBottom>Hedge Distance from ATM</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Minimum"
                    value={hedgeDistanceFromATMMin}
                    onChange={(e) => setHedgeDistanceFromATMMin(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Maximum"
                    value={hedgeDistanceFromATMMax}
                    onChange={(e) => setHedgeDistanceFromATMMax(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Optional Fields</Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Min Premium"
                    value={minPremium}
                    onChange={(e) => setMinPremium(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Max Premium"
                    value={maxPremium}
                    onChange={(e) => setMaxPremium(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="SqOff Leg On"
                    value={sqOffLegOn}
                    onChange={(e) => setSqOffLegOn(e.target.value)}
                  />
                </Grid>
                <Grid item xs={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Unsatisfied Condition Action</InputLabel>
                    <Select
                      value={unsatisfiedConditionAction}
                      label="Unsatisfied Condition Action"
                      onChange={(e) => setUnsatisfiedConditionAction(e.target.value)}
                    >
                      <MenuItem value="IgnoreField">IgnoreField</MenuItem>
                      <MenuItem value="StopExecution">StopExecution</MenuItem>
                      <MenuItem value="SkipLeg">SkipLeg</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
              
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={selectOnly500Strikes}
                    onChange={(e) => setSelectOnly500Strikes(e.target.checked)}
                  />
                }
                label="Select Only 500 strikes for NIFTY & BANKNIFTY"
                sx={{ mt: 2 }}
              />
              
              <Button 
                variant="contained" 
                color="primary" 
                sx={{ mt: 2, float: 'right' }}
              >
                Suggest Legs
              </Button>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Target Settings Tab */}
      <TabPanel value={activeTab} index={6}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <FormControl fullWidth size="small">
              <InputLabel>Target Type</InputLabel>
              <Select
                value={targetType}
                label="Target Type"
                onChange={(e) => setTargetType(e.target.value)}
              >
                <MenuItem value="None">None</MenuItem>
                <MenuItem value="Fixed">Fixed</MenuItem>
                <MenuItem value="Trailing">Trailing</MenuItem>
                <MenuItem value="Combined">Combined</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Portfolio Profit Protection</Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label="If Profit Reaches"
                    value={ifProfitReaches}
                    onChange={(e) => setIfProfitReaches(e.target.value)}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Lock Minimum Profit At"
                    value={lockMinimumProfitAt}
                    onChange={(e) => setLockMinimumProfitAt(e.target.value)}
                  />
                </Grid>
              </Grid>
              
              <Grid container spacing={2} sx={{ mt: 1 }}>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label="For Every Increase In Profit By"
                    value={forEveryIncreaseInProfitBy}
                    onChange={(e) => setForEveryIncreaseInProfitBy(e.target.value)}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Trail Profit By"
                    value={trailProfitBy}
                    onChange={(e) => setTrailProfitBy(e.target.value)}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Stoploss Settings Tab */}
      <TabPanel value={activeTab} index={7}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <FormControl fullWidth size="small">
              <InputLabel>Stoploss Type</InputLabel>
              <Select
                value={stoplossType}
                label="Stoploss Type"
                onChange={(e) => setStoplossType(e.target.value)}
              >
                <MenuItem value="None">None</MenuItem>
                <MenuItem value="Fixed">Fixed</MenuItem>
                <MenuItem value="Trailing">Trailing</MenuItem>
                <MenuItem value="Combined">Combined</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Move SL to Cost Settings</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={moveSLToCost}
                    onChange={(e) => setMoveSLToCost(e.target.checked)}
                  />
                }
                label="Move SL to Cost"
              />
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* Exit Settings Tab */}
      <TabPanel value={activeTab} index={8}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <FormControl fullWidth size="small">
              <InputLabel>Exit Type</InputLabel>
              <Select
                value={exitType}
                label="Exit Type"
                onChange={(e) => setExitType(e.target.value)}
              >
                <MenuItem value="Time">Time-based</MenuItem>
                <MenuItem value="Condition">Condition-based</MenuItem>
                <MenuItem value="Manual">Manual</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          
          {exitType === 'Time' && (
            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                size="small"
                label="Exit Time"
                value={exitTime}
                onChange={(e) => setExitTime(e.target.value)}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton size="small">
                        <HelpOutlineIcon fontSize="small" />
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
          )}
          
          <Grid item xs={12}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Profit/Loss Exit Conditions</Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={exitOnProfit}
                        onChange={(e) => setExitOnProfit(e.target.checked)}
                      />
                    }
                    label="Exit on Profit"
                  />
                  <TextField
                    fullWidth
                    size="small"
                    label="Profit Amount (₹)"
                    value={profitAmount}
                    onChange={(e) => setProfitAmount(e.target.value)}
                    disabled={!exitOnProfit}
                    sx={{ mt: 1 }}
                  />
                </Grid>
                <Grid item xs={12} md={6}>
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={exitOnLoss}
                        onChange={(e) => setExitOnLoss(e.target.checked)}
                      />
                    }
                    label="Exit on Loss"
                  />
                  <TextField
                    fullWidth
                    size="small"
                    label="Loss Amount (₹)"
                    value={lossAmount}
                    onChange={(e) => setLossAmount(e.target.value)}
                    disabled={!exitOnLoss}
                    sx={{ mt: 1 }}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Partial Exit Settings</Typography>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={partialExitEnabled}
                    onChange={(e) => setPartialExitEnabled(e.target.checked)}
                  />
                }
                label="Enable Partial Exit"
              />
              
              <Grid container spacing={2} sx={{ mt: 1 }}>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Exit Percentage"
                    value={partialExitPercentage}
                    onChange={(e) => setPartialExitPercentage(e.target.value)}
                    disabled={!partialExitEnabled}
                    InputProps={{
                      endAdornment: <InputAdornment position="end">%</InputAdornment>,
                    }}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <FormControl fullWidth size="small" disabled={!partialExitEnabled}>
                    <InputLabel>Trigger</InputLabel>
                    <Select
                      value={partialExitTrigger}
                      label="Trigger"
                      onChange={(e) => setPartialExitTrigger(e.target.value)}
                    >
                      <MenuItem value="Profit">Profit</MenuItem>
                      <MenuItem value="Loss">Loss</MenuItem>
                      <MenuItem value="Time">Time</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    size="small"
                    label={partialExitTrigger === 'Time' ? 'Time' : 'Amount (₹)'}
                    value={partialExitValue}
                    onChange={(e) => setPartialExitValue(e.target.value)}
                    disabled={!partialExitEnabled}
                  />
                </Grid>
              </Grid>
            </Paper>
          </Grid>
          
          <Grid item xs={12}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Exit Execution Settings</Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Exit Strategy</InputLabel>
                    <Select
                      value={exitStrategy}
                      label="Exit Strategy"
                      onChange={(e) => setExitStrategy(e.target.value)}
                    >
                      <MenuItem value="CloseAll">Close All Positions</MenuItem>
                      <MenuItem value="CloseProfitable">Close Only Profitable</MenuItem>
                      <MenuItem value="CloseLossMaking">Close Only Loss Making</MenuItem>
                      <MenuItem value="ReversePositions">Reverse Positions</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Exit Order Type</InputLabel>
                    <Select
                      value={exitOrder}
                      label="Exit Order Type"
                      onChange={(e) => setExitOrder(e.target.value)}
                    >
                      <MenuItem value="MARKET">MARKET</MenuItem>
                      <MenuItem value="LIMIT">LIMIT</MenuItem>
                      <MenuItem value="SL">SL</MenuItem>
                      <MenuItem value="SL-M">SL-M</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>
      
      {/* At Broker Tab */}
      <TabPanel value={activeTab} index={9}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={logSLAtBroker}
                    onChange={(e) => setLogSLAtBroker(e.target.checked)}
                  />
                }
                label="Log SL at Broker"
              />
              <Typography variant="caption" display="block" gutterBottom>
                If this is ticked, then for the legs where you have specified the clear SL, the bridge will place the SL orders at the broker end.
                Clear SL means where the bridge can calculate the SL Price which will be required for the SL Order. SL based on Greeks cannot be used to Place SL at Broker.
              </Typography>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={logTargetAtBroker}
                    onChange={(e) => setLogTargetAtBroker(e.target.checked)}
                  />
                }
                label="Log Target at Broker"
              />
              <Typography variant="caption" display="block" gutterBottom>
                If this is ticked, then for the legs where you have specified the clear Target, the bridge will place the Target orders at the broker end.
                Clear Target means where the bridge can calculate the exact Target Price which will be required for the Order. Target based on Greeks cannot be used to Place Target at Broker.
              </Typography>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={legReEntryAtBroker}
                    onChange={(e) => setLegReEntryAtBroker(e.target.checked)}
                  />
                }
                label="Leg ReEntry at Broker"
              />
              <Typography variant="caption" display="block" gutterBottom>
                If this is ticked, then ReEntry (if any specified for the leg) Orders will be sent to the broker in advance.
              </Typography>
              <Typography variant="caption" display="block" gutterBottom>
                This may consume higher Margins, For more details check the documentation.
              </Typography>
            </Paper>
          </Grid>
          
          <Grid item xs={12} md={6}>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <FormControlLabel
                control={
                  <Checkbox 
                    checked={legWaitTradeAtBroker}
                    onChange={(e) => setLegWaitTradeAtBroker(e.target.checked)}
                  />
                }
                label="Leg Wait & Trade at Broker"
              />
              <Typography variant="caption" display="block" gutterBottom>
                If this is ticked, then Leg's Wait and Trade setting under the leg) Orders will be sent to the broker in advance.
              </Typography>
              <Typography variant="caption" display="block" gutterBottom>
                This may consume higher Margins, For more details check the documentation. If ticked, Leg's Wait and Trade setting under Extra Conditions will not work.
              </Typography>
            </Paper>
          </Grid>
        </Grid>
      </TabPanel>

      {/* SL Action Configuration Dialog */}
      <Dialog open={slActionDialogOpen} onClose={handleCloseSLActionDialog} maxWidth="sm" fullWidth>
        <DialogTitle>
          <Typography variant="h6" sx={{ fontSize: '1rem', fontWeight: 'bold' }}>
            Stop Loss Action Configuration
          </Typography>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <Typography variant="body2" sx={{ mb: 2, color: '#666' }}>
              Configure how the system should respond when stop loss is triggered.
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                  <InputLabel>Action Type</InputLabel>
                  <Select
                    value={advancedSLConfig.actionType}
                    label="Action Type"
                    onChange={(e) => handleSLConfigChange('actionType', e.target.value)}
                  >
                    <MenuItem value="SqOff">Square Off Position</MenuItem>
                    <MenuItem value="Trail">Trail Stop Loss</MenuItem>
                    <MenuItem value="None">No Action</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              
              {advancedSLConfig.actionType === 'Trail' && (
                <>
                  <Grid item xs={12} sm={6}>
                    <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                      <InputLabel>Trail Type</InputLabel>
                      <Select
                        value={advancedSLConfig.trailType}
                        label="Trail Type"
                        onChange={(e) => handleSLConfigChange('trailType', e.target.value)}
                      >
                        <MenuItem value="Percentage">Percentage</MenuItem>
                        <MenuItem value="Points">Points</MenuItem>
                        <MenuItem value="Absolute">Absolute Value</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                  
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      size="small"
                      label="Trail Value"
                      value={advancedSLConfig.trailValue}
                      onChange={(e) => handleSLConfigChange('trailValue', e.target.value)}
                      sx={{ mb: 2 }}
                    />
                  </Grid>
                </>
              )}
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedSLConfig.reentryEnabled}
                      onChange={(e) => handleSLConfigChange('reentryEnabled', e.target.checked)}
                    />
                  }
                  label="Enable Re-entry"
                />
              </Grid>
              
              {advancedSLConfig.reentryEnabled && (
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Re-entry Delay (seconds)"
                    value={advancedSLConfig.reentryDelay}
                    onChange={(e) => handleSLConfigChange('reentryDelay', e.target.value)}
                    sx={{ mb: 2 }}
                  />
                </Grid>
              )}
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedSLConfig.notifyUser}
                      onChange={(e) => handleSLConfigChange('notifyUser', e.target.checked)}
                    />
                  }
                  label="Notify User"
                />
              </Grid>
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedSLConfig.executeImmediately}
                      onChange={(e) => handleSLConfigChange('executeImmediately', e.target.checked)}
                    />
                  }
                  label="Execute Immediately"
                />
              </Grid>
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedSLConfig.reversePosition}
                      onChange={(e) => handleSLConfigChange('reversePosition', e.target.checked)}
                    />
                  }
                  label="Reverse Position"
                />
              </Grid>
            </Grid>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseSLActionDialog}>Cancel</Button>
          <Button onClick={handleSaveSLConfig} variant="contained">Save Configuration</Button>
        </DialogActions>
      </Dialog>

      {/* Target Action Configuration Dialog */}
      <Dialog open={targetActionDialogOpen} onClose={handleCloseTargetActionDialog} maxWidth="sm" fullWidth>
        <DialogTitle>
          <Typography variant="h6" sx={{ fontSize: '1rem', fontWeight: 'bold' }}>
            Target Action Configuration
          </Typography>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <Typography variant="body2" sx={{ mb: 2, color: '#666' }}>
              Configure how the system should respond when target is reached.
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                  <InputLabel>Action Type</InputLabel>
                  <Select
                    value={advancedTargetConfig.actionType}
                    label="Action Type"
                    onChange={(e) => handleTargetConfigChange('actionType', e.target.value)}
                  >
                    <MenuItem value="SqOff">Square Off Position</MenuItem>
                    <MenuItem value="Trail">Trail Target</MenuItem>
                    <MenuItem value="None">No Action</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              
              {advancedTargetConfig.actionType === 'Trail' && (
                <>
                  <Grid item xs={12} sm={6}>
                    <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                      <InputLabel>Trail Type</InputLabel>
                      <Select
                        value={advancedTargetConfig.trailType}
                        label="Trail Type"
                        onChange={(e) => handleTargetConfigChange('trailType', e.target.value)}
                      >
                        <MenuItem value="Percentage">Percentage</MenuItem>
                        <MenuItem value="Points">Points</MenuItem>
                        <MenuItem value="Absolute">Absolute Value</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                  
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      size="small"
                      label="Trail Value"
                      value={advancedTargetConfig.trailValue}
                      onChange={(e) => handleTargetConfigChange('trailValue', e.target.value)}
                      sx={{ mb: 2 }}
                    />
                  </Grid>
                </>
              )}
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedTargetConfig.reentryEnabled}
                      onChange={(e) => handleTargetConfigChange('reentryEnabled', e.target.checked)}
                    />
                  }
                  label="Enable Re-entry"
                />
              </Grid>
              
              {advancedTargetConfig.reentryEnabled && (
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Re-entry Delay (seconds)"
                    value={advancedTargetConfig.reentryDelay}
                    onChange={(e) => handleTargetConfigChange('reentryDelay', e.target.value)}
                    sx={{ mb: 2 }}
                  />
                </Grid>
              )}
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedTargetConfig.notifyUser}
                      onChange={(e) => handleTargetConfigChange('notifyUser', e.target.checked)}
                    />
                  }
                  label="Notify User"
                />
              </Grid>
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedTargetConfig.executeImmediately}
                      onChange={(e) => handleTargetConfigChange('executeImmediately', e.target.checked)}
                    />
                  }
                  label="Execute Immediately"
                />
              </Grid>
              
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Checkbox 
                      checked={advancedTargetConfig.reversePosition}
                      onChange={(e) => handleTargetConfigChange('reversePosition', e.target.checked)}
                    />
                  }
                  label="Reverse Position"
                />
              </Grid>
            </Grid>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseTargetActionDialog}>Cancel</Button>
          <Button onClick={handleSaveTargetConfig} variant="contained">Save Configuration</Button>
        </DialogActions>
      </Dialog>

      {/* Feedback Snackbar */}
      <Snackbar
        open={feedbackOpen}
        autoHideDuration={3000}
        onClose={handleFeedbackClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={handleFeedbackClose} severity={feedbackSeverity} sx={{ width: '100%' }}>
          {feedbackMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
};

// TabPanel component
function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`tabpanel-${index}`}
      aria-labelledby={`tab-${index}`}
      {...other}
      style={{ padding: '16px 0' }}
    >
      {value === index && (
        <Box>
          {children}
        </Box>
      )}
    </div>
  );
}

export default Panel3Tabs;
