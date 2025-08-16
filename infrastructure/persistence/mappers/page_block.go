package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// PageBlockMapper handles conversion between domain entities and GORM models
type PageBlockMapper struct{}

// NewPageBlockMapper creates a new PageBlockMapper
func NewPageBlockMapper() *PageBlockMapper {
	return &PageBlockMapper{}
}

// ToModel converts a domain PageBlock to a GORM models.PageBlock
func (m *PageBlockMapper) ToModel(block *entities.PageBlock) (*models.PageBlock, error) {
	if block == nil {
		return nil, nil
	}

	return &models.PageBlock{
		Base: models.Base{
			ID:        block.ID().Value(),
			CreatedAt: block.CreatedAt(),
			UpdatedAt: block.UpdatedAt(),
		},
		PageVersionID: block.PageVersionID().Value(),
		BlockKey:      block.BlockKey(),
		Index:         block.Index(),
		ContentType:   block.ContentType(),
		Content:       block.Content(),
	}, nil
}

// ToDomain converts a GORM models.PageBlock to a domain PageBlock
func (m *PageBlockMapper) ToDomain(model *models.PageBlock) (*entities.PageBlock, error) {
	if model == nil {
		return nil, nil
	}

	block, err := entities.NewPageBlock(
		entities.NewPageVersionID(model.PageVersionID),
		model.BlockKey,
		model.Index,
		model.ContentType,
		model.Content,
	)
	if err != nil {
		return nil, err
	}

	block.SetID(entities.NewPageBlockID(model.ID))
	block.SetTimestamps(model.CreatedAt, model.UpdatedAt)

	return block, nil
}

// ToModels converts a slice of domain PageBlocks to GORM models
func (m *PageBlockMapper) ToModels(blocks []*entities.PageBlock) ([]*models.PageBlock, error) {
	if blocks == nil {
		return nil, nil
	}

	result := make([]*models.PageBlock, len(blocks))
	for i, block := range blocks {
		model, err := m.ToModel(block)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

// ToDomains converts a slice of GORM models to domain PageBlocks
func (m *PageBlockMapper) ToDomains(modelList []*models.PageBlock) ([]*entities.PageBlock, error) {
	if modelList == nil {
		return nil, nil
	}

	result := make([]*entities.PageBlock, len(modelList))
	for i, model := range modelList {
		block, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = block
	}

	return result, nil
}
