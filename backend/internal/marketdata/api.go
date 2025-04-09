package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// APIHandler handles API requests for market data
type APIHandler struct {
	marketDataService *MarketDataService
	realTimeManager   *RealTimeUpdateManager
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(marketDataService *MarketDataService, realTimeManager *RealTimeUpdateManager) *APIHandler {
	return &APIHandler{
		marketDataService: marketDataService,
		realTimeManager:   realTimeManager,
	}
}

// RegisterRoutes registers API routes
func (h *APIHandler) RegisterRoutes(router *mux.Router) {
	// Market data endpoints
	router.HandleFunc("/api/v1/market-data/symbols", h.GetSymbols).Methods("GET")
	router.HandleFunc("/api/v1/market-data/quote/{symbol}", h.GetQuote).Methods("GET")
	router.HandleFunc("/api/v1/market-data/quotes", h.GetQuotes).Methods("GET")
	
	// Historical data endpoints
	router.HandleFunc("/api/v1/market-data/historical/{symbol}", h.GetHistoricalData).Methods("GET")
	
	// Technical indicators endpoints
	router.HandleFunc("/api/v1/market-data/indicators/{indicator}/{symbol}", h.GetIndicator).Methods("GET")
	
	// WebSocket endpoint
	router.HandleFunc("/api/v1/market-data/stream", h.WebSocketHandler)
}

// GetSymbols handles requests for available symbols
func (h *APIHandler) GetSymbols(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would fetch symbols from a database or API
	// For now, we'll return a static list of symbols
	symbols := []map[string]string{
		{"symbol": "AAPL", "name": "Apple Inc.", "exchange": "NASDAQ"},
		{"symbol": "MSFT", "name": "Microsoft Corporation", "exchange": "NASDAQ"},
		{"symbol": "GOOG", "name": "Alphabet Inc.", "exchange": "NASDAQ"},
		{"symbol": "AMZN", "name": "Amazon.com, Inc.", "exchange": "NASDAQ"},
		{"symbol": "FB", "name": "Meta Platforms, Inc.", "exchange": "NASDAQ"},
		{"symbol": "TSLA", "name": "Tesla, Inc.", "exchange": "NASDAQ"},
		{"symbol": "NFLX", "name": "Netflix, Inc.", "exchange": "NASDAQ"},
		{"symbol": "NVDA", "name": "NVIDIA Corporation", "exchange": "NASDAQ"},
		{"symbol": "PYPL", "name": "PayPal Holdings, Inc.", "exchange": "NASDAQ"},
		{"symbol": "ADBE", "name": "Adobe Inc.", "exchange": "NASDAQ"},
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"symbols": symbols,
	})
}

// GetQuote handles requests for a single quote
func (h *APIHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	// Get symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	// Get market data
	data, err := h.marketDataService.GetMarketData(r.Context(), []string{symbol})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting market data: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if we got data for the symbol
	quote, ok := data[symbol]
	if !ok {
		http.Error(w, fmt.Sprintf("No data found for symbol: %s", symbol), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"quote":  quote,
	})
}

// GetQuotes handles requests for multiple quotes
func (h *APIHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	// Get symbols from query parameters
	symbolsParam := r.URL.Query().Get("symbols")
	if symbolsParam == "" {
		http.Error(w, "Missing symbols parameter", http.StatusBadRequest)
		return
	}

	// Split symbols
	symbols := splitCSV(symbolsParam)

	// Get market data
	data, err := h.marketDataService.GetMarketData(r.Context(), symbols)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting market data: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"quotes": data,
	})
}

// GetHistoricalData handles requests for historical data
func (h *APIHandler) GetHistoricalData(w http.ResponseWriter, r *http.Request) {
	// Get symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	// Get query parameters
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1d" // Default to daily
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Parse from and to dates
	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid from date: %s", fromStr), http.StatusBadRequest)
			return
		}
	} else {
		// Default to 30 days ago
		from = time.Now().AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid to date: %s", toStr), http.StatusBadRequest)
			return
		}
	} else {
		// Default to today
		to = time.Now()
	}

	// Get historical data
	data, err := h.marketDataService.GetHistoricalData(r.Context(), symbol, interval, from, to)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting historical data: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"symbol":   symbol,
			"interval": interval,
			"from":     from.Format("2006-01-02"),
			"to":       to.Format("2006-01-02"),
			"candles":  data,
		},
	})
}

// GetIndicator handles requests for technical indicators
func (h *APIHandler) GetIndicator(w http.ResponseWriter, r *http.Request) {
	// Get indicator and symbol from URL
	vars := mux.Vars(r)
	indicator := vars["indicator"]
	symbol := vars["symbol"]

	// Get query parameters
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1d" // Default to daily
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Parse from and to dates
	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid from date: %s", fromStr), http.StatusBadRequest)
			return
		}
	} else {
		// Default to 30 days ago
		from = time.Now().AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid to date: %s", toStr), http.StatusBadRequest)
			return
		}
	} else {
		// Default to today
		to = time.Now()
	}

	// Get parameters for the indicator
	params := make(map[string]interface{})
	for key, values := range r.URL.Query() {
		if key != "interval" && key != "from" && key != "to" {
			// Try to parse as int
			if intVal, err := strconv.Atoi(values[0]); err == nil {
				params[key] = intVal
			} else {
				// Try to parse as float
				if floatVal, err := strconv.ParseFloat(values[0], 64); err == nil {
					params[key] = floatVal
				} else {
					// Use as string
					params[key] = values[0]
				}
			}
		}
	}

	// Get historical data
	data, err := h.marketDataService.GetHistoricalData(r.Context(), symbol, interval, from, to)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting historical data: %v", err), http.StatusInternalServerError)
		return
	}

	// Calculate indicator
	// In a real implementation, this would use a technical indicator library
	// For now, we'll return a simple moving average if the indicator is "sma"
	var indicatorValues []map[string]interface{}

	if indicator == "sma" {
		// Get period parameter
		period := 14 // Default
		if periodParam, ok := params["period"]; ok {
			if periodInt, ok := periodParam.(int); ok {
				period = periodInt
			}
		}

		// Calculate SMA
		if len(data) < period {
			http.Error(w, fmt.Sprintf("Not enough data to calculate SMA with period %d", period), http.StatusBadRequest)
			return
		}

		for i := period - 1; i < len(data); i++ {
			sum := 0.0
			for j := 0; j < period; j++ {
				sum += data[i-j].Close
			}
			sma := sum / float64(period)

			indicatorValues = append(indicatorValues, map[string]interface{}{
				"timestamp": data[i].Timestamp,
				"value":     sma,
			})
		}
	} else {
		http.Error(w, fmt.Sprintf("Unsupported indicator: %s", indicator), http.StatusBadRequest)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"symbol":    symbol,
			"indicator": indicator,
			"interval":  interval,
			"from":      from.Format("2006-01-02"),
			"to":        to.Format("2006-01-02"),
			"params":    params,
			"values":    indicatorValues,
		},
	})
}

// WebSocketHandler handles WebSocket connections
func (h *APIHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error upgrading to WebSocket: %v", err), http.StatusInternalServerError)
		return
	}

	// Create WebSocket handler
	wsHandler := NewWebSocketHandler(h.realTimeManager)

	// Handle connection
	wsHandler.HandleConnection(conn)
}

// Helper function to split comma-separated values
func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
