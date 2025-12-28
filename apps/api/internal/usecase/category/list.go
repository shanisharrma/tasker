package category

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type ListCategories struct {
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewListCategories(categoryRepo category.Repository, logger *zerolog.Logger) *ListCategories {
	return &ListCategories{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *ListCategories) Execute(ctx context.Context, userID string,
	query *category.GetCategoriesQuery,
) (*types.PaginatedResponse[category.Category], error) {

	categories, err := uc.categoryRepo.GetCategories(ctx, userID, query)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch categories")
		return nil, err
	}

	return categories, nil
}
