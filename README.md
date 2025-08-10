# byFood Library Management System

A comprehensive, production-ready library management system built with Clean Architecture principles, featuring a modern Go backend and React frontend.

## Features

### Backend (Go + Echo)
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **UUID Support**: All entities use UUIDs for better scalability and security
- **PostgreSQL**: Robust database with advanced querying and UUID support
- **sqlx Integration**: Enhanced SQL operations with named queries and struct scanning
- **Comprehensive Testing**: Unit tests, integration tests, and mocks with 56.2% coverage
- **Swagger Documentation**: Complete API documentation with interactive testing
- **Request Tracing**: UUID-based request tracking for observability
- **Error Handling**: Centralized error handling with structured logging
- **Security**: Rate limiting, API key auth, security headers, and CORS
- **Monitoring**: Prometheus metrics and health checks
- **Production Ready**: Docker, resource limits, and optimized configurations

### Frontend (React + Next.js + TypeScript)
- **Modern React**: React 18 with TypeScript for type safety
- **Real-time Search**: Instant book search and filtering
- **Responsive Design**: Tailwind CSS with mobile-first approach
- **State Management**: Context API with optimized performance
- **CRUD Operations**: Complete book management with modals
- **Data Persistence**: Seamless integration with backend API
- **Performance Optimized**: Memoization, lazy loading, and code splitting

### DevOps & Infrastructure
- **Docker Containerization**: Full stack containerized with persistent storage
- **Production Config**: Optimized Docker Compose for production deployment
- **Monitoring Stack**: Prometheus + Grafana for comprehensive monitoring
- **Reverse Proxy**: Nginx configuration for production load balancing
- **Health Checks**: Application and container health monitoring
- **Resource Management**: Memory and CPU limits for efficient scaling

## Architecture

### Clean Architecture Implementation
```
┌─────────────────────────────────────────────────────────────┐
│                        Delivery Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   HTTP Handler  │  │   Middleware    │                  │
│  │   - Book CRUD   │  │   - Request ID  │                  │
│  │   - Swagger     │  │   - Error Hdlr  │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                        Use Case Layer                       │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Book UseCase   │  │   Validation    │                  │
│  │   - Business    │  │   - Domain      │                  │
│  │   - Logic       │  │   - Rules       │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                        Domain Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │    Entities     │  │  Repositories   │                  │
│  │   - Book        │  │   - Interface   │                  │
│  │   - DTOs        │  │   - Contracts   │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                     Infrastructure Layer                    │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   PostgreSQL    │  │    Logging      │                  │
│  │   - sqlx        │  │   - Zap         │                  │
│  │   - Migrations  │  │   - Structured  │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

## Technology Stack

### Backend
- **Go 1.21+**: High-performance, statically typed language
- **Echo Framework**: Fast and lightweight web framework
- **PostgreSQL 15**: Advanced relational database with UUID support
- **sqlx**: Enhanced SQL toolkit with named queries and struct scanning
- **Zap**: Structured, high-performance logging
- **Testify**: Comprehensive testing framework with mocks
- **Swagger**: API documentation with echo-swagger integration
- **Prometheus**: Metrics collection and monitoring
- **Docker**: Containerization with multi-stage builds

### Frontend  
- **React 18**: Modern UI library with concurrent features
- **Next.js 14**: Full-stack React framework with App Router
- **TypeScript 5**: Strong typing and enhanced developer experience
- **Tailwind CSS**: Utility-first CSS framework for rapid styling
- **Context API**: Efficient state management without external dependencies

### Infrastructure
- **Docker Compose**: Multi-container orchestration
- **Nginx**: Reverse proxy and load balancer
- **Prometheus**: Metrics collection
- **Grafana**: Monitoring and visualization
- **PostgreSQL**: Persistent data storage with named volumes

## Prerequisites

- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Git**: For cloning the repository
- **Make**: For running build commands (optional)

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd byfood-project
```

### 2. Start Development Environment
```bash
# Using Make (recommended)
make up

# Or using Docker Compose directly
docker-compose up -d
```

### 3. Access the Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080  
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Database**: localhost:5432 (postgres/password)

### 4. Test the System
```bash
# Run backend tests with coverage
make test-coverage

# Test data persistence
make test-persistence
```

## Development

### Available Make Commands
```bash
make build     # Build all Docker images
make up        # Start development environment  
make down      # Stop and remove containers
make test      # Run backend unit tests
make logs      # View application logs
make clean     # Clean up containers and volumes
```

### Testing
The system includes comprehensive testing:

**Backend Tests** (56.2% Coverage):
- Unit tests for entities, use cases, and handlers
- Integration tests with real database
- Mock testing for isolated components

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run integration tests only
make test-integration
```

**Data Persistence Testing**:
```bash
# Test that data survives container restarts
make test-persistence
```

## 📖 API Documentation

Complete API documentation is available via Swagger UI:
- **Development**: http://localhost:8080/swagger/index.html
- **Production**: https://your-domain.com/swagger/index.html

### Core Endpoints
```
GET    /api/v1/books       # List all books with pagination
POST   /api/v1/books       # Create a new book
GET    /api/v1/books/{id}  # Get book by UUID
PUT    /api/v1/books/{id}  # Update book by UUID  
DELETE /api/v1/books/{id}  # Delete book by UUID
GET    /health             # Application health check
GET    /metrics            # Prometheus metrics
```

### Example API Usage
```bash
# Create a book
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code","author":"Robert Martin","year":2008}'

# Get all books
curl http://localhost:8080/api/v1/books

# Get specific book
curl http://localhost:8080/api/v1/books/{uuid}
```

## Production Deployment

### 1. Configure Environment
```bash
cp .env.example .env
# Edit .env with production values
```

### 2. Deploy Production Stack
```bash
# Deploy with monitoring
docker-compose -f docker-compose.prod.yml up -d

# Scale backend services
docker-compose -f docker-compose.prod.yml up -d --scale backend=3
```

### 3. Access Production Services  
- **Application**: https://your-domain.com
- **Prometheus**: http://your-domain.com:9090
- **Grafana**: http://your-domain.com:3001

## 📊 Monitoring & Observability

### Metrics Collected
- **HTTP Metrics**: Request duration, count, status codes
- **Database Metrics**: Query performance, connection pool status
- **Application Metrics**: Book count, operation success rates
- **System Metrics**: Memory usage, CPU utilization

### Health Checks
- **Application**: `/health` endpoint with dependency checking
- **Database**: Connection and query validation  
- **Container**: Docker health check configurations

## Security Features

### Authentication & Authorization
- **API Key Authentication**: Configurable API keys for production
- **Request Validation**: Input sanitization and validation
- **Rate Limiting**: Configurable requests per second and burst limits

### Security Headers
- **HSTS**: HTTP Strict Transport Security
- **CSP**: Content Security Policy
- **XSS Protection**: Cross-site scripting prevention
- **CORS**: Configurable cross-origin resource sharing

### Data Protection
- **UUID Primary Keys**: No sequential ID exposure
- **SQL Injection Prevention**: Parameterized queries with sqlx
- **Input Validation**: Comprehensive request validation

## Project Structure

```
byfood-project/
├── backend/                    # Go backend application
│   ├── internal/
│   │   ├── domain/            # Domain layer (entities, interfaces)
│   │   ├── usecases/          # Business logic layer
│   │   ├── repositories/      # Data access layer
│   │   ├── delivery/          # Presentation layer (handlers)
│   │   ├── middleware/        # HTTP middleware components
│   │   └── infrastructure/    # External concerns (database)
│   ├── test/                  # Integration tests
│   ├── docs/                  # Generated API documentation
│   └── Dockerfile             # Backend container configuration
├── frontend/                  # React frontend application
│   ├── src/
│   │   ├── app/              # Next.js App Router pages
│   │   ├── components/       # Reusable React components
│   │   ├── contexts/         # React Context providers
│   │   ├── types/           # TypeScript type definitions
│   │   └── api/             # API client functions
│   └── Dockerfile           # Frontend container configuration
├── docker-compose.yml        # Development environment
├── docker-compose.prod.yml   # Production environment  
├── .env.example              # Environment variable template
├── Makefile                  # Build automation
└── README.md                 # This file
```

## Contributing

1. **Fork the Repository**
2. **Create Feature Branch**: `git checkout -b feature/amazing-feature`
3. **Commit Changes**: `git commit -m 'Add amazing feature'`
4. **Push to Branch**: `git push origin feature/amazing-feature`  
5. **Open Pull Request**

### Development Guidelines
- Follow Clean Architecture principles
- Write comprehensive tests for new features
- Update API documentation for endpoint changes
- Use meaningful commit messages
- Ensure all tests pass before submitting PR

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.

## Achievements

- **Clean Architecture**: Proper separation of concerns
- **56.2% Test Coverage**: Comprehensive testing suite
- **Production Ready**: Docker, monitoring, security
- **API Documentation**: Complete Swagger documentation
- **Data Persistence**: Verified across container restarts
- **Performance Optimized**: Efficient queries and caching
- **Security Hardened**: Multiple security layers
- **Monitoring**: Prometheus metrics and health checks

## Acknowledgments

- **Clean Architecture** principles by Robert C. Martin
- **Echo Framework** for Go web development
- **React & Next.js** communities for frontend excellence
- **PostgreSQL** team for robust database technology
- **Docker** for containerization platform
- **Open Source** contributors and maintainers

---

**Built with ❤️ using Clean Architecture, Go, React, and modern DevOps practices.**
