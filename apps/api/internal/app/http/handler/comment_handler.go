package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanisharrma/tasker/internal/app/http/middleware"
	commentModule "github.com/shanisharrma/tasker/internal/app/modules/comment"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/server"
)

type CommentHandler struct {
	Handler
	module *commentModule.Module
}

func NewCommentHandler(s *server.Server, module *commentModule.Module) *CommentHandler {
	return &CommentHandler{
		Handler: NewHandler(s),
		module:  module,
	}
}

func (h *CommentHandler) AddComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.AddCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.module.CreateCommentUC.Execute(c.Request().Context(), userID, payload.TodoID, payload)
		},
		http.StatusCreated,
		&comment.AddCommentPayload{},
	)(c)
}

func (h *CommentHandler) GetCommentsByTodoID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.GetCommentsByTodoIDPayload) ([]comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.module.GetByTodoID.Execute(c.Request().Context(), userID, payload.TodoID)
		},
		http.StatusOK,
		&comment.GetCommentsByTodoIDPayload{},
	)(c)
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.module.UpdateComment.Execute(c.Request().Context(), userID, payload.ID, payload.Content)
		},
		http.StatusOK,
		&comment.UpdateCommentPayload{},
	)(c)
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *comment.DeleteCommentPayload) error {
			userID := middleware.GetUserID(c)
			return h.module.DeleteComment.Execute(c.Request().Context(), userID, payload.ID)
		},
		http.StatusNoContent,
		&comment.DeleteCommentPayload{},
	)(c)
}
