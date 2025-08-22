package app

import (
	internalMiddleware "twitter_clone/internal/middleware"
	"twitter_clone/internal/modules/auth"
	"twitter_clone/internal/modules/tweet"
	tweetaction "twitter_clone/internal/modules/tweet/action"
	tweetbookmark "twitter_clone/internal/modules/tweet/bookmark"
	tweetreply "twitter_clone/internal/modules/tweet/reply"
	"twitter_clone/internal/modules/user"
	useraction "twitter_clone/internal/modules/user/action"
	userconnection "twitter_clone/internal/modules/user/connection"
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

	userConnectionRepo := userconnection.NewUserConnectionRepository(db)
	userConnectionService := userconnection.NewUserConnectionService(userConnectionRepo)

	userHandler := user.NewUserHandler(userProfileService, userSearchService, userActionService, userConnectionService)

	// Create tweet Dependency
	tweetRepo := tweet.NewRepository(db)
	tweetService := tweet.NewTweetService(tweetRepo)
	tweetActionRepo := tweetaction.NewUserActionRepository(db)
	tweetActionService := tweetaction.NewTweetActionService(tweetActionRepo)
	tweetReplyRepo := tweetreply.NewReplyRepository(db)
	tweetReplyService := tweetreply.NewReplyService(tweetReplyRepo)
	tweetBookMarkRepo := tweetbookmark.NewBookMarkRepository(db)
	tweetBookMarkService := tweetbookmark.NewBookMarkService(tweetBookMarkRepo)

	tweetHandler := tweet.NewTweetHandler(tweetService, tweetActionService, tweetReplyService, tweetBookMarkService)

	// Routs
	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	// route need protection and token
	authGroup := e.Group("")
	authGroup.Use(internalMiddleware.JWTauthentication)

	// user
	authGroup.PUT("users/update-profile", userHandler.UpdateProfile)
	authGroup.GET("users/get-profile", userHandler.GetProfile)
	authGroup.GET("users/search-by-user-name", userHandler.SearchUsersByUserName)
	authGroup.GET("users/get-by-user-name", userHandler.GetUserByUsername)
	authGroup.POST("users/follow", userHandler.Follow)
	authGroup.POST("users/unfollow", userHandler.Unfollow)
	authGroup.GET("users/followers", userHandler.GetFollowers)
	authGroup.GET("users/followings", userHandler.GetFollowings)

	// tweet
	authGroup.POST("tweets/create-new-tweet", tweetHandler.CreateNewTweet)
	authGroup.POST("tweets/:tweet_id/like", tweetHandler.Like)
	authGroup.POST("tweets/:tweet_id/dislike", tweetHandler.Dislike)
	authGroup.DELETE("tweets/:tweet_id/reaction", tweetHandler.RemoveReaction)
	authGroup.POST("tweets/:tweet_id/reply", tweetHandler.CreateReply)
	authGroup.GET("tweets/:tweet_id/replies", tweetHandler.GetReplies)
	authGroup.POST("tweets/:tweet_id/bookmark", tweetHandler.Bookmark)
	authGroup.DELETE("tweets/:tweet_id/bookmark", tweetHandler.Unbookmark)
	authGroup.GET("tweets/bookmarks", tweetHandler.ListBookmarks)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
