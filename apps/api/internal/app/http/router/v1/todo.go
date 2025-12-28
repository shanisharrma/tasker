package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shanisharrma/tasker/internal/app/http/handler"
	"github.com/shanisharrma/tasker/internal/app/http/middleware"
)

func registerTodoRoutes(r *echo.Group, h *handler.TodoHandler, ch *handler.CommentHandler, ah *handler.AttachmentHandler, auth *middleware.AuthMiddleware) {
	// Todo operations
	todos := r.Group("/todos")
	todos.Use(auth.RequireAuth)

	// Collection operations
	todos.POST("", h.CreateTodo)
	todos.GET("", h.GetTodos)
	todos.GET("/stats", h.GetTodoStats)

	// Individual todo operations
	dynamicTodo := todos.Group("/:id")
	dynamicTodo.GET("", h.GetTodoByID)
	dynamicTodo.PATCH("", h.UpdateTodo)
	dynamicTodo.DELETE("", h.DeleteTodo)

	// Todo comments
	todoComments := dynamicTodo.Group("/comments")
	todoComments.POST("", ch.AddComment)
	todoComments.GET("", ch.GetCommentsByTodoID)

	// Todo attachments
	todoAttachments := dynamicTodo.Group("/attachments")
	todoAttachments.POST("", ah.UploadTodoAttachment)
	todoAttachments.DELETE("/:attachmentId", ah.DeleteTodoAttachment)
	todoAttachments.GET("/:attachmentId/download", ah.GetAttachmentPresignedURL)
}
