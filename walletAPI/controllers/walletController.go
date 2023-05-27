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
	CustomerServiceImpl      services.PostgresCustomerDBHandler
	WalletServiceImpl        services.PostgresWalletDBHandler
	WalletTrackerServiceImpl services.PostgresWalletTrackerDBHandler
}

type ValidationInfo struct {
	Check struct {
		CheckID      string `json:"check_id"`
		CreationDate string `json:"creation_date"`
		Score        int    `json:"score"`
	} `json:"check"`
}

const (
	StatusRejected  = "RECHAZADO"
	StatusCompleted = "COMPLETADO"
)

func (c *WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//TODO validar que el usuario no tiene una billetera ya asociada
	var dto model.WalletDTO = model.WalletDTO{
		CustomerId:   nil,
		CustomerDTO:  nil,
		WalletNumber: nil,
		CreationDate: new(time.Time),
		Balance:      0.0,
	}
	err := json.NewDecoder(r.Body).Decode(&dto)

	defer r.Body.Close()

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	*dto.CreationDate = time.Now()

	if dto.CustomerId == nil {
		customerId, err := c.CustomerServiceImpl.CreateCustomer(*dto.CustomerDTO)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		*dto.CustomerId = customerId
	}

	customer, err := c.CustomerServiceImpl.GetCustomerById(*dto.CustomerId)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if customer == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cliente no encontrado"))
		return
	}

	validation := c.validateScore(customer.NationalIdentityNumber, customer.CountryId)

	if validation.Check.Score != 1 {
		_, err = c.WalletTrackerServiceImpl.CreateWalletTracker(model.WalletTrackerDTO{CustomerId: *dto.CustomerId, RecordDate: dto.CreationDate, CreationStatus: StatusRejected})
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("El documento consultado no pasa las validaciones para la creación de la billetera virtual"))
		return
	}

	_, err = c.WalletServiceImpl.CreateWallet(dto)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.WalletTrackerServiceImpl.CreateWalletTracker(model.WalletTrackerDTO{CustomerId: *dto.CustomerId, RecordDate: dto.CreationDate, CreationStatus: StatusCompleted})
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Billetera creada correctamente"))

}

func (c *WalletController) validateScore(nationalIdentityNumber, countryId string) ValidationInfo {
	data := url.Values{}
	data.Set("national_id", nationalIdentityNumber)
	data.Set("country", countryId)
	data.Set("type", "person")
	data.Set("user_authorized", strconv.FormatBool(true))

	// Codificar los datos en una cadena en formato application/x-www-form-urlencoded
	body := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", c.WalletServiceImpl.Config.BackgroundChecksUrl, body)
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

	var responseMap ValidationInfo

	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		log.Fatal(err)
	}

	return responseMap

}
