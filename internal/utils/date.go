package utils

import (
	"time"

	appErrors "golang-boilerplate/internal/errors"
)

// DateRange represents a parsed date range with validation
type DateRange struct {
	StartDate *time.Time
	EndDate   *time.Time
}

// ParseDateRange parses start_date and end_date query parameters and validates them
// Returns a DateRange struct and any validation errors
func ParseDateRange(startDateStr, endDateStr string) (*DateRange, error) {
	var startDate, endDate *time.Time

	// Parse start date if provided
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339Nano, startDateStr)
		if err != nil {
			return nil, appErrors.ValidationError("Invalid start_date format. Must be RFC3339Nano format (e.g., 2025-03-08T15:05:42.536581Z)", err)
		}
		startDate = &parsedStartDate
	}

	// Parse end date if provided
	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339Nano, endDateStr)
		if err != nil {
			return nil, appErrors.ValidationError("Invalid end_date format. Must be RFC3339Nano format (e.g., 2025-03-08T15:05:42.536581Z)", err)
		}
		endDate = &parsedEndDate
	}

	// Validate date range if both dates are provided
	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return nil, appErrors.ValidationError("end_date must be greater than or equal to start_date", nil)
		}
	}

	return &DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// ParseDate parses a single date string in RFC3339Nano format
func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		return nil, appErrors.ValidationError("Invalid date format. Must be RFC3339Nano format (e.g., 2025-03-08T15:05:42.536581Z)", err)
	}

	return &parsedDate, nil
}
