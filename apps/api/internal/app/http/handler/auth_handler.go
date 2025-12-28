package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authModule "github.com/shanisharrma/tasker/internal/app/modules/auth"
)

type AuthHandler struct {
	module *authModule.Module
}

func NewAuthHandler(module *authModule.Module) *AuthHandler {
	return &AuthHandler{module: module}
}

func (h *AuthHandler) GetUserEmail(c echo.Context) error {
	userID := c.Param("userId")

	email, err := h.module.GetUserEmail.Execute(
		c.Request().Context(),
		userID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"email": email,
	})
}
