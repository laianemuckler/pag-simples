package wallet

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type MemoryWalletRepository struct {
	wallets map[int]*Wallet
}

func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets: make(map[int]*Wallet),
	}
}

func (r *MemoryWalletRepository) GetBalance(userID int) (decimal.Decimal, error) {
	wallet, exists := r.wallets[userID]
	if !exists {
		return decimal.Zero, fmt.Errorf("wallet não encontrada para o usuário %d", userID)
	}
	return wallet.Balance, nil
}

func (r *MemoryWalletRepository) UpdateBalance(userID int, amount decimal.Decimal) error {
	wallet, exists := r.wallets[userID]
	if !exists {
		return fmt.Errorf("wallet não encontrada para o usuário %d", userID)
	}
	wallet.Balance = wallet.Balance.Add(amount)
	return nil
}

func (r *MemoryWalletRepository) CreateWallet(userID int, initialBalance decimal.Decimal) error {
	r.wallets[userID] = &Wallet{
		ID:      userID,
		Balance: initialBalance,
	}
	return nil
}
