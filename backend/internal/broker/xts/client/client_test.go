package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
)

// setupMockServer creates a mock HTTP server for testing
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *XTSClientImpl) {
	server := httptest.NewServer(handler)
	
	config := &common.XTSClientConfig{
		APIKey:    "test_api_key",
		SecretKey: "test_secret_key",
		Source:    "WEBAPI",
		BaseURL:   server.URL,
	}
	
	client, err := NewXTSClientImpl(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	
	// Set a test token for authenticated requests
	client.token = "test_token"
	client.userID = "test_user_id"
	
	return server, client
}

// TestNewXTSClientImpl tests the creation of a new XTS Client implementation
func TestNewXTSClientImpl(t *testing.T) {
	// Test with valid config
	config := &common.XTSClientConfig{
		APIKey:    "test_api_key",
		SecretKey: "test_secret_key",
		Source:    "WEBAPI",
	}
	
	client, err := NewXTSClientImpl(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test_api_key", client.apiKey)
	assert.Equal(t, "test_secret_key", client.secretKey)
	assert.Equal(t, "WEBAPI", client.source)
	assert.Equal(t, "https://developers.symphonyfintech.in", client.baseURL)
	
	// Test with custom base URL
	config.BaseURL = "https://custom.api.url"
	client, err = NewXTSClientImpl(config)
	assert.NoError(t, err)
	assert.Equal(t, "https://custom.api.url", client.baseURL)
	
	// Test with nil config
	client, err = NewXTSClientImpl(nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	
	// Test with missing API key
	config = &common.XTSClientConfig{
		SecretKey: "test_secret_key",
	}
	client, err = NewXTSClientImpl(config)
	assert.Error(t, err)
	assert.Nil(t, client)
	
	// Test with missing secret key
	config = &common.XTSClientConfig{
		APIKey: "test_api_key",
	}
	client, err = NewXTSClientImpl(config)
	assert.Error(t, err)
	assert.Nil(t, client)
}

// TestLogin tests the Login method
func TestLogin(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/interactive/user/session", r.URL.Path)
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "test_api_key", query.Get("appKey"))
		assert.Equal(t, "test_secret_key", query.Get("secretKey"))
		assert.Equal(t, "WEBAPI", query.Get("source"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Login successful",
			"result": {
				"token": "test_token_response",
				"userID": "test_user_id_response",
				"isInvestorClient": true,
				"expiresIn": 3600
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test login with default credentials
	session, err := client.Login(nil)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "test_token_response", session.Token)
	assert.Equal(t, "test_user_id_response", session.UserID)
	
	// Test login with provided credentials
	credentials := &common.Credentials{
		APIKey:    "custom_api_key",
		SecretKey: "custom_secret_key",
	}
	
	// Update mock server to expect custom credentials
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "custom_api_key", query.Get("appKey"))
		assert.Equal(t, "custom_secret_key", query.Get("secretKey"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Login successful",
			"result": {
				"token": "custom_token_response",
				"userID": "custom_user_id_response",
				"isInvestorClient": false,
				"expiresIn": 3600
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	session, err = client.Login(credentials)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, "custom_token_response", session.Token)
	assert.Equal(t, "custom_user_id_response", session.UserID)
	
	// Test login failure
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := `{
			"type": "error",
			"code": 401,
			"description": "Invalid credentials"
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	session, err = client.Login(nil)
	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "Invalid credentials")
}

// TestLogout tests the Logout method
func TestLogout(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/interactive/user/session", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Logout successful"
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test logout
	err := client.Logout()
	assert.NoError(t, err)
	assert.Empty(t, client.token)
	assert.Empty(t, client.userID)
	
	// Test logout when not logged in
	err = client.Logout()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test logout failure
	client.token = "test_token" // Reset token for test
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := `{
			"type": "error",
			"code": 500,
			"description": "Internal server error"
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	err = client.Logout()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Internal server error")
}

// TestGetOrderBook tests the GetOrderBook method
func TestGetOrderBook(t *testing.T) {
	// Setup mock server for investor client
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Ensure no clientID parameter for investor client
		query := r.URL.Query()
		assert.Empty(t, query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order book retrieved successfully",
			"result": [
				{
					"OrderID": "order1",
					"ExchangeOrderID": "ex_order1",
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"OrderSide": "BUY",
					"OrderType": "LIMIT",
					"ProductType": "MIS",
					"TimeInForce": "DAY",
					"OrderQuantity": 10,
					"FilledQuantity": 5,
					"RemainingQuantity": 5,
					"LimitPrice": 100.5,
					"StopPrice": 0,
					"OrderStatus": "PARTIALLY_FILLED",
					"OrderTimestamp": 1617345678000,
					"LastUpdateTimestamp": 1617345679000,
					"CancelTimestamp": 0,
					"ClientID": ""
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as investor
	client.isInvestor = true
	
	// Test get order book for investor client
	orderBook, err := client.GetOrderBook("")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.Len(t, orderBook.Orders, 1)
	assert.Equal(t, "order1", orderBook.Orders[0].OrderID)
	assert.Equal(t, "ex_order1", orderBook.Orders[0].ExchangeOrderID)
	assert.Equal(t, "NSECM", orderBook.Orders[0].ExchangeSegment)
	assert.Equal(t, "123456", orderBook.Orders[0].ExchangeInstrumentID)
	assert.Equal(t, "BUY", orderBook.Orders[0].OrderSide)
	assert.Equal(t, "LIMIT", orderBook.Orders[0].OrderType)
	assert.Equal(t, "MIS", orderBook.Orders[0].ProductType)
	assert.Equal(t, "DAY", orderBook.Orders[0].TimeInForce)
	assert.Equal(t, 10, orderBook.Orders[0].OrderQuantity)
	assert.Equal(t, 5, orderBook.Orders[0].FilledQuantity)
	assert.Equal(t, 5, orderBook.Orders[0].RemainingQuantity)
	assert.Equal(t, 100.5, orderBook.Orders[0].LimitPrice)
	assert.Equal(t, 0.0, orderBook.Orders[0].StopPrice)
	assert.Equal(t, "PARTIALLY_FILLED", orderBook.Orders[0].OrderStatus)
	assert.Equal(t, int64(1617345678000), orderBook.Orders[0].OrderTimestamp)
	assert.Equal(t, int64(1617345679000), orderBook.Orders[0].LastUpdateTimestamp)
	assert.Equal(t, int64(0), orderBook.Orders[0].CancelTimestamp)
	assert.Equal(t, "", orderBook.Orders[0].ClientID)
	
	// Setup mock server for dealer client
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		
		// Check for clientID parameter for dealer client
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order book retrieved successfully",
			"result": [
				{
					"OrderID": "order2",
					"ExchangeOrderID": "ex_order2",
					"ExchangeSegment": "NSEFO",
					"ExchangeInstrumentID": "789012",
					"OrderSide": "SELL",
					"OrderType": "MARKET",
					"ProductType": "NRML",
					"TimeInForce": "DAY",
					"OrderQuantity": 20,
					"FilledQuantity": 20,
					"RemainingQuantity": 0,
					"LimitPrice": 0,
					"StopPrice": 0,
					"OrderStatus": "FILLED",
					"OrderTimestamp": 1617345680000,
					"LastUpdateTimestamp": 1617345681000,
					"CancelTimestamp": 0,
					"ClientID": "client123"
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	// Test get order book for dealer client
	orderBook, err = client.GetOrderBook("client123")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.Len(t, orderBook.Orders, 1)
	assert.Equal(t, "order2", orderBook.Orders[0].OrderID)
	assert.Equal(t, "client123", orderBook.Orders[0].ClientID)
	
	// Test get order book when not logged in
	client.token = ""
	orderBook, err = client.GetOrderBook("client123")
	assert.Error(t, err)
	assert.Nil(t, orderBook)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestPlaceOrder tests the PlaceOrder method
func TestPlaceOrder(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "NSECM", query.Get("exchangeSegment"))
		assert.Equal(t, "123456", query.Get("exchangeInstrumentID"))
		assert.Equal(t, "MIS", query.Get("productType"))
		assert.Equal(t, "LIMIT", query.Get("orderType"))
		assert.Equal(t, "BUY", query.Get("orderSide"))
		assert.Equal(t, "DAY", query.Get("timeInForce"))
		assert.Equal(t, "0", query.Get("disclosedQuantity"))
		assert.Equal(t, "10", query.Get("orderQuantity"))
		assert.Equal(t, "100.500000", query.Get("limitPrice"))
		assert.Equal(t, "0.000000", query.Get("stopPrice"))
		assert.Equal(t, "test123", query.Get("orderUniqueIdentifier"))
		assert.Equal(t, "WEBAPI", query.Get("apiOrderSource"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order placed successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "PLACED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test place order
	order := &common.Order{
		ExchangeSegment:       "NSECM",
		ExchangeInstrumentID:  "123456",
		ProductType:           "MIS",
		OrderType:             "LIMIT",
		OrderSide:             "BUY",
		TimeInForce:           "DAY",
		DisclosedQuantity:     0,
		OrderQuantity:         10,
		LimitPrice:            100.5,
		StopPrice:             0,
		OrderUniqueIdentifier: "test123",
		APIOrderSource:        "WEBAPI",
	}
	
	orderResponse, err := client.PlaceOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "ex_order123", orderResponse.ExchangeOrderID)
	assert.Equal(t, "PLACED", orderResponse.Status)
	assert.Empty(t, orderResponse.RejectionReason)
	
	// Test place order with default apiOrderSource
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "WEBAPI", query.Get("apiOrderSource"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order placed successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "PLACED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	order.APIOrderSource = ""
	orderResponse, err = client.PlaceOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	
	// Test place order with rejected order
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order rejected",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "",
				"OrderStatus": "REJECTED",
				"OrderRejected": true,
				"RejectReason": "Insufficient funds"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	orderResponse, err = client.PlaceOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "REJECTED", orderResponse.Status)
	assert.Equal(t, "Insufficient funds", orderResponse.RejectionReason)
	
	// Test place order when not logged in
	client.token = ""
	orderResponse, err = client.PlaceOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test place order with nil order
	client.token = "test_token"
	orderResponse, err = client.PlaceOrder(nil)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order is required")
}

// TestGetDealerOrderBook tests the GetDealerOrderBook method
func TestGetDealerOrderBook(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/orders/dealerorderbook", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Dealer order book retrieved successfully",
			"result": [
				{
					"OrderID": "order1",
					"ExchangeOrderID": "ex_order1",
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"OrderSide": "BUY",
					"OrderType": "LIMIT",
					"ProductType": "MIS",
					"TimeInForce": "DAY",
					"OrderQuantity": 10,
					"FilledQuantity": 5,
					"RemainingQuantity": 5,
					"LimitPrice": 100.5,
					"StopPrice": 0,
					"OrderStatus": "PARTIALLY_FILLED",
					"OrderTimestamp": 1617345678000,
					"LastUpdateTimestamp": 1617345679000,
					"CancelTimestamp": 0,
					"ClientID": "client123"
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	// Test get dealer order book
	orderBook, err := client.GetDealerOrderBook("client123")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.Len(t, orderBook.Orders, 1)
	assert.Equal(t, "order1", orderBook.Orders[0].OrderID)
	assert.Equal(t, "client123", orderBook.Orders[0].ClientID)
	
	// Test get dealer order book as investor client
	client.isInvestor = true
	orderBook, err = client.GetDealerOrderBook("client123")
	assert.Error(t, err)
	assert.Nil(t, orderBook)
	assert.Contains(t, err.Error(), "dealer endpoints are not available for investor clients")
	
	// Test get dealer order book when not logged in
	client.isInvestor = false
	client.token = ""
	orderBook, err = client.GetDealerOrderBook("client123")
	assert.Error(t, err)
	assert.Nil(t, orderBook)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestGetDealerTrades tests the GetDealerTrades method
func TestGetDealerTrades(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/orders/dealertradebook", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Dealer trades retrieved successfully",
			"result": [
				{
					"OrderID": "order1",
					"ExchangeOrderID": "ex_order1",
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"OrderSide": "BUY",
					"OrderType": "LIMIT",
					"ProductType": "MIS",
					"TimeInForce": "DAY",
					"OrderQuantity": 10,
					"FilledQuantity": 10,
					"RemainingQuantity": 0,
					"LimitPrice": 100.5,
					"StopPrice": 0,
					"OrderStatus": "FILLED",
					"OrderTimestamp": 1617345678000,
					"LastUpdateTimestamp": 1617345679000,
					"CancelTimestamp": 0,
					"ClientID": "client123",
					"TradeID": "trade1",
					"TradePrice": 100.5,
					"TradeQuantity": 10,
					"TradeTimestamp": 1617345680000
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	// Test get dealer trades
	trades, err := client.GetDealerTrades("client123")
	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.Len(t, trades, 1)
	assert.Equal(t, "order1", trades[0].OrderID)
	assert.Equal(t, "client123", trades[0].ClientID)
	
	// Test get dealer trades as investor client
	client.isInvestor = true
	trades, err = client.GetDealerTrades("client123")
	assert.Error(t, err)
	assert.Nil(t, trades)
	assert.Contains(t, err.Error(), "dealer endpoints are not available for investor clients")
	
	// Test get dealer trades when not logged in
	client.isInvestor = false
	client.token = ""
	trades, err = client.GetDealerTrades("client123")
	assert.Error(t, err)
	assert.Nil(t, trades)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestGetDealerPositions tests the GetDealerPositions method
func TestGetDealerPositions(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/portfolio/dealerpositions", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Dealer positions retrieved successfully",
			"result": [
				{
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"ProductType": "MIS",
					"Quantity": 10,
					"BuyQuantity": 10,
					"SellQuantity": 0,
					"NetQuantity": 10,
					"AveragePrice": 100.5,
					"LastPrice": 101.0,
					"RealizedProfit": 0,
					"UnrealizedProfit": 5.0,
					"ClientID": "client123"
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	// Test get dealer positions
	positions, err := client.GetDealerPositions("client123")
	assert.NoError(t, err)
	assert.NotNil(t, positions)
	assert.Len(t, positions, 1)
	assert.Equal(t, "NSECM", positions[0].ExchangeSegment)
	assert.Equal(t, "123456", positions[0].ExchangeInstrumentID)
	assert.Equal(t, "MIS", positions[0].ProductType)
	assert.Equal(t, 10, positions[0].Quantity)
	assert.Equal(t, 10, positions[0].BuyQuantity)
	assert.Equal(t, 0, positions[0].SellQuantity)
	assert.Equal(t, 10, positions[0].NetQuantity)
	assert.Equal(t, 100.5, positions[0].AveragePrice)
	assert.Equal(t, 101.0, positions[0].LastPrice)
	assert.Equal(t, 0.0, positions[0].RealizedProfit)
	assert.Equal(t, 5.0, positions[0].UnrealizedProfit)
	assert.Equal(t, "client123", positions[0].ClientID)
	
	// Test get dealer positions as investor client
	client.isInvestor = true
	positions, err = client.GetDealerPositions("client123")
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "dealer endpoints are not available for investor clients")
	
	// Test get dealer positions when not logged in
	client.isInvestor = false
	client.token = ""
	positions, err = client.GetDealerPositions("client123")
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestPlaceDealerOrder tests the PlaceDealerOrder method
func TestPlaceDealerOrder(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "NSECM", query.Get("exchangeSegment"))
		assert.Equal(t, "123456", query.Get("exchangeInstrumentID"))
		assert.Equal(t, "MIS", query.Get("productType"))
		assert.Equal(t, "LIMIT", query.Get("orderType"))
		assert.Equal(t, "BUY", query.Get("orderSide"))
		assert.Equal(t, "DAY", query.Get("timeInForce"))
		assert.Equal(t, "0", query.Get("disclosedQuantity"))
		assert.Equal(t, "10", query.Get("orderQuantity"))
		assert.Equal(t, "100.500000", query.Get("limitPrice"))
		assert.Equal(t, "0.000000", query.Get("stopPrice"))
		assert.Equal(t, "test123", query.Get("orderUniqueIdentifier"))
		assert.Equal(t, "WEBAPI", query.Get("apiOrderSource"))
		assert.Equal(t, "client123", query.Get("clientID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order placed successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "PLACED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	// Test place dealer order
	order := &common.Order{
		ExchangeSegment:       "NSECM",
		ExchangeInstrumentID:  "123456",
		ProductType:           "MIS",
		OrderType:             "LIMIT",
		OrderSide:             "BUY",
		TimeInForce:           "DAY",
		DisclosedQuantity:     0,
		OrderQuantity:         10,
		LimitPrice:            100.5,
		StopPrice:             0,
		OrderUniqueIdentifier: "test123",
		APIOrderSource:        "WEBAPI",
		ClientID:              "client123",
	}
	
	orderResponse, err := client.PlaceDealerOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "ex_order123", orderResponse.ExchangeOrderID)
	assert.Equal(t, "PLACED", orderResponse.Status)
	assert.Empty(t, orderResponse.RejectionReason)
	
	// Test place dealer order as investor client
	client.isInvestor = true
	orderResponse, err = client.PlaceDealerOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "dealer endpoints are not available for investor clients")
	
	// Test place dealer order when not logged in
	client.isInvestor = false
	client.token = ""
	orderResponse, err = client.PlaceDealerOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test place dealer order with nil order
	client.token = "test_token"
	orderResponse, err = client.PlaceDealerOrder(nil)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order is required")
	
	// Test place dealer order without client ID
	order.ClientID = ""
	orderResponse, err = client.PlaceDealerOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "client ID is required for dealer orders")
}
