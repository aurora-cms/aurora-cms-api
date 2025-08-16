package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"reflect"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

func TestTenantMapper_ToModel(t *testing.T) {
	mapper := NewTenantMapper()
	desc := value_objects.NewNullableString("TestDescription").Value()
	now := time.Now()
	tests := []struct {
		name     string
		input    *entities.Tenant
		expected *models.Tenant
		wantErr  bool
	}{
		{
			name:     "NilInput",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "ValidInput",
			input: func() *entities.Tenant {
				tenant, _ := entities.NewTenant("TestTenant", desc)
				tenant.SetID(entities.NewTenantID(1))
				tenant.SetTimestamps(now, now)
				tenant.EnableBilling()
				tenant.Activate()
				return tenant
			}(),
			expected: &models.Tenant{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:             "TestTenant",
				Description:      desc,
				IsActive:         true,
				IsBillingEnabled: true,
			},
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

func TestTenantMapper_ToDomain(t *testing.T) {
	mapper := NewTenantMapper()
	desc := value_objects.NewNullableString("TestDescription").Value()
	desc2 := value_objects.NewNullableString("").Value()
	now := time.Now()
	tests := []struct {
		name     string
		input    *models.Tenant
		expected *entities.Tenant
		wantErr  bool
	}{
		{
			name:     "NilInput",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "ValidInput",
			input: &models.Tenant{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Name:             "TestTenant",
				Description:      desc,
				IsActive:         true,
				IsBillingEnabled: true,
			},
			expected: func() *entities.Tenant {
				tenant, _ := entities.NewTenant("TestTenant", desc)
				tenant.SetID(entities.NewTenantID(1))
				tenant.SetTimestamps(now, now)
				tenant.EnableBilling()
				tenant.Activate()
				return tenant
			}(),
			wantErr: false,
		},
		{
			name: "InvalidEntity",
			input: &models.Tenant{
				Base: models.Base{
					ID: 1,
				},
				Name:             "",
				Description:      desc2,
				IsActive:         true,
				IsBillingEnabled: true,
			},
			expected: nil,
			wantErr:  true,
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
			if got == nil {
				if tt.expected != nil {
					t.Errorf("ToDomain() = nil, want %v", tt.expected)
				}
				return
			}
			if got.ID().Value() != tt.expected.ID().Value() ||
				got.Name() != tt.expected.Name() ||
				got.Description() != tt.expected.Description() ||
				got.IsActive() != tt.expected.IsActive() ||
				got.IsBillingEnabled() != tt.expected.IsBillingEnabled() {
				t.Errorf("ToDomain() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTenantMapper_ToModels(t *testing.T) {
	mapper := NewTenantMapper()
	now := time.Now()
	desc1 := value_objects.NewNullableString("Description1").Value()
	desc2 := value_objects.NewNullableString("Description2").Value()

	tests := []struct {
		name      string
		input     []*entities.Tenant
		want      []*models.Tenant
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
			input: []*entities.Tenant{
				func() *entities.Tenant {
					tenant, _ := entities.NewTenant("Tenant1", desc1)
					tenant.SetID(entities.NewTenantID(1))
					tenant.SetTimestamps(now, now)
					return tenant
				}(),
				func() *entities.Tenant {
					tenant, _ := entities.NewTenant("Tenant2", desc2)
					tenant.SetID(entities.NewTenantID(2))
					tenant.SetTimestamps(now, now)
					return tenant
				}(),
			},
			want: []*models.Tenant{
				{
					Base:             models.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
					Name:             "Tenant1",
					Description:      desc1,
					IsActive:         true,
					IsBillingEnabled: false,
				},
				{
					Base:             models.Base{ID: 2, CreatedAt: now, UpdatedAt: now},
					Name:             "Tenant2",
					Description:      desc2,
					IsActive:         true,
					IsBillingEnabled: false,
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToModels(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToModels() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToModels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantMapper_ToDomains(t *testing.T) {
	mapper := NewTenantMapper()
	desc1 := value_objects.NewNullableString("Description1").Value()
	desc2 := value_objects.NewNullableString("Description2").Value()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*models.Tenant
		want      []*entities.Tenant
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
			input: []*models.Tenant{
				{
					Base:             models.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
					Name:             "Tenant1",
					Description:      desc1,
					IsActive:         true,
					IsBillingEnabled: false,
				},
				{
					Base:             models.Base{ID: 2, CreatedAt: now, UpdatedAt: now},
					Name:             "Tenant2",
					Description:      desc2,
					IsActive:         true,
					IsBillingEnabled: false,
				},
			},
			want: []*entities.Tenant{
				func() *entities.Tenant {
					tenant, _ := entities.NewTenant("Tenant1", desc1)
					tenant.SetID(entities.NewTenantID(1))
					tenant.SetTimestamps(now, now)
					return tenant
				}(),
				func() *entities.Tenant {
					tenant, _ := entities.NewTenant("Tenant2", desc2)
					tenant.SetID(entities.NewTenantID(2))
					tenant.SetTimestamps(now, now)
					return tenant
				}(),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.ToDomains(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToDomains() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomains() = %v, want %v", got, tt.want)
			}
		})
	}
}
