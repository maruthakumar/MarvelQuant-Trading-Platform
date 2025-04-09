package models

// UserType represents the type of user in the system
type UserType string

const (
	// UserTypeStandard represents a standard user with live trading capabilities
	UserTypeStandard UserType = "STANDARD"
	
	// UserTypeAdmin represents an administrator user
	UserTypeAdmin UserType = "ADMIN"
	
	// UserTypeSIM represents a simulation user for paper trading and backtesting
	UserTypeSIM UserType = "SIM"
)

// Environment represents the trading environment
type Environment string

const (
	// EnvironmentLive represents the live trading environment
	EnvironmentLive Environment = "LIVE"
	
	// EnvironmentSIM represents the simulation/paper trading environment
	EnvironmentSIM Environment = "SIM"
)
