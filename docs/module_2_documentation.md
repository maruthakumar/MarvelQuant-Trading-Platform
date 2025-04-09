# Module 2: API Endpoints for Orderbook Documentation

## Overview
This document provides detailed information about the API endpoints implemented for order management in the trading platform. The implementation follows a layered architecture with handlers, services, and repositories.

## Architecture

### Layered Design
The API implementation follows a clean, layered architecture:

1. **Handlers Layer**: Responsible for HTTP request/response handling, parameter parsing, and input validation
2. **Services Layer**: Contains business logic, validation, and orchestration of operations
3. **Repositories Layer**: Handles data persistence and retrieval from the database

This separation of concerns ensures maintainability, testability, and scalability of the codebase.

## API Endpoints

### Order Management Endpoints

#### Create Order
- **Endpoint**: `POST /api/orders`
- **Description**: Creates a new order in the system
- **Request Body**: Order object with required fields
- **Response**: Created order with assigned ID and status
- **Status Codes**:
  - 201: Order created successfully
  - 400: Invalid request payload or validation error
  - 500: Internal server error

#### Get Order by ID
- **Endpoint**: `GET /api/orders/{id}`
- **Description**: Retrieves a specific order by its ID
- **URL Parameters**: `id` - The unique identifier of the order
- **Response**: Order object with complete details
- **Status Codes**:
  - 200: Order retrieved successfully
  - 404: Order not found
  - 500: Internal server error

#### Get Orders (with filtering)
- **Endpoint**: `GET /api/orders`
- **Description**: Retrieves a list of orders with optional filtering and pagination
- **Query Parameters**:
  - `userId`: Filter by user ID
  - `symbol`: Filter by symbol
  - `status`: Filter by order status
  - `direction`: Filter by order direction
  - `productType`: Filter by product type
  - `instrumentType`: Filter by instrument type
  - `portfolioId`: Filter by portfolio ID
  - `strategyId`: Filter by strategy ID
  - `fromDate`: Filter by creation date (start)
  - `toDate`: Filter by creation date (end)
  - `page`: Page number for pagination (default: 1)
  - `limit`: Number of items per page (default: 50)
- **Response**: List of orders with pagination metadata
- **Status Codes**:
  - 200: Orders retrieved successfully
  - 500: Internal server error

#### Update Order
- **Endpoint**: `PUT /api/orders/{id}`
- **Description**: Updates an existing order
- **URL Parameters**: `id` - The unique identifier of the order
- **Request Body**: Updated order object
- **Response**: Updated order with complete details
- **Status Codes**:
  - 200: Order updated successfully
  - 400: Invalid request payload or validation error
  - 404: Order not found
  - 500: Internal server error

#### Cancel Order
- **Endpoint**: `POST /api/orders/{id}/cancel`
- **Description**: Cancels an existing order
- **URL Parameters**: `id` - The unique identifier of the order
- **Response**: Success message
- **Status Codes**:
  - 200: Order cancelled successfully
  - 400: Order cannot be cancelled (e.g., already executed)
  - 404: Order not found
  - 500: Internal server error

### User-Specific Order Endpoints

#### Get Orders by User
- **Endpoint**: `GET /api/users/{userId}/orders`
- **Description**: Retrieves all orders for a specific user
- **URL Parameters**: `userId` - The unique identifier of the user
- **Query Parameters**: Same filtering and pagination options as the main Get Orders endpoint
- **Response**: List of orders with pagination metadata
- **Status Codes**:
  - 200: Orders retrieved successfully
  - 500: Internal server error

### Strategy-Specific Order Endpoints

#### Get Orders by Strategy
- **Endpoint**: `GET /api/strategies/{strategyId}/orders`
- **Description**: Retrieves all orders for a specific strategy
- **URL Parameters**: `strategyId` - The unique identifier of the strategy
- **Query Parameters**: Same filtering and pagination options as the main Get Orders endpoint
- **Response**: List of orders with pagination metadata
- **Status Codes**:
  - 200: Orders retrieved successfully
  - 500: Internal server error

### Portfolio-Specific Order Endpoints

#### Get Orders by Portfolio
- **Endpoint**: `GET /api/portfolios/{portfolioId}/orders`
- **Description**: Retrieves all orders for a specific portfolio
- **URL Parameters**: `portfolioId` - The unique identifier of the portfolio
- **Query Parameters**: Same filtering and pagination options as the main Get Orders endpoint
- **Response**: List of orders with pagination metadata
- **Status Codes**:
  - 200: Orders retrieved successfully
  - 500: Internal server error

## Implementation Details

### Handlers
The handlers layer is responsible for:
- Parsing HTTP requests
- Validating input data
- Calling the appropriate service methods
- Formatting and returning HTTP responses

Key features:
- Comprehensive error handling
- Proper HTTP status codes
- JSON response formatting
- Pagination support

### Services
The services layer contains the business logic:
- Order validation
- Status management
- Business rule enforcement
- Orchestration of repository operations

Key features:
- Input validation
- Business rule enforcement (e.g., only pending orders can be cancelled)
- Pagination parameter normalization
- Error handling and propagation

### Repositories
The repositories layer handles data persistence:
- MongoDB integration
- CRUD operations
- Query building for filtering
- Pagination implementation

Key features:
- MongoDB integration with proper error handling
- Filter construction based on query parameters
- Pagination with offset and limit
- Sorting options

## Data Models
The API endpoints use the Order model defined in Module 1, with all its validation rules and business logic.

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

## Pagination
All list endpoints support pagination:
- Default page size is 50 items
- Maximum page size is 100 items
- Response includes metadata:
  - total: Total number of items matching the filter
  - page: Current page number
  - limit: Number of items per page
  - totalPages: Total number of pages
  - hasNextPage: Boolean indicating if there are more pages

## Usage Examples

### Creating an Order
```
POST /api/orders
Content-Type: application/json

{
  "userId": "user123",
  "symbol": "NIFTY",
  "exchange": "NSE",
  "orderType": "LIMIT",
  "direction": "BUY",
  "quantity": 10,
  "price": 500.50,
  "productType": "MIS",
  "instrumentType": "OPTION",
  "optionType": "CE",
  "strikePrice": 18000,
  "expiry": "2025-05-03T00:00:00Z"
}
```

### Retrieving Orders with Filtering
```
GET /api/orders?userId=user123&status=PENDING&page=1&limit=20
```

### Cancelling an Order
```
POST /api/orders/order123/cancel
```

## Future Enhancements
Potential future enhancements for the API endpoints include:
- Authentication and authorization middleware
- Rate limiting
- Caching for frequently accessed data
- WebSocket integration for real-time updates
- Bulk operations for orders
- Advanced filtering options
- Sorting options for list endpoints

## Implementation Notes
- All endpoints follow RESTful principles
- JSON is used for request and response bodies
- Proper HTTP methods are used for different operations
- Consistent URL structure and naming conventions
- Comprehensive error messages for better client experience
