package validation

import "strings"

// IsValidEmail performs basic email format validation.
// Checks for non-empty local and domain parts with at least one dot in domain.
func IsValidEmail(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	parts := strings.Split(value, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}
