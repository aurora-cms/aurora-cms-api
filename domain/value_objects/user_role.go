package value_objects

import "github.com/h4rdc0m/aurora-api/domain/errors"

type UserRole struct {
	value string
}

const (
	RoleSuperAdmin   = "super_admin"
	RoleAdmin        = "admin"
	RoleTenantAdmin  = "tenant_admin"
	RoleTenantEditor = "tenant_editor"
	RoleUser         = "user"
)

var validRoles = map[string]bool{
	RoleSuperAdmin:   true,
	RoleAdmin:        true,
	RoleTenantAdmin:  true,
	RoleTenantEditor: true,
	RoleUser:         true,
}

func NewUserRole(role string) (*UserRole, error) {
	if role == "" {
		return nil, errors.ErrUserRoleEmpty
	}

	if !validRoles[role] {
		return nil, errors.ErrUserRoleInvalid
	}

	return &UserRole{value: role}, nil
}

func (r UserRole) Value() string {
	return r.value
}

func (r UserRole) String() string {
	return r.value
}

func (r UserRole) Equals(other UserRole) bool {
	return r.value == other.value
}

func (r UserRole) IsSuperAdmin() bool {
	return r.value == RoleSuperAdmin
}

func (r UserRole) IsAdmin() bool {
	return r.value == RoleAdmin
}

func (r UserRole) IsTenantAdmin() bool {
	return r.value == RoleTenantAdmin
}

func (r UserRole) IsTenantEditor() bool {
	return r.value == RoleTenantEditor
}

func (r UserRole) IsUser() bool {
	return r.value == RoleUser
}

func (r UserRole) CanManageTenant() bool {
	return r.IsSuperAdmin() || r.IsAdmin() || r.IsTenantAdmin()
}

func (r UserRole) CanEditContent() bool {
	return r.IsSuperAdmin() || r.IsAdmin() || r.IsTenantAdmin() || r.IsTenantEditor()
}
