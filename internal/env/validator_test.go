package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateName(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		envVarName  string
		expectError bool
	}{
		{"valid name", "FOO", false},
		{"valid name with underscore", "FOO_BAR", false},
		{"valid name with numbers", "FOO123", false},
		{"empty name", "", true},
		{"starts with number", "1FOO", true},
		{"contains hyphen", "FOO-BAR", true},
		{"contains space", "FOO BAR", true},
		{"reserved name", "PATH", true},
		{"reserved name case insensitive", "path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateName(tt.envVarName)
			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, &ValidationError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateValue(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid value", "foo", false},
		{"empty value", "", false},
		{"value with spaces", "foo bar", false},
		{"value with special chars", "foo@bar", false},
		{"value with null byte", "foo\x00bar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateValue(tt.value)
			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, &ValidationError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateMap(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		envMap      map[string]string
		expectError bool
	}{
		{
			"valid map",
			map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
			false,
		},
		{
			"invalid name",
			map[string]string{
				"FOO-BAR": "baz",
			},
			true,
		},
		{
			"invalid value",
			map[string]string{
				"FOO": "bar\x00baz",
			},
			true,
		},
		{
			"reserved name",
			map[string]string{
				"PATH": "/custom/path",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateMap(tt.envMap)
			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, &ValidationError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
