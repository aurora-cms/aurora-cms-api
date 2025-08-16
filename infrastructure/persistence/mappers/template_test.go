package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"reflect"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

func TestTemplateMapper_ToModel(t *testing.T) {
	mapper := NewTemplateMapper()
	now := time.Now()
	tests := []struct {
		name     string
		input    *entities.Template
		expected *models.Template
		wantErr  bool
	}{
		{
			name:    "Nil input",
			input:   nil,
			wantErr: false,
		},
		{
			name: "Valid input",
			input: func() *entities.Template {
				t, _ := entities.NewTemplate("template1", "/path/to/file", value_objects.NewNullableString("description").Value())
				t.SetID(entities.NewTemplateID(123))
				t.SetTimestamps(now, now)
				t.Enable()
				return t
			}(),
			expected: func() *models.Template {
				return &models.Template{
					Base: models.Base{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:        "template1",
					Description: value_objects.NewNullableString("description").Value(),
					FilePath:    "/path/to/file",
					Enabled:     true,
				}
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ToModel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTemplateMapper_ToDomain(t *testing.T) {
	mapper := NewTemplateMapper()
	desc := value_objects.NewNullableString("description").Value()
	now := time.Now()
	tests := []struct {
		name     string
		input    *models.Template
		expected *entities.Template
		wantErr  bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "Valid input",
			input: &models.Template{
				Base: models.Base{
					ID:        123,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:        "template1",
				Description: desc,
				FilePath:    "/path/to/file",
				Enabled:     true,
			},
			expected: func() *entities.Template {
				t, _ := entities.NewTemplate("template1", "/path/to/file", desc)
				t.SetID(entities.NewTemplateID(123))
				t.Enable()
				return t
			}(),
			wantErr: false,
		},
		{
			name: "Disabled template",
			input: &models.Template{
				Base: models.Base{
					ID:        456,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:        "template2",
				Description: desc,
				FilePath:    "/path/to/file2",
				Enabled:     false,
			},
			expected: func() *entities.Template {
				t, _ := entities.NewTemplate("template2", "/path/to/file2", desc)
				t.SetID(entities.NewTemplateID(456))
				t.Disable()
				return t
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomain(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.input == nil {
				if got != nil {
					t.Errorf("ToDomain() = %v, want nil", got)
				}
				return
			}
			if got.ID().Value() != tt.expected.ID().Value() ||
				got.Name() != tt.expected.Name() ||
				got.Description() != tt.expected.Description() ||
				got.FilePath() != tt.expected.FilePath() ||
				got.IsEnabled() != tt.expected.IsEnabled() {
				t.Errorf("ToDomain() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTemplateMapper_ToModels(t *testing.T) {
	mapper := NewTemplateMapper()
	now := time.Now()
	desc := value_objects.NewNullableString("description").Value()
	desc2 := value_objects.NewNullableString("description2").Value()
	tests := []struct {
		name     string
		input    []*entities.Template
		expected []*models.Template
		wantErr  bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "Valid templates",
			input: []*entities.Template{
				func() *entities.Template {
					t, _ := entities.NewTemplate("template1", "/path/to/file", desc)
					t.SetID(entities.NewTemplateID(123))
					t.SetTimestamps(now, now)
					t.Enable()
					return t
				}(),
				func() *entities.Template {
					t, _ := entities.NewTemplate("template2", "/path/to/file2", desc2)
					t.SetID(entities.NewTemplateID(456))
					t.SetTimestamps(now, now)
					t.Disable()
					return t
				}(),
			},
			expected: []*models.Template{
				{
					Base: models.Base{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:        "template1",
					Description: desc,
					FilePath:    "/path/to/file",
					Enabled:     true,
				},
				{
					Base: models.Base{
						ID:        456,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:        "template2",
					Description: desc2,
					FilePath:    "/path/to/file2",
					Enabled:     false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModels(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ToModels() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewTemplateMapper_ToDomains(t *testing.T) {
	mapper := NewTemplateMapper()
	now := time.Now()
	desc := value_objects.NewNullableString("description").Value()
	desc2 := value_objects.NewNullableString("description2").Value()
	tests := []struct {
		name     string
		input    []*models.Template
		expected []*entities.Template
		wantErr  bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "Valid templates",
			input: []*models.Template{
				{
					Base: models.Base{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:        "template1",
					Description: desc,
					FilePath:    "/path/to/file",
					Enabled:     true,
				},
				{
					Base: models.Base{
						ID:        456,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:        "template2",
					Description: desc2,
					FilePath:    "/path/to/file2",
					Enabled:     false,
				},
			},
			expected: []*entities.Template{
				func() *entities.Template {
					t, _ := entities.NewTemplate("template1", "/path/to/file", desc)
					t.SetID(entities.NewTemplateID(123))
					t.SetTimestamps(now, now)
					t.Enable()
					return t
				}(),
				func() *entities.Template {
					t, _ := entities.NewTemplate("template2", "/path/to/file2", desc2)
					t.SetID(entities.NewTemplateID(456))
					t.SetTimestamps(now, now)
					t.Disable()
					return t
				}(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomains(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.input == nil {
				if got != nil {
					t.Errorf("ToDomains() = %v, want nil", got)
				}
				return
			}
			for i := range got {
				if got[i].ID().Value() != tt.expected[i].ID().Value() ||
					got[i].Name() != tt.expected[i].Name() ||
					got[i].Description() != tt.expected[i].Description() ||
					got[i].FilePath() != tt.expected[i].FilePath() ||
					got[i].IsEnabled() != tt.expected[i].IsEnabled() {
					t.Errorf("ToDomains()[%d] = %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
