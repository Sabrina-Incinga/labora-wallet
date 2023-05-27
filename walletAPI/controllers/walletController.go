package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type WalletController struct {
	CustomerServiceImpl      interfaces.CustomerDBHandler
	WalletServiceImpl        interfaces.WalletDBHandler
	WalletTrackerServiceImpl interfaces.WalletTrackerDBHandler
	WalletCreationServiceImpl interfaces.WalletCreationDBHandler
}

//Method that creates a new wallet if validation requirements are met
func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto model.WalletDTO = model.InitializeWallet()
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		ThrowError(err, w, http.StatusBadRequest)
		return
	}

	validation, rowsAffected, err := c.WalletCreationServiceImpl.AttemptWalletCreation(dto)

	if validation == model.StatusRejected{
		Ok(w, http.StatusOK, "El usuario no pasa las validaciones para la creación de la billetera")
		return
	}
	if rowsAffected == 0{
		ThrowError(fmt.Errorf("La billetera no pudo ser creada"), w, http.StatusBadRequest)
	}
	Ok(w, http.StatusOK, "Billetera creada correctamente")
}

func ThrowError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err)
	http.Error(w, err.Error(), statusCode)
}

func Ok(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// //Method that validates if a customer exists and returns it or creates it 
// func (c *WalletController) createCustomer(dto model.WalletDTO, w http.ResponseWriter) (*model.Customer, bool) {
// 	if dto.CustomerId == 0 {
// 		customerId, err := c.CustomerServiceImpl.CreateCustomer(dto.CustomerDTO, nil)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return nil, true
// 		}

// 		dto.CustomerId = customerId
// 	}

// 	customer, err := c.CustomerServiceImpl.GetCustomerById(dto.CustomerId)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return nil, true
// 	}
// 	if customer == nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Cliente no encontrado"))
// 		return nil, true
// 	}
// 	return customer, false
// }

