package main

import (
	"aauth/internal/db"
	"aauth/internal/handler"
	"aauth/internal/service"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("GOOSE_DBSTRING"))

	if err != nil {
		log.Fatalf("unable to connect to db")
	}

	queries := db.New(pool)

	authService := service.NewAuthService(queries)

	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()

	authHandler.RegisterRoutes(mux)

	log.Println("Starting http server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
