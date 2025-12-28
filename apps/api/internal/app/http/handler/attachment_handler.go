package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanisharrma/tasker/internal/app/http/middleware"
	attachmentModule "github.com/shanisharrma/tasker/internal/app/modules/attachment"
	"github.com/shanisharrma/tasker/internal/domain/attachment"
	"github.com/shanisharrma/tasker/internal/server"
	"github.com/shanisharrma/tasker/internal/shared/errs"
)

type AttachmentHandler struct {
	Handler
	module *attachmentModule.Module
}

func NewAttachmentHandler(s *server.Server, attachmentModule *attachmentModule.Module) *AttachmentHandler {
	return &AttachmentHandler{
		Handler: NewHandler(s),
		module:  attachmentModule,
	}
}

func (h *AttachmentHandler) UploadTodoAttachment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *attachment.UploadTodoAttachmentPayload) (*attachment.TodoAttachment, error) {
			userID := middleware.GetUserID(c)

			form, err := c.MultipartForm()
			if err != nil {
				return nil, errs.NewBadRequestError("multipart form not found", false, nil, nil, nil)
			}

			files := form.File["file"]
			if len(files) == 0 {
				return nil, errs.NewBadRequestError("no file found", false, nil, nil, nil)
			}

			if len(files) > 1 {
				return nil, errs.NewBadRequestError("only one file allowed per upload", false, nil, nil, nil)
			}

			return h.module.UploadAttachmentUC.Execute(c.Request().Context(), userID, payload.TodoID, files[0])
		},
		http.StatusCreated,
		&attachment.UploadTodoAttachmentPayload{},
	)(c)
}

func (h *AttachmentHandler) DeleteTodoAttachment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *attachment.DeleteTodoAttachmentPayload) error {
			userID := middleware.GetUserID(c)
			return h.module.DeleteAttachmentUC.Execute(c.Request().Context(), userID, payload.TodoID, payload.AttachmentID)
		},
		http.StatusNoContent,
		&attachment.DeleteTodoAttachmentPayload{},
	)(c)
}

func (h *AttachmentHandler) GetAttachmentPresignedURL(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *attachment.GetAttachmentPresignedURLPayload) (*struct {
			URL string `json:"url"`
		}, error,
		) {
			userID := middleware.GetUserID(c)
			url, err := h.module.GetPresignedURLUC.Execute(c.Request().Context(), userID, payload.TodoID, payload.AttachmentID)
			if err != nil {
				return nil, err
			}
			return &struct {
				URL string `json:"url"`
			}{URL: url}, nil
		},
		http.StatusOK,
		&attachment.GetAttachmentPresignedURLPayload{},
	)(c)
}
