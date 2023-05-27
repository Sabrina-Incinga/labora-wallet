package services

import (
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"time"
	"github.com/labora-wallet/walletAPI/db/variablesHandler"
	"github.com/labora-wallet/walletAPI/model"
)

type PostgresWalletDBHandler struct {
	Db *sql.DB;
	Config variablesHandler.DbConfig;
}

func (p *PostgresWalletDBHandler) CreateWallet(wallet model.WalletDTO) (int64, error) {
	var rowsAffected int64
	maxAttempts := 10
	attempt := 0
	walletNumber := generateWalletNumber()

	for {
		existentWallet, err := p.GetWalletByNumber(walletNumber)
		if err != nil {
			return rowsAffected, err
		}

		if existentWallet == nil {
			break // No existe una billetera con este número, se puede utilizar
		}

		attempt++
		if attempt >= maxAttempts {
			return rowsAffected, fmt.Errorf("se excedió el límite de intentos para generar un número único de billetera")
		}

		walletNumber = generateWalletNumber()
	}

	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		} else if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	response, err := transaction.Exec(`INSERT INTO public.wallet(
								customer_id, wallet_number, creation_date, balance)
								VALUES ($1, $2, $3, $4);`, wallet.CustomerId, walletNumber, wallet.CreationDate, wallet.Balance)

	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresWalletDBHandler) GetWalletByNumber(walletNumber string) (*model.Wallet, error) {
	row := p.Db.QueryRow(`SELECT 
						id
						, customer_id
						, wallet_number
						, creation_date
						, balance
						FROM public.wallet
						WHERE wallet_number=$1;`, walletNumber)

	var wallet model.Wallet

	err := row.Scan(&wallet.ID, &wallet.CustomerId, &wallet.WalletNumber, &wallet.CreationDate, &wallet.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &wallet, nil
}

func (p *PostgresWalletDBHandler) DeleteWallet(id int) (int64, error) {
	var rowsAffected int64

	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		} else if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	response, err := transaction.Exec(`DELETE FROM public.wallet
								WHERE id=$1;`, id)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func generateWalletNumber() string {
	rand.Seed(time.Now().UnixNano())

	minStr := "1000000000000000000000"
	maxStr := "9999999999999999999999"

	min := new(big.Int)
	min.SetString(minStr, 10)

	max := new(big.Int)
	max.SetString(maxStr, 10)

	randomNumber := new(big.Int).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), new(big.Int).Sub(max, min))
	randomNumber.Add(randomNumber, min)

	return fmt.Sprintf("%022s", randomNumber.String())
}
