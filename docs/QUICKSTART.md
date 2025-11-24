# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- Go 1.21 or higher
- MySql or PostgreSQL 12+ (optional, can use SQLite for development)
- Git

### Step 1: Install Dependencies

```bash
# Clone the repository
git clone https://github.com/budimanlai/go-core.git
cd go-core

# Install required packages
go get github.com/budimanlai/go-pkg
go get github.com/gofiber/fiber/v2
go get github.com/golang-jwt/jwt/v5
go get github.com/go-playground/validator/v10
go get golang.org/x/crypto/bcrypt
go get gorm.io/gorm
go get gorm.io/driver/postgres

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
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
JWT_SECRET=your-super-secret-key-min-32-chars-long
JWT_ISSUER=go-core
JWT_EXPIRATION_HOURS=24

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_core_db
DB_SSLMODE=disable
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

### Step 4: Run Example Application

```bash
# Navigate to examples
cd examples/fiber

# Run the application
go run main.go
```

You should see:
```
[2025-11-21 10:30:00] Starting application...
[2025-11-21 10:30:00] Database connected successfully
[2025-11-21 10:30:00] INFO: Server starting on 0.0.0.0:8080
```

### Step 5: Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "go-core"
}
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

**Response:**
```json
{
  "meta": {
    "success": true,
    "message": "Account registered successfully"
  },
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "role": "user",
    "is_active": true,
    "created_at": "2025-11-21T10:30:15Z",
    "updated_at": "2025-11-21T10:30:15Z"
  }
}
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

**Response:**
```json
{
  "meta": {
    "success": true,
    "message": "Login successful"
  },
  "data": {
    "account": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "username": "johndoe"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

Save the `access_token` from the response.

**Get Account (Protected):**
```bash
curl http://localhost:8080/api/v1/accounts/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**List Accounts (Protected):**
```bash
curl "http://localhost:8080/api/v1/accounts?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📦 Using in Your Project

### Import as Module

```go
import (
    accountHTTP "github.com/budimanlai/go-core/account/platform/http"
    accountRepository "github.com/budimanlai/go-core/account/platform/repository"
    accountSecurity "github.com/budimanlai/go-core/account/platform/security"
    accountUsecase "github.com/budimanlai/go-core/account/platform/usecase"
    "github.com/budimanlai/go-core/config"
    "github.com/budimanlai/go-core/middleware/auth"
)

func main() {
    // Load config
    cfg := config.LoadConfig()
    
    // Setup database
    db := setupDatabase(cfg)
    
    // Initialize dependencies
    passwordHasher := accountSecurity.NewBcryptHasher()
    accountRepo := accountRepository.NewAccountRepository(db)
    accountUC := accountUsecase.NewAccountUsecase(accountRepo, passwordHasher)
    accountHandler := accountHTTP.NewAccountHandler(accountUC)
    
    // Setup JWT
    jwtService := auth.NewJWTService(auth.JWTConfig{
        SecretKey:       cfg.JWTSecret,
        Issuer:          cfg.JWTIssuer,
        ExpirationHours: cfg.JWTExpirationHours,
    })
    
    // Setup Fiber app
    app := fiber.New()
    
    // Public routes
    public := app.Group("/api/v1/public")
    public.Post("/register", accountHandler.Register)
    public.Post("/login", accountHandler.Login)
    
    // Protected routes
    accounts := app.Group("/api/v1/accounts")
    accounts.Use(auth.FiberJWTMiddleware(jwtService))
    accounts.Get("/", accountHandler.List)
    accounts.Get("/:id", accountHandler.GetByID)
    accounts.Delete("/:id", accountHandler.Delete)
    
    app.Listen(":8080")
}
```

## 🎯 Common Use Cases

### Use Case 1: Add JWT Authentication

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

### Use Case 2: Add All Middlewares

```go
import (
    "github.com/budimanlai/go-core/middleware/cors"
    "github.com/budimanlai/go-core/middleware/logging"
    "github.com/budimanlai/go-core/middleware/ratelimit"
    "github.com/budimanlai/go-core/middleware/recovery"
)

app.Use(recovery.FiberRecoveryMiddleware(recovery.DefaultConfig()))
app.Use(cors.FiberCORSMiddleware(cors.DefaultConfig()))
app.Use(logging.FiberLoggerMiddleware(logging.DefaultConfig()))
app.Use(ratelimit.FiberRateLimitMiddleware(ratelimit.DefaultConfig()))
```

### Use Case 3: Implement Custom Business Logic

```go
// 1. Define interface in domain/usecase
type MyUsecase interface {
    DoSomething(ctx context.Context, input string) (*Result, error)
}

// 2. Implement in platform/usecase
type myUsecaseImpl struct {
    repo MyRepository
}

func (u *myUsecaseImpl) DoSomething(ctx context.Context, input string) (*Result, error) {
    // Your business logic here
    return &Result{}, nil
}

// 3. Create handler in platform/http
func (h *MyHandler) HandleRequest(c *fiber.Ctx) error {
    result, err := h.usecase.DoSomething(c.Context(), input)
    if err != nil {
        return response.Error(c, 500, err.Error())
    }
    return response.Success(c, "Success", result)
}
```

## 🔧 Troubleshooting

### Database Connection Failed
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U postgres -h localhost
```

### Port Already in Use
```bash
# Change SERVER_PORT in .env
SERVER_PORT=8081

# Or kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### JWT Token Invalid
- Check JWT_SECRET in .env (minimum 32 characters)
- Ensure token is not expired
- Verify Authorization header format: `Bearer TOKEN`

## 📚 Next Steps

1. **Read Architecture Docs:** `docs/ARCHITECTURE.md`
2. **Learn Security Best Practices:** `docs/SECURITY.md`
3. **Write Tests:** `docs/TESTING.md`
4. **Explore Folder Structure:** `docs/STRUCTURE_SUMMARY.md`

## 🆘 Need Help?

- Check existing examples in `examples/`
- Read inline code documentation
- Open an issue on GitHub
- Consult `.clinerules` for development guidelines

**Happy coding! 🚀**
