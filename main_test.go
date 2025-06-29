package main

import (
	"testing"
)

func TestPathToDatadogMetrics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path with uuid parameter",
			input:    "/bm/branding/buying-guide/landing-pages/<uuid:landing_page_uuid>",
			expected: "/bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid",
		},
		{
			name:     "path with multiple dynamic parameters",
			input:    "/api/v1/users/<int:user_id>/posts/<uuid:post_id>",
			expected: "/api/v1/users/_int:user_id/posts/_uuid:post_id",
		},
		{
			name:     "path without dynamic parameters",
			input:    "/api/v1/health/status",
			expected: "/api/v1/health/status",
		},
		{
			name:     "path with string parameter",
			input:    "/products/<string:category>/items",
			expected: "/products/_string:category/items",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "/",
		},
		{
			name:     "path with complex dynamic parameter",
			input:    "/api/<version:v2>/resources/<slug:resource_name>",
			expected: "/api/_version:v2/resources/_slug:resource_name",
		},
		{
			name:     "admin auth user path",
			input:    "/admin/auth/user",
			expected: "/admin/auth/user",
		},
		{
			name:     "path with nested parameters",
			input:    "/customer/<int:id>/order/<uuid:order_id>",
			expected: "/customer/_int:id/order/_uuid:order_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathToDatadogMetrics(tt.input)
			if result != tt.expected {
				t.Errorf("pathToDatadogMetrics(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name         string
		dateStr      string
		defaultValue int64
		expected     int64
	}{
		{
			name:         "valid timestamp",
			dateStr:      "1719525600000",
			defaultValue: 0,
			expected:     1719525600000,
		},
		{
			name:         "empty string returns default",
			dateStr:      "",
			defaultValue: 1234567890000,
			expected:     1234567890000,
		},
		{
			name:         "invalid timestamp returns default",
			dateStr:      "not-a-number",
			defaultValue: 9876543210000,
			expected:     9876543210000,
		},
		{
			name:         "negative timestamp",
			dateStr:      "-1000",
			defaultValue: 0,
			expected:     -1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.dateStr, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("parseDate(%q, %d) = %d, want %d", tt.dateStr, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestMetricResult(t *testing.T) {
	// Test that MetricResult struct contains expected fields
	mr := MetricResult{
		Metric:            "resource_name:get_admin/auth/user",
		LastSeenTimestamp: 1719525600000,
		LastSeenDate:      "28 June 2024 00:00:00 UTC",
		Value:             float64(100),
	}
	
	if mr.Metric != "resource_name:get_admin/auth/user" {
		t.Errorf("MetricResult.Metric = %v, want %v", mr.Metric, "resource_name:get_admin/auth/user")
	}
	
	if mr.LastSeenTimestamp != 1719525600000 {
		t.Errorf("MetricResult.LastSeenTimestamp = %v, want %v", mr.LastSeenTimestamp, 1719525600000)
	}
	
	if mr.LastSeenDate != "28 June 2024 00:00:00 UTC" {
		t.Errorf("MetricResult.LastSeenDate = %v, want %v", mr.LastSeenDate, "28 June 2024 00:00:00 UTC")
	}
	
	if mr.Value != float64(100) {
		t.Errorf("MetricResult.Value = %v, want %v", mr.Value, float64(100))
	}
}

func TestResult(t *testing.T) {
	// Test Result wrapper struct
	result := Result{
		Result: []MetricResult{
			{
				Metric:            "resource_name:get_admin/auth/user",
				LastSeenTimestamp: 1719525600000,
				LastSeenDate:      "28 June 2024 00:00:00 UTC",
				Value:             float64(100),
			},
			{
				Metric:            "resource_name:post_admin/auth/user",
				LastSeenTimestamp: 1719439200000,
				LastSeenDate:      "27 June 2024 00:00:00 UTC",
				Value:             float64(50),
			},
		},
	}
	
	if len(result.Result) != 2 {
		t.Errorf("Result should contain 2 MetricResults, got %d", len(result.Result))
	}
	
	if result.Result[0].Metric != "resource_name:get_admin/auth/user" {
		t.Errorf("First result metric = %v, want %v", result.Result[0].Metric, "resource_name:get_admin/auth/user")
	}
	
	if result.Result[1].Value != float64(50) {
		t.Errorf("Second result value = %v, want %v", result.Result[1].Value, float64(50))
	}
}