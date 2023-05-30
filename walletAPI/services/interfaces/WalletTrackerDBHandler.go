package interfaces

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model"
)

type WalletTrackerDBHandler interface{
	CreateWalletTracker(tracker model.WalletTrackerDTO, tx *sql.Tx) (int64, error)
	GetWalletTrackByCustomerId(customerId int64) ([]model.WalletTracker, error)
}