package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"reflect"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

func TestPageVersionMapper_ToModel(t *testing.T) {
	mapper := NewPageVersionMapper()
	now := time.Now()
	tests := []struct {
		name        string
		input       *entities.PageVersion
		expected    *models.PageVersion
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectError: false,
		},
		{
			name: "valid input",
			input: func() *entities.PageVersion {
				id := entities.NewPageVersionID(123)
				pageID := entities.NewPageID(456)
				v, err := entities.NewPageVersion(pageID, 1, "Title", value_objects.NewNullableString("Description").Value())
				if err != nil {
					t.Fatalf("failed to create PageVersion: %v", err)
				}
				if v == nil {
					t.Fatal("PageVersion is nil")
				}
				v.SetID(id)
				v.SetTimestamps(now, now)
				v.Publish()
				return v
			}(),
			expected: &models.PageVersion{
				Base: models.Base{
					ID:        123,
					CreatedAt: now,
					UpdatedAt: now,
				},
				PageID:      456,
				Version:     1,
				Title:       "Title",
				Description: value_objects.NewNullableString("Description").Value(),
				IsPublished: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := mapper.ToModel(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, actual)
			}
		})
	}
}

func TestPageVersionMapper_ToDomain(t *testing.T) {
	mapper := NewPageVersionMapper()

	tests := []struct {
		name        string
		input       *models.PageVersion
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expectError: false,
		},
		{
			name: "valid input",
			input: &models.PageVersion{
				Base: models.Base{
					ID:        123,
					CreatedAt: time.Unix(0, 0),
					UpdatedAt: time.Unix(0, 0),
				},
				PageID:      456,
				Version:     1,
				Title:       "Title",
				Description: value_objects.NewNullableString("Description").Value(),
				IsPublished: true,
			},
			expectError: false,
		},
		{
			name: "invalid input",
			input: &models.PageVersion{
				PageID: 456,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mapper.ToDomain(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPageVersionMapper_ToModels(t *testing.T) {
	mapper := NewPageVersionMapper()

	now := time.Now()
	tests := []struct {
		name        string
		input       []*entities.PageVersion
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expectError: false,
		},
		{
			name: "valid input",
			input: []*entities.PageVersion{
				func() *entities.PageVersion {
					id := entities.NewPageVersionID(123)
					pageID := entities.NewPageID(456)
					v, _ := entities.NewPageVersion(pageID, 1, "Title", value_objects.NewNullableString("Description").Value())
					v.SetID(id)
					v.SetTimestamps(now, now)
					return v
				}(),
			},
			expectError: false,
		},
		{
			name: "contains nil version",
			input: []*entities.PageVersion{
				nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mapper.ToModels(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPageVersionMapper_ToDomains(t *testing.T) {
	mapper := NewPageVersionMapper()

	tests := []struct {
		name        string
		input       []*models.PageVersion
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expectError: false,
		},
		{
			name: "valid input",
			input: []*models.PageVersion{
				{
					Base: models.Base{
						ID:        123,
						CreatedAt: time.Unix(0, 0),
						UpdatedAt: time.Unix(0, 0),
					},
					PageID:      456,
					Version:     1,
					Title:       "Title",
					Description: value_objects.NewNullableString("Description").Value(),
				},
			},
			expectError: false,
		},
		{
			name: "contains invalid model",
			input: []*models.PageVersion{
				nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mapper.ToDomains(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
