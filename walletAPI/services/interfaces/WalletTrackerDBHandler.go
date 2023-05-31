package interfaces

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type WalletTrackerDBHandler interface{
	CreateWalletTracker(tracker dtos.WalletTrackerDTO, tx *sql.Tx) (int64, error)
	GetWalletTrackByCustomerId(customerId int64) ([]model.WalletTracker, error)
}