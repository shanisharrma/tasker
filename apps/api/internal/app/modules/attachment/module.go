package attachment

import (
	"github.com/rs/zerolog"
	"github.com/shanisharrma/tasker/internal/domain/attachment"
	"github.com/shanisharrma/tasker/internal/domain/storage"
	"github.com/shanisharrma/tasker/internal/domain/todo"
	"github.com/shanisharrma/tasker/internal/shared/config"
	attachmentUC "github.com/shanisharrma/tasker/internal/usecase/attachment"
)

type Module struct {
	UploadAttachmentUC *attachmentUC.UploadAttachment
	DeleteAttachmentUC *attachmentUC.DeleteAttachment
	GetPresignedURLUC  *attachmentUC.GetPresignedURL
}

func NewModule(todoRepo todo.Repository, attachmentRepo attachment.Repository, s3Storage storage.ObjectStorage, logger *zerolog.Logger, cfg *config.Config) *Module {

	return &Module{
		UploadAttachmentUC: attachmentUC.NewUploadAttachment(todoRepo, attachmentRepo, s3Storage, logger, cfg),
		DeleteAttachmentUC: attachmentUC.NewDeleteAttachment(todoRepo, attachmentRepo, s3Storage, logger, cfg),
		GetPresignedURLUC:  attachmentUC.NewGetPresignedURL(todoRepo, attachmentRepo, s3Storage, logger, cfg),
	}

}
