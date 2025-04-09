package gateway

import (
	"context"
	"fmt"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"trading_platform/backend/internal/models"
)

// TestSecurityIsolation tests the security isolation between SIM and LIVE environments
func TestSecurityIsolation(t *testing.T) {
	// Create mock services
	mockExecutionPlatform := new(MockExecutionPlatform)
	
	// Create API Gateway with mock execution platform
	gateway := NewAPIGateway(mockExecutionPlatform)
	
	// Add permissions for test users
	gateway.accessControlList["sim_user"] = []string{
		"simulation:account:read",
		"simulation:order:create",
		"simulation:market:read",
	}
	
	gateway.accessControlList["live_user"] = []string{
		"live:account:read",
		"live:order:create",
		"live:market:read",
	}
	
	gateway.accessControlList["admin_user"] = []string{
		"*", // Admin has all permissions
	}
	
	t.Run("SIM User Cannot Access LIVE Resources", func(t *testing.T) {
		// Create context with SIM user
		simCtx := context.WithValue(context.Background(), "userID", "sim_user")
		simCtx = context.WithValue(simCtx, "userType", "SIM")
		
		// Try to access a LIVE resource with SIM user
		// We'll use a custom method for this test
		err := gateway.checkPermission(simCtx, "live:account:read")
		
		// Assert that access is denied
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SIM users can only access simulation resources")
	})
	
	t.Run("LIVE User Can Access LIVE Resources", func(t *testing.T) {
		// Create context with LIVE user
		liveCtx := context.WithValue(context.Background(), "userID", "live_user")
		liveCtx = context.WithValue(liveCtx, "userType", "STANDARD")
		
		// Try to access a LIVE resource with LIVE user
		err := gateway.checkPermission(liveCtx, "live:account:read")
		
		// Assert that access is granted
		assert.NoError(t, err)
	})
	
	t.Run("LIVE User Cannot Access SIM Resources", func(t *testing.T) {
		// Create context with LIVE user
		liveCtx := context.WithValue(context.Background(), "userID", "live_user")
		liveCtx = context.WithValue(liveCtx, "userType", "STANDARD")
		
		// Try to access a SIM resource with LIVE user
		err := gateway.checkPermission(liveCtx, "simulation:account:read")
		
		// Assert that access is denied
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user does not have required permission")
	})
	
	t.Run("Admin User Can Access All Resources", func(t *testing.T) {
		// Create context with Admin user
		adminCtx := context.WithValue(context.Background(), "userID", "admin_user")
		adminCtx = context.WithValue(adminCtx, "userType", "ADMIN")
		
		// Try to access both LIVE and SIM resources with Admin user
		err1 := gateway.checkPermission(adminCtx, "live:account:read")
		err2 := gateway.checkPermission(adminCtx, "simulation:account:read")
		
		// Assert that access is granted for both
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})
}

// TestErrorHandling tests the error handling functionality
func TestErrorHandling(t *testing.T) {
	// Create API Gateway
	gateway := NewAPIGateway(nil)
	
	// Create context
	ctx := context.Background()
	
	t.Run("Authentication Error Handling", func(t *testing.T) {
		// Create an authentication error
		originalError := fmt.Errorf("invalid token signature")
		
		// Handle the error
		handledError := gateway.handleError(ctx, "authentication", originalError)
		
		// Assert that the error is properly handled
		assert.Error(t, handledError)
		assert.Contains(t, handledError.Error(), "authentication failed")
		assert.NotContains(t, handledError.Error(), "invalid token signature") // Original error details should be hidden
	})
	
	t.Run("Authorization Error Handling", func(t *testing.T) {
		// Create an authorization error
		originalError := fmt.Errorf("user does not have permission to access resource")
		
		// Handle the error
		handledError := gateway.handleError(ctx, "authorization", originalError)
		
		// Assert that the error is properly handled
		assert.Error(t, handledError)
		assert.Contains(t, handledError.Error(), "authorization failed")
		assert.NotContains(t, handledError.Error(), "user does not have permission") // Original error details should be hidden
	})
	
	t.Run("Validation Error Handling", func(t *testing.T) {
		// Create a validation error
		originalError := fmt.Errorf("invalid order quantity: must be greater than zero")
		
		// Handle the error
		handledError := gateway.handleError(ctx, "validation", originalError)
		
		// Assert that the error is properly handled
		assert.Error(t, handledError)
		assert.Contains(t, handledError.Error(), "invalid order quantity") // Validation errors should preserve details
	})
	
	t.Run("System Error Handling", func(t *testing.T) {
		// Create a system error
		originalError := fmt.Errorf("database connection failed: connection refused")
		
		// Handle the error
		handledError := gateway.handleError(ctx, "system", originalError)
		
		// Assert that the error is properly handled
		assert.Error(t, handledError)
		assert.Contains(t, handledError.Error(), "internal system error")
		assert.NotContains(t, handledError.Error(), "database connection failed") // System details should be hidden
	})
	
	t.Run("Unknown Error Category", func(t *testing.T) {
		// Create an error with unknown category
		originalError := fmt.Errorf("some unknown error")
		
		// Handle the error
		handledError := gateway.handleError(ctx, "unknown_category", originalError)
		
		// Assert that the error is handled as a system error (fallback)
		assert.Error(t, handledError)
		assert.Contains(t, handledError.Error(), "internal system error")
	})
}

// TestRateLimiting tests the rate limiting functionality
func TestRateLimiting(t *testing.T) {
	// Create API Gateway
	gateway := NewAPIGateway(nil)
	
	// Create context with user ID
	ctx := context.WithValue(context.Background(), "userID", "test_user")
	
	t.Run("Within Rate Limit", func(t *testing.T) {
		// Set up a rate limit
		gateway.rateLimits["test_category"] = RateLimit{
			MaxRequests:     5,
			TimeWindow:      time.Minute,
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// Make multiple requests within the limit
		for i := 0; i < 5; i++ {
			err := gateway.checkRateLimit(ctx, "test_category")
			assert.NoError(t, err)
		}
	})
	
	t.Run("Exceeding Rate Limit", func(t *testing.T) {
		// Set up a rate limit
		gateway.rateLimits["test_category"] = RateLimit{
			MaxRequests:     3,
			TimeWindow:      time.Minute,
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// Make requests up to the limit
		for i := 0; i < 3; i++ {
			err := gateway.checkRateLimit(ctx, "test_category")
			assert.NoError(t, err)
		}
		
		// Next request should exceed the limit
		err := gateway.checkRateLimit(ctx, "test_category")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
	})
	
	t.Run("Expired Requests Don't Count", func(t *testing.T) {
		// Set up a rate limit with a very short window
		gateway.rateLimits["test_category"] = RateLimit{
			MaxRequests:     2,
			TimeWindow:      10 * time.Millisecond, // Very short for testing
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// Make requests up to the limit
		for i := 0; i < 2; i++ {
			err := gateway.checkRateLimit(ctx, "test_category")
			assert.NoError(t, err)
		}
		
		// Wait for the window to expire
		time.Sleep(20 * time.Millisecond)
		
		// Next request should succeed because previous ones expired
		err := gateway.checkRateLimit(ctx, "test_category")
		assert.NoError(t, err)
	})
	
	t.Run("Different Users Have Separate Limits", func(t *testing.T) {
		// Set up a rate limit
		gateway.rateLimits["test_category"] = RateLimit{
			MaxRequests:     2,
			TimeWindow:      time.Minute,
			CurrentRequests: make(map[string][]time.Time),
		}
		
		// Create contexts for two different users
		ctx1 := context.WithValue(context.Background(), "userID", "user1")
		ctx2 := context.WithValue(context.Background(), "userID", "user2")
		
		// User 1 makes requests up to the limit
		for i := 0; i < 2; i++ {
			err := gateway.checkRateLimit(ctx1, "test_category")
			assert.NoError(t, err)
		}
		
		// User 1's next request should exceed the limit
		err1 := gateway.checkRateLimit(ctx1, "test_category")
		assert.Error(t, err1)
		
		// User 2 should still be able to make requests
		for i := 0; i < 2; i++ {
			err := gateway.checkRateLimit(ctx2, "test_category")
			assert.NoError(t, err)
		}
	})
}

// TestPermissionChecking tests the permission checking functionality
func TestPermissionChecking(t *testing.T) {
	// Create API Gateway
	gateway := NewAPIGateway(nil)
	
	// Add permissions for test users
	gateway.accessControlList["restricted_user"] = []string{
		"resource1:action1",
		"resource1:action2",
		"resource2:action1",
	}
	
	gateway.accessControlList["wildcard_user"] = []string{
		"resource1:*", // All actions on resource1
	}
	
	gateway.accessControlList["super_user"] = []string{
		"*", // All permissions
	}
	
	t.Run("User Has Specific Permission", func(t *testing.T) {
		// Create context with user ID
		ctx := context.WithValue(context.Background(), "userID", "restricted_user")
		ctx = context.WithValue(ctx, "userType", "STANDARD")
		
		// Check permissions
		err1 := gateway.checkPermission(ctx, "resource1:action1")
		err2 := gateway.checkPermission(ctx, "resource1:action2")
		err3 := gateway.checkPermission(ctx, "resource2:action1")
		
		// Assert that permissions are granted
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
	})
	
	t.Run("User Doesn't Have Permission", func(t *testing.T) {
		// Create context with user ID
		ctx := context.WithValue(context.Background(), "userID", "restricted_user")
		ctx = context.WithValue(ctx, "userType", "STANDARD")
		
		// Check permission that user doesn't have
		err := gateway.checkPermission(ctx, "resource2:action2")
		
		// Assert that permission is denied
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user does not have required permission")
	})
	
	t.Run("User Has Wildcard Permission", func(t *testing.T) {
		// Create context with user ID
		ctx := context.WithValue(context.Background(), "userID", "wildcard_user")
		ctx = context.WithValue(ctx, "userType", "STANDARD")
		
		// Check permissions covered by wildcard
		err1 := gateway.checkPermission(ctx, "resource1:action1")
		err2 := gateway.checkPermission(ctx, "resource1:action2")
		err3 := gateway.checkPermission(ctx, "resource1:action3") // Not explicitly granted but covered by wildcard
		
		// Assert that permissions are granted
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		
		// Check permission not covered by wildcard
		err4 := gateway.checkPermission(ctx, "resource2:action1")
		
		// Assert that permission is denied
		assert.Error(t, err4)
	})
	
	t.Run("Super User Has All Permissions", func(t *testing.T) {
		// Create context with user ID
		ctx := context.WithValue(context.Background(), "userID", "super_user")
		ctx = context.WithValue(ctx, "userType", "STANDARD")
		
		// Check various permissions
		err1 := gateway.checkPermission(ctx, "resource1:action1")
		err2 := gateway.checkPermission(ctx, "resource2:action2")
		err3 := gateway.checkPermission(ctx, "resource3:action3") // Not explicitly granted but covered by wildcard
		
		// Assert that all permissions are granted
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
	})
	
	t.Run("Admin User Type Has All Permissions", func(t *testing.T) {
		// Create context with user ID and admin type
		ctx := context.WithValue(context.Background(), "userID", "restricted_user") // Even with restricted permissions
		ctx = context.WithValue(ctx, "userType", "ADMIN")
		
		// Check various permissions
		err1 := gateway.checkPermission(ctx, "resource1:action1")
		err2 := gateway.checkPermission(ctx, "resource2:action2") // Not explicitly granted
		err3 := gateway.checkPermission(ctx, "resource3:action3") // Not explicitly granted
		
		// Assert that all permissions are granted due to admin type
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
	})
	
	t.Run("Missing User ID in Context", func(t *testing.T) {
		// Create context without user ID
		ctx := context.Background()
		
		// Check permission
		err := gateway.checkPermission(ctx, "resource1:action1")
		
		// Assert that error is returned
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID not found in context")
	})
	
	t.Run("Missing User Type in Context", func(t *testing.T) {
		// Create context with user ID but without user type
		ctx := context.WithValue(context.Background(), "userID", "restricted_user")
		
		// Check permission
		err := gateway.checkPermission(ctx, "resource1:action1")
		
		// Assert that error is returned
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user type not found in context")
	})
}

// TestSimulationPermissionCheck tests the isSimulationPermission function
func TestSimulationPermissionCheck(t *testing.T) {
	testCases := []struct {
		permission string
		expected   bool
	}{
		{"simulation:account:read", true},
		{"simulation:order:create", true},
		{"backtest:session:run", true},
		{"papertrading:position:close", true},
		{"live:account:read", false},
		{"system:status:read", false},
		{"random:permission", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.permission, func(t *testing.T) {
			result := isSimulationPermission(tc.permission)
			assert.Equal(t, tc.expected, result)
		})
	}
}
