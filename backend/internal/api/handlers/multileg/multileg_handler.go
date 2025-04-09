package multileg

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/multileg"
	"github.com/trading-platform/backend/pkg/utils"
)

// MultilegHandler handles HTTP requests related to multileg strategies
type MultilegHandler struct {
	multilegService multileg.MultilegService
}

// NewMultilegHandler creates a new MultilegHandler
func NewMultilegHandler(multilegService multileg.MultilegService) *MultilegHandler {
	return &MultilegHandler{
		multilegService: multilegService,
	}
}

// CreateMultilegStrategy handles the creation of a new multileg strategy
func (h *MultilegHandler) CreateMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	var strategy models.MultilegStrategy
	if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	createdStrategy, err := h.multilegService.CreateMultilegStrategy(&strategy)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdStrategy)
}

// GetMultilegStrategy handles the retrieval of a multileg strategy by ID
func (h *MultilegHandler) GetMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	strategy, err := h.multilegService.GetMultilegStrategyByID(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategy)
}

// GetUserMultilegStrategies handles the retrieval of all multileg strategies for a user
func (h *MultilegHandler) GetUserMultilegStrategies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	strategies, err := h.multilegService.GetMultilegStrategiesByUser(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategies)
}

// GetPortfolioMultilegStrategies handles the retrieval of all multileg strategies for a portfolio
func (h *MultilegHandler) GetPortfolioMultilegStrategies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID := vars["portfolioId"]

	strategies, err := h.multilegService.GetMultilegStrategiesByPortfolio(portfolioID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategies)
}

// UpdateMultilegStrategy handles the update of a multileg strategy
func (h *MultilegHandler) UpdateMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var strategy models.MultilegStrategy
	if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Ensure the strategy ID matches the URL parameter
	if strategy.ID != strategyID {
		utils.RespondWithError(w, http.StatusBadRequest, "Strategy ID mismatch")
		return
	}

	updatedStrategy, err := h.multilegService.UpdateMultilegStrategy(&strategy)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedStrategy)
}

// DeleteMultilegStrategy handles the deletion of a multileg strategy
func (h *MultilegHandler) DeleteMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.multilegService.DeleteMultilegStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy deleted successfully"})
}

// AddLeg handles the addition of a leg to a multileg strategy
func (h *MultilegHandler) AddLeg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var leg models.Leg
	if err := json.NewDecoder(r.Body).Decode(&leg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	createdLeg, err := h.multilegService.AddLeg(strategyID, &leg)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdLeg)
}

// UpdateLeg handles the update of a leg in a multileg strategy
func (h *MultilegHandler) UpdateLeg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]
	legID := vars["legId"]

	var leg models.Leg
	if err := json.NewDecoder(r.Body).Decode(&leg); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Ensure the leg ID matches the URL parameter
	if leg.ID != legID {
		utils.RespondWithError(w, http.StatusBadRequest, "Leg ID mismatch")
		return
	}

	updatedLeg, err := h.multilegService.UpdateLeg(strategyID, &leg)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedLeg)
}

// RemoveLeg handles the removal of a leg from a multileg strategy
func (h *MultilegHandler) RemoveLeg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]
	legID := vars["legId"]

	err := h.multilegService.RemoveLeg(strategyID, legID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Leg removed successfully"})
}

// GetLegs handles the retrieval of all legs for a multileg strategy
func (h *MultilegHandler) GetLegs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	legs, err := h.multilegService.GetLegsByStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, legs)
}

// ExecuteMultilegStrategy handles the execution of a multileg strategy
func (h *MultilegHandler) ExecuteMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.multilegService.ExecuteMultilegStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy execution started"})
}

// PauseMultilegStrategy handles the pausing of a multileg strategy
func (h *MultilegHandler) PauseMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.multilegService.PauseMultilegStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy paused successfully"})
}

// ResumeMultilegStrategy handles the resumption of a multileg strategy
func (h *MultilegHandler) ResumeMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.multilegService.ResumeMultilegStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy resumed successfully"})
}

// CancelMultilegStrategy handles the cancellation of a multileg strategy
func (h *MultilegHandler) CancelMultilegStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.multilegService.CancelMultilegStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy canceled successfully"})
}

// GetMultilegStrategyStatus handles the retrieval of a multileg strategy's status
func (h *MultilegHandler) GetMultilegStrategyStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	status, err := h.multilegService.GetMultilegStrategyStatus(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": status})
}

// GetMultilegStrategyPerformance handles the retrieval of a multileg strategy's performance
func (h *MultilegHandler) GetMultilegStrategyPerformance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	performance, err := h.multilegService.GetMultilegStrategyPerformance(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, performance)
}
