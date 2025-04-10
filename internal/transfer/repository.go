package transfer

import "fmt"

type TransferRepository interface {
	CreateTransfer(transfer *Transfer) error
	CreateTransaction(transaction *Transaction) error
	UpdateTransactionStatus(transactionID string, status string) error
}

type MemoryTransferRepository struct {
	transfers    map[string]Transfer
	transactions map[string]Transaction
}

func NewMemoryTransferRepository() *MemoryTransferRepository {
	return &MemoryTransferRepository{
		transfers:    make(map[string]Transfer),
		transactions: make(map[string]Transaction),
	}
}

func (r *MemoryTransferRepository) CreateTransfer(transfer *Transfer) error {
	r.transfers[transfer.ID] = *transfer
	return nil
}

func (r *MemoryTransferRepository) CreateTransaction(transaction *Transaction) error {
	r.transactions[transaction.ID] = *transaction
	return nil
}

func (r *MemoryTransferRepository) UpdateTransactionStatus(transactionID string, status string) error {
	transaction, exists := r.transactions[transactionID]
	if !exists {
		return fmt.Errorf("transaction not found")
	}
	transaction.Status = status
	r.transactions[transactionID] = transaction
	return nil
}
