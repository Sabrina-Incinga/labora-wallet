package services

import (
	"database/sql"
	"fmt"
	"sync"
	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type PostgresWalletTransactionDBHandler struct {
	Db                *sql.DB
	WalletServiceImpl interfaces.WalletDBHandler
	Mutex             *sync.Mutex
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

	p.Mutex.Lock()

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

	query := updateBalanceQuery()

	if tx != nil {
		response, err = tx.Exec(query, amount, wallet.ID)
	} else {
		response, err = p.Db.Exec(query, amount, wallet.ID)
	}
	if err != nil {
		return rowsAffected, err
	}
	p.Mutex.Unlock()
	
	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func updateBalanceQuery() string {
	query := `UPDATE public.wallet
						SET balance=balance+$1
						WHERE id=$2;`

	return query
}
