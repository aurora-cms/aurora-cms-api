package mappers

import (
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/stretchr/testify/assert"
)

func TestPageMapper_ToModel(t *testing.T) {
	mapper := NewPageMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *entities.Page
		want      *models.Page
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid content page",
			input: func() *entities.Page {
				key, _ := value_objects.NewPageKey("content-key")
				page, _ := entities.NewPage(key, value_objects.NewNullableString("/content-path").Value(), entities.NewSiteID(1), entities.PageTypeContent)
				page.SetID(entities.NewPageID(1))
				page.SetTimestamps(now, now)
				page.UpdateIndex(1)
				return page
			}(),
			want: &models.Page{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:      "content-key",
				Path:     value_objects.NewNullableString("/content-path").Value(),
				Index:    1,
				SiteID:   1,
				Type:     models.PageTypeContent,
				LinkURL:  nil,
				ParentID: nil,
			},
			expectErr: false,
		},
		{
			name: "valid link page with URL",
			input: func() *entities.Page {
				key, _ := value_objects.NewPageKey("link-key")
				page, _ := entities.NewPage(key, value_objects.NewNullableString("/link-path").Value(), entities.NewSiteID(1), entities.PageTypeLink)
				page.SetID(entities.NewPageID(2))
				page.SetTimestamps(now, now)
				page.UpdateIndex(2)
				_ = page.SetLinkURL(value_objects.NewNullableString("https://example.com").Value())
				return page
			}(),
			want: &models.Page{
				Base: models.Base{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:      "link-key",
				Path:     value_objects.NewNullableString("/link-path").Value(),
				Index:    2,
				SiteID:   1,
				Type:     models.PageTypeLink,
				LinkURL:  value_objects.NewNullableString("https://example.com").Value(),
				ParentID: nil,
			},
			expectErr: false,
		},
		{
			name: "page with parent",
			input: func() *entities.Page {
				key, _ := value_objects.NewPageKey("child-key")
				page, _ := entities.NewPage(key, value_objects.NewNullableString("/child-path").Value(), entities.NewSiteID(1), entities.PageTypeContent)
				page.SetID(entities.NewPageID(3))
				page.SetTimestamps(now, now)
				page.UpdateIndex(3)
				parentID := entities.NewPageID(1)
				page.SetParent(&parentID)
				return page
			}(),
			want: &models.Page{
				Base: models.Base{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:      "child-key",
				Path:     value_objects.NewNullableString("/child-path").Value(),
				Index:    3,
				SiteID:   1,
				Type:     models.PageTypeContent,
				LinkURL:  nil,
				ParentID: func() *uint64 { id := uint64(1); return &id }(),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModel(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPageMapper_ToDomain(t *testing.T) {
	mapper := NewPageMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *models.Page
		expectErr bool
		validate  func(t *testing.T, result *entities.Page)
	}{
		{
			name:      "nil input",
			input:     nil,
			expectErr: false,
			validate: func(t *testing.T, result *entities.Page) {
				assert.Nil(t, result)
			},
		},
		{
			name: "valid content page",
			input: &models.Page{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:      "content-key",
				Path:     value_objects.NewNullableString("/content-path").Value(),
				Index:    1,
				SiteID:   1,
				Type:     models.PageTypeContent,
				LinkURL:  nil,
				ParentID: nil,
			},
			expectErr: false,
			validate: func(t *testing.T, result *entities.Page) {
				assert.NotNil(t, result)
				assert.Equal(t, uint64(1), result.ID().Value())
				assert.Equal(t, "content-key", result.Key().Value())

				// Safe comparison for nullable string
				expectedPath := value_objects.NewNullableString("/content-path").Value()
				if expectedPath != nil && result.Path() != nil {
					assert.Equal(t, *expectedPath, *result.Path())
				} else {
					assert.Equal(t, expectedPath, result.Path())
				}

				assert.Equal(t, 1, result.Index())
				assert.Equal(t, uint64(1), result.SiteID().Value())
				assert.Equal(t, entities.PageTypeContent, result.Type())
				assert.Nil(t, result.LinkURL())
				assert.Nil(t, result.ParentID())
				assert.Equal(t, now, result.CreatedAt())
				assert.Equal(t, now, result.UpdatedAt())
			},
		},
		{
			name: "valid link page with URL",
			input: &models.Page{
				Base: models.Base{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:     "link-key",
				Path:    value_objects.NewNullableString("/link-path").Value(),
				Index:   2,
				SiteID:  1,
				Type:    models.PageTypeLink,
				LinkURL: value_objects.NewNullableString("https://example.com").Value(),
			},
			expectErr: false,
			validate: func(t *testing.T, result *entities.Page) {
				assert.NotNil(t, result)
				assert.Equal(t, uint64(2), result.ID().Value())
				assert.Equal(t, "link-key", result.Key().Value())

				// Safe comparison for nullable string path
				expectedPath := value_objects.NewNullableString("/link-path").Value()
				if expectedPath != nil && result.Path() != nil {
					assert.Equal(t, *expectedPath, *result.Path())
				} else {
					assert.Equal(t, expectedPath, result.Path())
				}

				assert.Equal(t, 2, result.Index())
				assert.Equal(t, uint64(1), result.SiteID().Value())
				assert.Equal(t, entities.PageTypeLink, result.Type())

				// Safe comparison for nullable string linkURL
				expectedURL := value_objects.NewNullableString("https://example.com").Value()
				if expectedURL != nil && result.LinkURL() != nil {
					assert.Equal(t, *expectedURL, *result.LinkURL())
				} else {
					assert.Equal(t, expectedURL, result.LinkURL())
				}

				assert.Nil(t, result.ParentID())
			},
		},
		{
			name: "page with parent",
			input: &models.Page{
				Base: models.Base{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Key:      "child-key",
				Path:     value_objects.NewNullableString("/child-path").Value(),
				Index:    3,
				SiteID:   1,
				Type:     models.PageTypeContent,
				ParentID: func() *uint64 { id := uint64(1); return &id }(),
			},
			expectErr: false,
			validate: func(t *testing.T, result *entities.Page) {
				assert.NotNil(t, result)
				assert.Equal(t, uint64(3), result.ID().Value())
				assert.Equal(t, "child-key", result.Key().Value())
				assert.Equal(t, 3, result.Index())
				assert.NotNil(t, result.ParentID())
				assert.Equal(t, uint64(1), result.ParentID().Value())
			},
		},
		{
			name: "invalid key",
			input: &models.Page{
				Base: models.Base{
					ID: 1,
				},
				Key:    "invalid key**",
				SiteID: 1,
				Type:   models.PageTypeContent,
			},
			expectErr: true,
			validate: func(t *testing.T, result *entities.Page) {
				assert.Nil(t, result)
			},
		},
		{
			name: "empty key should fail",
			input: &models.Page{
				Base: models.Base{
					ID: 1,
				},
				Key:    "",
				SiteID: 1,
				Type:   models.PageTypeContent,
			},
			expectErr: true,
			validate: func(t *testing.T, result *entities.Page) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomain(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			tt.validate(t, got)
		})
	}
}

func TestPageMapper_ToModels(t *testing.T) {
	mapper := NewPageMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*entities.Page
		expectErr bool
		validate  func(t *testing.T, result []*models.Page)
	}{
		{
			name:      "nil input",
			input:     nil,
			expectErr: false,
			validate: func(t *testing.T, result []*models.Page) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "empty slice",
			input:     []*entities.Page{},
			expectErr: false,
			validate: func(t *testing.T, result []*models.Page) {
				assert.NotNil(t, result)
				assert.Len(t, result, 0)
			},
		},
		{
			name: "valid multiple pages",
			input: []*entities.Page{
				func() *entities.Page {
					key, _ := value_objects.NewPageKey("key1")
					page, _ := entities.NewPage(key, value_objects.NewNullableString("/path1").Value(), entities.NewSiteID(1), entities.PageTypeContent)
					page.SetID(entities.NewPageID(1))
					page.SetTimestamps(now, now)
					page.UpdateIndex(1)
					return page
				}(),
				func() *entities.Page {
					key, _ := value_objects.NewPageKey("key2")
					page, _ := entities.NewPage(key, value_objects.NewNullableString("/path2").Value(), entities.NewSiteID(1), entities.PageTypeLink)
					page.SetID(entities.NewPageID(2))
					page.SetTimestamps(now, now)
					page.UpdateIndex(2)
					_ = page.SetLinkURL(value_objects.NewNullableString("https://example.com").Value())
					return page
				}(),
			},
			expectErr: false,
			validate: func(t *testing.T, result []*models.Page) {
				assert.NotNil(t, result)
				assert.Len(t, result, 2)

				// Validate first page
				assert.Equal(t, uint64(1), result[0].ID)
				assert.Equal(t, "key1", result[0].Key)
				assert.Equal(t, models.PageTypeContent, result[0].Type)
				assert.Nil(t, result[0].LinkURL)

				// Validate second page
				assert.Equal(t, uint64(2), result[1].ID)
				assert.Equal(t, "key2", result[1].Key)
				assert.Equal(t, models.PageTypeLink, result[1].Type)

				expectedURL := value_objects.NewNullableString("https://example.com").Value()
				assert.Equal(t, expectedURL, result[1].LinkURL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModels(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			tt.validate(t, got)
		})
	}
}

func TestPageMapper_ToDomains(t *testing.T) {
	mapper := NewPageMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*models.Page
		expectErr bool
		validate  func(t *testing.T, result []*entities.Page)
	}{
		{
			name:      "nil input",
			input:     nil,
			expectErr: false,
			validate: func(t *testing.T, result []*entities.Page) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "empty slice",
			input:     []*models.Page{},
			expectErr: false,
			validate: func(t *testing.T, result []*entities.Page) {
				assert.NotNil(t, result)
				assert.Len(t, result, 0)
			},
		},
		{
			name: "valid multiple models",
			input: []*models.Page{
				{
					Base: models.Base{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Key:    "key1",
					Path:   value_objects.NewNullableString("/path1").Value(),
					Index:  1,
					SiteID: 1,
					Type:   models.PageTypeContent,
				},
				{
					Base: models.Base{
						ID:        2,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Key:     "key2",
					Path:    value_objects.NewNullableString("/path2").Value(),
					Index:   2,
					SiteID:  1,
					Type:    models.PageTypeLink,
					LinkURL: value_objects.NewNullableString("https://example.com").Value(),
				},
			},
			expectErr: false,
			validate: func(t *testing.T, result []*entities.Page) {
				assert.NotNil(t, result)
				assert.Len(t, result, 2)

				// Validate first page
				assert.Equal(t, uint64(1), result[0].ID().Value())
				assert.Equal(t, "key1", result[0].Key().Value())
				assert.Equal(t, entities.PageTypeContent, result[0].Type())
				assert.Nil(t, result[0].LinkURL())

				// Validate second page
				assert.Equal(t, uint64(2), result[1].ID().Value())
				assert.Equal(t, "key2", result[1].Key().Value())
				assert.Equal(t, entities.PageTypeLink, result[1].Type())

				expectedURL := value_objects.NewNullableString("https://example.com").Value()
				if expectedURL != nil && result[1].LinkURL() != nil {
					assert.Equal(t, *expectedURL, *result[1].LinkURL())
				} else {
					assert.Equal(t, expectedURL, result[1].LinkURL())
				}
			},
		},
		{
			name: "error with invalid key",
			input: []*models.Page{
				{
					Base: models.Base{
						ID: 1,
					},
					Key:    "invalid key**",
					SiteID: 1,
					Type:   models.PageTypeContent,
				},
			},
			expectErr: true,
			validate: func(t *testing.T, result []*entities.Page) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomains(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			tt.validate(t, got)
		})
	}
}
