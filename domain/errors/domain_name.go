package errors

import "errors"

var ErrDomainNameEmpty = errors.New("domain name cannot be empty")
var ErrDomainNameTooLong = errors.New("domain name is too long, must be less than 253 characters")
var ErrDomainNameInvalid = errors.New("domain name is invalid")
