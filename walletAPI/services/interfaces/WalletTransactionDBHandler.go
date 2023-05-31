package interfaces

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type WalletTransactionDBHandler interface {
	Transfer(transactionData dtos.WalletTransactionDTO) (int64, error)
	Withdraw(transactionData dtos.WalletTransactionDTO) (int64, error)
	AddToAccount(transactionData dtos.WalletTransactionDTO) (int64, error)
	UpdateWalletBalance(walletNumber, transactionType string, amount float64, tx *sql.Tx) (int64, error)
}
