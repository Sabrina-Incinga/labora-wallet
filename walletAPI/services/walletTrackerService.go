package services

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model"
)

type PostgresWalletTrackerDBHandler struct {
	Db *sql.DB
}

func (p *PostgresWalletTrackerDBHandler) CreateWalletTracker(tracker model.WalletTrackerDTO) (int64, error) {
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

	response, err := transaction.Exec(`INSERT INTO public.wallet_tracker(
										customer_id, record_date, creation_status)
										VALUES ($1, $2, $3);`, tracker.CustomerId, tracker.RecordDate, tracker.CreationStatus)
	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresWalletTrackerDBHandler) GetWalletTrackByCustomerId(customerId int64) (*[]model.WalletTracker, error) {
	var trackers []model.WalletTracker = make([]model.WalletTracker, 0)
	rows, err := p.Db.Query(`SELECT 
						id
						, customer_id
						, record_date
						, creation_status
						FROM public.wallet_tracker
						WHERE customer_id=$1;`, customerId)
	if err != nil {
		return &trackers, err
	}
	defer rows.Close()

	for rows.Next() {
		var tracker model.WalletTracker

		err = rows.Scan(&tracker.ID, &tracker.CustomerId, &tracker.RecordDate, &tracker.CreationStatus)
		if err != nil {
			return &trackers, err
		}

		trackers = append(trackers, tracker)
	}

	return &trackers, nil
}
