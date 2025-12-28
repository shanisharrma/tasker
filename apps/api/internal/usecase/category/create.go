package category

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
)

type CreateCategory struct {
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewCreateCategory(categoryRepo category.Repository, logger *zerolog.Logger) *CreateCategory {
	return &CreateCategory{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *CreateCategory) Execute(ctx context.Context, userID string,
	payload *category.CreateCategoryPayload,
) (*category.Category, error) {

	categoryItem, err := uc.categoryRepo.CreateCategory(ctx, userID, payload)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to create category")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "category_created").
		Str("category_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Str("color", categoryItem.Color).
		Msg("Category created successfully")

	return categoryItem, nil
}
