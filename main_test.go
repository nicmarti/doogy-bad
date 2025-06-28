package main

import (
	"reflect"
	"testing"
)

func TestPathToDatadogMetrics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "path with uuid parameter",
			input: "/bm/branding/buying-guide/landing-pages/<uuid:landing_page_uuid>",
			expected: []string{
				"get_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid",
				"post_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid",
				"options_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid",
			},
		},
		{
			name:  "path with multiple dynamic parameters",
			input: "/api/v1/users/<int:user_id>/posts/<uuid:post_id>",
			expected: []string{
				"get_api/v1/users/_int:user_id/posts/_uuid:post_id",
				"post_api/v1/users/_int:user_id/posts/_uuid:post_id",
				"options_api/v1/users/_int:user_id/posts/_uuid:post_id",
			},
		},
		{
			name:  "path without dynamic parameters",
			input: "/api/v1/health/status",
			expected: []string{
				"get_api/v1/health/status",
				"post_api/v1/health/status",
				"options_api/v1/health/status",
			},
		},
		{
			name:  "path with string parameter",
			input: "/products/<string:category>/items",
			expected: []string{
				"get_products/_string:category/items",
				"post_products/_string:category/items",
				"options_products/_string:category/items",
			},
		},
		{
			name:  "root path",
			input: "/",
			expected: []string{
				"get_",
				"post_",
				"options_",
			},
		},
		{
			name:  "path with complex dynamic parameter",
			input: "/api/<version:v2>/resources/<slug:resource_name>",
			expected: []string{
				"get_api/_version:v2/resources/_slug:resource_name",
				"post_api/_version:v2/resources/_slug:resource_name",
				"options_api/_version:v2/resources/_slug:resource_name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathToDatadogMetrics(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("pathToDatadogMetrics(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}