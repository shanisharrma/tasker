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

type GetPresignedURL struct {
	todoRepo       todoDomain.Repository
	attachmentRepo attachmentDomain.Repository
	storage        storage.ObjectStorage
	logger         *zerolog.Logger
	cfg            *config.Config
}

func NewGetPresignedURL(todoRepo todoDomain.Repository, attachmentRepo attachmentDomain.Repository, storage storage.ObjectStorage, logger *zerolog.Logger, cfg *config.Config) *GetPresignedURL {
	return &GetPresignedURL{
		todoRepo:       todoRepo,
		attachmentRepo: attachmentRepo,
		storage:        storage,
		logger:         logger,
		cfg:            cfg,
	}
}

func (uc *GetPresignedURL) Execute(
	ctx context.Context,
	userID string,
	todoID uuid.UUID,
	attachmentID uuid.UUID,
) (string, error) {

	// Verify todo exists and belongs to user
	_, err := uc.todoRepo.CheckExists(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("todo validation failed")
		return "", err
	}

	// Get attachment details
	attachment, err := uc.attachmentRepo.GetAttachment(
		ctx,
		todoID,
		attachmentID,
	)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to get attachment details")
		return "", err
	}

	// Generate presigned URL
	url, err := uc.storage.Presign(ctx, uc.cfg.AWS.UploadBucket, attachment.DownloadKey)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to generate presigned URL")
		return "", err
	}

	return url, nil
}
