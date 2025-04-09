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

// OrderHandler handles order-related API endpoints
type OrderHandler struct {
	orderRepo *database.OrderRepository
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderRepo *database.OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
	}
}

// CreateOrder handles the creation of a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID
	order.UserID = userID

	// Validate order
	if err := order.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create order
	id, err := h.orderRepo.Create(&order)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating order")
		return
	}

	// Set ID in response
	order.ID = id

	utils.RespondWithJSON(w, http.StatusCreated, order)
}

// GetOrder handles retrieving an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get order
	order, err := h.orderRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order")
		}
		return
	}

	// Check if user has access to this order
	if order.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}

// UpdateOrder handles updating an order
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing order
	existingOrder, err := h.orderRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order")
		}
		return
	}

	// Check if user has access to this order
	if existingOrder.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set ID and user ID
	updatedOrder.ID = id
	updatedOrder.UserID = userID
	updatedOrder.CreatedAt = existingOrder.CreatedAt

	// Validate order
	if err := updatedOrder.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update order
	if err := h.orderRepo.Update(&updatedOrder); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating order")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedOrder)
}

// DeleteOrder handles deleting an order
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing order
	existingOrder, err := h.orderRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order")
		}
		return
	}

	// Check if user has access to this order
	if existingOrder.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete order
	if err := h.orderRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting order")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}

// GetOrders handles retrieving orders with filtering and pagination
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
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
	filter := models.OrderFilter{
		UserID: userID,
	}

	// Add optional filters
	if symbol := query.Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	if status := query.Get("status"); status != "" {
		filter.Status = models.OrderStatus(status)
	}

	if direction := query.Get("direction"); direction != "" {
		filter.Direction = models.OrderDirection(direction)
	}

	if productType := query.Get("productType"); productType != "" {
		filter.ProductType = models.ProductType(productType)
	}

	if instrumentType := query.Get("instrumentType"); instrumentType != "" {
		filter.InstrumentType = models.InstrumentType(instrumentType)
	}

	if portfolioID := query.Get("portfolioId"); portfolioID != "" {
		filter.PortfolioID = portfolioID
	}

	if strategyID := query.Get("strategyId"); strategyID != "" {
		filter.StrategyID = strategyID
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

	// Parse tags
	if tags := query.Get("tags"); tags != "" {
		filter.Tags = utils.SplitAndTrim(tags, ",")
	}

	// Get orders
	orders, total, err := h.orderRepo.Find(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving orders")
		return
	}

	// Build response with pagination
	response := map[string]interface{}{
		"data":       orders,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// CancelOrder handles cancelling an order
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing order
	existingOrder, err := h.orderRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving order")
		}
		return
	}

	// Check if user has access to this order
	if existingOrder.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if order can be cancelled
	if existingOrder.Status != models.OrderStatusPending && existingOrder.Status != models.OrderStatusPartial {
		utils.RespondWithError(w, http.StatusBadRequest, "Order cannot be cancelled")
		return
	}

	// Update order status
	existingOrder.Status = models.OrderStatusCancelled
	existingOrder.UpdatedAt = time.Now()

	// Update order
	if err := h.orderRepo.Update(existingOrder); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error cancelling order")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, existingOrder)
}

// RegisterOrderRoutes registers order-related routes
func RegisterOrderRoutes(router *mux.Router, orderRepo *database.OrderRepository, authMiddleware func(http.Handler) http.Handler) {
	handler := NewOrderHandler(orderRepo)
	
	// Apply auth middleware to all routes
	orderRouter := router.PathPrefix("/orders").Subrouter()
	orderRouter.Use(authMiddleware)

	// Register routes
	orderRouter.HandleFunc("", handler.CreateOrder).Methods("POST")
	orderRouter.HandleFunc("", handler.GetOrders).Methods("GET")
	orderRouter.HandleFunc("/{id}", handler.GetOrder).Methods("GET")
	orderRouter.HandleFunc("/{id}", handler.UpdateOrder).Methods("PUT")
	orderRouter.HandleFunc("/{id}", handler.DeleteOrder).Methods("DELETE")
	orderRouter.HandleFunc("/{id}/cancel", handler.CancelOrder).Methods("POST")
}
