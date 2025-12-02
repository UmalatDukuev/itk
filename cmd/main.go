package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"itk/internal/handler"
	"itk/internal/repo"
	"itk/internal/service"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	port := getenv("APP_PORT", "8080")

	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "wallet")
	dbPassword := getenv("DB_PASSWORD", "wallet")
	dbName := getenv("DB_NAME", "wallet")
	dbSSLMode := getenv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	walletRepo := repo.NewPostgresWalletRepository(pool)
	walletService := service.NewWalletService(walletRepo)
	h := handler.NewHandler(walletService)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	addr := ":" + port
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
