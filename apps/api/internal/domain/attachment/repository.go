package attachment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetAttachment(ctx context.Context, todoID, attachmentID uuid.UUID) (*TodoAttachment, error)
	GetAttachments(ctx context.Context, todoID uuid.UUID) ([]TodoAttachment, error)
	UploadAttachment(ctx context.Context, todoID uuid.UUID, userID, key, name string, size int64, mime string) (*TodoAttachment, error)
	DeleteAttachment(ctx context.Context, todoID, attachmentID uuid.UUID) error
}
