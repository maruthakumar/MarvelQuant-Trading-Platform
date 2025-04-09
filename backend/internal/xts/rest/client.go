package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/trade-execution-platform/backend/internal/xts/config"
	"github.com/trade-execution-platform/backend/internal/xts/errors"
	"github.com/trade-execution-platform/backend/internal/xts/models"
)

// Client represents an XTS REST API client
type Client struct {
	BaseURL    string
	APIKey     string
	SecretKey  string
	Token      string
	UserID     string
	HTTPClient *http.Client
	Routes     *config.Routes
}

// NewClient creates a new XTS REST API client
func NewClient(cfg *config.XTSConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:    cfg.BaseURL,
		APIKey:     cfg.APIKey,
		SecretKey:  cfg.SecretKey,
		HTTPClient: &http.Client{Timeout: cfg.Timeout},
		Routes:     config.DefaultRoutes(),
	}

	return client, nil
}

// SetToken sets the authentication token
func (c *Client) SetToken(token, userID string) {
	c.Token = token
	c.UserID = userID
}

// Login authenticates with the XTS API
func (c *Client) Login() (*models.Session, error) {
	params := map[string]interface{}{
		"appKey":    c.APIKey,
		"secretKey": c.SecretKey,
		"source":    "WEB",
	}

	var response models.LoginResponse
	err := c.doRequest(http.MethodPost, c.Routes.UserLogin, params, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Login failed",
			Description: response.Description,
			HTTPStatus:  http.StatusUnauthorized,
		}
	}

	session := &models.Session{
		Token:           response.Result.Token,
		UserID:          response.Result.UserID,
		IsInvestorClient: response.Result.IsInvestorClient,
		ExpiresAt:       time.Now().Add(24 * time.Hour), // Token typically valid for 24 hours
	}

	c.SetToken(session.Token, session.UserID)
	return session, nil
}

// Logout terminates the current session
func (c *Client) Logout() error {
	if c.Token == "" {
		return errors.ErrSessionInvalid
	}

	var response map[string]interface{}
	err := c.doRequest(http.MethodDelete, c.Routes.UserLogout, nil, &response)
	if err != nil {
		return err
	}

	c.Token = ""
	c.UserID = ""
	return nil
}

// GetOrderBook retrieves the current order book
func (c *Client) GetOrderBook(clientID string) (*models.OrderBook, error) {
	if c.Token == "" {
		return nil, errors.ErrSessionInvalid
	}

	params := make(map[string]string)
	if clientID != "" {
		params["clientID"] = clientID
	}

	var response struct {
		Type        string `json:"type"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Result      models.OrderBook `json:"result"`
	}

	err := c.doRequest(http.MethodGet, c.Routes.OrderStatus, params, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Failed to get order book",
			Description: response.Description,
			HTTPStatus:  http.StatusBadRequest,
		}
	}

	return &response.Result, nil
}

// PlaceOrder places a new order
func (c *Client) PlaceOrder(order *models.Order) (*models.OrderResponse, error) {
	if c.Token == "" {
		return nil, errors.ErrSessionInvalid
	}

	var response struct {
		Type        string `json:"type"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Result      models.OrderResponse `json:"result"`
	}

	err := c.doRequest(http.MethodPost, c.Routes.OrderPlace, order, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Order placement failed",
			Description: response.Description,
			HTTPStatus:  http.StatusBadRequest,
		}
	}

	return &response.Result, nil
}

// ModifyOrder modifies an existing order
func (c *Client) ModifyOrder(modifyOrder *models.ModifyOrder) (*models.OrderResponse, error) {
	if c.Token == "" {
		return nil, errors.ErrSessionInvalid
	}

	var response struct {
		Type        string `json:"type"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Result      models.OrderResponse `json:"result"`
	}

	err := c.doRequest(http.MethodPut, c.Routes.OrderModify, modifyOrder, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Order modification failed",
			Description: response.Description,
			HTTPStatus:  http.StatusBadRequest,
		}
	}

	return &response.Result, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(appOrderID int, clientID string) (*models.OrderResponse, error) {
	if c.Token == "" {
		return nil, errors.ErrSessionInvalid
	}

	params := map[string]interface{}{
		"appOrderID": appOrderID,
	}

	if clientID != "" {
		params["clientID"] = clientID
	}

	var response struct {
		Type        string `json:"type"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Result      models.OrderResponse `json:"result"`
	}

	err := c.doRequest(http.MethodDelete, c.Routes.OrderCancel, params, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Order cancellation failed",
			Description: response.Description,
			HTTPStatus:  http.StatusBadRequest,
		}
	}

	return &response.Result, nil
}

// GetPositions retrieves current positions
func (c *Client) GetPositions(clientID string) ([]models.Position, error) {
	if c.Token == "" {
		return nil, errors.ErrSessionInvalid
	}

	params := make(map[string]string)
	if clientID != "" {
		params["clientID"] = clientID
	}

	var response struct {
		Type        string `json:"type"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Result      []models.Position `json:"result"`
	}

	err := c.doRequest(http.MethodGet, c.Routes.PortfolioPositions, params, &response)
	if err != nil {
		return nil, err
	}

	if response.Type != "success" {
		return nil, &errors.XTSError{
			Code:        response.Code,
			Message:     "Failed to get positions",
			Description: response.Description,
			HTTPStatus:  http.StatusBadRequest,
		}
	}

	return response.Result, nil
}

// doRequest performs an HTTP request to the XTS API
func (c *Client) doRequest(method, endpoint string, params interface{}, result interface{}) error {
	var reqBody io.Reader
	var queryParams string

	// Prepare request body or query parameters
	if params != nil {
		switch method {
		case http.MethodGet, http.MethodDelete:
			// For GET and DELETE, use query parameters
			values := url.Values{}
			
			// Handle different parameter types
			switch p := params.(type) {
			case map[string]string:
				for k, v := range p {
					values.Add(k, v)
				}
			case map[string]interface{}:
				for k, v := range p {
					values.Add(k, fmt.Sprintf("%v", v))
				}
			default:
				return fmt.Errorf("unsupported params type for GET/DELETE: %T", params)
			}
			
			if len(values) > 0 {
				queryParams = "?" + values.Encode()
			}
		case http.MethodPost, http.MethodPut:
			// For POST and PUT, use JSON body
			jsonData, err := json.Marshal(params)
			if err != nil {
				return err
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	}

	// Create request
	fullURL := c.BaseURL + endpoint + queryParams
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}

	// Execute request with context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return errors.ErrRequestTimeout
		}
		return errors.Wrap(err, "network_error", "Failed to connect to XTS API", http.StatusServiceUnavailable)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read_error", "Failed to read response body", resp.StatusCode)
	}

	// Handle non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err == nil {
			return errors.FromResponse(errResp)
		}
		return &errors.XTSError{
			Code:        fmt.Sprintf("%d", resp.StatusCode),
			Message:     "HTTP request failed",
			Description: string(body),
			HTTPStatus:  resp.StatusCode,
		}
	}

	// Parse response
	if err := json.Unmarshal(body, result); err != nil {
		return errors.Wrap(err, "parse_error", "Failed to parse response", http.StatusOK)
	}

	return nil
}
