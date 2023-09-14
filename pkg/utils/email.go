package utils

import "strings"

// Normalize Email
func Normalize(email string) string {
	email = strings.TrimSpace(email)
	email = strings.TrimRight(email, ".")
	email = strings.ToLower(email)
	return email
}
