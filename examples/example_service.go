package base

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/budimanlai/go-core/base"
	"gorm.io/gorm"
)

// ============================================
// Example: User Domain
// ============================================

// Entity - Domain layer
type UserEntity struct {
	ID        uint
	Email     string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Model - Persistence layer
type UserModel struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	Status    string         `gorm:"default:'pending'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}

// ============================================
// Example: User Service with Custom Logic
// ============================================

// Custom service interface
type UserService interface {
	base.BaseService[UserEntity]

	// Custom business methods
	FindByEmail(ctx context.Context, email string) (*UserEntity, error)
	ActivateUser(ctx context.Context, id uint) error
	RegisterUserWithProfile(ctx context.Context, user *UserEntity, profileData map[string]interface{}) error
}

// Custom service implementation
type userServiceImpl struct {
	base.BaseService[UserEntity]
	repo UserRepository // Type-safe repository
}

// Custom repository interface (implements DomainRepository[UserEntity])
type UserRepository interface {
	base.DomainRepository[UserEntity]
	// Add custom repo methods if needed
}

// Constructor
func NewUserService(repo UserRepository, db *gorm.DB) UserService {
	baseService := base.NewBaseService[UserEntity](repo, db)

	return &userServiceImpl{
		BaseService: baseService,
		repo:        repo,
	}
}

// ============================================
// Override Hooks for Business Logic
// ============================================

// BeforeCreate - Validation example
func (s *userServiceImpl) BeforeCreate(ctx context.Context, entity *UserEntity) error {
	// Validate email
	if entity.Email == "" {
		return errors.New("email is required")
	}

	// Check duplicate email
	existing, err := s.FindByEmail(ctx, entity.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already exists")
	}

	// Set default status
	if entity.Status == "" {
		entity.Status = "pending"
	}

	return nil
}

// AfterCreate - Notification example
func (s *userServiceImpl) AfterCreate(ctx context.Context, entity *UserEntity) error {
	// Send welcome email (async)
	go func() {
		fmt.Printf("Sending welcome email to %s\n", entity.Email)
	}()

	// Log audit trail
	fmt.Printf("User created: ID=%d, Email=%s\n", entity.ID, entity.Email)

	return nil
}

// BeforeUpdate - Status change validation
func (s *userServiceImpl) BeforeUpdate(ctx context.Context, entity *UserEntity) error {
	// Prevent status change from active to banned without approval
	existing, err := s.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}

	if existing.Status == "active" && entity.Status == "banned" {
		// Check if admin approval exists (example)
		return errors.New("status change from active to banned requires admin approval")
	}

	return nil
}

// AfterUpdate - Cache invalidation example
func (s *userServiceImpl) AfterUpdate(ctx context.Context, entity *UserEntity) error {
	// Invalidate cache, send event, etc.
	fmt.Printf("User updated: ID=%d, Status=%s\n", entity.ID, entity.Status)
	return nil
}

// ============================================
// Custom Business Methods
// ============================================

// FindByEmail - Custom query
func (s *userServiceImpl) FindByEmail(ctx context.Context, email string) (*UserEntity, error) {
	return s.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	})
}

// ActivateUser - Business logic example
func (s *userServiceImpl) ActivateUser(ctx context.Context, id uint) error {
	user, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Business rule: cannot activate banned user
	if user.Status == "banned" {
		return errors.New("cannot activate banned user")
	}

	if user.Status == "active" {
		return errors.New("user already active")
	}

	user.Status = "active"
	return s.Update(ctx, user)
}

// RegisterUserWithProfile - Transaction example
func (s *userServiceImpl) RegisterUserWithProfile(ctx context.Context, user *UserEntity, profileData map[string]interface{}) error {
	return s.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create user (BeforeCreate & AfterCreate hooks will run)
		if err := s.Create(txCtx, user); err != nil {
			return err
		}

		// Create profile with user ID
		profileData["user_id"] = user.ID
		// Assuming you have profileRepo...
		// if err := profileRepo.CreateFromMap(txCtx, profileData); err != nil {
		//     return err // Auto rollback
		// }

		fmt.Printf("User and profile created in transaction: UserID=%d\n", user.ID)
		return nil // Auto commit
	})
}

// ============================================
// Usage Example
// ============================================

func ExampleUserService() {
	// Setup (assuming db and repo initialized)
	var db *gorm.DB
	var userRepo UserRepository

	// Create service
	userService := NewUserService(userRepo, db)
	ctx := context.Background()

	// Example 1: Create user (with validation hooks)
	user := &UserEntity{
		Email: "john@example.com",
		Name:  "John Doe",
	}
	err := userService.Create(ctx, user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Created user ID: %d\n", user.ID)

	// Example 2: Find by email (custom method)
	found, _ := userService.FindByEmail(ctx, "john@example.com")
	fmt.Printf("Found user: %s\n", found.Name)

	// Example 3: Activate user (business logic)
	err = userService.ActivateUser(ctx, user.ID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 4: Pagination with filter
	result, _ := userService.FindAll(ctx, 1, 10, func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active")
	})
	fmt.Printf("Found %d active users\n", len(result.Data))

	// Example 5: Batch operations
	newUsers := []*UserEntity{
		{Email: "user1@example.com", Name: "User 1"},
		{Email: "user2@example.com", Name: "User 2"},
	}
	userService.CreateBatch(ctx, newUsers)

	// Example 6: Transaction
	userService.RegisterUserWithProfile(ctx, user, map[string]interface{}{
		"bio": "Software Engineer",
	})

	// Example 7: Soft delete & restore
	userService.Delete(ctx, user.ID)
	userService.Restore(ctx, user.ID)
	userService.ForceDelete(ctx, user.ID) // Permanent
}
