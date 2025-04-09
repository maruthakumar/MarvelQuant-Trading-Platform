# MarvelQuant Trading Platform - Repository Structure Documentation

This document provides a detailed overview of the MarvelQuant Trading Platform repository structure after the integration of frontend v10.3.3 and backend v10.2.0 into version v10.3.4.

## Directory Structure

```
MarvelQuant_v10.3.4/
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Continuous Integration workflow
│       └── cd.yml                 # Continuous Deployment workflow
├── backend/
│   ├── cmd/                       # Command-line applications
│   ├── internal/                  # Private application and library code
│   ├── pkg/                       # Public API packages
│   └── tests/                     # Backend test suites
├── config/                        # Configuration files
├── docs/                          # Documentation files
│   ├── paper_trading/             # Paper trading documentation
│   ├── testing/                   # Testing documentation
│   └── user_guide/                # User guides
├── frontend/
│   ├── build/                     # Production build output (not tracked in git)
│   ├── public/                    # Static assets and HTML template
│   │   └── assets/                # Images and other static assets
│   ├── src/                       # Source code for the frontend
│   │   ├── components/            # React components
│   │   │   ├── auth/              # Authentication components
│   │   │   ├── common/            # Common UI components
│   │   │   ├── dashboard/         # Dashboard components
│   │   │   ├── multileg/          # Multi-leg trading components
│   │   │   ├── orderbook/         # Order book components
│   │   │   ├── positions/         # Positions management components
│   │   │   ├── strategies/        # Strategy components
│   │   │   └── user/              # User settings components
│   │   ├── db/                    # Database connection code
│   │   ├── routes/                # API routes
│   │   └── services/              # Service layer for API calls
│   ├── .env                       # Development environment variables
│   ├── .env.production            # Production environment variables
│   ├── package.json               # Frontend dependencies and scripts
│   └── server.js                  # Express server for frontend
├── scripts/                       # Utility scripts
├── tests/                         # End-to-end and integration tests
├── .gitignore                     # Git ignore rules
└── README.md                      # Main repository documentation
```

## Key Components

### Frontend

The frontend is a React application with the following key features:
- Modern React with functional components and hooks
- Material UI for consistent styling
- React Router for navigation
- WebSocket integration for real-time data
- Modular component structure for maintainability

### Backend

The backend is written in Go with the following architecture:
- Clean architecture with separation of concerns
- RESTful API endpoints
- WebSocket server for real-time data
- Database integration with MySQL
- Comprehensive test coverage

### CI/CD Workflows

The repository includes GitHub Actions workflows for:
- Continuous Integration (CI): Building and testing both frontend and backend
- Continuous Deployment (CD): Deploying frontend to AWS S3 and backend to AWS infrastructure

### Version Control

- All version references have been updated to v10.3.4
- The repository is initialized with Git
- Remote origin is set to https://github.com/maruthakumar/MarvelQuant-Trading-Platform

## Integration Changes

The following key changes were made during the integration process:
1. Updated frontend UI components from v10.3.3
2. Integrated backend code from v10.2.0
3. Resolved dependencies and conflicts
4. Updated version references to v10.3.4
5. Added GitHub workflows for CI/CD
6. Created comprehensive documentation

## Next Steps

To complete the setup:
1. Push the repository to GitHub
2. Set up GitHub secrets for CI/CD workflows
3. Configure AWS credentials for deployment
4. Set up branch protection rules
5. Configure code owners for review requirements
