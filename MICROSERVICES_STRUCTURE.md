# TradeIT Microservices Architecture

## Overview

The TradeIT project has been restructured into a proper microservices architecture with two independent services that communicate via Redis queues.

## Repository Structure

```
TradeIT/
├── cmd/                                    # Service entrypoints
│   ├── user-order-service/
│   │   └── main.go                        # User-Order service entry point
│   └── engine-service/
│       └── main.go                        # Engine service entry point
│
├── services/                              # Service-specific implementations
│   ├── user-order-service/                # HTTP API for user and order management
│   │   ├── controller/                    # HTTP request handlers
│   │   │   ├── user.go
│   │   │   └── orders.go
│   │   ├── routes/                        # Route definitions
│   │   │   └── userRoutes.go
│   │   ├── services/                      # Business logic
│   │   │   ├── user.go
│   │   │   ├── orders.go
│   │   │   ├── orderpool.go
│   │   │   └── util.go
│   │   ├── Dockerfile                     # Docker configuration
│   │   └── main.go                        # (moved to cmd/)
│   │
│   └── engine-service/                    # Order matching and execution engine
│       ├── engine/                        # Core engine logic
│       │   ├── orderengine.go            # Queue processors
│       │   ├── orderbook.go              # Order matching
│       │   ├── ledger.go                 # Trade tracking
│       │   ├── trade.go                  # Trade execution
│       │   ├── struct.go                 # Data structures
│       │   ├── errors.json               # Error messages
│       │   ├── err.go                    # Error handling
│       │   └── init.go
│       ├── repo/                         # Data repository
│       │   ├── order.go                  # Order DB operations
│       │   └── orderPool.go              # Order pooling
│       ├── services/                     # Engine-specific services
│       │   ├── orders.go
│       │   ├── orderpool.go
│       │   └── util.go
│       ├── Dockerfile                    # Docker configuration
│       └── main.go                       # (moved to cmd/)
│
├── internal/                              # Shared packages (internal use)
│   ├── database/                         # Database initialization and setup
│   │   └── init.go
│   ├── models/                           # Shared data models
│   │   └── models.go
│   ├── middleware/                       # Shared middleware
│   │   ├── auth.go
│   │   ├── jwt.go
│   │   ├── userDetails.go
│   │   └── order.go
│   └── testing/                          # Shared test utilities
│       ├── orderdata.go
│       ├── cancel_scenarios_test.go
│       └── engine_test/
│
├── docker-compose.yml                    # Docker Compose orchestration
├── go.mod                                # Go module file (root level)
├── go.sum
├── main.go                               # (deprecated - use cmd/)
├── test.go                               # (deprecated)
└── [config, docs, etc.]
```

## Services

### 1. User-Order Service (Port 8000)

**Purpose:** Provides HTTP API for user authentication and order management.

**Key Responsibilities:**
- User registration and login
- Order placement and retrieval
- Account management (email, password, funds)

**Technology Stack:**
- Gin (web framework)
- GORM (ORM)
- PostgreSQL (data storage)
- Redis (inter-service communication)

**Endpoints:**
- `POST /register` - Register new user
- `POST /login` - User login
- `POST /order/create` - Place new order
- `GET /order/` - Get all orders
- `POST /order/sort` - Get orders with parameters

**Entry Point:** `cmd/user-order-service/main.go`

### 2. Engine Service (Port 8001)

**Purpose:** Handles order matching, execution, and trade processing.

**Key Responsibilities:**
- Process incoming orders from Redis queue
- Match buy and sell orders
- Execute trades
- Handle order cancellations and modifications
- Maintain ledger of trades

**Technology Stack:**
- Pure Go (no web framework)
- Redis (queue-based communication)
- PostgreSQL (persistence)
- Sync pools (order pooling for performance)

**Queue Operations:**
- Reads from: `OrderQueue`, `CancelQueue`, `ModifyQueue`
- Writes to: Response queues (DB 2 in Redis)

**Entry Point:** `cmd/engine-service/main.go`

## Import Paths

After restructuring, all import paths have been updated to reflect the new structure:

**Shared Packages (internal):**
```go
import (
    "TradeIT/internal/database"
    "TradeIT/internal/models"
    "TradeIT/internal/middleware"
    "TradeIT/internal/testing"
)
```

**User-Order Service:**
```go
import (
    "TradeIT/services/user-order-service/controller"
    "TradeIT/services/user-order-service/routes"
    "TradeIT/services/user-order-service/services"
)
```

**Engine Service:**
```go
import (
    "TradeIT/services/engine-service/engine"
    "TradeIT/services/engine-service/repo"
    "TradeIT/services/engine-service/services"
)
```

## Building and Deployment

### Build Locally

```bash
# Build user-order-service
go build ./cmd/user-order-service

# Build engine-service
go build ./cmd/engine-service
```

### Docker Deployment

```bash
# Build and run both services with dependencies
docker-compose -f docker-compose.yml up --build

# Services will start on:
# - User-Order Service: http://localhost:8000
# - Engine Service: (internal, no HTTP)
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```

## Communication Flow

```
Client HTTP Request
        ↓
[User-Order Service:8000]
        ↓
Validates & places order to Redis (OrderQueue)
        ↓
[Engine Service]
Processes from queue, matches orders
        ↓
Updates database, publishes response
        ↓
User-Order Service reads response
        ↓
Returns HTTP response to client
```

## Configuration Files

- `docker-compose.yml` - Orchestrates both services, PostgreSQL, and Redis
- `.env` - Environment variables (create from `.env.example`)
- `internal/database/init.go` - Database connection configuration

## Testing

Tests are located in `internal/testing/` with comprehensive scenarios for:
- Order book operations
- Cancellation handling
- Order modifications
- Market order matching

Run tests:
```bash
go test ./internal/testing/...
```

## Migration from Monolith

The repository has been migrated from a monolithic architecture to microservices while maintaining all existing functionality:

1. **Code Organization:** Separated concerns into service-specific and shared packages
2. **Imports Updated:** All import paths adjusted to new structure
3. **Build System:** Docker builds updated to use cmd/ entrypoints
4. **Data Flow:** Unchanged - services communicate via Redis queues as before
5. **Database:** Single PostgreSQL instance shared by both services

All existing business logic has been preserved without functional changes - only packaging and import paths were modified.

## Dependencies

- Go 1.21+
- PostgreSQL 13+
- Redis 7+
- Docker & Docker Compose (for containerized deployment)

## Next Steps

1. Update your CI/CD pipeline to build from `cmd/` directories
2. Configure environment variables in `.env`
3. Run `docker-compose up --build` to start the system
4. Test endpoints against the running User-Order Service
