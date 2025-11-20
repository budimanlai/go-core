# Architectural Decision Records (ADR)

## ADR-001: Clean Architecture Pattern

**Status:** Accepted

**Context:**
We need a scalable, testable, and maintainable architecture for microservice modules that will be used across various projects.

**Decision:**
Adopt Clean Architecture with layer separation:
1. Domain Layer - Core business logic
2. Application Layer - Use cases
3. Infrastructure Layer - External dependencies
4. Presentation Layer - Handlers

**Consequences:**
- ✅ High testability - business logic independent
- ✅ Framework independence - easy to switch frameworks
- ✅ Database independence - easy to change DB
- ⚠️ More boilerplate code
- ⚠️ Learning curve for new developers

---

## ADR-002: Repository Pattern

**Status:** Accepted

**Context:**
Need abstraction for data access so business logic is not tied to database implementation.

**Decision:**
Use Repository Pattern with interfaces in domain layer and implementations in infrastructure layer.

**Consequences:**
- ✅ Testability - mock repositories for testing
- ✅ Flexibility - change DB without changing business logic
- ✅ Single Responsibility - clear separation
- ⚠️ Additional abstraction layer

---

## ADR-003: JWT for Authentication

**Status:** Accepted

**Context:**
Microservice architecture requires stateless authentication mechanism.

**Decision:**
Use JWT (JSON Web Tokens) for authentication with:
- HMAC-SHA256 signing
- 24 hours expiration (configurable)
- Claims: user_id, email, role

**Consequences:**
- ✅ Stateless - no session storage
- ✅ Scalable - works across multiple services
- ✅ Standard - widely supported
- ⚠️ Token size - larger than session ID
- ⚠️ Revocation - need blacklist mechanism

---

## ADR-004: Middleware Architecture

**Status:** Accepted

**Context:**
Cross-cutting concerns (logging, auth, cors, etc.) need to be handled consistently.

**Decision:**
Implement middleware pattern with:
- Framework-agnostic interfaces
- Framework-specific adapters (Fiber, Gin, etc.)
- Composable and reusable

**Consequences:**
- ✅ Reusability across frameworks
- ✅ Separation of concerns
- ✅ Easy to test
- ⚠️ Need adapters for each framework

---

## ADR-005: Soft Delete by Default

**Status:** Accepted

**Context:**
Data recovery and audit trail requirements.

**Decision:**
Implement soft delete as default behavior:
- DeletedAt field in all entities
- Hard delete available but explicit

**Consequences:**
- ✅ Data recovery capability
- ✅ Audit trail
- ✅ Compliance requirements
- ⚠️ Query complexity
- ⚠️ Storage overhead

---

## ADR-006: Error Handling Strategy

**Status:** Accepted

**Context:**
Need consistent error handling across all modules.

**Decision:**
- Domain errors (ErrAccountNotFound, etc.)
- AppError wrapper with error codes
- Generic error messages for external users
- Detailed errors for logging

**Consequences:**
- ✅ Security - no information leakage
- ✅ Consistency - same pattern everywhere
- ✅ Debuggability - detailed server logs
- ⚠️ Multiple error types to maintain

---

## ADR-007: Bcrypt for Password Hashing

**Status:** Accepted

**Context:**
Password security requirements.

**Decision:**
Use bcrypt with cost factor 10+ for password hashing.

**Consequences:**
- ✅ Industry standard
- ✅ Built-in salt
- ✅ Slow by design (prevents brute force)
- ⚠️ Higher CPU usage
- ⚠️ Slower than SHA256

---

## ADR-008: GORM as Default ORM

**Status:** Accepted

**Context:**
Need mature and well-supported ORM for Go.

**Decision:**
Use GORM as default ORM with options for:
- Raw SQL queries
- Other ORMs (sqlx, etc.)

**Consequences:**
- ✅ Productivity - auto-migration, associations
- ✅ Community support
- ✅ Multiple DB support
- ⚠️ Performance overhead
- ⚠️ Learning curve

---

## ADR-009: Environment-based Configuration

**Status:** Accepted

**Context:**
Configuration management for different environments.

**Decision:**
- Environment variables for configuration
- .env files for development
- Config struct with defaults

**Consequences:**
- ✅ 12-factor app compliant
- ✅ Container-friendly
- ✅ Secure - no secrets in code
- ⚠️ Need to manage env vars

---

## ADR-010: API Versioning

**Status:** Accepted

**Context:**
Need backward compatibility when API changes.

**Decision:**
URL-based versioning: `/api/v1/`, `/api/v2/`

**Consequences:**
- ✅ Clear and explicit
- ✅ Easy to route
- ✅ Multiple versions can coexist
- ⚠️ URL duplication for same resources

---

## ADR-011: Use go-pkg for Common Utilities

**Status:** Accepted

**Context:**
Need to avoid duplicating common utility functions and follow DRY principle. External package `github.com/budimanlai/go-pkg` already provides well-tested utilities.

**Decision:**
Use `go-pkg` for:
- Security (password hashing with bcrypt)
- Response formatting (standardized API responses)
- Logging (structured logging with timestamps)
- Validation (input validation with i18n)
- i18n (internationalization support)

**Consequences:**
- ✅ No code duplication
- ✅ Well-tested utilities
- ✅ Consistent patterns across projects
- ✅ Follows .clinerules guidelines
- ⚠️ External dependency to maintain
- ⚠️ Need to keep go-pkg updated

---

## ADR-012: Platform Layer Structure

**Status:** Accepted

**Context:**
Need clear separation between domain logic and infrastructure implementations. Previous structure had inconsistent placement of handlers and implementations.

**Decision:**
All infrastructure implementations must be in `platform/` subdirectory:
- `platform/http/` - HTTP/REST handlers
- `platform/grpc/` - gRPC service handlers
- `platform/repository/` - Database implementations
- `platform/usecase/` - Business logic implementations
- `platform/security/` - Security adapters

Domain layer contains only:
- `domain/entity/` - Business entities
- `domain/repository/` - Repository interfaces
- `domain/usecase/` - Use case interfaces

**Consequences:**
- ✅ Clear architectural boundaries
- ✅ Easy to add new delivery methods (CLI, GraphQL, WebSocket)
- ✅ Testable domain layer (no infrastructure dependencies)
- ✅ Framework independence
- ⚠️ More directory nesting
- ⚠️ Requires discipline to maintain separation

---

## Architecture Overview

### Layer Dependencies

```
┌─────────────────────────────────────┐
│         Delivery Layer              │
│  (HTTP, gRPC, CLI, GraphQL)         │
│    platform/http/                   │
│    platform/grpc/                   │
└──────────────┬──────────────────────┘
               │ depends on
┌──────────────▼──────────────────────┐
│      Application Layer              │
│    (Use Case Implementations)       │
│    platform/usecase/                │
└──────────────┬──────────────────────┘
               │ depends on
┌──────────────▼──────────────────────┐
│         Domain Layer                │
│    (Business Logic & Rules)         │
│    domain/entity/                   │
│    domain/repository/               │
│    domain/usecase/                  │
└──────────────┬──────────────────────┘
               │ depends on
┌──────────────▼──────────────────────┐
│    Infrastructure Layer             │
│  (External Dependencies)            │
│    platform/repository/             │
│    platform/security/               │
└─────────────────────────────────────┘
```

### Data Flow

```
HTTP Request
    ↓
[platform/http Handler]
    ↓
[platform/usecase Implementation]
    ↓ (calls interface)
[domain/usecase Interface]
    ↓
[platform/repository Implementation]
    ↓ (calls interface)
[domain/repository Interface]
    ↓
[domain/entity Business Logic]
```

### Key Principles

1. **Dependency Rule:** Dependencies point inward. Domain layer has no dependencies on outer layers.
2. **Interface Segregation:** Domain defines interfaces, platform implements them.
3. **Single Responsibility:** Each layer has one reason to change.
4. **Open/Closed:** Open for extension (new delivery methods), closed for modification (domain logic).

### Module Structure Example

```
account/
├── domain/              # Business rules (no external deps)
│   ├── entity/         # Business entities
│   ├── repository/     # Repository interfaces
│   └── usecase/        # Use case interfaces
├── dto/                # Data transfer objects
├── models/             # Database models (GORM)
└── platform/           # Infrastructure implementations
    ├── http/          # HTTP handlers (Fiber)
    ├── grpc/          # gRPC handlers
    ├── repository/    # Database implementation
    ├── security/      # Security adapters
    └── usecase/       # Business logic implementation
```

### Adding New Features

1. **New Module:** Follow account/ structure
2. **New Delivery Method:** Add to platform/ (e.g., platform/cli/)
3. **New Database:** Implement domain/repository interface
4. **New Use Case:** Define interface in domain/usecase, implement in platform/usecase

---

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [The Twelve-Factor App](https://12factor.net/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
