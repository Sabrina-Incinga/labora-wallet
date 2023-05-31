package model

import "time"

type WalletMovement struct {
	ID               int64     `json:"id"`
	SenderWalletId   int64     `json:"sender_wallet_id"`
	ReceiverWalletId int64     `json:"receiver_wallet_id,omitempty"`
	MovementDate     time.Time `json:"movement_date,omitempty"`
	MovementType     string    `json:"movement_type"`
	Amount           float64   `json:"amount"`
}
