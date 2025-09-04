package timeline

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewTimeLineHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) GetTestUser(c echo.Context) error {
	userID := c.Get("userID").(int64)

	user, appErr := h.service.AddDataToRedis(userID)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusCreated, user)
}
