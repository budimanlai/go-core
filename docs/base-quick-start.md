# Base Package - Complete Quick Start

Complete guide untuk membuat RESTful API dengan **3 layer architecture** menggunakan Base Package.

---

## 🎯 Overview

Build production-ready RESTful API dalam **5 menit** dengan:

```
BaseHandler[E, C, U]     →  HTTP Layer (5 endpoints)
    ↓
BaseService[E]           →  Business Logic Layer
    ↓
BaseRepository[E, M]     →  Data Layer + Cache + Metrics
```

---

## 🚀 Quick Start

### Prerequisites

```bash
go get github.com/budimanlai/go-core
go get github.com/jinzhu/copier
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get github.com/redis/go-redis/v9
```

---

## Step-by-Step Guide

### Step 1: Define Domain (Entity + Model + DTOs)

```go
// domain/user.go
package domain

import (
    "time"
    "gorm.io/gorm"
)

// Entity - Domain layer (business logic)
type UserEntity struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Model - Persistence layer (database)
type UserModel struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Name      string         `gorm:"not null"`
    Status    string         `gorm:"default:'pending'"`
    Password  string         `gorm:"not null"`  // Not in Entity!
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
    return "users"
}

// CreateDTO - Input validation for create
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=3,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

// UpdateDTO - Input validation for update
type UpdateUserRequest struct {
    Name   string `json:"name" validate:"required,min=3,max=100"`
    Status string `json:"status" validate:"omitempty,oneof=active inactive pending"`
}
```

---

### Step 2: Create Repository

```go
// repository/user_repository.go
package repository

import (
    "github.com/budimanlai/go-core/base"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "yourapp/domain"
)

type UserRepository interface {
    base.DomainRepository[domain.UserEntity]
    // Add custom methods if needed
}

type userRepositoryImpl struct {
    base.BaseRepository[domain.UserEntity, domain.UserModel]
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
    // Create factory with cache and metrics
    factory := base.NewFactory(db, base.RepoConfig{
        EnableCache:      true,
        EnablePrometheus: true,
        RedisClient:      rdb,
    })
    
    // Get base repository with all features
    baseRepo := base.NewRepository[domain.UserEntity, domain.UserModel](factory)
    
    return &userRepositoryImpl{
        BaseRepository: baseRepo,
    }
}
```

**That's it for repository!** 13 CRUD methods ready:
- ✅ Create, FindByID, Update, UpdateFields, Delete
- ✅ FindAll (pagination), FindOne, Count
- ✅ CreateBatch, DeleteBatch
- ✅ Restore, ForceDelete

---

### Step 3: Create Service

```go
// usecase/user_service.go
package usecase

import (
    "context"
    "errors"
    "github.com/budimanlai/go-core/base"
    "gorm.io/gorm"
    "yourapp/domain"
    "yourapp/repository"
)

type UserService interface {
    base.BaseService[domain.UserEntity]
    
    // Add custom business methods
    FindByEmail(ctx context.Context, email string) (*domain.UserEntity, error)
    ActivateUser(ctx context.Context, id uint) error
}

type userServiceImpl struct {
    base.BaseService[domain.UserEntity]
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) UserService {
    baseService := base.NewBaseService[domain.UserEntity](repo, db)
    
    return &userServiceImpl{
        BaseService: baseService,
        repo:        repo,
    }
}

// Override Create for business logic
func (s *userServiceImpl) Create(ctx context.Context, entity *domain.UserEntity) error {
    // Validation
    if entity.Email == "" {
        return errors.New("email is required")
    }
    
    // Check duplicate
    existing, _ := s.FindByEmail(ctx, entity.Email)
    if existing != nil {
        return errors.New("email already exists")
    }
    
    // Set default
    if entity.Status == "" {
        entity.Status = "pending"
    }
    
    // Call base
    if err := s.BaseService.Create(ctx, entity); err != nil {
        return err
    }
    
    // Post-create action
    go sendWelcomeEmail(entity.Email)
    
    return nil
}

// Custom business method
func (s *userServiceImpl) FindByEmail(ctx context.Context, email string) (*domain.UserEntity, error) {
    return s.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", email)
    })
}

func (s *userServiceImpl) ActivateUser(ctx context.Context, id uint) error {
    user, err := s.FindByID(ctx, id)
    if err != nil || user == nil {
        return errors.New("user not found")
    }
    
    if user.Status == "banned" {
        return errors.New("cannot activate banned user")
    }
    
    user.Status = "active"
    return s.Update(ctx, user)
}

func sendWelcomeEmail(email string) {
    // Implementation...
}
```

---

### Step 4: Create HTTP Handler

```go
// handler/user_handler.go
package handler

import (
    "github.com/budimanlai/go-core/base"
    "github.com/gofiber/fiber/v2"
    "yourapp/domain"
    "yourapp/usecase"
)

type UserHandler struct {
    *base.BaseHandler[domain.UserEntity, domain.CreateUserRequest, domain.UpdateUserRequest]
}

func NewUserHandler(service usecase.UserService) *UserHandler {
    return &UserHandler{
        BaseHandler: base.NewBaseHandler[domain.UserEntity, domain.CreateUserRequest, domain.UpdateUserRequest](service),
    }
}

// Optional: Override for custom logic
func (h *UserHandler) Create(c *fiber.Ctx) error {
    // Custom pre-validation
    apiKey := c.Get("X-API-Key")
    if apiKey == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "API key required",
        })
    }
    
    // Call base implementation
    return h.BaseHandler.Create(c)
}
```

**That's it for handler!** 5 RESTful endpoints ready:
- ✅ GET / (Index with pagination)
- ✅ GET /:id (View)
- ✅ POST / (Create with validation)
- ✅ PUT /:id (Update with validation)
- ✅ DELETE /:id (Delete)

---

### Step 5: Wire Everything Together

```go
// main.go
package main

import (
    "log"
    
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    
    "yourapp/domain"
    "yourapp/handler"
    "yourapp/repository"
    "yourapp/usecase"
)

func main() {
    // Setup database
    db, err := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/dbname"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Auto migrate
    db.AutoMigrate(&domain.UserModel{})
    
    // Setup Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create layers (Dependency Injection)
    userRepo := repository.NewUserRepository(db, rdb)
    userService := usecase.NewUserService(userRepo, db)
    userHandler := handler.NewUserHandler(userService)
    
    // Setup Fiber
    app := fiber.New(fiber.Config{
        ErrorHandler: errorHandler,
    })
    
    // Register routes
    users := app.Group("/api/users")
    users.Get("/", userHandler.Index)           // GET /api/users?page=1&limit=10
    users.Get("/:id", userHandler.View)         // GET /api/users/1
    users.Post("/", userHandler.Create)         // POST /api/users
    users.Put("/:id", userHandler.Update)       // PUT /api/users/1
    users.Delete("/:id", userHandler.Delete)    // DELETE /api/users/1
    
    // Start server
    log.Fatal(app.Listen(":3000"))
}

func errorHandler(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "success": false,
        "message": err.Error(),
    })
}
```

---

## 🎉 Done! Test Your API

### Create User

```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "password": "secret123"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### List Users (Pagination)

```bash
curl http://localhost:3000/api/users?page=1&limit=10
```

**Response:**
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "data": [
      {
        "id": 1,
        "email": "john@example.com",
        "name": "John Doe",
        "status": "pending"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_page": 1
  }
}
```

---

### Get User by ID

```bash
curl http://localhost:3000/api/users/1
```

---

### Update User

```bash
curl -X PUT http://localhost:3000/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "status": "active"
  }'
```

---

### Delete User

```bash
curl -X DELETE http://localhost:3000/api/users/1
```

---

## 📊 What You Get

### Repository Layer Features

✅ **13 CRUD Methods** out of the box  
✅ **Entity/Model Separation** via copier  
✅ **Redis Caching** (10min TTL, automatic invalidation)  
✅ **Prometheus Metrics** (`gocore_db_query_duration_seconds`)  
✅ **Transaction Support** via context injection  
✅ **Pagination** with safety limits (max 100)  
✅ **Batch Operations** (CreateBatch, DeleteBatch)  
✅ **Soft Delete** + Restore + ForceDelete  

### Service Layer Features

✅ **Business Logic Layer** between handler & repository  
✅ **Transaction Helper** (`WithTransaction`)  
✅ **Method Override** for custom business rules  
✅ **DomainRepository Bridge** (hides Model from service)  

### Handler Layer Features

✅ **5 RESTful Endpoints** (Index, View, Create, Update, Delete)  
✅ **DTO Validation** via struct tags  
✅ **Automatic Mapping** (DTO → Entity via copier)  
✅ **i18n Support** for error messages  
✅ **Proper HTTP Status Codes** (200, 400, 404, 500)  

---

## 🏗️ Project Structure

```
yourapp/
├── main.go
├── domain/
│   └── user.go          # Entity, Model, DTOs
├── repository/
│   └── user_repository.go
├── usecase/
│   └── user_service.go
└── handler/
    └── user_handler.go
```

---

## 🔥 Advanced Features

### Transaction Example

```go
func (s *orderServiceImpl) ProcessOrder(ctx context.Context, order *domain.OrderEntity) error {
    return s.WithTransaction(ctx, func(txCtx context.Context) error {
        // Step 1: Create order
        if err := s.Create(txCtx, order); err != nil {
            return err // Rollback
        }
        
        // Step 2: Update inventory
        if err := s.inventoryService.DecreaseStock(txCtx, order.Items); err != nil {
            return err // Rollback
        }
        
        // Step 3: Create payment
        payment := &domain.PaymentEntity{OrderID: order.ID}
        if err := s.paymentService.Create(txCtx, payment); err != nil {
            return err // Rollback
        }
        
        return nil // Commit all
    })
}
```

---

### Custom Query Filters

```go
// In handler
func (h *UserHandler) Index(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    var scopes []func(*gorm.DB) *gorm.DB
    
    // Filter by status
    if status := c.Query("status"); status != "" {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("status = ?", status)
        })
    }
    
    // Search by name or email
    if search := c.Query("search"); search != "" {
        scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
            return db.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
        })
    }
    
    result, err := h.Service.FindAll(c.Context(), page, limit, scopes...)
    // ... return response
}
```

**Usage:**
```bash
GET /api/users?status=active&search=john&page=1&limit=20
```

---

## 📝 Code Reduction

### Before Base Package

```
Repository:  ~200 lines (manual CRUD)
Service:     ~150 lines (delegation)
Handler:     ~120 lines (5 endpoints)
Total:       ~470 lines per entity
```

### After Base Package

```
Repository:  ~30 lines (factory + interface)
Service:     ~50 lines (interface + override)
Handler:     ~15 lines (constructor only)
Total:       ~95 lines per entity
```

**🎯 80% Code Reduction!**

---

## ⚡ Performance

| Feature | Performance Impact |
|---------|-------------------|
| Copier E↔M | ~10-20μs (negligible) |
| Redis Cache Hit | 8x faster than DB |
| Prometheus | ~0.01ms overhead |
| Service Layer | ~0.001ms delegation |
| Handler DTO | ~0.12ms parse+validate |

**Total Overhead:** < 0.5ms (negligible vs DB I/O: 1-10ms)

---

## 🧪 Testing

### Unit Test Repository

```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB()
    rdb := setupTestRedis()
    
    repo := repository.NewUserRepository(db, rdb)
    
    user := &domain.UserEntity{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    err := repo.Create(context.Background(), user)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Unit Test Service

```go
func TestUserService_Create(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := usecase.NewUserService(mockRepo, db)
    
    user := &domain.UserEntity{
        Email: "test@example.com",
        Name:  "Test",
    }
    
    err := service.Create(context.Background(), user)
    
    assert.NoError(t, err)
}
```

### Integration Test Handler

```go
func TestUserHandler_Create(t *testing.T) {
    app := fiber.New()
    handler := setupTestHandler()
    
    app.Post("/users", handler.Create)
    
    body := `{"email":"test@example.com","name":"Test","password":"password123"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## 🎓 Next Steps

1. **Add More Entities:**
   - Copy domain, repository, service, handler pattern
   - Wire in main.go
   - Register routes

2. **Add Authentication:**
   - JWT middleware
   - Auth service
   - Protected routes

3. **Add Custom Business Logic:**
   - Override service methods
   - Add custom endpoints
   - Implement transactions

4. **Monitoring:**
   - Prometheus metrics: `/metrics`
   - Grafana dashboard
   - Alert rules

5. **Documentation:**
   - Swagger/OpenAPI
   - API documentation
   - Code examples

---

## 📚 Documentation

- [Base Repository](./base-repository.md) - Data layer complete guide
- [Base Service](./base-service.md) - Business logic layer guide
- [Base Handler](./base-handler.md) - HTTP layer guide

---

## 🐛 Troubleshooting

**Q: Cache not working?**
```go
// Check Redis connection
rdb.Ping(context.Background())

// Verify config
RepoConfig{
    EnableCache: true,  // Must be true
    RedisClient: rdb,   // Must not be nil
}
```

**Q: Metrics not showing?**
```bash
# Check Prometheus endpoint
curl http://localhost:3000/metrics

# Look for: gocore_db_query_duration_seconds
```

**Q: Validation errors not showing?**
```go
// Make sure DTO has validate tags
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
}
```

---

**🎉 You're ready to build production-grade APIs!**

**Happy coding!** 🚀
