package transfer

import "github.com/shopspring/decimal"

type TransferUsecase interface {
	Transfer(value decimal.Decimal, payerID int, payeeID int) error
}