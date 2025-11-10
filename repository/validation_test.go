package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		fieldName string
		wantErr   bool
		errType   error
	}{
		{
			name:      "valid ID with letters and numbers",
			id:        "abc123",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "valid ID with hyphens",
			id:        "test-id-123",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "valid ID with underscores",
			id:        "test_id_123",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "valid ID with dots",
			id:        "test.id.123",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "valid UUID-like ID",
			id:        "550e8400-e29b-41d4-a716-446655440000",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "empty ID",
			id:        "",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrEmptyInput,
		},
		{
			name:      "ID with dollar sign",
			id:        "test$injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with curly braces",
			id:        "test{injection}",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with square brackets",
			id:        "test[injection]",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with parentheses",
			id:        "test(injection)",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with asterisk",
			id:        "test*injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with plus sign",
			id:        "test+injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with question mark",
			id:        "test?injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with pipe",
			id:        "test|injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with caret",
			id:        "test^injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with backslash",
			id:        "test\\injection",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID starting with hyphen",
			id:        "-test123",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID starting with underscore",
			id:        "_test123",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID starting with dot",
			id:        ".test123",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with MongoDB operator at start",
			id:        "$gt",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "ID with MongoDB operator in middle",
			id:        "test$ne123",
			fieldName: "testId",
			wantErr:   true,
			errType:   ErrInvalidInput,
		},
		{
			name:      "valid numeric ID",
			id:        "123456789",
			fieldName: "testId",
			wantErr:   false,
		},
		{
			name:      "valid alphabetic ID",
			id:        "abcdefgh",
			fieldName: "testId",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id, tt.fieldName)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.Is(err, tt.errType), "expected error type %v, got %v", tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSortField(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   string
		expected string
	}{
		{
			name:     "valid datetime field",
			sortBy:   "datetime",
			expected: "datetime",
		},
		{
			name:     "valid name field",
			sortBy:   "name",
			expected: "person.name",
		},
		{
			name:     "valid email field",
			sortBy:   "email",
			expected: "person.email",
		},
		{
			name:     "valid person.name field",
			sortBy:   "person.name",
			expected: "person.name",
		},
		{
			name:     "valid person.email field",
			sortBy:   "person.email",
			expected: "person.email",
		},
		{
			name:     "invalid field defaults to datetime",
			sortBy:   "invalid_field",
			expected: "datetime",
		},
		{
			name:     "empty field defaults to datetime",
			sortBy:   "",
			expected: "datetime",
		},
		{
			name:     "dangerous field name defaults to datetime",
			sortBy:   "$where",
			expected: "datetime",
		},
		{
			name:     "injection attempt defaults to datetime",
			sortBy:   "'; DROP TABLE updates; --",
			expected: "datetime",
		},
		{
			name:     "MongoDB operator defaults to datetime",
			sortBy:   "$gt",
			expected: "datetime",
		},
		{
			name:     "field with special characters defaults to datetime",
			sortBy:   "field.name[0]",
			expected: "datetime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSortField(tt.sortBy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		size        int
		wantErr     bool
		expectedPage int
		expectedSize int
		errMsg      string
	}{
		{
			name:        "valid pagination",
			page:        1,
			size:        10,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:        "valid pagination with max page",
			page:        10000,
			size:        10,
			wantErr:     false,
			expectedPage: 10000,
			expectedSize: 10,
		},
		{
			name:        "valid pagination with max size",
			page:        1,
			size:        100,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 100,
		},
		{
			name:        "page below minimum defaults to 1",
			page:        0,
			size:        10,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:        "page negative defaults to 1",
			page:        -1,
			size:        10,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:        "size below minimum defaults to 1",
			page:        1,
			size:        0,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 1,
		},
		{
			name:        "size negative defaults to 1",
			page:        1,
			size:        -5,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 1,
		},
		{
			name:        "page exceeds maximum",
			page:        10001,
			size:        10,
			wantErr:     true,
			errMsg:      "page cannot exceed 10000",
		},
		{
			name:        "size exceeds maximum",
			page:        1,
			size:        101,
			wantErr:     true,
			errMsg:      "size cannot exceed 100",
		},
		{
			name:        "both page and size exceed maximum",
			page:        10001,
			size:        101,
			wantErr:     true,
			errMsg:      "page cannot exceed 10000",
		},
		{
			name:        "large page number",
			page:        9999,
			size:        50,
			wantErr:     false,
			expectedPage: 9999,
			expectedSize: 50,
		},
		{
			name:        "large size number",
			page:        1,
			size:        99,
			wantErr:     false,
			expectedPage: 1,
			expectedSize: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, size, err := ValidatePagination(tt.page, tt.size)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPage, page)
				assert.Equal(t, tt.expectedSize, size)
			}
		})
	}
}

