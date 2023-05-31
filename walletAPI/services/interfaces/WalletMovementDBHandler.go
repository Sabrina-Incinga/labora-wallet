package interfaces

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type WalletMovementDBHandler interface {
	CreateWalletMovement(movementData dtos.WalletMovementDTO, tx *sql.Tx) (int64, error)
	GetWalletMovementByWalletId(walletId int64) ([]model.WalletMovement, error)
}