package interfaces

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/db/variablesHandler"
	"github.com/labora-wallet/walletAPI/model"
)

type WalletDBHandler interface{
	CreateWallet(wallet model.WalletDTO, tx *sql.Tx) (int64, error)
	GetWalletByNumber(walletNumber string) (*model.Wallet, error)
	GetWalletStatusById(id int64) (string, error)
	DeleteWallet(id int64, tx *sql.Tx) (int64, error)
	GetConfig() variablesHandler.DbConfig
}