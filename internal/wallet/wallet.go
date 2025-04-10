package wallet

import "github.com/shopspring/decimal"

type Wallet struct {
	ID      int
	Balance decimal.Decimal
}