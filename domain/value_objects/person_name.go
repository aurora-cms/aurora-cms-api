package value_objects

import "github.com/h4rdc0m/aurora-api/domain/errors"

type PersonName struct {
	firstName string
	lastName  string
}

// NewPersonName creates a new PersonName instance with the provided first and last names.
// It returns an error if either name is empty.
func NewPersonName(firstName, lastName string) (*PersonName, error) {
	if firstName == "" {
		return nil, errors.ErrPersonFirstNameEmpty
	}
	if lastName == "" {
		return nil, errors.ErrPersonLastNameEmpty
	}
	return &PersonName{firstName: firstName, lastName: lastName}, nil
}

// FirstName returns the first name of the person.
func (pn PersonName) FirstName() string {
	return pn.firstName
}

// LastName returns the last name of the person.
func (pn PersonName) LastName() string {
	return pn.lastName
}

// FullName returns the full name of the person in the format "First Last".
func (pn PersonName) FullName() string {
	return pn.firstName + " " + pn.lastName
}

// String returns the full name of the person as a string.
func (pn PersonName) String() string {
	return pn.FullName()
}

// Equals checks if two PersonName instances are equal based on their first and last names.
func (pn PersonName) Equals(other PersonName) bool {
	return pn.firstName == other.firstName && pn.lastName == other.lastName
}
