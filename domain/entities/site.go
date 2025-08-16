package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"time"
)

// SiteID represents a unique identifier for a site entity.
// It encapsulates an unsigned integer value as its underlying data.
type SiteID struct {
	value uint64
}

// NewSiteID creates a new SiteID instance with the specified unsigned integer value.
func NewSiteID(id uint64) SiteID {
	return SiteID{value: id}
}

// Value retrieves the internal `value` field of the SiteID.
func (s SiteID) Value() uint64 {
	return s.value
}

// IsEmpty checks if the SiteID is empty, which is defined as having a value of 0.
func (s SiteID) IsEmpty() bool {
	return s.value == 0
}

// Site represents a web platform or application containing various pages managed by a tenant.
type Site struct {
	id            SiteID
	name          string
	description   *string
	domain        *value_objects.DomainName
	titleTemplate *string
	enabled       bool
	templateID    TemplateID
	tenantID      TenantID
	createdAt     time.Time
	updatedAt     time.Time
	pages         []*Page
}

// NewSite creates a new Site entity
func NewSite(name string, description *string, domain *value_objects.DomainName, templateID TemplateID, tenantID TenantID) (*Site, error) {
	if name == "" {
		return nil, errors.ErrSiteNameEmpty
	}

	if domain == nil {
		return nil, errors.ErrDomainNameEmpty
	}

	now := time.Now()

	return &Site{
		name:        name,
		description: description,
		domain:      domain,
		enabled:     true,
		templateID:  templateID,
		tenantID:    tenantID,
		createdAt:   now,
		updatedAt:   now,
		pages:       make([]*Page, 0),
	}, nil
}

// ID returns the unique identifier of the Site as a SiteID.
func (s *Site) ID() SiteID {
	return s.id
}

// Name returns the name of the Site as a string.
func (s *Site) Name() string {
	return s.name
}

// Description returns the description of the Site.
func (s *Site) Description() *string {
	return s.description
}

// Domain retrieves the DomainName associated with the Site.
func (s *Site) Domain() *value_objects.DomainName {
	return s.domain
}

// TitleTemplate retrieves the title template of the Site, which can be used for generating page titles.
func (s *Site) TitleTemplate() *string {
	return s.titleTemplate
}

// IsEnabled determines whether the site is currently active and returns true if active, otherwise false.
func (s *Site) IsEnabled() bool {
	return s.enabled
}

// TemplateID retrieves the unique identifier of the template associated with the Site.
func (s *Site) TemplateID() TemplateID {
	return s.templateID
}

// TenantID retrieves the unique identifier of the tenant associated with the site.
func (s *Site) TenantID() TenantID {
	return s.tenantID
}

// CreatedAt retrieves the creation timestamp of the Site instance.
func (s *Site) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt retrieves the timestamp of the most recent update made to the Site.
func (s *Site) UpdatedAt() time.Time {
	return s.updatedAt
}

// Pages returns a slice of pointers to the Page objects associated with the Site.
func (s *Site) Pages() []*Page {
	return s.pages
}

// UpdateName updates the name of the Site. Returns an error if the provided name is empty.
func (s *Site) UpdateName(name string) error {
	if name == "" {
		return errors.ErrSiteNameEmpty
	}

	s.name = name

	return nil
}

// UpdateDescription updates the description of the Site with the provided value and returns an error if the operation fails.
func (s *Site) UpdateDescription(description *string) error {
	s.description = description

	return nil
}

// UpdateDomain updates the domain of a Site. Returns an error if the provided domain is nil or has an empty value.
func (s *Site) UpdateDomain(domain *value_objects.DomainName) error {
	if domain == nil || domain.Value() == "" {
		return errors.ErrSiteDomainEmpty
	}

	s.domain = domain

	return nil
}

// UpdateTitleTemplate updates the title template
func (s *Site) UpdateTitleTemplate(titleTemplate *string) {
	s.titleTemplate = titleTemplate
	s.updatedAt = time.Now()
}

// Enable sets the `isActive` field of the Site to `true`, marking the site as active.
func (s *Site) Enable() {
	s.enabled = true
}

// Disable sets the site's isActive property to false, marking it as inactive.
func (s *Site) Disable() {
	s.enabled = false
}

// UpdateTemplate updates the template ID
func (s *Site) UpdateTemplate(templateID TemplateID) {
	s.templateID = templateID
}

// AddPage adds a new page to the Site. Returns an error if the page is nil or if the page's slug already exists.
func (s *Site) AddPage(page *Page) error {
	if page == nil {
		return errors.ErrSitePageEmpty
	}

	for _, existingPage := range s.pages {
		if existingPage.Key().Equals(*page.Key()) {
			return errors.ErrSitePageWithSlugAlreadyExists
		}
	}
	s.pages = append(s.pages, page)

	return nil
}

// RemovePage removes a page from the Site by its PageID. Returns an error if the page is not found in the Site's pages.
func (s *Site) RemovePage(pageID PageID) error {
	for i, page := range s.pages {
		if page.ID().Value() == pageID.Value() {
			s.pages = append(s.pages[:i], s.pages[i+1:]...)
			return nil
		}
	}
	return errors.ErrSitePageNotFound // or a more specific error for page not found
}

// SetID updates the unique identifier of the Site with the given SiteID. It returns an error if the operation fails.
func (s *Site) SetID(id SiteID) error {
	s.id = id
	return nil
}

// SetTimestamps sets the creation and last updated timestamps for the Site instance.
func (s *Site) SetTimestamps(createdAt, updatedAt time.Time) {
	s.createdAt = createdAt
	s.updatedAt = updatedAt
}
