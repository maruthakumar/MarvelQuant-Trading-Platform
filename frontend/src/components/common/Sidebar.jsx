import React, { useState } from 'react';
import { 
  Box, 
  Drawer, 
  List, 
  ListItem, 
  ListItemButton, 
  ListItemIcon, 
  ListItemText, 
  Toolbar, 
  IconButton, 
  Divider,
  Tooltip
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

// Reduced drawer width to match Quantiply design
const drawerWidth = 220;
const collapsedWidth = 56; // Reduced from 64 to match Quantiply's more compact design

const Sidebar = () => {
  const [open, setOpen] = useState(true);
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
          width: open ? drawerWidth : collapsedWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: open ? drawerWidth : collapsedWidth,
            boxSizing: 'border-box',
            transition: theme => theme.transitions.create('width', {
              easing: theme.transitions.easing.easeInOut, // Changed to easeInOut for smoother transition
              duration: theme.transitions.duration.standard, // Standard duration for better feel
            }),
            overflowX: 'hidden',
            backgroundColor: '#f5f5f5', // Updated to plain light gray color as in Quantiply
            color: '#333', // Updated text color to dark gray for better contrast with light background
            boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)', // Added subtle shadow for depth
            display: 'flex',
            flexDirection: 'column', // Added to support placing toggle at bottom
          },
        }}
      >
        <Toolbar 
          sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: open ? 'flex-start' : 'center', // Changed to flex-start since toggle is now at bottom
            px: [1],
            minHeight: '56px !important', // Reduced from 64px to match Quantiply's more compact design
          }}
        >
          {/* Logo container - only visible when sidebar is expanded */}
          {open && (
            <Box 
              component="img"
              src="/assets/images/MQ Logo-Main.svg"
              alt="MarvelQuant"
              sx={{ 
                height: 32, 
                width: 32,
                ml: 1,
                display: open ? 'block' : 'none'
              }}
            />
          )}
        </Toolbar>
        
        <Divider sx={{ borderColor: 'rgba(0, 0, 0, 0.08)', margin: '0' }} />
        
        <List sx={{ pt: 0.5, pb: 0.5, flexGrow: 1 }}> {/* Added flexGrow to push toggle to bottom */}
          {menuItems.map((item) => (
            <ListItem 
              key={item.text} 
              disablePadding 
              sx={{ 
                display: 'block',
                mb: 0.5, // Added small margin between items
              }}
            >
              <Tooltip title={open ? "" : item.text} placement="right">
                <ListItemButton
                  sx={{
                    minHeight: 40, // Reduced from 48 to match Quantiply's compact design
                    justifyContent: open ? 'initial' : 'center',
                    px: open ? 2 : 1.5, // Reduced padding
                    py: 0.5, // Added vertical padding
                    mx: open ? 1 : 0.5, // Added margin for better spacing
                    borderRadius: '4px', // Added rounded corners like Quantiply
                    backgroundColor: location.pathname === item.path 
                      ? 'rgba(0, 0, 0, 0.05)' 
                      : 'transparent',
                    '&:hover': {
                      backgroundColor: 'rgba(0, 0, 0, 0.05)',
                    },
                    transition: 'background-color 0.2s ease', // Smooth transition for hover effect
                  }}
                  onClick={() => navigate(item.path)}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      mr: open ? 2 : 'auto', // Reduced margin
                      justifyContent: 'center',
                      color: location.pathname === item.path 
                        ? 'rgba(0, 0, 0, 0.9)' 
                        : 'rgba(0, 0, 0, 0.7)', // Updated colors for better contrast with light background
                      fontSize: '1.2rem', // Slightly smaller icons
                    }}
                  >
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText 
                    primary={item.text} 
                    sx={{ 
                      opacity: open ? 1 : 0,
                      '& .MuiTypography-root': {
                        fontWeight: location.pathname === item.path ? 'bold' : 'normal',
                        fontSize: '0.9rem', // Smaller font size
                      }
                    }} 
                  />
                </ListItemButton>
              </Tooltip>
            </ListItem>
          ))}
        </List>
        
        {/* Toggle button moved to bottom */}
        <Box sx={{ 
          p: 1, 
          display: 'flex', 
          justifyContent: open ? 'flex-end' : 'center',
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
