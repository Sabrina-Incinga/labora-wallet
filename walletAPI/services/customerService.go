package services

import (
	"database/sql"
	"github.com/labora-wallet/walletAPI/model"
)

type PostgresCustomerDBHandler struct {
	Db *sql.DB
}

func (p *PostgresCustomerDBHandler) CreateCustomer(customer model.CustomerDTO, tx *sql.Tx) (int64, error) {
	var id int64
	var row *sql.Row

	existentCustomer, err := p.GetCustomerByIdentityNumber(customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)
	if err != nil {
		return id, err
	} else if existentCustomer != nil {
		return existentCustomer.ID, nil
	}

	query := createCustomerQuery()

	if tx != nil {
		row = tx.QueryRow(query, customer.FirstName, customer.LastName, customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)
	}else{
		row = p.Db.QueryRow(query, customer.FirstName, customer.LastName, customer.NationalIdentityNumber, customer.NationalIdentityType, customer.CountryId)
	}

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

func (p *PostgresCustomerDBHandler) UpdateCustomer(dto model.CustomerDTO, id int64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result
	query := updateCustomerQuery()

	if tx != nil {
		response, err = tx.Exec(query, dto.FirstName, dto.LastName, dto.NationalIdentityNumber, dto.NationalIdentityType, dto.CountryId, id)
	}else{
		response, err = p.Db.Exec(query, dto.FirstName, dto.LastName, dto.NationalIdentityNumber, dto.NationalIdentityType, dto.CountryId, id)
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

func (p *PostgresCustomerDBHandler) DeleteCustomer(id int64, tx *sql.Tx) (int64, error) {
	var rowsAffected int64
	var err error
	var response sql.Result
	query := deleteCustomerQuery()

	if tx != nil {
		response, err = tx.Exec(query, id)
	}else{
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

func createCustomerQuery()string{
	query := `INSERT INTO public.customer(
							first_name
							, last_name
							, national_identity_number
							, national_identity_type
							, country_id)
							VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	return query
}

func updateCustomerQuery() string {
	query := `UPDATE public.customer 
						first_name=$1
						, last_name=$2
						, national_identity_number=$3
						, national_identity_type=$4
						, country_id=$5
						WHERE id=$6;`

	return query
}

func deleteCustomerQuery() string {
	query := `DELETE FROM public.customer
						WHERE id=$1;`

	return query
}