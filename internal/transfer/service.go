package transfer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"pag-simples/internal/user"
	"pag-simples/internal/wallet"
	"pag-simples/pkg/authorization"
	"pag-simples/pkg/notification"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var mu sync.Mutex

type TransferService struct {
	userUsecase          user.UserUsecase
	walletService        wallet.WalletUseCase
	transferRepo         TransferRepository
	authorizationService authorization.AuthorizationService
}

func NewTransferService(
	userUsecase user.UserUsecase,
	walletService wallet.WalletUseCase,
	transferRepo TransferRepository,
	authorizationService authorization.AuthorizationService,
) TransferUsecase {
	return &TransferService{
		userUsecase:          userUsecase,
		walletService:        walletService,
		transferRepo:         transferRepo,
		authorizationService: authorizationService,
	}
}

func (s *TransferService) Transfer(value decimal.Decimal, payerID int, payeeID int) error {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Iniciando transferência de %.2f de %d para %d", value.InexactFloat64(), payerID, payeeID)

	payer, err := s.userUsecase.GetUser(payerID)
	if err != nil {
		log.Printf("Erro ao encontrar pagador %d: %v", payerID, err)
		return fmt.Errorf("pagador não encontrado: %v", err)
	}

	payee, err := s.userUsecase.GetUser(payeeID)
	if err != nil {
		log.Printf("Erro ao encontrar recebedor %d: %v", payeeID, err)
		return fmt.Errorf("recebedor não encontrado: %v", err)
	}

	if payer.UserType == "merchant" {
		log.Printf("Erro: usuário %d é um lojista e não pode realizar transferência", payerID)
		return fmt.Errorf("um lojista não pode realizar transferências")
	}

	payerBalance, err := s.walletService.GetBalance(payerID)
	if err != nil {
		log.Printf("Falha ao obter saldo do pagador %d: %v", payerID, err)
		return fmt.Errorf("falha ao obter o saldo do pagador: %v", err)
	}

	payeeBalance, err := s.walletService.GetBalance(payeeID)
	if err != nil {
		log.Printf("Falha ao obter saldo do recebedor %d: %v", payeeID, err)
		return fmt.Errorf("falha ao obter o saldo do recebedor: %v", err)
	}

	if payerBalance.LessThan(value) {
		log.Printf("Erro: saldo insuficiente para a transferência de %.2f de %d para %d", value.InexactFloat64(), payerID, payeeID)
		return fmt.Errorf("saldo insuficiente para a transferência")
	}

	authorized, err := s.authorizationService.CheckAuthorization()
	if err != nil {
		log.Printf("Falha na autorização: %v", err)
		return fmt.Errorf("falha na autorização: %v", err)
	}

	if !authorized {
		log.Printf("Falha na autorização da transferência de %.2f de %d para %d", value.InexactFloat64(), payerID, payeeID)
		return fmt.Errorf("transferência não autorizada")
	}

	transfer := &Transfer{
		ID:    generateID(),
		Value: value,
		Payer: payerID,
		Payee: payeeID,
	}

	err = s.transferRepo.CreateTransfer(transfer)
	if err != nil {
		log.Printf("Falha ao salvar a transferência de %.2f: %v", value.InexactFloat64(), err)
		return fmt.Errorf("falha ao salvar a transferência: %v", err)
	}

	err = s.walletService.UpdateBalance(payerID, payerBalance.Sub(value))
	if err != nil {
		log.Printf("Falha ao atualizar saldo do pagador %d: %v", payerID, err)
		return fmt.Errorf("falha ao atualizar o saldo do pagador %d: %v", payerID, err)
	}

	err = s.walletService.UpdateBalance(payeeID, payeeBalance.Add(value))
	if err != nil {
		log.Printf("Falha ao atualizar saldo do recebedor %d: %v", payeeID, err)
		return fmt.Errorf("falha ao atualizar o saldo do recebedor %d: %v", payeeID, err)
	}

	transaction := &Transaction{
		ID:         generateID(),
		TransferID: transfer.ID,
		Amount:     value,
		Status:     "sucesso",
		CreatedAt:  time.Now(),
	}

	err = s.transferRepo.CreateTransaction(transaction)
	if err != nil {
		log.Printf("Falha ao salvar transação de %.2f: %v", value.InexactFloat64(), err)
		return fmt.Errorf("falha ao salvar a transação: %v", err)
	}

	log.Printf("Transferência de %.2f realizada com sucesso de %d para %d", value.InexactFloat64(), payerID, payeeID)

	go s.notifyUser(payer, fmt.Sprintf("Transferência de %.2f para %s foi realizada com sucesso", value.InexactFloat64(), payee.FullName))
	go s.notifyUser(payee, fmt.Sprintf("Você recebeu %.2f de %s", value.InexactFloat64(), payer.FullName))

	return nil
}

func generateID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *TransferService) notifyUser(user *user.User, message string) error {
	notificationRequest := notification.NotificationRequest{
		Email:   user.Email,
		Message: message,
	}

	err := notification.SendNotification(notificationRequest)
	if err != nil {
		log.Printf("Falha ao enviar notificação para o usuário %d: %v", user.ID, err)
		return fmt.Errorf("falha ao enviar a notificação: %v", err)
	}

	log.Printf("Notificação enviada ao usuário %d: %s", user.ID, message)
	return nil
}
