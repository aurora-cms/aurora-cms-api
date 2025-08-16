// page_block_test.go
package mappers

import (
	"errors"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/stretchr/testify/assert"
)

func TestPageBlockMapper_ToModel(t *testing.T) {
	mapper := NewPageBlockMapper()

	tests := []struct {
		name     string
		input    *entities.PageBlock
		expected *models.PageBlock
		wantErr  error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  nil,
		},
		{
			name: "valid input",
			input: func() *entities.PageBlock {
				block, _ := entities.NewPageBlock(
					entities.NewPageVersionID(1),
					"block_key",
					1,
					"text",
					"content",
				)
				block.SetID(entities.NewPageBlockID(1))
				block.SetTimestamps(time.Now(), time.Now())
				return block
			}(),
			expected: &models.PageBlock{
				Base: models.Base{
					ID: 1,
				},
				PageVersionID: 1,
				BlockKey:      "block_key",
				Index:         1,
				ContentType:   "text",
				Content:       "content",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.ToModel(tt.input)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err)
			}
			if tt.expected != nil && result != nil {
				assert.Equal(t, tt.expected.PageVersionID, result.PageVersionID)
				assert.Equal(t, tt.expected.BlockKey, result.BlockKey)
				assert.Equal(t, tt.expected.Index, result.Index)
				assert.Equal(t, tt.expected.ContentType, result.ContentType)
				assert.Equal(t, tt.expected.Content, result.Content)
			}
		})
	}
}

func TestPageBlockMapper_ToDomain(t *testing.T) {
	mapper := NewPageBlockMapper()
	now := time.Now()
	tests := []struct {
		name     string
		input    *models.PageBlock
		expected *entities.PageBlock
		wantErr  error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  nil,
		},
		{
			name: "valid input",
			input: &models.PageBlock{
				Base: models.Base{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				PageVersionID: 1,
				BlockKey:      "block_key",
				Index:         1,
				ContentType:   "text",
				Content:       "content",
			},
			expected: func() *entities.PageBlock {
				block, _ := entities.NewPageBlock(
					entities.NewPageVersionID(1),
					"block_key",
					1,
					"text",
					"content",
				)
				block.SetTimestamps(now, now)
				block.SetID(entities.NewPageBlockID(1))
				return block
			}(),
			wantErr: nil,
		},
		{
			name: "invalid input",
			input: &models.PageBlock{
				PageVersionID: 0,
				BlockKey:      "",
				Index:         -1,
				ContentType:   "",
				Content:       "",
			},
			expected: nil,
			wantErr:  errors.New("invalid input"), // Assuming entities.NewPageBlock enforces validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.ToDomain(tt.input)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestPageBlockMapper_ToModels(t *testing.T) {
	mapper := NewPageBlockMapper()
	now := time.Now()
	tests := []struct {
		name     string
		input    []*entities.PageBlock
		expected []*models.PageBlock
		wantErr  error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  nil,
		},
		{
			name: "valid input",
			input: []*entities.PageBlock{
				func() *entities.PageBlock {
					block, _ := entities.NewPageBlock(
						entities.NewPageVersionID(1),
						"block_key",
						1,
						"text",
						"content",
					)
					block.SetID(entities.NewPageBlockID(1))
					block.SetTimestamps(now, now)
					return block
				}(),
			},
			expected: []*models.PageBlock{
				{
					Base: models.Base{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					PageVersionID: 1,
					BlockKey:      "block_key",
					Index:         1,
					ContentType:   "text",
					Content:       "content",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.ToModels(tt.input)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPageBlockMapper_ToDomains(t *testing.T) {
	mapper := NewPageBlockMapper()
	now := time.Now()
	tests := []struct {
		name     string
		input    []*models.PageBlock
		expected []*entities.PageBlock
		wantErr  error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  nil,
		},
		{
			name: "valid input",
			input: []*models.PageBlock{
				{
					Base: models.Base{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					PageVersionID: 1,
					BlockKey:      "block_key",
					Index:         1,
					ContentType:   "text",
					Content:       "content",
				},
			},
			expected: []*entities.PageBlock{
				func() *entities.PageBlock {
					block, _ := entities.NewPageBlock(
						entities.NewPageVersionID(1),
						"block_key",
						1,
						"text",
						"content",
					)
					block.SetTimestamps(now, now)
					block.SetID(entities.NewPageBlockID(1))
					return block
				}(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.ToDomains(tt.input)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
