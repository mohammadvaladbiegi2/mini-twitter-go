package app

import (
	_ "twitter_clone/docs"

	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	e := echo.New()

	return e
}
