# Testing Guide

## 🧪 Testing Strategy

### Testing Pyramid
```
        E2E Tests (5%)
       ┌─────────────┐
      │   Integration  │ (15%)
     └─────────────────┘
    ┌───────────────────┐
   │    Unit Tests       │ (80%)
  └─────────────────────┘
```

## 📝 Unit Testing

### 1. Testing Domain Entities
```go
package entity_test

import (
    "testing"
    "time"

    "github.com/budimanlai/go-core/account/domain/entity"
    "github.com/stretchr/testify/assert"
)

func TestAccount_Activate(t *testing.T) {
    // Arrange
    account := &entity.Account{
        ID:       "test-id",
        IsActive: false,
    }

    // Act
    account.Activate()

    // Assert
    assert.True(t, account.IsActive)
    assert.WithinDuration(t, time.Now(), account.UpdatedAt, time.Second)
}

func TestAccount_Deactivate(t *testing.T) {
    account := &entity.Account{
        ID:       "test-id",
        IsActive: true,
    }

    account.Deactivate()

    assert.False(t, account.IsActive)
}

func TestAccount_SoftDelete(t *testing.T) {
    account := &entity.Account{
        ID: "test-id",
    }

    account.SoftDelete()

    assert.NotNil(t, account.DeletedAt)
    assert.True(t, account.IsDeleted())
}
```

### 2. Testing Use Cases with Mocks
```go
package usecase_test

import (
    "context"
    "testing"

    "github.com/budimanlai/go-core/account/domain/entity"
    "github.com/budimanlai/go-core/account/domain/usecase"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock Repository
type MockAccountRepository struct {
    mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *entity.Account) error {
    args := m.Called(ctx, account)
    return args.Error(0)
}

func (m *MockAccountRepository) FindByEmail(ctx context.Context, email string) (*entity.Account, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Account), args.Error(1)
}

// Mock Password Hasher
type MockPasswordHasher struct {
    mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
    args := m.Called(password)
    return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(hashedPassword, password string) bool {
    args := m.Called(hashedPassword, password)
    return args.Bool(0)
}

// Test Register Success
func TestAccountUsecase_Register_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockAccountRepository)
    mockHasher := new(MockPasswordHasher)
    uc := usecase.NewAccountUsecase(mockRepo, mockHasher)

    ctx := context.Background()
    email := "test@example.com"
    username := "testuser"
    password := "password123"
    fullName := "Test User"

    // Mock expectations
    mockRepo.On("FindByEmail", ctx, email).Return(nil, nil)
    mockRepo.On("FindByUsername", ctx, username).Return(nil, nil)
    mockHasher.On("Hash", password).Return("hashed_password", nil)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Account")).Return(nil)

    // Act
    account, err := uc.Register(ctx, email, username, password, fullName)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, account)
    assert.Equal(t, email, account.Email)
    assert.Equal(t, username, account.Username)
    assert.Equal(t, "hashed_password", account.Password)
    assert.True(t, account.IsActive)

    mockRepo.AssertExpectations(t)
    mockHasher.AssertExpectations(t)
}

// Test Register - Email Already Exists
func TestAccountUsecase_Register_EmailExists(t *testing.T) {
    mockRepo := new(MockAccountRepository)
    mockHasher := new(MockPasswordHasher)
    uc := usecase.NewAccountUsecase(mockRepo, mockHasher)

    ctx := context.Background()
    existingAccount := &entity.Account{Email: "test@example.com"}

    mockRepo.On("FindByEmail", ctx, "test@example.com").Return(existingAccount, nil)

    account, err := uc.Register(ctx, "test@example.com", "testuser", "password", "Test User")

    assert.Error(t, err)
    assert.Equal(t, usecase.ErrAccountAlreadyExists, err)
    assert.Nil(t, account)
}

// Test Login Success
func TestAccountUsecase_Login_Success(t *testing.T) {
    mockRepo := new(MockAccountRepository)
    mockHasher := new(MockPasswordHasher)
    uc := usecase.NewAccountUsecase(mockRepo, mockHasher)

    ctx := context.Background()
    existingAccount := &entity.Account{
        ID:       "test-id",
        Email:    "test@example.com",
        Password: "hashed_password",
        IsActive: true,
    }

    mockRepo.On("FindByEmail", ctx, "test@example.com").Return(existingAccount, nil)
    mockHasher.On("Verify", "hashed_password", "password123").Return(true)

    account, err := uc.Login(ctx, "test@example.com", "password123")

    assert.NoError(t, err)
    assert.NotNil(t, account)
    assert.Equal(t, "test-id", account.ID)
}

// Test Login - Invalid Credentials
func TestAccountUsecase_Login_InvalidCredentials(t *testing.T) {
    mockRepo := new(MockAccountRepository)
    mockHasher := new(MockPasswordHasher)
    uc := usecase.NewAccountUsecase(mockRepo, mockHasher)

    ctx := context.Background()
    existingAccount := &entity.Account{
        Email:    "test@example.com",
        Password: "hashed_password",
        IsActive: true,
    }

    mockRepo.On("FindByEmail", ctx, "test@example.com").Return(existingAccount, nil)
    mockHasher.On("Verify", "hashed_password", "wrong_password").Return(false)

    account, err := uc.Login(ctx, "test@example.com", "wrong_password")

    assert.Error(t, err)
    assert.Equal(t, usecase.ErrInvalidCredentials, err)
    assert.Nil(t, account)
}
```

### 3. Testing Handlers
```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"

    "github.com/budimanlai/go-core/account/dto"
    "github.com/budimanlai/go-core/account/handler"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockAccountUsecase struct {
    mock.Mock
}

func (m *MockAccountUsecase) Register(ctx context.Context, email, username, password, fullName string) (*entity.Account, error) {
    args := m.Called(ctx, email, username, password, fullName)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Account), args.Error(1)
}

func TestAccountHandler_Register_Success(t *testing.T) {
    // Arrange
    mockUC := new(MockAccountUsecase)
    handler := handler.NewAccountHandler(mockUC)

    app := fiber.New()
    app.Post("/register", handler.Register)

    account := &entity.Account{
        ID:       "test-id",
        Email:    "test@example.com",
        Username: "testuser",
        FullName: "Test User",
        IsActive: true,
    }

    mockUC.On("Register", mock.Anything, "test@example.com", "testuser", "password123", "Test User").
        Return(account, nil)

    reqBody := dto.RegisterRequest{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "password123",
        FullName: "Test User",
    }
    body, _ := json.Marshal(reqBody)

    // Act
    req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)

    // Assert
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🔗 Integration Testing

### 1. Database Integration Tests
```go
package integration_test

import (
    "context"
    "testing"

    "github.com/budimanlai/go-core/account/domain/entity"
    "github.com/budimanlai/go-core/account/models"
    "github.com/budimanlai/go-core/account/platform/persistence"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type RepositoryTestSuite struct {
    suite.Suite
    db   *gorm.DB
    repo repository.AccountRepository
}

func (suite *RepositoryTestSuite) SetupTest() {
    // Setup in-memory SQLite for testing
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    suite.NoError(err)

    // Auto-migrate
    err = db.AutoMigrate(&models.AccountModel{})
    suite.NoError(err)

    suite.db = db
    suite.repo = persistence.NewAccountRepository(db)
}

func (suite *RepositoryTestSuite) TearDownTest() {
    sqlDB, _ := suite.db.DB()
    sqlDB.Close()
}

func (suite *RepositoryTestSuite) TestCreate() {
    ctx := context.Background()
    account := &entity.Account{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "hashed_password",
        FullName: "Test User",
        IsActive: true,
    }

    err := suite.repo.Create(ctx, account)
    suite.NoError(err)
    suite.NotEmpty(account.ID)
}

func (suite *RepositoryTestSuite) TestFindByEmail() {
    ctx := context.Background()
    account := &entity.Account{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "hashed_password",
        FullName: "Test User",
    }

    suite.repo.Create(ctx, account)

    found, err := suite.repo.FindByEmail(ctx, "test@example.com")
    suite.NoError(err)
    suite.NotNil(found)
    suite.Equal(account.Email, found.Email)
}

func TestRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}
```

## 🚀 E2E Testing

### 1. API E2E Tests
```go
package e2e_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestE2E_AccountRegistrationFlow(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"

    // 1. Register new account
    registerReq := map[string]string{
        "email":     "e2e@example.com",
        "username":  "e2euser",
        "password":  "SecureP@ss123",
        "full_name": "E2E User",
    }
    body, _ := json.Marshal(registerReq)
    
    resp, err := http.Post(baseURL+"/public/register", "application/json", bytes.NewReader(body))
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)

    var registerResp map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&registerResp)
    accountID := registerResp["data"].(map[string]interface{})["id"].(string)

    // 2. Login with created account
    loginReq := map[string]string{
        "identifier": "e2e@example.com",
        "password":   "SecureP@ss123",
    }
    body, _ = json.Marshal(loginReq)
    
    resp, err = http.Post(baseURL+"/public/login", "application/json", bytes.NewReader(body))
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)

    var loginResp map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&loginResp)
    token := loginResp["data"].(map[string]interface{})["access_token"].(string)

    // 3. Get account details (authenticated)
    req, _ := http.NewRequest("GET", baseURL+"/accounts/"+accountID, nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err = http.DefaultClient.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 📊 Test Coverage

### Run Tests with Coverage
```bash
# Run all tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage in browser
open coverage.html

# Coverage by package
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Coverage Goals
- Unit tests: 80%+ coverage
- Integration tests: 60%+ coverage
- E2E tests: Critical paths covered

## 🎯 Testing Best Practices

### 1. Test Organization
```go
// Use table-driven tests
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 2. Test Fixtures
```go
// fixtures/account.go
package fixtures

func NewTestAccount() *entity.Account {
    return &entity.Account{
        ID:       "test-id",
        Email:    "test@example.com",
        Username: "testuser",
        FullName: "Test User",
        IsActive: true,
    }
}
```

### 3. Test Helpers
```go
// testhelpers/database.go
package testhelpers

func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })
    
    return db
}
```

## 🔍 Benchmarking

```go
func BenchmarkPasswordHash(b *testing.B) {
    hasher := crypto.NewBcryptHasher(10)
    password := "test_password_123"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        hasher.Hash(password)
    }
}

func BenchmarkJWTGeneration(b *testing.B) {
    jwtService := auth.NewJWTService(auth.JWTConfig{
        SecretKey:       "test-secret-key",
        ExpirationHours: 24,
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        jwtService.GenerateToken("user-id", "test@example.com", "user")
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

---

**Remember: Good tests are your safety net!**
