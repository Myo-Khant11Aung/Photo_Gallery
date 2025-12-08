package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectDatabase() *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
        // This is NORMAL in production where Railway provides env vars.
        log.Println("No .env file found, using environment variables from the system")
    }

	connStr := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

