package attachment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	attachmentDomain "github.com/shanisharrma/tasker/internal/domain/attachment"
	"github.com/shanisharrma/tasker/internal/domain/storage"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/config"
)

type DeleteAttachment struct {
	todoRepo       todoDomain.Repository
	attachmentRepo attachmentDomain.Repository
	storage        storage.ObjectStorage
	logger         *zerolog.Logger
	cfg            *config.Config
}

func NewDeleteAttachment(todoRepo todoDomain.Repository, attachmentRepo attachmentDomain.Repository, storage storage.ObjectStorage, logger *zerolog.Logger, cfg *config.Config) *DeleteAttachment {
	return &DeleteAttachment{
		todoRepo:       todoRepo,
		attachmentRepo: attachmentRepo,
		storage:        storage,
		logger:         logger,
		cfg:            cfg,
	}
}

func (uc *DeleteAttachment) Execute(
	ctx context.Context,
	userID string,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) error {

	// Verify todo exists and belongs to user
	_, err := uc.todoRepo.CheckExists(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("todo validation failed")
		return err
	}

	// Get attachment details for S3 deletion
	attachment, err := uc.attachmentRepo.GetAttachment(
		ctx,
		todoID,
		attachmentID,
	)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to get attachment details")
		return err
	}

	// Delete attachment record
	err = uc.attachmentRepo.DeleteAttachment(
		ctx,
		todoID,
		attachmentID,
	)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to delete attachment record")
		return err
	}

	// Delete from S3 asynchronously
	go func() {
		err := uc.storage.Delete(ctx, uc.cfg.AWS.UploadBucket, attachment.DownloadKey)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Str("s3_key", attachment.DownloadKey).
				Msg("failed to delete attachment from S3")
		}
	}()

	uc.logger.Info().Msg("deleted todo attachment")

	return nil
}
