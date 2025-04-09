package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
	"net/http"
	"net/http/httptest"
)

// TestModifyOrder tests the ModifyOrder method
func TestModifyOrder(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "order123", query.Get("appOrderID"))
		assert.Equal(t, "20", query.Get("orderQuantity"))
		assert.Equal(t, "LIMIT", query.Get("orderType"))
		assert.Equal(t, "105.500000", query.Get("limitPrice"))
		assert.Equal(t, "0.000000", query.Get("stopPrice"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order modified successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "MODIFIED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test modify order
	order := &common.ModifyOrder{
		OrderID:         "order123",
		OrderType:       "LIMIT",
		OrderQuantity:   20,
		LimitPrice:      105.5,
		StopPrice:       0,
	}
	
	orderResponse, err := client.ModifyOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "ex_order123", orderResponse.ExchangeOrderID)
	assert.Equal(t, "MODIFIED", orderResponse.Status)
	assert.Empty(t, orderResponse.RejectionReason)
	
	// Test modify order with rejected order
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
				"RejectReason": "Invalid order quantity"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	orderResponse, err = client.ModifyOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "REJECTED", orderResponse.Status)
	assert.Equal(t, "Invalid order quantity", orderResponse.RejectionReason)
	
	// Test modify order when not logged in
	client.token = ""
	orderResponse, err = client.ModifyOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test modify order with nil order
	client.token = "test_token"
	orderResponse, err = client.ModifyOrder(nil)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order is required")
	
	// Test modify order with missing order ID
	order.OrderID = ""
	orderResponse, err = client.ModifyOrder(order)
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order ID is required")
}

// TestCancelOrder tests the CancelOrder method
func TestCancelOrder(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/interactive/orders", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "order123", query.Get("appOrderID"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order cancelled successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "CANCELLED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test cancel order
	orderResponse, err := client.CancelOrder("order123", "")
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "order123", orderResponse.OrderID)
	assert.Equal(t, "ex_order123", orderResponse.ExchangeOrderID)
	assert.Equal(t, "CANCELLED", orderResponse.Status)
	assert.Empty(t, orderResponse.RejectionReason)
	
	// Test cancel order with client ID
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "order123", query.Get("appOrderID"))
		assert.Equal(t, "client123", query.Get("clientID"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order cancelled successfully",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "ex_order123",
				"OrderStatus": "CANCELLED",
				"OrderRejected": false,
				"RejectReason": ""
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	orderResponse, err = client.CancelOrder("order123", "client123")
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	
	// Test cancel order with rejected cancellation
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Order cancellation rejected",
			"result": {
				"AppOrderID": "order123",
				"OrderGeneratedID": "",
				"OrderStatus": "REJECTED",
				"OrderRejected": true,
				"RejectReason": "Order already executed"
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	orderResponse, err = client.CancelOrder("order123", "")
	assert.NoError(t, err)
	assert.NotNil(t, orderResponse)
	assert.Equal(t, "REJECTED", orderResponse.Status)
	assert.Equal(t, "Order already executed", orderResponse.RejectionReason)
	
	// Test cancel order when not logged in
	client.token = ""
	orderResponse, err = client.CancelOrder("order123", "")
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "not logged in")
	
	// Test cancel order with missing order ID
	client.token = "test_token"
	orderResponse, err = client.CancelOrder("", "")
	assert.Error(t, err)
	assert.Nil(t, orderResponse)
	assert.Contains(t, err.Error(), "order ID is required")
}

// TestGetPositions tests the GetPositions method
func TestGetPositions(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/portfolio/positions", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Positions retrieved successfully",
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
					"ClientID": ""
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test get positions
	positions, err := client.GetPositions("")
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
	assert.Equal(t, "", positions[0].ClientID)
	
	// Test get positions with client ID
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Positions retrieved successfully",
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
	
	positions, err = client.GetPositions("client123")
	assert.NoError(t, err)
	assert.NotNil(t, positions)
	assert.Len(t, positions, 1)
	assert.Equal(t, "client123", positions[0].ClientID)
	
	// Test get positions when not logged in
	client.token = ""
	positions, err = client.GetPositions("")
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestGetHoldings tests the GetHoldings method
func TestGetHoldings(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/interactive/portfolio/holdings", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Holdings retrieved successfully",
			"result": [
				{
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"TradingSymbol": "RELIANCE",
					"ISIN": "INE002A01018",
					"Quantity": 10,
					"AveragePrice": 2000.5,
					"LastPrice": 2050.0,
					"RealizedProfit": 0,
					"UnrealizedProfit": 495.0,
					"ClientID": ""
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test get holdings
	holdings, err := client.GetHoldings("")
	assert.NoError(t, err)
	assert.NotNil(t, holdings)
	assert.Len(t, holdings, 1)
	assert.Equal(t, "NSECM", holdings[0].ExchangeSegment)
	assert.Equal(t, "123456", holdings[0].ExchangeInstrumentID)
	assert.Equal(t, "RELIANCE", holdings[0].TradingSymbol)
	assert.Equal(t, "INE002A01018", holdings[0].ISIN)
	assert.Equal(t, 10, holdings[0].Quantity)
	assert.Equal(t, 2000.5, holdings[0].AveragePrice)
	assert.Equal(t, 2050.0, holdings[0].LastPrice)
	assert.Equal(t, 0.0, holdings[0].RealizedProfit)
	assert.Equal(t, 495.0, holdings[0].UnrealizedProfit)
	assert.Equal(t, "", holdings[0].ClientID)
	
	// Test get holdings with client ID
	server, client = setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "client123", query.Get("clientID"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Holdings retrieved successfully",
			"result": [
				{
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"TradingSymbol": "RELIANCE",
					"ISIN": "INE002A01018",
					"Quantity": 10,
					"AveragePrice": 2000.5,
					"LastPrice": 2050.0,
					"RealizedProfit": 0,
					"UnrealizedProfit": 495.0,
					"ClientID": "client123"
				}
			]
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Set client as dealer
	client.isInvestor = false
	
	holdings, err = client.GetHoldings("client123")
	assert.NoError(t, err)
	assert.NotNil(t, holdings)
	assert.Len(t, holdings, 1)
	assert.Equal(t, "client123", holdings[0].ClientID)
	
	// Test get holdings when not logged in
	client.token = ""
	holdings, err = client.GetHoldings("")
	assert.Error(t, err)
	assert.Nil(t, holdings)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestGetQuote tests the GetQuote method
func TestGetQuote(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/marketdata/instruments/quotes", r.URL.Path)
		assert.Equal(t, "test_token", r.Header.Get("Authorization"))
		
		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "NSECM:123456,NSEFO:789012", query.Get("instruments"))
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"type": "success",
			"code": 200,
			"description": "Quotes retrieved successfully",
			"result": {
				"NSECM:123456": {
					"ExchangeSegment": "NSECM",
					"ExchangeInstrumentID": "123456",
					"TradingSymbol": "RELIANCE",
					"LastPrice": 2050.0,
					"Open": 2000.0,
					"High": 2060.0,
					"Low": 1990.0,
					"Close": 2000.0,
					"Volume": 1000000,
					"BidPrice": 2049.0,
					"BidSize": 100,
					"AskPrice": 2051.0,
					"AskSize": 150,
					"Timestamp": 1617345678000
				},
				"NSEFO:789012": {
					"ExchangeSegment": "NSEFO",
					"ExchangeInstrumentID": "789012",
					"TradingSymbol": "RELIANCE25APR2100CE",
					"LastPrice": 50.0,
					"Open": 45.0,
					"High": 55.0,
					"Low": 40.0,
					"Close": 45.0,
					"Volume": 50000,
					"BidPrice": 49.5,
					"BidSize": 50,
					"AskPrice": 50.5,
					"AskSize": 75,
					"Timestamp": 1617345678000
				}
			}
		}`
		w.Write([]byte(response))
	})
	defer server.Close()
	
	// Test get quote
	quotes, err := client.GetQuote([]string{"NSECM:123456", "NSEFO:789012"})
	assert.NoError(t, err)
	assert.NotNil(t, quotes)
	assert.Len(t, quotes, 2)
	
	// Check RELIANCE quote
	quote, ok := quotes["NSECM:123456"]
	assert.True(t, ok)
	assert.Equal(t, "NSECM", quote.ExchangeSegment)
	assert.Equal(t, "123456", quote.ExchangeInstrumentID)
	assert.Equal(t, "RELIANCE", quote.TradingSymbol)
	assert.Equal(t, 2050.0, quote.LastPrice)
	assert.Equal(t, 2000.0, quote.Open)
	assert.Equal(t, 2060.0, quote.High)
	assert.Equal(t, 1990.0, quote.Low)
	assert.Equal(t, 2000.0, quote.Close)
	assert.Equal(t, int64(1000000), quote.Volume)
	assert.Equal(t, 2049.0, quote.BidPrice)
	assert.Equal(t, 100, quote.BidSize)
	assert.Equal(t, 2051.0, quote.AskPrice)
	assert.Equal(t, 150, quote.AskSize)
	assert.Equal(t, int64(1617345678000), quote.Timestamp)
	
	// Check RELIANCE option quote
	quote, ok = quotes["NSEFO:789012"]
	assert.True(t, ok)
	assert.Equal(t, "NSEFO", quote.ExchangeSegment)
	assert.Equal(t, "789012", quote.ExchangeInstrumentID)
	assert.Equal(t, "RELIANCE25APR2100CE", quote.TradingSymbol)
	
	// Test get quote with empty symbols
	quotes, err = client.GetQuote([]string{})
	assert.Error(t, err)
	assert.Nil(t, quotes)
	assert.Contains(t, err.Error(), "at least one symbol is required")
	
	// Test get quote when not logged in
	client.token = ""
	quotes, err = client.GetQuote([]string{"NSECM:123456"})
	assert.Error(t, err)
	assert.Nil(t, quotes)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestSubscribeToQuotes tests the SubscribeToQuotes method
func TestSubscribeToQuotes(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder since we're not actually implementing WebSocket in the tests
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()
	
	// Test subscribe to quotes
	quoteChan, err := client.SubscribeToQuotes([]string{"NSECM:123456"})
	assert.Error(t, err) // We expect an error since this is not implemented yet
	assert.NotNil(t, quoteChan)
	assert.Contains(t, err.Error(), "not yet implemented")
	
	// Test subscribe to quotes with empty symbols
	quoteChan, err = client.SubscribeToQuotes([]string{})
	assert.Error(t, err)
	assert.Nil(t, quoteChan)
	assert.Contains(t, err.Error(), "at least one symbol is required")
	
	// Test subscribe to quotes when not logged in
	client.token = ""
	quoteChan, err = client.SubscribeToQuotes([]string{"NSECM:123456"})
	assert.Error(t, err)
	assert.Nil(t, quoteChan)
	assert.Contains(t, err.Error(), "not logged in")
}

// TestUnsubscribeFromQuotes tests the UnsubscribeFromQuotes method
func TestUnsubscribeFromQuotes(t *testing.T) {
	// Setup mock server
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// This is a placeholder since we're not actually implementing WebSocket in the tests
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()
	
	// Test unsubscribe from quotes
	err := client.UnsubscribeFromQuotes([]string{"NSECM:123456"})
	assert.Error(t, err) // We expect an error since this is not implemented yet
	assert.Contains(t, err.Error(), "not yet implemented")
	
	// Test unsubscribe from quotes with empty symbols
	err = client.UnsubscribeFromQuotes([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one symbol is required")
	
	// Test unsubscribe from quotes when not logged in
	client.token = ""
	err = client.UnsubscribeFromQuotes([]string{"NSECM:123456"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")
}
