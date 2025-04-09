package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret []byte
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents a request to login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents a response from authentication
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, request RegisterRequest) (*AuthResponse, error) {
	// This would typically involve:
	// 1. Validating the request data
	// 2. Checking if the username or email already exists
	// 3. Hashing the password
	// 4. Storing the user in the database
	// 5. Generating a JWT token

	// For this implementation, we'll create a placeholder user and token
	log.Printf("Registering user %s", request.Username)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a placeholder user
	user := User{
		ID:           fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Generate a token
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, request LoginRequest) (*AuthResponse, error) {
	// This would typically involve:
	// 1. Retrieving the user from the database
	// 2. Verifying the password
	// 3. Generating a JWT token

	// For this implementation, we'll create a placeholder user and token
	log.Printf("Logging in user %s", request.Username)

	// Create a placeholder user (in a real implementation, this would be retrieved from the database)
	user := User{
		ID:           "user-123",
		Username:     request.Username,
		Email:        fmt.Sprintf("%s@example.com", request.Username),
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz0123456789", // Placeholder hash
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// In a real implementation, we would verify the password
	// For this placeholder, we'll just generate a token
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// VerifyToken verifies a JWT token
func (s *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user User) (string, int64, error) {
	// Set the expiration time
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	// Create the claims
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "trading-platform",
			Subject:   user.ID,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// Middleware for authentication
type Middleware struct {
	authService *AuthService
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(authService *AuthService) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

// AuthenticateUser is a middleware function to authenticate a user
func (m *Middleware) AuthenticateUser(next func(ctx context.Context, userID string) error) func(ctx context.Context, token string) error {
	return func(ctx context.Context, token string) error {
		// Verify the token
		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Call the next handler with the user ID
		return next(ctx, claims.UserID)
	}
}
