package dtos

import (
	"time"

	"github.com/labora-wallet/walletAPI/model"
)

type WalletDTO struct {
	CustomerId   int64               `json:"customer_id,omitempty"`
	CustomerDTO  *CustomerDTO         `json:"customer,omitempty"`
	WalletNumber string              `json:"wallet_number,omitempty"`
	CreationDate time.Time           `json:"creation_date,omitempty"`
	Balance      float64             `json:"balance"`
	Movements    []WalletMovementDTO `json:"movements,omitempty"`
}

type WalletTransactionDTO struct {
	OriginWalletNumber      string  `json:"origin_wallet_number"`
	DestinationWalletNumber string  `json:"destination_wallet_number,omitempty"`
	Amount                  float64 `json:"amount"`
}

type WalletStatusDTO struct {
	Wallet model.Wallet `json:"wallet"`
	Status string       `json:"status"`
}

func InitializeWallet() WalletDTO {
	var customerDto CustomerDTO
	return WalletDTO{
		CustomerId:   0,
		CustomerDTO:  &customerDto,
		WalletNumber: "",
		CreationDate: time.Now(),
		Balance:      0,
	}
}
