# Base Repository - Quick Start Guide

Panduan cepat untuk memulai menggunakan Base Repository pattern dengan **Entity/Model separation** untuk Clean Architecture.

## 🚀 Installation

```bash
go get github.com/budimanlai/go-core
go get github.com/jinzhu/copier  # For automatic Entity↔Model conversion
```

## 📝 3 Steps to Get Started

### Step 1: Define Entity and Model

```go
package domain

// Entity - Domain layer (business logic, clean)
type UserEntity struct {
    ID        uint
    Email     string
    Name      string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Model - Persistence layer (database, GORM tags)
type UserModel struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex"`
    Name      string
    Status    string         `gorm:"default:'active'"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
    return "users"
}
```

**Why separate?**
- ✅ Domain logic independent of database
- ✅ Better testability (Entity mocks easier)
- ✅ Clean Architecture compliance
- ✅ Zero manual mapping (copier handles conversion)

### Step 2: Create Repository

```go
package repository

import (
    "github.com/budimanlai/go-core/base"
    "yourapp/domain"
)

type UserRepository interface {
    base.BaseRepository[domain.UserEntity, domain.UserModel]
    // Add custom methods here
}

type userRepoImpl struct {
    base.BaseRepository[domain.UserEntity, domain.UserModel]
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
    factory := base.NewFactory(db, base.RepoConfig{
        EnableCache:      true,
        EnablePrometheus: true,
        RedisClient:      rdb,
    })
    
    return &userRepoImpl{
        BaseRepository: base.NewRepository[domain.UserEntity, domain.UserModel](factory),
    }
}
```

### Step 3: Use It!

```go
repo := NewUserRepository(db, rdb)
ctx := context.Background()

// Create - Entity → Model conversion automatic
user := &domain.UserEntity{Email: "john@example.com", Name: "John"}
repo.Create(ctx, user)
fmt.Printf("Created with ID: %d\n", user.ID) // ID populated!

// Read (cached!) - Model → Entity conversion automatic
found, _ := repo.FindByID(ctx, user.ID)

// Update - Entity → Model conversion automatic
repo.UpdateFields(ctx, user.ID, map[string]interface{}{
    "name": "John Doe",
})

// Delete
repo.Delete(ctx, user.ID)
```

**Conversion happens automatically via copier!** You work with clean Entities, repository handles Model internally.

## 🎯 Common Operations

### Bulk Insert
```go
users := []*domain.UserEntity{{...}, {...}, {...}}
repo.CreateBatch(ctx, users) // Efficient batching, IDs populated automatically
```

### Pagination
```go
result, _ := repo.FindAll(ctx, 1, 20, // page 1, 20 items
    func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", "active")
    },
)
// result.Data contains []UserEntity (converted from []UserModel)
fmt.Printf("Page %d of %d\n", result.Page, result.TotalPage)
```

### Find by Condition
```go
user, _ := repo.FindOne(ctx,
    func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", "john@example.com")
    },
)
// Returns UserEntity
```

### Count
```go
total, _ := repo.Count(ctx,
    func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", "active")
    },
)
```

### Soft Delete & Restore
```go
repo.Delete(ctx, 1)    // Soft delete
repo.Restore(ctx, 1)   // Restore
repo.ForceDelete(ctx, 1) // Permanent delete
```

## 🔥 Advanced Features

### Transactions
```go
db.Transaction(func(tx *gorm.DB) error {
    ctx = base.InjectTx(ctx, tx)
    
    repo.Create(ctx, user)
    repo.Create(ctx, profile)
    
    return nil // Auto commit/rollback
})
```

### Reusable Scopes
```go
var ActiveScope = func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

users, _ := repo.FindAll(ctx, 1, 10, ActiveScope)
```

### Custom Methods
```go
func (r *userRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.UserEntity, error) {
    return r.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", email)
    })
}
```

## 🧬 Entity vs Model

**Entity (E)** - Application layer:
- Clean domain objects
- No database knowledge
- Business logic focus
- Easy to test/mock

**Model (M)** - Persistence layer:
- GORM tags for database
- Table mapping
- Technical concerns
- Used internally

**Copier** automatically converts between E↔M:
- `Create`: E→M→DB, then copy ID back to E
- `FindByID`: DB→M→E
- `Update`: E→M→DB
- `FindAll`: DB→[]M→[]E

**Zero manual mapping code needed!**

## 📊 What You Get

✅ **13 CRUD Methods** out of the box (Create, FindByID, Update, Delete, etc.)  
✅ **Entity/Model Separation** for Clean Architecture  
✅ **Automatic Conversion** via copier (zero manual mapping)  
✅ **Redis Caching** automatic (stores Entity, not Model)  
✅ **Prometheus Metrics** built-in  
✅ **Transaction Support** via context  
✅ **Type Safety** with generics `[E any, M any]`  
✅ **85% Less Code** than manual implementation  
✅ **Copier Overhead** only ~10-20μs (negligible vs DB I/O)  

## 📖 Full Documentation

See [base-repository.md](./base-repository.md) for complete API reference and advanced usage.

## 🐛 Troubleshooting

**Cache not working?**
- Check Redis connection: `rdb.Ping(ctx)`
- Verify `EnableCache: true`

**Metrics not showing?**
- Check `/metrics` endpoint
- Verify `EnablePrometheus: true`

**Transaction not rolling back?**
- Use `ctx = base.InjectTx(ctx, tx)`
- Return error to trigger rollback

---

**Questions?** Open an issue on GitHub.
