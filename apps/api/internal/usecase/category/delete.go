package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
)

type DeleteCategory struct {
	categoryRepo category.Repository
	logger       *zerolog.Logger
}

func NewDeleteCategory(categoryRepo category.Repository, logger *zerolog.Logger) *DeleteCategory {
	return &DeleteCategory{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *DeleteCategory) Execute(ctx context.Context, userID string, categoryID uuid.UUID) error {

	err := uc.categoryRepo.DeleteCategory(ctx, userID, categoryID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to delete category")
		return err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "category_deleted").
		Str("category_id", categoryID.String()).
		Msg("Category deleted successfully")

	return nil
}
