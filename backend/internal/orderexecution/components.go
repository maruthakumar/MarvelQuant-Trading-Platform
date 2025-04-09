	// Set start and end times to today
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startTime = today.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute + time.Duration(startTime.Second())*time.Second)
	endTime = today.Add(time.Duration(endTime.Hour())*time.Hour + time.Duration(endTime.Minute())*time.Minute + time.Duration(endTime.Second())*time.Second)

	// Check if current time is within execution window
	if now.Before(startTime) {
		// Schedule execution for start time
		go func() {
			timer := time.NewTimer(startTime.Sub(now))
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				// Execute order at start time
				s.engine.defaultExecute(ctx, order)
			}
		}()

		// Update order status
		order.Status = OrderStatusPending
		order.UpdatedAt = time.Now()
		s.engine.orderStore.UpdateOrder(order)

		return OrderResponse{
			Order:  order,
			Status: "SCHEDULED",
		}, nil
	} else if now.After(endTime) {
		// Execution window has passed
		order.Status = OrderStatusRejected
		order.ErrorMessage = "Execution window has passed"
		order.UpdatedAt = time.Now()
		s.engine.orderStore.UpdateOrder(order)

		return OrderResponse{
			Order:  order,
			Error:  order.ErrorMessage,
			Status: "ERROR",
		}, fmt.Errorf(order.ErrorMessage)
	}

	// Current time is within execution window, execute immediately
	return s.engine.defaultExecute(ctx, order)
}

// CanHandle checks if this strategy can handle an order
func (s *TimeBasedExecutionStrategy) CanHandle(order Order) bool {
	return order.Metadata.ExecutionMode == "Time-based"
}

// SignalBasedExecutionStrategy implements the ExecutionStrategy interface for signal-based execution
type SignalBasedExecutionStrategy struct {
	engine *OrderExecutionEngine
}

// NewSignalBasedExecutionStrategy creates a new signal-based execution strategy
func NewSignalBasedExecutionStrategy(engine *OrderExecutionEngine) *SignalBasedExecutionStrategy {
	return &SignalBasedExecutionStrategy{
		engine: engine,
	}
}

// Execute executes an order using signal-based execution
func (s *SignalBasedExecutionStrategy) Execute(ctx context.Context, order Order) (OrderResponse, error) {
	// Check if this is a signal-based execution
	if order.Metadata.ExecutionMode != "Signal-based" {
		return OrderResponse{
			Error:  "Not a signal-based execution",
			Status: "ERROR",
		}, fmt.Errorf("not a signal-based execution")
	}

	// In a real implementation, this would:
	// 1. Register for the specified signal
	// 2. Execute the order when the signal is received
	// 3. Handle signal timeout

	// For now, we'll simulate signal-based execution by executing immediately
	return s.engine.defaultExecute(ctx, order)
}

// CanHandle checks if this strategy can handle an order
func (s *SignalBasedExecutionStrategy) CanHandle(order Order) bool {
	return order.Metadata.ExecutionMode == "Signal-based"
}

// CombinedPremiumExecutionStrategy implements the ExecutionStrategy interface for combined premium execution
type CombinedPremiumExecutionStrategy struct {
	engine *OrderExecutionEngine
}

// NewCombinedPremiumExecutionStrategy creates a new combined premium execution strategy
func NewCombinedPremiumExecutionStrategy(engine *OrderExecutionEngine) *CombinedPremiumExecutionStrategy {
	return &CombinedPremiumExecutionStrategy{
		engine: engine,
	}
}

// Execute executes an order using combined premium execution
func (s *CombinedPremiumExecutionStrategy) Execute(ctx context.Context, order Order) (OrderResponse, error) {
	// Check if this is a combined premium execution
	if order.Metadata.ExecutionMode != "Combined Premium" {
		return OrderResponse{
			Error:  "Not a combined premium execution",
			Status: "ERROR",
		}, fmt.Errorf("not a combined premium execution")
	}

	// In a real implementation, this would:
	// 1. Calculate the combined premium of all legs
	// 2. Monitor the premium until it reaches the target value
	// 3. Execute all legs simultaneously when the target is reached

	// For now, we'll simulate combined premium execution by executing immediately
	return s.engine.defaultExecute(ctx, order)
}

// CanHandle checks if this strategy can handle an order
func (s *CombinedPremiumExecutionStrategy) CanHandle(order Order) bool {
	return order.Metadata.ExecutionMode == "Combined Premium"
}