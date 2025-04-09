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
  Avatar,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Checkbox,
  Tooltip,
  Alert,
  Snackbar
} from '@mui/material';
import SaveIcon from '@mui/icons-material/Save';
import SecurityIcon from '@mui/icons-material/Security';
import NotificationsIcon from '@mui/icons-material/Notifications';
import AccountBalanceIcon from '@mui/icons-material/AccountBalance';
import ApiIcon from '@mui/icons-material/Api';
import SettingsIcon from '@mui/icons-material/Settings';
import PersonIcon from '@mui/icons-material/Person';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RefreshIcon from '@mui/icons-material/Refresh';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import WarningIcon from '@mui/icons-material/Warning';
import LogsPanel from '../logs/LogsPanel';

const UserSettings = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [userProfile, setUserProfile] = useState({
    name: 'Admin User',
    email: 'admin@marvelquant.com',
    phone: '+91 9876543210',
    role: 'Administrator',
    lastLogin: '2025-04-07 09:15:30',
    twoFactorEnabled: true,
    emailNotifications: true,
    smsNotifications: false,
    pushNotifications: true,
    defaultBroker: 'Zerodha',
    apiKey: 'XXXXXXXXXXXX',
    apiSecret: '••••••••••••',
    theme: 'light',
    dataRefreshInterval: 5,
    orderConfirmation: true,
    autoSquareOff: true,
    autoSquareOffTime: '15:15'
  });
  
  const [openPasswordDialog, setOpenPasswordDialog] = useState(false);
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  
  const [brokerAccounts, setBrokerAccounts] = useState([
    { 
      id: 1, 
      enabled: true, 
      broker: 'Zerodha', 
      userId: 'ZM9343', 
      apiKey: 'bpt7scgatyn8nvq1', 
      apiSecret: 'ty61tnuqqpqp3p8h...', 
      availableMargin: '0.00',
      mtm: '0.00',
      status: 'Allowed',
      squareOffTime: '00:00:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    },
    { 
      id: 2, 
      enabled: true, 
      broker: 'NIFSG100X', 
      userId: 'SIM2', 
      apiKey: 'APITest', 
      apiSecret: '••••••••••••', 
      availableMargin: '100,000,000.00',
      mtm: '79,485.00',
      status: 'Allowed',
      squareOffTime: '00:00:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    },
    { 
      id: 3, 
      enabled: false, 
      broker: 'JOSHUA_FIN/ASIA', 
      userId: 'FA161611', 
      apiKey: 'Finvasia', 
      apiSecret: '4c295bd39bc96d...', 
      availableMargin: '0.00',
      mtm: '0.00',
      status: 'Allowed',
      squareOffTime: '00:00:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    },
    { 
      id: 4, 
      enabled: false, 
      broker: 'JOSHUA', 
      userId: 'DLT1182', 
      apiKey: 'JainamBroking', 
      apiSecret: '0fc36e21efa0094f...', 
      availableMargin: '0.00',
      mtm: '0.00',
      status: 'Allowed',
      squareOffTime: '00:00:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    },
    { 
      id: 5, 
      enabled: true, 
      broker: 'SIMULATED1', 
      userId: 'SIM1', 
      apiKey: 'APITest', 
      apiSecret: '••••••••••••', 
      availableMargin: '100,000,000.00',
      mtm: '298,018.00',
      status: 'Allowed',
      squareOffTime: '23:55:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    }
  ]);
  
  const [openBrokerDialog, setOpenBrokerDialog] = useState(false);
  const [selectedBroker, setSelectedBroker] = useState(null);
  const [editMode, setEditMode] = useState(false);
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success'
  });
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  const handleInputChange = (field, value) => {
    setUserProfile({
      ...userProfile,
      [field]: value
    });
  };
  
  const handleSaveProfile = () => {
    // In a real implementation, this would save the profile data to the backend
    console.log('Saving profile:', userProfile);
    // Show a success message or notification
    setSnackbar({
      open: true,
      message: 'Profile saved successfully',
      severity: 'success'
    });
  };
  
  const handleChangePassword = () => {
    setOpenPasswordDialog(true);
    setPasswordData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
  };
  
  const handleClosePasswordDialog = () => {
    setOpenPasswordDialog(false);
  };
  
  const handleSavePassword = () => {
    // Validate passwords
    if (!passwordData.currentPassword) {
      setSnackbar({
        open: true,
        message: 'Current password is required',
        severity: 'error'
      });
      return;
    }
    
    if (!passwordData.newPassword) {
      setSnackbar({
        open: true,
        message: 'New password is required',
        severity: 'error'
      });
      return;
    }
    
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setSnackbar({
        open: true,
        message: 'New password and confirm password do not match',
        severity: 'error'
      });
      return;
    }
    
    // In a real implementation, this would send the password change request to the backend
    console.log('Changing password');
    
    // Close dialog and show success message
    setOpenPasswordDialog(false);
    setSnackbar({
      open: true,
      message: 'Password changed successfully',
      severity: 'success'
    });
  };
  
  const handleAddBroker = () => {
    setSelectedBroker({
      id: brokerAccounts.length + 1,
      enabled: true,
      broker: '',
      userId: '',
      apiKey: '',
      apiSecret: '',
      availableMargin: '0.00',
      mtm: '0.00',
      status: 'Allowed',
      squareOffTime: '00:00:00',
      enableNRML: true,
      enableCNC: true,
      orderType: 'MARKET'
    });
    setEditMode(true);
    setOpenBrokerDialog(true);
  };
  
  const handleEditBroker = (broker) => {
    setSelectedBroker(broker);
    setEditMode(true);
    setOpenBrokerDialog(true);
  };
  
  const handleDeleteBroker = (brokerId) => {
    if (window.confirm('Are you sure you want to delete this broker account?')) {
      const updatedBrokers = brokerAccounts.filter(broker => broker.id !== brokerId);
      setBrokerAccounts(updatedBrokers);
      setSnackbar({
        open: true,
        message: 'Broker account deleted successfully',
        severity: 'success'
      });
    }
  };
  
  const handleCloseBrokerDialog = () => {
    setOpenBrokerDialog(false);
  };
  
  const handleSaveBroker = () => {
    if (!selectedBroker.broker) {
      setSnackbar({
        open: true,
        message: 'Broker name is required',
        severity: 'error'
      });
      return;
    }
    
    if (!selectedBroker.userId) {
      setSnackbar({
        open: true,
        message: 'User ID is required',
        severity: 'error'
      });
      return;
    }
    
    if (!selectedBroker.apiKey) {
      setSnackbar({
        open: true,
        message: 'API Key is required',
        severity: 'error'
      });
      return;
    }
    
    // In a real implementation, this would save the broker data to the backend
    console.log('Saving broker:', selectedBroker);
    
    // Update the brokers list
    if (editMode) {
      const updatedBrokers = brokerAccounts.map(broker => 
        broker.id === selectedBroker.id ? selectedBroker : broker
      );
      setBrokerAccounts(updatedBrokers);
    } else {
      setBrokerAccounts([...brokerAccounts, selectedBroker]);
    }
    
    // Close dialog and show success message
    setOpenBrokerDialog(false);
    setSnackbar({
      open: true,
      message: `Broker account ${editMode ? 'updated' : 'added'} successfully`,
      severity: 'success'
    });
  };
  
  const handleBrokerInputChange = (field, value) => {
    setSelectedBroker({
      ...selectedBroker,
      [field]: value
    });
  };
  
  const handleToggleBrokerStatus = (brokerId) => {
    const updatedBrokers = brokerAccounts.map(broker => 
      broker.id === brokerId ? { ...broker, enabled: !broker.enabled } : broker
    );
    setBrokerAccounts(updatedBrokers);
  };
  
  const handleCloseSnackbar = () => {
    setSnackbar({
      ...snackbar,
      open: false
    });
  };
  
  const handleRefreshBrokers = () => {
    // In a real implementation, this would fetch the latest broker data from the backend
    console.log('Refreshing broker accounts');
    setSnackbar({
      open: true,
      message: 'Broker accounts refreshed',
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
        id={`settings-tabpanel-${index}`}
        aria-labelledby={`settings-tab-${index}`}
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
      <Typography variant="h6" sx={{ mb: 2 }}>User Settings</Typography>
      
      <Grid container spacing={2} sx={{ flexGrow: 1, mb: 2 }}>
        {/* Left sidebar for navigation */}
        <Grid item xs={12} md={3}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mb: 3 }}>
              <Avatar 
                sx={{ 
                  width: 80, 
                  height: 80, 
                  mb: 2,
                  bgcolor: 'primary.main'
                }}
              >
                <PersonIcon sx={{ fontSize: 40 }} />
              </Avatar>
              <Typography variant="h6">{userProfile.name}</Typography>
              <Typography variant="body2" color="textSecondary">{userProfile.role}</Typography>
            </Box>
            
            <Divider sx={{ mb: 2 }} />
            
            <List component="nav" sx={{ width: '100%' }}>
              <ListItem 
                button 
                selected={activeTab === 0}
                onClick={(e) => handleTabChange(e, 0)}
              >
                <ListItemIcon>
                  <PersonIcon />
                </ListItemIcon>
                <ListItemText primary="Profile" />
              </ListItem>
              
              <ListItem 
                button 
                selected={activeTab === 1}
                onClick={(e) => handleTabChange(e, 1)}
              >
                <ListItemIcon>
                  <SecurityIcon />
                </ListItemIcon>
                <ListItemText primary="Security" />
              </ListItem>
              
              <ListItem 
                button 
                selected={activeTab === 2}
                onClick={(e) => handleTabChange(e, 2)}
              >
                <ListItemIcon>
                  <NotificationsIcon />
                </ListItemIcon>
                <ListItemText primary="Notifications" />
              </ListItem>
              
              <ListItem 
                button 
                selected={activeTab === 3}
                onClick={(e) => handleTabChange(e, 3)}
              >
                <ListItemIcon>
                  <AccountBalanceIcon />
                </ListItemIcon>
                <ListItemText primary="Broker Settings" />
              </ListItem>
              
              <ListItem 
                button 
                selected={activeTab === 4}
                onClick={(e) => handleTabChange(e, 4)}
              >
                <ListItemIcon>
                  <ApiIcon />
                </ListItemIcon>
                <ListItemText primary="API Configuration" />
              </ListItem>
              
              <ListItem 
                button 
                selected={activeTab === 5}
                onClick={(e) => handleTabChange(e, 5)}
              >
                <ListItemIcon>
                  <SettingsIcon />
                </ListItemIcon>
                <ListItemText primary="Preferences" />
              </ListItem>
            </List>
          </Paper>
        </Grid>
        
        {/* Right side: Settings content */}
        <Grid item xs={12} md={9}>
          <Paper sx={{ p: 2, height: '100%' }}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider', display: { xs: 'block', md: 'none' } }}>
              <Tabs 
                value={activeTab} 
                onChange={handleTabChange}
                variant="scrollable"
                scrollButtons="auto"
              >
                <Tab label="Profile" />
                <Tab label="Security" />
                <Tab label="Notifications" />
                <Tab label="Broker Settings" />
                <Tab label="API Configuration" />
                <Tab label="Preferences" />
              </Tabs>
            </Box>
            
            {/* Tab 1: Profile */}
            <TabPanel value={activeTab} index={0}>
              <Typography variant="h6" gutterBottom>Profile Information</Typography>
              <Divider sx={{ mb: 3 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Full Name"
                    value={userProfile.name}
                    onChange={(e) => handleInputChange('name', e.target.value)}
                    sx={{ mb: 2 }}
                  />
                  
                  <TextField
                    fullWidth
                    label="Email"
                    type="email"
                    value={userProfile.email}
                    onChange={(e) => handleInputChange('email', e.target.value)}
                    sx={{ mb: 2 }}
                  />
                  
                  <TextField
                    fullWidth
                    label="Phone"
                    value={userProfile.phone}
                    onChange={(e) => handleInputChange('phone', e.target.value)}
                    sx={{ mb: 2 }}
                  />
                </Grid>
                
                <Grid item xs={12} sm={6}>
                  <FormControl fullWidth sx={{ mb: 2 }}>
                    <InputLabel>Role</InputLabel>
                    <Select
                      value={userProfile.role}
                      label="Role"
                      onChange={(e) => handleInputChange('role', e.target.value)}
                      disabled
                    >
                      <MenuItem value="Administrator">Administrator</MenuItem>
                      <MenuItem value="Trader">Trader</MenuItem>
                      <MenuItem value="Analyst">Analyst</MenuItem>
                      <MenuItem value="Viewer">Viewer</MenuItem>
                    </Select>
                  </FormControl>
                  
                  <TextField
                    fullWidth
                    label="Last Login"
                    value={userProfile.lastLogin}
                    disabled
                    sx={{ mb: 2 }}
                  />
                </Grid>
                
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                    <Button 
                      variant="contained" 
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save Profile
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </TabPanel>
            
            {/* Tab 2: Security */}
            <TabPanel value={activeTab} index={1}>
              <Typography variant="h6" gutterBottom>Security Settings</Typography>
              <Divider sx={{ mb: 3 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Password</Typography>
                    <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                      It's a good idea to use a strong password that you don't use elsewhere
                    </Typography>
                    <Button 
                      variant="outlined" 
                      onClick={handleChangePassword}
                    >
                      Change Password
                    </Button>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Two-Factor Authentication</Typography>
                    <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                      Add an extra layer of security to your account
                    </Typography>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.twoFactorEnabled}
                          onChange={(e) => handleInputChange('twoFactorEnabled', e.target.checked)}
                        />
                      }
                      label="Enable Two-Factor Authentication"
                    />
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle1" gutterBottom>Session Management</Typography>
                    <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                      Manage your active sessions and sign out from other devices
                    </Typography>
                    <Button 
                      variant="outlined" 
                      color="error"
                    >
                      Sign Out From All Other Devices
                    </Button>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                    <Button 
                      variant="contained" 
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save Security Settings
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </TabPanel>
            
            {/* Tab 3: Notifications */}
            <TabPanel value={activeTab} index={2}>
              <Typography variant="h6" gutterBottom>Notification Settings</Typography>
              <Divider sx={{ mb: 3 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Email Notifications</Typography>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.emailNotifications}
                          onChange={(e) => handleInputChange('emailNotifications', e.target.checked)}
                        />
                      }
                      label="Receive Email Notifications"
                    />
                    
                    <Box sx={{ ml: 4, mt: 1 }}>
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.emailNotifications} />}
                        label="Order Execution"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.emailNotifications} />}
                        label="Strategy Alerts"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.emailNotifications} />}
                        label="P&L Updates"
                      />
                      <FormControlLabel
                        control={<Switch checked={false} disabled={!userProfile.emailNotifications} />}
                        label="Market News"
                      />
                    </Box>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>SMS Notifications</Typography>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.smsNotifications}
                          onChange={(e) => handleInputChange('smsNotifications', e.target.checked)}
                        />
                      }
                      label="Receive SMS Notifications"
                    />
                    
                    <Box sx={{ ml: 4, mt: 1 }}>
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.smsNotifications} />}
                        label="Order Execution"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.smsNotifications} />}
                        label="Strategy Alerts"
                      />
                      <FormControlLabel
                        control={<Switch checked={false} disabled={!userProfile.smsNotifications} />}
                        label="P&L Updates"
                      />
                    </Box>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle1" gutterBottom>Push Notifications</Typography>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.pushNotifications}
                          onChange={(e) => handleInputChange('pushNotifications', e.target.checked)}
                        />
                      }
                      label="Receive Push Notifications"
                    />
                    
                    <Box sx={{ ml: 4, mt: 1 }}>
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.pushNotifications} />}
                        label="Order Execution"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.pushNotifications} />}
                        label="Strategy Alerts"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.pushNotifications} />}
                        label="P&L Updates"
                      />
                      <FormControlLabel
                        control={<Switch checked={true} disabled={!userProfile.pushNotifications} />}
                        label="Market News"
                      />
                    </Box>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                    <Button 
                      variant="contained" 
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save Notification Settings
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </TabPanel>
            
            {/* Tab 4: Broker Settings */}
            <TabPanel value={activeTab} index={3}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6">Broker Accounts</Typography>
                <Box>
                  <Button 
                    variant="contained" 
                    startIcon={<AddIcon />}
                    onClick={handleAddBroker}
                    size="small"
                    sx={{ mr: 1 }}
                  >
                    Add Broker
                  </Button>
                  <IconButton size="small" onClick={handleRefreshBrokers}>
                    <RefreshIcon />
                  </IconButton>
                </Box>
              </Box>
              <Divider sx={{ mb: 3 }} />
              
              <TableContainer sx={{ maxHeight: 440 }}>
                <Table stickyHeader size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell padding="checkbox">
                        <Typography variant="subtitle2">Enabled</Typography>
                      </TableCell>
                      <TableCell>Delete</TableCell>
                      <TableCell>Logout</TableCell>
                      <TableCell>Manual Square Off</TableCell>
                      <TableCell>LoggedIn</TableCell>
                      <TableCell>MTM (All)</TableCell>
                      <TableCell>Available Margin</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>User Alias</TableCell>
                      <TableCell>User ID</TableCell>
                      <TableCell>Broker</TableCell>
                      <TableCell>API Key</TableCell>
                      <TableCell>API Secret</TableCell>
                      <TableCell>Historical API</TableCell>
                      <TableCell>SqOff Time</TableCell>
                      <TableCell>Enable NRML SqOff</TableCell>
                      <TableCell>Enable CNC SqOff</TableCell>
                      <TableCell>SqOff Order Type</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {brokerAccounts.map((broker) => (
                      <TableRow 
                        key={broker.id}
                        hover
                        sx={{ 
                          bgcolor: broker.enabled ? 'rgba(76, 175, 80, 0.05)' : 'rgba(244, 67, 54, 0.05)'
                        }}
                      >
                        <TableCell padding="checkbox">
                          <Checkbox
                            checked={broker.enabled}
                            onChange={() => handleToggleBrokerStatus(broker.id)}
                          />
                        </TableCell>
                        <TableCell>
                          <IconButton 
                            size="small" 
                            color="error"
                            onClick={() => handleDeleteBroker(broker.id)}
                          >
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </TableCell>
                        <TableCell>
                          <IconButton size="small" color="error">
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </TableCell>
                        <TableCell>
                          <IconButton size="small" color="primary">
                            <CheckCircleIcon fontSize="small" />
                          </IconButton>
                        </TableCell>
                        <TableCell>
                          {broker.enabled ? (
                            <IconButton size="small" color="success">
                              <CheckCircleIcon fontSize="small" />
                            </IconButton>
                          ) : (
                            <IconButton size="small" color="error">
                              <ErrorIcon fontSize="small" />
                            </IconButton>
                          )}
                        </TableCell>
                        <TableCell>{broker.mtm}</TableCell>
                        <TableCell>{broker.availableMargin}</TableCell>
                        <TableCell>{broker.status}</TableCell>
                        <TableCell>{broker.broker}</TableCell>
                        <TableCell>{broker.userId}</TableCell>
                        <TableCell>{broker.broker}</TableCell>
                        <TableCell>{broker.apiKey}</TableCell>
                        <TableCell>{broker.apiSecret}</TableCell>
                        <TableCell>
                          <Checkbox checked={true} disabled />
                        </TableCell>
                        <TableCell>{broker.squareOffTime}</TableCell>
                        <TableCell>
                          <Checkbox checked={broker.enableNRML} disabled />
                        </TableCell>
                        <TableCell>
                          <Checkbox checked={broker.enableCNC} disabled />
                        </TableCell>
                        <TableCell>{broker.orderType}</TableCell>
                        <TableCell>
                          <IconButton 
                            size="small" 
                            color="primary"
                            onClick={() => handleEditBroker(broker)}
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))}
                    {brokerAccounts.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={19} align="center" sx={{ py: 3 }}>
                          <Typography variant="body2" color="textSecondary">
                            No broker accounts found
                          </Typography>
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </TableContainer>
              
              <Box sx={{ mt: 3 }}>
                <Typography variant="subtitle2" gutterBottom>Notes:</Typography>
                <Typography variant="body2" color="textSecondary">
                  1. Bridge never collects user id or related details. Any provided ID, Password or any other sensitive information will be saved only in user's computer with encryption.
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  2. Password and Pin is only required if you have selected for Auto Login. Auto login internally fills user details in browser for easy login. It is totally optional feature.
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  3. If you are facing Login issue with Zerodha, AliceBlue, Upstox then just Un-Tick the Auto Login and then proceed with Manual Login.
                </Typography>
              </Box>
            </TabPanel>
            
            {/* Tab 5: API Configuration */}
            <TabPanel value={activeTab} index={4}>
              <Typography variant="h6" gutterBottom>API Configuration</Typography>
              <Divider sx={{ mb: 3 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Default Broker API</Typography>
                    <FormControl fullWidth sx={{ mb: 2 }}>
                      <InputLabel>Default Broker</InputLabel>
                      <Select
                        value={userProfile.defaultBroker}
                        label="Default Broker"
                        onChange={(e) => handleInputChange('defaultBroker', e.target.value)}
                      >
                        {brokerAccounts.map(broker => (
                          <MenuItem key={broker.id} value={broker.broker}>{broker.broker}</MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                    
                    <TextField
                      fullWidth
                      label="API Key"
                      value={userProfile.apiKey}
                      onChange={(e) => handleInputChange('apiKey', e.target.value)}
                      sx={{ mb: 2 }}
                    />
                    
                    <TextField
                      fullWidth
                      label="API Secret"
                      type="password"
                      value={userProfile.apiSecret}
                      onChange={(e) => handleInputChange('apiSecret', e.target.value)}
                      sx={{ mb: 2 }}
                    />
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>API Rate Limits</Typography>
                    <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                      Configure API rate limits to avoid hitting broker's API limits
                    </Typography>
                    
                    <Grid container spacing={2}>
                      <Grid item xs={12} sm={6}>
                        <TextField
                          fullWidth
                          label="Max API Calls Per Minute"
                          type="number"
                          defaultValue={60}
                          sx={{ mb: 2 }}
                        />
                      </Grid>
                      
                      <Grid item xs={12} sm={6}>
                        <TextField
                          fullWidth
                          label="Max Order API Calls Per Minute"
                          type="number"
                          defaultValue={10}
                          sx={{ mb: 2 }}
                        />
                      </Grid>
                    </Grid>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle1" gutterBottom>WebSocket Configuration</Typography>
                    <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                      Configure WebSocket settings for real-time data
                    </Typography>
                    
                    <FormControlLabel
                      control={<Switch defaultChecked />}
                      label="Enable WebSocket Connection"
                    />
                    
                    <FormControlLabel
                      control={<Switch defaultChecked />}
                      label="Auto Reconnect on Disconnection"
                    />
                    
                    <TextField
                      fullWidth
                      label="Reconnect Interval (seconds)"
                      type="number"
                      defaultValue={5}
                      sx={{ mt: 2 }}
                    />
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                    <Button 
                      variant="contained" 
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save API Settings
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </TabPanel>
            
            {/* Tab 6: Preferences */}
            <TabPanel value={activeTab} index={5}>
              <Typography variant="h6" gutterBottom>Preferences</Typography>
              <Divider sx={{ mb: 3 }} />
              
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Theme Settings</Typography>
                    <FormControl fullWidth sx={{ mb: 2 }}>
                      <InputLabel>Theme</InputLabel>
                      <Select
                        value={userProfile.theme}
                        label="Theme"
                        onChange={(e) => handleInputChange('theme', e.target.value)}
                      >
                        <MenuItem value="light">Light</MenuItem>
                        <MenuItem value="dark">Dark</MenuItem>
                        <MenuItem value="system">System Default</MenuItem>
                      </Select>
                    </FormControl>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2, mb: 3 }}>
                    <Typography variant="subtitle1" gutterBottom>Data Settings</Typography>
                    <FormControl fullWidth sx={{ mb: 2 }}>
                      <InputLabel>Data Refresh Interval (seconds)</InputLabel>
                      <Select
                        value={userProfile.dataRefreshInterval}
                        label="Data Refresh Interval (seconds)"
                        onChange={(e) => handleInputChange('dataRefreshInterval', e.target.value)}
                      >
                        <MenuItem value={1}>1 second</MenuItem>
                        <MenuItem value={2}>2 seconds</MenuItem>
                        <MenuItem value={5}>5 seconds</MenuItem>
                        <MenuItem value={10}>10 seconds</MenuItem>
                        <MenuItem value={30}>30 seconds</MenuItem>
                        <MenuItem value={60}>1 minute</MenuItem>
                      </Select>
                    </FormControl>
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle1" gutterBottom>Order Settings</Typography>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.orderConfirmation}
                          onChange={(e) => handleInputChange('orderConfirmation', e.target.checked)}
                        />
                      }
                      label="Show Order Confirmation Dialog"
                    />
                    
                    <FormControlLabel
                      control={
                        <Switch
                          checked={userProfile.autoSquareOff}
                          onChange={(e) => handleInputChange('autoSquareOff', e.target.checked)}
                        />
                      }
                      label="Enable Auto Square Off"
                    />
                    
                    <TextField
                      fullWidth
                      label="Auto Square Off Time"
                      type="time"
                      value={userProfile.autoSquareOffTime}
                      onChange={(e) => handleInputChange('autoSquareOffTime', e.target.value)}
                      sx={{ mt: 2 }}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      inputProps={{
                        step: 300, // 5 min
                      }}
                    />
                  </Paper>
                </Grid>
                
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                    <Button 
                      variant="contained" 
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save Preferences
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </TabPanel>
          </Paper>
        </Grid>
      </Grid>
      
      {/* Password Change Dialog */}
      <Dialog open={openPasswordDialog} onClose={handleClosePasswordDialog}>
        <DialogTitle>Change Password</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Current Password"
            type="password"
            fullWidth
            value={passwordData.currentPassword}
            onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
            sx={{ mb: 2, mt: 1 }}
          />
          <TextField
            margin="dense"
            label="New Password"
            type="password"
            fullWidth
            value={passwordData.newPassword}
            onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Confirm New Password"
            type="password"
            fullWidth
            value={passwordData.confirmPassword}
            onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClosePasswordDialog}>Cancel</Button>
          <Button onClick={handleSavePassword} variant="contained">Change Password</Button>
        </DialogActions>
      </Dialog>
      
      {/* Broker Dialog */}
      <Dialog open={openBrokerDialog} onClose={handleCloseBrokerDialog} maxWidth="md" fullWidth>
        <DialogTitle>{editMode ? 'Edit Broker Account' : 'Add Broker Account'}</DialogTitle>
        <DialogContent>
          {selectedBroker && (
            <Grid container spacing={2} sx={{ mt: 1 }}>
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Broker Name"
                  value={selectedBroker.broker}
                  onChange={(e) => handleBrokerInputChange('broker', e.target.value)}
                  sx={{ mb: 2 }}
                />
                
                <TextField
                  fullWidth
                  label="User ID"
                  value={selectedBroker.userId}
                  onChange={(e) => handleBrokerInputChange('userId', e.target.value)}
                  sx={{ mb: 2 }}
                />
                
                <TextField
                  fullWidth
                  label="API Key"
                  value={selectedBroker.apiKey}
                  onChange={(e) => handleBrokerInputChange('apiKey', e.target.value)}
                  sx={{ mb: 2 }}
                />
                
                <TextField
                  fullWidth
                  label="API Secret"
                  type="password"
                  value={selectedBroker.apiSecret}
                  onChange={(e) => handleBrokerInputChange('apiSecret', e.target.value)}
                  sx={{ mb: 2 }}
                />
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Status</InputLabel>
                  <Select
                    value={selectedBroker.status}
                    label="Status"
                    onChange={(e) => handleBrokerInputChange('status', e.target.value)}
                  >
                    <MenuItem value="Allowed">Allowed</MenuItem>
                    <MenuItem value="Blocked">Blocked</MenuItem>
                  </Select>
                </FormControl>
                
                <TextField
                  fullWidth
                  label="Square Off Time"
                  type="time"
                  value={selectedBroker.squareOffTime}
                  onChange={(e) => handleBrokerInputChange('squareOffTime', e.target.value)}
                  sx={{ mb: 2 }}
                  InputLabelProps={{
                    shrink: true,
                  }}
                />
                
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Order Type</InputLabel>
                  <Select
                    value={selectedBroker.orderType}
                    label="Order Type"
                    onChange={(e) => handleBrokerInputChange('orderType', e.target.value)}
                  >
                    <MenuItem value="MARKET">MARKET</MenuItem>
                    <MenuItem value="LIMIT">LIMIT</MenuItem>
                    <MenuItem value="SL">STOP LOSS</MenuItem>
                    <MenuItem value="SL-M">STOP LOSS MARKET</MenuItem>
                  </Select>
                </FormControl>
                
                <Box>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={selectedBroker.enableNRML}
                        onChange={(e) => handleBrokerInputChange('enableNRML', e.target.checked)}
                      />
                    }
                    label="Enable NRML Square Off"
                  />
                  
                  <FormControlLabel
                    control={
                      <Switch
                        checked={selectedBroker.enableCNC}
                        onChange={(e) => handleBrokerInputChange('enableCNC', e.target.checked)}
                      />
                    }
                    label="Enable CNC Square Off"
                  />
                </Box>
              </Grid>
            </Grid>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseBrokerDialog}>Cancel</Button>
          <Button onClick={handleSaveBroker} variant="contained">Save</Button>
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

export default UserSettings;
