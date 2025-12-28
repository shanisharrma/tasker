package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
)

type UpdateCategory struct {
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewUpdateCategory(categoryRepo category.Repository, logger *zerolog.Logger) *UpdateCategory {
	return &UpdateCategory{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *UpdateCategory) Execute(ctx context.Context, userID string, categoryID uuid.UUID,
	payload *category.UpdateCategoryPayload,
) (*category.Category, error) {

	categoryItem, err := uc.categoryRepo.UpdateCategory(ctx, userID, categoryID, payload)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to update category")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "category_updated").
		Str("category_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Msg("Category updated successfully")

	return categoryItem, nil
}
