package user

import (
	"net/http"
	"twitter_clone/internal/modules/user/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var req dtos.UserSignUpReq

	// Parse JSON request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	// Call service layer
	message, appErr := h.service.SignUp(req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, message)
}
