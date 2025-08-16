package entities

import (
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"time"
)

type PageBlockID struct {
	value uint64
}

func NewPageBlockID(value uint64) PageBlockID {
	return PageBlockID{value: value}
}

func (id PageBlockID) Value() uint64 {
	return id.value
}

type PageBlock struct {
	id            PageBlockID
	pageVersionID PageVersionID
	blockKey      string
	index         int
	contentType   string
	content       string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPageBlock(pageVersionID PageVersionID, blockKey string, index int, contentType string, content string) (*PageBlock, error) {
	if blockKey == "" {
		return nil, errors.ErrInvalidBlockKey
	}
	if contentType == "" {
		return nil, errors.ErrInvalidContentType
	}

	now := time.Now()

	return &PageBlock{
		pageVersionID: pageVersionID,
		blockKey:      blockKey,
		index:         index,
		contentType:   contentType,
		content:       content,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func (pb *PageBlock) ID() PageBlockID {
	return pb.id
}

func (pb *PageBlock) PageVersionID() PageVersionID {
	return pb.pageVersionID
}

func (pb *PageBlock) BlockKey() string {
	return pb.blockKey
}

func (pb *PageBlock) Index() int {
	return pb.index
}

func (pb *PageBlock) ContentType() string {
	return pb.contentType
}

func (pb *PageBlock) Content() string {
	return pb.content
}

func (pb *PageBlock) CreatedAt() time.Time {
	return pb.createdAt
}

func (pb *PageBlock) UpdatedAt() time.Time {
	return pb.updatedAt
}

// UpdateBlockKey updates the block key
func (pb *PageBlock) UpdateBlockKey(blockKey string) error {
	if blockKey == "" {
		return errors.ErrInvalidBlockKey
	}

	pb.blockKey = blockKey
	pb.updatedAt = time.Now()
	return nil
}

// UpdateIndex updates the block index
func (pb *PageBlock) UpdateIndex(index int) {
	pb.index = index
	pb.updatedAt = time.Now()
}

// UpdateContentType updates the content type
func (pb *PageBlock) UpdateContentType(contentType string) error {
	if contentType == "" {
		return errors.ErrInvalidContentType
	}

	pb.contentType = contentType
	pb.updatedAt = time.Now()
	return nil
}

// UpdateContent updates the block content
func (pb *PageBlock) UpdateContent(content string) {
	pb.content = content
	pb.updatedAt = time.Now()
}

// SetID sets the page block ID (used by repository when loading from database)
func (pb *PageBlock) SetID(id PageBlockID) {
	pb.id = id
}

// SetTimestamps sets the timestamps (used by repository when loading from database)
func (pb *PageBlock) SetTimestamps(createdAt, updatedAt time.Time) {
	pb.createdAt = createdAt
	pb.updatedAt = updatedAt
}
