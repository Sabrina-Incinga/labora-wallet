package model

import "time"

type WalletTracker struct {
	ID             int64     `json:"id"`
	CustomerId     int64     `json:"customer_id"`
	RecordDate     time.Time `json:"record_date,omitempty"`
	CreationStatus string    `json:"creation_status"`
}

type WalletTrackerDTO struct {
	CustomerId     int64     `json:"customer_id"`
	RecordDate     time.Time `json:"record_date,omitempty"`
	CreationStatus string    `json:"creation_status"`
}

func InitializeWalletTracker() WalletTrackerDTO{
	return WalletTrackerDTO{
		CustomerId: 0,
		RecordDate: time.Now(),
		CreationStatus: "",
	}
}