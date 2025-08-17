package user

import (
	"net/http"
	userdtos "twitter_clone/internal/modules/user/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewUserHandler(service Service) *Handler {
	return &Handler{service: service}
}

// User godoc
// @Summary      update user profile
// @Description  update user name and bio user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param data body userdtos.UpdateProfileReq true "Update profile credentials (bio is optional)"
// @Success      200 {object} userdtos.UpdateProfileRes
// @Failure      400 {object} apperror.AppError
// @Router       /users/update-profile [put]
func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var req userdtos.UpdateProfileReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("Invalid request body", nil, err))
	}

	userupdated, appErr := h.service.UpdateProfile(userID, req)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, userupdated)
}
