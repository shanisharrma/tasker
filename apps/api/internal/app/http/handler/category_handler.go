package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanisharrma/tasker/internal/app/http/middleware"
	categoryModule "github.com/shanisharrma/tasker/internal/app/modules/category"
	"github.com/shanisharrma/tasker/internal/domain/category"
	"github.com/shanisharrma/tasker/internal/server"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type CategoryHandler struct {
	Handler
	module *categoryModule.Module
}

func NewCategoryHandler(s *server.Server, module *categoryModule.Module) *CategoryHandler {
	return &CategoryHandler{
		Handler: NewHandler(s),
		module:  module,
	}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *category.CreateCategoryPayload) (*category.Category, error) {
			userID := middleware.GetUserID(c)
			return h.module.CreateCategoryUC.Execute(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&category.CreateCategoryPayload{},
	)(c)
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *category.GetCategoriesQuery) (
			*types.PaginatedResponse[category.Category], error,
		) {
			userID := middleware.GetUserID(c)
			return h.module.ListCategoriesUC.Execute(c.Request().Context(), userID, query)
		},
		http.StatusOK,
		&category.GetCategoriesQuery{},
	)(c)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *category.UpdateCategoryPayload) (*category.Category, error) {
			userID := middleware.GetUserID(c)
			return h.module.UpdateCategory.Execute(c.Request().Context(), userID, payload.ID, payload)
		},
		http.StatusOK,
		&category.UpdateCategoryPayload{},
	)(c)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *category.DeleteCategoryPayload) error {
			userID := middleware.GetUserID(c)
			return h.module.DeleteCategory.Execute(c.Request().Context(), userID, payload.ID)
		},
		http.StatusNoContent,
		&category.DeleteCategoryPayload{},
	)(c)
}
