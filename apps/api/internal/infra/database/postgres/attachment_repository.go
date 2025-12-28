package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shanisharrma/tasker/internal/domain/attachment"
	"github.com/shanisharrma/tasker/internal/infra/database"
	"github.com/shanisharrma/tasker/internal/shared/errs"
)

type AttachmentRepository struct {
	db *database.Database
}

func NewAttachmentRepository(db *database.Database) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) GetAttachment(
	ctx context.Context,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) (*attachment.TodoAttachment, error) {
	stmt := `
		SELECT
			*
		FROM
			todo_attachments
		WHERE
			todo_id = @todo_id
			AND id = @attachment_id
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id":       todoID,
		"attachment_id": attachmentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get todo attachment: %w", err)
	}

	attachment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[attachment.TodoAttachment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			code := "ATTACHMENT_NOT_FOUND"
			return nil, errs.NewNotFoundError("attachment not found", false, &code)
		}
		return nil, fmt.Errorf("failed to collect row from table:todo_attachments: %w", err)
	}

	return &attachment, nil
}

func (r *AttachmentRepository) GetAttachments(
	ctx context.Context,
	todoID uuid.UUID,
) ([]attachment.TodoAttachment, error) {
	stmt := `
		SELECT
			*
		FROM
			todo_attachments
		WHERE
			todo_id = @todo_id
		ORDER BY
			created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id": todoID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get todo attachments: %w", err)
	}

	attachments, err := pgx.CollectRows(rows, pgx.RowToStructByName[attachment.TodoAttachment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []attachment.TodoAttachment{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:todo_attachments: %w", err)
	}

	return attachments, nil
}

func (r *AttachmentRepository) DeleteAttachment(
	ctx context.Context,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) error {
	stmt := `
		DELETE FROM todo_attachments
		WHERE
			todo_id = @todo_id
			AND id = @attachment_id
	`

	result, err := r.db.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"todo_id":       todoID,
		"attachment_id": attachmentID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete todo attachment: %w", err)
	}

	if result.RowsAffected() == 0 {
		code := "ATTACHMENT_NOT_FOUND"
		return errs.NewNotFoundError("attachment not found", false, &code)
	}

	return nil
}

func (r *AttachmentRepository) UploadAttachment(
	ctx context.Context,
	todoID uuid.UUID,
	userID string,
	s3Key string,
	fileName string,
	fileSize int64,
	mimeType string,
) (*attachment.TodoAttachment, error) {
	stmt := `
		INSERT INTO
			todo_attachments (
				todo_id,
				name,
				uploaded_by,
				download_key,
				file_size,
				mime_type
			)
		VALUES
			(
				@todo_id,
				@name,
				@uploaded_by,
				@download_key,
				@file_size,
				@mime_type
			)
		RETURNING
			*
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id":      todoID,
		"name":         fileName,
		"uploaded_by":  userID,
		"download_key": s3Key,
		"file_size":    fileSize,
		"mime_type":    mimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create todo attachment for todo_id=%s: %w", todoID.String(), err)
	}

	attachment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[attachment.TodoAttachment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todo_attachments: %w", err)
	}

	return &attachment, nil
}
