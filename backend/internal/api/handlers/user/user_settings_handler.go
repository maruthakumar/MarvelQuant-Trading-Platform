package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services/user"
	"github.com/trading-platform/backend/pkg/utils"
)

// UserSettingsHandler handles HTTP requests related to user settings
type UserSettingsHandler struct {
	userService user.UserService
}

// NewUserSettingsHandler creates a new UserSettingsHandler
func NewUserSettingsHandler(userService user.UserService) *UserSettingsHandler {
	return &UserSettingsHandler{
		userService: userService,
	}
}

// GetUserSettings handles the retrieval of user settings
func (h *UserSettingsHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	settings, err := h.userService.GetUserSettings(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User settings not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, settings)
}

// UpdateUserSettings handles the update of user settings
func (h *UserSettingsHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var settings models.UserSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	settings.UserID = userID
	settings.UpdatedAt = time.Now()

	updatedSettings, err := h.userService.UpdateUserSettings(userID, &settings)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedSettings)
}

// GetUserPreferences handles the retrieval of user preferences
func (h *UserSettingsHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	preferences, err := h.userService.GetUserPreferences(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User preferences not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, preferences)
}

// UpdateUserPreferences handles the update of user preferences
func (h *UserSettingsHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var preferences models.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	preferences.UserID = userID
	preferences.UpdatedAt = time.Now()

	updatedPreferences, err := h.userService.UpdateUserPreferences(userID, &preferences)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPreferences)
}

// GetUserTheme handles the retrieval of user theme settings
func (h *UserSettingsHandler) GetUserTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	theme, err := h.userService.GetUserTheme(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User theme not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, theme)
}

// UpdateUserTheme handles the update of user theme settings
func (h *UserSettingsHandler) UpdateUserTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var theme models.UserTheme
	if err := json.NewDecoder(r.Body).Decode(&theme); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	theme.UserID = userID
	theme.UpdatedAt = time.Now()

	updatedTheme, err := h.userService.UpdateUserTheme(userID, &theme)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedTheme)
}

// GetUserLayout handles the retrieval of user layout settings
func (h *UserSettingsHandler) GetUserLayout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	layoutName := vars["layoutName"]

	layout, err := h.userService.GetUserLayout(userID, layoutName)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User layout not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, layout)
}

// GetAllUserLayouts handles the retrieval of all user layout settings
func (h *UserSettingsHandler) GetAllUserLayouts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	layouts, err := h.userService.GetAllUserLayouts(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User layouts not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, layouts)
}

// SaveUserLayout handles the saving of user layout settings
func (h *UserSettingsHandler) SaveUserLayout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var layout models.UserLayout
	if err := json.NewDecoder(r.Body).Decode(&layout); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	layout.UserID = userID
	layout.UpdatedAt = time.Now()

	savedLayout, err := h.userService.SaveUserLayout(userID, &layout)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, savedLayout)
}

// DeleteUserLayout handles the deletion of user layout settings
func (h *UserSettingsHandler) DeleteUserLayout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	layoutName := vars["layoutName"]

	err := h.userService.DeleteUserLayout(userID, layoutName)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Layout deleted successfully"})
}

// GetUserApiKeys handles the retrieval of user API keys
func (h *UserSettingsHandler) GetUserApiKeys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	apiKeys, err := h.userService.GetUserApiKeys(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User API keys not found")
		return
	}

	// Mask sensitive information
	for i := range apiKeys {
		apiKeys[i].ApiSecret = "********"
	}

	utils.RespondWithJSON(w, http.StatusOK, apiKeys)
}

// AddUserApiKey handles the addition of a new user API key
func (h *UserSettingsHandler) AddUserApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var apiKey models.UserApiKey
	if err := json.NewDecoder(r.Body).Decode(&apiKey); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	apiKey.UserID = userID
	apiKey.CreatedAt = time.Now()

	addedApiKey, err := h.userService.AddUserApiKey(userID, &apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Mask sensitive information
	addedApiKey.ApiSecret = "********"

	utils.RespondWithJSON(w, http.StatusCreated, addedApiKey)
}

// UpdateUserApiKey handles the update of a user API key
func (h *UserSettingsHandler) UpdateUserApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	keyID := vars["keyId"]

	var apiKey models.UserApiKey
	if err := json.NewDecoder(r.Body).Decode(&apiKey); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID and key ID
	apiKey.UserID = userID
	apiKey.ID = keyID
	apiKey.UpdatedAt = time.Now()

	updatedApiKey, err := h.userService.UpdateUserApiKey(userID, keyID, &apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Mask sensitive information
	updatedApiKey.ApiSecret = "********"

	utils.RespondWithJSON(w, http.StatusOK, updatedApiKey)
}

// DeleteUserApiKey handles the deletion of a user API key
func (h *UserSettingsHandler) DeleteUserApiKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	keyID := vars["keyId"]

	err := h.userService.DeleteUserApiKey(userID, keyID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "API key deleted successfully"})
}

// GetUserNotificationSettings handles the retrieval of user notification settings
func (h *UserSettingsHandler) GetUserNotificationSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	settings, err := h.userService.GetUserNotificationSettings(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User notification settings not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, settings)
}

// UpdateUserNotificationSettings handles the update of user notification settings
func (h *UserSettingsHandler) UpdateUserNotificationSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var settings models.UserNotificationSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set the user ID
	settings.UserID = userID
	settings.UpdatedAt = time.Now()

	updatedSettings, err := h.userService.UpdateUserNotificationSettings(userID, &settings)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedSettings)
}
