package comment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddComment(ctx context.Context, userID string, todoID uuid.UUID, payload *AddCommentPayload) (*Comment, error)
	GetCommentsByTodoID(ctx context.Context, userID string, todoID uuid.UUID) ([]Comment, error)
	GetCommentByID(ctx context.Context, userID string, commentID uuid.UUID) (*Comment, error)
	UpdateComment(ctx context.Context, userID string, commentID uuid.UUID, content string) (*Comment, error)
	DeleteComment(ctx context.Context, userID string, commentID uuid.UUID) error
}
