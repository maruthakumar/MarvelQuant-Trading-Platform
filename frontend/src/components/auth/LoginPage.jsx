import React, { useState } from 'react';
import { 
  Box, 
  Container, 
  TextField, 
  Button, 
  Typography, 
  Paper,
  InputAdornment,
  Link,
  CssBaseline,
  Alert,
  Snackbar
} from '@mui/material';
import { useNavigate } from 'react-router-dom';

const LoginPage = () => {
  const [mobileNumber, setMobileNumber] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [showError, setShowError] = useState(false);
  const navigate = useNavigate();

  const handleLogin = (e) => {
    e.preventDefault();
    console.log('Login attempt with:', mobileNumber, password);
    
    // Check for hardcoded credentials
    if (mobileNumber === '9986666444' && password === 'password123') {
      // Set authentication in localStorage
      localStorage.setItem('isAuthenticated', 'true');
      localStorage.setItem('user', JSON.stringify({ mobile: mobileNumber }));
      
      // Redirect to main dashboard
      navigate('/');
    } else {
      // Show error message
      setError('Invalid credentials. Please try again.');
      setShowError(true);
    }
  };

  const handleCloseError = () => {
    setShowError(false);
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#f5f5f5',
      }}
    >
      <CssBaseline />
      <Container maxWidth="xs" sx={{ mb: 4 }}>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            mb: 4
          }}
        >
          <Box 
            component="img"
            src="/images/MQ-Logo-Main.svg"
            alt="MarvelQuant"
            sx={{ height: 80, mb: 2 }}
          />
          {/* Removed MarvelQuant text as requested */}
        </Box>
        
        <Paper
          elevation={3}
          sx={{
            p: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            borderRadius: 2
          }}
        >
          <Typography variant="h5" component="h1" gutterBottom align="center" sx={{ mb: 3 }}>
            SIGN IN
          </Typography>
          
          <Box component="form" onSubmit={handleLogin} sx={{ width: '100%' }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="mobile"
              label="Mobile Number"
              name="mobile"
              autoComplete="tel"
              autoFocus
              value={mobileNumber}
              onChange={(e) => setMobileNumber(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Box 
                      component="span" 
                      sx={{ 
                        display: 'flex', 
                        alignItems: 'center',
                        mr: 0.5
                      }}
                    >
                      <Box 
                        component="img" 
                        src="https://flagcdn.com/w20/in.png" 
                        alt="India" 
                        sx={{ 
                          width: 20, 
                          mr: 0.5 
                        }} 
                      />
                      +91
                    </Box>
                  </InputAdornment>
                ),
              }}
              sx={{ mb: 2 }}
            />
            
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              sx={{ mb: 3 }}
            />
            
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ 
                mt: 1, 
                mb: 2, 
                py: 1.5,
                backgroundColor: '#304FFE',
                '&:hover': {
                  backgroundColor: '#1E40FE',
                }
              }}
            >
              Login
            </Button>
          </Box>
        </Paper>
        
        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Link href="#" variant="body2" sx={{ display: 'block', mb: 1, color: '#304FFE' }}>
            Forgot Password
          </Link>
          <Link href="#" variant="body2" sx={{ color: '#304FFE' }}>
            Sign Up
          </Link>
        </Box>
      </Container>
      
      <Snackbar open={showError} autoHideDuration={6000} onClose={handleCloseError}>
        <Alert onClose={handleCloseError} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default LoginPage;
