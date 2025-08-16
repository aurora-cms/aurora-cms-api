package value_objects

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object, validating the email format.
func NewEmail(value string) (*Email, error) {
	email := strings.TrimSpace(value)
	if email == "" {
		return nil, errors.ErrEmailEmpty
	}

	if !emailRegex.MatchString(email) {
		return nil, errors.ErrInvalidEmailFormat
	}

	return &Email{value: email}, nil
}

// Value returns the string representation of the Email.
func (e Email) Value() string {
	return e.value
}

// String returns the string representation of the Email.
// This method is used to implement the Stringer interface, allowing the Email to be printed directly
// in logs or other output formats.
func (e Email) String() string {
	return e.value
}

// Equals checks if two Email values are equal.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
