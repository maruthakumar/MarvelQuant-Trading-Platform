package handlers

import (
	"encoding/json"
	"net/http"

	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/services/user"
	"trading_platform/backend/internal/utils"
)

// EnvironmentHandler handles environment-related requests
type EnvironmentHandler struct {
	environmentService *user.EnvironmentService
}

// NewEnvironmentHandler creates a new environment handler
func NewEnvironmentHandler(environmentService *user.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{
		environmentService: environmentService,
	}
}

// SwitchEnvironment handles requests to switch environment
func (h *EnvironmentHandler) SwitchEnvironment(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Environment string `json:"environment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate environment
	environment := models.Environment(req.Environment)
	if environment != models.EnvironmentLive && environment != models.EnvironmentSIM {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid environment")
		return
	}

	// Switch environment
	token, err := h.environmentService.SwitchEnvironment(r.Context(), environment)
	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	// Respond with new token
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token":       token,
		"environment": string(environment),
	})
}

// GetEnvironmentStatus handles requests to get environment status
func (h *EnvironmentHandler) GetEnvironmentStatus(w http.ResponseWriter, r *http.Request) {
	// Get environment status
	status := h.environmentService.GetEnvironmentStatus(r.Context())

	// Respond with status
	utils.RespondWithJSON(w, http.StatusOK, status)
}
