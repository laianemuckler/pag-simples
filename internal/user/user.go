package user

import "pag-simples/internal/wallet"

type UserType string

const (
	CommonUser UserType = "common_user"
	Merchant   UserType = "merchant"
)

type User struct {
	ID             int
	FullName       string
	DocumentNumber string
	Email          string
	Password       string
	UserType       UserType
	Wallet         wallet.Wallet
}
