package model

import "time"

type Wallet struct{
	ID            int    		`json:"id"`
	CustomerId    *int    		`json:"customer_id"`
	Customer	  *Customer		`json:"customer"`
	WalletNumber  string		`json:"wallet_number"`
	CreationDate  *time.Time	`json:"creation_date"`
	Balance 	  float64  		`json:"balance"`
}

type WalletDTO struct{
	CustomerId    *int    		`json:"customer_id"`
	CustomerDTO	  *CustomerDTO	`json:"customer"`
	WalletNumber  *string		`json:"wallet_number"`
	CreationDate  *time.Time	`json:"creation_date"`
	Balance 	  float64  		`json:"balance"`
}