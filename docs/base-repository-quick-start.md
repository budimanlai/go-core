# Base Repository - Quick Start Guide

Panduan cepat untuk memulai menggunakan Base Repository pattern.

## 🚀 Installation

```bash
go get github.com/budimanlai/go-core
```

## 📝 3 Steps to Get Started

### Step 1: Define Entity

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex"`
    Name      string
    Status    string         `gorm:"default:'active'"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### Step 2: Create Repository

```go
type UserRepository interface {
    base.BaseRepository[User]
    // Add custom methods here
}

type userRepoImpl struct {
    base.BaseRepository[User]
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
    factory := base.NewFactory(db, base.RepoConfig{
        EnableCache:      true,
        EnablePrometheus: true,
        RedisClient:      rdb,
    })
    
    return &userRepoImpl{
        BaseRepository: base.NewRepository[User](factory),
    }
}
```

### Step 3: Use It!

```go
repo := NewUserRepository(db, rdb)
ctx := context.Background()

// Create
user := &User{Email: "john@example.com", Name: "John"}
repo.Create(ctx, user)

// Read (cached!)
found, _ := repo.FindByID(ctx, user.ID)

// Update
repo.UpdateFields(ctx, user.ID, map[string]interface{}{
    "name": "John Doe",
})

// Delete
repo.Delete(ctx, user.ID)
```

## 🎯 Common Operations

### Bulk Insert
```go
users := []*User{{...}, {...}, {...}}
repo.CreateBatch(ctx, users) // Efficient batching
```

### Pagination
```go
result, _ := repo.FindAll(ctx, 1, 20, // page 1, 20 items
    func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", "active")
    },
)
fmt.Printf("Page %d of %d\n", result.Page, result.TotalPage)
```

### Find by Condition
```go
user, _ := repo.FindOne(ctx,
    func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", "john@example.com")
    },
)
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
func (r *userRepoImpl) FindByEmail(ctx context.Context, email string) (*User, error) {
    return r.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", email)
    })
}
```

## 📊 What You Get

✅ **12 CRUD Methods** out of the box  
✅ **Redis Caching** automatic  
✅ **Prometheus Metrics** built-in  
✅ **Transaction Support** via context  
✅ **Type Safety** with generics  
✅ **85% Less Code** than manual implementation  

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
