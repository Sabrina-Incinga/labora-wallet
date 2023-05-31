package services

import (
	"database/sql"
	"fmt"

	"github.com/labora-wallet/walletAPI/model"
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type PostgresWalletMovementDBHandler struct {
	Db *sql.DB
}

func (p *PostgresWalletMovementDBHandler) CreateWalletMovement(movementData dtos.WalletMovementDTO, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result

	query := createWalletMovementQuery()

	var receiverWalletID sql.NullInt64
	if *movementData.ReceiverWalletId != 0 {
		receiverWalletID.Int64 = *movementData.ReceiverWalletId
		receiverWalletID.Valid = true
	} else {
		receiverWalletID.Valid = false
	}

	fmt.Println(receiverWalletID)
	fmt.Println(*movementData.ReceiverWalletId)

	if tx != nil {
		response, err = tx.Exec(query, movementData.SenderWalletId, receiverWalletID, movementData.MovementDate, movementData.MovementType, movementData.Amount)
	} else {
		response, err = p.Db.Exec(query, movementData.SenderWalletId, receiverWalletID, movementData.MovementDate, movementData.MovementType, movementData.Amount)
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

func (p *PostgresWalletMovementDBHandler) GetWalletMovementByWalletId(walletId int64) ([]model.WalletMovement, error) {
	var movements []model.WalletMovement = make([]model.WalletMovement, 0)
	rows, err := p.Db.Query(`SELECT id
								, sender_wallet_id
								, receiver_wallet_id
								, movement_date
								, movement_type
								, amount
								FROM public.wallet_movement
								WHERE wallet_id = $1;`, walletId)

	if err != nil {
		return movements, err
	}

	defer rows.Close()

	for rows.Next() {
		var movement model.WalletMovement

		err = rows.Scan(&movement.ID, &movement.SenderWalletId, &movement.ReceiverWalletId, &movement.MovementDate, &movement.MovementType, &movement.Amount)
		if err != nil {
			return movements, err
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func createWalletMovementQuery() string {
	return `INSERT INTO public.wallet_movement(
			sender_wallet_id
			, receiver_wallet_id
			, movement_date
			, movement_type
			, amount)
			VALUES ($1, $2, $3, $4, $5);`
}
