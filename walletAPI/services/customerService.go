package services

import (
	"database/sql"

	"github.com/labora-wallet/walletAPI/model"
)

type PostgresCustomerDBHandler struct {
	Db *sql.DB
}

func (p *PostgresCustomerDBHandler) CreateCustomer(customer model.CustomerDTO) (int64, error) {
	var id int64

	existentCustomer, err := p.GetCustomerByIdentityNumber(customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)
	if err != nil {
		return id, err
	} else if existentCustomer != nil {
		return int64(existentCustomer.ID), nil
	}

	transaction, err := p.Db.Begin()
	if err != nil {
		return id, err
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

	row := transaction.QueryRow(`INSERT INTO public.customer(
						first_name, last_name, national_identity_number, national_identity_type, country_id)
						VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`, customer.FirstName, customer.LastName, customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)

	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresCustomerDBHandler) GetCustomerByIdentityNumber(nationalIdentityNumber, nationalIdentityType, countryId string) (*model.Customer, error) {
	row := p.Db.QueryRow(`SELECT
							id 
							, first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id
							FROM public.customer WHERE national_identity_number=$1 AND national_identity_type=$2 AND country_id=$3;`, nationalIdentityNumber, nationalIdentityType, countryId)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.NationalIdentityNumber, &customer.NationalIdentityType, &customer.CountryId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &customer, nil
}

func (p *PostgresCustomerDBHandler) GetCustomerById(id int64) (*model.Customer, error) {
	row := p.Db.QueryRow(`SELECT
							id 
							, first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id
							FROM public.customer WHERE id=$1;`, id)

	var customer model.Customer

	err := row.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.NationalIdentityNumber, &customer.NationalIdentityType, &customer.CountryId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &customer, nil
}

func (p *PostgresCustomerDBHandler) UpdateCustomer(dto model.CustomerDTO, id int64) (int64, error) {
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

	response, err := transaction.Exec(`UPDATE public.customer 
								first_name=$1
								, last_name=$2
								, national_identity_number=$3
								, national_identity_type=$4
								, country_id=$5
								WHERE id=$5;`, dto.FirstName, dto.LastName, dto.NationalIdentityNumber, dto.NationalIdentityType, dto.CountryId, id)

	if err != nil {
		return rowsAffected, err
	}

	rowsAffected, err = response.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (p *PostgresCustomerDBHandler) DeleteCustomer(id int64) (int64, error) {
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

	response, err := transaction.Exec(`DELETE FROM public.customer
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
