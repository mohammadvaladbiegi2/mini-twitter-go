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

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dtos.UserSignUpReq true "User signup data"
// @Success      200 {object} dtos.UserSignUpRes
// @Failure      400 {object} apperror.AppError
// @Failure      500 {object} apperror.AppError
// @Router       /signup [post]
func (h *UserHandler) SignUp(c echo.Context) error {
	var req dtos.UserSignUpReq

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
func (h *UserHandler) Login(c echo.Context) error {
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
