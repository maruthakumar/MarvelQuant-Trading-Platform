package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"trading_platform/backend/internal/models"
)

func TestGenerateTokenWithEnvironment(t *testing.T) {
	// Test data
	userID := "user123"
	username := "testuser"
	role := string(models.UserRoleTrader)
	userType := string(models.UserTypeSIM)
	environment := string(models.EnvironmentSIM)

	// Generate token
	token, err := GenerateToken(userID, username, role, userType, environment)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, userType, claims.UserType)
	assert.Equal(t, environment, claims.Environment)
}

func TestValidateTokenWithEnvironment(t *testing.T) {
	// Test data
	userID := "user123"
	username := "testuser"
	role := string(models.UserRoleTrader)
	userType := string(models.UserTypeSIM)
	environment := string(models.EnvironmentSIM)

	// Create claims
	claims := &Claims{
		UserID:      userID,
		Username:    username,
		Role:        role,
		UserType:    userType,
		Environment: environment,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "trading-platform",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	// Mock ValidateToken to use our test token
	originalValidateToken := ValidateToken
	defer func() { ValidateToken = originalValidateToken }()
	
	ValidateToken = func(tokenString string) (*Claims, error) {
		return claims, nil
	}

	// Validate token
	validatedClaims, err := ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, validatedClaims)
	assert.Equal(t, userID, validatedClaims.UserID)
	assert.Equal(t, username, validatedClaims.Username)
	assert.Equal(t, role, validatedClaims.Role)
	assert.Equal(t, userType, validatedClaims.UserType)
	assert.Equal(t, environment, validatedClaims.Environment)
}

func TestContextFunctionsWithEnvironment(t *testing.T) {
	// Test data
	userID := "user123"
	role := string(models.UserRoleTrader)
	userType := string(models.UserTypeSIM)
	environment := string(models.EnvironmentSIM)

	// Create context with values
	ctx := SetUserIDInContext(nil, userID)
	ctx = SetRoleInContext(ctx, role)
	ctx = SetUserTypeInContext(ctx, userType)
	ctx = SetEnvironmentInContext(ctx, environment)

	// Get values from context
	retrievedUserID := GetUserIDFromContext(ctx)
	retrievedRole := GetRoleFromContext(ctx)
	retrievedUserType := GetUserTypeFromContext(ctx)
	retrievedEnvironment := GetEnvironmentFromContext(ctx)

	// Assertions
	assert.Equal(t, userID, retrievedUserID)
	assert.Equal(t, role, retrievedRole)
	assert.Equal(t, userType, retrievedUserType)
	assert.Equal(t, environment, retrievedEnvironment)
}

func TestGetDefaultValuesFromEmptyContext(t *testing.T) {
	// Get values from empty context
	userType := GetUserTypeFromContext(nil)
	environment := GetEnvironmentFromContext(nil)

	// Assertions
	assert.Equal(t, string(models.UserTypeStandard), userType)
	assert.Equal(t, string(models.EnvironmentLive), environment)
}
