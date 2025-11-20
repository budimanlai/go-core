# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12+ (optional, can use SQLite for development)
- Git

### Step 1: Install Dependencies

```bash
# Initialize Go module (if not already initialized)
go mod init github.com/budimanlai/go-core

# Install required packages
go get github.com/gofiber/fiber/v2
go get github.com/golang-jwt/jwt/v5
go get github.com/go-playground/validator/v10
go get golang.org/x/crypto/bcrypt
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get gorm.io/driver/sqlite  # For development/testing

# Download all dependencies
go mod tidy
```

### Step 2: Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your settings
nano .env
```

**Minimum configuration (.env):**
```env
SERVER_PORT=8080
JWT_SECRET=your-super-secret-key-min-32-chars
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_core_db
```

### Step 3: Setup Database (PostgreSQL)

```bash
# Create database
createdb go_core_db

# Or using psql
psql -U postgres
CREATE DATABASE go_core_db;
\q
```

**For SQLite (Development):**
No setup needed! Just change the database driver in your code.

### Step 4: Run Example Application

```bash
# Navigate to examples
cd examples/fiber

# Run the application
go run main.go
```

You should see:
```
[INFO] Starting application...
[INFO] Database connected successfully
[INFO] Server starting on 0.0.0.0:8080
```

### Step 5: Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Register Account:**
```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "john@example.com",
    "password": "SecurePass123!"
  }'
```

Save the `access_token` from the response.

**Get Account (Protected):**
```bash
curl http://localhost:8080/api/v1/accounts/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📦 Using in Your Project

### As a Library

```bash
# In your project directory
go get github.com/budimanlai/go-core
```

### Import and Use

```go
package main

import (
    "github.com/budimanlai/go-core/account/domain/usecase"
    "github.com/budimanlai/go-core/account/platform/persistence"
    "github.com/budimanlai/go-core/middleware/auth"
    "github.com/budimanlai/go-core/pkg/crypto"
)

func main() {
    // Initialize components
    db := setupDatabase()
    
    // Create repository
    accountRepo := persistence.NewAccountRepository(db)
    
    // Create use case
    passwordHasher := crypto.NewBcryptHasher(10)
    accountUC := usecase.NewAccountUsecase(accountRepo, passwordHasher)
    
    // Use JWT middleware
    jwtService := auth.NewJWTService(auth.JWTConfig{
        SecretKey: "your-secret",
        Issuer: "your-app",
        ExpirationHours: 24,
    })
    
    // Add to your Fiber app
    app.Use(auth.FiberJWTMiddleware(jwtService))
}
```

## 🎯 Common Use Cases

### Use Case 1: Add JWT Authentication to Existing App

```go
import "github.com/budimanlai/go-core/middleware/auth"

jwtService := auth.NewJWTService(auth.JWTConfig{
    SecretKey:       os.Getenv("JWT_SECRET"),
    Issuer:          "my-app",
    ExpirationHours: 24,
})

// Protect routes
app.Use(auth.FiberJWTMiddleware(jwtService))
```

### Use Case 2: Add Account Management Module

```go
import (
    accountHandler "github.com/budimanlai/go-core/account/handler"
    accountPersistence "github.com/budimanlai/go-core/account/platform/persistence"
    accountUsecase "github.com/budimanlai/go-core/account/domain/usecase"
)

// Setup
accountRepo := accountPersistence.NewAccountRepository(db)
accountUC := accountUsecase.NewAccountUsecase(accountRepo, passwordHasher)
accountHandler := accountHandler.NewAccountHandler(accountUC)

// Routes
app.Post("/register", accountHandler.Register)
app.Post("/login", accountHandler.Login)
```

### Use Case 3: Add All Middlewares

```go
import (
    "github.com/budimanlai/go-core/middleware/cors"
    "github.com/budimanlai/go-core/middleware/logging"
    "github.com/budimanlai/go-core/middleware/recovery"
    "github.com/budimanlai/go-core/middleware/ratelimit"
)

app.Use(recovery.FiberRecoveryMiddleware(recovery.DefaultConfig()))
app.Use(cors.FiberCORSMiddleware(cors.DefaultConfig()))
app.Use(logging.FiberLoggerMiddleware(logging.LoggerConfig{...}))
app.Use(ratelimit.FiberRateLimitMiddleware(ratelimit.DefaultConfig()))
```

## 🔧 Development Workflow

### 1. Add New Module

```bash
# Create structure (similar to account/)
mkdir -p mymodule/{domain/{entity,repository,usecase},dto,models,platform/{http,grpc,persistence},handler}

# Create files following the same pattern
# See account module as reference
```

### 2. Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./account/domain/usecase/...
```

### 3. Format Code

```bash
# Format all Go files
go fmt ./...

# Run linter
golangci-lint run
```

### 4. Build

```bash
# Build binary
go build -o bin/myapp examples/fiber/main.go

# Run binary
./bin/myapp
```

## 🐳 Docker Deployment (Optional)

### Create Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main examples/fiber/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Build and Run

```bash
# Build image
docker build -t go-core:latest .

# Run container
docker run -p 8080:8080 \
  -e JWT_SECRET=your-secret \
  -e DB_HOST=host.docker.internal \
  go-core:latest
```

## 🔍 Troubleshooting

### Database Connection Error
```
Error: Failed to connect to database
```
**Solution:**
- Check PostgreSQL is running: `pg_isready`
- Verify credentials in `.env`
- Check firewall settings

### JWT Error
```
Error: invalid or expired token
```
**Solution:**
- Ensure JWT_SECRET is set and >= 32 characters
- Check token hasn't expired
- Verify Authorization header format: `Bearer <token>`

### Import Error
```
Error: package github.com/budimanlai/go-core/xxx not found
```
**Solution:**
```bash
go mod tidy
go get github.com/budimanlai/go-core
```

## 📚 Next Steps

1. ✅ Read [README.md](../README.md) for full documentation
2. ✅ Study [ARCHITECTURE.md](ARCHITECTURE.md) for design decisions
3. ✅ Review [SECURITY.md](SECURITY.md) for security best practices
4. ✅ Check [TESTING.md](TESTING.md) for testing strategies
5. ✅ Explore example code in `examples/fiber/`

## 💡 Tips

- Start with the example application in `examples/fiber/`
- Use SQLite for quick development/testing
- Follow the existing module patterns (account, region)
- Write tests as you develop
- Keep domain layer free from external dependencies

## 🤝 Getting Help

- Check the documentation in `/docs`
- Review example implementations
- Look at test files for usage examples
- Open an issue on GitHub

---

**Happy coding! 🚀**
