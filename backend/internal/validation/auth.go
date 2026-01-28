package validation

import "strings"

// BearerToken extracts the token from an Authorization header value.
// Returns empty string if the value is not a valid Bearer token format.
func BearerToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
