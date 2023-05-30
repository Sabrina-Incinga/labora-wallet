package interfaces

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
)

type WalletMovementDBHandler interface {
	CreateWalletMovement(movementData model.WalletMovementDTO, tx *sql.Tx) (int64, error)
	GetWalletMovementByWalletId(walletId int64) ([]model.WalletMovement, error)
}