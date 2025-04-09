import React, { useState } from 'react';
import { 
  Box, 
  Drawer, 
  List, 
  ListItem, 
  ListItemButton, 
  ListItemIcon, 
  Divider,
  Tooltip,
  IconButton
} from '@mui/material';
import { 
  Dashboard as DashboardIcon,
  Settings as SettingsIcon,
  TrendingUp as TrendingUpIcon,
  Book as BookIcon,
  AccountBalance as AccountBalanceIcon,
  Menu as MenuIcon,
  ChevronLeft as ChevronLeftIcon
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

// Narrow sidebar width to match Quantiply design
const collapsedWidth = 56; // Icon-only sidebar like Quantiply

const Sidebar = () => {
  const [open, setOpen] = useState(false); // Default to collapsed state like Quantiply
  const navigate = useNavigate();
  const location = useLocation();

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const menuItems = [
    { text: 'OrderBook', icon: <BookIcon />, path: '/orderbook' },
    { text: 'Positions', icon: <AccountBalanceIcon />, path: '/positions' },
    { text: 'User Settings', icon: <SettingsIcon />, path: '/user-settings' },
    { text: 'Strategies', icon: <TrendingUpIcon />, path: '/strategies' },
    { text: 'Multi-leg', icon: <DashboardIcon />, path: '/multi-leg' },
  ];

  return (
    <Box sx={{ display: 'flex' }}>
      <Drawer
        variant="permanent"
        sx={{
          width: collapsedWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: collapsedWidth,
            boxSizing: 'border-box',
            overflowX: 'hidden',
            backgroundColor: '#f5f5f5', // Light background like Quantiply
            color: '#333',
            boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)',
            display: 'flex',
            flexDirection: 'column',
          },
        }}
      >
        <Box sx={{ height: '48px' }} /> {/* Empty space at top to align with header */}
        
        <Divider sx={{ borderColor: 'rgba(0, 0, 0, 0.08)', margin: '0' }} />
        
        <List sx={{ pt: 0.5, pb: 0.5, flexGrow: 1 }}>
          {menuItems.map((item) => (
            <ListItem 
              key={item.text} 
              disablePadding 
              sx={{ 
                display: 'block',
                mb: 0.5,
              }}
            >
              <Tooltip title={item.text} placement="right">
                <ListItemButton
                  sx={{
                    minHeight: 40,
                    justifyContent: 'center',
                    px: 1.5,
                    py: 0.5,
                    mx: 0.5,
                    borderRadius: '4px',
                    backgroundColor: location.pathname === item.path 
                      ? 'rgba(0, 0, 0, 0.05)' 
                      : 'transparent',
                    '&:hover': {
                      backgroundColor: 'rgba(0, 0, 0, 0.05)',
                    },
                    transition: 'background-color 0.2s ease',
                  }}
                  onClick={() => navigate(item.path)}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      justifyContent: 'center',
                      color: location.pathname === item.path 
                        ? 'rgba(0, 0, 0, 0.9)' 
                        : 'rgba(0, 0, 0, 0.7)',
                      fontSize: '1.2rem',
                    }}
                  >
                    {item.icon}
                  </ListItemIcon>
                </ListItemButton>
              </Tooltip>
            </ListItem>
          ))}
        </List>
        
        {/* Toggle button at bottom */}
        <Box sx={{ 
          p: 1, 
          display: 'flex', 
          justifyContent: 'center',
          borderTop: '1px solid rgba(0, 0, 0, 0.08)'
        }}>
          <IconButton 
            onClick={handleDrawerToggle} 
            sx={{ 
              color: '#333',
              backgroundColor: 'rgba(0, 0, 0, 0.03)',
              '&:hover': {
                backgroundColor: 'rgba(0, 0, 0, 0.08)',
              },
              padding: '4px',
              borderRadius: '4px',
            }}
          >
            {open ? <ChevronLeftIcon /> : <MenuIcon />}
          </IconButton>
        </Box>
      </Drawer>
    </Box>
  );
};

export default Sidebar;
