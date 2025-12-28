package comment

import (
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/domain/todo"
	commentUC "github.com/shanisharrma/tasker/internal/usecase/comment"
)

type Module struct {
	CreateCommentUC *commentUC.CreateComment
	GetByTodoID     *commentUC.GetByTodoID
	UpdateComment   *commentUC.UpdateComment
	DeleteComment   *commentUC.DeleteComment
}

func NewModule(todoRepo todo.Repository, commentRepo comment.Repository, logger *zerolog.Logger) *Module {
	return &Module{
		CreateCommentUC: commentUC.NewCreateComment(todoRepo, commentRepo, logger),
		GetByTodoID:     commentUC.NewGetByTodoID(todoRepo, commentRepo, logger),
		UpdateComment:   commentUC.NewUpdateComment(todoRepo, commentRepo, logger),
		DeleteComment:   commentUC.NewDeleteComment(todoRepo, commentRepo, logger),
	}
}
