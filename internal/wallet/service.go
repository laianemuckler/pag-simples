package wallet

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type WalletService struct {
	repo *MemoryWalletRepository
}

func NewWalletService(repo *MemoryWalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) CreateWallet(userID int, balance decimal.Decimal) error {
	if _, exists := s.repo.wallets[userID]; exists {
		return fmt.Errorf("wallet já existe para o usuário %d", userID)
	}

	wallet := &Wallet{
		ID:      userID,
		Balance: balance,
	}

	s.repo.wallets[userID] = wallet
	return nil
}

func (s *WalletService) GetBalance(userID int) (decimal.Decimal, error) {
	return s.repo.GetBalance(userID)
}

func (s *WalletService) UpdateBalance(userID int, amount decimal.Decimal) error {
	return s.repo.UpdateBalance(userID, amount)
}
