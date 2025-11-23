# Real-Time Delivery State Manager

A concurrent in-memory service in Go for managing delivery drivers and customer orders in real time.

## Features

- **Clean Architecture** with clear separation of concerns (models, repository, service, use case, handler)
- **Thread-safe state management** using `sync.RWMutex` for concurrent access
- **RESTful HTTP APIs** for driver and order management
- **Background matching engine** that automatically assigns pending orders to available drivers
- **State validation** with proper state transition enforcement
- **Debug endpoint** for real-time state inspection

## Design Overview

### Architecture

The service is built around a central **StateManager** that owns all driver and order data in memory. All reads and writes go through this manager, ensuring thread safety using `sync.RWMutex`.

```
┌─────────────┐
│   HTTP API  │
└──────┬──────┘
       │
       ▼
┌─────────────────┐     ┌──────────────┐
│  StateManager   │◄────┤   Matcher    │
│  (RWMutex)      │     │  (goroutine) │
└─────────────────┘     └──────────────┘
       │
       ▼
  ┌────────┐  ┌────────┐
  │Drivers │  │Orders  │
  └────────┘  └────────┘
```

### Concurrency Strategy

- **RWMutex**: Allows multiple concurrent readers while ensuring exclusive write access
- **No direct map access**: All data access goes through StateManager methods
- **Atomic operations**: Order-driver assignment is atomic to prevent race conditions
- **Background goroutine**: Matcher runs independently every 3 seconds

### State Transitions

**Driver Status:**
- `available` ↔ `busy` ↔ `offline`

**Order Status:**
- `pending` → `assigned` → `picked_up` → `delivered`
- Any status → `canceled` (except `delivered`)

Invalid transitions are rejected by the StateManager.

## Setup and Run

### Prerequisites

- Go 1.21 or later

### Installation

```bash
# Clone or download the project
cd delivery-state-manager

# Download dependencies
go mod tidy

# Run the service
go run .

# Run with race detector (recommended for development)
go run -race .
```

The server will start on port **8080**.

## API Documentation

### Driver Endpoints

#### Create or Update Driver
```bash
POST /drivers
Content-Type: application/json

{
  "id": "driver-1",
  "name": "John Doe",
  "status": "available",
  "location": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

#### List All Drivers
```bash
GET /drivers
```

#### Get Driver Details
```bash
GET /drivers/{id}
```

#### Update Driver Status
```bash
PATCH /drivers/{id}/status
Content-Type: application/json

{
  "status": "busy"
}
```

**Valid statuses:** `available`, `busy`, `offline`

---

### Order Endpoints

#### Create Order
```bash
POST /orders
Content-Type: application/json

{
  "id": "order-1",
  "customer": "Jane Smith",
  "pickup": {
    "lat": 37.7749,
    "lon": -122.4194
  },
  "dropoff": {
    "lat": 37.8044,
    "lon": -122.2712
  }
}
```
**Note:** Orders are created with `status: "pending"` and will be automatically assigned by the matcher.

#### List All Orders
```bash
GET /orders
```

#### Get Order Details
```bash
GET /orders/{id}
```

#### Update Order Status
```bash
PATCH /orders/{id}/status
Content-Type: application/json

{
  "status": "picked_up"
}
```

**Valid statuses:** `pending`, `assigned`, `picked_up`, `delivered`, `canceled`

---

### Debug Endpoint

#### Get State Snapshot
```bash
GET /debug/state
```

Returns a complete snapshot of all drivers and orders with a timestamp.

## Example Workflow

```bash
# 1. Create an available driver
curl -X POST http://localhost:8080/drivers \
  -H "Content-Type: application/json" \
  -d '{
    "id": "driver-1",
    "name": "Alice",
    "status": "available",
    "location": {"lat": 37.7749, "lon": -122.4194}
  }'

# 2. Create a pending order
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-1",
    "customer": "Bob",
    "pickup": {"lat": 37.7749, "lon": -122.4194},
    "dropoff": {"lat": 37.8044, "lon": -122.2712}
  }'

# 3. Wait a few seconds for the matcher to run

# 4. Check the order status (should be "assigned" now)
curl http://localhost:8080/orders/order-1

# 5. Update order through its lifecycle
curl -X PATCH http://localhost:8080/orders/order-1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "picked_up"}'

curl -X PATCH http://localhost:8080/orders/order-1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "delivered"}'

# 6. View complete state
curl http://localhost:8080/debug/state
```

## Matching Engine

The background matcher runs every **3 seconds** and:

1. Finds all orders with `status: "pending"`
2. Finds all drivers with `status: "available"`
3. Matches them using **first-come-first-served** logic
4. Atomically updates:
   - Order: `status` → `assigned`, `driver_id` → driver's ID
   - Driver: `status` → `busy`

The matcher logs all matching activity for debugging.

## Testing

Run the service with the race detector to ensure thread safety:

```bash
go run -race .
```

Then run concurrent requests from multiple terminals to verify no race conditions occur.

## Project Structure

```
.
├── config/                      # Configuration management
│   └── config.go                # Environment-based configuration
├── internal/                    # Private application code
│   ├── models/                  # Domain models (entities)
│   │   └── models.go            # Driver, Order, Location, state machines
│   ├── repository/              # Data access layer
│   │   └── state_manager.go     # Thread-safe in-memory storage
│   ├── service/                 # Business services
│   │   └── matcher.go           # Background order-driver matching
│   ├── usecase/                 # Application business logic
│   │   ├── driver_usecase.go    # Driver operations
│   │   ├── order_usecase.go     # Order operations
│   │   └── debug_usecase.go     # Debug operations
│   └── handler/                 # HTTP presentation layer
│       ├── handlers.go          # REST API endpoints
│       └── handlers_test.go     # Comprehensive endpoint tests
├── test/                        # Integration tests
│   └── test.sh                  # End-to-end test script
├── docs/                        # Detailed documentation
│   ├── QUICK_START.md           # Getting started guide
│   ├── TECHNICAL_DOCUMENTATION.md # Comprehensive technical details
│   ├── INTERVIEW_GUIDE.md       # Interview Q&A preparation
│   └── ARCHITECTURE.md          # System architecture diagrams
├── main.go                      # Entry point with dependency wiring
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Clean Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:

**Dependency Flow:** `Handler → UseCase → Service/Repository → Models`

- **Models Layer**: Core business entities and rules, no external dependencies
- **Repository Layer**: Data access with thread-safe in-memory storage
- **Service Layer**: Background services (matcher runs every 3 seconds)
- **Use Case Layer**: Application business logic and orchestration
- **Handler Layer**: HTTP API presentation and request/response handling

**Benefits:**
- ✅ Each layer is independently testable
- ✅ Easy to swap implementations (e.g., in-memory → PostgreSQL)
- ✅ Business logic is framework-independent
- ✅ Clear separation enables parallel development
- ✅ Scales well as complexity grows

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed diagrams and explanations.

## Implementation Highlights

- **Clean Architecture**: Layered design with dependency inversion
- **Thread-safe**: `sync.RWMutex` for concurrent map access with defensive copying
- **Zero external database**: All data kept in memory with atomic operations
- **Validated transitions**: State machine enforces valid order status changes
- **Comprehensive tests**: All endpoints tested with success and error cases
- **Production-ready**: Environment config, structured logging, proper error handling
- **RESTful API**: Resource-based design with proper HTTP status codes
- **Background processing**: Goroutine-based matcher with safe concurrent access
