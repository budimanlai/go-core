# Architectural Decision Records (ADR)

## ADR-001: Clean Architecture Pattern

**Status:** Accepted

**Context:**
Kami memerlukan arsitektur yang scalable, testable, dan maintainable untuk microservice modules yang akan digunakan di berbagai project.

**Decision:**
Mengadopsi Clean Architecture dengan layer separation:
1. Domain Layer - Core business logic
2. Application Layer - Use cases
3. Infrastructure Layer - External dependencies
4. Presentation Layer - Handlers

**Consequences:**
- ✅ High testability - business logic independent
- ✅ Framework independence - mudah switch framework
- ✅ Database independence - mudah ganti DB
- ⚠️ Lebih banyak boilerplate code
- ⚠️ Learning curve untuk new developers

---

## ADR-002: Repository Pattern

**Status:** Accepted

**Context:**
Perlu abstraction untuk data access agar business logic tidak terikat dengan database implementation.

**Decision:**
Menggunakan Repository Pattern dengan interface di domain layer dan implementation di infrastructure layer.

**Consequences:**
- ✅ Testability - mock repositories untuk testing
- ✅ Flexibility - ganti DB tanpa ubah business logic
- ✅ Single Responsibility - clear separation
- ⚠️ Additional abstraction layer

---

## ADR-003: JWT for Authentication

**Status:** Accepted

**Context:**
Microservice architecture memerlukan stateless authentication mechanism.

**Decision:**
Menggunakan JWT (JSON Web Tokens) untuk authentication dengan:
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
Cross-cutting concerns (logging, auth, cors, etc.) perlu di-handle secara consistent.

**Decision:**
Implement middleware pattern yang:
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
Data recovery dan audit trail requirements.

**Decision:**
Implement soft delete sebagai default behavior:
- DeletedAt field di semua entities
- Hard delete tersedia tapi explicit

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
Perlu consistent error handling across all modules.

**Decision:**
- Domain errors (ErrAccountNotFound, etc.)
- AppError wrapper dengan error codes
- Generic error messages untuk external users
- Detailed errors untuk logging

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
Menggunakan bcrypt dengan cost factor 10+ untuk password hashing.

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
Perlu ORM yang mature dan well-supported untuk Go.

**Decision:**
Menggunakan GORM sebagai default ORM dengan options untuk:
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
Configuration management untuk different environments.

**Decision:**
- Environment variables untuk configuration
- .env files untuk development
- Config struct dengan defaults

**Consequences:**
- ✅ 12-factor app compliant
- ✅ Container-friendly
- ✅ Secure - no secrets in code
- ⚠️ Need to manage env vars

---

## ADR-010: API Versioning

**Status:** Accepted

**Context:**
Perlu backward compatibility saat API changes.

**Decision:**
URL-based versioning: `/api/v1/`, `/api/v2/`

**Consequences:**
- ✅ Clear and explicit
- ✅ Easy to route
- ✅ Multiple versions can coexist
- ⚠️ URL duplication for same resources
