package strategy

import (
	"errors"
	"time"

	"github.com/trading-platform/backend/internal/models"
	"github.com/trading-platform/backend/internal/repositories"
)

// StrategyService defines the interface for strategy-related operations
type StrategyService interface {
	// Strategy CRUD operations
	CreateStrategy(strategy *models.Strategy) (*models.Strategy, error)
	GetStrategyByID(id string) (*models.Strategy, error)
	GetStrategiesByUser(userID string) ([]models.Strategy, error)
	UpdateStrategy(strategy *models.Strategy) (*models.Strategy, error)
	DeleteStrategy(id string) error
	
	// Strategy execution operations
	ExecuteStrategy(strategyID string) error
	PauseStrategy(strategyID string) error
	ResumeStrategy(strategyID string) error
	StopStrategy(strategyID string) error
	
	// Strategy monitoring operations
	GetStrategyStatus(strategyID string) (models.StrategyStatus, error)
	GetStrategyPerformance(strategyID string) (*models.StrategyPerformance, error)
	
	// Strategy scheduling operations
	ScheduleStrategy(strategyID string, schedule *models.StrategySchedule) error
	GetStrategySchedule(strategyID string) (*models.StrategySchedule, error)
	UpdateStrategySchedule(strategyID string, schedule *models.StrategySchedule) error
	DeleteStrategySchedule(strategyID string) error
	
	// Strategy tagging operations
	AddStrategyTag(strategyID string, tag string) error
	RemoveStrategyTag(strategyID string, tag string) error
	GetStrategiesByTag(tag string) ([]models.Strategy, error)
}

// StrategyServiceImpl implements the StrategyService interface
type StrategyServiceImpl struct {
	strategyRepo repositories.StrategyRepository
	orderRepo    repositories.OrderRepository
	positionRepo repositories.PositionRepository
}

// NewStrategyService creates a new StrategyService
func NewStrategyService(
	strategyRepo repositories.StrategyRepository,
	orderRepo repositories.OrderRepository,
	positionRepo repositories.PositionRepository,
) StrategyService {
	return &StrategyServiceImpl{
		strategyRepo: strategyRepo,
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
	}
}

// CreateStrategy creates a new strategy
func (s *StrategyServiceImpl) CreateStrategy(strategy *models.Strategy) (*models.Strategy, error) {
	if strategy == nil {
		return nil, errors.New("strategy cannot be nil")
	}
	
	// Validate strategy
	if err := strategy.Validate(); err != nil {
		return nil, err
	}
	
	// Set timestamps
	now := time.Now()
	strategy.CreatedAt = now
	strategy.UpdatedAt = now
	
	// Set initial status
	strategy.Status = models.StrategyStatusDraft
	
	// Create strategy
	return s.strategyRepo.Create(strategy)
}

// GetStrategyByID retrieves a strategy by ID
func (s *StrategyServiceImpl) GetStrategyByID(id string) (*models.Strategy, error) {
	if id == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	return s.strategyRepo.GetByID(id)
}

// GetStrategiesByUser retrieves all strategies for a user
func (s *StrategyServiceImpl) GetStrategiesByUser(userID string) ([]models.Strategy, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	return s.strategyRepo.GetByUser(userID)
}

// UpdateStrategy updates an existing strategy
func (s *StrategyServiceImpl) UpdateStrategy(strategy *models.Strategy) (*models.Strategy, error) {
	if strategy == nil {
		return nil, errors.New("strategy cannot be nil")
	}
	
	// Validate strategy
	if err := strategy.Validate(); err != nil {
		return nil, err
	}
	
	// Check if strategy exists
	existingStrategy, err := s.strategyRepo.GetByID(strategy.ID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Preserve creation time and status
	strategy.CreatedAt = existingStrategy.CreatedAt
	strategy.Status = existingStrategy.Status
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	return s.strategyRepo.Update(strategy)
}

// DeleteStrategy deletes a strategy
func (s *StrategyServiceImpl) DeleteStrategy(id string) error {
	if id == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	_, err := s.strategyRepo.GetByID(id)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active
	strategy, _ := s.strategyRepo.GetByID(id)
	if strategy.Status == models.StrategyStatusActive {
		return errors.New("cannot delete active strategy")
	}
	
	// Delete strategy
	return s.strategyRepo.Delete(id)
}

// ExecuteStrategy executes a strategy
func (s *StrategyServiceImpl) ExecuteStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is already active
	if strategy.Status == models.StrategyStatusActive {
		return errors.New("strategy is already active")
	}
	
	// Update strategy status
	strategy.Status = models.StrategyStatusActive
	strategy.UpdatedAt = time.Now()
	strategy.LastExecutedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	// Execute strategy logic
	// This would typically involve creating orders based on the strategy rules
	// For now, we'll just update the status
	
	return nil
}

// PauseStrategy pauses a strategy
func (s *StrategyServiceImpl) PauseStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active
	if strategy.Status != models.StrategyStatusActive {
		return errors.New("strategy is not active")
	}
	
	// Update strategy status
	strategy.Status = models.StrategyStatusPaused
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// ResumeStrategy resumes a paused strategy
func (s *StrategyServiceImpl) ResumeStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is paused
	if strategy.Status != models.StrategyStatusPaused {
		return errors.New("strategy is not paused")
	}
	
	// Update strategy status
	strategy.Status = models.StrategyStatusActive
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// StopStrategy stops a strategy
func (s *StrategyServiceImpl) StopStrategy(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if strategy is active or paused
	if strategy.Status != models.StrategyStatusActive && strategy.Status != models.StrategyStatusPaused {
		return errors.New("strategy is not active or paused")
	}
	
	// Update strategy status
	strategy.Status = models.StrategyStatusStopped
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// GetStrategyStatus retrieves the status of a strategy
func (s *StrategyServiceImpl) GetStrategyStatus(strategyID string) (models.StrategyStatus, error) {
	if strategyID == "" {
		return "", errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return "", errors.New("strategy not found")
	}
	
	return strategy.Status, nil
}

// GetStrategyPerformance retrieves the performance of a strategy
func (s *StrategyServiceImpl) GetStrategyPerformance(strategyID string) (*models.StrategyPerformance, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
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

// ScheduleStrategy schedules a strategy for execution
func (s *StrategyServiceImpl) ScheduleStrategy(strategyID string, schedule *models.StrategySchedule) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	if schedule == nil {
		return errors.New("schedule cannot be nil")
	}
	
	// Validate schedule
	if err := schedule.Validate(); err != nil {
		return err
	}
	
	// Check if strategy exists
	_, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Set strategy ID and timestamps
	schedule.StrategyID = strategyID
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()
	
	// Create or update schedule
	return s.strategyRepo.SaveSchedule(schedule)
}

// GetStrategySchedule retrieves the schedule for a strategy
func (s *StrategyServiceImpl) GetStrategySchedule(strategyID string) (*models.StrategySchedule, error) {
	if strategyID == "" {
		return nil, errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	_, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return nil, errors.New("strategy not found")
	}
	
	// Get schedule
	return s.strategyRepo.GetSchedule(strategyID)
}

// UpdateStrategySchedule updates the schedule for a strategy
func (s *StrategyServiceImpl) UpdateStrategySchedule(strategyID string, schedule *models.StrategySchedule) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	if schedule == nil {
		return errors.New("schedule cannot be nil")
	}
	
	// Validate schedule
	if err := schedule.Validate(); err != nil {
		return err
	}
	
	// Check if strategy exists
	_, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if schedule exists
	existingSchedule, err := s.strategyRepo.GetSchedule(strategyID)
	if err != nil {
		return errors.New("schedule not found")
	}
	
	// Set strategy ID and timestamps
	schedule.StrategyID = strategyID
	schedule.CreatedAt = existingSchedule.CreatedAt
	schedule.UpdatedAt = time.Now()
	
	// Update schedule
	return s.strategyRepo.SaveSchedule(schedule)
}

// DeleteStrategySchedule deletes the schedule for a strategy
func (s *StrategyServiceImpl) DeleteStrategySchedule(strategyID string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	// Check if strategy exists
	_, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Delete schedule
	return s.strategyRepo.DeleteSchedule(strategyID)
}

// AddStrategyTag adds a tag to a strategy
func (s *StrategyServiceImpl) AddStrategyTag(strategyID string, tag string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	if tag == "" {
		return errors.New("tag cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Check if tag already exists
	for _, existingTag := range strategy.Tags {
		if existingTag == tag {
			return errors.New("tag already exists")
		}
	}
	
	// Add tag
	strategy.Tags = append(strategy.Tags, tag)
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// RemoveStrategyTag removes a tag from a strategy
func (s *StrategyServiceImpl) RemoveStrategyTag(strategyID string, tag string) error {
	if strategyID == "" {
		return errors.New("strategy ID cannot be empty")
	}
	
	if tag == "" {
		return errors.New("tag cannot be empty")
	}
	
	// Check if strategy exists
	strategy, err := s.strategyRepo.GetByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Find and remove tag
	found := false
	var updatedTags []string
	for _, existingTag := range strategy.Tags {
		if existingTag != tag {
			updatedTags = append(updatedTags, existingTag)
		} else {
			found = true
		}
	}
	
	if !found {
		return errors.New("tag not found")
	}
	
	// Update tags
	strategy.Tags = updatedTags
	strategy.UpdatedAt = time.Now()
	
	// Update strategy
	_, err = s.strategyRepo.Update(strategy)
	if err != nil {
		return err
	}
	
	return nil
}

// GetStrategiesByTag retrieves all strategies with a specific tag
func (s *StrategyServiceImpl) GetStrategiesByTag(tag string) ([]models.Strategy, error) {
	if tag == "" {
		return nil, errors.New("tag cannot be empty")
	}
	
	return s.strategyRepo.GetByTag(tag)
}
