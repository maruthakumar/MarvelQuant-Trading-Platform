package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"trading_platform/backend/internal/database"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/utils"
)

// PositionHandler handles position-related API endpoints
type PositionHandler struct {
	positionRepo *database.PositionRepository
}

// NewPositionHandler creates a new PositionHandler
func NewPositionHandler(positionRepo *database.PositionRepository) *PositionHandler {
	return &PositionHandler{
		positionRepo: positionRepo,
	}
}

// CreatePosition handles the creation of a new position
func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var position models.Position
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID
	position.UserID = userID

	// Validate position
	if err := position.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create position
	id, err := h.positionRepo.Create(&position)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating position")
		return
	}

	// Set ID in response
	position.ID = id

	utils.RespondWithJSON(w, http.StatusCreated, position)
}

// GetPosition handles retrieving a position by ID
func (h *PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get position ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get position
	position, err := h.positionRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving position")
		}
		return
	}

	// Check if user has access to this position
	if position.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, position)
}

// UpdatePosition handles updating a position
func (h *PositionHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get position ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing position
	existingPosition, err := h.positionRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving position")
		}
		return
	}

	// Check if user has access to this position
	if existingPosition.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var updatedPosition models.Position
	if err := json.NewDecoder(r.Body).Decode(&updatedPosition); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and user ID
	updatedPosition.ID = id
	updatedPosition.UserID = userID
	updatedPosition.CreatedAt = existingPosition.CreatedAt
	updatedPosition.EntryTime = existingPosition.EntryTime

	// Validate position
	if err := updatedPosition.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update position
	if err := h.positionRepo.Update(&updatedPosition); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating position")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPosition)
}

// DeletePosition handles deleting a position
func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get position ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing position
	existingPosition, err := h.positionRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving position")
		}
		return
	}

	// Check if user has access to this position
	if existingPosition.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete position
	if err := h.positionRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting position")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Position deleted successfully"})
}

// GetPositions handles retrieving positions with filtering and pagination
func (h *PositionHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	
	// Parse pagination parameters
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// Build filter
	filter := models.PositionFilter{
		UserID: userID,
	}

	// Add optional filters
	if symbol := query.Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	if status := query.Get("status"); status != "" {
		filter.Status = status
	}

	if direction := query.Get("direction"); direction != "" {
		filter.Direction = direction
	}

	if productType := query.Get("productType"); productType != "" {
		filter.ProductType = productType
	}

	if instrumentType := query.Get("instrumentType"); instrumentType != "" {
		filter.InstrumentType = instrumentType
	}

	if portfolioID := query.Get("portfolioId"); portfolioID != "" {
		filter.PortfolioID = portfolioID
	}

	if strategyID := query.Get("strategyId"); strategyID != "" {
		filter.StrategyID = strategyID
	}

	// Parse date range
	if fromDate := query.Get("fromDate"); fromDate != "" {
		parsedFromDate, err := time.Parse(time.RFC3339, fromDate)
		if err == nil {
			filter.FromDate = parsedFromDate
		}
	}

	if toDate := query.Get("toDate"); toDate != "" {
		parsedToDate, err := time.Parse(time.RFC3339, toDate)
		if err == nil {
			filter.ToDate = parsedToDate
		}
	}

	// Parse tags
	if tags := query.Get("tags"); tags != "" {
		filter.Tags = utils.SplitAndTrim(tags, ",")
	}

	// Get positions
	positions, total, err := h.positionRepo.Find(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving positions")
		return
	}

	// Build response with pagination
	response := map[string]interface{}{
		"data":       positions,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// ClosePosition handles closing a position
func (h *PositionHandler) ClosePosition(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get position ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing position
	existingPosition, err := h.positionRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Position not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving position")
		}
		return
	}

	// Check if user has access to this position
	if existingPosition.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if position can be closed
	if existingPosition.Status != models.PositionStatusOpen {
		utils.RespondWithError(w, http.StatusBadRequest, "Position cannot be closed")
		return
	}

	// Parse request body for exit price
	var closeRequest struct {
		ExitPrice float64 `json:"exitPrice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&closeRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Update position
	existingPosition.Status = models.PositionStatusClosed
	existingPosition.ExitPrice = closeRequest.ExitPrice
	existingPosition.ExitTime = time.Now()
	existingPosition.UpdatedAt = time.Now()

	// Calculate realized P&L
	existingPosition.CalculateRealizedPnL()

	// Update position
	if err := h.positionRepo.Update(existingPosition); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error closing position")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPosition)
}

// RegisterPositionRoutes registers position-related routes
func RegisterPositionRoutes(router *mux.Router, positionRepo *database.PositionRepository, authMiddleware func(http.Handler) http.Handler) {
	handler := NewPositionHandler(positionRepo)
	
	// Apply auth middleware to all routes
	positionRouter := router.PathPrefix("/positions").Subrouter()
	positionRouter.Use(authMiddleware)

	// Register routes
	positionRouter.HandleFunc("", handler.CreatePosition).Methods("POST")
	positionRouter.HandleFunc("", handler.GetPositions).Methods("GET")
	positionRouter.HandleFunc("/{id}", handler.GetPosition).Methods("GET")
	positionRouter.HandleFunc("/{id}", handler.UpdatePosition).Methods("PUT")
	positionRouter.HandleFunc("/{id}", handler.DeletePosition).Methods("DELETE")
	positionRouter.HandleFunc("/{id}/close", handler.ClosePosition).Methods("POST")
}
