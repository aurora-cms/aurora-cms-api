package value_objects

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"regexp"
	"strings"
)

type DomainName struct {
	value string
}

var domainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// NewDomainName creates a new DomainName value object.
// It validates the domain name against a regex pattern.
func NewDomainName(value string) (*DomainName, error) {
	domain := strings.TrimSpace(strings.ToLower(value))
	if domain == "" {
		return nil, errors.ErrDomainNameEmpty
	}

	if len(domain) > 253 {
		return nil, errors.ErrDomainNameTooLong
	}

	if !domainNameRegex.MatchString(domain) {
		return nil, errors.ErrDomainNameInvalid
	}

	return &DomainName{value: domain}, nil
}

// Value returns the string representation of the DomainName.
func (d DomainName) Value() string {
	return d.value
}

// String returns the string representation of the DomainName.
// It implements the fmt.Stringer interface.
func (d DomainName) String() string {
	return d.value
}

// Equals checks if two DomainName instances are equal.
func (d DomainName) Equals(other DomainName) bool {
	return d.value == other.value
}
