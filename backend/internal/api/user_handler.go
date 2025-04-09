package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"trading_platform/backend/internal/database"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/utils"
)

// UserHandler handles user-related API endpoints
type UserHandler struct {
	userRepo      *database.UserRepository
	preferenceRepo *database.UserPreferenceRepository
	apiKeyRepo    *database.APIKeyRepository
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	userRepo *database.UserRepository,
	preferenceRepo *database.UserPreferenceRepository,
	apiKeyRepo *database.APIKeyRepository,
) *UserHandler {
	return &UserHandler{
		userRepo:      userRepo,
		preferenceRepo: preferenceRepo,
		apiKeyRepo:    apiKeyRepo,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate user
	if err := user.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if username already exists
	existingUser, err := h.userRepo.GetByUsername(user.Username)
	if err == nil && existingUser != nil {
		utils.RespondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Check if email already exists
	existingUser, err = h.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.Password = string(hashedPassword)

	// Set default values
	user.Role = models.RoleUser
	user.Active = true
	user.EmailVerified = false
	user.PasswordChangedAt = time.Now()

	// Create user
	id, err := h.userRepo.Create(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Set ID in response
	user.ID = id

	// Create default preferences
	preferences := models.GetDefaultPreferences(id)
	_, err = h.preferenceRepo.Create(preferences)
	if err != nil {
		// Log error but don't fail registration
		// TODO: Add proper logging
	}

	// Remove password from response
	user.Password = ""

	utils.RespondWithJSON(w, http.StatusCreated, user)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user by username
	user, err := h.userRepo.GetByUsername(loginRequest.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	// Check if user is active
	if !user.Active {
		utils.RespondWithError(w, http.StatusUnauthorized, "Account is inactive")
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	// Update last login time
	user.LastLoginAt = time.Now()
	err = h.userRepo.Update(user)
	if err != nil {
		// Log error but don't fail login
		// TODO: Add proper logging
	}

	// Remove password from response
	user.Password = ""

	// Build response
	response := map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
		"user":         user,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate refresh token
	claims, err := auth.ValidateRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Get user by ID
	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	// Check if user is active
	if !user.Active {
		utils.RespondWithError(w, http.StatusUnauthorized, "Account is inactive")
		return
	}

	// Generate new JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Generate new refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	// Build response
	response := map[string]interface{}{
		"token":        token,
		"refreshToken": refreshToken,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetProfile handles retrieving the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	// Remove password from response
	user.Password = ""

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get existing user
	existingUser, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	// Parse request body
	var profileUpdate struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if email is being changed and if it already exists
	if profileUpdate.Email != existingUser.Email {
		userWithEmail, err := h.userRepo.GetByEmail(profileUpdate.Email)
		if err == nil && userWithEmail != nil {
			utils.RespondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
	}

	// Update user fields
	existingUser.FirstName = profileUpdate.FirstName
	existingUser.LastName = profileUpdate.LastName
	existingUser.Email = profileUpdate.Email
	existingUser.Phone = profileUpdate.Phone
	existingUser.UpdatedAt = time.Now()

	// Update user
	if err := h.userRepo.Update(existingUser); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating profile")
		return
	}

	// Remove password from response
	existingUser.Password = ""

	utils.RespondWithJSON(w, http.StatusOK, existingUser)
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get existing user
	existingUser, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	// Parse request body
	var passwordChange struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&passwordChange); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check current password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(passwordChange.CurrentPassword))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Validate new password
	if len(passwordChange.NewPassword) < 8 {
		utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChange.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Update user
	existingUser.Password = string(hashedPassword)
	existingUser.PasswordChangedAt = time.Now()
	existingUser.UpdatedAt = time.Now()

	if err := h.userRepo.Update(existingUser); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error changing password")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// GetPreferences handles retrieving the user's preferences
func (h *UserHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get preferences
	preferences, err := h.preferenceRepo.GetByUserID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving preferences")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, preferences)
}

// UpdatePreferences handles updating the user's preferences
func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get existing preferences
	existingPreferences, err := h.preferenceRepo.GetByUserID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving preferences")
		return
	}

	// Parse request body
	var updatedPreferences models.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&updatedPreferences); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Ensure user ID is not changed
	updatedPreferences.ID = existingPreferences.ID
	updatedPreferences.UserID = userID
	updatedPreferences.CreatedAt = existingPreferences.CreatedAt
	updatedPreferences.UpdatedAt = time.Now()

	// Update preferences
	if err := h.preferenceRepo.Update(&updatedPreferences); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating preferences")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPreferences)
}

// GetAPIKeys handles retrieving the user's API keys
func (h *UserHandler) GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get API keys
	apiKeys, err := h.apiKeyRepo.GetByUserID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving API keys")
		return
	}

	// Mask sensitive data
	for _, key := range apiKeys {
		key.MaskSensitiveData()
	}

	utils.RespondWithJSON(w, http.StatusOK, apiKeys)
}

// CreateAPIKey handles creating a new API key
func (h *UserHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var apiKey models.APIKey
	if err := json.NewDecoder(r.Body).Decode(&apiKey); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID
	apiKey.UserID = userID

	// Validate API key
	if err := apiKey.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create API key
	id, err := h.apiKeyRepo.Create(&apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating API key")
		return
	}

	// Set ID in response
	apiKey.ID = id

	utils.RespondWithJSON(w, http.StatusCreated, apiKey)
}

// UpdateAPIKey handles updating an API key
func (h *UserHandler) UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get API key ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing API key
	existingAPIKey, err := h.apiKeyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "API key not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving API key")
		}
		return
	}

	// Check if user has access to this API key
	if existingAPIKey.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var updatedAPIKey models.APIKey
	if err := json.NewDecoder(r.Body).Decode(&updatedAPIKey); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and user ID
	updatedAPIKey.ID = id
	updatedAPIKey.UserID = userID
	updatedAPIKey.CreatedAt = existingAPIKey.CreatedAt

	// Validate API key
	if err := updatedAPIKey.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update API key
	if err := h.apiKeyRepo.Update(&updatedAPIKey); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating API key")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedAPIKey)
}

// DeleteAPIKey handles deleting an API key
func (h *UserHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get API key ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing API key
	existingAPIKey, err := h.apiKeyRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "API key not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving API key")
		}
		return
	}

	// Check if user has access to this API key
	if existingAPIKey.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete API key
	if err := h.apiKeyRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting API key")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "API key deleted successfully"})
}

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(
	router *mux.Router,
	userRepo *database.UserRepository,
	preferenceRepo *database.UserPreferenceRepository,
	apiKeyRepo *database.APIKeyRepository,
	authMiddleware func(http.Handler) http.Handler,
) {
	handler := NewUserHandler(userRepo, preferenceRepo, apiKeyRepo)

	// Public routes
	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/refresh", handler.RefreshToken).Methods("POST")

	// Protected routes
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.Use(authMiddleware)

	userRouter.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	userRouter.HandleFunc("/profile", handler.UpdateProfile).Methods("PUT")
	userRouter.HandleFunc("/password", handler.ChangePassword).Methods("PUT")
	userRouter.HandleFunc("/preferences", handler.GetPreferences).Methods("GET")
	userRouter.HandleFunc("/preferences", handler.UpdatePreferences).Methods("PUT")
	userRouter.HandleFunc("/apikeys", handler.GetAPIKeys).Methods("GET")
	userRouter.HandleFunc("/apikeys", handler.CreateAPIKey).Methods("POST")
	userRouter.HandleFunc("/apikeys/{id}", handler.UpdateAPIKey).Methods("PUT")
	userRouter.HandleFunc("/apikeys/{id}", handler.DeleteAPIKey).Methods("DELETE")
}
