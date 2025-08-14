package main

import (
	"fmt"
	"log"
	"os"
	"twitter_clone/internal/app"
	"twitter_clone/internal/repository"
)

func main() {
	// connect to database
	db := repository.NewPoolReqToSQLDB()

	// create echo Servers
	e := app.NewServer()

	// add routes
	app.RegisterRoutes(e, db)

	serverPort := os.Getenv("SERVER_PORT")

	// start server
	if err := e.Start(fmt.Sprintf(":%s", serverPort)); err != nil {
		log.Fatal(err)
	}
}
