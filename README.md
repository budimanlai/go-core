# Go-Core: Reusable Microservice Modules

🚀 Go-Core adalah repository modular yang berisi komponen-komponen yang sering digunakan dalam pengembangan microservice dengan Go, mengadopsi **Clean Architecture**, **best practices**, dan **security-first approach**.

## 📋 Table of Contents

- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Modules](#modules)
- [Cara Penggunaan](#cara-penggunaan)
- [Best Practices](#best-practices)
- [Security](#security)
- [Microservice Ready](#microservice-ready)
- [Contoh Implementasi](#contoh-implementasi)

## 🏗️ Arsitektur

Repository ini mengadopsi **Clean Architecture** dengan pemisahan layer yang jelas:

```
┌─────────────────────────────────────────────────┐
│              Presentation Layer                 │
│  (Handlers: HTTP, gRPC, CLI, etc.)             │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│              Application Layer                  │
│  (Use Cases / Business Logic)                  │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│               Domain Layer                      │
│  (Entities, Repository Interfaces)             │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│           Infrastructure Layer                  │
│  (Database, External Services, etc.)           │
└─────────────────────────────────────────────────┘
```

### Prinsip Clean Architecture:

1. **Independence of Frameworks** - Tidak terikat dengan framework tertentu
2. **Testability** - Business logic dapat ditest tanpa UI, database, dll
3. **Independence of UI** - UI dapat diganti tanpa mengubah business logic
4. **Independence of Database** - Database dapat diganti dengan mudah
5. **Independence of External Agency** - Business logic tidak tahu tentang outside world

## 📁 Struktur Folder

```
go-core/
├── account/                    # Module Account
│   ├── domain/                 # Domain layer (core business logic)
│   │   ├── entity/            # Business entities
│   │   ├── repository/        # Repository interfaces
│   │   └── usecase/           # Use case interfaces & implementations
│   ├── dto/                   # Data Transfer Objects
│   ├── models/                # Database models
│   ├── platform/              # Infrastructure implementations
│   │   ├── http/             # HTTP client implementations
│   │   ├── grpc/             # gRPC implementations
│   │   └── persistence/      # Database implementations
│   └── handler/               # Request handlers (HTTP, gRPC, etc.)
│
├── region/                    # Module Region
│   └── [same structure as account]
│
├── middleware/                # Reusable middlewares
│   ├── auth/                 # JWT & Basic Auth
│   ├── logging/              # Request logging
│   ├── cors/                 # CORS handling
│   ├── recovery/             # Panic recovery
│   └── ratelimit/            # Rate limiting
│
├── pkg/                       # Shared packages
│   ├── errors/               # Custom error handling
│   ├── response/             # Standard API responses
│   ├── validator/            # Input validation
│   ├── crypto/               # Hashing & encryption
│   └── logger/               # Logging utilities
│
├── config/                    # Configuration management
├── docs/                      # Documentation
└── examples/                  # Usage examples
    ├── fiber/                # Fiber framework example
    └── grpc/                 # gRPC example
```

### Penjelasan Layer:

#### 1. Domain Layer (`domain/`)
- **Entity**: Core business objects, tidak bergantung pada apapun
- **Repository Interface**: Kontrak untuk data access
- **Usecase**: Business logic dan orchestration

#### 2. DTO Layer (`dto/`)
- Request/Response objects untuk API
- Validation tags
- Serialization logic

#### 3. Models Layer (`models/`)
- Database-specific models (GORM, SQL tags)
- Mapping antara Entity dan Database

#### 4. Platform Layer (`platform/`)
- **persistence**: Database implementations
- **http**: HTTP client implementations
- **grpc**: gRPC implementations
- Implementasi konkrit dari Repository interfaces

#### 5. Handler Layer (`handler/`)
- HTTP handlers (Fiber, Gin, Echo, etc.)
- gRPC handlers
- CLI handlers
- Convert DTO ↔ Entity

## 🔧 Modules

### 1. Account Module
Manajemen user/account dengan fitur:
- ✅ Registration & Authentication
- ✅ Password hashing (bcrypt)
- ✅ Role-based access control
- ✅ Soft delete support
- ✅ Account activation/deactivation

### 2. Region Module
Manajemen data region/lokasi (template sama dengan Account)

### 3. Middleware Collection

#### Authentication
```go
// JWT Authentication
jwtService := auth.NewJWTService(auth.JWTConfig{
    SecretKey: "your-secret",
    Issuer: "your-app",
    ExpirationHours: 24,
})
app.Use(auth.FiberJWTMiddleware(jwtService))

// Basic Authentication
basicAuth := auth.NewBasicAuthService(auth.BasicAuthConfig{
    Users: map[string]string{
        "admin": "password",
    },
})
app.Use(auth.FiberBasicAuthMiddleware(basicAuth))
```

#### Logging
```go
app.Use(logging.FiberLoggerMiddleware(logging.LoggerConfig{
    SkipPaths: []string{"/health"},
    LogFunc: func(entry logging.LogEntry) {
        log.Printf("%s %s - %d", entry.Method, entry.Path, entry.StatusCode)
    },
}))
```

#### CORS
```go
app.Use(cors.FiberCORSMiddleware(cors.DefaultConfig()))
```

#### Recovery
```go
app.Use(recovery.FiberRecoveryMiddleware(recovery.RecoveryConfig{
    EnableStackTrace: true, // for development
}))
```

#### Rate Limiting
```go
app.Use(ratelimit.FiberRateLimitMiddleware(ratelimit.RateLimitConfig{
    Max: 100,
    Expiration: 1 * time.Minute,
}))
```

## 🚀 Cara Penggunaan

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

```bash
# PostgreSQL example
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=yourdb
export DB_SSL_MODE=disable
```

### 3. Setup JWT Secret

```bash
export JWT_SECRET=your-super-secret-key
export JWT_EXPIRATION_HOURS=24
export JWT_ISSUER=go-core
```

### 4. Run Example Application

```bash
cd examples/fiber
go run main.go
```

### 5. Integrate ke Project Anda

```go
import (
    "github.com/budimanlai/go-core/account/domain/usecase"
    "github.com/budimanlai/go-core/account/platform/persistence"
    "github.com/budimanlai/go-core/middleware/auth"
)

// Initialize repository
accountRepo := persistence.NewAccountRepository(db)

// Initialize use case
accountUsecase := usecase.NewAccountUsecase(accountRepo, passwordHasher)

// Use in your handlers
accountHandler := handler.NewAccountHandler(accountUsecase)
```

## ✨ Best Practices

### 1. Dependency Injection
Gunakan constructor injection untuk semua dependencies:

```go
type accountUsecase struct {
    repo           repository.AccountRepository
    passwordHasher PasswordHasher
}

func NewAccountUsecase(repo repository.AccountRepository, hasher PasswordHasher) AccountUsecase {
    return &accountUsecase{
        repo:           repo,
        passwordHasher: hasher,
    }
}
```

### 2. Interface Segregation
Definisikan interface yang spesifik dan minimal:

```go
type AccountRepository interface {
    Create(ctx context.Context, account *entity.Account) error
    FindByID(ctx context.Context, id string) (*entity.Account, error)
}
```

### 3. Error Handling
Gunakan custom errors untuk business logic:

```go
var (
    ErrAccountNotFound = errors.New("account not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)
```

### 4. Context Propagation
Selalu pass context untuk cancellation dan timeout:

```go
func (u *accountUsecase) GetByID(ctx context.Context, id string) (*entity.Account, error) {
    return u.repo.FindByID(ctx, id)
}
```

### 5. Validation
Validate input di DTO level menggunakan struct tags:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

## 🔒 Security

### 1. Password Security
- ✅ Bcrypt hashing (cost 10+)
- ✅ Never log passwords
- ✅ Constant-time comparison

### 2. JWT Security
- ✅ Strong secret keys
- ✅ Token expiration
- ✅ HMAC-SHA256 signing
- ✅ Validate all claims

### 3. Input Validation
- ✅ Validate all user inputs
- ✅ Sanitize before database operations
- ✅ Use parameterized queries

### 4. Rate Limiting
- ✅ Prevent brute force attacks
- ✅ Per-IP limiting
- ✅ Custom limits per endpoint

### 5. CORS
- ✅ Configure allowed origins
- ✅ Limit allowed methods
- ✅ Set proper headers

### 6. Error Handling
- ✅ Don't expose internal errors
- ✅ Use generic error messages
- ✅ Log detailed errors server-side

## 🌐 Microservice Ready

### 1. Stateless Design
- Tidak ada session storage
- JWT untuk authentication
- Database untuk persistence

### 2. Service Discovery Ready
- Health check endpoints
- Graceful shutdown
- Configuration via environment

### 3. Horizontal Scaling
- No shared state
- Database connection pooling
- Idempotent operations

### 4. Monitoring & Observability
- Structured logging
- Request tracing
- Metrics collection ready

### 5. API Versioning
```go
api := app.Group("/api/v1")
```

### 6. Protocol Agnostic
Support multiple protocols:
- HTTP/REST (via Fiber, Gin, Echo)
- gRPC
- GraphQL (add your implementation)

## 📖 Contoh Implementasi

### Register Account

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "SecureP@ss123",
    "full_name": "John Doe"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "email": "user@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Login

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user@example.com",
    "password": "SecureP@ss123"
  }'
```

**Response:**
```json
{
  "data": {
    "account": { ... },
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### Get Account (Protected)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/accounts/uuid-here \
  -H "Authorization: Bearer eyJhbGc..."
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./account/domain/usecase/...
```

## 📝 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or support, please open an issue in the repository.

---

**Built with ❤️ using Clean Architecture principles**
