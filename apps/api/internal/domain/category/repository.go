package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type Repository interface {
	CreateCategory(ctx context.Context, userID string, payload *CreateCategoryPayload) (*Category, error)
	GetCategoryByID(ctx context.Context, userID string, categoryID uuid.UUID) (*Category, error)
	GetCategories(ctx context.Context, userID string, query *GetCategoriesQuery) (*types.PaginatedResponse[Category], error)
	UpdateCategory(ctx context.Context, userID string, categoryID uuid.UUID, payload *UpdateCategoryPayload) (*Category, error)
	DeleteCategory(ctx context.Context, userID string, categoryID uuid.UUID) error
}
