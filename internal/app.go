package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mr-Rafael/finance-calculator/internal/api"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type App struct {
	Handler http.Handler
	DB      *pgxpool.Pool
}

func New() *App {
	ctx := context.Background()

	// environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Error reading .env: %v\n", err)
	}
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	dbURL := os.Getenv("POSTGRES_CONNECTION_STRING")

	// database
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(pool)

	// repos
	usersRepo := repository.NewUsersRepo(queries)
	authRepo := repository.NewAuthRepo(queries)
	savingsRepo := repository.NewSavingsRepo(queries)
	loansRepo := repository.NewLoansRepo(queries)

	// services
	userService := service.NewUserService(usersRepo)
	authService := service.NewAuthService(authRepo, usersRepo, accessSecret, refreshSecret)
	savingsService := service.NewSavingsService(savingsRepo)
	loanService := service.NewLoansService(loansRepo)

	// handlers
	adminHandler := api.NewAdminHandler()
	userHandler := api.NewUsersHandler(userService)
	authHandler := api.NewAuthHandler(authService)
	savingsHandler := api.NewSavingsHandler(savingsService)
	loansHandler := api.NewLoanHandler(loanService)

	// middlewares
	authMW := api.NewAuthMiddleware(authService)

	// mux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", adminHandler.HandlerHealthZ)
	mux.HandleFunc("POST /app/users/create", userHandler.CreateUser)
	mux.HandleFunc("POST /app/login", authHandler.Login)
	mux.HandleFunc("POST /app/refresh", authHandler.Refresh)
	mux.Handle("POST /app/savings/calculate", http.HandlerFunc(savingsHandler.HandleCalculateSavings))
	mux.Handle("POST /app/loans/calculate", http.HandlerFunc(loansHandler.HandleCalculateLoan))
	mux.Handle("POST /app/savings/save", authMW.Handle(http.HandlerFunc(savingsHandler.HandleSaveSavings)))
	mux.Handle("POST /app/loans/save", authMW.Handle(http.HandlerFunc(loansHandler.HandleSaveLoan)))
	mux.Handle("GET /app/savings/list", authMW.Handle(http.HandlerFunc(savingsHandler.HandleListSavings)))
	mux.Handle("GET /app/loans/list", authMW.Handle(http.HandlerFunc(loansHandler.HandleListLoans)))
	mux.Handle("GET /app/savings/{id}", authMW.Handle(http.HandlerFunc(savingsHandler.HandleGetSavings)))
	mux.Handle("GET /app/loans/{id}", authMW.Handle(http.HandlerFunc(loansHandler.HandleGetLoan)))
	mux.Handle("PATCH /app/savings/{id}", authMW.Handle(http.HandlerFunc(savingsHandler.HandleUpdateSavings)))
	mux.Handle("PATCH /app/loans/{id}", authMW.Handle(http.HandlerFunc(loansHandler.HandleUpdateLoan)))
	mux.Handle("DELETE /app/savings/{id}", authMW.Handle(http.HandlerFunc(savingsHandler.HandleDeleteSavings)))
	mux.Handle("DELETE /app/loans/{id}", authMW.Handle(http.HandlerFunc(loansHandler.HandleDeleteLoan)))

	return &App{
		Handler: mux,
		DB:      pool,
	}
}

func (a *App) Run() {
	defer a.DB.Close()
	port := ":8080"
	http.ListenAndServe(port, a.Handler)
}
