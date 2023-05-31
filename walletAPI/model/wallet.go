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


