package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/trading-platform/backend/internal/middleware"
	"github.com/trading-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, request auth.RegisterRequest) (*auth.AuthResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*auth.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, request auth.LoginRequest) (*auth.AuthResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*auth.AuthResponse), args.Error(1)
}

func (m *MockAuthService) VerifyToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	// Create a mock auth service
	mockAuthService := new(MockAuthService)
	
	// Create an auth middleware with the mock service
	authMiddleware := middleware.NewAuthMiddleware(mockAuthService)
	
	// Set up Gin for testing
	gin.SetMode(gin.TestMode)
	
	// Test requiring authentication with a valid token
	t.Run("RequireAuth_ValidToken", func(t *testing.T) {
		// Create a new Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Create a test request with an Authorization header
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		c.Request = req
		
		// Set up the mock to return valid claims
		validClaims := &auth.Claims{
			UserID:   "user123",
			Username: "testuser",
		}
		mockAuthService.On("VerifyToken", "valid-token").Return(validClaims, nil)
		
		// Create a test handler that will be called after the middleware
		var handlerCalled bool
		var userID string
		var username string
		testHandler := func(c *gin.Context) {
			handlerCalled = true
			userID = middleware.GetUserID(c)
			username = middleware.GetUsername(c)
			c.Status(http.StatusOK)
		}
		
		// Apply the middleware
		handler := authMiddleware.RequireAuth()
		handler(c)
		
		// If the middleware passed, call the test handler
		if !c.IsAborted() {
			testHandler(c)
		}
		
		// Check that the handler was called
		assert.True(t, handlerCalled)
		
		// Check that the user ID and username were set correctly
		assert.Equal(t, "user123", userID)
		assert.Equal(t, "testuser", username)
		
		// Check that the response status is OK
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify that the mock method was called
		mockAuthService.AssertExpectations(t)
	})
	
	// Test requiring authentication with an invalid token
	t.Run("RequireAuth_InvalidToken", func(t *testing.T) {
		// Create a new Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Create a test request with an Authorization header
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		c.Request = req
		
		// Set up the mock to return an error
		mockAuthService.On("VerifyToken", "invalid-token").Return(nil, assert.AnError)
		
		// Create a test handler that will be called after the middleware
		var handlerCalled bool
		testHandler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}
		
		// Apply the middleware
		handler := authMiddleware.RequireAuth()
		handler(c)
		
		// If the middleware passed, call the test handler
		if !c.IsAborted() {
			testHandler(c)
		}
		
		// Check that the handler was not called
		assert.False(t, handlerCalled)
		
		// Check that the response status is Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		// Verify that the mock method was called
		mockAuthService.AssertExpectations(t)
	})
	
	// Test requiring authentication with no token
	t.Run("RequireAuth_NoToken", func(t *testing.T) {
		// Create a new Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Create a test request with no Authorization header
		req := httptest.NewRequest("GET", "/", nil)
		c.Request = req
		
		// Create a test handler that will be called after the middleware
		var handlerCalled bool
		testHandler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}
		
		// Apply the middleware
		handler := authMiddleware.RequireAuth()
		handler(c)
		
		// If the middleware passed, call the test handler
		if !c.IsAborted() {
			testHandler(c)
		}
		
		// Check that the handler was not called
		assert.False(t, handlerCalled)
		
		// Check that the response status is Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	
	// Test requiring authentication with an invalid Authorization header format
	t.Run("RequireAuth_InvalidHeaderFormat", func(t *testing.T) {
		// Create a new Gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Create a test request with an invalid Authorization header
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		c.Request = req
		
		// Create a test handler that will be called after the middleware
		var handlerCalled bool
		testHandler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}
		
		// Apply the middleware
		handler := authMiddleware.RequireAuth()
		handler(c)
		
		// If the middleware passed, call the test handler
		if !c.IsAborted() {
			testHandler(c)
		}
		
		// Check that the handler was not called
		assert.False(t, handlerCalled)
		
		// Check that the response status is Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
