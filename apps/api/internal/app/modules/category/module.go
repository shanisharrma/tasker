package category

import (
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/category"
	categoryUC "github.com/shanisharrma/tasker/internal/usecase/category"
)

type Module struct {
	CreateCategoryUC *categoryUC.CreateCategory
	ListCategoriesUC *categoryUC.ListCategories
	GetByIdUC        *categoryUC.GetByID
	UpdateCategory   *categoryUC.UpdateCategory
	DeleteCategory   *categoryUC.DeleteCategory
}

func NewModule(categoryRepo category.Repository, logger *zerolog.Logger) *Module {
	return &Module{
		CreateCategoryUC: categoryUC.NewCreateCategory(categoryRepo, logger),
		ListCategoriesUC: categoryUC.NewListCategories(categoryRepo, logger),
		GetByIdUC:        categoryUC.NewGetByID(categoryRepo, logger),
		UpdateCategory:   categoryUC.NewUpdateCategory(categoryRepo, logger),
		DeleteCategory:   categoryUC.NewDeleteCategory(categoryRepo, logger),
	}
}
