package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/position"
	"github.com/trading-platform/backend/pkg/utils"
)

// PositionHandler handles HTTP requests related to positions
type PositionHandler struct {
	positionService position.PositionService
}

// NewPositionHandler creates a new PositionHandler
func NewPositionHandler(positionService position.PositionService) *PositionHandler {
	return &PositionHandler{
		positionService: positionService,
	}
}

// CreatePositionFromOrder handles the creation of a new position from an order
func (h *PositionHandler) CreatePositionFromOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create the position
	createdPosition, err := h.positionService.CreatePositionFromOrder(&order)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdPosition)
}

// GetPosition handles the retrieval of a position by ID
func (h *PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	position, err := h.positionService.GetPositionByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, position)
}

// GetPositions handles the retrieval of all positions with optional filtering
func (h *PositionHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filter := models.PositionFilter{}
	
	if userID := r.URL.Query().Get("userId"); userID != "" {
		filter.UserID = userID
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.PositionStatus(status)
	}
	if direction := r.URL.Query().Get("direction"); direction != "" {
		filter.Direction = models.PositionDirection(direction)
	}
	if productType := r.URL.Query().Get("productType"); productType != "" {
		filter.ProductType = models.ProductType(productType)
	}
	if instrumentType := r.URL.Query().Get("instrumentType"); instrumentType != "" {
		filter.InstrumentType = models.InstrumentType(instrumentType)
	}
	if portfolioID := r.URL.Query().Get("portfolioId"); portfolioID != "" {
		filter.PortfolioID = portfolioID
	}
	if strategyID := r.URL.Query().Get("strategyId"); strategyID != "" {
		filter.StrategyID = strategyID
	}
	if orderID := r.URL.Query().Get("orderId"); orderID != "" {
		filter.OrderID = orderID
	}

	// Parse date range if provided
	if fromDate := r.URL.Query().Get("fromDate"); fromDate != "" {
		parsedFromDate, err := time.Parse(time.RFC3339, fromDate)
		if err == nil {
			filter.FromDate = parsedFromDate
		}
	}
	if toDate := r.URL.Query().Get("toDate"); toDate != "" {
		parsedToDate, err := time.Parse(time.RFC3339, toDate)
		if err == nil {
			filter.ToDate = parsedToDate
		}
	}

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get positions with pagination
	positions, total, err := h.positionService.GetPositions(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"positions":   positions,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// UpdatePosition handles the update of an existing position
func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse update data
	var positionUpdate models.Position
	if err := json.NewDecoder(r.Body).Decode(&positionUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set ID
	positionUpdate.ID = id

	// Update the position
	updatedPosition, err := h.positionService.UpdatePosition(&positionUpdate)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPosition)
}

// ClosePosition handles the closing of a position
func (h *PositionHandler) ClosePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse close parameters
	var closeParams struct {
		ExitPrice    float64 `json:"exitPrice"`
		ExitQuantity int     `json:"exitQuantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&closeParams); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Close the position
	closedPosition, err := h.positionService.ClosePosition(id, closeParams.ExitPrice, closeParams.ExitQuantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, closedPosition)
}

// CalculatePnL handles the calculation of P&L for a position
func (h *PositionHandler) CalculatePnL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the position
	position, err := h.positionService.GetPositionByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		return
	}

	// Calculate P&L
	pnl, err := h.positionService.CalculatePnL(position)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"positionId": id,
		"pnl":        pnl,
	})
}

// CalculateGreeks handles the calculation of Greeks for a position
func (h *PositionHandler) CalculateGreeks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the position
	position, err := h.positionService.GetPositionByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		return
	}

	// Calculate Greeks
	greeks, err := h.positionService.CalculateGreeks(position)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"positionId": id,
		"greeks":     greeks,
	})
}

// CalculateExposure handles the calculation of exposure for a user's positions
func (h *PositionHandler) CalculateExposure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Get the user's positions
	filter := models.PositionFilter{
		UserID: userID,
		Status: models.PositionStatusOpen, // Only consider open positions
	}
	positions, _, err := h.positionService.GetPositions(filter, 1, 1000) // Get up to 1000 positions
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate exposure
	exposure, err := h.positionService.CalculateExposure(positions)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"userId":   userID,
		"exposure": exposure,
	})
}

// AggregatePositions handles the aggregation of positions by the specified grouping
func (h *PositionHandler) AggregatePositions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	userID := r.URL.Query().Get("userId")
	groupBy := r.URL.Query().Get("groupBy")
	if groupBy == "" {
		groupBy = "symbol" // Default grouping
	}

	// Get the user's positions
	filter := models.PositionFilter{
		UserID: userID,
	}
	positions, _, err := h.positionService.GetPositions(filter, 1, 1000) // Get up to 1000 positions
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Aggregate positions
	aggregated, err := h.positionService.AggregatePositions(positions, groupBy)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert map to slice for better JSON response
	var result []models.AggregatedPosition
	for _, agg := range aggregated {
		result = append(result, agg)
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

// GetPositionsByUser handles the retrieval of all positions for a specific user
func (h *PositionHandler) GetPositionsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for user
	filter := models.PositionFilter{
		UserID: userID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.PositionStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get positions with pagination
	positions, total, err := h.positionService.GetPositions(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"positions":   positions,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetPositionsByStrategy handles the retrieval of all positions for a specific strategy
func (h *PositionHandler) GetPositionsByStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for strategy
	filter := models.PositionFilter{
		StrategyID: strategyID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.PositionStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get positions with pagination
	positions, total, err := h.positionService.GetPositions(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"positions":   positions,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetPositionsByPortfolio handles the retrieval of all positions for a specific portfolio
func (h *PositionHandler) GetPositionsByPortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID := vars["portfolioId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for portfolio
	filter := models.PositionFilter{
		PortfolioID: portfolioID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.PositionStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get positions with pagination
	positions, total, err := h.positionService.GetPositions(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"positions":   positions,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
