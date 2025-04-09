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

// StrategyHandler handles strategy-related API endpoints
type StrategyHandler struct {
	strategyRepo *database.StrategyRepository
}

// NewStrategyHandler creates a new StrategyHandler
func NewStrategyHandler(strategyRepo *database.StrategyRepository) *StrategyHandler {
	return &StrategyHandler{
		strategyRepo: strategyRepo,
	}
}

// CreateStrategy handles the creation of a new strategy
func (h *StrategyHandler) CreateStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var strategy models.Strategy
	if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID
	strategy.UserID = userID

	// Validate strategy
	if err := strategy.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create strategy
	id, err := h.strategyRepo.Create(&strategy)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating strategy")
		return
	}

	// Set ID in response
	strategy.ID = id

	utils.RespondWithJSON(w, http.StatusCreated, strategy)
}

// GetStrategy handles retrieving a strategy by ID
func (h *StrategyHandler) GetStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get strategy ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get strategy
	strategy, err := h.strategyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
		}
		return
	}

	// Check if user has access to this strategy
	if strategy.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategy)
}

// UpdateStrategy handles updating a strategy
func (h *StrategyHandler) UpdateStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get strategy ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing strategy
	existingStrategy, err := h.strategyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
		}
		return
	}

	// Check if user has access to this strategy
	if existingStrategy.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var updatedStrategy models.Strategy
	if err := json.NewDecoder(r.Body).Decode(&updatedStrategy); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and user ID
	updatedStrategy.ID = id
	updatedStrategy.UserID = userID
	updatedStrategy.CreatedAt = existingStrategy.CreatedAt

	// Validate strategy
	if err := updatedStrategy.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update strategy
	if err := h.strategyRepo.Update(&updatedStrategy); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating strategy")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedStrategy)
}

// DeleteStrategy handles deleting a strategy
func (h *StrategyHandler) DeleteStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get strategy ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing strategy
	existingStrategy, err := h.strategyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
		}
		return
	}

	// Check if user has access to this strategy
	if existingStrategy.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete strategy
	if err := h.strategyRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting strategy")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy deleted successfully"})
}

// GetStrategies handles retrieving strategies with filtering and pagination
func (h *StrategyHandler) GetStrategies(w http.ResponseWriter, r *http.Request) {
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
	filter := models.StrategyFilter{
		UserID: userID,
	}

	// Add optional filters
	if name := query.Get("name"); name != "" {
		filter.Name = name
	}

	if strategyType := query.Get("type"); strategyType != "" {
		filter.Type = strategyType
	}

	if tag := query.Get("tag"); tag != "" {
		filter.Tag = tag
	}

	if activeStr := query.Get("active"); activeStr != "" {
		active := activeStr == "true"
		filter.Active = &active
	}

	if symbol := query.Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	if productType := query.Get("productType"); productType != "" {
		filter.ProductType = productType
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

	// Get strategies
	strategies, total, err := h.strategyRepo.Find(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategies")
		return
	}

	// Build response with pagination
	response := map[string]interface{}{
		"data":       strategies,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// ActivateStrategy handles activating a strategy
func (h *StrategyHandler) ActivateStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get strategy ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing strategy
	existingStrategy, err := h.strategyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
		}
		return
	}

	// Check if user has access to this strategy
	if existingStrategy.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if strategy is already active
	if existingStrategy.Active {
		utils.RespondWithError(w, http.StatusBadRequest, "Strategy is already active")
		return
	}

	// Update strategy
	existingStrategy.Active = true
	existingStrategy.UpdatedAt = time.Now()

	if err := h.strategyRepo.Update(existingStrategy); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error activating strategy")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingStrategy)
}

// DeactivateStrategy handles deactivating a strategy
func (h *StrategyHandler) DeactivateStrategy(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get strategy ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing strategy
	existingStrategy, err := h.strategyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
		}
		return
	}

	// Check if user has access to this strategy
	if existingStrategy.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if strategy is already inactive
	if !existingStrategy.Active {
		utils.RespondWithError(w, http.StatusBadRequest, "Strategy is already inactive")
		return
	}

	// Update strategy
	existingStrategy.Active = false
	existingStrategy.UpdatedAt = time.Now()

	if err := h.strategyRepo.Update(existingStrategy); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deactivating strategy")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingStrategy)
}

// RegisterStrategyRoutes registers strategy-related routes
func RegisterStrategyRoutes(router *mux.Router, strategyRepo *database.StrategyRepository, authMiddleware func(http.Handler) http.Handler) {
	handler := NewStrategyHandler(strategyRepo)
	
	// Apply auth middleware to all routes
	strategyRouter := router.PathPrefix("/strategies").Subrouter()
	strategyRouter.Use(authMiddleware)

	// Register routes
	strategyRouter.HandleFunc("", handler.CreateStrategy).Methods("POST")
	strategyRouter.HandleFunc("", handler.GetStrategies).Methods("GET")
	strategyRouter.HandleFunc("/{id}", handler.GetStrategy).Methods("GET")
	strategyRouter.HandleFunc("/{id}", handler.UpdateStrategy).Methods("PUT")
	strategyRouter.HandleFunc("/{id}", handler.DeleteStrategy).Methods("DELETE")
	strategyRouter.HandleFunc("/{id}/activate", handler.ActivateStrategy).Methods("POST")
	strategyRouter.HandleFunc("/{id}/deactivate", handler.DeactivateStrategy).Methods("POST")
}
