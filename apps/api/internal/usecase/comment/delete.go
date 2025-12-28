package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/comment"
	"github.com/shanisharrma/tasker/internal/domain/todo"
)

type DeleteComment struct {
	todoRepo    todo.Repository
	commentRepo comment.Repository
	logger      *zerolog.Logger
}

func NewDeleteComment(todoRepo todo.Repository, commentRepo comment.Repository, logger *zerolog.Logger) *DeleteComment {
	return &DeleteComment{
		todoRepo:    todoRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (uc *DeleteComment) Execute(ctx context.Context, userID string, commentID uuid.UUID) error {

	// Validate comment exists and belongs to user
	_, err := uc.commentRepo.GetCommentByID(ctx, userID, commentID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("comment validation failed")
		return err
	}

	err = uc.commentRepo.DeleteComment(ctx, userID, commentID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	// Business event log
	uc.logger.Info().
		Str("event", "comment_deleted").
		Str("comment_id", commentID.String()).
		Msg("Comment deleted successfully")

	return nil
}
