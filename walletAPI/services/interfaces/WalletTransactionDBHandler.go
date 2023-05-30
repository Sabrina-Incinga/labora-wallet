package interfaces

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
)

type WalletTransactionDBHandler interface {
	Transfer(transactionData model.WalletTransactionDTO) (int64, error)
	Withdraw(transactionData model.WalletTransactionDTO) (int64, error)
	AddToAccount(transactionData model.WalletTransactionDTO) (int64, error)
	UpdateWalletBalance(walletNumber, transactionType string, amount float64, tx *sql.Tx) (int64, error)
}
