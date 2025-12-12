# Dash - Uber Clone Application

A production-ready Uber-clone application built with Go, designed with microservice architecture in mind.

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters) and **Domain-Driven Design** principles. The structure is organized to make it easy to split into microservices in the future.

### Directory Structure

```
Dash/
├── cmd/
│   ├── api/                    # Application entry point
│   └── migrate/                # Migration command-line tool
│
├── core/                       # Business logic (Domain layer)
│   └── user/                   # User domain
│       ├── entity.go
│       ├── repository.go      # Repository interface (port)
│       └── service.go         # Business logic
│
├── db/                         # Database layer (adapter)
│   └── repository/            # Repository implementations
│       └── user/
│           └── postgresql/
│               ├── migrations/ # Database migrations
│               └── repository.go
│
├── pkg/                        # Public packages
│   ├── auth/                  # JWT authentication
│   ├── env/                   # Environment variable handling
│   ├── logger/                # Structured logging
│   └── middleware/            # HTTP middleware
│
└── internal/                   # Internal packages
    ├── api/                   # HTTP layer (adapter)
    │   └── user/
    │       ├── handler.go
    │       └── dto.go
    ├── database/              # Database connection
    └── router/                # Route setup
```

## Current Features

### Authentication
- **User Login**: Phone number-based authentication with password hashing
- **JWT Tokens**: Secure token-based authentication with access and refresh tokens

### API Endpoints

All API endpoints end with a trailing slash (`/`) and return `405 Method Not Allowed` for unsupported HTTP methods.

- `GET /api/health/` - Health check endpoint
- `POST /api/user/auth/login/` - User login endpoint

#### Login Request
```json
{
  "phone": "+998901234567",
  "password": "your_password"
}
```

#### Login Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Validation Errors
When validation fails, the API returns field-specific error messages:
```json
{
  "error": "validation_error",
  "fields": {
    "phone": "This field is required",
    "password": "This field is required"
  }
}
```

## Future Domain Entities

The application is designed to support the following entities:

1. **User** - Base authentication entity (phone-based auth)
2. **Passenger** - Passenger profiles
3. **DriverProfile** - Driver profiles (one-to-one with User)
4. **Vehicle** - Vehicle information (belongs to DriverProfile, unique plate)
5. **Agent** - Agents who can move passengers (related to Vehicle)
6. **Ride** - Ride bookings
7. **RideRequest** - Ride requests that drivers can accept/reject
8. **Payment** - Payment records
9. **Transaction** - Payment transactions (nested under Payment domain)

## Technology Stack

- **Go 1.25+**
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **Squirrel** - SQL query builder
- **JWT** - JSON Web Tokens for authentication

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL database(version 15)
- Environment variables (see Configuration)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Dash
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (all are optional with defaults):
```bash
export DB_HOST=localhost          # Default: localhost
export DB_PORT=5432               # Default: 5432
export DB_USER=dash               # Default: dash
export DB_PASSWORD=dash           # Default: dash
export DB_NAME=dash              # Default: dash
export DB_SSLMODE=disable        # Default: disable
export JWT_SECRET=your-secret-key-change-in-production  # Required
export CORS_ALLOW_ORIGIN=*       # Default: *
export PORT=8080                 # Default: 8080
```

**Note**: The default database credentials are `dash/dash` for local development. Make sure to set proper credentials in production.

4. Run database migrations:
```bash
# Using the built-in migration command
go run cmd/migrate/main.go -command up

# Or build and run
go build -o migrate cmd/migrate/main.go
./migrate -command up

# Other commands:
# -command up      # Apply all pending migrations
# -command down    # Rollback all migrations
# -command version # Show current migration version
# -command force -version N # Force migration to version N
# -steps N         # Apply/rollback N steps (use with up/down)
```

5. Run the application:
```bash
go run cmd/api/main.go
```

## Architecture Principles

### Hexagonal Architecture
- **Core**: Business logic with no external dependencies
- **Ports**: Interfaces defined in core (e.g., Repository interfaces)
- **Adapters**: Implementations in infrastructure layer (e.g., database, HTTP)

### Domain-Driven Design
- Each domain is self-contained
- Clear boundaries between domains
- Easy to extract into microservices

### Microservice-Ready
- Each domain in `core/` can become its own microservice
- Shared code in `pkg/` for cross-service communication
- Independent deployment capability

## Future Roadmap

### Phase 1: Core Entities
- [x] User authentication
- [ ] Passenger management
- [ ] DriverProfile management
- [ ] Vehicle management
- [ ] Agent management

### Phase 2: Ride Management
- [ ] Ride creation and management
- [ ] RideRequest system
- [ ] Driver ride acceptance/rejection

### Phase 3: Payment System
- [ ] Payment processing
- [ ] Transaction management
- [ ] Payment history

### Phase 4: Advanced Features
- [ ] Real-time location tracking
- [ ] Ride matching algorithm
- [ ] Rating and review system
- [ ] Notification system

### Phase 5: Microservices Migration
- [ ] Split into separate services:
  - User Service
  - Ride Service
  - Payment Service
  - Notification Service

## Development Guidelines

### Code Organization
- Business logic goes in `core/`
- HTTP handlers in `internal/api/`
- Database implementations in `db/repository/`
- Shared utilities in `pkg/` or `shared/`

### Authentication
- Phone number is used as username
- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- JWT secret must be set via `JWT_SECRET` environment variable

### Database
- Use Squirrel for query building
- All migrations in `db/repository/{domain}/postgresql/migrations/`
- Repository implementations in `db/repository/{domain}/postgresql/`
- Use PostgreSQL-specific features:
  - **UUIDv7** for primary keys (improved indexing and performance)
  - **TIMESTAMP WITH TIME ZONE** for all timestamp columns
  - **Automatic `updated_at` triggers** for timestamp updates

### API Conventions
- All endpoints must end with a trailing slash (`/`)
- Unsupported HTTP methods return `405 Method Not Allowed`
- Validation errors return field-specific messages
- Structured JSON logging with file and line numbers
- Error codes follow `FunctionName:Code` format


