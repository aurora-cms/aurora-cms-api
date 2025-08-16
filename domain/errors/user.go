package errors

import "errors"

var ErrUserAlreadyOnTenant = errors.New("user is already on the tenant")
var ErrUserRoleEmpty = errors.New("user role cannot be empty")
var ErrUserRoleInvalid = errors.New("user role is invalid")
