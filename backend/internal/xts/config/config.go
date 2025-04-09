package config

import (
	"encoding/json"
	"os"
	"time"
)

// XTSConfig represents the configuration for XTS API connectivity
type XTSConfig struct {
	BaseURL     string        `json:"baseUrl"`
	APIKey      string        `json:"apiKey"`
	SecretKey   string        `json:"secretKey"`
	Source      string        `json:"source"`
	Timeout     time.Duration `json:"timeout"`
	RetryCount  int           `json:"retryCount"`
	RetryDelay  time.Duration `json:"retryDelay"`
	DisableSSL  bool          `json:"disableSSL"`
}

// NewXTSConfig creates a new XTS configuration with default values
func NewXTSConfig() *XTSConfig {
	return &XTSConfig{
		BaseURL:    "https://xts-api.trading.com",
		Timeout:    7 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
		DisableSSL: false,
		Source:     "WEB",
	}
}

// LoadFromFile loads configuration from a JSON file
func (c *XTSConfig) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

// SaveToFile saves configuration to a JSON file
func (c *XTSConfig) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// Validate checks if the configuration is valid
func (c *XTSConfig) Validate() error {
	if c.BaseURL == "" {
		return ErrEmptyBaseURL
	}
	if c.APIKey == "" {
		return ErrEmptyAPIKey
	}
	if c.SecretKey == "" {
		return ErrEmptySecretKey
	}
	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}

// Routes defines the API endpoints for XTS
type Routes struct {
	// Interactive API endpoints
	InteractivePrefix     string
	UserLogin             string
	UserLogout            string
	UserProfile           string
	UserBalance           string
	Orders                string
	Trades                string
	OrderStatus           string
	OrderPlace            string
	BracketOrderPlace     string
	BracketOrderModify    string
	BracketOrderCancel    string
	OrderPlaceCover       string
	OrderExitCover        string
	OrderModify           string
	OrderCancel           string
	OrderCancelAll        string
	OrderHistory          string
	PortfolioPositions    string
	PortfolioHoldings     string
	PortfolioPositionsConvert string
	PortfolioSquareoff    string
	PortfolioDealerPositions string
	OrderDealerStatus     string
	DealerTrades          string

	// Market API endpoints
	MarketDataPrefix      string
	MarketLogin           string
	MarketLogout          string
	MarketConfig          string
	MarketInstrumentsMaster string
	MarketInstrumentsSubscription string
	MarketInstrumentsUnsubscription string
	MarketInstrumentsOHLC string
	MarketInstrumentsIndexList string
	MarketInstrumentsQuotes string
	MarketSearchInstrumentsByID string
	MarketSearchInstrumentsByString string
	MarketInstrumentsSeries string
	MarketInstrumentsEquitySymbol string
	MarketInstrumentsFutureSymbol string
	MarketInstrumentsOptionSymbol string
	MarketInstrumentsOptionType string
	MarketInstrumentsExpiryDate string
}

// DefaultRoutes returns the default API routes for XTS
func DefaultRoutes() *Routes {
	return &Routes{
		// Interactive API endpoints
		InteractivePrefix:     "interactive",
		UserLogin:             "/interactive/user/session",
		UserLogout:            "/interactive/user/session",
		UserProfile:           "/interactive/user/profile",
		UserBalance:           "/interactive/user/balance",
		Orders:                "/interactive/orders",
		Trades:                "/interactive/orders/trades",
		OrderStatus:           "/interactive/orders",
		OrderPlace:            "/interactive/orders",
		BracketOrderPlace:     "/interactive/orders/bracket",
		BracketOrderModify:    "/interactive/orders/bracket",
		BracketOrderCancel:    "/interactive/orders/bracket",
		OrderPlaceCover:       "/interactive/orders/cover",
		OrderExitCover:        "/interactive/orders/cover",
		OrderModify:           "/interactive/orders",
		OrderCancel:           "/interactive/orders",
		OrderCancelAll:        "/interactive/orders/cancelall",
		OrderHistory:          "/interactive/orders",
		PortfolioPositions:    "/interactive/portfolio/positions",
		PortfolioHoldings:     "/interactive/portfolio/holdings",
		PortfolioPositionsConvert: "/interactive/portfolio/positions/convert",
		PortfolioSquareoff:    "/interactive/portfolio/squareoff",
		PortfolioDealerPositions: "interactive/portfolio/dealerpositions",
		OrderDealerStatus:     "/interactive/orders/dealerorderbook",
		DealerTrades:          "/interactive/orders/dealertradebook",

		// Market API endpoints
		MarketDataPrefix:      "apimarketdata",
		MarketLogin:           "/apimarketdata/auth/login",
		MarketLogout:          "/apimarketdata/auth/logout",
		MarketConfig:          "/apimarketdata/config/clientConfig",
		MarketInstrumentsMaster: "/apimarketdata/instruments/master",
		MarketInstrumentsSubscription: "/apimarketdata/instruments/subscription",
		MarketInstrumentsUnsubscription: "/apimarketdata/instruments/subscription",
		MarketInstrumentsOHLC: "/apimarketdata/instruments/ohlc",
		MarketInstrumentsIndexList: "/apimarketdata/instruments/indexlist",
		MarketInstrumentsQuotes: "/apimarketdata/instruments/quotes",
		MarketSearchInstrumentsByID: "/apimarketdata/search/instrumentsbyid",
		MarketSearchInstrumentsByString: "/apimarketdata/search/instruments",
		MarketInstrumentsSeries: "/apimarketdata/instruments/instrument/series",
		MarketInstrumentsEquitySymbol: "/apimarketdata/instruments/instrument/symbol",
		MarketInstrumentsFutureSymbol: "/apimarketdata/instruments/instrument/futureSymbol",
		MarketInstrumentsOptionSymbol: "/apimarketdata/instruments/instrument/optionsymbol",
		MarketInstrumentsOptionType: "/apimarketdata/instruments/instrument/optionType",
		MarketInstrumentsExpiryDate: "/apimarketdata/instruments/instrument/expiryDate",
	}
}
