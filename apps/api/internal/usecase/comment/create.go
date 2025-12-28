package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/domain/todo"
)

type CreateComment struct {
	todoRepo    todo.Repository
	commentRepo comment.Repository
	logger      *zerolog.Logger
}

func NewCreateComment(todoRepo todo.Repository, commentRepo comment.Repository, logger *zerolog.Logger) *CreateComment {
	return &CreateComment{
		todoRepo:    todoRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (uc *CreateComment) Execute(ctx context.Context, userID string, todoID uuid.UUID,
	payload *comment.AddCommentPayload,
) (*comment.Comment, error) {

	// Validate todo exists and belongs to user
	_, err := uc.todoRepo.CheckExists(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	commentItem, err := uc.commentRepo.AddComment(ctx, userID, todoID, payload)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to add comment")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "comment_added").
		Str("comment_id", commentItem.ID.String()).
		Str("todo_id", todoID.String()).
		Msg("Comment added successfully")

	return commentItem, nil
}
