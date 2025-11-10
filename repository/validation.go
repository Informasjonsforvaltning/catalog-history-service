package repository

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	// ErrInvalidInput indicates that the input contains invalid characters
	ErrInvalidInput = errors.New("input contains invalid characters")
	// ErrEmptyInput indicates that the input is empty
	ErrEmptyInput = errors.New("input cannot be empty")
)

// ValidateID validates that an ID parameter is safe to use in database queries.
// It ensures the ID is not empty and contains only safe characters.
// This prevents NoSQL injection attacks by rejecting potentially dangerous input.
func ValidateID(id, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%w: %s", ErrEmptyInput, fieldName)
	}
	
	// Check for dangerous MongoDB operators that could be used for injection
	dangerousPatterns := []string{"$", "{", "}", "[", "]", "(", ")", "*", "+", "?", "|", "^", "\\"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(id, pattern) {
			return fmt.Errorf("%w: %s contains potentially dangerous character: %s", ErrInvalidInput, fieldName, pattern)
		}
	}
	
	// Additional check: ensure the ID doesn't start with special characters
	if len(id) > 0 && !unicode.IsLetter(rune(id[0])) && !unicode.IsDigit(rune(id[0])) {
		return fmt.Errorf("%w: %s must start with a letter or digit", ErrInvalidInput, fieldName)
	}
	
	return nil
}

// ValidateSortField validates that a sort field name is whitelisted and safe.
// Only whitelisted field names are allowed to prevent injection attacks.
// Returns the validated field name (always safe, defaults to "datetime" if invalid).
func ValidateSortField(sortBy string) string {
	// Whitelist of allowed sort fields
	allowedFields := map[string]string{
		"datetime":     "datetime",
		"name":         "person.name",
		"email":        "person.email",
		"person.name":  "person.name",
		"person.email": "person.email",
	}
	
	if validated, ok := allowedFields[sortBy]; ok {
		return validated
	}
	
	// Default to datetime if not specified or invalid (safe fallback)
	return "datetime"
}

// ValidatePagination validates and limits pagination parameters to prevent DoS attacks.
func ValidatePagination(page, size int) (int, int, error) {
	// Set reasonable limits
	const maxPage = 10000
	const maxSize = 100
	const minPage = 1
	const minSize = 1
	
	if page < minPage {
		page = minPage
	}
	if page > maxPage {
		return 0, 0, fmt.Errorf("page cannot exceed %d", maxPage)
	}
	
	if size < minSize {
		size = minSize
	}
	if size > maxSize {
		return 0, 0, fmt.Errorf("size cannot exceed %d", maxSize)
	}
	
	return page, size, nil
}

