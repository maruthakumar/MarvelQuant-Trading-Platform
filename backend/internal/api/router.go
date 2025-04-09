package api

import (
	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/api/handlers"
	"github.com/trading-platform/backend/internal/services"
	"github.com/trading-platform/backend/internal/services/position"
)

// Router sets up the API routes
type Router struct {
	router         *mux.Router
	orderHandler   *handlers.OrderHandler
	positionHandler *handlers.PositionHandler
}

// NewRouter creates a new Router
func NewRouter(orderService services.OrderService, positionService position.PositionService) *Router {
	router := mux.NewRouter()
	orderHandler := handlers.NewOrderHandler(orderService)
	positionHandler := handlers.NewPositionHandler(positionService)

	return &Router{
		router:         router,
		orderHandler:   orderHandler,
		positionHandler: positionHandler,
	}
}

// SetupRoutes configures all the routes for the API
func (r *Router) SetupRoutes() *mux.Router {
	// Order routes
	r.router.HandleFunc("/api/orders", r.orderHandler.CreateOrder).Methods("POST")
	r.router.HandleFunc("/api/orders", r.orderHandler.GetOrders).Methods("GET")
	r.router.HandleFunc("/api/orders/{id}", r.orderHandler.GetOrder).Methods("GET")
	r.router.HandleFunc("/api/orders/{id}", r.orderHandler.UpdateOrder).Methods("PUT")
	r.router.HandleFunc("/api/orders/{id}/cancel", r.orderHandler.CancelOrder).Methods("POST")
	
	// User-specific order routes
	r.router.HandleFunc("/api/users/{userId}/orders", r.orderHandler.GetOrdersByUser).Methods("GET")
	
	// Strategy-specific order routes
	r.router.HandleFunc("/api/strategies/{strategyId}/orders", r.orderHandler.GetOrdersByStrategy).Methods("GET")
	
	// Portfolio-specific order routes
	r.router.HandleFunc("/api/portfolios/{portfolioId}/orders", r.orderHandler.GetOrdersByPortfolio).Methods("GET")

	// Position routes
	r.router.HandleFunc("/api/positions/create-from-order", r.positionHandler.CreatePositionFromOrder).Methods("POST")
	r.router.HandleFunc("/api/positions", r.positionHandler.GetPositions).Methods("GET")
	r.router.HandleFunc("/api/positions/{id}", r.positionHandler.GetPosition).Methods("GET")
	r.router.HandleFunc("/api/positions/{id}", r.positionHandler.UpdatePosition).Methods("PUT")
	r.router.HandleFunc("/api/positions/{id}/close", r.positionHandler.ClosePosition).Methods("POST")
	r.router.HandleFunc("/api/positions/{id}/pnl", r.positionHandler.CalculatePnL).Methods("GET")
	r.router.HandleFunc("/api/positions/{id}/greeks", r.positionHandler.CalculateGreeks).Methods("GET")
	
	// User-specific position routes
	r.router.HandleFunc("/api/users/{userId}/positions", r.positionHandler.GetPositionsByUser).Methods("GET")
	r.router.HandleFunc("/api/users/{userId}/exposure", r.positionHandler.CalculateExposure).Methods("GET")
	
	// Strategy-specific position routes
	r.router.HandleFunc("/api/strategies/{strategyId}/positions", r.positionHandler.GetPositionsByStrategy).Methods("GET")
	
	// Portfolio-specific position routes
	r.router.HandleFunc("/api/portfolios/{portfolioId}/positions", r.positionHandler.GetPositionsByPortfolio).Methods("GET")
	
	// Position aggregation route
	r.router.HandleFunc("/api/positions/aggregate", r.positionHandler.AggregatePositions).Methods("GET")

	return r.router
}

// GetRouter returns the configured router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
