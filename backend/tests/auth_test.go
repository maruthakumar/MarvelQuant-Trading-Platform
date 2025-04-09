package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/auth"
)

func TestAuthService(t *testing.T) {
	// Create a new auth service with a test secret
	authService := auth.NewAuthService("test-secret-key")

	// Test user registration
	t.Run("Register", func(t *testing.T) {
		ctx := context.Background()
		
		// Create a register request
		request := auth.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		
		// Register the user
		response, err := authService.Register(ctx, request)
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the response is not nil
		require.NotNil(t, response)
		
		// Check that the token is not empty
		assert.NotEmpty(t, response.Token)
		
		// Check that the expiration time is in the future
		assert.Greater(t, response.ExpiresAt, time.Now().Unix())
		
		// Check that the user details are correct
		assert.Equal(t, request.Username, response.User.Username)
		assert.Equal(t, request.Email, response.User.Email)
		assert.NotEmpty(t, response.User.ID)
	})

	// Test user login
	t.Run("Login", func(t *testing.T) {
		ctx := context.Background()
		
		// Create a login request
		request := auth.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		
		// Login the user
		response, err := authService.Login(ctx, request)
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the response is not nil
		require.NotNil(t, response)
		
		// Check that the token is not empty
		assert.NotEmpty(t, response.Token)
		
		// Check that the expiration time is in the future
		assert.Greater(t, response.ExpiresAt, time.Now().Unix())
		
		// Check that the user details are correct
		assert.Equal(t, request.Username, response.User.Username)
		assert.NotEmpty(t, response.User.ID)
	})

	// Test token verification
	t.Run("VerifyToken", func(t *testing.T) {
		ctx := context.Background()
		
		// Create a login request to get a token
		loginRequest := auth.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		
		// Login the user
		loginResponse, err := authService.Login(ctx, loginRequest)
		require.NoError(t, err)
		require.NotNil(t, loginResponse)
		
		// Verify the token
		claims, err := authService.VerifyToken(loginResponse.Token)
		
		// Check that there was no error
		require.NoError(t, err)
		
		// Check that the claims are not nil
		require.NotNil(t, claims)
		
		// Check that the user details in the claims are correct
		assert.Equal(t, loginResponse.User.ID, claims.UserID)
		assert.Equal(t, loginResponse.User.Username, claims.Username)
	})

	// Test invalid token verification
	t.Run("VerifyInvalidToken", func(t *testing.T) {
		// Verify an invalid token
		claims, err := authService.VerifyToken("invalid-token")
		
		// Check that there was an error
		require.Error(t, err)
		
		// Check that the claims are nil
		assert.Nil(t, claims)
	})
}

func TestMain(m *testing.M) {
	// Setup code before running tests
	
	// Run tests
	exitCode := m.Run()
	
	// Cleanup code after running tests
	
	// Exit with the same code as the tests
	os.Exit(exitCode)
}
