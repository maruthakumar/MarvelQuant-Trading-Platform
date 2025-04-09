package auth

import (
        "net/http"
        "strings"

        "trading_platform/backend/internal/models"
        "trading_platform/backend/internal/utils"
)

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Get token from Authorization header
                authHeader := r.Header.Get("Authorization")
                if authHeader == "" {
                        utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
                        return
                }

                // Check if the header has the Bearer prefix
                if !strings.HasPrefix(authHeader, "Bearer ") {
                        utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
                        return
                }

                // Extract token
                tokenString := strings.TrimPrefix(authHeader, "Bearer ")

                // Validate token
                claims, err := ValidateToken(tokenString)
                if err != nil {
                        utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
                        return
                }

                // Set user ID, role, user type, and environment in context
                ctx := SetUserIDInContext(r.Context(), claims.UserID)
                ctx = SetRoleInContext(ctx, claims.Role)
                ctx = SetUserTypeInContext(ctx, claims.UserType)
                ctx = SetEnvironmentInContext(ctx, claims.Environment)

                // Call next handler with updated context
                next.ServeHTTP(w, r.WithContext(ctx))
        })
}

// RoleMiddleware is a middleware for role-based authorization
func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        // Get role from context
                        role := GetRoleFromContext(r.Context())
                        if role == "" {
                                utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
                                return
                        }

                        // Check if user has required role
                        hasRole := false
                        for _, allowedRole := range roles {
                                if role == allowedRole {
                                        hasRole = true
                                        break
                                }
                        }

                        if !hasRole {
                                utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
                                return
                        }

                        // Call next handler
                        next.ServeHTTP(w, r)
                })
        }
}

// UserTypeMiddleware is a middleware for user type-based authorization
func UserTypeMiddleware(userTypes ...string) func(http.Handler) http.Handler {
        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        // Get user type from context
                        userType := GetUserTypeFromContext(r.Context())
                        if userType == "" {
                                utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
                                return
                        }

                        // Check if user has required user type
                        hasUserType := false
                        for _, allowedUserType := range userTypes {
                                if userType == allowedUserType {
                                        hasUserType = true
                                        break
                                }
                        }

                        if !hasUserType {
                                utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
                                return
                        }

                        // Call next handler
                        next.ServeHTTP(w, r)
                })
        }
}

// EnvironmentMiddleware is a middleware for environment-based authorization
func EnvironmentMiddleware(environment string) func(http.Handler) http.Handler {
        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        // Get environment from context
                        currentEnv := GetEnvironmentFromContext(r.Context())
                        if currentEnv == "" {
                                utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
                                return
                        }

                        // Check if environment matches
                        if currentEnv != environment {
                                utils.RespondWithError(w, http.StatusForbidden, "Operation not allowed in current environment")
                                return
                        }

                        // Call next handler
                        next.ServeHTTP(w, r)
                })
        }
}

// SimUserMiddleware is a middleware that ensures only SIM users can access simulation environment
func SimUserMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Get user type and environment from context
                userType := GetUserTypeFromContext(r.Context())
                environment := GetEnvironmentFromContext(r.Context())

                // If environment is SIM, ensure user type is SIM
                if environment == string(models.EnvironmentSIM) && userType != string(models.UserTypeSIM) {
                        utils.RespondWithError(w, http.StatusForbidden, "Simulation environment is restricted to SIM users only")
                        return
                }

                // Call next handler
                next.ServeHTTP(w, r)
        })
}
