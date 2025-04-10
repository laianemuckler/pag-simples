package transfer

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID        string    `json:"id"`
	Value     decimal.Decimal   `json:"value"`
	Payer     int       `json:"payer"`
	Payee     int       `json:"payee"`
}

type Transaction struct {
	ID				string  `json:"id"`
	TransferID string    `json:"transfer_id"`
	Amount     decimal.Decimal `json:"amount"`
	Status     string  `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Notification struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}