package dtos

import "time"

type WalletTrackerDTO struct {
	CustomerId     int64     `json:"customer_id"`
	RecordDate     time.Time `json:"record_date,omitempty"`
	CreationStatus string    `json:"creation_status"`
	TrackType      string    `json:"track_type"`
	RequestStatus  string    `json:"request_status"`
}

func InitializeWalletTracker() WalletTrackerDTO {
	return WalletTrackerDTO{
		CustomerId:     0,
		RecordDate:     time.Now(),
		CreationStatus: "",
		TrackType:      "",
		RequestStatus:  "",
	}
}

const (
	WALLETMOVEMENT = "WALLET MOVEMENT"
	WALLETCREATION = "WALLET CREATION"
	WALLETREMOVAL  = "WALLET REMOVAL"
)

const (
	FAILEDREQUEST     = "FAILED"
	SUCCESSFULREQUEST = "SUCCESSFUL"
)