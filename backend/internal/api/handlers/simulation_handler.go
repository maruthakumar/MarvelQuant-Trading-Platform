package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/services/simulation"
)

// SimulationHandler handles API requests related to simulation accounts
type SimulationHandler struct {
	simulationAccountService *simulation.SimulationAccountService
	virtualBalanceService    *simulation.VirtualBalanceService
	simulationOrderService   *simulation.SimulationOrderService
	marketSimulationService  *simulation.MarketSimulationService
	backtestService          *simulation.BacktestService
}

// NewSimulationHandler creates a new instance of SimulationHandler
func NewSimulationHandler() *SimulationHandler {
	return &SimulationHandler{
		simulationAccountService: simulation.NewSimulationAccountService(),
		virtualBalanceService:    simulation.NewVirtualBalanceService(),
		simulationOrderService:   simulation.NewSimulationOrderService(),
		marketSimulationService:  simulation.NewMarketSimulationService(),
		backtestService:          simulation.NewBacktestService(),
	}
}

// CreateSimulationAccount handles the creation of a new simulation account
func (h *SimulationHandler) CreateSimulationAccount(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(string)
	
	// Parse request body
	var accountData models.SimulationAccount
	if err := json.NewDecoder(r.Body).Decode(&accountData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Create simulation account
	account, err := h.simulationAccountService.CreateSimulationAccount(userID, accountData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return created account
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetSimulationAccount handles the retrieval of a simulation account
func (h *SimulationHandler) GetSimulationAccount(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Get simulation account
	account, err := h.simulationAccountService.GetSimulationAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return account
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// GetSimulationAccountsByUser handles the retrieval of all simulation accounts for a user
func (h *SimulationHandler) GetSimulationAccountsByUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(string)
	
	// Get simulation accounts
	accounts, err := h.simulationAccountService.GetSimulationAccountsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return accounts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// AddFunds handles adding funds to a simulation account
func (h *SimulationHandler) AddFunds(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse request body
	var requestData struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Add funds
	transaction, err := h.simulationAccountService.AddFunds(accountID, requestData.Amount, requestData.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return transaction
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// WithdrawFunds handles withdrawing funds from a simulation account
func (h *SimulationHandler) WithdrawFunds(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse request body
	var requestData struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Withdraw funds
	transaction, err := h.simulationAccountService.WithdrawFunds(accountID, requestData.Amount, requestData.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return transaction
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// GetTransactions handles the retrieval of transactions for a simulation account
func (h *SimulationHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse query parameters
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() // Default to now
	}
	
	// Get transactions
	transactions, err := h.simulationAccountService.GetTransactions(accountID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return transactions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// ResetAccount handles resetting a simulation account to its initial state
func (h *SimulationHandler) ResetAccount(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Reset account
	err := h.simulationAccountService.ResetAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account reset successfully"})
}

// GetAccountPerformance handles the retrieval of performance metrics for a simulation account
func (h *SimulationHandler) GetAccountPerformance(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse query parameters
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() // Default to now
	}
	
	// Get performance metrics
	metrics, err := h.simulationAccountService.GetAccountPerformance(accountID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return metrics
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// CreateOrder handles the creation of a new simulation order
func (h *SimulationHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(string)
	
	// Parse request body
	var orderData models.SimulationOrder
	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set user ID
	orderData.UserID = userID
	
	// Create order
	order, err := h.simulationOrderService.CreateOrder(accountID, orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder handles the retrieval of a simulation order
func (h *SimulationHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	
	// Get order
	order, err := h.simulationOrderService.GetOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrdersByAccount handles the retrieval of all simulation orders for an account
func (h *SimulationHandler) GetOrdersByAccount(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Get orders
	orders, err := h.simulationOrderService.GetOrdersByAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// CancelOrder handles cancelling a simulation order
func (h *SimulationHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	
	// Cancel order
	order, err := h.simulationOrderService.CancelOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return cancelled order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// ModifyOrder handles modifying a simulation order
func (h *SimulationHandler) ModifyOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	
	// Parse request body
	var orderData models.SimulationOrder
	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Modify order
	order, err := h.simulationOrderService.ModifyOrder(orderID, orderData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return modified order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrderHistory handles the retrieval of order history for a simulation account
func (h *SimulationHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse query parameters
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	symbol := r.URL.Query().Get("symbol")
	
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() // Default to now
	}
	
	// Get order history
	orders, err := h.simulationOrderService.GetOrderHistory(accountID, startDate, endDate, symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrderStatistics handles the retrieval of order statistics for a simulation account
func (h *SimulationHandler) GetOrderStatistics(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse query parameters
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() // Default to now
	}
	
	// Get order statistics
	statistics, err := h.simulationOrderService.GetOrderStatistics(accountID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return statistics
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}

// GetCurrentMarketPrice handles the retrieval of the current market price for a symbol
func (h *SimulationHandler) GetCurrentMarketPrice(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	
	// Get current market price
	marketData, err := h.marketSimulationService.GetCurrentMarketPrice(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return market data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketData)
}

// GetHistoricalMarketData handles the retrieval of historical market data for a symbol
func (h *SimulationHandler) GetHistoricalMarketData(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	
	// Parse query parameters
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	timeframe := r.URL.Query().Get("timeframe")
	
	var startDate, endDate time.Time
	var err error
	
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}
	
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() // Default to now
	}
	
	if timeframe == "" {
		timeframe = "1d" // Default to daily
	}
	
	// Get historical market data
	marketData, err := h.marketSimulationService.GetHistoricalMarketData(symbol, startDate, endDate, timeframe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return market data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketData)
}

// SimulateMarketMovement handles simulating market movement for a symbol
func (h *SimulationHandler) SimulateMarketMovement(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	
	// Parse request body
	var requestData struct {
		Timeframe string        `json:"timeframe"`
		Duration  time.Duration `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Simulate market movement
	marketData, err := h.marketSimulationService.SimulateMarketMovement(symbol, requestData.Timeframe, requestData.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return market data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketData)
}

// GetMarketDepth handles the retrieval of market depth for a symbol
func (h *SimulationHandler) GetMarketDepth(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	
	// Parse query parameters
	levelsStr := r.URL.Query().Get("levels")
	levels := 5 // Default to 5 levels
	
	if levelsStr != "" {
		var err error
		levels, err = strconv.Atoi(levelsStr)
		if err != nil {
			http.Error(w, "Invalid levels parameter", http.StatusBadRequest)
			return
		}
	}
	
	// Get market depth
	marketDepth, err := h.marketSimulationService.GetMarketDepth(symbol, levels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return market depth
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketDepth)
}

// SimulateMarketCondition handles simulating a specific market condition
func (h *SimulationHandler) SimulateMarketCondition(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	
	// Parse request body
	var requestData struct {
		Condition string        `json:"condition"`
		Duration  time.Duration `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Simulate market condition
	marketData, err := h.marketSimulationService.SimulateMarketCondition(symbol, requestData.Condition, requestData.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return market data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marketData)
}

// CreateBacktestSession handles the creation of a new backtest session
func (h *SimulationHandler) CreateBacktestSession(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Parse request body
	var sessionData models.BacktestSession
	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Create backtest session
	session, err := h.backtestService.CreateBacktestSession(accountID, sessionData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return created session
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// GetBacktestSession handles the retrieval of a backtest session
func (h *SimulationHandler) GetBacktestSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Get backtest session
	session, err := h.backtestService.GetBacktestSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetBacktestSessionsByAccount handles the retrieval of all backtest sessions for an account
func (h *SimulationHandler) GetBacktestSessionsByAccount(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from URL
	vars := mux.Vars(r)
	accountID := vars["accountID"]
	
	// Get backtest sessions
	sessions, err := h.backtestService.GetBacktestSessionsByAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return sessions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// RunBacktest handles running a backtest session
func (h *SimulationHandler) RunBacktest(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Run backtest
	err := h.backtestService.RunBacktest(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Backtest started successfully"})
}

// StopBacktest handles stopping a running backtest session
func (h *SimulationHandler) StopBacktest(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Stop backtest
	err := h.backtestService.StopBacktest(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Backtest stopped successfully"})
}

// GetBacktestResults handles the retrieval of backtest results
func (h *SimulationHandler) GetBacktestResults(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Get backtest results
	results, err := h.backtestService.GetBacktestResults(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetBacktestPerformanceMetrics handles the retrieval of backtest performance metrics
func (h *SimulationHandler) GetBacktestPerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Get backtest performance metrics
	metrics, err := h.backtestService.GetBacktestPerformanceMetrics(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return metrics
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetBacktestTrades handles the retrieval of backtest trades
func (h *SimulationHandler) GetBacktestTrades(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Get backtest trades
	trades, err := h.backtestService.GetBacktestTrades(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return trades
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// CompareBacktestSessions handles comparing multiple backtest sessions
func (h *SimulationHandler) CompareBacktestSessions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData struct {
		SessionIDs []string `json:"sessionIDs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Compare backtest sessions
	comparison, err := h.backtestService.CompareBacktestSessions(requestData.SessionIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return comparison
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparison)
}

// OptimizeStrategy handles optimizing a strategy based on historical data
func (h *SimulationHandler) OptimizeStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract strategy ID from URL
	vars := mux.Vars(r)
	strategyID := vars["strategyID"]
	
	// Parse request body
	var requestData struct {
		ParameterRanges   map[string]map[string]interface{} `json:"parameterRanges"`
		OptimizationMetric string                           `json:"optimizationMetric"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Optimize strategy
	results, err := h.backtestService.OptimizeStrategy(strategyID, requestData.ParameterRanges, requestData.OptimizationMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ExportBacktestResults handles exporting backtest results to a specified format
func (h *SimulationHandler) ExportBacktestResults(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	
	// Parse query parameters
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv" // Default to CSV
	}
	
	// Export backtest results
	filePath, err := h.backtestService.ExportBacktestResults(sessionID, format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return file path
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"filePath": filePath})
}
