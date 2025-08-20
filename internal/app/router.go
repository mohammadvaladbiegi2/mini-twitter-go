package app

import (
	internalMiddleware "twitter_clone/internal/middleware"
	"twitter_clone/internal/modules/auth"
	"twitter_clone/internal/modules/tweet"
	tweetaction "twitter_clone/internal/modules/tweet/action"
	"twitter_clone/internal/modules/user"
	useraction "twitter_clone/internal/modules/user/action"
	userprofile "twitter_clone/internal/modules/user/profile"
	usersearch "twitter_clone/internal/modules/user/search"

	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create auth Dependency
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Create user Dependency

	userProfileRepo := userprofile.NewProfileRepository(db)
	userProfileService := userprofile.NewProfileService(userProfileRepo)

	userActionRepo := useraction.NewUserActionRepository(db)
	userActionService := useraction.NewUserActionService(userActionRepo)

	userSearchRepo := usersearch.NewSearchRepository(db)
	userSearchService := usersearch.NewSearchService(userSearchRepo)

	userHandler := user.NewUserHandler(userProfileService, userSearchService, userActionService)

	// Create tweet Dependency
	tweetRepo := tweet.NewRepository(db)
	tweetService := tweet.NewTweetService(tweetRepo)
	tweetActionRepo := tweetaction.NewUserActionRepository(db)
	tweetActionService := tweetaction.NewTweetActionService(tweetActionRepo)

	tweetHandler := tweet.NewTweetHandler(tweetService, tweetActionService)

	// Routs
	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	// route need protection and token
	authGroup := e.Group("")
	authGroup.Use(internalMiddleware.JWTauthentication)

	// user
	authGroup.PUT("users/update-profile", userHandler.UpdateProfile)
	authGroup.GET("users/get-profile", userHandler.GetProfile)
	authGroup.GET("users/search-by-user-name/username", userHandler.SearchUsersByUserName)
	authGroup.GET("users/get-by-user-name/username", userHandler.GetUserByUsername)
	authGroup.POST("users/follow/:target_id", userHandler.Follow)
	authGroup.POST("users/unfollow/:target_id", userHandler.Unfollow)

	// tweet
	authGroup.POST("tweets/create-new-tweet", tweetHandler.CreateNewTweet)
	authGroup.POST("tweets/:tweet_id/like", tweetHandler.Like)
	authGroup.POST("tweets/:tweet_id/dislike", tweetHandler.Dislike)
	authGroup.DELETE("tweets/:tweet_id/reaction", tweetHandler.RemoveReaction)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
