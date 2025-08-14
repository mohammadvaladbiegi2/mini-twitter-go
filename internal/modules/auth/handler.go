package auth

import (
	"net/http"
	"twitter_clone/internal/modules/auth/dtos"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewAuthHandler(service Service) *Handler {
	return &Handler{service: service}
}

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dtos.SignUpReq true "User signup data"
// @Success      200 {object} dtos.SignUpRes
// @Failure      400 {object} apperror.AppError
// @Failure      500 {object} apperror.AppError
// @Router       /signup [post]
func (h *Handler) SignUp(c echo.Context) error {
	var req dtos.SignUpReq

	// Parse JSON request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	// Call service layer
	token, appErr := h.service.SignUp(req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, token)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dtos.LoginReq true "Login credentials"
// @Success      200 {object} dtos.LoginRes
// @Failure      400 {object} apperror.AppError
// @Failure      401 {object} apperror.AppError
// @Router       /login [post]
func (h *Handler) Login(c echo.Context) error {
	var req dtos.LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}
	token, appErr := h.service.Login(req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}
	return c.JSON(http.StatusOK, token)
}
