package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewPoolReqToSQLDB() *pgxpool.Pool {

	//  load envirment varible
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ error to read .env file: %v", err)
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// create Data Base URL
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbName,
	)

	// Create new Pool requst to postgres db
	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("❌ error to create Pool: %v", err)
	}

	return pool
}

//  example use
// pool := repository.NewPoolReqToSQLDB()
// 	defer pool.Close()

// 	var greeting string
// 	err := pool.QueryRow(context.Background(), "SELECT email from users").Scan(&greeting)
// 	if err != nil {
// 		log.Fatal("erro in  query", err)
// 	}

// 	fmt.Println(greeting)
