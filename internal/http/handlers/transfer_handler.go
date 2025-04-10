package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"pag-simples/internal/transfer"

	"github.com/shopspring/decimal"
)

type TransferHandler struct {
	transferService transfer.TransferUsecase
}

func NewTransferHandler(transferService transfer.TransferUsecase) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var transferRequest struct {
		Value decimal.Decimal `json:"value"`
		Payer int             `json:"payer"`
		Payee int             `json:"payee"`
	}

	if err := json.NewDecoder(r.Body).Decode(&transferRequest); err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	err := h.transferService.Transfer(transferRequest.Value, transferRequest.Payer, transferRequest.Payee)
	if err != nil {
		if strings.Contains(err.Error(), "falha na autorização") {
			http.Error(w, "Você não tem permissão para realizar essa ação.", http.StatusForbidden)
			return
		}

		http.Error(w, fmt.Sprintf("Erro ao realizar a transferência: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transferência realizada com sucesso"))
}
