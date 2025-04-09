package models

import (
        "errors"
        "regexp"
        "time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
        UserRoleAdmin  UserRole = "ADMIN"
        UserRoleTrader UserRole = "TRADER"
        UserRoleViewer UserRole = "VIEWER"
)

// User represents a user in the trading system
type User struct {
        ID                string    `json:"id" bson:"_id,omitempty"`
        Username          string    `json:"username" bson:"username"`
        Email             string    `json:"email" bson:"email"`
        PasswordHash      string    `json:"-" bson:"passwordHash"`
        FirstName         string    `json:"firstName" bson:"firstName"`
        LastName          string    `json:"lastName" bson:"lastName"`
        Role              UserRole  `json:"role" bson:"role"`
        UserType          UserType  `json:"userType" bson:"userType"`
        Active            bool      `json:"active" bson:"active"`
        Phone             string    `json:"phone,omitempty" bson:"phone,omitempty"`
        TwoFactorEnabled  bool      `json:"twoFactorEnabled" bson:"twoFactorEnabled"`
        TwoFactorSecret   string    `json:"-" bson:"twoFactorSecret,omitempty"`
        LastLogin         time.Time `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
        FailedLoginCount  int       `json:"-" bson:"failedLoginCount"`
        LockedUntil       time.Time `json:"-" bson:"lockedUntil,omitempty"`
        PasswordChangedAt time.Time `json:"-" bson:"passwordChangedAt"`
        CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
        UpdatedAt         time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserPreferences represents user-specific settings and preferences
type UserPreferences struct {
        ID                   string            `json:"id" bson:"_id,omitempty"`
        UserID               string            `json:"userId" bson:"userId"`
        Theme                string            `json:"theme" bson:"theme"`
        DefaultEnvironment   Environment       `json:"defaultEnvironment" bson:"defaultEnvironment"`
        DefaultProductType   ProductType       `json:"defaultProductType" bson:"defaultProductType"`
        DefaultOrderType     OrderType         `json:"defaultOrderType" bson:"defaultOrderType"`
        DefaultQuantity      int               `json:"defaultQuantity" bson:"defaultQuantity"`
        DefaultSquareOffTime string            `json:"defaultSquareOffTime" bson:"defaultSquareOffTime"`
        AutoSquareOff        bool              `json:"autoSquareOff" bson:"autoSquareOff"`
        MaxDailyLoss         float64           `json:"maxDailyLoss" bson:"maxDailyLoss"`
        MaxPositionSize      int               `json:"maxPositionSize" bson:"maxPositionSize"`
        MaxOrdersPerMinute   int               `json:"maxOrdersPerMinute" bson:"maxOrdersPerMinute"`
        CircuitBreaker       float64           `json:"circuitBreaker" bson:"circuitBreaker"`
        SlippageTolerance    float64           `json:"slippageTolerance" bson:"slippageTolerance"`
        PriceAdjustmentBuffer float64          `json:"priceAdjustmentBuffer" bson:"priceAdjustmentBuffer"`
        OrderPlacementDelay  int               `json:"orderPlacementDelay" bson:"orderPlacementDelay"`
        DefaultTrailingSettings map[string]float64 `json:"defaultTrailingSettings" bson:"defaultTrailingSettings"`
        NotificationSettings map[string]bool   `json:"notificationSettings" bson:"notificationSettings"`
        NotificationChannels map[string]bool   `json:"notificationChannels" bson:"notificationChannels"`
        DisplaySettings      map[string]interface{} `json:"displaySettings" bson:"displaySettings"`
        GridLayouts          map[string]interface{} `json:"gridLayouts" bson:"gridLayouts"`
        ColumnVisibility     map[string]bool   `json:"columnVisibility" bson:"columnVisibility"`
        DefaultTimeframes    map[string]string `json:"defaultTimeframes" bson:"defaultTimeframes"`
        DataRefreshRate      int               `json:"dataRefreshRate" bson:"dataRefreshRate"`
        GreeksPrecision      int               `json:"greeksPrecision" bson:"greeksPrecision"`
        PriceFormatting      string            `json:"priceFormatting" bson:"priceFormatting"`
        PnLFormatting        string            `json:"pnLFormatting" bson:"pnLFormatting"`
        FavoriteSymbols      []string          `json:"favoriteSymbols" bson:"favoriteSymbols"`
        RecentSymbols        []string          `json:"recentSymbols" bson:"recentSymbols"`
        CustomShortcuts      map[string]string `json:"customShortcuts" bson:"customShortcuts"`
        SessionTimeout       int               `json:"sessionTimeout" bson:"sessionTimeout"`
        CreatedAt            time.Time         `json:"createdAt" bson:"createdAt"`
        UpdatedAt            time.Time         `json:"updatedAt" bson:"updatedAt"`
}

// APIKey represents an API key for broker integration
type APIKey struct {
        ID          string    `json:"id" bson:"_id,omitempty"`
        UserID      string    `json:"userId" bson:"userId"`
        BrokerName  string    `json:"brokerName" bson:"brokerName"`
        KeyName     string    `json:"keyName" bson:"keyName"`
        APIKey      string    `json:"-" bson:"apiKey"`
        APISecret   string    `json:"-" bson:"apiSecret"`
        AccessToken string    `json:"-" bson:"accessToken,omitempty"`
        IsActive    bool      `json:"isActive" bson:"isActive"`
        ExpiresAt   time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
        LastUsed    time.Time `json:"lastUsed,omitempty" bson:"lastUsed,omitempty"`
        Permissions []string  `json:"permissions" bson:"permissions"`
        IPWhitelist []string  `json:"ipWhitelist,omitempty" bson:"ipWhitelist,omitempty"`
        CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
        UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserFilter represents filters for querying users
type UserFilter struct {
        Username  string    `json:"username,omitempty"`
        Email     string    `json:"email,omitempty"`
        Role      UserRole  `json:"role,omitempty"`
        UserType  UserType  `json:"userType,omitempty"`
        Active    *bool     `json:"active,omitempty"`
        FromDate  time.Time `json:"fromDate,omitempty"`
        ToDate    time.Time `json:"toDate,omitempty"`
}

// Validate validates the user data
func (u *User) Validate() error {
        // Check required fields
        if u.Username == "" {
                return errors.New("username is required")
        }
        if u.Email == "" {
                return errors.New("email is required")
        }
        if u.PasswordHash == "" {
                return errors.New("password hash is required")
        }
        if u.FirstName == "" {
                return errors.New("first name is required")
        }
        if u.LastName == "" {
                return errors.New("last name is required")
        }

        // Validate username format (alphanumeric and underscore only, 3-30 chars)
        usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
        if !usernameRegex.MatchString(u.Username) {
                return errors.New("username must be 3-30 characters and contain only letters, numbers, and underscores")
        }

        // Validate email format
        emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
        if !emailRegex.MatchString(u.Email) {
                return errors.New("invalid email format")
        }

        // Validate phone number format if provided
        if u.Phone != "" {
                phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
                if !phoneRegex.MatchString(u.Phone) {
                        return errors.New("invalid phone number format")
                }
        }

        // Validate user role
        switch u.Role {
        case UserRoleAdmin, UserRoleTrader, UserRoleViewer:
                // Valid roles
        default:
                return errors.New("invalid user role")
        }

        // Validate user type
        switch u.UserType {
        case UserTypeStandard, UserTypeAdmin, UserTypeSIM:
                // Valid user types
        default:
                return errors.New("invalid user type")
        }

        // Validate failed login count
        if u.FailedLoginCount < 0 {
                return errors.New("failed login count cannot be negative")
        }

        return nil
}

// ValidateUserPreferences validates the user preferences data
func (p *UserPreferences) Validate() error {
        // Check required fields
        if p.UserID == "" {
                return errors.New("user ID is required")
        }

        // Validate theme
        if p.Theme != "light" && p.Theme != "dark" && p.Theme != "system" {
                return errors.New("invalid theme value")
        }

        // Validate default environment
        switch p.DefaultEnvironment {
        case EnvironmentLive, EnvironmentSIM:
                // Valid environments
        default:
                return errors.New("invalid default environment")
        }

        // Validate default product type
        switch p.DefaultProductType {
        case ProductTypeMIS, ProductTypeNRML, ProductTypeCNC:
                // Valid product types
        default:
                return errors.New("invalid default product type")
        }

        // Validate default order type
        switch p.DefaultOrderType {
        case OrderTypeMarket, OrderTypeLimit, OrderTypeSLLimit:
                // Valid order types
        default:
                return errors.New("invalid default order type")
        }

        // Validate default quantity
        if p.DefaultQuantity <= 0 {
                return errors.New("default quantity must be greater than zero")
        }

        // Validate default square off time format (HH:MM:SS)
        timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)
        if !timeRegex.MatchString(p.DefaultSquareOffTime) {
                return errors.New("invalid default square off time format (use HH:MM:SS)")
        }

        // Validate max daily loss
        if p.MaxDailyLoss < 0 {
                return errors.New("max daily loss cannot be negative")
        }

        // Validate max position size
        if p.MaxPositionSize <= 0 {
                return errors.New("max position size must be greater than zero")
        }

        // Validate max orders per minute
        if p.MaxOrdersPerMinute <= 0 {
                return errors.New("max orders per minute must be greater than zero")
        }

        // Validate circuit breaker
        if p.CircuitBreaker < 0 {
                return errors.New("circuit breaker cannot be negative")
        }

        // Validate slippage tolerance
        if p.SlippageTolerance < 0 {
                return errors.New("slippage tolerance cannot be negative")
        }

        // Validate price adjustment buffer
        if p.PriceAdjustmentBuffer < 0 {
                return errors.New("price adjustment buffer cannot be negative")
        }

        // Validate order placement delay
        if p.OrderPlacementDelay < 0 {
                return errors.New("order placement delay cannot be negative")
        }

        // Validate data refresh rate
        if p.DataRefreshRate <= 0 {
                return errors.New("data refresh rate must be greater than zero")
        }

        // Validate Greeks precision
        if p.GreeksPrecision < 0 || p.GreeksPrecision > 10 {
                return errors.New("Greeks precision must be between 0 and 10")
        }

        // Validate session timeout
        if p.SessionTimeout <= 0 {
                return errors.New("session timeout must be greater than zero")
        }

        return nil
}
