# Deployment Configuration for Trading Platform

## Overview

This document outlines the deployment configuration for the trading platform. The platform is designed to be deployed using Docker containers for consistency across development, staging, and production environments.

## Prerequisites

- Docker and Docker Compose
- Access to a PostgreSQL/TimescaleDB instance
- Access to Redis and RabbitMQ instances
- SSL certificates for secure communication

## Environment Variables

The following environment variables should be configured for deployment:

```
# Server Configuration
PORT=8080
WEBSOCKET_PORT=8081
API_BASE_URL=https://api.trading-platform.com
CORS_ALLOWED_ORIGINS=https://trading-platform.com

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=trading_platform
DB_USER=postgres
DB_PASSWORD=<secure-password>
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<secure-password>

# RabbitMQ Configuration
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=<secure-password>
RABBITMQ_VHOST=/

# Broker Configuration
XTS_API_URL=https://xts-api.com
XTS_API_KEY=<api-key>
XTS_API_SECRET=<api-secret>

# Authentication
JWT_SECRET=<secure-jwt-secret>
JWT_EXPIRATION=24h
```

## Docker Compose for Production

```yaml
version: '3.8'

services:
  backend:
    image: trading-platform/backend:latest
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=trading_platform
      - DB_USER=postgres
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - RABBITMQ_HOST=rabbitmq
      - RABBITMQ_PORT=5672
      - RABBITMQ_USER=guest
      - RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
      - rabbitmq
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  websocket:
    image: trading-platform/websocket:latest
    build:
      context: ./backend
      dockerfile: Dockerfile.websocket
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
    depends_on:
      - redis
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  frontend:
    image: trading-platform/frontend:latest
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        - REACT_APP_API_URL=https://api.trading-platform.com
        - REACT_APP_WEBSOCKET_URL=wss://ws.trading-platform.com
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./ssl:/etc/nginx/ssl
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: timescale/timescaledb:latest-pg14
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=trading_platform
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: always
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    restart: always
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      - RABBITMQ_DEFAULT_USER=guest
      - RABBITMQ_DEFAULT_PASS=${RABBITMQ_PASSWORD}
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    restart: always
    healthcheck:
      test: ["CMD", "rabbitmqctl", "status"]
      interval: 30s
      timeout: 10s
      retries: 5

volumes:
  postgres_data:
  redis_data:
  rabbitmq_data:
```

## Deployment Steps

1. **Build Docker Images**:
   ```bash
   docker-compose build
   ```

2. **Set Environment Variables**:
   Create a `.env` file with all required environment variables.

3. **Deploy Services**:
   ```bash
   docker-compose up -d
   ```

4. **Initialize Database**:
   ```bash
   docker-compose exec backend go run cmd/server/main.go --migrate
   ```

5. **Verify Deployment**:
   ```bash
   docker-compose ps
   ```

## Scaling Considerations

For production environments, consider the following scaling strategies:

1. **Horizontal Scaling**:
   - Deploy multiple instances of the backend and websocket services
   - Use a load balancer to distribute traffic

2. **Database Scaling**:
   - Implement read replicas for PostgreSQL
   - Use TimescaleDB's hypertable partitioning for efficient time-series data

3. **Caching Strategy**:
   - Implement Redis caching for frequently accessed data
   - Configure appropriate TTL values for cached data

## Monitoring and Logging

1. **Prometheus** for metrics collection
2. **Grafana** for visualization
3. **ELK Stack** (Elasticsearch, Logstash, Kibana) for centralized logging
4. **Healthchecks** for service availability monitoring

## Backup Strategy

1. **Database Backups**:
   - Daily full backups
   - Hourly incremental backups
   - Point-in-time recovery configuration

2. **Configuration Backups**:
   - Version control for all configuration files
   - Regular backups of environment variables

## Security Considerations

1. **Network Security**:
   - Use private networks for inter-service communication
   - Implement proper firewall rules

2. **Authentication and Authorization**:
   - Secure JWT implementation
   - Role-based access control

3. **Data Protection**:
   - Encrypt sensitive data at rest
   - Use SSL/TLS for all communications

## Disaster Recovery

1. **Failover Strategy**:
   - Configure service redundancy
   - Implement automatic failover

2. **Recovery Time Objective (RTO)**:
   - Target RTO: 15 minutes

3. **Recovery Point Objective (RPO)**:
   - Target RPO: 5 minutes
