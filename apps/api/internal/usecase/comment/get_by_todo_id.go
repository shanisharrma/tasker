package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/domain/todo"
)

type GetByTodoID struct {
	todoRepo    todo.Repository
	commentRepo comment.Repository
	logger      *zerolog.Logger
}

func NewGetByTodoID(todoRepo todo.Repository, commentRepo comment.Repository, logger *zerolog.Logger) *GetByTodoID {
	return &GetByTodoID{
		todoRepo:    todoRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (uc *GetByTodoID) Execute(ctx context.Context, userID string, todoID uuid.UUID) ([]comment.Comment, error) {

	// Validate todo exists and belongs to user
	_, err := uc.todoRepo.CheckExists(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	comments, err := uc.commentRepo.GetCommentsByTodoID(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to fetch comments by todo ID")
		return nil, err
	}

	return comments, nil
}
