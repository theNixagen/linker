package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/theNixagen/linker/internal/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	r := chi.NewMux()

	api := api.API{
		Router:    r,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.Router,
	}

	log.Println("Iniciando servidor na porta 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
