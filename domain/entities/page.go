package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"time"
)

// PageID represents a unique identifier for a page entity.
type PageID struct {
	value uint64
}

// NewPageID creates and returns a new PageID initialized with the given id value.
func NewPageID(id uint64) PageID {
	return PageID{value: id}
}

// Value returns the underlying unsigned integer value of the PageID.
func (p PageID) Value() uint64 {
	return p.value
}

func (p PageID) ValuePtr() *uint64 {
	return &p.value
}

// IsEmpty checks if the PageID is empty (i.e., has a value of 0).
func (p PageID) IsEmpty() bool {
	return p.value == 0
}

type PageType string

const (
	PageTypeContent  PageType = "content"  // Represents a content page type.
	PageTypeLink     PageType = "link"     // Represents a link page type.
	PageTypeHardLink PageType = "hardlink" // Represents a hard link page type.
	PageTypeSnippet  PageType = "snippet"  // Represents a snippet page type.
)

// Page represents a page entity within a site
type Page struct {
	id             PageID
	key            *value_objects.PageKey
	path           *string
	index          int
	parentID       *PageID
	siteID         SiteID
	pageType       PageType
	linkURL        *string
	hardLinkPageID *PageID
	createdAt      time.Time
	updatedAt      time.Time
	children       []*Page
	versions       []*PageVersion
}

func NewPage(key *value_objects.PageKey, path *string, siteID SiteID, pageType PageType) (*Page, error) {
	if key == nil {
		return nil, errors.ErrPageKeyEmpty
	}

	now := time.Now()

	return &Page{
		key:       key,
		path:      path,
		index:     0,
		siteID:    siteID,
		pageType:  pageType,
		createdAt: now,
		updatedAt: now,
		children:  make([]*Page, 0),
		versions:  make([]*PageVersion, 0),
	}, nil
}

// ID returns the unique identifier of the Page.
func (p *Page) ID() PageID {
	return p.id
}

// Key returns the page key
func (p *Page) Key() *value_objects.PageKey {
	return p.key
}

// Path returns the page path
func (p *Page) Path() *string {
	return p.path
}

func (p *Page) FullPath() string {
	if p.path != nil {
		return *p.path + "/" + p.key.Value()
	}
	return p.key.Value()
}

// Index returns the page index
func (p *Page) Index() int {
	return p.index
}

// ParentID returns the parent page ID
func (p *Page) ParentID() *PageID {
	return p.parentID
}

// SiteID returns the unique identifier of the site associated with the page.
func (p *Page) SiteID() SiteID {
	return p.siteID
}

// Type returns the page type
func (p *Page) Type() PageType {
	return p.pageType
}

// LinkURL returns the link URL (for link type pages)
func (p *Page) LinkURL() *string {
	return p.linkURL
}

// HardLinkPageID returns the hard link page ID (for hard link type pages)
func (p *Page) HardLinkPageID() *PageID {
	return p.hardLinkPageID
}

// CreatedAt returns the creation timestamp of the Page.
func (p *Page) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the timestamp of the last update to the Page.
func (p *Page) UpdatedAt() time.Time {
	return p.updatedAt
}

// Versions returns page versions
func (p *Page) Versions() []*PageVersion {
	return p.versions
}

// UpdateKey updates the page key
func (p *Page) UpdateKey(key *value_objects.PageKey) error {
	if key == nil {
		return errors.ErrPageKeyEmpty
	}

	p.key = key
	return nil
}

// UpdatePath updates the page path
func (p *Page) UpdatePath(path *string) {
	p.path = path
}

// UpdateIndex updates the page index
func (p *Page) UpdateIndex(index int) {
	p.index = index
	p.updatedAt = time.Now()
}

// SetParent sets the parent page
func (p *Page) SetParent(parentID *PageID) {
	p.parentID = parentID
	p.updatedAt = time.Now()
}

// UpdateType updates the page type
func (p *Page) UpdateType(pageType PageType) {
	p.pageType = pageType
	p.updatedAt = time.Now()
}

// SetLinkURL sets the link URL (for link type pages)
func (p *Page) SetLinkURL(linkURL *string) error {
	if p.pageType != PageTypeLink {
		return errors.ErrPageTypeInvalid
	}
	p.linkURL = linkURL
	p.updatedAt = time.Now()
	return nil
}

// SetHardLinkPageID sets the hard link page ID (for hard link type pages)
func (p *Page) SetHardLinkPageID(pageID *PageID) error {
	if p.pageType != PageTypeHardLink {
		return errors.ErrPageTypeInvalid
	}
	p.hardLinkPageID = pageID
	p.updatedAt = time.Now()
	return nil
}

// AddChild adds a child page
func (p *Page) AddChild(child *Page) error {
	if child == nil {
		return errors.ErrSitePageEmpty
	}

	child.SetParent(&p.id)
	p.children = append(p.children, child)
	p.updatedAt = time.Now()
	return nil
}

// AddVersion adds a new page version
func (p *Page) AddVersion(version *PageVersion) error {
	if version == nil {
		return errors.ErrPageVersionEmpty
	}

	p.versions = append(p.versions, version)
	p.updatedAt = time.Now()
	return nil
}

// SetID updates the Page's unique identifier with the provided PageID.
func (p *Page) SetID(id PageID) {
	p.id = id
}

// SetTimestamps updates the createdAt and updatedAt timestamps of the Page.
func (p *Page) SetTimestamps(createdAt, updatedAt time.Time) {
	p.createdAt = createdAt
	p.updatedAt = updatedAt
}
