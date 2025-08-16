package models

import (
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleTenantAdmin  UserRole = "tenant_admin"
	RoleTenantEditor UserRole = "tenant_editor"
	RoleUser         UserRole = "user"
)

type User struct {
	Base
	KeycloakID uuid.UUID
	Role       UserRole
	Tenants    []Tenant
}
