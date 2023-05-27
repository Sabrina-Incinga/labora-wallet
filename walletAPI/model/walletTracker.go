package model

import "time"

type WalletTracker struct{
	ID            	int    		`json:"id"`
	CustomerId    	int    		`json:"customer_id"`
	RecordDate		*time.Time	`json:"record_date"`
	CreationStatus 	string  	`json:"creation_status"`
}

type WalletTrackerDTO struct{
	CustomerId    	int    		`json:"customer_id"`
	RecordDate		*time.Time	`json:"record_date"`
	CreationStatus 	string  	`json:"creation_status"`
}