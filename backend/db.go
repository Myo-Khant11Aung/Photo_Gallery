package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectDatabase() *pgxpool.Pool {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

