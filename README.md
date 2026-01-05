# TradeIT - Microservices Trading Engine

A microservices-based trading system with separate services for user management/order placement and order matching/execution.

## Architecture

The project is organized into two independent microservices:

- **User-Order Service** (Port 8000): HTTP REST API for user authentication and order management
- **Engine Service**: Queue-based order matching and execution engine

For detailed architecture documentation, see [MICROSERVICES_STRUCTURE.md](MICROSERVICES_STRUCTURE.md).

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 13+
- Redis 7+

### Build

```bash
# Build user-order-service
go build ./cmd/user-order-service

# Build engine-service
go build ./cmd/engine-service
```

### Run with Docker

```bash
docker-compose -f docker-compose.yml up --build
```

Services will be available at:
- User-Order Service: `http://localhost:8000`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

## Project Structure

```
├── cmd/                    Service entry points
├── services/              Service implementations
├── internal/              Shared packages
├── docker-compose.yml     Docker orchestration
└── go.mod                Go module file
```

## Documentation

- [MICROSERVICES_STRUCTURE.md](MICROSERVICES_STRUCTURE.md) - Complete architecture and design
- [RESTRUCTURE_CHECKLIST.md](RESTRUCTURE_CHECKLIST.md) - Implementation checklist

## Environment Setup

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Update with your database and Redis credentials.

## API Endpoints

### User Management
- `POST /register` - Register new user
- `POST /login` - Login user

### Order Management  
- `POST /order/create` - Place new order
- `GET /order` - Get all orders
- `POST /order/sort` - Get orders with filters

All endpoints except register/login require JWT authentication.

## Development

The codebase follows the standard Go project layout:
- Business logic in `services/*/` directories
- Shared utilities in `internal/` directory
- Entry points in `cmd/` directory

See [MICROSERVICES_STRUCTURE.md](MICROSERVICES_STRUCTURE.md) for detailed information.