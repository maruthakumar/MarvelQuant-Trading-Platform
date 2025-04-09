import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  AppBar, 
  Toolbar,
  Container,
  Grid
} from '@mui/material';

const LiveRateDisplay = () => {
  // Mock data for market indices - in a real app, this would come from an API
  const [marketData, setMarketData] = useState([
    { name: 'NIFTY', value: '23,250.10', change: '-0.35%', isPositive: false },
    { name: 'BANKNIFTY', value: '51,597.35', change: '+0.49%', isPositive: true },
    { name: 'FINNIFTY', value: '24,724.95', change: '-0.10%', isPositive: false },
    { name: 'MIDCPNIFTY', value: '11,514.10', change: '-0.80%', isPositive: false },
    { name: 'SENSEX', value: '76,295.36', change: '-0.42%', isPositive: false },
    { name: 'BANKEX', value: '59,202.40', change: '+0.09%', isPositive: true }
  ]);

  // Simulate live updates
  useEffect(() => {
    const interval = setInterval(() => {
      setMarketData(prevData => 
        prevData.map(item => {
          // Generate random small changes to simulate live updates
          const randomChange = (Math.random() * 0.2 - 0.1).toFixed(2);
          const currentValue = parseFloat(item.value.replace(',', ''));
          const newValue = (currentValue * (1 + parseFloat(randomChange) / 100)).toFixed(2);
          const formattedValue = newValue.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
          
          const isPositive = randomChange >= 0;
          return {
            ...item,
            value: formattedValue,
            change: `${isPositive ? '+' : ''}${randomChange}%`,
            isPositive
          };
        })
      );
    }, 5000); // Update every 5 seconds

    return () => clearInterval(interval);
  }, []);

  return (
    <AppBar 
      position="static" 
      color="default" 
      elevation={0}
      sx={{ 
        borderBottom: '1px solid #e0e0e0',
        backgroundColor: 'white'
      }}
    >
      <Toolbar variant="dense" sx={{ minHeight: '48px !important' }}>
        <Container maxWidth="xl" sx={{ display: 'flex', alignItems: 'center' }}>
          {/* Logo in top bar - larger and clearer as in Quantiply design */}
          <Box 
            component="img"
            src="/images/MQ-Logo-Main.svg"
            alt="MarvelQuant"
            sx={{ 
              height: 40, 
              width: 40,
              mr: 3,
              display: 'block' // Always display the logo
            }}
          />
          
          <Grid container spacing={2} alignItems="center">
            {marketData.map((item, index) => (
              <Grid item key={index}>
                <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                  <Typography variant="caption" color="textSecondary" sx={{ fontSize: '0.7rem', fontWeight: 'medium' }}>
                    {item.name}
                  </Typography>
                  <Typography variant="body2" fontWeight="bold" sx={{ fontSize: '0.8rem' }}>
                    {item.value}
                  </Typography>
                  <Typography 
                    variant="caption" 
                    sx={{ 
                      color: item.isPositive ? 'success.main' : 'error.main',
                      fontWeight: 'medium',
                      fontSize: '0.7rem'
                    }}
                  >
                    {item.change}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Toolbar>
    </AppBar>
  );
};

export default LiveRateDisplay;
