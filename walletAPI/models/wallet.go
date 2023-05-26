package models

import "time"

type Wallet struct{
	ID            int    		`json:"id"`
	CustomerId    int    		`json:"customer_id"`
	CreationDate  *time.Time 	`json:"creation_date"`
	Balance float64  			`json:"balance"`
}