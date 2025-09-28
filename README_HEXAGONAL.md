# LINE OA Backend - Hexagonal Architecture

This project has been refactored to use **Hexagonal Architecture** (also known as Ports and Adapters Architecture), which provides better separation of concerns, testability, and maintainability.

## Architecture Overview

```
internal/
├── domain/              # Core business logic (entities, value objects)
│   ├── entities/        # Business entities (User, Booking)
│   └── valueobjects/    # Value objects (LINEProfile, AuthToken)
├── ports/               # Interfaces (contracts)
│   ├── repositories.go  # Data persistence interfaces
│   ├── services.go      # External service interfaces
│   └── handlers.go      # HTTP handler interfaces
├── application/         # Use cases (business logic orchestration)
│   ├── auth_service.go  # Authentication use cases
│   └── booking_service.go # Booking use cases
├── adapters/            # External adapters (implementations)
│   ├── repositories/    # Database adapters
│   ├── services/        # External service adapters
│   ├── handlers/        # HTTP handlers
│   └── middleware/      # HTTP middleware
└── infrastructure/      # Infrastructure concerns
    ├── config/          # Configuration
    └── database/        # Database connection
```

## Key Benefits

### 1. **Separation of Concerns**
- **Domain Layer**: Contains pure business logic without external dependencies
- **Application Layer**: Orchestrates business operations (use cases)
- **Adapters Layer**: Handles external integrations (HTTP, database, APIs)
- **Infrastructure Layer**: Configuration and technical concerns

### 2. **Dependency Inversion**
- Core business logic doesn't depend on external frameworks
- Dependencies flow inward toward the domain
- Easy to swap implementations (e.g., change from MongoDB to PostgreSQL)

### 3. **Testability**
- Business logic can be tested in isolation
- Mock implementations can be easily injected
- Clear boundaries between layers

### 4. **Maintainability**
- Changes in external services don't affect business logic
- Clear interfaces make the system more predictable
- Easier to add new features or modify existing ones

## Layer Responsibilities

### Domain Layer (`internal/domain/`)
- **Entities**: Core business objects with behavior (`User`, `Booking`)
- **Value Objects**: Immutable objects that describe aspects of the domain (`LINEProfile`, `AuthToken`)
- **Business Rules**: Domain-specific validation and logic

### Ports Layer (`internal/ports/`)
- **Primary Ports**: Interfaces for incoming requests (HTTP handlers)
- **Secondary Ports**: Interfaces for outgoing requests (repositories, external services)
- **Contracts**: Define what the application needs without specifying how

### Application Layer (`internal/application/`)
- **Use Cases**: Orchestrate business operations
- **Application Services**: Coordinate between domain entities and external services
- **Business Workflows**: Implement complex business processes

### Adapters Layer (`internal/adapters/`)
- **Inbound Adapters**: HTTP handlers, CLI commands, message consumers
- **Outbound Adapters**: Database repositories, external API clients, message producers
- **Implementation Details**: Concrete implementations of port interfaces

### Infrastructure Layer (`internal/infrastructure/`)
- **Configuration**: Environment variables, settings
- **Database**: Connection management, migrations
- **Dependency Injection**: Wiring up the application

## Running the Application

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your actual values
   ```

3. **Run the application:**
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Generate LINE OAuth URL
- `POST /api/v1/auth/callback` - Handle LINE OAuth callback
- `POST /api/v1/auth/refresh` - Refresh JWT token

### User (Protected)
- `GET /api/v1/me` - Get current user profile

### Bookings (Protected)
- `POST /api/v1/bookings` - Create new booking
- `GET /api/v1/bookings` - Get user's bookings
- `GET /api/v1/bookings/:id` - Get specific booking
- `PUT /api/v1/bookings/:id` - Update booking
- `DELETE /api/v1/bookings/:id` - Cancel booking

## Testing Strategy

### Unit Tests
- Test domain entities and value objects in isolation
- Test application services with mocked dependencies
- Test individual adapters

### Integration Tests
- Test complete workflows through the application layer
- Test database operations with test database
- Test external API integrations with mock servers

### Example Test Structure
```go
// Domain entity test
func TestUser_UpdateProfile(t *testing.T) {
    user := entities.NewUser("line123", "John", "", "pic.jpg")
    user.UpdateProfile("John Doe", "john@example.com", "new-pic.jpg")
    
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}

// Application service test with mocks
func TestAuthService_AuthenticateWithLINE(t *testing.T) {
    mockUserRepo := &MockUserRepository{}
    mockLINEService := &MockLINEOAuthService{}
    mockJWTService := &MockJWTService{}
    
    authService := application.NewAuthService(mockUserRepo, mockLINEService, mockJWTService)
    
    // Test implementation...
}
```

## Migration from Old Architecture

The old architecture has been preserved in the root directories. To fully migrate:

1. **Update imports** in any remaining files to use the new internal structure
2. **Remove old directories** once migration is complete:
   - `handlers/`
   - `services/`
   - `models/`
   - `middleware/`
   - `database/`
   - `config/`

3. **Update deployment scripts** to use `cmd/main.go` instead of `_cmd/main.go`

## Future Enhancements

1. **Add more domain events** for better decoupling
2. **Implement CQRS pattern** for read/write separation
3. **Add domain services** for complex business operations
4. **Implement event sourcing** for audit trails
5. **Add more comprehensive validation** using domain rules

This architecture provides a solid foundation for scaling the application while maintaining clean separation of concerns and high testability.
