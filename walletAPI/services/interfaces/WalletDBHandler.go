package interfaces

import (
	"github.com/labora-wallet/walletAPI/model"
)

type WalletDBHandler interface{
	CreateWallet(wallet model.WalletDTO) (int64, error)
	GetWalletByNumber(walletNumber string) (*model.Wallet, error)
	DeleteWallet(id int) (int64, error)
}