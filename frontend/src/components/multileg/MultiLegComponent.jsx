import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Paper, Grid } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import PortfolioComponent from './PortfolioComponent';

const MultiLegComponent = () => {
  const navigate = useNavigate();
  const [portfolios, setPortfolios] = useState([]);
  const [showSuccessMessage, setShowSuccessMessage] = useState(false);

  // Handle creating a new portfolio
  const handleNewPortfolio = () => {
    // Navigate to the dedicated new portfolio page instead of showing a dialog
    navigate('/multi-leg/new-portfolio');
  };

  // This would be called when a portfolio is deleted
  const handleDeletePortfolio = (index) => {
    const updatedPortfolios = [...portfolios];
    updatedPortfolios.splice(index, 1);
    setPortfolios(updatedPortfolios);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Multi-Leg Portfolios</Typography>
        <Button 
          variant="contained" 
          startIcon={<AddIcon />} 
          onClick={handleNewPortfolio}
          sx={{ 
            backgroundColor: '#1976d2',
            '&:hover': {
              backgroundColor: '#1565c0',
            }
          }}
        >
          New Portfolio
        </Button>
      </Box>

      {showSuccessMessage && (
        <Paper 
          sx={{ 
            p: 2, 
            mb: 3, 
            backgroundColor: '#e8f5e9',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}
        >
          <Typography variant="body1" color="success.main">
            Portfolio created successfully!
          </Typography>
          <Button 
            size="small" 
            onClick={() => setShowSuccessMessage(false)}
            sx={{ color: '#2e7d32' }}
          >
            Dismiss
          </Button>
        </Paper>
      )}

      {portfolios.length > 0 ? (
        <Grid container spacing={3}>
          {portfolios.map((portfolio, index) => (
            <Grid item xs={12} key={index}>
              <PortfolioComponent 
                portfolio={portfolio} 
                onDelete={() => handleDeletePortfolio(index)} 
              />
            </Grid>
          ))}
        </Grid>
      ) : (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No portfolios yet
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph>
            Create your first portfolio by clicking the "New Portfolio" button above.
          </Typography>
          <Button 
            variant="outlined" 
            startIcon={<AddIcon />} 
            onClick={handleNewPortfolio}
            sx={{ mt: 2 }}
          >
            Create Portfolio
          </Button>
        </Paper>
      )}
    </Box>
  );
};

export default MultiLegComponent;
