package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanisharrma/tasker/internal/app/http/middleware"
	todoModule "github.com/shanisharrma/tasker/internal/app/modules/todo"
	"github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/server"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type TodoHandler struct {
	Handler
	module *todoModule.Module
}

func NewTodoHandler(s *server.Server, module *todoModule.Module) *TodoHandler {
	return &TodoHandler{
		Handler: NewHandler(s),
		module:  module,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.module.CreateTodoUC.Execute(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&todo.CreateTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodoByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoByIDPayload) (*todo.PopulatedTodo, error) {
			userID := middleware.GetUserID(c)
			return h.module.GetTodoByIdUC.Execute(c.Request().Context(), userID, payload.ID)
		},
		http.StatusOK,
		&todo.GetTodoByIDPayload{},
	)(c)
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *todo.GetTodosQuery) (*types.PaginatedResponse[todo.PopulatedTodo], error) {
			userID := middleware.GetUserID(c)
			return h.module.ListTodosUC.Execute(c.Request().Context(), userID, query)
		},
		http.StatusOK,
		&todo.GetTodosQuery{},
	)(c)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.module.UpdateTodoUC.Execute(c.Request().Context(), userID, payload)
		},
		http.StatusOK,
		&todo.UpdateTodoPayload{},
	)(c)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *todo.DeleteTodoPayload) error {
			userID := middleware.GetUserID(c)
			return h.module.DeleteTodoUC.Execute(c.Request().Context(), userID, payload.ID)
		},
		http.StatusNoContent,
		&todo.DeleteTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodoStats(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoStatsPayload) (*todo.TodoStats, error) {
			userID := middleware.GetUserID(c)
			return h.module.GetStatsUC.Execute(c.Request().Context(), userID)
		},
		http.StatusOK,
		&todo.GetTodoStatsPayload{},
	)(c)
}
