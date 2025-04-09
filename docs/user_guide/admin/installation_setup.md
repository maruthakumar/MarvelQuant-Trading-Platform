# Installation and Setup

## Introduction

This guide provides comprehensive instructions for installing and setting up the Trading Platform in various environments. It covers system requirements, installation procedures, configuration options, and initial setup steps. This document is intended for system administrators responsible for deploying and maintaining the platform.

## System Requirements

### Hardware Requirements

The Trading Platform can be deployed on various hardware configurations depending on your scale and performance needs. Below are the recommended specifications for different deployment scenarios:

#### Development/Testing Environment
- **CPU**: 4+ cores (Intel Core i5/i7 or AMD Ryzen 5/7)
- **RAM**: 16GB minimum, 32GB recommended
- **Storage**: 100GB SSD
- **Network**: 100Mbps+ internet connection

#### Small Production Environment (up to 100 concurrent users)
- **CPU**: 8+ cores (Intel Xeon or AMD EPYC)
- **RAM**: 32GB minimum, 64GB recommended
- **Storage**: 500GB SSD (RAID configuration recommended)
- **Network**: 1Gbps+ internet connection with low latency

#### Medium Production Environment (100-500 concurrent users)
- **CPU**: 16+ cores across multiple servers
- **RAM**: 128GB+ across the cluster
- **Storage**: 1TB+ SSD in RAID configuration
- **Network**: Redundant 1Gbps+ connections

#### Large Production Environment (500+ concurrent users)
- **CPU**: 32+ cores across multiple servers
- **RAM**: 256GB+ across the cluster
- **Storage**: 2TB+ SSD in RAID configuration with separate database servers
- **Network**: Redundant 10Gbps+ connections

### Software Requirements

#### Operating System
- **Linux**: Ubuntu 20.04 LTS or newer, CentOS 8+, Red Hat Enterprise Linux 8+
- **Windows**: Windows Server 2019 or newer (less recommended)
- **macOS**: macOS 11 (Big Sur) or newer (development only)

#### Database
- **PostgreSQL**: Version 13.0 or newer
- **TimescaleDB**: Version 2.0 or newer (for time-series data)
- **Redis**: Version 6.0 or newer

#### Message Queue
- **Kafka**: Version 2.8.0 or newer
- **ZooKeeper**: Version 3.6.0 or newer

#### Container Platform (for distributed deployment)
- **Docker**: Version 20.10 or newer
- **Kubernetes**: Version 1.20 or newer

#### Additional Software
- **Node.js**: Version 16.0 or newer
- **Python**: Version 3.9 or newer
- **Go**: Version 1.17 or newer
- **C++ Compiler**: GCC 10+ or Clang 12+
- **NGINX**: Version 1.20 or newer (for API Gateway)

### Network Requirements

- **Firewall Configuration**: Specific ports must be open for communication between components
- **Load Balancer**: Required for distributed deployments
- **SSL Certificates**: Required for secure communication
- **DNS Configuration**: Proper DNS setup for service discovery
- **Latency**: Low-latency connections between components (<5ms recommended)

## Pre-Installation Checklist

Before beginning the installation process, ensure you have completed the following preparations:

1. **System Access**
   - Root or sudo access to all servers
   - SSH access configured
   - Necessary firewall ports opened

2. **Software Prerequisites**
   - Base operating system installed and updated
   - Package managers configured (apt, yum, etc.)
   - Docker and Docker Compose installed (if using containerized deployment)
   - Kubernetes cluster configured (if using Kubernetes deployment)

3. **Database Preparation**
   - PostgreSQL installed and secured
   - Database user created with appropriate permissions
   - TimescaleDB extension installed
   - Redis server installed and secured

4. **Network Configuration**
   - DNS entries created for all services
   - SSL certificates obtained
   - Load balancer configured (if applicable)
   - Network security groups or firewall rules established

5. **Resource Planning**
   - Storage volumes provisioned
   - Backup strategy defined
   - Monitoring solution prepared
   - Logging infrastructure ready

## Installation Methods

The Trading Platform supports several installation methods to accommodate different deployment scenarios and preferences.

### Containerized Deployment (Recommended)

The containerized deployment uses Docker and Docker Compose to simplify installation and ensure consistency across environments.

#### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+
- Git

#### Installation Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/tradingplatform/trading-platform.git
   cd trading-platform
   ```

2. **Configure Environment Variables**
   ```bash
   cp .env.example .env
   # Edit .env file with your specific configuration
   nano .env
   ```

3. **Build and Start the Containers**
   ```bash
   docker-compose build
   docker-compose up -d
   ```

4. **Initialize the Database**
   ```bash
   docker-compose exec backend ./init-db.sh
   ```

5. **Verify Installation**
   ```bash
   docker-compose ps
   # All services should be in the "Up" state
   ```

6. **Access the Platform**
   - Web UI: https://your-server-address
   - API: https://your-server-address/api
   - Admin Panel: https://your-server-address/admin

### Kubernetes Deployment

For production environments with high availability and scalability requirements, Kubernetes deployment is recommended.

#### Prerequisites
- Kubernetes cluster 1.20+
- kubectl configured
- Helm 3.0+
- Persistent volume provisioner

#### Installation Steps

1. **Add the Trading Platform Helm Repository**
   ```bash
   helm repo add tradingplatform https://charts.tradingplatform.example.com
   helm repo update
   ```

2. **Create Configuration Values File**
   ```bash
   # Create a values.yaml file with your specific configuration
   nano values.yaml
   ```

3. **Install the Helm Chart**
   ```bash
   helm install trading-platform tradingplatform/trading-platform -f values.yaml
   ```

4. **Verify Deployment**
   ```bash
   kubectl get pods
   # All pods should be in the "Running" state
   ```

5. **Configure Ingress**
   ```bash
   kubectl apply -f ingress.yaml
   ```

6. **Access the Platform**
   - Access the platform using the Ingress URL configured in your values.yaml

### Manual Installation

For environments where containers cannot be used, manual installation is available.

#### Prerequisites
- All required software installed (PostgreSQL, Redis, Node.js, etc.)
- Git

#### Installation Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/tradingplatform/trading-platform.git
   cd trading-platform
   ```

2. **Install Backend Dependencies**
   ```bash
   cd backend
   go mod download
   cd ..
   ```

3. **Install Frontend Dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Build the C++ Execution Engine**
   ```bash
   cd cpp
   mkdir build && cd build
   cmake ..
   make
   cd ../..
   ```

5. **Configure the Platform**
   ```bash
   cp config/config.example.yaml config/config.yaml
   # Edit config.yaml with your specific configuration
   nano config/config.yaml
   ```

6. **Initialize the Database**
   ```bash
   cd backend
   ./scripts/init-db.sh
   cd ..
   ```

7. **Start the Services**
   ```bash
   # Start backend services
   cd backend
   ./scripts/start-services.sh
   cd ..
   
   # Start frontend
   cd frontend
   npm run build
   npm run start
   cd ..
   ```

8. **Configure NGINX**
   ```bash
   cp nginx/nginx.example.conf /etc/nginx/sites-available/trading-platform.conf
   # Edit the configuration file
   nano /etc/nginx/sites-available/trading-platform.conf
   
   # Enable the site
   ln -s /etc/nginx/sites-available/trading-platform.conf /etc/nginx/sites-enabled/
   
   # Test and reload NGINX
   nginx -t
   systemctl reload nginx
   ```

9. **Access the Platform**
   - Web UI: https://your-server-address
   - API: https://your-server-address/api
   - Admin Panel: https://your-server-address/admin

## Configuration

The Trading Platform offers extensive configuration options to adapt to different environments and requirements.

### Core Configuration

The core configuration is managed through environment variables or a configuration file, depending on the installation method.

#### Environment Variables

For containerized deployments, key configuration options are set in the `.env` file:

```
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=tradingplatform
DB_USER=tradinguser
DB_PASSWORD=securepassword

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=securepassword

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_PREFIX=tradingplatform

# Authentication
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRATION=86400
ENABLE_2FA=true

# API Configuration
API_RATE_LIMIT=120
API_TIMEOUT=30

# Execution Engine
EXECUTION_ENGINE_HOST=execution-engine
EXECUTION_ENGINE_PORT=8085
```

#### Configuration File

For manual installations, configuration is managed through `config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  name: tradingplatform
  user: tradinguser
  password: securepassword
  pool_size: 20
  timeout: 5

redis:
  host: localhost
  port: 6379
  password: securepassword
  db: 0

kafka:
  brokers:
    - localhost:9092
  topic_prefix: tradingplatform
  consumer_group: trading-platform

authentication:
  jwt_secret: your-jwt-secret-key
  jwt_expiration: 86400
  enable_2fa: true
  allowed_origins:
    - https://tradingplatform.example.com

api:
  host: 0.0.0.0
  port: 8080
  rate_limit: 120
  timeout: 30
  cors_enabled: true

execution_engine:
  host: localhost
  port: 8085
  max_connections: 100
  timeout: 5
```

### Database Configuration

The database configuration includes settings for PostgreSQL and TimescaleDB:

#### PostgreSQL Settings

```yaml
postgresql:
  max_connections: 200
  shared_buffers: 4GB
  effective_cache_size: 12GB
  maintenance_work_mem: 1GB
  checkpoint_completion_target: 0.9
  wal_buffers: 16MB
  default_statistics_target: 100
  random_page_cost: 1.1
  effective_io_concurrency: 200
  work_mem: 20MB
  min_wal_size: 1GB
  max_wal_size: 4GB
```

#### TimescaleDB Settings

```yaml
timescaledb:
  max_background_workers: 8
  chunk_time_interval: 86400000  # 1 day in milliseconds
  retention_policy:
    enabled: true
    interval: 90d  # 90 days retention
```

### Security Configuration

Security settings control authentication, authorization, and encryption:

```yaml
security:
  ssl:
    enabled: true
    cert_file: /path/to/cert.pem
    key_file: /path/to/key.pem
  
  authentication:
    password_policy:
      min_length: 8
      require_uppercase: true
      require_lowercase: true
      require_numbers: true
      require_special_chars: true
    
    lockout_policy:
      max_attempts: 5
      lockout_duration: 30m
    
    session:
      idle_timeout: 30m
      absolute_timeout: 12h
  
  authorization:
    rbac_enabled: true
    default_role: user
```

### Logging Configuration

Logging settings control the verbosity and destination of logs:

```yaml
logging:
  level: info  # debug, info, warn, error
  format: json  # json, text
  output: file  # file, stdout, both
  file_path: /var/log/trading-platform
  rotation:
    max_size: 100  # MB
    max_age: 7  # days
    max_backups: 10
    compress: true
```

### Performance Tuning

Performance settings optimize the platform for different workloads:

```yaml
performance:
  worker_pools:
    order_processor:
      min_workers: 10
      max_workers: 50
      queue_size: 1000
    
    market_data:
      min_workers: 5
      max_workers: 20
      queue_size: 500
  
  caching:
    market_data:
      enabled: true
      ttl: 60s
    
    user_data:
      enabled: true
      ttl: 300s
  
  rate_limiting:
    enabled: true
    strategy: token_bucket
    refill_rate: 100
    bucket_size: 200
```

## Initial Setup

After installation, several initial setup steps are required to prepare the platform for use.

### Administrator Account Creation

Create the initial administrator account:

```bash
# For containerized deployment
docker-compose exec backend ./create-admin.sh

# For Kubernetes deployment
kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./create-admin.sh

# For manual installation
cd backend
./scripts/create-admin.sh
```

Follow the prompts to set the administrator username, email, and password.

### System Configuration

1. **Log in to the Admin Panel**
   - Access https://your-server-address/admin
   - Log in with the administrator credentials created earlier

2. **Configure System Settings**
   - General Settings: Set platform name, contact information, etc.
   - Email Settings: Configure SMTP server for notifications
   - Integration Settings: Set up external service connections
   - Security Settings: Configure authentication policies
   - Trading Settings: Set trading hours, order types, etc.

3. **Create User Roles**
   - Navigate to User Management > Roles
   - Create roles with appropriate permissions
   - Assign roles to users

4. **Configure Market Data Sources**
   - Navigate to Market Data > Sources
   - Add and configure market data providers
   - Set up data refresh intervals
   - Configure historical data retention

5. **Set Up Broker Connections**
   - Navigate to Trading > Brokers
   - Add and configure broker connections
   - Set up API credentials
   - Test connectivity

### Data Import

Import initial data required for the platform:

1. **Import Instrument Data**
   ```bash
   # For containerized deployment
   docker-compose exec backend ./import-instruments.sh /path/to/instruments.csv

   # For Kubernetes deployment
   kubectl cp /path/to/instruments.csv $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}'):/tmp/
   kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./import-instruments.sh /tmp/instruments.csv

   # For manual installation
   cd backend
   ./scripts/import-instruments.sh /path/to/instruments.csv
   ```

2. **Import Historical Market Data**
   ```bash
   # For containerized deployment
   docker-compose exec backend ./import-market-data.sh /path/to/market-data.csv

   # For Kubernetes deployment
   kubectl cp /path/to/market-data.csv $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}'):/tmp/
   kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./import-market-data.sh /tmp/market-data.csv

   # For manual installation
   cd backend
   ./scripts/import-market-data.sh /path/to/market-data.csv
   ```

### Verification

Verify that the installation and setup are complete and functioning correctly:

1. **System Health Check**
   ```bash
   # For containerized deployment
   docker-compose exec backend ./health-check.sh

   # For Kubernetes deployment
   kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./health-check.sh

   # For manual installation
   cd backend
   ./scripts/health-check.sh
   ```

2. **Service Connectivity Test**
   ```bash
   # For containerized deployment
   docker-compose exec backend ./connectivity-test.sh

   # For Kubernetes deployment
   kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./connectivity-test.sh

   # For manual installation
   cd backend
   ./scripts/connectivity-test.sh
   ```

3. **User Interface Test**
   - Access the web UI at https://your-server-address
   - Verify that all pages load correctly
   - Test basic functionality (login, navigation, etc.)

4. **API Test**
   ```bash
   # Test API connectivity
   curl -X GET https://your-server-address/api/v1/health
   
   # Test authentication
   curl -X POST https://your-server-address/api/v1/auth/token \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"your-password"}'
   ```

## Upgrading

The Trading Platform is regularly updated with new features, improvements, and bug fixes. Follow these procedures to upgrade your installation.

### Backup Before Upgrading

Always create a backup before upgrading:

```bash
# For containerized deployment
docker-compose exec postgres pg_dump -U tradinguser tradingplatform > backup.sql

# For Kubernetes deployment
kubectl exec -it $(kubectl get pods -l app=postgres -o jsonpath='{.items[0].metadata.name}') -- pg_dump -U tradinguser tradingplatform > backup.sql

# For manual installation
pg_dump -U tradinguser tradingplatform > backup.sql
```

### Upgrading Containerized Deployment

1. **Pull the Latest Changes**
   ```bash
   cd trading-platform
   git pull
   ```

2. **Update Environment Variables**
   ```bash
   # Check for new environment variables
   diff .env.example .env
   # Update .env file as needed
   ```

3. **Rebuild and Restart Containers**
   ```bash
   docker-compose down
   docker-compose build
   docker-compose up -d
   ```

4. **Run Database Migrations**
   ```bash
   docker-compose exec backend ./run-migrations.sh
   ```

5. **Verify Upgrade**
   ```bash
   docker-compose exec backend ./version-check.sh
   ```

### Upgrading Kubernetes Deployment

1. **Update Helm Repository**
   ```bash
   helm repo update
   ```

2. **Check for Changes in Values**
   ```bash
   helm show values tradingplatform/trading-platform > new-values.yaml
   diff values.yaml new-values.yaml
   # Update values.yaml as needed
   ```

3. **Upgrade the Helm Release**
   ```bash
   helm upgrade trading-platform tradingplatform/trading-platform -f values.yaml
   ```

4. **Verify Upgrade**
   ```bash
   kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./version-check.sh
   ```

### Upgrading Manual Installation

1. **Pull the Latest Changes**
   ```bash
   cd trading-platform
   git pull
   ```

2. **Update Dependencies**
   ```bash
   cd backend
   go mod download
   cd ../frontend
   npm install
   cd ..
   ```

3. **Rebuild Components**
   ```bash
   cd cpp
   mkdir -p build && cd build
   cmake ..
   make
   cd ../../frontend
   npm run build
   cd ..
   ```

4. **Run Database Migrations**
   ```bash
   cd backend
   ./scripts/run-migrations.sh
   cd ..
   ```

5. **Restart Services**
   ```bash
   cd backend
   ./scripts/restart-services.sh
   cd ../frontend
   npm run restart
   cd ..
   ```

6. **Verify Upgrade**
   ```bash
   cd backend
   ./scripts/version-check.sh
   cd ..
   ```

## Troubleshooting

### Common Installation Issues

#### Database Connection Errors

**Issue**: Services cannot connect to the database.

**Solutions**:
- Verify database credentials in configuration
- Check that the database server is running
- Ensure network connectivity between services and database
- Check database logs for authentication failures
- Verify that the database user has appropriate permissions

```bash
# Test database connection
psql -h <db_host> -p <db_port> -U <db_user> -d <db_name>
```

#### Permission Issues

**Issue**: Services fail to start due to permission errors.

**Solutions**:
- Check file permissions on configuration files and directories
- Ensure the service user has appropriate permissions
- For containerized deployments, check volume mount permissions
- Verify that SSL certificate files are readable

```bash
# Fix permissions on configuration files
chmod 644 config/*.yaml
chmod 600 config/*_key.pem
```

#### Network Connectivity Issues

**Issue**: Services cannot communicate with each other.

**Solutions**:
- Verify that all required ports are open in firewalls
- Check network configuration in Docker or Kubernetes
- Ensure DNS resolution is working correctly
- Test connectivity between services

```bash
# Test connectivity between services
nc -zv <service_host> <service_port>
```

#### Memory or CPU Resource Issues

**Issue**: Services crash or perform poorly due to resource constraints.

**Solutions**:
- Check system resource usage (CPU, memory, disk I/O)
- Adjust resource limits in configuration
- Scale up hardware resources if necessary
- Optimize performance settings

```bash
# Check system resource usage
top
free -m
df -h
```

### Diagnostic Tools

The Trading Platform includes several diagnostic tools to help troubleshoot issues:

#### Health Check Tool

```bash
# For containerized deployment
docker-compose exec backend ./health-check.sh

# For Kubernetes deployment
kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./health-check.sh

# For manual installation
cd backend
./scripts/health-check.sh
```

The health check tool verifies:
- Database connectivity
- Redis connectivity
- Kafka connectivity
- Execution engine connectivity
- API functionality
- Authentication services
- File system access

#### Log Analysis Tool

```bash
# For containerized deployment
docker-compose exec backend ./analyze-logs.sh

# For Kubernetes deployment
kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./analyze-logs.sh

# For manual installation
cd backend
./scripts/analyze-logs.sh
```

The log analysis tool:
- Identifies common error patterns
- Summarizes error frequencies
- Suggests potential solutions
- Generates a diagnostic report

#### Configuration Validator

```bash
# For containerized deployment
docker-compose exec backend ./validate-config.sh

# For Kubernetes deployment
kubectl exec -it $(kubectl get pods -l app=backend -o jsonpath='{.items[0].metadata.name}') -- ./validate-config.sh

# For manual installation
cd backend
./scripts/validate-config.sh
```

The configuration validator:
- Checks configuration syntax
- Validates configuration values
- Identifies missing required settings
- Suggests optimal values based on system resources

### Getting Help

If you encounter issues that you cannot resolve:

1. **Check Documentation**
   - Review this installation guide
   - Check the troubleshooting section of the user guide
   - Review release notes for known issues

2. **Check Logs**
   - Review application logs for error messages
   - Check system logs for related issues
   - Analyze database logs for query problems

3. **Contact Support**
   - Email: support@tradingplatform.example.com
   - Support Portal: https://support.tradingplatform.example.com
   - Include detailed information about your environment and the issue
   - Attach relevant logs and configuration files (with sensitive information redacted)

## Next Steps

After completing the installation and initial setup, explore these related guides:

- [System Configuration](./system_configuration.md) - Detailed configuration options
- [Performance Monitoring](./performance_monitoring.md) - Monitoring and optimizing performance
- [Backup and Recovery](./backup_recovery.md) - Data backup and disaster recovery procedures
- [Security Management](./security_management.md) - Security best practices and configuration
- [User Management](./user_management.md) - Managing user accounts and permissions
