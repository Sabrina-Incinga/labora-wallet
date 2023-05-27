package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services"
)

type WalletController struct {
	CustomerServiceImpl services.PostgresCustomerDBHandler;
	WalletServiceImpl services.PostgresWalletDBHandler;
	WalletTrackerServiceImpl services.PostgresWalletTrackerDBHandler;
}

func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var dto model.WalletDTO
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()
	
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	*dto.CreationDate = time.Now()

	if *dto.CustomerId != 0 {
		customer, err := c.CustomerServiceImpl.GetCustomerById(*dto.CustomerId)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if customer != nil {
			
		}
	}
}

func validateScore(customer model.CustomerDTO){
	data := []byte(`{"clave":"valor"}`)
	resp, err := http.Post("https://api.checks.truora.com/v1/checks", "application/www-x-form-urlencoded", bytes.NewBuffer(data))

	if err != nil {
        log.Fatal(err)
    }

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(string(body))

}