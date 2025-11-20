package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrInvalidBasicAuth  = errors.New("invalid basic auth credentials")
	ErrMissingAuthHeader = errors.New("missing authorization header")
)

type BasicAuthConfig struct {
	Users map[string]string
}

type BasicAuthService interface {
	Validate(authHeader string) (string, error)
	ValidateCredentials(username, password string) bool
}

type basicAuthService struct {
	config BasicAuthConfig
}

func NewBasicAuthService(config BasicAuthConfig) BasicAuthService {
	return &basicAuthService{
		config: config,
	}
}

func (s *basicAuthService) Validate(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidBasicAuth
	}

	encoded := authHeader[len(prefix):]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidBasicAuth
	}

	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return "", ErrInvalidBasicAuth
	}

	username, password := parts[0], parts[1]

	if !s.ValidateCredentials(username, password) {
		return "", ErrInvalidBasicAuth
	}

	return username, nil
}

func (s *basicAuthService) ValidateCredentials(username, password string) bool {
	expectedPassword, ok := s.config.Users[username]
	if !ok {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1
}
