package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func initializeQueries(ctx context.Context) *db.Queries {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("Error reading .env: %v", err)
	}
	dbURL := os.Getenv("POSTGRES_CONNECTION_STRING")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return db.New(pool)
}
