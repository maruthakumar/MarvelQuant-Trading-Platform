package multileg

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// MultilegService defines the interface for multileg strategy operations
type MultilegService interface {
	// Multileg strategy CRUD operations
	CreateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error)
	GetMultilegStrategyByID(id string) (*models.MultilegStrategy, error)
	GetMultilegStrategiesByUser(userID string) ([]models.MultilegStrategy, error)
	GetMultilegStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error)
	UpdateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error)
	DeleteMultilegStrategy(id string) error
	
	// Leg operations
	AddLeg(strategyID string, leg *models.Leg) (*models.Leg, error)
	UpdateLeg(strategyID string, leg *models.Leg) (*models.Leg, error)
	RemoveLeg(strategyID string, legID string) error
	GetLegsByStrategy(strategyID string) ([]models.Leg, error)
	
	// Execution operations
	ExecuteMultilegStrategy(strategyID string) error
	PauseMultilegStrategy(strategyID string) error
	ResumeMultilegStrategy(strategyID string) error
	CancelMultilegStrategy(strategyID string) error
	
	// Monitoring operations
	GetMultilegStrategyStatus(strategyID string) (string, error)
	GetMultilegStrategyPerformance(strategyID string) (*models.StrategyPerformance, error)
}

// MultilegServiceImpl implements the MultilegService interface
type MultilegServiceImpl struct {
	multilegRepo repositories.MultilegRepository
	orderRepo    repositories.OrderRepository
	positionRepo repositories.PositionRepository
	portfolioRepo repositories.PortfolioRepository
}

// NewMultilegService creates a new MultilegService
func NewMultilegService(
	multilegRepo repositories.MultilegRepository,
	orderRepo repositories.OrderRepository,
	positionRepo repositories.PositionRepository,
	portfolioRepo repositories.PortfolioRepository,
) MultilegService {
	return &MultilegServiceImpl{
		multilegRepo: multilegRepo,
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
		portfolioRepo: portfolioRepo,
	}
}

// CreateMultilegStrategy creates a new multileg strategy
func (s *MultilegServiceImpl) CreateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	if strategy == nil {
		return nil, errors.New("strategy cannot be nil")
	}
	
	// Validate strategy
	if err := strategy.Validate(); err != nil {
		return nil, err
	}
	
	// Verify portfolio exists
	_, err := s.portfolioRepo.GetByID(strategy.PortfolioID)
	if err != nil {
		return nil, errors.New("portfolio not found")
	}
	
	// Set timestamps
	now := time.Now()
	strategy.CreatedAt = now
	strategy.UpdatedAt = now
	
	// Set initial status
	strategy.Status = "DRAFT"
	
	// Create strategy
	return s.multilegRepo.CreateStrategy(strategy)
}

// GetMultilegStrategyByID retrieves a multileg strategy by ID
func (s *MultilegServiceImpl) GetMultilegStrategyByID(id string) (*models.MultilegStrategy, error) {
	if id == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	return s.multilegRepo.GetStrategyByID(id)
}

// GetMultilegStrategiesByUser retrieves all multileg strategies for a user
func (s *MultilegServiceImpl) GetMultilegStrategiesByUser(userID string) ([]models.MultilegStrategy, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	return s.multilegRepo.GetStrategiesByUser(userID)
}

// GetMultilegStrategiesByPortfolio retrieves all multileg strategies for a portfolio
func (s *MultilegServiceImpl) GetMultilegStrategiesByPortfolio(portfolioID string) ([]models.MultilegStrategy, error) {
	if portfolioID == "" {
		return nil, errors.New("portfolio ID cannot be empty")
	}
	
	return s.multilegRepo.GetStrategiesByPortfolio(portfolioID)
}

// UpdateMultilegStrategy updates an existing multileg strategy
func (s *MultilegServiceImpl) UpdateMultilegStrategy(strategy *models.MultilegStrategy) (*models.MultilegStrategy, error) {
	if strategy == nil {
		return nil, errors.New("strategy cannot be nil")
	}
	
	// Validate strategy
	if err := strategy.Validate(); err != nil {
		return nil, err
	}
	
	// Check if strategy exists
	existingStrategy, err := s.multilegRepo.GetStrategyByID(strategy.ID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Verify portfolio exists
	_, err = s.portfolioRepo.GetByID(strategy.PortfolioID)
	if err != nil {
		return nil, errors.New("portfolio not found")
	}
	
	// Preserve creation time and update the update time
	strategy.CreatedAt = existingStrategy.CreatedAt
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	return s.multilegRepo.UpdateStrategy(strategy)
}

// DeleteMultilegStrategy deletes a multileg strategy
func (s *MultilegServiceImpl) DeleteMultilegStrategy(id string) error {
	if id == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(id)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active
	if strategy.Status == "ACTIVE" {
		return errors.New("cannot delete active strategy")
	}
	
	// Delete strategy
	return s.multilegRepo.DeleteStrategy(id)
}

// AddLeg adds a leg to a multileg strategy
func (s *MultilegServiceImpl) AddLeg(strategyID string, leg *models.Leg) (*models.Leg, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	if leg == nil {
		return nil, errors.New("leg cannot be nil")
	}
	
	// Validate leg
	if err := leg.Validate(); err != nil {
		return nil, err
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Set strategy ID and timestamps
	leg.StrategyID = strategyID
	now := time.Now()
	leg.CreatedAt = now
	leg.UpdatedAt = now
	
	// Set initial status
	leg.Status = models.LegStatusPending
	
	// Add leg to strategy
	strategy.Legs = append(strategy.Legs, *leg)
	strategy.UpdatedAt = now
	
	// Update strategy
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return nil, err
	}
	
	return leg, nil
}

// UpdateLeg updates a leg in a multileg strategy
func (s *MultilegServiceImpl) UpdateLeg(strategyID string, leg *models.Leg) (*models.Leg, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	if leg == nil {
		return nil, errors.New("leg cannot be nil")
	}
	
	// Validate leg
	if err := leg.Validate(); err != nil {
		return nil, err
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Find and update the leg
	legFound := false
	for i, existingLeg := range strategy.Legs {
		if existingLeg.ID == leg.ID {
			// Preserve creation time and update the update time
			leg.CreatedAt = existingLeg.CreatedAt
			leg.UpdatedAt = time.Now()
			
			// Update the leg
			strategy.Legs[i] = *leg
			legFound = true
			break
		}
	}
	
	if !legFound {
		return nil, errors.New("leg not found")
	}
	
	// Update strategy
	strategy.UpdatedAt = time.Now()
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return nil, err
	}
	
	return leg, nil
}

// RemoveLeg removes a leg from a multileg strategy
func (s *MultilegServiceImpl) RemoveLeg(strategyID string, legID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	if legID == "" {
		return errors.New("leg ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Find and remove the leg
	legFound := false
	var updatedLegs []models.Leg
	for _, existingLeg := range strategy.Legs {
		if existingLeg.ID != legID {
			updatedLegs = append(updatedLegs, existingLeg)
		} else {
			legFound = true
		}
	}
	
	if !legFound {
		return errors.New("leg not found")
	}
	
	// Update strategy
	strategy.Legs = updatedLegs
	strategy.UpdatedAt = time.Now()
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// GetLegsByStrategy retrieves all legs for a multileg strategy
func (s *MultilegServiceImpl) GetLegsByStrategy(strategyID string) ([]models.Leg, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	return strategy.Legs, nil
}

// ExecuteMultilegStrategy executes a multileg strategy
func (s *MultilegServiceImpl) ExecuteMultilegStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is already active
	if strategy.Status == "ACTIVE" {
		return errors.New("strategy is already active")
	}
	
	// Update strategy status
	strategy.Status = "ACTIVE"
	strategy.UpdatedAt = time.Now()
	strategy.LastExecutedAt = time.Now()
	
	// Update strategy
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	// Execute strategy logic
	// This would typically involve creating orders for each leg based on the execution parameters
	// For now, we'll just update the status
	
	return nil
}

// PauseMultilegStrategy pauses a multileg strategy
func (s *MultilegServiceImpl) PauseMultilegStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active
	if strategy.Status != "ACTIVE" {
		return errors.New("strategy is not active")
	}
	
	// Update strategy status
	strategy.Status = "PAUSED"
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// ResumeMultilegStrategy resumes a paused multileg strategy
func (s *MultilegServiceImpl) ResumeMultilegStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is paused
	if strategy.Status != "PAUSED" {
		return errors.New("strategy is not paused")
	}
	
	// Update strategy status
	strategy.Status = "ACTIVE"
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// CancelMultilegStrategy cancels a multileg strategy
func (s *MultilegServiceImpl) CancelMultilegStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active or paused
	if strategy.Status != "ACTIVE" && strategy.Status != "PAUSED" {
		return errors.New("strategy is not active or paused")
	}
	
	// Update strategy status
	strategy.Status = "CANCELED"
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.multilegRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}
	
	// Cancel any active orders
	// This would typically involve canceling orders for each leg
	// For now, we'll just update the status
	
	return nil
}

// GetMultilegStrategyStatus retrieves the status of a multileg strategy
func (s *MultilegServiceImpl) GetMultilegStrategyStatus(strategyID string) (string, error) {
	if strategyID == "" {
		return "", errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return "", errors.New("strategy not found")
	}
	
	return strategy.Status, nil
}

// GetMultilegStrategyPerformance retrieves the performance of a multileg strategy
func (s *MultilegServiceImpl) GetMultilegStrategyPerformance(strategyID string) (*models.StrategyPerformance, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.multilegRepo.GetStrategyByID(strategyID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Get all orders for the strategy
	orders, err := s.orderRepo.GetByStrategy(strategyID)
	if err != nil {
		return nil, err
	}
	
	// Get all positions for the strategy
	positions, err := s.positionRepo.GetByStrategy(strategyID)
	if err != nil {
		return nil, err
	}
	
	// Calculate performance metrics
	var totalPnL float64
	var winCount, lossCount int
	var totalTrades int
	var maxDrawdown float64
	
	for _, position := range positions {
		totalPnL += position.RealizedPnL
		totalTrades++
		
		if position.RealizedPnL > 0 {
			winCount++
		} else if position.RealizedPnL < 0 {
			lossCount++
			
			// Calculate drawdown
			if position.RealizedPnL < maxDrawdown {
				maxDrawdown = position.RealizedPnL
			}
		}
	}
	
	// Calculate win rate
	var winRate float64
	if totalTrades > 0 {
		winRate = float64(winCount) / float64(totalTrades) * 100
	}
	
	// Create performance object
	performance := &models.StrategyPerformance{
		StrategyID:   strategyID,
		TotalPnL:     totalPnL,
		WinCount:     winCount,
		LossCount:    lossCount,
		TotalTrades:  totalTrades,
		WinRate:      winRate,
		MaxDrawdown:  maxDrawdown,
		OrderCount:   len(orders),
		PositionCount: len(positions),
		StartDate:    strategy.CreatedAt,
		EndDate:      time.Now(),
	}
	
	return performance, nil
}
