package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/domain/todo"
)

type UpdateComment struct {
	todoRepo    todo.Repository
	commentRepo comment.Repository
	logger      *zerolog.Logger
}

func NewUpdateComment(todoRepo todo.Repository, commentRepo comment.Repository, logger *zerolog.Logger) *UpdateComment {
	return &UpdateComment{
		todoRepo:    todoRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (uc *UpdateComment) Execute(ctx context.Context, userID string, commentID uuid.UUID, content string) (*comment.Comment, error) {
	// Validate comment exists and belongs to user
	_, err := uc.commentRepo.GetCommentByID(ctx, userID, commentID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("comment validation failed")
		return nil, err
	}

	commentItem, err := uc.commentRepo.UpdateComment(ctx, userID, commentID, content)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "comment_updated").
		Str("comment_id", commentItem.ID.String()).
		Msg("Comment updated successfully")

	return commentItem, nil
}
