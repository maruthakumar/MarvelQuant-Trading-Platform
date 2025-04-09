package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Standard error types
var (
	// Configuration errors
	ErrEmptyBaseURL   = errors.New("base URL cannot be empty")
	ErrEmptyAPIKey    = errors.New("API key cannot be empty")
	ErrEmptySecretKey = errors.New("secret key cannot be empty")
	ErrInvalidTimeout = errors.New("timeout must be greater than zero")

	// Authentication errors
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidCredentials   = errors.New("invalid API key or secret key")
	ErrSessionInvalid       = errors.New("session is invalid or expired")

	// Network errors
	ErrConnectionFailed = errors.New("connection to XTS API failed")
	ErrRequestTimeout   = errors.New("request to XTS API timed out")
	ErrNetworkFailure   = errors.New("network failure occurred")

	// API errors
	ErrInvalidRequest  = errors.New("invalid request parameters")
	ErrAPIRateLimited  = errors.New("API rate limit exceeded")
	ErrAPIUnavailable  = errors.New("API service is unavailable")
	ErrInvalidResponse = errors.New("invalid response from API")

	// Order errors
	ErrOrderRejected      = errors.New("order was rejected")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderParams = errors.New("invalid order parameters")
	ErrOrderModifyFailed  = errors.New("order modification failed")
	ErrOrderCancelFailed  = errors.New("order cancellation failed")

	// WebSocket errors
	ErrWebSocketConnFailed = errors.New("WebSocket connection failed")
	ErrWebSocketClosed     = errors.New("WebSocket connection closed")
	ErrSubscriptionFailed  = errors.New("subscription to market data failed")
)

// XTSError represents an error from the XTS API
type XTSError struct {
	Code        string
	Message     string
	Description string
	HTTPStatus  int
	Err         error
}

// Error returns the error message
func (e *XTSError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *XTSError) Unwrap() error {
	return e.Err
}

// New creates a new XTS error
func New(code, message, description string, httpStatus int) *XTSError {
	return &XTSError{
		Code:        code,
		Message:     message,
		Description: description,
		HTTPStatus:  httpStatus,
	}
}

// Wrap wraps an existing error in an XTSError
func Wrap(err error, code, message string, httpStatus int) *XTSError {
	return &XTSError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// FromResponse creates an XTSError from an API response
func FromResponse(resp map[string]interface{}) *XTSError {
	code, _ := resp["code"].(string)
	message, _ := resp["message"].(string)
	description, _ := resp["description"].(string)
	
	return &XTSError{
		Code:        code,
		Message:     message,
		Description: description,
		HTTPStatus:  http.StatusBadRequest, // Default status
	}
}

// IsAuthError checks if the error is an authentication error
func IsAuthError(err error) bool {
	var xtsErr *XTSError
	if errors.As(err, &xtsErr) {
		return xtsErr.HTTPStatus == http.StatusUnauthorized
	}
	return errors.Is(err, ErrAuthenticationFailed) ||
		errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrSessionInvalid)
}

// IsNetworkError checks if the error is a network error
func IsNetworkError(err error) bool {
	var xtsErr *XTSError
	if errors.As(err, &xtsErr) {
		return xtsErr.HTTPStatus >= 500 || xtsErr.HTTPStatus == 0
	}
	return errors.Is(err, ErrConnectionFailed) ||
		errors.Is(err, ErrRequestTimeout) ||
		errors.Is(err, ErrNetworkFailure) ||
		errors.Is(err, ErrAPIUnavailable)
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	var xtsErr *XTSError
	if errors.As(err, &xtsErr) {
		return xtsErr.HTTPStatus == http.StatusTooManyRequests
	}
	return errors.Is(err, ErrAPIRateLimited)
}

// IsOrderError checks if the error is related to order operations
func IsOrderError(err error) bool {
	return errors.Is(err, ErrOrderRejected) ||
		errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrInvalidOrderParams) ||
		errors.Is(err, ErrOrderModifyFailed) ||
		errors.Is(err, ErrOrderCancelFailed)
}

// IsWebSocketError checks if the error is related to WebSocket operations
func IsWebSocketError(err error) bool {
	return errors.Is(err, ErrWebSocketConnFailed) ||
		errors.Is(err, ErrWebSocketClosed) ||
		errors.Is(err, ErrSubscriptionFailed)
}
