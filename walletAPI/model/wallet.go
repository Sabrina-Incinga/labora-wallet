package model

import "time"

type Wallet struct {
	ID           int64     `json:"id"`
	CustomerId   int64     `json:"customer_id,omitempty"`
	Customer     Customer  `json:"customer"`
	WalletNumber string    `json:"wallet_number"`
	CreationDate time.Time `json:"creation_date,omitempty"`
	Balance      float64   `json:"balance"`
}

type WalletDTO struct {
	CustomerId   int64       `json:"customer_id,omitempty"`
	CustomerDTO  CustomerDTO `json:"customer"`
	WalletNumber string      `json:"wallet_number,omitempty"`
	CreationDate time.Time   `json:"creation_date,omitempty"`
	Balance      float64     `json:"balance"`
}

func InitializeWallet() WalletDTO {
	var customerDto CustomerDTO
	return WalletDTO{
		CustomerId:   0,
		CustomerDTO:  customerDto,
		WalletNumber: "",
		CreationDate: time.Now(),
		Balance:      0,
	}
}
