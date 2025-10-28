package models

import (
	"time"
)

// AuthToken represents an authentication token
type AuthToken struct {
	BaseModel
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired checks if the token is expired
func (t *AuthToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// UserSession represents a user session
type UserSession struct {
	BaseModel
	UserID     string    `json:"user_id" db:"user_id"`
	Token      string    `json:"token" db:"token"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastUsedAt time.Time `json:"last_used_at" db:"last_used_at"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
}

// IsValid checks if the session is valid
func (s *UserSession) IsValid() bool {
	return s.UserID != "" && s.Token != "" && !s.IsExpired()
}

// IsExpired checks if the session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
