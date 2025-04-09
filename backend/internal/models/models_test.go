package models

import (
	"testing"
	"time"
)

func TestOrderValidation(t *testing.T) {
	// Test valid order
	validOrder := &Order{
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		OrderType:      OrderTypeLimit,
		Direction:      OrderDirectionBuy,
		Quantity:       10,
		FilledQuantity: 0,
		Price:          500.50,
		Status:         OrderStatusPending,
		ProductType:    ProductTypeMIS,
		InstrumentType: InstrumentTypeOption,
		OptionType:     OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
	}
	if err := validOrder.Validate(); err != nil {
		t.Errorf("Valid order failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name        string
		modifyOrder func(*Order)
	}{
		{
			name: "Missing UserID",
			modifyOrder: func(o *Order) {
				o.UserID = ""
			},
		},
		{
			name: "Missing Symbol",
			modifyOrder: func(o *Order) {
				o.Symbol = ""
			},
		},
		{
			name: "Missing Exchange",
			modifyOrder: func(o *Order) {
				o.Exchange = ""
			},
		},
		{
			name: "Invalid Quantity",
			modifyOrder: func(o *Order) {
				o.Quantity = 0
			},
		},
		{
			name: "Invalid OrderType",
			modifyOrder: func(o *Order) {
				o.OrderType = "INVALID"
			},
		},
		{
			name: "Invalid Direction",
			modifyOrder: func(o *Order) {
				o.Direction = "INVALID"
			},
		},
		{
			name: "Invalid ProductType",
			modifyOrder: func(o *Order) {
				o.ProductType = "INVALID"
			},
		},
		{
			name: "Invalid InstrumentType",
			modifyOrder: func(o *Order) {
				o.InstrumentType = "INVALID"
			},
		},
		{
			name: "Missing StrikePrice for Option",
			modifyOrder: func(o *Order) {
				o.StrikePrice = 0
			},
		},
		{
			name: "Missing Expiry for Option",
			modifyOrder: func(o *Order) {
				o.Expiry = time.Time{}
			},
		},
		{
			name: "Invalid OptionType",
			modifyOrder: func(o *Order) {
				o.OptionType = "INVALID"
			},
		},
		{
			name: "Missing Price for Limit Order",
			modifyOrder: func(o *Order) {
				o.Price = 0
			},
		},
		{
			name: "Invalid Status",
			modifyOrder: func(o *Order) {
				o.Status = "INVALID"
			},
		},
		{
			name: "Inconsistent Status and FilledQuantity",
			modifyOrder: func(o *Order) {
				o.Status = OrderStatusExecuted
				o.FilledQuantity = 5
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testOrder := *validOrder // Create a copy
			tc.modifyOrder(&testOrder)
			if err := testOrder.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestPositionValidation(t *testing.T) {
	// Test valid position
	validPosition := &Position{
		UserID:         "user123",
		Symbol:         "NIFTY",
		Exchange:       "NSE",
		Direction:      OrderDirectionBuy,
		Quantity:       10,
		EntryPrice:     500.50,
		CurrentPrice:   505.75,
		ProductType:    ProductTypeMIS,
		InstrumentType: InstrumentTypeOption,
		OptionType:     OptionTypeCall,
		StrikePrice:    18000,
		Expiry:         time.Now().AddDate(0, 1, 0),
		Status:         PositionStatusOpen,
		EntryTime:      time.Now(),
	}
	if err := validPosition.Validate(); err != nil {
		t.Errorf("Valid position failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name           string
		modifyPosition func(*Position)
	}{
		{
			name: "Missing UserID",
			modifyPosition: func(p *Position) {
				p.UserID = ""
			},
		},
		{
			name: "Missing Symbol",
			modifyPosition: func(p *Position) {
				p.Symbol = ""
			},
		},
		{
			name: "Missing Exchange",
			modifyPosition: func(p *Position) {
				p.Exchange = ""
			},
		},
		{
			name: "Invalid Quantity",
			modifyPosition: func(p *Position) {
				p.Quantity = 0
			},
		},
		{
			name: "Invalid EntryPrice",
			modifyPosition: func(p *Position) {
				p.EntryPrice = 0
			},
		},
		{
			name: "Invalid CurrentPrice",
			modifyPosition: func(p *Position) {
				p.CurrentPrice = 0
			},
		},
		{
			name: "Missing EntryTime",
			modifyPosition: func(p *Position) {
				p.EntryTime = time.Time{}
			},
		},
		{
			name: "Invalid Direction",
			modifyPosition: func(p *Position) {
				p.Direction = "INVALID"
			},
		},
		{
			name: "Invalid ProductType",
			modifyPosition: func(p *Position) {
				p.ProductType = "INVALID"
			},
		},
		{
			name: "Invalid InstrumentType",
			modifyPosition: func(p *Position) {
				p.InstrumentType = "INVALID"
			},
		},
		{
			name: "Missing StrikePrice for Option",
			modifyPosition: func(p *Position) {
				p.StrikePrice = 0
			},
		},
		{
			name: "Missing Expiry for Option",
			modifyPosition: func(p *Position) {
				p.Expiry = time.Time{}
			},
		},
		{
			name: "Invalid OptionType",
			modifyPosition: func(p *Position) {
				p.OptionType = "INVALID"
			},
		},
		{
			name: "Invalid Status",
			modifyPosition: func(p *Position) {
				p.Status = "INVALID"
			},
		},
		{
			name: "Closed Position without ExitTime",
			modifyPosition: func(p *Position) {
				p.Status = PositionStatusClosed
				p.ExitTime = time.Time{}
			},
		},
		{
			name: "Negative MarginUsed",
			modifyPosition: func(p *Position) {
				p.MarginUsed = -100
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPosition := *validPosition // Create a copy
			tc.modifyPosition(&testPosition)
			if err := testPosition.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestUserValidation(t *testing.T) {
	// Test valid user
	validUser := &User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "Test",
		LastName:     "User",
		Role:         UserRoleTrader,
		Phone:        "+1234567890",
	}
	if err := validUser.Validate(); err != nil {
		t.Errorf("Valid user failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name       string
		modifyUser func(*User)
	}{
		{
			name: "Missing Username",
			modifyUser: func(u *User) {
				u.Username = ""
			},
		},
		{
			name: "Missing Email",
			modifyUser: func(u *User) {
				u.Email = ""
			},
		},
		{
			name: "Missing PasswordHash",
			modifyUser: func(u *User) {
				u.PasswordHash = ""
			},
		},
		{
			name: "Missing FirstName",
			modifyUser: func(u *User) {
				u.FirstName = ""
			},
		},
		{
			name: "Missing LastName",
			modifyUser: func(u *User) {
				u.LastName = ""
			},
		},
		{
			name: "Invalid Username Format",
			modifyUser: func(u *User) {
				u.Username = "user@name"
			},
		},
		{
			name: "Invalid Email Format",
			modifyUser: func(u *User) {
				u.Email = "invalid-email"
			},
		},
		{
			name: "Invalid Phone Format",
			modifyUser: func(u *User) {
				u.Phone = "123"
			},
		},
		{
			name: "Invalid Role",
			modifyUser: func(u *User) {
				u.Role = "INVALID"
			},
		},
		{
			name: "Negative FailedLoginCount",
			modifyUser: func(u *User) {
				u.FailedLoginCount = -1
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testUser := *validUser // Create a copy
			tc.modifyUser(&testUser)
			if err := testUser.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}

	// Test user preferences validation
	validPrefs := &UserPreferences{
		UserID:               "user123",
		Theme:                "light",
		DefaultProductType:   ProductTypeMIS,
		DefaultOrderType:     OrderTypeMarket,
		DefaultQuantity:      1,
		DefaultSquareOffTime: "15:20:00",
		MaxDailyLoss:         5000,
		MaxPositionSize:      10,
		MaxOrdersPerMinute:   10,
		CircuitBreaker:       10,
		SlippageTolerance:    0.1,
		DataRefreshRate:      1000,
		GreeksPrecision:      4,
		SessionTimeout:       30,
	}
	if err := validPrefs.Validate(); err != nil {
		t.Errorf("Valid user preferences failed validation: %v", err)
	}

	// Test invalid user preferences
	prefsTests := []struct {
		name        string
		modifyPrefs func(*UserPreferences)
	}{
		{
			name: "Missing UserID",
			modifyPrefs: func(p *UserPreferences) {
				p.UserID = ""
			},
		},
		{
			name: "Invalid Theme",
			modifyPrefs: func(p *UserPreferences) {
				p.Theme = "invalid"
			},
		},
		{
			name: "Invalid DefaultProductType",
			modifyPrefs: func(p *UserPreferences) {
				p.DefaultProductType = "INVALID"
			},
		},
		{
			name: "Invalid DefaultOrderType",
			modifyPrefs: func(p *UserPreferences) {
				p.DefaultOrderType = "INVALID"
			},
		},
		{
			name: "Invalid DefaultQuantity",
			modifyPrefs: func(p *UserPreferences) {
				p.DefaultQuantity = 0
			},
		},
		{
			name: "Invalid DefaultSquareOffTime",
			modifyPrefs: func(p *UserPreferences) {
				p.DefaultSquareOffTime = "25:00:00"
			},
		},
		{
			name: "Negative MaxDailyLoss",
			modifyPrefs: func(p *UserPreferences) {
				p.MaxDailyLoss = -1
			},
		},
		{
			name: "Invalid MaxPositionSize",
			modifyPrefs: func(p *UserPreferences) {
				p.MaxPositionSize = 0
			},
		},
		{
			name: "Invalid MaxOrdersPerMinute",
			modifyPrefs: func(p *UserPreferences) {
				p.MaxOrdersPerMinute = 0
			},
		},
		{
			name: "Negative CircuitBreaker",
			modifyPrefs: func(p *UserPreferences) {
				p.CircuitBreaker = -1
			},
		},
		{
			name: "Negative SlippageTolerance",
			modifyPrefs: func(p *UserPreferences) {
				p.SlippageTolerance = -0.1
			},
		},
		{
			name: "Invalid DataRefreshRate",
			modifyPrefs: func(p *UserPreferences) {
				p.DataRefreshRate = 0
			},
		},
		{
			name: "Invalid GreeksPrecision",
			modifyPrefs: func(p *UserPreferences) {
				p.GreeksPrecision = 11
			},
		},
		{
			name: "Invalid SessionTimeout",
			modifyPrefs: func(p *UserPreferences) {
				p.SessionTimeout = 0
			},
		},
	}

	for _, tc := range prefsTests {
		t.Run(tc.name, func(t *testing.T) {
			testPrefs := *validPrefs // Create a copy
			tc.modifyPrefs(&testPrefs)
			if err := testPrefs.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestStrategyValidation(t *testing.T) {
	// Test valid strategy
	validStrategy := &Strategy{
		UserID:                "user123",
		Name:                  "Test Strategy",
		Symbol:                "NIFTY",
		Exchange:              "NSE",
		Type:                  StrategyTypeDirectional,
		ExecutionMode:         ExecutionModeTime,
		ProductType:           ProductTypeMIS,
		StartTime:             "09:15:00",
		EndTime:               "15:15:00",
		SquareOffTime:         "15:20:00",
		RunOnDays:             []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"},
		FailureAction:         FailureActionExitPlacedLegs,
		LegExecutionMode:      LegExecutionModeParallel,
		EntryOrderType:        OrderTypeMarket,
		ExitOrderType:         OrderTypeMarket,
		ExitMode:              ExitModeNormal,
		MaxRetries:            3,
		RetryInterval:         1,
		MaxLots:               5,
		MaxLossPerStrategy:    5000,
		MaxPositionSize:       10,
		LegMonitoringType:     MonitoringTypeRealtime,
		CombinedMonitoringType: MonitoringTypeRealtime,
		TargetType:            TargetTypeCombinedProfit,
		TargetValue:           1000,
		StopLossType:          StopLossTypeCombinedLoss,
		StopLossValue:         500,
	}
	if err := validStrategy.Validate(); err != nil {
		t.Errorf("Valid strategy failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name           string
		modifyStrategy func(*Strategy)
	}{
		{
			name: "Missing UserID",
			modifyStrategy: func(s *Strategy) {
				s.UserID = ""
			},
		},
		{
			name: "Missing Name",
			modifyStrategy: func(s *Strategy) {
				s.Name = ""
			},
		},
		{
			name: "Missing Symbol",
			modifyStrategy: func(s *Strategy) {
				s.Symbol = ""
			},
		},
		{
			name: "Missing Exchange",
			modifyStrategy: func(s *Strategy) {
				s.Exchange = ""
			},
		},
		{
			name: "Invalid Type",
			modifyStrategy: func(s *Strategy) {
				s.Type = "INVALID"
			},
		},
		{
			name: "Invalid ExecutionMode",
			modifyStrategy: func(s *Strategy) {
				s.ExecutionMode = "INVALID"
			},
		},
		{
			name: "Invalid ProductType",
			modifyStrategy: func(s *Strategy) {
				s.ProductType = "INVALID"
			},
		},
		{
			name: "Invalid StartTime",
			modifyStrategy: func(s *Strategy) {
				s.StartTime = "25:00:00"
			},
		},
		{
			name: "Invalid EndTime",
			modifyStrategy: func(s *Strategy) {
				s.EndTime = "25:00:00"
			},
		},
		{
			name: "Invalid SquareOffTime",
			modifyStrategy: func(s *Strategy) {
				s.SquareOffTime = "25:00:00"
			},
		},
		{
			name: "Invalid RunOnDays",
			modifyStrategy: func(s *Strategy) {
				s.RunOnDays = []string{"INVALID_DAY"}
			},
		},
		{
			name: "Invalid FailureAction",
			modifyStrategy: func(s *Strategy) {
				s.FailureAction = "INVALID"
			},
		},
		{
			name: "Invalid LegExecutionMode",
			modifyStrategy: func(s *Strategy) {
				s.LegExecutionMode = "INVALID"
			},
		},
		{
			name: "Invalid EntryOrderType",
			modifyStrategy: func(s *Strategy) {
				s.EntryOrderType = "INVALID"
			},
		},
		{
			name: "Invalid ExitOrderType",
			modifyStrategy: func(s *Strategy) {
				s.ExitOrderType = "INVALID"
			},
		},
		{
			name: "Invalid ExitMode",
			modifyStrategy: func(s *Strategy) {
				s.ExitMode = "INVALID"
			},
		},
		{
			name: "Negative ExitPriceBuffer",
			modifyStrategy: func(s *Strategy) {
				s.ExitPriceBuffer = -1
			},
		},
		{
			name: "Negative MaxRetries",
			modifyStrategy: func(s *Strategy) {
				s.MaxRetries = -1
			},
		},
		{
			name: "Negative RetryInterval",
			modifyStrategy: func(s *Strategy) {
				s.RetryInterval = -1
			},
		},
		{
			name: "Invalid MaxLots",
			modifyStrategy: func(s *Strategy) {
				s.MaxLots = 0
			},
		},
		{
			name: "Negative MaxLossPerStrategy",
			modifyStrategy: func(s *Strategy) {
				s.MaxLossPerStrategy = -1
			},
		},
		{
			name: "Invalid MaxPositionSize",
			modifyStrategy: func(s *Strategy) {
				s.MaxPositionSize = 0
			},
		},
		{
			name: "Invalid LegMonitoringType",
			modifyStrategy: func(s *Strategy) {
				s.LegMonitoringType = "INVALID"
			},
		},
		{
			name: "Invalid CombinedMonitoringType",
			modifyStrategy: func(s *Strategy) {
				s.CombinedMonitoringType = "INVALID"
			},
		},
		{
			name: "Missing MonitoringInterval for Interval Monitoring",
			modifyStrategy: func(s *Strategy) {
				s.LegMonitoringType = MonitoringTypeInterval
				s.MonitoringInterval = 0
			},
		},
		{
			name: "Invalid TargetType",
			modifyStrategy: func(s *Strategy) {
				s.TargetType = "INVALID"
			},
		},
		{
			name: "Invalid TargetValue",
			modifyStrategy: func(s *Strategy) {
				s.TargetValue = 0
			},
		},
		{
			name: "Invalid StopLossType",
			modifyStrategy: func(s *Strategy) {
				s.StopLossType = "INVALID"
			},
		},
		{
			name: "Invalid StopLossValue",
			modifyStrategy: func(s *Strategy) {
				s.StopLossValue = 0
			},
		},
		{
			name: "Negative StopLossWaitSeconds",
			modifyStrategy: func(s *Strategy) {
				s.StopLossWaitSeconds = -1
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testStrategy := *validStrategy // Create a copy
			tc.modifyStrategy(&testStrategy)
			if err := testStrategy.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestPortfolioValidation(t *testing.T) {
	// Create a valid leg for the portfolio
	validLeg := &Leg{
		ID:                1,
		PortfolioID:       "portfolio123",
		Symbol:            "NIFTY",
		Exchange:          "NSE",
		Type:              LegTypeOption,
		BuySell:           string(OrderDirectionBuy),
		OptionType:        string(OptionTypeCall),
		StrikePrice:       18000,
		Expiry:            time.Now().AddDate(0, 1, 0),
		Lots:              1,
		LotSize:           50,
		StrikeSelectionMode: StrikeSelectionModeNormal,
		EntryOrderType:    OrderTypeMarket,
		ExitOrderType:     OrderTypeMarket,
		EntryPriceBuffer:  0.1,
		ExitPriceBuffer:   0.1,
		MaxEntryRetries:   3,
		MaxExitRetries:    3,
		EntryRetryInterval: 1,
		ExitRetryInterval: 1,
		ExecutionPriority: 1,
		Status:            "PENDING",
	}

	// Test valid portfolio
	validPortfolio := &Portfolio{
		UserID:                "user123",
		Name:                  "Test Portfolio",
		Symbol:                "NIFTY",
		Exchange:              "NSE",
		Expiry:                time.Now().AddDate(0, 1, 0),
		DefaultLots:           1,
		Status:                PortfolioStatusPending,
		StrikeSelection:       StrikeSelectionModeNormal,
		UnderlyingRef:         UnderlyingReferenceFuture,
		PriceType:             PriceTypeLTP,
		StrikeStep:            50,
		ProductType:           ProductTypeMIS,
		FailureAction:         FailureActionExitPlacedLegs,
		LegExecutionMode:      LegExecutionModeParallel,
		MaxLots:               5,
		RunOnDays:             []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"},
		StartTime:             "09:15:00",
		EndTime:               "15:15:00",
		SquareOffTime:         "15:20:00",
		ExecutionMode:         ExecutionModeTime,
		EntryOrderType:        OrderTypeMarket,
		LegMonitoringType:     MonitoringTypeRealtime,
		CombinedMonitoringType: MonitoringTypeRealtime,
		TargetType:            TargetTypeCombinedProfit,
		TargetValue:           1000,
		StopLossType:          StopLossTypeCombinedLoss,
		StopLossValue:         500,
		ExitMode:              ExitModeNormal,
		ExitOrderType:         OrderTypeMarket,
		ExitPriceBuffer:       0.1,
		MaxExitRetries:        3,
		ExitRetryInterval:     1,
		Legs:                  []Leg{*validLeg},
	}
	if err := validPortfolio.Validate(); err != nil {
		t.Errorf("Valid portfolio failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name            string
		modifyPortfolio func(*Portfolio)
	}{
		{
			name: "Missing UserID",
			modifyPortfolio: func(p *Portfolio) {
				p.UserID = ""
			},
		},
		{
			name: "Missing Name",
			modifyPortfolio: func(p *Portfolio) {
				p.Name = ""
			},
		},
		{
			name: "Missing Symbol",
			modifyPortfolio: func(p *Portfolio) {
				p.Symbol = ""
			},
		},
		{
			name: "Missing Exchange",
			modifyPortfolio: func(p *Portfolio) {
				p.Exchange = ""
			},
		},
		{
			name: "Missing Expiry",
			modifyPortfolio: func(p *Portfolio) {
				p.Expiry = time.Time{}
			},
		},
		{
			name: "Invalid DefaultLots",
			modifyPortfolio: func(p *Portfolio) {
				p.DefaultLots = 0
			},
		},
		{
			name: "Invalid Status",
			modifyPortfolio: func(p *Portfolio) {
				p.Status = "INVALID"
			},
		},
		{
			name: "Invalid StrikeSelection",
			modifyPortfolio: func(p *Portfolio) {
				p.StrikeSelection = "INVALID"
			},
		},
		{
			name: "Invalid UnderlyingRef",
			modifyPortfolio: func(p *Portfolio) {
				p.UnderlyingRef = "INVALID"
			},
		},
		{
			name: "Invalid PriceType",
			modifyPortfolio: func(p *Portfolio) {
				p.PriceType = "INVALID"
			},
		},
		{
			name: "Invalid StrikeStep",
			modifyPortfolio: func(p *Portfolio) {
				p.StrikeStep = 0
			},
		},
		{
			name: "Invalid ProductType",
			modifyPortfolio: func(p *Portfolio) {
				p.ProductType = "INVALID"
			},
		},
		{
			name: "Invalid FailureAction",
			modifyPortfolio: func(p *Portfolio) {
				p.FailureAction = "INVALID"
			},
		},
		{
			name: "Invalid LegExecutionMode",
			modifyPortfolio: func(p *Portfolio) {
				p.LegExecutionMode = "INVALID"
			},
		},
		{
			name: "Invalid MaxLots",
			modifyPortfolio: func(p *Portfolio) {
				p.MaxLots = 0
			},
		},
		{
			name: "Invalid RunOnDays",
			modifyPortfolio: func(p *Portfolio) {
				p.RunOnDays = []string{"INVALID_DAY"}
			},
		},
		{
			name: "Invalid StartTime",
			modifyPortfolio: func(p *Portfolio) {
				p.StartTime = "25:00:00"
			},
		},
		{
			name: "Invalid EndTime",
			modifyPortfolio: func(p *Portfolio) {
				p.EndTime = "25:00:00"
			},
		},
		{
			name: "Invalid SquareOffTime",
			modifyPortfolio: func(p *Portfolio) {
				p.SquareOffTime = "25:00:00"
			},
		},
		{
			name: "Invalid ExecutionMode",
			modifyPortfolio: func(p *Portfolio) {
				p.ExecutionMode = "INVALID"
			},
		},
		{
			name: "Invalid EntryOrderType",
			modifyPortfolio: func(p *Portfolio) {
				p.EntryOrderType = "INVALID"
			},
		},
		{
			name: "Invalid LegMonitoringType",
			modifyPortfolio: func(p *Portfolio) {
				p.LegMonitoringType = "INVALID"
			},
		},
		{
			name: "Invalid CombinedMonitoringType",
			modifyPortfolio: func(p *Portfolio) {
				p.CombinedMonitoringType = "INVALID"
			},
		},
		{
			name: "Missing MonitoringInterval for Interval Monitoring",
			modifyPortfolio: func(p *Portfolio) {
				p.LegMonitoringType = MonitoringTypeInterval
				p.MonitoringInterval = 0
			},
		},
		{
			name: "Invalid TargetType",
			modifyPortfolio: func(p *Portfolio) {
				p.TargetType = "INVALID"
			},
		},
		{
			name: "Invalid TargetValue",
			modifyPortfolio: func(p *Portfolio) {
				p.TargetValue = 0
			},
		},
		{
			name: "Invalid StopLossType",
			modifyPortfolio: func(p *Portfolio) {
				p.StopLossType = "INVALID"
			},
		},
		{
			name: "Invalid StopLossValue",
			modifyPortfolio: func(p *Portfolio) {
				p.StopLossValue = 0
			},
		},
		{
			name: "Invalid ExitMode",
			modifyPortfolio: func(p *Portfolio) {
				p.ExitMode = "INVALID"
			},
		},
		{
			name: "Invalid ExitOrderType",
			modifyPortfolio: func(p *Portfolio) {
				p.ExitOrderType = "INVALID"
			},
		},
		{
			name: "Negative ExitPriceBuffer",
			modifyPortfolio: func(p *Portfolio) {
				p.ExitPriceBuffer = -1
			},
		},
		{
			name: "Negative MaxExitRetries",
			modifyPortfolio: func(p *Portfolio) {
				p.MaxExitRetries = -1
			},
		},
		{
			name: "Negative ExitRetryInterval",
			modifyPortfolio: func(p *Portfolio) {
				p.ExitRetryInterval = -1
			},
		},
		{
			name: "Invalid Partial Exit Settings",
			modifyPortfolio: func(p *Portfolio) {
				p.EnablePartialExits = true
				p.MinExitPercentage = 0
			},
		},
		{
			name: "No Legs",
			modifyPortfolio: func(p *Portfolio) {
				p.Legs = []Leg{}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testPortfolio := *validPortfolio // Create a copy
			tc.modifyPortfolio(&testPortfolio)
			if err := testPortfolio.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestLegValidation(t *testing.T) {
	// Test valid leg
	validLeg := &Leg{
		PortfolioID:       "portfolio123",
		Symbol:            "NIFTY",
		Exchange:          "NSE",
		Type:              LegTypeOption,
		BuySell:           string(OrderDirectionBuy),
		OptionType:        string(OptionTypeCall),
		StrikePrice:       18000,
		Expiry:            time.Now().AddDate(0, 1, 0),
		Lots:              1,
		LotSize:           50,
		StrikeSelectionMode: StrikeSelectionModeNormal,
		EntryOrderType:    OrderTypeMarket,
		ExitOrderType:     OrderTypeMarket,
		EntryPriceBuffer:  0.1,
		ExitPriceBuffer:   0.1,
		MaxEntryRetries:   3,
		MaxExitRetries:    3,
		EntryRetryInterval: 1,
		ExitRetryInterval: 1,
		ExecutionPriority: 1,
		Status:            "PENDING",
	}
	if err := validLeg.Validate(); err != nil {
		t.Errorf("Valid leg failed validation: %v", err)
	}

	// Test missing required fields
	tests := []struct {
		name      string
		modifyLeg func(*Leg)
	}{
		{
			name: "Missing PortfolioID",
			modifyLeg: func(l *Leg) {
				l.PortfolioID = ""
			},
		},
		{
			name: "Missing Symbol",
			modifyLeg: func(l *Leg) {
				l.Symbol = ""
			},
		},
		{
			name: "Missing Exchange",
			modifyLeg: func(l *Leg) {
				l.Exchange = ""
			},
		},
		{
			name: "Invalid Lots",
			modifyLeg: func(l *Leg) {
				l.Lots = 0
			},
		},
		{
			name: "Invalid LotSize",
			modifyLeg: func(l *Leg) {
				l.LotSize = 0
			},
		},
		{
			name: "Invalid Type",
			modifyLeg: func(l *Leg) {
				l.Type = "INVALID"
			},
		},
		{
			name: "Invalid BuySell",
			modifyLeg: func(l *Leg) {
				l.BuySell = "INVALID"
			},
		},
		{
			name: "Missing StrikePrice for Option",
			modifyLeg: func(l *Leg) {
				l.StrikePrice = 0
			},
		},
		{
			name: "Missing Expiry for Option",
			modifyLeg: func(l *Leg) {
				l.Expiry = time.Time{}
			},
		},
		{
			name: "Invalid OptionType",
			modifyLeg: func(l *Leg) {
				l.OptionType = "INVALID"
			},
		},
		{
			name: "Invalid StrikeSelectionMode",
			modifyLeg: func(l *Leg) {
				l.StrikeSelectionMode = "INVALID"
			},
		},
		{
			name: "Invalid EntryOrderType",
			modifyLeg: func(l *Leg) {
				l.EntryOrderType = "INVALID"
			},
		},
		{
			name: "Invalid ExitOrderType",
			modifyLeg: func(l *Leg) {
				l.ExitOrderType = "INVALID"
			},
		},
		{
			name: "Missing EntryLimitPrice for Limit Order",
			modifyLeg: func(l *Leg) {
				l.EntryOrderType = OrderTypeLimit
				l.EntryLimitPrice = 0
			},
		},
		{
			name: "Missing EntryTriggerPrice for SL Limit Order",
			modifyLeg: func(l *Leg) {
				l.EntryOrderType = OrderTypeSLLimit
				l.EntryTriggerPrice = 0
			},
		},
		{
			name: "Negative EntryPriceBuffer",
			modifyLeg: func(l *Leg) {
				l.EntryPriceBuffer = -1
			},
		},
		{
			name: "Negative ExitPriceBuffer",
			modifyLeg: func(l *Leg) {
				l.ExitPriceBuffer = -1
			},
		},
		{
			name: "Negative MaxEntryRetries",
			modifyLeg: func(l *Leg) {
				l.MaxEntryRetries = -1
			},
		},
		{
			name: "Negative EntryRetryInterval",
			modifyLeg: func(l *Leg) {
				l.EntryRetryInterval = -1
			},
		},
		{
			name: "Negative MaxExitRetries",
			modifyLeg: func(l *Leg) {
				l.MaxExitRetries = -1
			},
		},
		{
			name: "Negative ExitRetryInterval",
			modifyLeg: func(l *Leg) {
				l.ExitRetryInterval = -1
			},
		},
		{
			name: "Negative ExecutionPriority",
			modifyLeg: func(l *Leg) {
				l.ExecutionPriority = -1
			},
		},
		{
			name: "Invalid IndividualTarget",
			modifyLeg: func(l *Leg) {
				l.IndividualTarget = -1
			},
		},
		{
			name: "Invalid IndividualStopLoss",
			modifyLeg: func(l *Leg) {
				l.IndividualStopLoss = -1
			},
		},
		{
			name: "Missing TrailValue when Trailing Enabled",
			modifyLeg: func(l *Leg) {
				l.TrailTarget = true
				l.TrailValue = 0
			},
		},
		{
			name: "Invalid Status",
			modifyLeg: func(l *Leg) {
				l.Status = "INVALID"
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testLeg := *validLeg // Create a copy
			tc.modifyLeg(&testLeg)
			if err := testLeg.Validate(); err == nil {
				t.Errorf("Expected validation error for %s, but got none", tc.name)
			}
		})
	}
}

func TestCalculationMethods(t *testing.T) {
	// Test Order calculation methods
	order := &Order{
		Direction:      OrderDirectionBuy,
		Quantity:       10,
		FilledQuantity: 5,
		Price:          100,
		AveragePrice:   102,
		OrderType:      OrderTypeLimit,
		Status:         OrderStatusPartial,
	}

	// Test CalculateSlippage
	slippage := order.CalculateSlippage()
	if slippage != 2 {
		t.Errorf("Expected slippage to be 2, got %f", slippage)
	}

	// Test IsComplete
	if order.IsComplete() {
		t.Errorf("Expected IsComplete to return false for partial order")
	}

	// Test RemainingQuantity
	if order.RemainingQuantity() != 5 {
		t.Errorf("Expected RemainingQuantity to be 5, got %d", order.RemainingQuantity())
	}

	// Test Position calculation methods
	position := &Position{
		Direction:    OrderDirectionBuy,
		Quantity:     10,
		EntryPrice:   100,
		CurrentPrice: 110,
		RealizedPnL:  50,
	}

	// Test CalculatePnL
	position.CalculatePnL()
	if position.UnrealizedPnL != 100 {
		t.Errorf("Expected UnrealizedPnL to be 100, got %f", position.UnrealizedPnL)
	}
	if position.TotalPnL != 150 {
		t.Errorf("Expected TotalPnL to be 150, got %f", position.TotalPnL)
	}
	if position.PnLPercentage != 15 {
		t.Errorf("Expected PnLPercentage to be 15, got %f", position.PnLPercentage)
	}

	// Test UpdateCurrentPrice
	position.UpdateCurrentPrice(120)
	if position.CurrentPrice != 120 {
		t.Errorf("Expected CurrentPrice to be 120, got %f", position.CurrentPrice)
	}
	if position.UnrealizedPnL != 200 {
		t.Errorf("Expected UnrealizedPnL to be 200, got %f", position.UnrealizedPnL)
	}

	// Test PartialClose
	closedPosition, err := position.PartialClose(5, 120)
	if err != nil {
		t.Errorf("PartialClose failed: %v", err)
	}
	if position.Quantity != 5 {
		t.Errorf("Expected remaining quantity to be 5, got %d", position.Quantity)
	}
	if position.Status != PositionStatusPartial {
		t.Errorf("Expected status to be PARTIAL, got %s", position.Status)
	}
	if closedPosition.Quantity != 5 {
		t.Errorf("Expected closed position quantity to be 5, got %d", closedPosition.Quantity)
	}
	if closedPosition.Status != PositionStatusClosed {
		t.Errorf("Expected closed position status to be CLOSED, got %s", closedPosition.Status)
	}
	if closedPosition.RealizedPnL != 100 {
		t.Errorf("Expected closed position RealizedPnL to be 100, got %f", closedPosition.RealizedPnL)
	}

	// Test Leg calculation methods
	leg := &Leg{
		BuySell:      string(OrderDirectionBuy),
		Lots:         2,
		LotSize:      50,
		EntryPrice:   100,
		CurrentPrice: 110,
		RealizedPnL:  500,
	}

	// Test CalculateQuantity
	leg.CalculateQuantity()
	if leg.Quantity != 100 {
		t.Errorf("Expected Quantity to be 100, got %d", leg.Quantity)
	}

	// Test CalculatePnL
	leg.CalculatePnL()
	if leg.UnrealizedPnL != 1000 {
		t.Errorf("Expected UnrealizedPnL to be 1000, got %f", leg.UnrealizedPnL)
	}
	if leg.TotalPnL != 1500 {
		t.Errorf("Expected TotalPnL to be 1500, got %f", leg.TotalPnL)
	}
	if leg.PnLPercentage != 15 {
		t.Errorf("Expected PnLPercentage to be 15, got %f", leg.PnLPercentage)
	}

	// Test UpdateCurrentPrice
	leg.UpdateCurrentPrice(120)
	if leg.CurrentPrice != 120 {
		t.Errorf("Expected CurrentPrice to be 120, got %f", leg.CurrentPrice)
	}
	if leg.UnrealizedPnL != 2000 {
		t.Errorf("Expected UnrealizedPnL to be 2000, got %f", leg.UnrealizedPnL)
	}
}
