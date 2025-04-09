package auth

import (
        "context"
        "errors"
        "fmt"
        "time"

        "github.com/golang-jwt/jwt/v4"
        "trading_platform/backend/internal/config"
        "trading_platform/backend/internal/models"
)

// Claims represents the JWT claims
type Claims struct {
        UserID      string `json:"userId"`
        Username    string `json:"username"`
        Role        string `json:"role"`
        UserType    string `json:"userType"`
        Environment string `json:"environment"`
        jwt.RegisteredClaims
}

// RefreshClaims represents the refresh token claims
type RefreshClaims struct {
        UserID string `json:"userId"`
        jwt.RegisteredClaims
}

// contextKey is a custom type for context keys
type contextKey string

// UserIDKey is the context key for user ID
const UserIDKey contextKey = "userId"

// RoleKey is the context key for user role
const RoleKey contextKey = "role"

// UserTypeKey is the context key for user type
const UserTypeKey contextKey = "userType"

// EnvironmentKey is the context key for environment
const EnvironmentKey contextKey = "environment"

// GenerateToken generates a JWT token
func GenerateToken(userID, username, role, userType string, environment string) (string, error) {
        // Load config
        cfg := config.DefaultConfig()

        // Create claims
        claims := &Claims{
                UserID:      userID,
                Username:    username,
                Role:        role,
                UserType:    userType,
                Environment: environment,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpirationTime)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                        NotBefore: jwt.NewNumericDate(time.Now()),
                        Issuer:    "trading-platform",
                        Subject:   userID,
                },
        }

        // Create token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        // Sign token
        tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
        if err != nil {
                return "", err
        }

        return tokenString, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
        // Load config
        cfg := config.DefaultConfig()

        // Parse token
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                // Validate signing method
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }

                return []byte(cfg.JWT.Secret), nil
        })

        if err != nil {
                return nil, err
        }

        // Extract claims
        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
                return claims, nil
        }

        return nil, errors.New("invalid token")
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userID string) (string, error) {
        // Load config
        cfg := config.DefaultConfig()

        // Create claims
        claims := &RefreshClaims{
                UserID: userID,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshExpiryTime)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                        NotBefore: jwt.NewNumericDate(time.Now()),
                        Issuer:    "trading-platform",
                        Subject:   userID,
                },
        }

        // Create token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        // Sign token
        tokenString, err := token.SignedString([]byte(cfg.JWT.RefreshSecret))
        if err != nil {
                return "", err
        }

        return tokenString, nil
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
        // Load config
        cfg := config.DefaultConfig()

        // Parse token
        token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
                // Validate signing method
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }

                return []byte(cfg.JWT.RefreshSecret), nil
        })

        if err != nil {
                return nil, err
        }

        // Extract claims
        if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
                return claims, nil
        }

        return nil, errors.New("invalid refresh token")
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) string {
        userID, ok := ctx.Value(UserIDKey).(string)
        if !ok {
                return ""
        }
        return userID
}

// GetRoleFromContext gets the user role from the context
func GetRoleFromContext(ctx context.Context) string {
        role, ok := ctx.Value(RoleKey).(string)
        if !ok {
                return ""
        }
        return role
}

// GetUserTypeFromContext gets the user type from the context
func GetUserTypeFromContext(ctx context.Context) string {
        userType, ok := ctx.Value(UserTypeKey).(string)
        if !ok {
                return string(models.UserTypeStandard)
        }
        return userType
}

// GetEnvironmentFromContext gets the environment from the context
func GetEnvironmentFromContext(ctx context.Context) string {
        environment, ok := ctx.Value(EnvironmentKey).(string)
        if !ok {
                return string(models.EnvironmentLive)
        }
        return environment
}

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
        return context.WithValue(ctx, UserIDKey, userID)
}

// SetRoleInContext sets the user role in the context
func SetRoleInContext(ctx context.Context, role string) context.Context {
        return context.WithValue(ctx, RoleKey, role)
}

// SetUserTypeInContext sets the user type in the context
func SetUserTypeInContext(ctx context.Context, userType string) context.Context {
        return context.WithValue(ctx, UserTypeKey, userType)
}

// SetEnvironmentInContext sets the environment in the context
func SetEnvironmentInContext(ctx context.Context, environment string) context.Context {
        return context.WithValue(ctx, EnvironmentKey, environment)
}
