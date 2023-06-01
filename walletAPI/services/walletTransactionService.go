package services

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/labora-wallet/walletAPI/model/dtos"
	"github.com/labora-wallet/walletAPI/services/interfaces"
)

type PostgresWalletTransactionDBHandler struct {
	Db                        *sql.DB
	WalletServiceImpl         interfaces.WalletDBHandler
	WalletMovementServiceImpl interfaces.WalletMovementDBHandler
	WalletTrackerServiceImpl  interfaces.WalletTrackerDBHandler
	Mutex                     *sync.Mutex
}

func (p *PostgresWalletTransactionDBHandler) Transfer(transactionData dtos.WalletTransactionDTO) (int64, error) {
	var rowsAffected, rowsAffected1, rowsAffected2 int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}
	trackerDto := dtos.InitializeWalletTracker()
	if transactionData.OriginWalletNumber == transactionData.DestinationWalletNumber {
		err = fmt.Errorf("No es posible realizar una transferencia a la misma cuenta")
		return rowsAffected, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()


	rowsAffected1, err = p.UpdateWalletBalance(transactionData.OriginWalletNumber, dtos.WITHDRAWMOVEMENT, transactionData.Amount, transaction)
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected2, err = p.UpdateWalletBalance(transactionData.DestinationWalletNumber, dtos.DEPOSITMOVEMENT, transactionData.Amount, transaction)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected = rowsAffected1 + rowsAffected2

	if rowsAffected < 2 {
		trackerDto.RequestStatus = dtos.FAILEDREQUEST
	}else{
		trackerDto.RequestStatus = dtos.SUCCESSFULREQUEST
	}		


	err = getWalletDataAndStoreTracks(p, transactionData, dtos.TRANSFERMOVEMENT, transaction, trackerDto)

	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresWalletTransactionDBHandler) Withdraw(transactionData dtos.WalletTransactionDTO) (int64, error) {
	return performWalletMovement(transactionData, p, dtos.WITHDRAWMOVEMENT)
}

func (p *PostgresWalletTransactionDBHandler) AddToAccount(transactionData dtos.WalletTransactionDTO) (int64, error) {
	return performWalletMovement(transactionData, p, dtos.DEPOSITMOVEMENT)
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
	if transactionType == dtos.WITHDRAWMOVEMENT {
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

func performWalletMovement(transactionData dtos.WalletTransactionDTO, p *PostgresWalletTransactionDBHandler, movementType string) (int64, error) {
	var rowsAffected int64
	transaction, err := p.Db.Begin()
	if err != nil {
		return rowsAffected, err
	}

	trackerDto := dtos.InitializeWalletTracker()
	transactionData.DestinationWalletNumber = ""
	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			err = transaction.Commit()
		}
	}()

	rowsAffected, err = p.UpdateWalletBalance(transactionData.OriginWalletNumber, movementType, transactionData.Amount, transaction)
	if err != nil {
		return rowsAffected, err
	}

	trackerDto.RequestStatus = dtos.SUCCESSFULREQUEST
	if rowsAffected == 0 {
		trackerDto.RequestStatus = dtos.FAILEDREQUEST
	}

	err = getWalletDataAndStoreTracks(p, transactionData, movementType, transaction, trackerDto)
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func getWalletDataAndStoreTracks(p *PostgresWalletTransactionDBHandler, transactionData dtos.WalletTransactionDTO, movementType string, transaction *sql.Tx, trackerDto dtos.WalletTrackerDTO) error {
	originWallet, err := p.WalletServiceImpl.GetWalletByNumber(transactionData.OriginWalletNumber)
	if err != nil {
		return err
	}

	destinationWallet, err := p.WalletServiceImpl.GetWalletByNumber(transactionData.DestinationWalletNumber)
	if err != nil {
		return err
	}

	trackerDto.CreationStatus = "NA"
	trackerDto.TrackType = dtos.WALLETMOVEMENT

	walletMovementDTO := dtos.InitializeWalletMovement()
	walletMovementDTO.ReceiverWalletId = new(int64)
	if originWallet != nil {
		trackerDto.CustomerId = originWallet.CustomerId
		walletMovementDTO.SenderWalletId = originWallet.ID
	}
	if destinationWallet != nil {
		*walletMovementDTO.ReceiverWalletId = destinationWallet.ID
	}
	walletMovementDTO.MovementType = movementType
	walletMovementDTO.Amount = transactionData.Amount

	_, err = p.WalletTrackerServiceImpl.CreateWalletTracker(trackerDto, nil)
	if err != nil {
		return err
	}

	_, err = p.WalletMovementServiceImpl.CreateWalletMovement(walletMovementDTO, transaction)
	if err != nil {
		return err
	}

	return nil
}

func updateBalanceQuery() string {
	query := `UPDATE public.wallet
						SET balance=balance+$1
						WHERE id=$2;`

	return query
}
