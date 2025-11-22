package security

import (
	pkgsecurity "github.com/budimanlai/go-pkg/security"
)

// PasswordHasher defines password hashing interface
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// BcryptHasher implements PasswordHasher using go-pkg/security
type BcryptHasher struct{}

// NewBcryptHasher creates a new BcryptHasher instance
func NewBcryptHasher() PasswordHasher {
	return &BcryptHasher{}
}

// HashPassword generates a bcrypt hash from the given password
func (h *BcryptHasher) HashPassword(password string) (string, error) {
	hash := pkgsecurity.HashPassword(password)
	if hash == "" {
		return "", nil
	}
	return hash, nil
}

// CheckPasswordHash compares a password with its hash
func (h *BcryptHasher) CheckPasswordHash(password, hash string) bool {
	valid, _ := pkgsecurity.CheckPasswordHash(password, hash)
	return valid
}
