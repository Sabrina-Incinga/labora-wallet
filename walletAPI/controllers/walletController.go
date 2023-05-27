package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

func (c *WalletController) validateScore(customer model.CustomerDTO){
	data := url.Values{}
	data.Set("national_id", customer.NationalIdentityNumber)
	data.Set("country", customer.CountryId)
	data.Set("type", "person")
	data.Set("user_authorized", strconv.FormatBool(true))

	// Codificar los datos en una cadena en formato application/x-www-form-urlencoded
	body := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", "https://api.checks.truora.com/v1/checks", body)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/www-x-form-urlencoded")
    req.Header.Set("Accept", "*/*")
	req.Header.Set("Truora-Api-Key", c.WalletServiceImpl.Config.ApiKey)

	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(string(responseBody))


}