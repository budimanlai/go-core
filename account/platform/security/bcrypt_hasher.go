package security

import (
	"errors"

	"github.com/budimanlai/go-core/account/domain/usecase"
	"github.com/budimanlai/go-pkg/security"
)

// BcryptHasher implements usecase.PasswordHasher using go-pkg/security
type BcryptHasher struct{}

// NewBcryptHasher creates a new BcryptHasher instance
func NewBcryptHasher() usecase.PasswordHasher {
	return &BcryptHasher{}
}

// Hash generates a bcrypt hash from the given password
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash := security.HashPassword(password)
	if hash == "" {
		return "", errors.New("failed to hash password")
	}
	return hash, nil
}

// Verify checks if the given password matches the hashed password
func (h *BcryptHasher) Verify(hashedPassword, password string) bool {
	valid, _ := security.CheckPasswordHash(password, hashedPassword)
	return valid
}
