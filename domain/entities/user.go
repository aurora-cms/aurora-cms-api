package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"time"
)

// UserID represents a user identifier
type UserID struct {
	value uint64
}

// NewUserID creates a new UserID
func NewUserID(id uint64) UserID {
	return UserID{value: id}
}

// Value returns the ID value
func (u UserID) Value() uint64 {
	return u.value
}

// User represents a user aggregate root
type User struct {
	id         UserID
	keycloakID *value_objects.KeycloakID
	role       *value_objects.UserRole
	createdAt  time.Time
	updatedAt  time.Time
	tenantIDs  []TenantID
}

// NewUser creates a new User aggregate
func NewUser(keycloakID *value_objects.KeycloakID, role *value_objects.UserRole) (*User, error) {
	if keycloakID == nil {
		return nil, errors.ErrKeycloakIDEmpty
	}

	if role == nil {
		return nil, errors.ErrUserRoleEmpty
	}

	now := time.Now()

	return &User{
		keycloakID: keycloakID,
		role:       role,
		createdAt:  now,
		updatedAt:  now,
		tenantIDs:  make([]TenantID, 0),
	}, nil
}

// ID returns the user ID
func (u *User) ID() UserID {
	return u.id
}

// KeycloakID returns the Keycloak ID
func (u *User) KeycloakID() *value_objects.KeycloakID {
	return u.keycloakID
}

// Role returns the user role
func (u *User) Role() *value_objects.UserRole {
	return u.role
}

// CreatedAt returns the creation time
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update time
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// TenantIDs returns the tenant IDs the user has access to
func (u *User) TenantIDs() []TenantID {
	return u.tenantIDs
}

func (u *User) UpdateKeycloakID(keycloakID *value_objects.KeycloakID) error {
	if keycloakID == nil {
		return errors.ErrKeycloakIDEmpty
	}

	if !keycloakID.Equals(*u.keycloakID) {
		u.keycloakID = keycloakID
	}

	return nil
}

func (u *User) UpdateRole(role *value_objects.UserRole) error {
	if role == nil {
		return errors.ErrUserRoleEmpty
	}

	if !role.Equals(*u.role) {
		u.role = role
	}

	return nil
}

// AddToTenant adds the user to a tenant
func (u *User) AddToTenant(tenantID TenantID) error {
	// Check if the user is already in this tenant
	for _, existingTenantID := range u.tenantIDs {
		if existingTenantID.Value() == tenantID.Value() {
			return errors.ErrUserAlreadyOnTenant
		}
	}

	u.tenantIDs = append(u.tenantIDs, tenantID)
	return nil
}

// RemoveFromTenant removes the user from a tenant
func (u *User) RemoveFromTenant(tenantID TenantID) {
	for i, existingTenantID := range u.tenantIDs {
		if existingTenantID.Value() == tenantID.Value() {
			u.tenantIDs = append(u.tenantIDs[:i], u.tenantIDs[i+1:]...)
			break
		}
	}
}

// HasAccessToTenant checks if the user has access to a specific tenant
func (u *User) HasAccessToTenant(tenantID TenantID) bool {
	for _, existingTenantID := range u.tenantIDs {
		if existingTenantID.Value() == tenantID.Value() {
			return true
		}
	}
	return false
}

// CanManageTenant checks if the user can manage a specific tenant
func (u *User) CanManageTenant(tenantID TenantID) bool {
	return u.role.CanManageTenant() && u.HasAccessToTenant(tenantID)
}

// CanEditContent checks if the user can edit content in a specific tenant
func (u *User) CanEditContent(tenantID TenantID) bool {
	return u.role.CanEditContent() && u.HasAccessToTenant(tenantID)
}

// SetID sets the user ID (used by repository when loading from the database)
func (u *User) SetID(id UserID) {
	u.id = id
}

// SetTimestamps sets the timestamps (used by repository when loading from the database)
func (u *User) SetTimestamps(createdAt, updatedAt time.Time) {
	u.createdAt = createdAt
	u.updatedAt = updatedAt
}
