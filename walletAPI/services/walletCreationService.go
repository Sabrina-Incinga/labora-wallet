package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type PostgresWalletCreationtDBHandler struct {
	Db                       *sql.DB
	CustomerServiceImpl      interfaces.CustomerDBHandler
	WalletServiceImpl        interfaces.WalletDBHandler
	WalletTrackerServiceImpl interfaces.WalletTrackerDBHandler
}

func (p *PostgresWalletCreationtDBHandler) AttemptWalletCreation(wallet model.WalletDTO) (string, int64, error) {
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return "", rowsAffected, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	customerId, err := p.CustomerServiceImpl.CreateCustomer(wallet.CustomerDTO, transaction)

	if err != nil {
		return "", rowsAffected, err
	}

	wallet.CustomerId = customerId

	validationResult := p.ValidateScore(wallet.CustomerDTO.NationalIdentityNumber, wallet.CustomerDTO.CountryId)

	trackerDto := model.InitializeWalletTracker()
	trackerDto.CreationStatus = validationResult
	trackerDto.CustomerId = customerId
	_, err = p.WalletTrackerServiceImpl.CreateWalletTracker(trackerDto, transaction)

	if err != nil {
		return validationResult, rowsAffected, err
	}

	if validationResult == model.StatusCompleted {
		rowsAffected, err = p.WalletServiceImpl.CreateWallet(wallet, transaction)

		if err != nil {
			return validationResult, rowsAffected, err
		}
	}

	if err != nil {
		return validationResult, rowsAffected, err
	}

	return validationResult, rowsAffected, nil
}

// Method that returns background check information from Truora API
func (p *PostgresWalletCreationtDBHandler) ValidateScore(nationalIdentityNumber, countryId string) string {
	data := url.Values{}
	config := p.WalletServiceImpl.GetConfig()
	data.Set("national_id", nationalIdentityNumber)
	data.Set("country", countryId)
	data.Set("type", "person")
	data.Set("user_authorized", strconv.FormatBool(true))

	body := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", config.BackgroundChecksUrl, body)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/www-x-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Truora-Api-Key", config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var responseMap model.ValidationInfo

	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		log.Println(err)
	}

	if responseMap.Check.Score != 1 {
		return model.StatusRejected
	}
	return model.StatusCompleted

}
