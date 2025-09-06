package timeline

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewTimeLineHandler(s Service) *Handler {
	return &Handler{s}
}

// GET /timeline/my_time_line?limit=&cursor=
// @Summary  Get user timeline (keyset pagination)
// @Tags     Timeline
// @Produce  json
// @Param    limit  query int    false "max 100 (default 20)"
// @Param    cursor query string false "opaque cursor for pagination"
// @Success  200    {object} timelinedto.TimelineResponse
// @Security Bearer
// @Router   /timeline [get]
func (h *Handler) MyTimeLine(c echo.Context) error {
	userID := c.Get("userID").(int64)

	limit := 20
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	cursor := c.QueryParam("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	resp, appErr := h.service.MyTimeline(userID, limit, cursorPtr)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, resp)
}
