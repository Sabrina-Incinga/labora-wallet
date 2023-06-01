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
	"github.com/labora-wallet/walletAPI/model/dtos"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type PostgresWalletAdministrator struct {
	Db                       *sql.DB
	CustomerServiceImpl      interfaces.CustomerDBHandler
	WalletServiceImpl        interfaces.WalletDBHandler
	WalletTrackerServiceImpl interfaces.WalletTrackerDBHandler
}

func (p *PostgresWalletAdministrator) AttemptWalletCreation(wallet dtos.WalletDTO) (string, int64, error) {
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

	customerId, err := p.CustomerServiceImpl.CreateCustomer(*wallet.CustomerDTO, transaction)

	if err != nil {
		return "", rowsAffected, err
	}

	wallet.CustomerId = customerId

	validationResult := p.ValidateScore(wallet.CustomerDTO.NationalIdentityNumber, wallet.CustomerDTO.CountryId)

	trackerDto := dtos.InitializeWalletTracker()
	trackerDto.CreationStatus = validationResult
	trackerDto.CustomerId = customerId
	trackerDto.TrackType = dtos.WALLETCREATION
	
	if validationResult == model.StatusCompleted {
		rowsAffected, err = p.WalletServiceImpl.CreateWallet(wallet, transaction)
		
		if err != nil {
			return validationResult, rowsAffected, err
		}
		trackerDto.RequestStatus = dtos.SUCCESSFULREQUEST
	}else{
		trackerDto.RequestStatus = dtos.FAILEDREQUEST
	}
	_, err = p.WalletTrackerServiceImpl.CreateWalletTracker(trackerDto, transaction)

	if err != nil {
		return validationResult, rowsAffected, err
	}
	
	return validationResult, rowsAffected, nil
}

// Method that returns background check information from Truora API
func (p *PostgresWalletAdministrator) ValidateScore(nationalIdentityNumber, countryId string) string {
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

func (p *PostgresWalletAdministrator) AttemptWalletRemoval(walletId int64, customerId int64) (int64, error){
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return  rowsAffected, err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	trackerDto := dtos.InitializeWalletTracker()
	trackerDto.CreationStatus = "NA"
	trackerDto.CustomerId = customerId
	trackerDto.TrackType = dtos.WALLETREMOVAL
	
	
	rowsAffected, err = p.WalletServiceImpl.DeleteWallet(walletId, transaction)
	
	if err != nil{
		trackerDto.RequestStatus = dtos.FAILEDREQUEST
		_, err = p.WalletTrackerServiceImpl.CreateWalletTracker(trackerDto, nil)
		return rowsAffected, err
	}

	trackerDto.RequestStatus = dtos.SUCCESSFULREQUEST
	if rowsAffected == 0 {
		trackerDto.RequestStatus = dtos.FAILEDREQUEST
	}

	_, err = p.WalletTrackerServiceImpl.CreateWalletTracker(trackerDto, nil)
	if err != nil{
		return rowsAffected, err
	}

	return rowsAffected, nil
}