package routes

import (
	"pag-simples/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func ConfigureUserRoutes(r chi.Router, userHandler *handlers.UserHandler) {
	r.Get("/users/{id}", userHandler.GetUser)
	r.Get("/users", userHandler.GetAllUsers)
	r.Post("/users", userHandler.CreateUser)
}
