package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"time"
)

type PageVersionID struct {
	value uint64
}

func NewPageVersionID(id uint64) PageVersionID {
	return PageVersionID{value: id}
}

func (p PageVersionID) Value() uint64 {
	return p.value
}

type PageVersion struct {
	id          PageVersionID
	pageID      PageID
	version     uint
	title       string
	description *string
	isPublished bool
	createdAt   time.Time
	updatedAt   time.Time
	blocks      []*PageBlock
}

func NewPageVersion(pageID PageID, version uint, title string, description *string) (*PageVersion, error) {
	if title == "" {
		return nil, errors.ErrPageVersionTitleEmpty
	}

	if version <= 0 {
		return nil, errors.ErrPageVersionInvalidVersion
	}

	now := time.Now()

	return &PageVersion{
		pageID:      pageID,
		version:     version,
		title:       title,
		description: description,
		isPublished: false,
		createdAt:   now,
		updatedAt:   now,
		blocks:      make([]*PageBlock, 0),
	}, nil
}

// ID returns the page version ID
func (p *PageVersion) ID() PageVersionID {
	return p.id
}

// PageID returns the page ID
func (p *PageVersion) PageID() PageID {
	return p.pageID
}

// Version returns the version number
func (p *PageVersion) Version() uint {
	return p.version
}

// Title returns the page version title
func (p *PageVersion) Title() string {
	return p.title
}

// Description returns the page version description
func (p *PageVersion) Description() *string {
	return p.description
}

// IsPublished returns whether the page version is published
func (p *PageVersion) IsPublished() bool {
	return p.isPublished
}

// CreatedAt returns the creation time
func (p *PageVersion) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update time
func (p *PageVersion) UpdatedAt() time.Time {
	return p.updatedAt
}

// Blocks returns the page blocks
func (p *PageVersion) Blocks() []*PageBlock {
	return p.blocks
}

// UpdateTitle updates the page version title
func (p *PageVersion) UpdateTitle(title string) error {
	if title == "" {
		return errors.ErrPageVersionTitleEmpty
	}

	p.title = title
	p.updatedAt = time.Now()

	return nil
}

// UpdateDescription updates the page version description
func (p *PageVersion) UpdateDescription(description *string) {
	p.description = description
	p.updatedAt = time.Now()
}

// Publish publishes the page version
func (p *PageVersion) Publish() {
	p.isPublished = true
	p.updatedAt = time.Now()
}

// Unpublish unpublishes the page version
func (p *PageVersion) Unpublish() {
	p.isPublished = false
	p.updatedAt = time.Now()
}

// AddBlock adds a page block
func (p *PageVersion) AddBlock(block *PageBlock) error {
	if block == nil {
		return errors.ErrPageBlockEmpty
	}

	p.blocks = append(p.blocks, block)
	p.updatedAt = time.Now()

	return nil
}

// RemoveBlock removes a page block and returns true if the block was found and removed
func (p *PageVersion) RemoveBlock(blockID PageBlockID) bool {
	for i, block := range p.blocks {
		if block != nil && block.ID().Value() == blockID.Value() {
			p.blocks = append(p.blocks[:i], p.blocks[i+1:]...)
			p.updatedAt = time.Now()
			return true
		}
	}
	return false
}

// SetID sets the page version ID (used by repository when loading from database)
func (p *PageVersion) SetID(id PageVersionID) {
	p.id = id
}

// SetTimestamps sets the timestamps (used by repository when loading from database)
func (p *PageVersion) SetTimestamps(createdAt, updatedAt time.Time) {
	p.createdAt = createdAt
	p.updatedAt = updatedAt
}
