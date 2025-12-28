package attachment

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	attachmentDomain "github.com/shanisharrma/tasker/internal/domain/attachment"
	"github.com/shanisharrma/tasker/internal/domain/storage"
	todoDomain "github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/config"
	"github.com/shanisharrma/tasker/internal/shared/errs"
)

type UploadAttachment struct {
	todoRepo       todoDomain.Repository
	attachmentRepo attachmentDomain.Repository
	storage        storage.ObjectStorage
	logger         *zerolog.Logger
	cfg            *config.Config
}

func NewUploadAttachment(todoRepo todoDomain.Repository, attachmentRepo attachmentDomain.Repository, storage storage.ObjectStorage, logger *zerolog.Logger, cfg *config.Config) *UploadAttachment {
	return &UploadAttachment{
		todoRepo:       todoRepo,
		attachmentRepo: attachmentRepo,
		storage:        storage,
		logger:         logger,
		cfg:            cfg,
	}
}

func (uc *UploadAttachment) Execute(
	ctx context.Context,
	userID string,
	todoID uuid.UUID,
	file *multipart.FileHeader,
) (*attachmentDomain.TodoAttachment, error) {

	// Verify todo exists and belongs to user
	_, err := uc.todoRepo.CheckExists(ctx, userID, todoID)
	if err != nil {
		uc.logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to open uploaded file")
		return nil, errs.NewBadRequestError("failed to open uploaded file", false, nil, nil, nil)
	}
	defer src.Close()

	// Upload to S3
	s3Key, err := uc.storage.Upload(ctx, uc.cfg.AWS.UploadBucket, "todos/attachments/"+file.Filename, src)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to upload file to S3")
		return nil, errors.Wrap(err, "failed to upload file")
	}

	// Detect MIME type
	src, err = file.Open()
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to reopen file for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to read file for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)
	}
	mimeType := http.DetectContentType(buffer)

	// Create attachment record
	attachment, err := uc.attachmentRepo.UploadAttachment(
		ctx,
		todoID,
		userID,
		s3Key,
		file.Filename,
		file.Size,
		mimeType,
	)
	if err != nil {
		uc.logger.Error().Err(err).Msg("failed to create attachment record")
		return nil, err
	}

	uc.logger.Info().
		Str("attachment_id", attachment.ID.String()).
		Str("s3_key", s3Key).
		Msg("uploaded todo attachment")

	return attachment, nil
}
