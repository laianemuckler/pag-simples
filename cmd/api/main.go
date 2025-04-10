package main

import (
	"log"
	"net/http"

	"pag-simples/internal/http/handlers"
	"pag-simples/internal/http/routes"
	"pag-simples/internal/transfer"
	"pag-simples/internal/user"
	"pag-simples/internal/wallet"
	"pag-simples/pkg/authorization"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shopspring/decimal"
)

func initializeData(userRepo *user.MemoryUserRepository, walletRepo *wallet.MemoryWalletRepository) {
	userRepo.SaveUser(&user.User{
		ID:             1,
		FullName:       "João Silva",
		DocumentNumber: "12345678901",
		Email:          "joao@email.com",
		Password:       "senha123",
		UserType:       user.CommonUser,
	})
	userRepo.SaveUser(&user.User{
		ID:             2,
		FullName:       "Maria Oliveira",
		DocumentNumber: "98765432100",
		Email:          "maria@email.com",
		Password:       "senha456",
		UserType:       user.CommonUser,
	})
	userRepo.SaveUser(&user.User{
		ID:             3,
		FullName:       "Loja Exemplo",
		DocumentNumber: "12345678000100",
		Email:          "loja@email.com",
		Password:       "senha789",
		UserType:       user.Merchant,
	})

	walletRepo.CreateWallet(1, decimal.NewFromFloat(1000.0))
	walletRepo.CreateWallet(2, decimal.NewFromFloat(500.0))
	walletRepo.CreateWallet(3, decimal.NewFromFloat(2000.0))
}

func main() {
	userRepo := user.NewMemoryUserRepository()
	walletRepo := wallet.NewMemoryWalletRepository()
	transferRepo := transfer.NewMemoryTransferRepository()
	authorizationService := authorization.NewAuthorizationService()

	userService := user.NewUserService(userRepo)
	walletService := wallet.NewWalletService(walletRepo)

	userHandler := handlers.NewUserHandler(userService, walletService)

	transferService := transfer.NewTransferService(userService, walletRepo, transferRepo, authorizationService)
	transferHandler := handlers.NewTransferHandler(transferService)

	initializeData(userRepo, walletRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.ConfigureUserRoutes(r, userHandler)
	routes.ConfigureTransferRoutes(r, transferHandler)

	log.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
