package repository

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrInvalidInput = errors.New("input contains invalid characters")
	ErrEmptyInput   = errors.New("input cannot be empty")
)

func ValidateID(id, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%w: %s", ErrEmptyInput, fieldName)
	}

	dangerousPatterns := []string{"$", "{", "}", "[", "]", "(", ")", "*", "+", "?", "|", "^", "\\"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(id, pattern) {
			return fmt.Errorf("%w: %s contains potentially dangerous character: %s", ErrInvalidInput, fieldName, pattern)
		}
	}

	if len(id) > 0 && !unicode.IsLetter(rune(id[0])) && !unicode.IsDigit(rune(id[0])) {
		return fmt.Errorf("%w: %s must start with a letter or digit", ErrInvalidInput, fieldName)
	}

	return nil
}

func ValidateSortField(sortBy string) string {
	allowedFields := map[string]string{
		"datetime":     "datetime",
		"name":         "person_name",
		"email":        "person_email",
		"person.name":  "person_name",
		"person.email": "person_email",
		"person_name":  "person_name",
		"person_email": "person_email",
	}

	if validated, ok := allowedFields[sortBy]; ok {
		return validated
	}

	return "datetime"
}

func ValidatePagination(page, size int) (int, int, error) {
	const maxPage = 10000
	const maxSize = 100
	const minPage = 0
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
