package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"trading_platform/backend/internal/api/handlers"
	"trading_platform/backend/internal/auth"
	"trading_platform/backend/internal/models"
	"trading_platform/backend/internal/repositories"
	"trading_platform/backend/internal/services/user"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *mux.Router, repos *repositories.Repositories) {
	// Create services
	userService := user.NewUserService(repos.UserRepository, repos.UserPreferencesRepository)
	environmentService := user.NewEnvironmentService(repos.UserRepository, repos.UserPreferencesRepository)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	environmentHandler := handlers.NewEnvironmentHandler(environmentService)

	// Public routes
	r.HandleFunc("/api/auth/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/refresh", userHandler.RefreshToken).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(auth.AuthMiddleware)

	// User routes
	protected.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/preferences", userHandler.GetPreferences).Methods("GET")
	protected.HandleFunc("/users/preferences", userHandler.UpdatePreferences).Methods("PUT")

	// API key routes
	protected.HandleFunc("/users/api-keys", userHandler.GetAPIKeys).Methods("GET")
	protected.HandleFunc("/users/api-keys", userHandler.CreateAPIKey).Methods("POST")
	protected.HandleFunc("/users/api-keys/{id}", userHandler.UpdateAPIKey).Methods("PUT")
	protected.HandleFunc("/users/api-keys/{id}", userHandler.DeleteAPIKey).Methods("DELETE")

	// Environment routes
	protected.HandleFunc("/environment/status", environmentHandler.GetEnvironmentStatus).Methods("GET")
	protected.HandleFunc("/environment/switch", environmentHandler.SwitchEnvironment).Methods("POST")

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(auth.RoleMiddleware(string(models.UserRoleAdmin)))
	admin.Use(auth.UserTypeMiddleware(string(models.UserTypeAdmin)))

	// SIM user management routes (admin only)
	admin.HandleFunc("/users/sim", userHandler.CreateSIMUser).Methods("POST")
	admin.HandleFunc("/users/sim", userHandler.GetSIMUsers).Methods("GET")
	admin.HandleFunc("/users/sim/{id}", userHandler.UpdateSIMUser).Methods("PUT")
	admin.HandleFunc("/users/sim/{id}", userHandler.DeleteSIMUser).Methods("DELETE")

	// SIM environment routes
	sim := protected.PathPrefix("/sim").Subrouter()
	sim.Use(auth.UserTypeMiddleware(string(models.UserTypeSIM)))
	sim.Use(auth.EnvironmentMiddleware(string(models.EnvironmentSIM)))
	sim.Use(auth.SimUserMiddleware)

	// Add SIM-specific routes here
}
