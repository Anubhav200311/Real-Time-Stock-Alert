# Stock Alerts Microservices

This project has been restructured into a microservices architecture with the following components:

## Architecture

```
stock-alerts/
├── api/                        # API service
│   └── main.go                 # HTTP API server
│
├── consumers/                  # Consumer microservices
│   ├── alert/                  # Alert processing
│   │   └── main.go
│   ├── persistence/            # Data persistence
│   │   └── main.go
│   └── analytics/              # Analytics aggregation
│       └── main.go
│
├── services/                   # Shared business logic
│   ├── producer.go             # Kafka producer utilities
│   └── stock.go                # Stock price fetching
│
├── models/                     # Database models
│   └── models.go
│
├── db/                         # Database connection
│   └── db.go
│
├── routes/                     # HTTP routes
│   └── routes.go
│
├── Dockerfile.api              # API service Docker
├── Dockerfile.alert            # Alert consumer Docker
├── Dockerfile.persistence      # Persistence consumer Docker
├── Dockerfile.analytics        # Analytics consumer Docker
├── docker-compose.yml          # Microservices orchestration
├── .env                        # Environment variables
├── go.mod
└── go.sum
```

## Microservices

### 1. API Service (`api/main.go`)
- **Port:** 8080
- **Responsibilities:**
  - Serves HTTP API endpoints
  - Handles user management, portfolio operations
  - Fetches stock prices and publishes to Kafka
  - Manages stock price thresholds

### 2. Alert Consumer (`consumers/alert/main.go`)
- **Kafka Group:** `stock-alerts-consumer`
- **Responsibilities:**
  - Consumes stock price events from Kafka
  - Checks price thresholds for user portfolios
  - Creates alerts when thresholds are exceeded
  - Stores alerts in database

### 3. Persistence Consumer (`consumers/persistence/main.go`)
- **Kafka Group:** `persistence-consumer-group`
- **Responsibilities:**
  - Consumes all stock price events
  - Stores historical price data in `stock_price_records` table
  - Provides audit trail for all price updates

### 4. Analytics Consumer (`consumers/analytics/main.go`)
- **Kafka Group:** `analytics-consumer-group`
- **Responsibilities:**
  - Consumes stock price events
  - Aggregates daily analytics (min, max, avg prices)
  - Tracks price change frequency
  - Stores analytics in `stock_analytics` table

## Database Tables

- `users` - User accounts
- `portfolios` - User portfolios (1:1 with users)
- `stocks` - Stocks in portfolios with thresholds
- `alerts` - Price threshold alerts
- `stock_price_records` - Historical price data
- `stock_analytics` - Daily aggregated analytics

## Running the Application

### Development (Local)
```bash
# Start infrastructure
docker compose up zookeeper kafka postgres kafka-ui -d

# Run services locally
go run api/main.go                    # API service
go run consumers/alert/main.go        # Alert consumer
go run consumers/persistence/main.go  # Persistence consumer
go run consumers/analytics/main.go    # Analytics consumer
```

### Production (Docker)
```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f api-service
docker compose logs -f alert-consumer
docker compose logs -f persistence-consumer
docker compose logs -f analytics-consumer
```

## API Endpoints

- `POST /users` - Create user
- `GET /users` - List users
- `POST /users/:id/portfolio` - Create portfolio
- `GET /users/:id/portfolio` - Get portfolio
- `POST /portfolio/:id/stocks` - Add stock to portfolio
- `GET /portfolio/:id/stocks` - List portfolio stocks
- `GET /users/:id/alerts` - Get user alerts

## Environment Variables

```bash
ALPHA_VANTAGE_API_KEY=your_api_key_here
```

## Monitoring

- **Kafka UI:** http://localhost:8081
- **API Service:** http://localhost:8080
- **Database:** localhost:5432

## Benefits of Microservices Architecture

1. **Scalability:** Each service can be scaled independently
2. **Isolation:** Failure in one service doesn't affect others
3. **Technology Independence:** Each service can use different tech stacks
4. **Team Autonomy:** Different teams can own different services
5. **Deployment Independence:** Services can be deployed separately
6. **Data Ownership:** Each service owns its data domain

## Migration from Monolith

The original `main.go` functionality has been split:
- HTTP API → `api/main.go`
- Alert processing → `consumers/alert/main.go` 
- Data persistence → `consumers/persistence/main.go`
- Analytics → `consumers/analytics/main.go` (new)

All services communicate via Kafka events and share the same database for simplicity.
