package errors

import "errors"

var ErrPersonFirstNameEmpty = errors.New("person name cannot be empty")
var ErrPersonLastNameEmpty = errors.New("person name cannot be empty")
