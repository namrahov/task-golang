package util

import "unicode"

type IPasswordChecker interface {
	IsMiddleStrength(password string) bool
}

type PasswordChecker struct{}

// IsMiddleStrength validates if a password meets the required criteria:
// - Minimum 8 characters
// - At least one lowercase letter
// - At least one uppercase letter
// - At least one digit
func (p *PasswordChecker) IsMiddleStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower, hasUpper, hasDigit := false, false, false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasLower && hasUpper && hasDigit
}
