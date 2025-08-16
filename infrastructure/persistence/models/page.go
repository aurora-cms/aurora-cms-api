package models

type PageType string

const (
	PageTypeContent  PageType = "content"
	PageTypeLink     PageType = "link"
	PageTypeHardLink PageType = "hard_link"
	PageTypeSnippet  PageType = "snippet"
)

type Page struct {
	Base
	Key            string
	Path           *string
	Index          int
	Type           PageType
	ParentID       *uint64
	SiteID         uint64
	LinkURL        *string
	HardLinkPageID *uint64
}

type PageVersion struct {
	Base
	PageID      uint64
	Version     uint
	Title       string
	Description *string
	IsPublished bool
}

type PageBlock struct {
	Base
	PageVersionID uint64
	BlockKey      string
	Index         int
	ContentType   string
	Content       string
}
