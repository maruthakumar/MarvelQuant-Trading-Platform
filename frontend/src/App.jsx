import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box, CssBaseline } from '@mui/material';
import Sidebar from './components/common/Sidebar';
import LiveRateDisplay from './components/common/LiveRateDisplay';
import './App.css';

function App() {
  return (
    <Box sx={{ display: 'flex', height: '100vh' }}>
      <CssBaseline />
      <Sidebar />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
          marginTop: '80px' // Updated margin to account for increased toolbar height
        }}
      >
        <LiveRateDisplay />
        <Box sx={{ p: 2, flexGrow: 1, overflow: 'auto' }}>
          <Outlet />
        </Box>
      </Box>
    </Box>
  );
}

export default App;
