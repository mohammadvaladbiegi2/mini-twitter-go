package app

import (
	"twitter_clone/internal/modules/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	// Create Dependency
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Routs
	e.POST("/signup", userHandler.SignUp)
	e.POST("/login", userHandler.Login)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
