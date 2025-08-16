package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"time"
)

type TenantID struct {
	value uint64
}

// NewTenantID creates a new TenantID instance with the provided value.
func NewTenantID(id uint64) TenantID {
	return TenantID{value: id}
}

// Value returns the underlying value of the TenantID.
func (t TenantID) Value() uint64 {
	return t.value
}

// Tenant represents a tenant aggregate root
type Tenant struct {
	id               TenantID
	name             string
	description      *string
	isActive         bool
	isBillingEnabled bool
	createdAt        time.Time
	updatedAt        time.Time
	sites            []*Site
	usersIDs         []UserID
}

// NewTenant creates a new Tenant instance with the provided name, description, and active status.
// It returns an error if the name or description is empty.
func NewTenant(name string, description *string) (*Tenant, error) {
	if name == "" {
		return nil, errors.ErrTenantNameEmpty
	}

	now := time.Now()

	return &Tenant{
		name:             name,
		description:      description,
		isActive:         true,
		isBillingEnabled: false,
		createdAt:        now,
		updatedAt:        now,
		sites:            make([]*Site, 0),
		usersIDs:         make([]UserID, 0),
	}, nil
}

func (t *Tenant) ID() TenantID {
	return t.id
}

func (t *Tenant) Name() string {
	return t.name
}

func (t *Tenant) Description() *string {
	return t.description
}

func (t *Tenant) IsActive() bool {
	return t.isActive
}

func (t *Tenant) IsBillingEnabled() bool {
	return t.isBillingEnabled
}

func (t *Tenant) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Tenant) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tenant) Sites() []*Site {
	return t.sites
}

func (t *Tenant) UsersIDs() []UserID {
	return t.usersIDs
}

func (t *Tenant) UpdateName(name string) error {
	if name == "" {
		return errors.ErrTenantNameEmpty
	}

	t.name = name

	return nil
}

func (t *Tenant) UpdateDescription(description *string) error {
	if description == nil {
		return errors.ErrTenantDescriptionEmpty
	}

	t.description = description

	return nil
}

func (t *Tenant) Activate() {
	t.isActive = true
}

func (t *Tenant) Deactivate() {
	t.isActive = false
}

func (t *Tenant) EnableBilling() {
	t.isBillingEnabled = true
}

func (t *Tenant) DisableBilling() {
	t.isBillingEnabled = false
}

func (t *Tenant) AddSite(site *Site) error {
	if site == nil {
		return errors.ErrTenantSiteEmpty
	}

	for _, existingSite := range t.sites {
		if existingSite.Domain().Equals(*site.Domain()) {
			return errors.ErrTenantSiteDomainAlreadyExists
		}
	}

	t.sites = append(t.sites, site)

	return nil
}

func (t *Tenant) RemoveSite(siteID SiteID) error {
	for i, site := range t.sites {
		if site.ID().Value() == siteID.Value() {
			t.sites = append(t.sites[:i], t.sites[i+1:]...)
			return nil
		}
	}
	return errors.ErrTenantSiteNotFound // or a more specific error for site not found
}

func (t *Tenant) AddUser(userID UserID) error {
	for _, existingUserID := range t.usersIDs {
		if existingUserID.Value() == userID.Value() {
			return errors.ErrUserAlreadyOnTenant
		}
	}
	t.usersIDs = append(t.usersIDs, userID)
	return nil
}

func (t *Tenant) RemoveUser(userID UserID) error {
	for i, existingUserID := range t.usersIDs {
		if existingUserID.Value() == userID.Value() {
			t.usersIDs = append(t.usersIDs[:i], t.usersIDs[i+1:]...)
			return nil
		}
	}
	return errors.ErrUserNotFoundOnTenant
}

func (t *Tenant) SetID(id TenantID) {
	t.id = id
}

func (t *Tenant) SetTimestamps(createdAt, updatedAt time.Time) {
	t.createdAt = createdAt
	t.updatedAt = updatedAt
}
