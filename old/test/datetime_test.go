package test

import (
	"testing"

	"github.com/Rhaqim/savannahtech/old/utils"
)

func TestValidateDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// Test with a valid RFC3339 date
		{
			input:    "2024-08-02T14:20:00Z",
			expected: "2024-08-02T14:20:00Z",
		},
		// Test with an empty string
		{
			input:    "",
			expected: "",
		},
		// Test with an invalid date format
		{
			input:    "2024-08-02 14:20:00",
			expected: "",
		},
		// Test with a different valid RFC3339 date
		{
			input:    "2023-01-01T00:00:00Z",
			expected: "2023-01-01T00:00:00Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := utils.ValidateDate(tc.input)
			if actual != tc.expected {
				t.Errorf("ValidateDate(%v) = %v; expected %v", tc.input, actual, tc.expected)
			}
		})
	}
}
