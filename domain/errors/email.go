package errors

import "errors"

var ErrEmailEmpty = errors.New("email cannot be empty")
var ErrInvalidEmailFormat = errors.New("invalid email format")
