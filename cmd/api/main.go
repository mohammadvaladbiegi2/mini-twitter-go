package main

import (
	"errors"
	"log/slog"
	"net/http"
	"twitter_clone/internal/modules/user"
	"twitter_clone/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	pool := repository.NewPoolReqToSQLDB()
	Userrepo := user.NewUserRepository(pool)
	Userservice := user.NewUserService(Userrepo)
	UserHandler := user.NewUserHandler(Userservice)

	// Routes

	e.POST("/signup", UserHandler.SignUp)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
