package mappers

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

func TestUserMapper_ToModel(t *testing.T) {
	mapper := NewUserMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *entities.User
		want      *models.User
		expectErr bool
	}{
		{
			name:      "nil input",
			input:     nil,
			want:      nil,
			expectErr: false,
		},
		{
			name: "valid user input",
			input: func() *entities.User {
				keycloakUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				keycloakID := value_objects.NewKeycloakIDFromUUID(keycloakUUID)
				role, _ := value_objects.NewUserRole("admin")
				user, _ := entities.NewUser(keycloakID, role)
				user.SetID(entities.NewUserID(123))
				user.SetTimestamps(now, now)
				return user
			}(),
			want: &models.User{
				Base: models.Base{
					ID:        123,
					CreatedAt: now,
					UpdatedAt: now,
				},
				KeycloakID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Role:       "admin",
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

func TestUserMapper_ToDomain(t *testing.T) {
	mapper := NewUserMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     *models.User
		want      *entities.User
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
			input: &models.User{
				Base: models.Base{
					ID:        123,
					CreatedAt: now,
					UpdatedAt: now,
				},
				KeycloakID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Role:       "admin",
			},
			want: func() *entities.User {
				keycloakUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				keycloakID := value_objects.NewKeycloakIDFromUUID(keycloakUUID)
				role, _ := value_objects.NewUserRole("admin")
				user, _ := entities.NewUser(keycloakID, role)
				user.SetID(entities.NewUserID(123))
				user.SetTimestamps(now, now)
				return user
			}(),
			expectErr: false,
		},
		{
			name: "invalid role",
			input: &models.User{
				Base: models.Base{
					ID: 123,
				},
				KeycloakID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Role:       "invalid_role",
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

func TestUserMapper_ToModels(t *testing.T) {
	mapper := NewUserMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*entities.User
		want      []*models.User
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
			input: []*entities.User{
				func() *entities.User {
					keycloakUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
					keycloakID := value_objects.NewKeycloakIDFromUUID(keycloakUUID)
					role, _ := value_objects.NewUserRole("admin")
					user, _ := entities.NewUser(keycloakID, role)
					user.SetID(entities.NewUserID(123))
					user.SetTimestamps(now, now)
					return user
				}(),
			},
			want: []*models.User{
				{
					Base: models.Base{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					KeycloakID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Role:       "admin",
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

func TestUserMapper_ToDomains(t *testing.T) {
	mapper := NewUserMapper()
	now := time.Now()

	tests := []struct {
		name      string
		input     []*models.User
		want      []*entities.User
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
			input: []*models.User{
				{
					Base: models.Base{
						ID:        123,
						CreatedAt: now,
						UpdatedAt: now,
					},
					KeycloakID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Role:       "admin",
				},
			},
			want: []*entities.User{
				func() *entities.User {
					keycloakUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
					keycloakID := value_objects.NewKeycloakIDFromUUID(keycloakUUID)
					role, _ := value_objects.NewUserRole("admin")
					user, _ := entities.NewUser(keycloakID, role)
					user.SetID(entities.NewUserID(123))
					user.SetTimestamps(now, now)
					return user
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
