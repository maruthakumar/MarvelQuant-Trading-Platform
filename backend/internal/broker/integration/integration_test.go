package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trading-platform/backend/internal/broker/common"
	"github.com/trading-platform/backend/internal/broker/factory"
)

// TestBrokerFactoryIntegration tests the broker factory integration
func TestBrokerFactoryIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=1 to run.")
	}

	// Test XTS Pro client creation
	xtsProConfig := &common.BrokerConfig{
		BrokerType: common.BrokerTypeXTSPro,
		XTSPro: &common.XTSProConfig{
			APIKey:    "test_api_key",
			SecretKey: "test_secret_key",
			Source:    "WEBAPI",
		},
	}

	xtsProClient, err := factory.NewBrokerClient(xtsProConfig)
	assert.NoError(t, err)
	assert.NotNil(t, xtsProClient)

	// Test XTS Client client creation
	xtsClientConfig := &common.BrokerConfig{
		BrokerType: common.BrokerTypeXTSClient,
		XTSClient: &common.XTSClientConfig{
			APIKey:    "test_api_key",
			SecretKey: "test_secret_key",
			Source:    "WEBAPI",
		},
	}

	xtsClientClient, err := factory.NewBrokerClient(xtsClientConfig)
	assert.NoError(t, err)
	assert.NotNil(t, xtsClientClient)

	// Test Zerodha client creation
	zerodhaConfig := &common.BrokerConfig{
		BrokerType: common.BrokerTypeZerodha,
		Zerodha: &common.ZerodhaConfig{
			APIKey:      "test_api_key",
			APISecret:   "test_api_secret",
			RedirectURL: "https://example.com/redirect",
		},
	}

	zerodhaClient, err := factory.NewBrokerClient(zerodhaConfig)
	assert.NoError(t, err)
	assert.NotNil(t, zerodhaClient)
}

// TestXTSClientIntegration tests the XTS Client integration
func TestXTSClientIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=1 to run.")
	}

	// Skip if XTS Client credentials are not provided
	apiKey := os.Getenv("XTS_CLIENT_API_KEY")
	secretKey := os.Getenv("XTS_CLIENT_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		t.Skip("Skipping XTS Client integration tests. Set XTS_CLIENT_API_KEY and XTS_CLIENT_SECRET_KEY to run.")
	}

	// Create XTS Client client
	config := &common.BrokerConfig{
		BrokerType: common.BrokerTypeXTSClient,
		XTSClient: &common.XTSClientConfig{
			APIKey:    apiKey,
			SecretKey: secretKey,
			Source:    "WEBAPI",
		},
	}

	client, err := factory.NewBrokerClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test login
	session, err := client.Login(nil)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.Token)
	assert.NotEmpty(t, session.UserID)

	// Test get order book
	orderBook, err := client.GetOrderBook("")
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)

	// Test get positions
	positions, err := client.GetPositions("")
	assert.NoError(t, err)
	assert.NotNil(t, positions)

	// Test get holdings
	holdings, err := client.GetHoldings("")
	assert.NoError(t, err)
	assert.NotNil(t, holdings)

	// Test logout
	err = client.Logout()
	assert.NoError(t, err)
}

// TestZerodhaIntegration tests the Zerodha integration
func TestZerodhaIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=1 to run.")
	}

	// Skip if Zerodha credentials are not provided
	apiKey := os.Getenv("ZERODHA_API_KEY")
	apiSecret := os.Getenv("ZERODHA_API_SECRET")
	requestToken := os.Getenv("ZERODHA_REQUEST_TOKEN")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping Zerodha integration tests. Set ZERODHA_API_KEY and ZERODHA_API_SECRET to run.")
	}

	// Create Zerodha client
	config := &common.BrokerConfig{
		BrokerType: common.BrokerTypeZerodha,
		Zerodha: &common.ZerodhaConfig{
			APIKey:      apiKey,
			APISecret:   apiSecret,
			RedirectURL: "https://example.com/redirect",
		},
	}

	client, err := factory.NewBrokerClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test login (only if request token is provided)
	if requestToken != "" {
		credentials := &common.Credentials{
			TwoFactorCode: requestToken,
		}
		session, err := client.Login(credentials)
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.Token)
		assert.NotEmpty(t, session.UserID)

		// Test get order book
		orderBook, err := client.GetOrderBook("")
		assert.NoError(t, err)
		assert.NotNil(t, orderBook)

		// Test get positions
		positions, err := client.GetPositions("")
		assert.NoError(t, err)
		assert.NotNil(t, positions)

		// Test get holdings
		holdings, err := client.GetHoldings("")
		assert.NoError(t, err)
		assert.NotNil(t, holdings)

		// Test logout
		err = client.Logout()
		assert.NoError(t, err)
	} else {
		t.Log("Skipping Zerodha login tests. Set ZERODHA_REQUEST_TOKEN to run.")
	}
}
