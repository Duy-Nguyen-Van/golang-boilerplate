package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateRange(t *testing.T) {
	tests := []struct {
		name          string
		startDateStr  string
		endDateStr    string
		expectError   bool
		expectStart   bool
		expectEnd     bool
		errorContains string
	}{
		{
			name:         "both dates valid",
			startDateStr: "2025-01-01T00:00:00.000000000Z",
			endDateStr:   "2025-01-31T23:59:59.999999999Z",
			expectError:  false,
			expectStart:  true,
			expectEnd:    true,
		},
		{
			name:         "only start date",
			startDateStr: "2025-01-01T00:00:00.000000000Z",
			endDateStr:   "",
			expectError:  false,
			expectStart:  true,
			expectEnd:    false,
		},
		{
			name:         "only end date",
			startDateStr: "",
			endDateStr:   "2025-01-31T23:59:59.999999999Z",
			expectError:  false,
			expectStart:  false,
			expectEnd:    true,
		},
		{
			name:         "no dates",
			startDateStr: "",
			endDateStr:   "",
			expectError:  false,
			expectStart:  false,
			expectEnd:    false,
		},
		{
			name:          "invalid start date format",
			startDateStr:  "2025-01-01",
			endDateStr:    "",
			expectError:   true,
			errorContains: "Invalid start_date format",
		},
		{
			name:          "invalid end date format",
			startDateStr:  "",
			endDateStr:    "2025-01-31",
			expectError:   true,
			errorContains: "Invalid end_date format",
		},
		{
			name:          "end date before start date",
			startDateStr:  "2025-01-31T00:00:00.000000000Z",
			endDateStr:    "2025-01-01T00:00:00.000000000Z",
			expectError:   true,
			errorContains: "end_date must be greater than or equal to start_date",
		},
		{
			name:         "end date equal to start date",
			startDateStr: "2025-01-01T00:00:00.000000000Z",
			endDateStr:   "2025-01-01T00:00:00.000000000Z",
			expectError:  false,
			expectStart:  true,
			expectEnd:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDateRange(tt.startDateStr, tt.endDateStr)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.expectStart {
				assert.NotNil(t, result.StartDate)
			} else {
				assert.Nil(t, result.StartDate)
			}

			if tt.expectEnd {
				assert.NotNil(t, result.EndDate)
			} else {
				assert.Nil(t, result.EndDate)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name          string
		dateStr       string
		expectError   bool
		expectNil     bool
		errorContains string
	}{
		{
			name:        "valid date",
			dateStr:     "2025-01-01T00:00:00.000000000Z",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "empty string",
			dateStr:     "",
			expectError: false,
			expectNil:   true,
		},
		{
			name:          "invalid format",
			dateStr:       "2025-01-01",
			expectError:   true,
			errorContains: "Invalid date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDate(tt.dateStr)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// Verify it's a valid time
				assert.IsType(t, &time.Time{}, result)
			}
		})
	}
}

func TestParseDateRange_TimeValues(t *testing.T) {
	startDateStr := "2025-01-01T12:30:45.123456789Z"
	endDateStr := "2025-01-31T18:45:30.987654321Z"

	result, err := ParseDateRange(startDateStr, endDateStr)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.StartDate)
	require.NotNil(t, result.EndDate)

	// Verify the parsed times are correct
	expectedStart, _ := time.Parse(time.RFC3339Nano, startDateStr)
	expectedEnd, _ := time.Parse(time.RFC3339Nano, endDateStr)

	assert.Equal(t, expectedStart, *result.StartDate)
	assert.Equal(t, expectedEnd, *result.EndDate)
	assert.True(t, result.EndDate.After(*result.StartDate))
}
