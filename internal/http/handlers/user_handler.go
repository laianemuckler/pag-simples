package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"

	"pag-simples/internal/user"
	"pag-simples/internal/wallet"
)

type UserHandler struct {
	userService   user.UserUsecase
	walletService wallet.WalletUseCase
}
func NewUserHandler(userService user.UserUsecase, walletService wallet.WalletUseCase) *UserHandler {
	return &UserHandler{
		userService:   userService,
		walletService: walletService,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	for i := range users {
			balance, err := h.walletService.GetBalance(users[i].ID)
			if err != nil {
					log.Printf("Erro ao obter saldo da carteira para o usuário %d: %v", users[i].ID, err)
					users[i].Wallet = wallet.Wallet{
							ID:      users[i].ID,
							Balance: decimal.Zero,
					}
			} else {
					users[i].Wallet = wallet.Wallet{
							ID:      users[i].ID,
							Balance: balance,
					}
			}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}


func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	if err := h.userService.SaveUser(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
