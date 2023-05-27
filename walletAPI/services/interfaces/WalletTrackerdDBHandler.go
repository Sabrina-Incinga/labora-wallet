package interfaces

import "github.com/labora-wallet/walletAPI/model"

type WalletTrackerDBHandler interface{
	CreateWalletTracker(tracker model.WalletTrackerDTO) (int64, error)
	GetWalletTrackByCustomerId(customerId int64) (*[]model.WalletTracker, error)
}