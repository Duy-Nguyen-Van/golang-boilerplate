package services

import (
	"golang-boilerplate/internal/errors"
	"golang-boilerplate/internal/integration/auth"

	"github.com/Nerzal/gocloak/v13"
)

// AuthService handles authentication business logic
type AuthService struct {
	authProvider auth.AuthService
}

// NewAuthService creates a new auth service
func ProvideAuthService(authProvider auth.AuthService) AuthService {
	return AuthService{
		authProvider: authProvider,
	}
}

// ValidateUserToken validates a user token and returns user information
func (s *AuthService) ValidateUserToken(token string) (*gocloak.IntroSpectTokenResult, error) {
	user, err := s.authProvider.ValidateToken(token)
	if err != nil {
		return nil, errors.UnauthorizedError("Token validation failed", err).
			WithOperation("validate_user_token").
			WithResource("token")
	}

	// Additional business logic validation can be added here
	// For example, check if user is active, not banned, etc.

	return user, nil
}

// HasRole checks if a user has a specific role
func (s *AuthService) HasRole(user *auth.User, role string) bool {
	for _, userRole := range user.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if a user has any of the specified roles
func (s *AuthService) HasAnyRole(user *auth.User, roles ...string) bool {
	for _, requiredRole := range roles {
		if s.HasRole(user, requiredRole) {
			return true
		}
	}
	return false
}
