# Real-Time Stock Alert System

A scalable microservices-based - Git** - Version control with proper gitignore configuration

## Getting Startedck monitoring and alert system built with Go, Apache Kafka, and PostgreSQL. The system tracks stock prices, analyzes market data, and sends real-time alerts to users based on their portfolio preferences.

## Project Overview

This project implements a distributed stock alert system that:
- Monitors stock prices in real-time
- Analyzes portfolio performance and market trends
- Sends automated alerts based on user-defined criteria
- Provides scalable microservices architecture for high availability

## Architecture

The system follows a microservices architecture with event-driven communication:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   API       │    │   Alert     │    │ Persistence │
│  Service    │    │  Consumer   │    │  Consumer   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                   ┌─────────────┐
                   │   Apache    │
                   │   Kafka     │
                   └─────────────┘
                           │
                   ┌─────────────┐
                   │ Analytics   │
                   │ Consumer    │
                   └─────────────┘
```

## Tech Stack

### Backend
- **Go 1.25.1** - Core application language
- **Gin Framework** - HTTP web framework for REST APIs
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database for data persistence

### Messaging & Communication
- **Apache Kafka** - Event streaming platform for real-time data processing
- **Kafka Producers** - For publishing stock events and user actions
- **Kafka Consumers** - For processing alerts, persistence, and analytics

### Microservices
- **API Service** - REST API for user interactions and portfolio management
- **Alert Consumer** - Processes and sends stock price alerts
- **Persistence Consumer** - Handles data storage and retrieval operations
- **Analytics Consumer** - Performs market analysis and trend calculations

### Infrastructure & DevOps
- **Docker** - Containerization of all microservices
- **Docker Compose** - Local development environment orchestration
- **GitHub Actions** - CI/CD pipeline for automated testing and deployment
- **Docker Hub** - Container registry for image distribution

### Development & Testing
- **Go Testing** - Comprehensive unit test suite (28+ tests)
- **Environment Configuration** - Secure environment variable management
- **Git** - Version control with proper gitignore configuration

## 🚀 Getting Started

### Prerequisites
- Go 1.25.1 or higher
- Docker and Docker Compose
- PostgreSQL
- Apache Kafka

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Anubhav200311/Real-Time-Stock-Alert.git
   cd Real-Time-Stock-Alert
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration:
   # - ALPHA_VANTAGE_API_KEY=your_api_key
   # - Database credentials
   # - Kafka broker configuration
   ```

3. **Start infrastructure services --> Basically the Entire Project**
   ```bash
   docker-compose up --build -d 
   ```

## Testing

The project includes comprehensive test coverage:

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./models -v     # Data model tests
go test ./services -v   # Business logic tests  
go test ./routes -v     # API endpoint tests
go test ./db -v         # Database tests
```

## API Endpoints

### User Management
- `POST /users` - Create new user
- `GET /users` - List all users
- `GET /users/:id/portfolio` - Get user portfolio

### Portfolio Management
- `POST /portfolio` - Create/update portfolio
- `GET /portfolio/:id` - Get portfolio details

### Stock Operations
- `GET /stocks/:symbol` - Get stock information
- `POST /alerts` - Create price alerts

## CI/CD Pipeline

The project uses GitHub Actions for automated CI/CD:

### Build Process
1. **Code Checkout** - Latest code from main branch
2. **Go Setup** - Configure Go 1.25.1 environment
3. **Dependency Caching** - Cache Go modules for faster builds
4. **Testing** - Run comprehensive test suite
5. **Build Verification** - Compile all microservices
6. **Docker Build** - Create optimized container images
7. **Registry Push** - Publish to Docker Hub

### Deployment
- Automated staging deployment on main branch pushes
- Production-ready Docker images with multi-stage builds
- Environment-specific configuration management

## Future Roadmap

### Short-term Improvements
- **WebSocket Integration** - Replace polling with real-time WebSocket connections for live price feeds
- **Enhanced Error Handling** - Implement circuit breakers and retry mechanisms
- **Monitoring & Observability** - Add Prometheus metrics and distributed tracing

### Medium-term Goals
- **AI-Powered Analytics** - Integrate machine learning models for:
  - Predictive price analysis
  - Market sentiment analysis
  - Automated trading strategies
  - Risk assessment algorithms

### Long-term Vision
- **Kubernetes Orchestration** - Migrate to Kubernetes for:
  - Auto-scaling based on market volatility
  - Rolling deployments with zero downtime
  - Advanced service mesh with Istio
  - Multi-region deployment capabilities

- **Advanced Features**
  - Real-time portfolio optimization
  - Social trading and copy-trading features
  - Integration with multiple stock exchanges
  - Mobile application with push notifications
  - Advanced charting and technical analysis tools

### Infrastructure Enhancements
- **Message Queue Optimization** - Implement Kafka Streams for complex event processing
- **Database Scaling** - Add read replicas and implement database sharding
- **Security Hardening** - Implement OAuth 2.0, rate limiting, and API encryption
- **Performance Optimization** - Redis caching layer and CDN integration

## Performance Metrics

System design targets and expectations:
- Designed to handle high concurrent user loads through microservices architecture
- Event-driven processing with Kafka enables scalable message throughput
- Asynchronous alert processing for responsive user experience
- Containerized deployment ready for horizontal scaling

*Note: Performance benchmarks will be established during production deployment and load testing phases.*

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Alpha Vantage API for real-time stock data
- Apache Kafka community for excellent event streaming platform
- Go community for robust ecosystem and libraries

---

Built with Go, Kafka, and Cloud-Native technologies
