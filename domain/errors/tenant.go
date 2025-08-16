package errors

import "errors"

var ErrTenantNameEmpty = errors.New("tenant name cannot be empty")
var ErrTenantDescriptionEmpty = errors.New("tenant description cannot be empty")
var ErrTenantSiteDomainAlreadyExists = errors.New("site domain already exists in tenant")
var ErrTenantSiteEmpty = errors.New("site cannot be empty")
var ErrTenantSiteNotFound = errors.New("site not found in tenant")
var ErrUserNotFoundOnTenant = errors.New("user not found on tenant")
var ErrTenantNotFound = errors.New("tenant not found")
