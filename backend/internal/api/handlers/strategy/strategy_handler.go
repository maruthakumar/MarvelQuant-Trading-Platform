package strategy

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/strategy"
	"github.com/trading-platform/backend/pkg/utils"
)

// StrategyHandler handles HTTP requests related to strategies
type StrategyHandler struct {
	strategyService strategy.StrategyService
}

// NewStrategyHandler creates a new StrategyHandler
func NewStrategyHandler(strategyService strategy.StrategyService) *StrategyHandler {
	return &StrategyHandler{
		strategyService: strategyService,
	}
}

// CreateStrategy handles the creation of a new strategy
func (h *StrategyHandler) CreateStrategy(w http.ResponseWriter, r *http.Request) {
	var strategy models.Strategy
	if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	createdStrategy, err := h.strategyService.CreateStrategy(&strategy)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdStrategy)
}

// GetStrategy handles the retrieval of a strategy by ID
func (h *StrategyHandler) GetStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	strategy, err := h.strategyService.GetStrategyByID(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Strategy not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategy)
}

// GetUserStrategies handles the retrieval of all strategies for a user
func (h *StrategyHandler) GetUserStrategies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	strategies, err := h.strategyService.GetStrategiesByUser(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategies)
}

// UpdateStrategy handles the update of a strategy
func (h *StrategyHandler) UpdateStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var strategy models.Strategy
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

	updatedStrategy, err := h.strategyService.UpdateStrategy(&strategy)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedStrategy)
}

// DeleteStrategy handles the deletion of a strategy
func (h *StrategyHandler) DeleteStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.DeleteStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy deleted successfully"})
}

// ExecuteStrategy handles the execution of a strategy
func (h *StrategyHandler) ExecuteStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.ExecuteStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy execution started"})
}

// PauseStrategy handles the pausing of a strategy
func (h *StrategyHandler) PauseStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.PauseStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy paused successfully"})
}

// ResumeStrategy handles the resumption of a strategy
func (h *StrategyHandler) ResumeStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.ResumeStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy resumed successfully"})
}

// StopStrategy handles the stopping of a strategy
func (h *StrategyHandler) StopStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.StopStrategy(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy stopped successfully"})
}

// GetStrategyStatus handles the retrieval of a strategy's status
func (h *StrategyHandler) GetStrategyStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	status, err := h.strategyService.GetStrategyStatus(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": string(status)})
}

// GetStrategyPerformance handles the retrieval of a strategy's performance
func (h *StrategyHandler) GetStrategyPerformance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	performance, err := h.strategyService.GetStrategyPerformance(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, performance)
}

// ScheduleStrategy handles the scheduling of a strategy
func (h *StrategyHandler) ScheduleStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var schedule models.StrategySchedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the strategy ID
	schedule.StrategyID = strategyID

	err := h.strategyService.ScheduleStrategy(strategyID, &schedule)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy scheduled successfully"})
}

// GetStrategySchedule handles the retrieval of a strategy's schedule
func (h *StrategyHandler) GetStrategySchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	schedule, err := h.strategyService.GetStrategySchedule(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, schedule)
}

// UpdateStrategySchedule handles the update of a strategy's schedule
func (h *StrategyHandler) UpdateStrategySchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var schedule models.StrategySchedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the strategy ID
	schedule.StrategyID = strategyID

	err := h.strategyService.UpdateStrategySchedule(strategyID, &schedule)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy schedule updated successfully"})
}

// DeleteStrategySchedule handles the deletion of a strategy's schedule
func (h *StrategyHandler) DeleteStrategySchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	err := h.strategyService.DeleteStrategySchedule(strategyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Strategy schedule deleted successfully"})
}

// AddStrategyTag handles the addition of a tag to a strategy
func (h *StrategyHandler) AddStrategyTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	var tagRequest struct {
		Tag string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&tagRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := h.strategyService.AddStrategyTag(strategyID, tagRequest.Tag)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Tag added successfully"})
}

// RemoveStrategyTag handles the removal of a tag from a strategy
func (h *StrategyHandler) RemoveStrategyTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]
	tag := vars["tag"]

	err := h.strategyService.RemoveStrategyTag(strategyID, tag)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Tag removed successfully"})
}

// GetStrategiesByTag handles the retrieval of all strategies with a specific tag
func (h *StrategyHandler) GetStrategiesByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]

	strategies, err := h.strategyService.GetStrategiesByTag(tag)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, strategies)
}
