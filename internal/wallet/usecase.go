package wallet

import "github.com/shopspring/decimal"

type WalletUseCase interface {
	CreateWallet(int, decimal.Decimal) error
	GetBalance(userID int) (decimal.Decimal, error)
	UpdateBalance(userID int, amount decimal.Decimal) error
}
