package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAndValidateSort(t *testing.T) {
	allowed := map[string]struct{}{
		"name":       {},
		"created_at": {},
		"email":      {},
	}

	tests := []struct {
		name            string
		inputs          []string
		expectedValid   []string
		expectedInvalid []string
	}{
		{
			name:            "all valid",
			inputs:          []string{"name", "-created_at", "email"},
			expectedValid:   []string{"name", "-created_at", "email"},
			expectedInvalid: []string{},
		},
		{
			name:            "some invalid",
			inputs:          []string{"name", "invalid_field", "-created_at", "another_invalid"},
			expectedValid:   []string{"name", "-created_at"},
			expectedInvalid: []string{"invalid_field", "another_invalid"},
		},
		{
			name:            "all invalid",
			inputs:          []string{"invalid1", "invalid2"},
			expectedValid:   []string{},
			expectedInvalid: []string{"invalid1", "invalid2"},
		},
		{
			name:            "empty input",
			inputs:          []string{},
			expectedValid:   []string{},
			expectedInvalid: []string{},
		},
		{
			name:            "nil input",
			inputs:          nil,
			expectedValid:   []string{},
			expectedInvalid: []string{},
		},
		{
			name:            "descending sort",
			inputs:          []string{"-name", "-created_at"},
			expectedValid:   []string{"-name", "-created_at"},
			expectedInvalid: []string{},
		},
		{
			name:            "mixed case with invalid",
			inputs:          []string{"name", "-invalid", "created_at"},
			expectedValid:   []string{"name", "created_at"},
			expectedInvalid: []string{"invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, invalid := NormalizeAndValidateSort(tt.inputs, allowed)

			assert.ElementsMatch(t, tt.expectedValid, valid, "Valid fields should match")
			assert.ElementsMatch(t, tt.expectedInvalid, invalid, "Invalid fields should match")
		})
	}
}

func TestNormalizeAndValidateSort_EdgeCases(t *testing.T) {
	allowed := map[string]struct{}{
		"name": {},
	}

	t.Run("double dash prefix", func(t *testing.T) {
		valid, invalid := NormalizeAndValidateSort([]string{"--name"}, allowed)
		// "--name" trimmed becomes "-name" which is not in allowed
		assert.Empty(t, valid)
		assert.Contains(t, invalid, "-name")
	})

	t.Run("empty string in input", func(t *testing.T) {
		valid, invalid := NormalizeAndValidateSort([]string{""}, allowed)
		assert.Empty(t, valid)
		assert.Contains(t, invalid, "")
	})

	t.Run("whitespace handling", func(t *testing.T) {
		// Note: The current implementation doesn't trim whitespace
		// This test documents current behavior
		valid, invalid := NormalizeAndValidateSort([]string{" name"}, allowed)
		assert.Empty(t, valid)
		assert.Contains(t, invalid, " name")
	})
}
