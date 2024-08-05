package utils

import (
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		n              uint
		maximumBackoff float64
		expectedMax    time.Duration
	}{
		{0, 1000.0, 1000 * time.Millisecond},
		{1, 1000.0, 1000 * time.Millisecond},
		{2, 1000.0, 1000 * time.Millisecond},
		{3, 2000.0, 2000 * time.Millisecond},
		{10, 10000.0, 10000 * time.Millisecond},
	}

	for _, tt := range tests {
		duration := ExponentialBackoff(tt.n, tt.maximumBackoff)
		if duration > tt.expectedMax {
			t.Errorf("expected duration to be less than or equal to %v, got %v", tt.expectedMax, duration)
		}
	}
}

func TestGetNextPageURL(t *testing.T) {
	tests := []struct {
		linkHeader  string
		expectedURL string
	}{
		{"<https://api.github.com/repositories/1/commits?page=2>; rel=\"next\", <https://api.github.com/repositories/1/commits?page=3>; rel=\"last\"", "https://api.github.com/repositories/1/commits?page=2"},
		{"<https://api.github.com/repositories/1/commits?page=3>; rel=\"last\"", ""},
		{"", ""},
	}

	for _, tt := range tests {
		nextURL := GetNextPageURL(tt.linkHeader)
		if nextURL != tt.expectedURL {
			t.Errorf("expected URL %v, got %v", tt.expectedURL, nextURL)
		}
	}
}

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
			actual := ValidateDate(tc.input)
			if actual != tc.expected {
				t.Errorf("ValidateDate(%v) = %v; expected %v", tc.input, actual, tc.expected)
			}
		})
	}
}
