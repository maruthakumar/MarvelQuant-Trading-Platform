# MarvelQuant Trading Platform v10.3.4

MarvelQuant Trading Platform is a comprehensive trading solution that combines advanced algorithmic trading capabilities with an intuitive user interface. This platform enables traders to execute complex trading strategies, monitor market data in real-time, and manage their portfolios efficiently.

## Repository Structure

The repository is organized into the following main directories:

- **frontend/**: Contains the React-based user interface code
  - `src/`: Source code for the frontend application
  - `public/`: Static assets and HTML template
  - `build/`: Production build output (not tracked in git)

- **backend/**: Contains the Go-based backend services
  - `cmd/`: Command-line applications
  - `internal/`: Private application and library code
  - `pkg/`: Public API packages
  - `tests/`: Backend test suites

- **docs/**: Documentation files
  - Release notes
  - API documentation
  - User guides

- **config/**: Configuration files for different environments

- **tests/**: End-to-end and integration tests

- **scripts/**: Utility scripts for development and deployment

## Technology Stack

### Frontend
- React.js
- Material UI
- React Router
- WebSockets for real-time data

### Backend
- Go
- RESTful API
- WebSocket server
- MySQL database

## Development Setup

### Prerequisites
- Node.js (v20.x or later)
- Go (v1.21 or later)
- MySQL

### Frontend Setup
```bash
cd frontend
npm install
npm start
```

### Backend Setup
```bash
cd backend
go mod download
go run cmd/server/main.go
```

## Deployment

The application is deployed using GitHub Actions workflows:

- CI workflow: Builds and tests the application on every push and pull request
- CD workflow: Deploys the frontend to AWS S3 and the backend to AWS infrastructure

## Version History

- v10.3.4: Current version with integrated frontend UI updates and backend optimizations
- v10.3.3: Frontend UI enhancements
- v10.2.0: Backend code improvements and API extensions

## License

Proprietary - All rights reserved
