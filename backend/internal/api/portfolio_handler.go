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

// PortfolioHandler handles portfolio-related API endpoints
type PortfolioHandler struct {
	portfolioRepo *database.PortfolioRepository
	strategyRepo  *database.StrategyRepository
}

// NewPortfolioHandler creates a new PortfolioHandler
func NewPortfolioHandler(portfolioRepo *database.PortfolioRepository, strategyRepo *database.StrategyRepository) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioRepo: portfolioRepo,
		strategyRepo:  strategyRepo,
	}
}

// CreatePortfolio handles the creation of a new portfolio
func (h *PortfolioHandler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var portfolio models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&portfolio); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID
	portfolio.UserID = userID

	// Validate portfolio
	if err := portfolio.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// If strategy ID is provided, check if it exists and belongs to the user
	if portfolio.StrategyID != "" {
		strategy, err := h.strategyRepo.GetByID(portfolio.StrategyID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				utils.RespondWithError(w, http.StatusBadRequest, "Strategy not found")
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
			}
			return
		}

		if strategy.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "Access denied to strategy")
			return
		}
	}

	// Create portfolio
	id, err := h.portfolioRepo.Create(&portfolio)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating portfolio")
		return
	}

	// Set ID in response
	portfolio.ID = id

	utils.RespondWithJSON(w, http.StatusCreated, portfolio)
}

// GetPortfolio handles retrieving a portfolio by ID
func (h *PortfolioHandler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get portfolio
	portfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if portfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, portfolio)
}

// UpdatePortfolio handles updating a portfolio
func (h *PortfolioHandler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var updatedPortfolio models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&updatedPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and user ID
	updatedPortfolio.ID = id
	updatedPortfolio.UserID = userID
	updatedPortfolio.CreatedAt = existingPortfolio.CreatedAt

	// Validate portfolio
	if err := updatedPortfolio.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// If strategy ID is provided, check if it exists and belongs to the user
	if updatedPortfolio.StrategyID != "" {
		strategy, err := h.strategyRepo.GetByID(updatedPortfolio.StrategyID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				utils.RespondWithError(w, http.StatusBadRequest, "Strategy not found")
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving strategy")
			}
			return
		}

		if strategy.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "Access denied to strategy")
			return
		}
	}

	// Update portfolio
	if err := h.portfolioRepo.Update(&updatedPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPortfolio)
}

// DeletePortfolio handles deleting a portfolio
func (h *PortfolioHandler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete portfolio
	if err := h.portfolioRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Portfolio deleted successfully"})
}

// GetPortfolios handles retrieving portfolios with filtering and pagination
func (h *PortfolioHandler) GetPortfolios(w http.ResponseWriter, r *http.Request) {
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
	filter := models.PortfolioFilter{
		UserID: userID,
	}

	// Add optional filters
	if name := query.Get("name"); name != "" {
		filter.Name = name
	}

	if strategyID := query.Get("strategyId"); strategyID != "" {
		filter.StrategyID = strategyID
	}

	if status := query.Get("status"); status != "" {
		filter.Status = status
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

	// Get portfolios
	portfolios, total, err := h.portfolioRepo.Find(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolios")
		return
	}

	// Build response with pagination
	response := map[string]interface{}{
		"data":       portfolios,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// ActivatePortfolio handles activating a portfolio
func (h *PortfolioHandler) ActivatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if portfolio is already active
	if existingPortfolio.Status == models.PortfolioStatusActive {
		utils.RespondWithError(w, http.StatusBadRequest, "Portfolio is already active")
		return
	}

	// Update portfolio
	existingPortfolio.Status = models.PortfolioStatusActive
	existingPortfolio.UpdatedAt = time.Now()

	if err := h.portfolioRepo.Update(existingPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error activating portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPortfolio)
}

// DeactivatePortfolio handles deactivating a portfolio
func (h *PortfolioHandler) DeactivatePortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if portfolio is already inactive
	if existingPortfolio.Status != models.PortfolioStatusActive {
		utils.RespondWithError(w, http.StatusBadRequest, "Portfolio is not active")
		return
	}

	// Update portfolio
	existingPortfolio.Status = models.PortfolioStatusInactive
	existingPortfolio.UpdatedAt = time.Now()

	if err := h.portfolioRepo.Update(existingPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deactivating portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPortfolio)
}

// AddLegToPortfolio handles adding a leg to a portfolio
func (h *PortfolioHandler) AddLegToPortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var leg models.Leg
	if err := json.NewDecoder(r.Body).Decode(&leg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set portfolio ID
	leg.PortfolioID = id

	// Validate leg
	if err := leg.CalculateQuantity(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Add leg to portfolio
	if existingPortfolio.Legs == nil {
		existingPortfolio.Legs = make([]models.Leg, 0)
	}
	
	// Set leg ID
	leg.ID = len(existingPortfolio.Legs) + 1
	leg.CreatedAt = time.Now()
	leg.UpdatedAt = time.Now()
	
	existingPortfolio.Legs = append(existingPortfolio.Legs, leg)
	existingPortfolio.UpdatedAt = time.Now()

	// Update portfolio
	if err := h.portfolioRepo.Update(existingPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error adding leg to portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPortfolio)
}

// UpdateLegInPortfolio handles updating a leg in a portfolio
func (h *PortfolioHandler) UpdateLegInPortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID and leg ID from URL
	vars := mux.Vars(r)
	portfolioID := vars["id"]
	legIDStr := vars["legId"]
	
	legID, err := strconv.Atoi(legIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid leg ID")
		return
	}

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Find leg in portfolio
	legIndex := -1
	for i, leg := range existingPortfolio.Legs {
		if leg.ID == legID {
			legIndex = i
			break
		}
	}

	if legIndex == -1 {
		utils.RespondWithError(w, http.StatusNotFound, "Leg not found in portfolio")
		return
	}

	// Parse request body
	var updatedLeg models.Leg
	if err := json.NewDecoder(r.Body).Decode(&updatedLeg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set portfolio ID and leg ID
	updatedLeg.PortfolioID = portfolioID
	updatedLeg.ID = legID
	updatedLeg.CreatedAt = existingPortfolio.Legs[legIndex].CreatedAt
	updatedLeg.UpdatedAt = time.Now()

	// Validate leg
	if err := updatedLeg.CalculateQuantity(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update leg in portfolio
	existingPortfolio.Legs[legIndex] = updatedLeg
	existingPortfolio.UpdatedAt = time.Now()

	// Update portfolio
	if err := h.portfolioRepo.Update(existingPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating leg in portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPortfolio)
}

// RemoveLegFromPortfolio handles removing a leg from a portfolio
func (h *PortfolioHandler) RemoveLegFromPortfolio(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get portfolio ID and leg ID from URL
	vars := mux.Vars(r)
	portfolioID := vars["id"]
	legIDStr := vars["legId"]
	
	legID, err := strconv.Atoi(legIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid leg ID")
		return
	}

	// Get existing portfolio
	existingPortfolio, err := h.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Portfolio not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving portfolio")
		}
		return
	}

	// Check if user has access to this portfolio
	if existingPortfolio.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Find leg in portfolio
	legIndex := -1
	for i, leg := range existingPortfolio.Legs {
		if leg.ID == legID {
			legIndex = i
			break
		}
	}

	if legIndex == -1 {
		utils.RespondWithError(w, http.StatusNotFound, "Leg not found in portfolio")
		return
	}

	// Remove leg from portfolio
	existingPortfolio.Legs = append(existingPortfolio.Legs[:legIndex], existingPortfolio.Legs[legIndex+1:]...)
	existingPortfolio.UpdatedAt = time.Now()

	// Update portfolio
	if err := h.portfolioRepo.Update(existingPortfolio); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error removing leg from portfolio")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingPortfolio)
}

// RegisterPortfolioRoutes registers portfolio-related routes
func RegisterPortfolioRoutes(
	router *mux.Router, 
	portfolioRepo *database.PortfolioRepository, 
	strategyRepo *database.StrategyRepository,
	authMiddleware func(http.Handler) http.Handler,
) {
	handler := NewPortfolioHandler(portfolioRepo, strategyRepo)
	
	// Apply auth middleware to all routes
	portfolioRouter := router.PathPrefix("/portfolios").Subrouter()
	portfolioRouter.Use(authMiddleware)

	// Register routes
	portfolioRouter.HandleFunc("", handler.CreatePortfolio).Methods("POST")
	portfolioRouter.HandleFunc("", handler.GetPortfolios).Methods("GET")
	portfolioRouter.HandleFunc("/{id}", handler.GetPortfolio).Methods("GET")
	portfolioRouter.HandleFunc("/{id}", handler.UpdatePortfolio).Methods("PUT")
	portfolioRouter.HandleFunc("/{id}", handler.DeletePortfolio).Methods("DELETE")
	portfolioRouter.HandleFunc("/{id}/activate", handler.ActivatePortfolio).Methods("POST")
	portfolioRouter.HandleFunc("/{id}/deactivate", handler.DeactivatePortfolio).Methods("POST")
	portfolioRouter.HandleFunc("/{id}/legs", handler.AddLegToPortfolio).Methods("POST")
	portfolioRouter.HandleFunc("/{id}/legs/{legId}", handler.UpdateLegInPortfolio).Methods("PUT")
	portfolioRouter.HandleFunc("/{id}/legs/{legId}", handler.RemoveLegFromPortfolio).Methods("DELETE")
}
