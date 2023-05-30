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
	Db     *sql.DB
	Config variablesHandler.DbConfig
}

func (p *PostgresWalletDBHandler) CreateWallet(wallet model.WalletDTO, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result

	walletNumber, err := generateUniqueWalletNumber(p)
	if err != nil {
		return rowsAffected, err
	}

	query := createWalletQuery()

	if tx != nil {
		response, err = tx.Exec(query, wallet.CustomerId, walletNumber, wallet.CreationDate, wallet.Balance)
	} else {
		response, err = p.Db.Exec(query, wallet.CustomerId, walletNumber, wallet.CreationDate, wallet.Balance)
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

func (p *PostgresWalletDBHandler) GetWalletStatusById(id int64) (string, error) {
	row := p.Db.QueryRow(`SELECT wt.creation_status
							FROM public.wallet w
							INNER JOIN public.wallet_tracker wt
							ON w.customer_id = wt.customer_id
							WHERE w.id = $1
							ORDER BY wt.id DESC
							LIMIT 1;`, id)

	var walletStatus string

	err := row.Scan(&walletStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			return "", err
		}
	}

	return walletStatus, nil
}

func (p *PostgresWalletDBHandler) DeleteWallet(id int64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result
	query := deleteWalletQuery()

	if tx != nil {
		response, err = tx.Exec(query, id)
	} else {
		response, err = p.Db.Exec(query, id)
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

func (p *PostgresWalletDBHandler) GetConfig() variablesHandler.DbConfig {
	return p.Config
}

func createWalletQuery() string{
	return `INSERT INTO public.wallet(
		customer_id, wallet_number, creation_date, balance)
		VALUES ($1, $2, $3, $4);`
}

func deleteWalletQuery() string {
	return `DELETE FROM public.wallet
	WHERE id=$1;`
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

func generateUniqueWalletNumber(p *PostgresWalletDBHandler) (string, error) {
	maxAttempts := 10
	attempt := 0
	walletNumber := generateWalletNumber()

	for {
		existentWallet, err := p.GetWalletByNumber(walletNumber)
		if err != nil {
			return "", err
		}

		if existentWallet == nil {
			break
		}

		attempt++
		if attempt >= maxAttempts {
			return "", fmt.Errorf("se excedió el límite de intentos para generar un número único de billetera")
		}

		walletNumber = generateWalletNumber()
	}
	return walletNumber, nil
}
