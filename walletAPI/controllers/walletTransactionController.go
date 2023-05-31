package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/labora-wallet/walletAPI/model/dtos"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type WalletTransactionController struct {
	WalletTransactionServiceImpl interfaces.WalletTransactionDBHandler
}

func (t WalletTransactionController) Transfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto dtos.WalletTransactionDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	_, err = t.WalletTransactionServiceImpl.Transfer(dto)

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	Ok(w, http.StatusOK, "Transacción realizada correctamente")
}

func (t WalletTransactionController) Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto dtos.WalletTransactionDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	_, err = t.WalletTransactionServiceImpl.Withdraw(dto)

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	Ok(w, http.StatusOK, "Transacción realizada correctamente")
}

func (t WalletTransactionController) AddToAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto dtos.WalletTransactionDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	_, err = t.WalletTransactionServiceImpl.AddToAccount(dto)

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	Ok(w, http.StatusOK, "Transacción realizada correctamente")
}
