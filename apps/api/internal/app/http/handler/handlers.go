package handler

import (
	"github.com/shanisharrma/tasker/internal/app/modules/attachment"
	"github.com/shanisharrma/tasker/internal/app/modules/auth"
	"github.com/shanisharrma/tasker/internal/app/modules/category"
	"github.com/shanisharrma/tasker/internal/app/modules/comment"
	"github.com/shanisharrma/tasker/internal/app/modules/todo"
	"github.com/shanisharrma/tasker/internal/server"
)

type Handlers struct {
	Health     *HealthHandler
	OpenAPI    *OpenAPIHandler
	Todo       *TodoHandler
	Attachment *AttachmentHandler
	Comment    *CommentHandler
	Category   *CategoryHandler
	Auth       *AuthHandler
}

func NewHandlers(
	s *server.Server,
	todoModule *todo.Module,
	attachmentModule *attachment.Module,
	categoryModule *category.Module,
	commentModule *comment.Module,
	authModule *auth.Module,
) *Handlers {
	return &Handlers{
		Health:     NewHealthHandler(s),
		OpenAPI:    NewOpenAPIHandler(s),
		Todo:       NewTodoHandler(s, todoModule),
		Attachment: NewAttachmentHandler(s, attachmentModule),
		Comment:    NewCommentHandler(s, commentModule),
		Category:   NewCategoryHandler(s, categoryModule),
		Auth:       NewAuthHandler(authModule),
	}
}
