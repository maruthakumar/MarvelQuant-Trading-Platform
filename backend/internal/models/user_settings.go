package models

import (
	"errors"
	"time"
)

// UserSettings represents user settings
type UserSettings struct {
	UserID    string    `json:"userId" bson:"userId"`
	Language  string    `json:"language" bson:"language"`
	TimeZone  string    `json:"timeZone" bson:"timeZone"`
	DateFormat string   `json:"dateFormat" bson:"dateFormat"`
	TimeFormat string   `json:"timeFormat" bson:"timeFormat"`
	Currency  string    `json:"currency" bson:"currency"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserPreferences represents user trading preferences
type UserPreferences struct {
	UserID                string      `json:"userId" bson:"userId"`
	DefaultOrderQuantity  int         `json:"defaultOrderQuantity" bson:"defaultOrderQuantity"`
	DefaultProductType    ProductType `json:"defaultProductType" bson:"defaultProductType"`
	DefaultExchange       string      `json:"defaultExchange" bson:"defaultExchange"`
	ShowConfirmationDialog bool        `json:"showConfirmationDialog" bson:"showConfirmationDialog"`
	DefaultInstrumentType InstrumentType `json:"defaultInstrumentType" bson:"defaultInstrumentType"`
	DefaultSymbols        []string    `json:"defaultSymbols" bson:"defaultSymbols"`
	CreatedAt             time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time   `json:"updatedAt" bson:"updatedAt"`
}

// UserTheme represents user theme settings
type UserTheme struct {
	UserID         string    `json:"userId" bson:"userId"`
	ThemeMode      string    `json:"themeMode" bson:"themeMode"`
	PrimaryColor   string    `json:"primaryColor" bson:"primaryColor"`
	SecondaryColor string    `json:"secondaryColor" bson:"secondaryColor"`
	ChartColors    []string  `json:"chartColors,omitempty" bson:"chartColors,omitempty"`
	FontSize       string    `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserLayout represents user layout settings
type UserLayout struct {
	UserID    string                 `json:"userId" bson:"userId"`
	Name      string                 `json:"name" bson:"name"`
	Type      string                 `json:"type" bson:"type"`
	IsDefault bool                   `json:"isDefault" bson:"isDefault"`
	Layout    map[string]interface{} `json:"layout" bson:"layout"`
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updatedAt"`
}

// UserApiKey represents a user API key
type UserApiKey struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	UserID      string    `json:"userId" bson:"userId"`
	Name        string    `json:"name" bson:"name"`
	ApiKey      string    `json:"apiKey" bson:"apiKey"`
	ApiSecret   string    `json:"apiSecret" bson:"apiSecret"`
	Broker      string    `json:"broker" bson:"broker"`
	IsActive    bool      `json:"isActive" bson:"isActive"`
	Permissions []string  `json:"permissions" bson:"permissions"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// UserNotificationSettings represents user notification settings
type UserNotificationSettings struct {
	UserID                 string    `json:"userId" bson:"userId"`
	EnableEmailNotifications bool      `json:"enableEmailNotifications" bson:"enableEmailNotifications"`
	EnablePushNotifications  bool      `json:"enablePushNotifications" bson:"enablePushNotifications"`
	OrderExecutionAlerts     bool      `json:"orderExecutionAlerts" bson:"orderExecutionAlerts"`
	PriceAlerts              bool      `json:"priceAlerts" bson:"priceAlerts"`
	MarginCallAlerts         bool      `json:"marginCallAlerts" bson:"marginCallAlerts"`
	NewsAlerts               bool      `json:"newsAlerts" bson:"newsAlerts"`
	CreatedAt              time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Validate validates the user settings
func (s *UserSettings) Validate() error {
	if s.UserID == "" {
		return errors.New("user ID is required")
	}
	if s.Language == "" {
		return errors.New("language is required")
	}
	if s.TimeZone == "" {
		return errors.New("time zone is required")
	}
	return nil
}

// Validate validates the user preferences
func (p *UserPreferences) Validate() error {
	if p.UserID == "" {
		return errors.New("user ID is required")
	}
	if p.DefaultOrderQuantity <= 0 {
		return errors.New("default order quantity must be greater than zero")
	}
	if p.DefaultExchange == "" {
		return errors.New("default exchange is required")
	}
	return nil
}

// Validate validates the user theme
func (t *UserTheme) Validate() error {
	if t.UserID == "" {
		return errors.New("user ID is required")
	}
	if t.ThemeMode != "light" && t.ThemeMode != "dark" {
		return errors.New("theme mode must be either 'light' or 'dark'")
	}
	return nil
}

// Validate validates the user layout
func (l *UserLayout) Validate() error {
	if l.UserID == "" {
		return errors.New("user ID is required")
	}
	if l.Name == "" {
		return errors.New("layout name is required")
	}
	if l.Type == "" {
		return errors.New("layout type is required")
	}
	if l.Layout == nil {
		return errors.New("layout configuration is required")
	}
	return nil
}

// Validate validates the user API key
func (k *UserApiKey) Validate() error {
	if k.UserID == "" {
		return errors.New("user ID is required")
	}
	if k.Name == "" {
		return errors.New("API key name is required")
	}
	if k.ApiKey == "" {
		return errors.New("API key is required")
	}
	if k.ApiSecret == "" {
		return errors.New("API secret is required")
	}
	if k.Broker == "" {
		return errors.New("broker is required")
	}
	return nil
}

// Validate validates the user notification settings
func (s *UserNotificationSettings) Validate() error {
	if s.UserID == "" {
		return errors.New("user ID is required")
	}
	return nil
}
