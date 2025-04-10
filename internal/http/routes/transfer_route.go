package routes

import (
	"pag-simples/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func ConfigureTransferRoutes(r chi.Router, transferHandler *handlers.TransferHandler) {
	r.Post("/transfer", transferHandler.Transfer)
}
