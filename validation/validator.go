package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

// ValidateStringLength check if the string length is between minLength and maxLength
func ValidateStringLength(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("string length is invalid, must be between %d and %d", minLength, maxLength)
	}

	return nil
}

// ValidateUsername check if the username is valid.
// It must be between 3 and 50 characters long
// and contain only letters, numbers and underscores.
func ValidateUsername(value string) error {
	if err := ValidateStringLength(value, 3, 50); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username is invalid, must contain only letters, numbers and underscores")
	}

	return nil
}

// ValidatePassword check if the password is valid.
// It must be between 6 and 100 characters long.
func ValidatePassword(value string) error {
	return ValidateStringLength(value, 6, 100)
}

// ValidateEmail check if the email is valid.
// It must be between 5 and 250 characters long
// and contain a valid email address.
func ValidateEmail(value string) error {
	if err := ValidateStringLength(value, 5, 250); err != nil {
		return err
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("email is invalid")
	}

	return nil
}
