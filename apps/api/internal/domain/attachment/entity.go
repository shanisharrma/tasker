package attachment

import (
	"github.com/google/uuid"
	"github.com/shanisharrma/tasker/internal/shared/types"
)

type TodoAttachment struct {
	types.Base

	TodoID      uuid.UUID `json:"todoId" db:"todo_id"`
	Name        string    `json:"name" db:"name"`
	UploadedBy  string    `json:"uploadedBy" db:"uploaded_by"`
	DownloadKey string    `json:"downloadKey" db:"download_key"`
	FileSize    *int64    `json:"fileSize" db:"file_size"`
	Mimetype    *string   `json:"mimeType" db:"mime_type"`
}
