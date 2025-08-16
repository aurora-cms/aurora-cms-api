// <llm-snippet-file>site_test.go</llm-snippet-file>
package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"reflect"
	"testing"
	"time"
)

func TestSiteMapper_ToModel(t *testing.T) {
	mapper := NewSiteMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *entities.Site
		want      *models.Site
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid input",
			input: func() *entities.Site {
				domain, _ := value_objects.NewDomainName("example.com")
				site, _ := entities.NewSite(
					"Test",
					value_objects.NewNullableString("Description").Value(),
					domain,
					entities.NewTemplateID(1),
					entities.NewTenantID(1),
				)
				site.SetTimestamps(now, now)
				_ = site.SetID(entities.NewSiteID(1))
				return site
			}(),
			want: &models.Site{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:          "Test",
				Description:   value_objects.NewNullableString("Description").Value(),
				Domain:        "example.com",
				TitleTemplate: nil,
				Enabled:       true,
				TemplateID:    1,
				TenantID:      1,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModel(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToModel() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSiteMapper_ToDomain(t *testing.T) {
	mapper := NewSiteMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *models.Site
		want      *entities.Site
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid input",
			input: &models.Site{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:          "Test",
				Description:   value_objects.NewNullableString("Description").Value(),
				Domain:        "example.com",
				TitleTemplate: nil,
				Enabled:       true,
				TemplateID:    1,
				TenantID:      1,
			},
			want: func() *entities.Site {
				domain, _ := value_objects.NewDomainName("example.com")
				site, _ := entities.NewSite(
					"Test",
					value_objects.NewNullableString("Description").Value(),
					domain,
					entities.NewTemplateID(1),
					entities.NewTenantID(1),
				)
				site.SetTimestamps(now, now)
				_ = site.SetID(entities.NewSiteID(1))
				return site
			}(),
			expectErr: false,
		},
		{
			name: "invalid domain",
			input: &models.Site{
				Domain: "invalid domain",
			},
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomain(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToDomain() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSiteMapper_ToModels(t *testing.T) {
	mapper := NewSiteMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*entities.Site
		want      []*models.Site
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid input",
			input: []*entities.Site{
				func() *entities.Site {
					domain, _ := value_objects.NewDomainName("example.com")
					site, _ := entities.NewSite(
						"Test",
						value_objects.NewNullableString("Description").Value(),
						domain,
						entities.NewTemplateID(1),
						entities.NewTenantID(1),
					)
					site.SetTimestamps(now, now)
					_ = site.SetID(entities.NewSiteID(1))
					return site
				}(),
			},
			want: []*models.Site{
				{
					Base: models.Base{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:          "Test",
					Description:   value_objects.NewNullableString("Description").Value(),
					Domain:        "example.com",
					TitleTemplate: nil,
					Enabled:       true,
					TemplateID:    1,
					TenantID:      1,
				},
			},
			expectErr: false,
		},
		{
			name: "error in one element",
			input: []*entities.Site{
				func() *entities.Site { return nil }(),
			},
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModels(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToModels() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToModels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSiteMapper_ToDomains(t *testing.T) {
	mapper := NewSiteMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*models.Site
		want      []*entities.Site
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid input",
			input: []*models.Site{
				{
					Base: models.Base{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:          "Test",
					Description:   value_objects.NewNullableString("Description").Value(),
					Domain:        "example.com",
					TitleTemplate: nil,
					Enabled:       true,
					TemplateID:    1,
					TenantID:      1,
				},
			},
			want: []*entities.Site{
				func() *entities.Site {
					domain, _ := value_objects.NewDomainName("example.com")
					site, _ := entities.NewSite(
						"Test",
						value_objects.NewNullableString("Description").Value(),
						domain,
						entities.NewTemplateID(1),
						entities.NewTenantID(1),
					)
					site.SetTimestamps(now, now)
					_ = site.SetID(entities.NewSiteID(1))
					return site
				}(),
			},
			expectErr: false,
		},
		{
			name: "error in one element",
			input: []*models.Site{
				{
					Domain: "invalid domain",
				},
			},
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomains(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToDomains() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomains() = %v, want %v", got, tt.want)
			}
		})
	}
}
