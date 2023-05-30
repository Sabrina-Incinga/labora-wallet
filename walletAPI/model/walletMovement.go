package model

import "time"

type WalletMovement struct {
	ID           int64     `json:"id"`
	WalletId     int64     `json:"wallet_id"`
	MovementDate time.Time `json:"movement_date,omitempty"`
	MovementType string    `json:"movement_type"`
	Amount       float64   `json:"amount"`
}

type WalletMovementDTO struct {
	WalletId     int64     `json:"wallet_id"`
	MovementDate time.Time `json:"movement_date,omitempty"`
	MovementType string    `json:"movement_type"`
	Amount       float64   `json:"amount"`
}

func InitializeWalletMovement() WalletMovementDTO {
	return WalletMovementDTO{
		WalletId:     0,
		MovementDate: time.Now(),
		MovementType: "",
		Amount:       0,
	}
}
