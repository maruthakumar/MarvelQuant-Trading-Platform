package strategy

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/trading-platform/backend/internal/models"
)

// StrategyScheduler handles the scheduling of strategy executions
type StrategyScheduler struct {
	strategyService StrategyService
	executionEngine *StrategyExecutionEngine
	cronScheduler   *cron.Cron
	scheduleIDs     map[string]cron.EntryID
}

// NewStrategyScheduler creates a new StrategyScheduler
func NewStrategyScheduler(strategyService StrategyService, executionEngine *StrategyExecutionEngine) *StrategyScheduler {
	return &StrategyScheduler{
		strategyService: strategyService,
		executionEngine: executionEngine,
		cronScheduler:   cron.New(),
		scheduleIDs:     make(map[string]cron.EntryID),
	}
}

// StartScheduler starts the strategy scheduler
func (s *StrategyScheduler) StartScheduler() error {
	// Start the cron scheduler
	s.cronScheduler.Start()
	
	// Load all existing schedules
	err := s.loadAllSchedules()
	if err != nil {
		return err
	}
	
	return nil
}

// StopScheduler stops the strategy scheduler
func (s *StrategyScheduler) StopScheduler() {
	s.cronScheduler.Stop()
}

// ScheduleStrategy schedules a strategy for execution
func (s *StrategyScheduler) ScheduleStrategy(strategyID string, schedule *models.StrategySchedule) error {
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
	_, err := s.strategyService.GetStrategyByID(strategyID)
	if err != nil {
		return errors.New("strategy not found")
	}
	
	// Remove existing schedule if any
	s.removeSchedule(strategyID)
	
	// Create cron expression based on schedule frequency
	cronExpr, err := s.createCronExpression(schedule)
	if err != nil {
		return err
	}
	
	// Add the schedule to the cron scheduler
	entryID, err := s.cronScheduler.AddFunc(cronExpr, func() {
		// Only execute if schedule is enabled
		if schedule.Enabled {
			s.executionEngine.ExecuteStrategy(strategyID)
		}
	})
	
	if err != nil {
		return err
	}
	
	// Store the entry ID
	s.scheduleIDs[strategyID] = entryID
	
	// Save the schedule in the database
	return s.strategyService.ScheduleStrategy(strategyID, schedule)
}

// UpdateStrategySchedule updates the schedule for a strategy
func (s *StrategyScheduler) UpdateStrategySchedule(strategyID string, schedule *models.StrategySchedule) error {
	// This is essentially the same as scheduling a strategy
	return s.ScheduleStrategy(strategyID, schedule)
}

// RemoveStrategySchedule removes the schedule for a strategy
func (s *StrategyScheduler) RemoveStrategySchedule(strategyID string) error {
	// Remove from cron scheduler
	s.removeSchedule(strategyID)
	
	// Remove from database
	return s.strategyService.DeleteStrategySchedule(strategyID)
}

// loadAllSchedules loads all existing schedules from the database
func (s *StrategyScheduler) loadAllSchedules() error {
	// This is a simplified implementation
	// In a real system, you would query the database for all strategy schedules
	// and add them to the cron scheduler
	
	// For demonstration purposes, we'll just return nil
	return nil
}

// removeSchedule removes a schedule from the cron scheduler
func (s *StrategyScheduler) removeSchedule(strategyID string) {
	// Check if schedule exists
	entryID, exists := s.scheduleIDs[strategyID]
	if exists {
		// Remove from cron scheduler
		s.cronScheduler.Remove(entryID)
		
		// Remove from map
		delete(s.scheduleIDs, strategyID)
	}
}

// createCronExpression creates a cron expression based on the schedule frequency
func (s *StrategyScheduler) createCronExpression(schedule *models.StrategySchedule) (string, error) {
	// Extract time components
	hour := schedule.StartTime.Hour()
	minute := schedule.StartTime.Minute()
	
	// Create cron expression based on frequency
	switch schedule.Frequency {
	case models.ScheduleFrequencyOnce:
		// For one-time schedules, we'll use the exact date and time
		return fmt.Sprintf("%d %d %d %d *", minute, hour, schedule.StartTime.Day(), schedule.StartTime.Month()), nil
	
	case models.ScheduleFrequencyDaily:
		// Run daily at the specified time
		return fmt.Sprintf("%d %d * * *", minute, hour), nil
	
	case models.ScheduleFrequencyWeekly:
		// Run weekly on specified days at the specified time
		if len(schedule.DaysOfWeek) == 0 {
			return "", errors.New("days of week are required for weekly frequency")
		}
		
		// Convert days of week to cron format (0-6 where 0 is Sunday)
		daysExpr := ""
		for i, day := range schedule.DaysOfWeek {
			if i > 0 {
				daysExpr += ","
			}
			daysExpr += fmt.Sprintf("%d", day)
		}
		
		return fmt.Sprintf("%d %d * * %s", minute, hour, daysExpr), nil
	
	case models.ScheduleFrequencyMonthly:
		// Run monthly on the same day at the specified time
		return fmt.Sprintf("%d %d %d * *", minute, hour, schedule.StartTime.Day()), nil
	
	case models.ScheduleFrequencyCustom:
		// For custom schedules, we would need additional information
		// For now, we'll just return an error
		return "", errors.New("custom frequency not implemented")
	
	default:
		return "", errors.New("invalid schedule frequency")
	}
}
