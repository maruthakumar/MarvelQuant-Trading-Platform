# Module 4: User Settings Backend Documentation

## Overview
This document provides detailed information about the User Settings Backend implemented in the trading platform. The implementation follows a clean, layered architecture with handlers, services, and repositories.

## Architecture

### Layered Design
The User Settings Backend follows a clean, layered architecture:

1. **Handlers Layer**: Responsible for HTTP request/response handling, parameter parsing, and input validation
2. **Services Layer**: Contains business logic, validation, and orchestration of operations
3. **Repositories Layer**: Handles data persistence and retrieval from the database

This separation of concerns ensures maintainability, testability, and scalability of the codebase.

## User Settings Models

The User Settings Backend includes several models to represent different aspects of user settings:

### UserSettings
Represents general user settings:
- Language preference
- Time zone
- Date and time formats
- Currency preference

### UserPreferences
Represents trading-specific preferences:
- Default order quantity
- Default product type
- Default exchange
- Confirmation dialog settings
- Default instrument type
- Default symbols

### UserTheme
Represents UI theme settings:
- Theme mode (light/dark)
- Primary color
- Secondary color
- Chart colors
- Font size

### UserLayout
Represents UI layout configurations:
- Layout name
- Layout type
- Default status
- Layout configuration (grid/flex settings, widgets, etc.)

### UserApiKey
Represents API keys for broker integrations:
- API key name
- API key and secret
- Broker name
- Active status
- Permissions
- Expiry date

### UserNotificationSettings
Represents notification preferences:
- Email notification settings
- Push notification settings
- Order execution alerts
- Price alerts
- Margin call alerts
- News alerts

## API Endpoints

### User Settings Endpoints

#### Get User Settings
- **Endpoint**: `GET /api/users/{userId}/settings`
- **Description**: Retrieves user settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: User settings object
- **Status Codes**:
  - 200: Settings retrieved successfully
  - 404: User settings not found
  - 500: Internal server error

#### Update User Settings
- **Endpoint**: `PUT /api/users/{userId}/settings`
- **Description**: Updates user settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: Updated user settings object
- **Response**: Updated user settings object
- **Status Codes**:
  - 200: Settings updated successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

### User Preferences Endpoints

#### Get User Preferences
- **Endpoint**: `GET /api/users/{userId}/preferences`
- **Description**: Retrieves user trading preferences
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: User preferences object
- **Status Codes**:
  - 200: Preferences retrieved successfully
  - 404: User preferences not found
  - 500: Internal server error

#### Update User Preferences
- **Endpoint**: `PUT /api/users/{userId}/preferences`
- **Description**: Updates user trading preferences
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: Updated user preferences object
- **Response**: Updated user preferences object
- **Status Codes**:
  - 200: Preferences updated successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

### User Theme Endpoints

#### Get User Theme
- **Endpoint**: `GET /api/users/{userId}/theme`
- **Description**: Retrieves user theme settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: User theme object
- **Status Codes**:
  - 200: Theme retrieved successfully
  - 404: User theme not found
  - 500: Internal server error

#### Update User Theme
- **Endpoint**: `PUT /api/users/{userId}/theme`
- **Description**: Updates user theme settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: Updated user theme object
- **Response**: Updated user theme object
- **Status Codes**:
  - 200: Theme updated successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

### User Layout Endpoints

#### Get User Layout
- **Endpoint**: `GET /api/users/{userId}/layouts/{layoutName}`
- **Description**: Retrieves a specific user layout
- **URL Parameters**: 
  - `userId` - The unique identifier of the user
  - `layoutName` - The name of the layout
- **Response**: User layout object
- **Status Codes**:
  - 200: Layout retrieved successfully
  - 404: User layout not found
  - 500: Internal server error

#### Get All User Layouts
- **Endpoint**: `GET /api/users/{userId}/layouts`
- **Description**: Retrieves all user layouts
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: Array of user layout objects
- **Status Codes**:
  - 200: Layouts retrieved successfully
  - 404: User not found
  - 500: Internal server error

#### Save User Layout
- **Endpoint**: `POST /api/users/{userId}/layouts`
- **Description**: Saves a user layout (creates new or updates existing)
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: User layout object
- **Response**: Saved user layout object
- **Status Codes**:
  - 200: Layout saved successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

#### Delete User Layout
- **Endpoint**: `DELETE /api/users/{userId}/layouts/{layoutName}`
- **Description**: Deletes a user layout
- **URL Parameters**: 
  - `userId` - The unique identifier of the user
  - `layoutName` - The name of the layout
- **Response**: Success message
- **Status Codes**:
  - 200: Layout deleted successfully
  - 400: Invalid request
  - 404: User layout not found
  - 500: Internal server error

### User API Key Endpoints

#### Get User API Keys
- **Endpoint**: `GET /api/users/{userId}/apikeys`
- **Description**: Retrieves all user API keys
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: Array of user API key objects (with masked secrets)
- **Status Codes**:
  - 200: API keys retrieved successfully
  - 404: User not found
  - 500: Internal server error

#### Add User API Key
- **Endpoint**: `POST /api/users/{userId}/apikeys`
- **Description**: Adds a new user API key
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: User API key object
- **Response**: Added user API key object (with masked secret)
- **Status Codes**:
  - 201: API key added successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

#### Update User API Key
- **Endpoint**: `PUT /api/users/{userId}/apikeys/{keyId}`
- **Description**: Updates a user API key
- **URL Parameters**: 
  - `userId` - The unique identifier of the user
  - `keyId` - The unique identifier of the API key
- **Request Body**: Updated user API key object
- **Response**: Updated user API key object (with masked secret)
- **Status Codes**:
  - 200: API key updated successfully
  - 400: Invalid request payload or validation error
  - 404: User API key not found
  - 500: Internal server error

#### Delete User API Key
- **Endpoint**: `DELETE /api/users/{userId}/apikeys/{keyId}`
- **Description**: Deletes a user API key
- **URL Parameters**: 
  - `userId` - The unique identifier of the user
  - `keyId` - The unique identifier of the API key
- **Response**: Success message
- **Status Codes**:
  - 200: API key deleted successfully
  - 400: Invalid request
  - 404: User API key not found
  - 500: Internal server error

### User Notification Settings Endpoints

#### Get User Notification Settings
- **Endpoint**: `GET /api/users/{userId}/notifications`
- **Description**: Retrieves user notification settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Response**: User notification settings object
- **Status Codes**:
  - 200: Notification settings retrieved successfully
  - 404: User notification settings not found
  - 500: Internal server error

#### Update User Notification Settings
- **Endpoint**: `PUT /api/users/{userId}/notifications`
- **Description**: Updates user notification settings
- **URL Parameters**: `userId` - The unique identifier of the user
- **Request Body**: Updated user notification settings object
- **Response**: Updated user notification settings object
- **Status Codes**:
  - 200: Notification settings updated successfully
  - 400: Invalid request payload or validation error
  - 404: User not found
  - 500: Internal server error

## Implementation Details

### User Service
The user service is responsible for:
- Retrieving and updating user settings
- Retrieving and updating user preferences
- Retrieving and updating user theme settings
- Managing user layouts (create, retrieve, update, delete)
- Managing user API keys (create, retrieve, update, delete)
- Retrieving and updating user notification settings

Key features:
- Comprehensive validation of all operations
- Default settings creation for new users
- Secure handling of API keys (masking secrets in responses)

### User Repository
The user repository handles data persistence:
- MongoDB integration
- CRUD operations for all user settings models
- Query building for filtering

Key features:
- MongoDB integration with proper error handling
- Separate collections for different types of settings
- Proper error handling for not found cases

### User Settings Handler
The user settings handler is responsible for:
- Parsing HTTP requests
- Validating input data
- Calling the appropriate service methods
- Formatting and returning HTTP responses

Key features:
- Comprehensive error handling
- Proper HTTP status codes
- JSON response formatting
- Security measures (masking sensitive information)

## Data Models

### UserSettings
```json
{
  "userId": "user123",
  "language": "en",
  "timeZone": "UTC",
  "dateFormat": "MM/DD/YYYY",
  "timeFormat": "HH:mm:ss",
  "currency": "USD",
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

### UserPreferences
```json
{
  "userId": "user123",
  "defaultOrderQuantity": 10,
  "defaultProductType": "MIS",
  "defaultExchange": "NSE",
  "showConfirmationDialog": true,
  "defaultInstrumentType": "OPTION",
  "defaultSymbols": ["NIFTY", "BANKNIFTY"],
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

### UserTheme
```json
{
  "userId": "user123",
  "themeMode": "dark",
  "primaryColor": "#1976d2",
  "secondaryColor": "#dc004e",
  "chartColors": ["#ff0000", "#00ff00", "#0000ff"],
  "fontSize": "medium",
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

### UserLayout
```json
{
  "userId": "user123",
  "name": "default",
  "type": "grid",
  "isDefault": true,
  "layout": {
    "columns": 12,
    "rows": 6,
    "widgets": [
      {
        "id": "widget1",
        "type": "chart",
        "x": 0,
        "y": 0,
        "w": 6,
        "h": 3
      },
      {
        "id": "widget2",
        "type": "orderbook",
        "x": 6,
        "y": 0,
        "w": 6,
        "h": 3
      }
    ]
  },
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

### UserApiKey
```json
{
  "id": "key123",
  "userId": "user123",
  "name": "Zerodha",
  "apiKey": "abc123",
  "apiSecret": "********", // Masked in responses
  "broker": "zerodha",
  "isActive": true,
  "permissions": ["trade", "data"],
  "expiresAt": "2026-04-03T09:30:00Z",
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

### UserNotificationSettings
```json
{
  "userId": "user123",
  "enableEmailNotifications": true,
  "enablePushNotifications": true,
  "orderExecutionAlerts": true,
  "priceAlerts": true,
  "marginCallAlerts": true,
  "newsAlerts": false,
  "createdAt": "2025-04-03T09:30:00Z",
  "updatedAt": "2025-04-03T09:30:00Z"
}
```

## Validation

Each model includes validation logic to ensure data integrity:

### UserSettings Validation
- User ID is required
- Language is required
- Time zone is required

### UserPreferences Validation
- User ID is required
- Default order quantity must be greater than zero
- Default exchange is required

### UserTheme Validation
- User ID is required
- Theme mode must be either 'light' or 'dark'

### UserLayout Validation
- User ID is required
- Layout name is required
- Layout type is required
- Layout configuration is required

### UserApiKey Validation
- User ID is required
- API key name is required
- API key is required
- API secret is required
- Broker is required

### UserNotificationSettings Validation
- User ID is required

## Default Settings

The service automatically creates default settings for new users:

### Default UserSettings
```json
{
  "language": "en",
  "timeZone": "UTC",
  "dateFormat": "MM/DD/YYYY",
  "timeFormat": "HH:mm:ss",
  "currency": "USD"
}
```

### Default UserPreferences
```json
{
  "defaultOrderQuantity": 1,
  "defaultProductType": "MIS",
  "defaultExchange": "NSE",
  "showConfirmationDialog": true
}
```

### Default UserTheme
```json
{
  "themeMode": "light",
  "primaryColor": "#1976d2",
  "secondaryColor": "#dc004e"
}
```

### Default UserNotificationSettings
```json
{
  "enableEmailNotifications": true,
  "enablePushNotifications": true,
  "orderExecutionAlerts": true,
  "priceAlerts": true,
  "marginCallAlerts": true,
  "newsAlerts": false
}
```

## Testing
Comprehensive unit tests have been implemented for all layers:
- Handler tests with mock services
- Service tests with mock repositories
- Repository tests with mock MongoDB interfaces

The tests cover:
- Happy path scenarios
- Error handling
- Edge cases
- Validation logic

## Error Handling
The API implements consistent error handling:
- Validation errors return 400 Bad Request with descriptive messages
- Not found errors return 404 Not Found
- Server errors return 500 Internal Server Error
- Custom error responses include an "error" field with a descriptive message

## Security Considerations
The User Settings Backend implements several security measures:
- API secrets are masked in responses
- User ID validation to prevent unauthorized access
- Input validation to prevent injection attacks
- Proper error handling to avoid information leakage

## Usage Examples

### Retrieving User Settings
```
GET /api/users/user123/settings
```

### Updating User Settings
```
PUT /api/users/user123/settings
Content-Type: application/json

{
  "userId": "user123",
  "language": "fr",
  "timeZone": "Europe/Paris",
  "dateFormat": "DD/MM/YYYY",
  "timeFormat": "HH:mm",
  "currency": "EUR"
}
```

### Retrieving User Theme
```
GET /api/users/user123/theme
```

### Saving a User Layout
```
POST /api/users/user123/layouts
Content-Type: application/json

{
  "userId": "user123",
  "name": "custom",
  "type": "grid",
  "isDefault": false,
  "layout": {
    "columns": 12,
    "rows": 6,
    "widgets": [
      {
        "id": "widget1",
        "type": "chart",
        "x": 0,
        "y": 0,
        "w": 12,
        "h": 6
      }
    ]
  }
}
```

### Adding a User API Key
```
POST /api/users/user123/apikeys
Content-Type: application/json

{
  "userId": "user123",
  "name": "Zerodha",
  "apiKey": "abc123",
  "apiSecret": "xyz789",
  "broker": "zerodha",
  "isActive": true,
  "permissions": ["trade", "data"]
}
```

### Updating User Notification Settings
```
PUT /api/users/user123/notifications
Content-Type: application/json

{
  "userId": "user123",
  "enableEmailNotifications": true,
  "enablePushNotifications": false,
  "orderExecutionAlerts": true,
  "priceAlerts": true,
  "marginCallAlerts": true,
  "newsAlerts": true
}
```

## Future Enhancements
Potential future enhancements for the User Settings Backend include:
- Multi-device synchronization
- Settings versioning and history
- Settings export/import functionality
- User preferences analytics
- Advanced theme customization options
- Layout templates and sharing
- API key usage analytics and monitoring
- Advanced notification rules and scheduling

## Implementation Notes
- All components follow a clean, layered architecture
- Proper separation of concerns between layers
- Comprehensive validation at all levels
- Thorough error handling
- MongoDB integration for data persistence
- RESTful API design principles
