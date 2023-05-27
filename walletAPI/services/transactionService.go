package services

import (
	"database/sql"
	"fmt"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type PostgresWalletTransactionDBHandler struct {
	Db                       *sql.DB
	WalletServiceImpl        interfaces.WalletDBHandler
}

func (p *PostgresWalletTransactionDBHandler) Transfer(transactionData model.WalletTransactionDTO) (int64, error) {
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	rowsAffected1, err := p.UpdateWalletBalance(transactionData.OriginWalletNumber, model.TRANSACTIONWITHDRAW, transactionData.Amount, transaction)

	if err != nil {
		return rowsAffected, err
	}

	rowsAffected2, err := p.UpdateWalletBalance(transactionData.DestinationWalletNumber, model.TRANSACTIONADD, transactionData.Amount, transaction)

	if err != nil {
		return rowsAffected, err
	}
	rowsAffected = rowsAffected1 + rowsAffected2

	return rowsAffected, nil
}

func (p *PostgresWalletTransactionDBHandler) Withdraw(transactionData model.WalletTransactionDTO) (int64, error) {
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	rowsAffected, err = p.UpdateWalletBalance(transactionData.OriginWalletNumber, model.TRANSACTIONWITHDRAW, transactionData.Amount, transaction)

	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresWalletTransactionDBHandler) AddToAccount(transactionData model.WalletTransactionDTO) (int64, error) {
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	rowsAffected, err = p.UpdateWalletBalance(transactionData.OriginWalletNumber, model.TRANSACTIONADD, transactionData.Amount, transaction)

	if err != nil {
		return rowsAffected, err
	}
	
	return rowsAffected, nil
}

func (p *PostgresWalletTransactionDBHandler) UpdateWalletBalance(walletNumber, transactionType string, amount float64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result

	wallet, err := p.WalletServiceImpl.GetWalletByNumber(walletNumber)
	if err != nil {
		return rowsAffected, err
	}
	
	if wallet == nil {
		return rowsAffected, fmt.Errorf("Billetera de número %s no encontrada", walletNumber)
	}
	if transactionType == model.TRANSACTIONWITHDRAW {
		amount = -amount
	}
	
	query := updateBalanceQuery(amount, wallet.ID)

	if tx != nil {
		response, err = tx.Exec(query)
	} else {
		response, err = p.Db.Exec(query)
	}
	if err != nil {
		return rowsAffected, err
	}
	
	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func updateBalanceQuery(amount float64, id int64) string {
	query := fmt.Sprintf(`UPDATE public.wallet
						SET balance=balance+%f
						WHERE id=%d;`, amount, id)

	return query
}
