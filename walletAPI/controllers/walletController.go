package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/labora-wallet/walletAPI/model/dtos"
	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type WalletController struct {
	CustomerServiceImpl            interfaces.CustomerDBHandler
	WalletServiceImpl              interfaces.WalletDBHandler
	WalletTrackerServiceImpl       interfaces.WalletTrackerDBHandler
	WalletAdministratorServiceImpl interfaces.WalletAdministratorHandler
}

// Method that creates a new wallet if validation requirements are met
func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto dtos.WalletDTO = dtos.InitializeWallet()
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	validation, rowsAffected, err := c.WalletAdministratorServiceImpl.AttemptWalletCreation(dto)

	if validation == model.StatusRejected {
		Ok(w, http.StatusOK, "El usuario no pasa las validaciones para la creación de la billetera")
		return
	}
	if rowsAffected == 0 {
		ThrowError(fmt.Errorf("La billetera no pudo ser creada"), w, http.StatusBadRequest)
		return
	}
	Ok(w, http.StatusOK, "Billetera creada correctamente")
}

func (c *WalletController) GetWalletStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	walletStatus, err := c.WalletServiceImpl.GetWalletStatusById(int64(id))

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}
	if walletStatus == nil {
		Ok(w, http.StatusOK, fmt.Sprintf(`La billetera de id %d no se encontró`, id))
		return
	}

	Ok(w, http.StatusOK, fmt.Sprintf(`El status de la billetera de id %d es: %s`, id, walletStatus.Status))
}

func (c *WalletController) GetWalletById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	wallet, err := c.WalletServiceImpl.GetFullWalletDataById(int64(id))

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}
	if wallet == nil {
		Ok(w, http.StatusOK, fmt.Sprintf(`La billetera de id %d no se encontró`, id))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

func (c *WalletController) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	rowsAffected, err := c.WalletAdministratorServiceImpl.AttemptWalletRemoval(int64(id))

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	if rowsAffected == 0 {
		Ok(w, http.StatusNoContent, "No se pudo eliminar la billetera seleccionada")
		return
	}
	Ok(w, http.StatusNoContent, "Billetera eliminada correctamente")
}

func ThrowError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}

func Ok(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
