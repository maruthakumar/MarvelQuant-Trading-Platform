# Trading Platform Project File Organization

## Directory Structure

The trading platform project is organized with the following directory structure:

```
trading-platform/
├── backend/
│   ├── cmd/              # Command-line applications and entry points
│   │   ├── server/       # Main API server
│   │   └── websocket/    # WebSocket server
│   ├── internal/         # Internal packages not meant for external use
│   │   ├── api/          # API routes and handlers
│   │   ├── auth/         # Authentication and authorization
│   │   ├── broker/       # Broker integration (XTS, etc.)
│   │   ├── core/         # Core business logic (orders, portfolios, strategies)
│   │   ├── database/     # Database connections and models
│   │   ├── messagequeue/ # Message queue services (Redis, RabbitMQ)
│   │   ├── middleware/   # HTTP middleware
│   │   └── websocket/    # WebSocket handling
│   └── tests/            # Backend tests
│       ├── broker/       # Broker integration tests
│       ├── database/     # Database performance tests
│       ├── integration/  # Integration tests
│       └── loadtest/     # Load testing
├── frontend/
│   ├── src/
│   │   ├── components/   # React components
│   │   │   ├── charts/   # Chart components (PriceChart, AnalyticsChart)
│   │   │   ├── dashboard/# Dashboard components
│   │   │   ├── layout/   # Layout components (Header, etc.)
│   │   │   ├── market/   # Market data components (MarketWatch, OptionChain)
│   │   │   ├── portfolio/# Portfolio components (PortfolioBuilder, PortfolioMonitor)
│   │   │   └── trading/  # Trading components (OrderForm, OrderMonitor)
│   │   ├── hooks/        # Custom React hooks
│   │   ├── services/     # API service clients
│   │   ├── store/        # Redux store
│   │   │   └── slices/   # Redux slices for state management
│   │   ├── tests/        # Frontend tests
│   │   └── utils/        # Utility functions
│   └── tsconfig.json     # TypeScript configuration
├── infrastructure/
│   └── docker/           # Docker configuration
│       └── docker-compose.yml # Docker Compose for local development
├── docs/                 # Documentation
├── plans/                # Implementation plans and roadmaps
├── code/                 # Additional code samples and utilities
├── previous_sessions/    # Archives from previous development sessions
└── python/               # Python utilities and scripts
```

## Key Components

### Backend

The backend is implemented in Go and follows a clean architecture approach:

1. **API Layer**: Handles HTTP requests and responses
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Manages data access

Key features:
- WebSocket server for real-time updates
- PostgreSQL/TimescaleDB for data storage
- Redis/RabbitMQ for message queuing
- Broker integration with XTS

### Frontend

The frontend is implemented in React with TypeScript and follows a component-based architecture:

1. **Components**: Reusable UI elements
2. **Services**: API clients for backend communication
3. **Store**: Redux for state management

Key features:
- Real-time data visualization with Recharts
- Material-UI for consistent styling
- Redux for state management
- WebSocket integration for real-time updates

### Infrastructure

The infrastructure is configured using Docker for consistent development and deployment:

1. **Docker Compose**: Local development environment
2. **PostgreSQL/TimescaleDB**: Database for time-series data
3. **Redis/RabbitMQ**: Message queuing

## File Organization Best Practices

1. **Modular Structure**: Each component is organized in its own directory
2. **Clear Separation**: Backend and frontend are clearly separated
3. **Consistent Naming**: Files and directories follow consistent naming conventions
4. **Test Organization**: Tests are organized alongside the code they test

## Deployment Configuration

The deployment configuration is designed for scalability and reliability:

1. **Docker Containers**: Each component runs in its own container
2. **Environment Variables**: Configuration is managed through environment variables
3. **Health Checks**: Each service includes health checks
4. **Logging**: Centralized logging for all components

## Development Workflow

1. Clone the repository
2. Install dependencies (without committing node_modules)
3. Run the development environment using Docker Compose
4. Make changes and run tests
5. Submit changes for review
