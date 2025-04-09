import React from 'react';
import { 
  Box, 
  Typography, 
  Grid, 
  Paper,
  Container
} from '@mui/material';
import { 
  Description as BookIcon,
  AccountBalance as AccountBalanceIcon,
  Settings as SettingsIcon,
  TrendingUp as TrendingUpIcon,
  Dashboard as DashboardIcon,
  Article as LogsIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const navigate = useNavigate();

  const menuItems = [
    { 
      text: 'OrderBook', 
      icon: <BookIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/orderbook',
      description: 'View and manage your orders'
    },
    { 
      text: 'Positions', 
      icon: <AccountBalanceIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/positions',
      description: 'Track your current positions'
    },
    { 
      text: 'User Settings', 
      icon: <SettingsIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/user-settings',
      description: 'Configure your account settings'
    },
    { 
      text: 'Strategies', 
      icon: <TrendingUpIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/strategies',
      description: 'Manage your trading strategies'
    },
    { 
      text: 'Multi-leg', 
      icon: <DashboardIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/multi-leg',
      description: 'Create and manage multi-leg option strategies'
    },
    { 
      text: 'Logs', 
      icon: <LogsIcon sx={{ fontSize: 40, color: '#1976d2' }} />, 
      path: '/logs',
      description: 'View system logs and activity'
    }
  ];

  return (
    <Container maxWidth="xl" sx={{ mt: 10, ml: { xs: 0, sm: 30 }, width: 'calc(100% - 240px)' }}>
      <Box sx={{ mt: 2, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Dashboard
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Welcome to the MarvelQuant Trading Platform. Use the navigation menu to access different components of the platform.
        </Typography>
      </Box>

      <Grid container spacing={3}>
        {menuItems.map((item, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
            <Paper 
              elevation={0}
              sx={{ 
                p: 3, 
                display: 'flex', 
                flexDirection: 'column', 
                alignItems: 'center',
                cursor: 'pointer',
                transition: 'all 0.2s',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: 1
                },
                height: '100%',
                borderRadius: '4px',
                border: '1px solid #e0e0e0'
              }}
              onClick={() => navigate(item.path)}
            >
              <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
                {item.icon}
              </Box>
              <Typography variant="h6" component="h2" align="center" gutterBottom>
                {item.text}
              </Typography>
              <Typography variant="body2" color="text.secondary" align="center">
                {item.description}
              </Typography>
            </Paper>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default Dashboard;
