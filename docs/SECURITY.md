# Security Best Practices

## 🔒 Password Security

### 1. Password Hashing
```go
// ✅ DO: Use bcrypt with appropriate cost
hasher := crypto.NewBcryptHasher(12) // cost 12 recommended for production
hashedPassword, err := hasher.Hash(plainPassword)

// ❌ DON'T: Use MD5 or SHA1
// ❌ DON'T: Store plain text passwords
```

### 2. Password Requirements
```go
// Implement strong password policy
type PasswordPolicy struct {
    MinLength      int  // minimum 8 characters
    RequireUpper   bool // at least 1 uppercase
    RequireLower   bool // at least 1 lowercase
    RequireDigit   bool // at least 1 digit
    RequireSpecial bool // at least 1 special character
}
```

### 3. Password Comparison
```go
// ✅ DO: Use constant-time comparison
isValid := hasher.Verify(storedHash, providedPassword)

// ❌ DON'T: Use simple string comparison
// This is vulnerable to timing attacks
```

## 🎫 JWT Security

### 1. Secret Key Management
```go
// ✅ DO: Use strong, random secret keys
// Minimum 256 bits (32 bytes) for HS256
jwtSecret := os.Getenv("JWT_SECRET")
if len(jwtSecret) < 32 {
    log.Fatal("JWT secret must be at least 32 characters")
}

// ✅ DO: Rotate secrets periodically
// ❌ DON'T: Hardcode secrets in source code
```

### 2. Token Expiration
```go
// ✅ DO: Set appropriate expiration times
JWTConfig{
    ExpirationHours: 1,  // Short-lived access tokens
    RefreshHours: 168,   // 7 days for refresh tokens
}

// ❌ DON'T: Create tokens that never expire
```

### 3. Token Validation
```go
// ✅ DO: Validate all claims
claims, err := jwtService.ValidateToken(token)
if err != nil {
    return ErrInvalidToken
}

// Check expiration
if time.Now().After(claims.ExpiresAt.Time) {
    return ErrExpiredToken
}

// Validate issuer
if claims.Issuer != expectedIssuer {
    return ErrInvalidIssuer
}
```

### 4. Token Storage (Client-side)
```
✅ DO:
- Store in httpOnly cookies
- Use secure flag in production
- Implement CSRF protection

❌ DON'T:
- Store in localStorage (XSS vulnerable)
- Store in sessionStorage
- Include sensitive data in token payload
```

## 🛡️ Input Validation

### 1. Always Validate User Input
```go
// ✅ DO: Use validation tags
type RegisterRequest struct {
    Email    string `validate:"required,email,max=255"`
    Username string `validate:"required,alphanum,min=3,max=50"`
    Password string `validate:"required,min=8,max=128"`
}

// ✅ DO: Validate in handler
validator := validator.NewValidator()
if err := validator.Validate(req); err != nil {
    return c.Status(400).JSON(fiber.Map{"error": err.Error()})
}
```

### 2. Sanitize Input
```go
import "html"

// ✅ DO: Escape HTML
sanitized := html.EscapeString(userInput)

// ✅ DO: Trim whitespace
username := strings.TrimSpace(req.Username)
```

### 3. SQL Injection Prevention
```go
// ✅ DO: Use parameterized queries
db.Where("email = ?", email).First(&user)

// ✅ DO: Use ORM properly
db.Model(&User{}).Where(&User{Email: email}).First(&user)

// ❌ DON'T: String concatenation in SQL
// sql := "SELECT * FROM users WHERE email = '" + email + "'"
```

## 🚦 Rate Limiting

### 1. Implement Rate Limiting
```go
// ✅ DO: Protect all public endpoints
app.Use(ratelimit.FiberRateLimitMiddleware(ratelimit.RateLimitConfig{
    Max:        100,              // 100 requests
    Expiration: 1 * time.Minute,  // per minute
}))

// ✅ DO: Stricter limits for auth endpoints
authLimiter := ratelimit.RateLimitConfig{
    Max:        5,                 // 5 attempts
    Expiration: 15 * time.Minute,  // per 15 minutes
}
app.Post("/login", ratelimit.FiberRateLimitMiddleware(authLimiter), loginHandler)
```

### 2. Per-User Rate Limiting
```go
// ✅ DO: Rate limit per user after authentication
KeyGenerator: func(c *fiber.Ctx) string {
    userID := c.Locals("user_id").(string)
    return userID
},
```

## 🔐 HTTPS & TLS

### 1. Force HTTPS in Production
```go
// ✅ DO: Redirect HTTP to HTTPS
if os.Getenv("ENVIRONMENT") == "production" {
    app.Use(func(c *fiber.Ctx) error {
        if c.Protocol() != "https" {
            return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
        }
        return c.Next()
    })
}
```

### 2. TLS Configuration
```go
// ✅ DO: Use strong TLS configuration
tlsConfig := &tls.Config{
    MinVersion:               tls.VersionTLS12,
    PreferServerCipherSuites: true,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    },
}
```

## 🌐 CORS Security

### 1. Configure CORS Properly
```go
// ✅ DO: Specific origins in production
corsConfig := cors.CORSConfig{
    AllowOrigins:     "https://yourdomain.com",
    AllowMethods:     "GET,POST,PUT,DELETE",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
}

// ❌ DON'T: Allow all origins in production
// AllowOrigins: "*"  // Only for development
```

### 2. Credentials Handling
```go
// ✅ DO: Be careful with credentials
if config.AllowCredentials {
    // Cannot use wildcard origin with credentials
    config.AllowOrigins = "https://specific-domain.com"
}
```

## 📝 Logging Security

### 1. Sensitive Data in Logs
```go
// ❌ DON'T: Log sensitive data
log.Printf("User login: %s with password: %s", email, password)

// ✅ DO: Log only non-sensitive data
log.Printf("User login attempt: %s", email)
log.Printf("Login successful for user: %s", userID)
```

### 2. Error Messages
```go
// ❌ DON'T: Expose internal errors to users
return c.Status(500).JSON(fiber.Map{
    "error": err.Error(), // May contain sensitive info
})

// ✅ DO: Generic error messages for users
log.Printf("Internal error: %v", err) // Log detailed error
return c.Status(500).JSON(fiber.Map{
    "error": "An internal error occurred",
})
```

## 🔑 API Key Security

### 1. API Key Storage
```go
// ✅ DO: Store API keys securely
// - Hash API keys before storing
// - Use environment variables
// - Never commit to version control

// ❌ DON'T: Store plain text API keys in database
```

### 2. API Key Validation
```go
// ✅ DO: Constant-time comparison
func ValidateAPIKey(provided, stored string) bool {
    return subtle.ConstantTimeCompare(
        []byte(provided),
        []byte(stored),
    ) == 1
}
```

## 🗄️ Database Security

### 1. Connection Security
```go
// ✅ DO: Use SSL/TLS for database connections
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
    host, port, user, password, dbname,
)

// ✅ DO: Use connection pooling
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 2. Least Privilege Principle
```
✅ DO:
- Create separate DB users for different services
- Grant only necessary permissions
- Use read-only users for read operations

❌ DON'T:
- Use root/admin for application
- Grant unnecessary permissions
```

## 🔍 Security Headers

### 1. Essential Security Headers
```go
app.Use(func(c *fiber.Ctx) error {
    // Prevent XSS attacks
    c.Set("X-Content-Type-Options", "nosniff")
    
    // Prevent clickjacking
    c.Set("X-Frame-Options", "DENY")
    
    // Enable XSS protection
    c.Set("X-XSS-Protection", "1; mode=block")
    
    // HSTS
    c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    
    // CSP
    c.Set("Content-Security-Policy", "default-src 'self'")
    
    return c.Next()
})
```

## 🎯 Session Security

### 1. Session Management
```go
// ✅ DO: Regenerate session ID after login
// ✅ DO: Set session timeout
// ✅ DO: Invalidate sessions on logout
// ✅ DO: Use secure session storage (Redis, etc.)

// ❌ DON'T: Use predictable session IDs
// ❌ DON'T: Store sensitive data in sessions
```

## 🚨 Security Monitoring

### 1. Log Security Events
```go
// ✅ DO: Log security-relevant events
logger.Info("User login successful", "user_id", userID, "ip", ip)
logger.Warn("Failed login attempt", "email", email, "ip", ip)
logger.Error("Suspicious activity detected", "details", details)
```

### 2. Implement Alerting
```
✅ DO:
- Monitor failed login attempts
- Track rate limit violations
- Alert on unusual patterns
- Regular security audits
```

## 📋 Security Checklist

### Before Production:
- [ ] All secrets in environment variables
- [ ] HTTPS enabled and enforced
- [ ] Rate limiting configured
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention verified
- [ ] XSS prevention verified
- [ ] CSRF protection enabled
- [ ] Security headers configured
- [ ] Error messages sanitized
- [ ] Logging configured properly
- [ ] JWT secrets are strong and rotated
- [ ] Password policy enforced
- [ ] Database connections secured
- [ ] CORS configured properly
- [ ] Dependencies up to date
- [ ] Security audit completed

---

**Remember: Security is not a feature, it's a requirement!**
