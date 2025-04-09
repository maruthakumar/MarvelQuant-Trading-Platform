package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"trading_platform/backend/internal/models"
)

func TestEnvironmentMiddleware(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		environment    string
		targetEnv      string
		expectedStatus int
	}{
		{
			name:           "Matching environment",
			environment:    string(models.EnvironmentLive),
			targetEnv:      string(models.EnvironmentLive),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-matching environment",
			environment:    string(models.EnvironmentLive),
			targetEnv:      string(models.EnvironmentSIM),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "SIM environment match",
			environment:    string(models.EnvironmentSIM),
			targetEnv:      string(models.EnvironmentSIM),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that always returns 200 OK
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := EnvironmentMiddleware(tc.targetEnv)(testHandler)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			
			// Set environment in context
			ctx := context.WithValue(req.Context(), EnvironmentKey, tc.environment)
			req = req.WithContext(ctx)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			middleware.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}

func TestSimUserMiddleware(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		userType       string
		environment    string
		expectedStatus int
	}{
		{
			name:           "SIM user in SIM environment",
			userType:       string(models.UserTypeSIM),
			environment:    string(models.EnvironmentSIM),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Standard user in LIVE environment",
			userType:       string(models.UserTypeStandard),
			environment:    string(models.EnvironmentLive),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Standard user in SIM environment",
			userType:       string(models.UserTypeStandard),
			environment:    string(models.EnvironmentSIM),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Admin user in SIM environment",
			userType:       string(models.UserTypeAdmin),
			environment:    string(models.EnvironmentSIM),
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that always returns 200 OK
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := SimUserMiddleware(testHandler)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			
			// Set user type and environment in context
			ctx := context.WithValue(req.Context(), UserTypeKey, tc.userType)
			ctx = context.WithValue(ctx, EnvironmentKey, tc.environment)
			req = req.WithContext(ctx)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			middleware.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}

func TestUserTypeMiddleware(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		userType       string
		allowedTypes   []string
		expectedStatus int
	}{
		{
			name:           "SIM user allowed",
			userType:       string(models.UserTypeSIM),
			allowedTypes:   []string{string(models.UserTypeSIM)},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Admin user allowed",
			userType:       string(models.UserTypeAdmin),
			allowedTypes:   []string{string(models.UserTypeAdmin), string(models.UserTypeSIM)},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Standard user not allowed",
			userType:       string(models.UserTypeStandard),
			allowedTypes:   []string{string(models.UserTypeSIM), string(models.UserTypeAdmin)},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Multiple allowed types",
			userType:       string(models.UserTypeStandard),
			allowedTypes:   []string{string(models.UserTypeStandard), string(models.UserTypeSIM)},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that always returns 200 OK
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := UserTypeMiddleware(tc.allowedTypes...)(testHandler)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)
			
			// Set user type in context
			ctx := context.WithValue(req.Context(), UserTypeKey, tc.userType)
			req = req.WithContext(ctx)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			middleware.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
