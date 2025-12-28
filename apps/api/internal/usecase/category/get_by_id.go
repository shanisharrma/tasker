package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
)

type GetByID struct {
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewGetByID(categoryRepo category.Repository, logger *zerolog.Logger) *GetByID {
	return &GetByID{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *GetByID) Execute(ctx context.Context, userID string, categoryID uuid.UUID) (*category.Category, error) {

	categoryItem, err := uc.categoryRepo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch category by ID")
		return nil, err
	}

	return categoryItem, nil
}
