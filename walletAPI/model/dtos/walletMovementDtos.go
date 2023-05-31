package dtos

import "time"

type WalletMovementDTO struct {
	SenderWalletId      int64     `json:"sender_wallet_id"`
	ReceiverWalletId 	*int64    `json:"receiver_wallet_id,omitempty"`
	MovementDate        time.Time `json:"movement_date,omitempty"`
	MovementType        string    `json:"movement_type"`
	Amount              float64   `json:"amount"`
}

const (
	DEPOSITMOVEMENT  = "DEPOSIT"
	WITHDRAWMOVEMENT = "WITHDRAWAL"
	TRANSFERMOVEMENT = "TRANSFER"
)

func InitializeWalletMovement() WalletMovementDTO {
	return WalletMovementDTO{
		SenderWalletId: 0,
		ReceiverWalletId: nil,
		MovementDate: time.Now(),
		MovementType: "",
		Amount: 0,
	}
}
