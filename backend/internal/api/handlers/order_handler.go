package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/services"
	"github.com/trading-platform/backend/pkg/utils"
)

// OrderHandler handles HTTP requests related to orders
type OrderHandler struct {
	orderService services.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateOrder handles the creation of a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set creation time
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Validate the order
	if err := order.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the order
	createdOrder, err := h.orderService.CreateOrder(&order)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdOrder)
}

// GetOrder handles the retrieval of an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}

// GetOrders handles the retrieval of all orders with optional filtering
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filter := models.OrderFilter{}
	
	if userID := r.URL.Query().Get("userId"); userID != "" {
		filter.UserID = userID
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.OrderStatus(status)
	}
	if direction := r.URL.Query().Get("direction"); direction != "" {
		filter.Direction = models.OrderDirection(direction)
	}
	if productType := r.URL.Query().Get("productType"); productType != "" {
		filter.ProductType = models.ProductType(productType)
	}
	if instrumentType := r.URL.Query().Get("instrumentType"); instrumentType != "" {
		filter.InstrumentType = models.InstrumentType(instrumentType)
	}
	if portfolioID := r.URL.Query().Get("portfolioId"); portfolioID != "" {
		filter.PortfolioID = portfolioID
	}
	if strategyID := r.URL.Query().Get("strategyId"); strategyID != "" {
		filter.StrategyID = strategyID
	}

	// Parse date range if provided
	if fromDate := r.URL.Query().Get("fromDate"); fromDate != "" {
		parsedFromDate, err := time.Parse(time.RFC3339, fromDate)
		if err == nil {
			filter.FromDate = parsedFromDate
		}
	}
	if toDate := r.URL.Query().Get("toDate"); toDate != "" {
		parsedToDate, err := time.Parse(time.RFC3339, toDate)
		if err == nil {
			filter.ToDate = parsedToDate
		}
	}

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get orders with pagination
	orders, total, err := h.orderService.GetOrders(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// UpdateOrder handles the update of an existing order
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if order exists
	existingOrder, err := h.orderService.GetOrderByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	// Parse update data
	var orderUpdate models.Order
	if err := json.NewDecoder(r.Body).Decode(&orderUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set ID and update time
	orderUpdate.ID = id
	orderUpdate.UpdatedAt = time.Now()
	orderUpdate.CreatedAt = existingOrder.CreatedAt

	// Validate the updated order
	if err := orderUpdate.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the order
	updatedOrder, err := h.orderService.UpdateOrder(&orderUpdate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedOrder)
}

// CancelOrder handles the cancellation of an order
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if order exists
	existingOrder, err := h.orderService.GetOrderByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	// Check if order can be cancelled
	if existingOrder.Status != models.OrderStatusPending && existingOrder.Status != models.OrderStatusPartial {
		utils.RespondWithError(w, http.StatusBadRequest, "Only pending or partially filled orders can be cancelled")
		return
	}

	// Cancel the order
	err = h.orderService.CancelOrder(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order cancelled successfully"})
}

// GetOrdersByUser handles the retrieval of all orders for a specific user
func (h *OrderHandler) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for user
	filter := models.OrderFilter{
		UserID: userID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.OrderStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get orders with pagination
	orders, total, err := h.orderService.GetOrders(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetOrdersByStrategy handles the retrieval of all orders for a specific strategy
func (h *OrderHandler) GetOrdersByStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strategyID := vars["strategyId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for strategy
	filter := models.OrderFilter{
		StrategyID: strategyID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.OrderStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get orders with pagination
	orders, total, err := h.orderService.GetOrders(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetOrdersByPortfolio handles the retrieval of all orders for a specific portfolio
func (h *OrderHandler) GetOrdersByPortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID := vars["portfolioId"]

	// Parse pagination parameters
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := utils.ParseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := utils.ParseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create filter for portfolio
	filter := models.OrderFilter{
		PortfolioID: portfolioID,
	}

	// Apply additional filters if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.OrderStatus(status)
	}
	if symbol := r.URL.Query().Get("symbol"); symbol != "" {
		filter.Symbol = symbol
	}

	// Get orders with pagination
	orders, total, err := h.orderService.GetOrders(filter, page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare response with pagination metadata
	response := map[string]interface{}{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"totalPages":  (total + limit - 1) / limit,
		"hasNextPage": page*limit < total,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
