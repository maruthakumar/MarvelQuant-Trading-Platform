package zerodha

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
)

// TestNewZerodhaAdapter tests the creation of a new Zerodha adapter
func TestNewZerodhaAdapter(t *testing.T) {
	// Test with valid config
	config := &common.ZerodhaConfig{
		APIKey:      "test_api_key",
		APISecret:   "test_api_secret",
		RedirectURL: "https://example.com/redirect",
	}
	
	adapter, err := NewZerodhaAdapter(config)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.Equal(t, "test_api_key", adapter.apiKey)
	assert.Equal(t, "test_api_secret", adapter.apiSecret)
	assert.Equal(t, "https://example.com/redirect", adapter.redirectURL)
	assert.Equal(t, "https://api.kite.trade", adapter.baseURL)
	
	// Test with custom base URL
	config.BaseURL = "https://custom.api.url"
	adapter, err = NewZerodhaAdapter(config)
	assert.NoError(t, err)
	assert.Equal(t, "https://custom.api.url", adapter.baseURL)
	
	// Test with nil config
	adapter, err = NewZerodhaAdapter(nil)
	assert.Error(t, err)
	assert.Nil(t, adapter)
	
	// Test with missing API key
	config = &common.ZerodhaConfig{
		APISecret: "test_api_secret",
	}
	adapter, err = NewZerodhaAdapter(config)
	assert.Error(t, err)
	assert.Nil(t, adapter)
	
	// Test with missing API secret
	config = &common.ZerodhaConfig{
		APIKey: "test_api_key",
	}
	adapter, err = NewZerodhaAdapter(config)
	assert.Error(t, err)
	assert.Nil(t, adapter)
}

// TestMapExchangeSegment tests the mapExchangeSegment function
func TestMapExchangeSegment(t *testing.T) {
	assert.Equal(t, "NSE", mapExchangeSegment("NSECM"))
	assert.Equal(t, "BSE", mapExchangeSegment("BSECM"))
	assert.Equal(t, "NFO", mapExchangeSegment("NSEFO"))
	assert.Equal(t, "BFO", mapExchangeSegment("BSEFO"))
	assert.Equal(t, "CDS", mapExchangeSegment("NSECD"))
	assert.Equal(t, "MCX", mapExchangeSegment("MCXFO"))
	assert.Equal(t, "UNKNOWN", mapExchangeSegment("UNKNOWN"))
}

// TestMapExchange tests the mapExchange function
func TestMapExchange(t *testing.T) {
	assert.Equal(t, "NSECM", mapExchange("NSE"))
	assert.Equal(t, "BSECM", mapExchange("BSE"))
	assert.Equal(t, "NSEFO", mapExchange("NFO"))
	assert.Equal(t, "BSEFO", mapExchange("BFO"))
	assert.Equal(t, "NSECD", mapExchange("CDS"))
	assert.Equal(t, "MCXFO", mapExchange("MCX"))
	assert.Equal(t, "UNKNOWN", mapExchange("UNKNOWN"))
}

// TestMapOrderSide tests the mapOrderSide function
func TestMapOrderSide(t *testing.T) {
	assert.Equal(t, "BUY", mapOrderSide("BUY"))
	assert.Equal(t, "SELL", mapOrderSide("SELL"))
	assert.Equal(t, "UNKNOWN", mapOrderSide("UNKNOWN"))
}

// TestMapTransactionType tests the mapTransactionType function
func TestMapTransactionType(t *testing.T) {
	assert.Equal(t, "BUY", mapTransactionType("BUY"))
	assert.Equal(t, "SELL", mapTransactionType("SELL"))
	assert.Equal(t, "UNKNOWN", mapTransactionType("UNKNOWN"))
}

// TestMapProductType tests the mapProductType function
func TestMapProductType(t *testing.T) {
	assert.Equal(t, "MIS", mapProductType("MIS"))
	assert.Equal(t, "NRML", mapProductType("NRML"))
	assert.Equal(t, "CNC", mapProductType("CNC"))
	assert.Equal(t, "UNKNOWN", mapProductType("UNKNOWN"))
}

// TestMapZerodhaProductType tests the mapZerodhaProductType function
func TestMapZerodhaProductType(t *testing.T) {
	assert.Equal(t, "MIS", mapZerodhaProductType("MIS"))
	assert.Equal(t, "NRML", mapZerodhaProductType("NRML"))
	assert.Equal(t, "CNC", mapZerodhaProductType("CNC"))
	assert.Equal(t, "UNKNOWN", mapZerodhaProductType("UNKNOWN"))
}

// TestMapOrderType tests the mapOrderType function
func TestMapOrderType(t *testing.T) {
	assert.Equal(t, "MARKET", mapOrderType("MARKET"))
	assert.Equal(t, "LIMIT", mapOrderType("LIMIT"))
	assert.Equal(t, "SL", mapOrderType("SL"))
	assert.Equal(t, "SL-M", mapOrderType("SL-M"))
	assert.Equal(t, "UNKNOWN", mapOrderType("UNKNOWN"))
}

// TestMapZerodhaOrderType tests the mapZerodhaOrderType function
func TestMapZerodhaOrderType(t *testing.T) {
	assert.Equal(t, "MARKET", mapZerodhaOrderType("MARKET"))
	assert.Equal(t, "LIMIT", mapZerodhaOrderType("LIMIT"))
	assert.Equal(t, "SL", mapZerodhaOrderType("SL"))
	assert.Equal(t, "SL-M", mapZerodhaOrderType("SL-M"))
	assert.Equal(t, "UNKNOWN", mapZerodhaOrderType("UNKNOWN"))
}

// TestMapTimeInForce tests the mapTimeInForce function
func TestMapTimeInForce(t *testing.T) {
	assert.Equal(t, "DAY", mapTimeInForce("DAY"))
	assert.Equal(t, "IOC", mapTimeInForce("IOC"))
	assert.Equal(t, "UNKNOWN", mapTimeInForce("UNKNOWN"))
}

// TestMapZerodhaValidity tests the mapZerodhaValidity function
func TestMapZerodhaValidity(t *testing.T) {
	assert.Equal(t, "DAY", mapZerodhaValidity("DAY"))
	assert.Equal(t, "IOC", mapZerodhaValidity("IOC"))
	assert.Equal(t, "UNKNOWN", mapZerodhaValidity("UNKNOWN"))
}

// TestMapZerodhaStatus tests the mapZerodhaStatus function
func TestMapZerodhaStatus(t *testing.T) {
	assert.Equal(t, "FILLED", mapZerodhaStatus("COMPLETE"))
	assert.Equal(t, "REJECTED", mapZerodhaStatus("REJECTED"))
	assert.Equal(t, "CANCELLED", mapZerodhaStatus("CANCELLED"))
	assert.Equal(t, "PENDING", mapZerodhaStatus("PENDING"))
	assert.Equal(t, "OPEN", mapZerodhaStatus("OPEN"))
	assert.Equal(t, "UNKNOWN", mapZerodhaStatus("UNKNOWN"))
}
