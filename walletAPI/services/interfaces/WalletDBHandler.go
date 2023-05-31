package interfaces

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/db/variablesHandler"
	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type WalletDBHandler interface{
	CreateWallet(wallet dtos.WalletDTO, tx *sql.Tx) (int64, error)
	GetWalletByNumber(walletNumber string) (*model.Wallet, error)
	GetWalletStatusById(id int64) (*dtos.WalletStatusDTO, error)
	GetFullWalletDataById(id int64) (*dtos.WalletDTO, error)
	DeleteWallet(id int64, tx *sql.Tx) (int64, error)
	GetConfig() variablesHandler.DbConfig
}