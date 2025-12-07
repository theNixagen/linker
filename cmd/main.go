package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/theNixagen/linker/internal/api"
	"github.com/theNixagen/linker/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("FINANSMART_DATABASE_USER"),
		os.Getenv("FINANSMART_DATABASE_PASSWORD"),
		os.Getenv("FINANSMART_DATABASE_HOST"),
		os.Getenv("FINANSMART_DATABASE_PORT"),
		os.Getenv("FINANSMART_DATABASE_NAME"),
	))

	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewMux()
	jwtSecret := os.Getenv("JWT_SECRET")

	api := api.API{
		Router:      r,
		Validator:   validator.New(validator.WithRequiredStructEnabled()),
		UserService: services.NewUserService(pool, jwtSecret),
		JwtSecret:   jwtSecret,
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.Router,
	}

	api.BindRoutes()

	log.Println("Iniciando servidor na porta 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
