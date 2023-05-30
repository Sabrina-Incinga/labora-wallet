package services

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model"
)

type PostgresWalletMovementDBHandler struct {
	Db *sql.DB
}

func (p *PostgresWalletMovementDBHandler) CreateWalletMovement(movementData model.WalletMovementDTO, tx *sql.Tx) (int64, error){
	var rowsAffected int64
	var err error
	var response sql.Result

	query := createWalletMovementQuery()
	if tx != nil {
		response, err = tx.Exec(query, movementData.WalletId, movementData.MovementDate, movementData.MovementType, movementData.Amount)
	} else{
		response, err = p.Db.Exec(query, movementData.WalletId, movementData.MovementDate, movementData.MovementType, movementData.Amount)
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
	
func (p *PostgresWalletMovementDBHandler) GetWalletMovementByWalletId(walletId int64) ([]model.WalletMovement, error){
	var movements []model.WalletMovement = make([]model.WalletMovement, 0)
	rows, err := p.Db.Query(`SELECT id
								, wallet_id
								, movement_date
								, movement_type
								, amount
								FROM public.wallet_movement
								WHERE wallet_id = $1;`, walletId)

	if err != nil {
		return movements, err
	}

	defer rows.Close()

	for rows.Next(){
		var movement model.WalletMovement

		err = rows.Scan(&movement.ID, &movement.WalletId, &movement.MovementDate, &movement.MovementType, &movement.Amount)
		if err != nil {
			return movements, err
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func createWalletMovementQuery() string{
	return `INSERT INTO public.wallet_movement(
			wallet_id
			, movement_date
			, movement_type
			, amount)
			VALUES ($1, $2, $3, $4);`
}