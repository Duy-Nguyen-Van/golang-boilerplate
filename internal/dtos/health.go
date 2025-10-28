package dtos

import "time"

// HealthResponse represents a health check response DTO
type HealthResponse struct {
	Status    string    `json:"status" example:"healthy"`
	Timestamp time.Time `json:"timestamp" example:"2021-01-01T00:00:00Z"`
	Version   string    `json:"version" example:"1.0.0"`
	Service   string    `json:"service" example:"golang-boilerplate"`
}
