package zerodha

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
	kiteconnect "github.com/zerodha/gokiteconnect"
)

// mockKiteConnectServer creates a mock HTTP server for testing Zerodha API calls
func mockKiteConnectServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

// TestLogin tests the Login method
func TestLogin(t *testing.T) {
	// Setup mock server for session generation
	server := mockKiteConnectServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/session/token", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		
		assert.Equal(t, "test_request_token", r.Form.Get("request_token"))
		assert.Equal(t, "test_api_secret", r.Form.Get("api_secret"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"status": "success",
			"data": {
				"user_id": "XX000",
				"user_name": "Test User",
				"user_shortname": "test",
				"email": "test@example.com",
				"user_type": "individual",
				"broker": "ZERODHA",
				"exchanges": ["NSE", "BSE", "NFO", "CDS", "MCX"],
				"products": ["CNC", "NRML", "MIS"],
				"order_types": ["MARKET", "LIMIT", "SL", "SL-M"],
				"access_token": "test_access_token",
				"public_token": "test_public_token",
				"refresh_token": "",
				"login_time": "2020-01-01 12:00:00"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Create adapter with mock server URL
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
		BaseURL:     server.URL,
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	
	// Override the client with one that points to our mock server
	adapter.client = kiteconnect.New("test_api_key")
	adapter.client.SetBaseURI(server.URL)
	
	// Test login with request token
	credentials := &common.Credentials{
		TwoFactorCode: "test_request_token",
	}
	
	session, err := adapter.Login(credentials)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test_access_token", session.Token)
	assert.Equal(t, "XX000", session.UserID)
	
	// Test login without request token
	credentials = &common.Credentials{}
	session, err = adapter.Login(credentials)
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "please visit")
	
	// Test login with nil credentials
	session, err = adapter.Login(nil)
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "please visit")
}

// TestLogout tests the Logout method
func TestLogout(t *testing.T) {
	// Setup mock server for session invalidation
	server := mockKiteConnectServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/session/token", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "test_access_token", r.Header.Get("Authorization"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"status": "success"
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Create adapter with mock server URL
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
		BaseURL:     server.URL,
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	
	// Override the client with one that points to our mock server
	adapter.client = kiteconnect.New("test_api_key")
	adapter.client.SetBaseURI(server.URL)
	
	// Set access token
	adapter.accessToken = "test_access_token"
	adapter.client.SetAccessToken("test_access_token")
	
	// Test logout
	err = adapter.Logout()
	assert.NoError(t, err)
	assert.Empty(t, adapter.accessToken)
	assert.Empty(t, adapter.userID)
	
	// Test logout when not logged in
	err = adapter.Logout()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestPlaceOrder tests the PlaceOrder method
func TestPlaceOrder(t *testing.T) {
	// Setup mock server for placing order
	server := mockKiteConnectServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/regular", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "test_access_token", r.Header.Get("Authorization"))
		
		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		
		assert.Equal(t, "NSE", r.Form.Get("exchange"))
		assert.Equal(t, "RELIANCE", r.Form.Get("tradingsymbol"))
		assert.Equal(t, "BUY", r.Form.Get("transaction_type"))
		assert.Equal(t, "10", r.Form.Get("quantity"))
		assert.Equal(t, "MIS", r.Form.Get("product"))
		assert.Equal(t, "LIMIT", r.Form.Get("order_type"))
		assert.Equal(t, "DAY", r.Form.Get("validity"))
		assert.Equal(t, "2000", r.Form.Get("price"))
		assert.Equal(t, "0", r.Form.Get("trigger_price"))
		assert.Equal(t, "test123", r.Form.Get("tag"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"status": "success",
			"data": {
				"order_id": "order123"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Create adapter with mock server URL
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
		BaseURL:     server.URL,
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	
	// Override the client with one that points to our mock server
	adapter.client = kiteconnect.New("test_api_key")
	adapter.client.SetBaseURI(server.URL)
	
	// Set access token
	adapter.accessToken = "test_access_token"
	adapter.client.SetAccessToken("test_access_token")
	
	// Test place order
	order := &common.Order{
		ExchangeSegment:       "NSECM",
		TradingSymbol:         "RELIANCE",
		OrderSide:             "BUY",
		OrderQuantity:         10,
		ProductType:           "MIS",
		OrderType:             "LIMIT",
		TimeInForce:           "DAY",
		LimitPrice:            2000,
		StopPrice:             0,
		OrderUniqueIdentifier: "test123",
	}
	
	orderResponse, err := adapter.PlaceOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "PLACED", orderResponse.Status)
	
	// Test place order when not logged in
	adapter.accessToken = ""
	orderResponse, err = adapter.PlaceOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test place order with nil order
	adapter.accessToken = "test_access_token"
	orderResponse, err = adapter.PlaceOrder(nil)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order is required")
}

// TestModifyOrder tests the ModifyOrder method
func TestModifyOrder(t *testing.T) {
	// Setup mock server for modifying order
	server := mockKiteConnectServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/regular/order123", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "test_access_token", r.Header.Get("Authorization"))
		
		// Parse form data
		err := r.ParseForm()
		assert.NoError(t, err)
		
		assert.Equal(t, "20", r.Form.Get("quantity"))
		assert.Equal(t, "LIMIT", r.Form.Get("order_type"))
		assert.Equal(t, "2100", r.Form.Get("price"))
		assert.Equal(t, "0", r.Form.Get("trigger_price"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"status": "success",
			"data": {
				"order_id": "order123"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Create adapter with mock server URL
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
		BaseURL:     server.URL,
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	
	// Override the client with one that points to our mock server
	adapter.client = kiteconnect.New("test_api_key")
	adapter.client.SetBaseURI(server.URL)
	
	// Set access token
	adapter.accessToken = "test_access_token"
	adapter.client.SetAccessToken("test_access_token")
	
	// Test modify order
	order := &common.ModifyOrder{
		OrderID:       "order123",
		OrderType:     "LIMIT",
		OrderQuantity: 20,
		LimitPrice:    2100,
		StopPrice:     0,
	}
	
	orderResponse, err := adapter.ModifyOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "MODIFIED", orderResponse.Status)
	
	// Test modify order when not logged in
	adapter.accessToken = ""
	orderResponse, err = adapter.ModifyOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test modify order with nil order
	adapter.accessToken = "test_access_token"
	orderResponse, err = adapter.ModifyOrder(nil)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order is required")
	
	// Test modify order with missing order ID
	order.OrderID = ""
	orderResponse, err = adapter.ModifyOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order ID is required")
}

// TestCancelOrder tests the CancelOrder method
func TestCancelOrder(t *testing.T) {
	// Setup mock server for cancelling order
	server := mockKiteConnectServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/orders/regular/order123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "test_access_token", r.Header.Get("Authorization"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"status": "success",
			"data": {
				"order_id": "order123"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Create adapter with mock server URL
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
		BaseURL:     server.URL,
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	
	// Override the client with one that points to our mock server
	adapter.client = kiteconnect.New("test_api_key")
	adapter.client.SetBaseURI(server.URL)
	
	// Set access token
	adapter.accessToken = "test_access_token"
	adapter.client.SetAccessToken("test_access_token")
	
	// Test cancel order
	orderResponse, err := adapter.CancelOrder("order123", "")
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "CANCELLED", orderResponse.Status)
	
	// Test cancel order when not logged in
	adapter.accessToken = ""
	orderResponse, err = adapter.CancelOrder("order123", "")
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test cancel order with missing order ID
	adapter.accessToken = "test_access_token"
	orderResponse, err = adapter.CancelOrder("", "")
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order ID is required")
}
