package transfer

import (
	"fmt"
	"pag-simples/internal/user"
	"pag-simples/pkg/authorization"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) GetUser(userID int) (*user.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserUsecase) SaveUser(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserUsecase) GetAllUsers() ([]user.User, error) {
	args := m.Called()
	users := args.Get(0).([]*user.User)
	result := make([]user.User, len(users))
	for i, u := range users {
		result[i] = *u
	}
	return result, args.Error(1)
}

func (m *MockUserUsecase) ValidateUniqueUser(param1 string, param2 string) error {
	args := m.Called(param1, param2)
	return args.Error(0)
}

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalance(userID int) (decimal.Decimal, error) {
	args := m.Called(userID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockWalletService) CreateWallet(userID int, initialBalance decimal.Decimal) error {
	args := m.Called(userID, initialBalance)
	return args.Error(0)
}

func (m *MockWalletService) UpdateBalance(userID int, newBalance decimal.Decimal) error {
	args := m.Called(userID, newBalance)
	return args.Error(0)
}

type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) UpdateTransactionStatus(transactionID string, status string) error {
	panic("error")
}

func (m *MockTransferRepository) CreateTransfer(transfer *Transfer) error {
	args := m.Called(transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) CreateTransaction(transaction *Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) CheckAuthorization() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

var _ authorization.AuthorizationService = (*MockAuthorizationService)(nil)

func TestTransferSuccess(t *testing.T) {
	userUsecase := new(MockUserUsecase)
	walletService := new(MockWalletService)
	transferRepo := new(MockTransferRepository)
	authorizationService := new(MockAuthorizationService)

	transferService := NewTransferService(userUsecase, walletService, transferRepo, authorizationService)

	payerID := 1
	payeeID := 2
	value := decimal.NewFromFloat(100.0)

	payer := &user.User{ID: payerID, UserType: "common_user", FullName: "Payer Name"}
	payee := &user.User{ID: payeeID, FullName: "Payee Name"}
	userUsecase.On("GetUser", payerID).Return(payer, nil)
	userUsecase.On("GetUser", payeeID).Return(payee, nil)
	walletService.On("GetBalance", payerID).Return(decimal.NewFromFloat(200.0), nil)
	walletService.On("GetBalance", payeeID).Return(decimal.NewFromFloat(50.0), nil)
	authorizationService.On("CheckAuthorization").Return(true, nil)
	transferRepo.On("CreateTransfer", mock.Anything).Return(nil)
	walletService.On("UpdateBalance", payerID, decimal.NewFromFloat(100.0)).Return(nil)
	walletService.On("UpdateBalance", payeeID, decimal.NewFromFloat(150.0)).Return(nil)
	transferRepo.On("CreateTransaction", mock.Anything).Return(nil)

	err := transferService.Transfer(value, payerID, payeeID)

	assert.NoError(t, err)
	userUsecase.AssertExpectations(t)
	walletService.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
	authorizationService.AssertExpectations(t)
}

func TestTransferErrorInsufficientBalance(t *testing.T) {
	userUsecase := new(MockUserUsecase)
	walletService := new(MockWalletService)
	transferRepo := new(MockTransferRepository)
	authorizationService := new(MockAuthorizationService)

	transferService := NewTransferService(userUsecase, walletService, transferRepo, authorizationService)

	payerID := 1
	payeeID := 2
	value := decimal.NewFromFloat(100.0)

	payer := &user.User{ID: payerID, UserType: "common_user", FullName: "Payer Name"}
	payee := &user.User{ID: payeeID, FullName: "Payee Name"}
	userUsecase.On("GetUser", payerID).Return(payer, nil)
	userUsecase.On("GetUser", payeeID).Return(payee, nil)

	walletService.On("GetBalance", payerID).Return(decimal.NewFromFloat(50.0), nil)
	walletService.On("GetBalance", payeeID).Return(decimal.NewFromFloat(50.0), nil)

	err := transferService.Transfer(value, payerID, payeeID)

	assert.Error(t, err)
	assert.Equal(t, "saldo insuficiente para a transferência", err.Error())

	userUsecase.AssertExpectations(t)
	walletService.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
	authorizationService.AssertExpectations(t)
}


func TestTransferErrorAuthorizationFailed(t *testing.T) {
	userUsecase := new(MockUserUsecase)
	walletService := new(MockWalletService)
	transferRepo := new(MockTransferRepository)
	authorizationService := new(MockAuthorizationService)

	transferService := NewTransferService(userUsecase, walletService, transferRepo, authorizationService)

	payerID := 1
	payeeID := 2
	value := decimal.NewFromFloat(100.0)

	payer := &user.User{ID: payerID, UserType: "common_user", FullName: "Payer Name"}
	payee := &user.User{ID: payeeID, FullName: "Payee Name"}
	userUsecase.On("GetUser", payerID).Return(payer, nil)
	userUsecase.On("GetUser", payeeID).Return(payee, nil)
	walletService.On("GetBalance", payerID).Return(decimal.NewFromFloat(200.0), nil)
	walletService.On("GetBalance", payeeID).Return(decimal.NewFromFloat(50.0), nil)
	authorizationService.On("CheckAuthorization").Return(false, nil)

	err := transferService.Transfer(value, payerID, payeeID)

	assert.Error(t, err)
	assert.Equal(t, "transferência não autorizada", err.Error())

	userUsecase.AssertExpectations(t)
	walletService.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
	authorizationService.AssertExpectations(t)
}

func TestTransferErrorRepositorySave(t *testing.T) {
	userUsecase := new(MockUserUsecase)
	walletService := new(MockWalletService)
	transferRepo := new(MockTransferRepository)
	authorizationService := new(MockAuthorizationService)

	transferService := NewTransferService(userUsecase, walletService, transferRepo, authorizationService)

	payerID := 1
	payeeID := 2
	value := decimal.NewFromFloat(100.0)

	payer := &user.User{ID: payerID, UserType: "common_user", FullName: "Payer Name"}
	payee := &user.User{ID: payeeID, FullName: "Payee Name"}
	userUsecase.On("GetUser", payerID).Return(payer, nil)
	userUsecase.On("GetUser", payeeID).Return(payee, nil)
	walletService.On("GetBalance", payerID).Return(decimal.NewFromFloat(200.0), nil)
	walletService.On("GetBalance", payeeID).Return(decimal.NewFromFloat(50.0), nil)
	authorizationService.On("CheckAuthorization").Return(true, nil)
	transferRepo.On("CreateTransfer", mock.Anything).Return(fmt.Errorf("erro ao salvar transferência"))

	err := transferService.Transfer(value, payerID, payeeID)

	assert.Error(t, err)
	assert.Equal(t, "falha ao salvar a transferência: erro ao salvar transferência", err.Error())

	userUsecase.AssertExpectations(t)
	walletService.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
	authorizationService.AssertExpectations(t)
}

func TestTransferErrorUpdateBalance(t *testing.T) {
	userUsecase := new(MockUserUsecase)
	walletService := new(MockWalletService)
	transferRepo := new(MockTransferRepository)
	authorizationService := new(MockAuthorizationService)

	transferService := NewTransferService(userUsecase, walletService, transferRepo, authorizationService)

	payerID := 1
	payeeID := 2
	value := decimal.NewFromFloat(100.0)

	payer := &user.User{ID: payerID, UserType: "common_user", FullName: "Payer Name"}
	payee := &user.User{ID: payeeID, FullName: "Payee Name"}
	userUsecase.On("GetUser", payerID).Return(payer, nil)
	userUsecase.On("GetUser", payeeID).Return(payee, nil)
	walletService.On("GetBalance", payerID).Return(decimal.NewFromFloat(200.0), nil)
	walletService.On("GetBalance", payeeID).Return(decimal.NewFromFloat(50.0), nil)
	authorizationService.On("CheckAuthorization").Return(true, nil)
	transferRepo.On("CreateTransfer", mock.Anything).Return(nil)
	walletService.On("UpdateBalance", payerID, decimal.NewFromFloat(100.0)).Return(fmt.Errorf("erro ao atualizar saldo do pagador"))

	err := transferService.Transfer(value, payerID, payeeID)

	assert.Error(t, err)
	assert.Equal(t, "falha ao atualizar o saldo do pagador 1: erro ao atualizar saldo do pagador", err.Error())

	userUsecase.AssertExpectations(t)
	walletService.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
	authorizationService.AssertExpectations(t)
}
